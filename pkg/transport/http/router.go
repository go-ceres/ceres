// Copyright 2022. ceres
// Author https://github.com/go-ceres/ceres
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"bytes"
	"fmt"
	"github.com/go-ceres/ceres/internal/bytesconv"
	"net/url"
	"strings"
)

const (
	staticKind kind = iota
	paramKind
	allKind
	paramLabel = byte(':')
	anyLabel   = byte('*')
	slash      = "/"
	nilString  = ""
)

var (
	strColon = []byte(":")
	strStar  = []byte("*")
)

// Router 方法路由管理器
type Router struct {
	method string
	root   *node
}

// Routers 路由管理器
type Routers []*Router

func (routers Routers) get(method string) *Router {
	for _, router := range routers {
		if router.method == method {
			return router
		}
	}
	return nil
}

func (routers Routers) MethodInt(method string) int {
	for i, router := range routers {
		if router.method == method {
			return i
		}
	}
	return -1
}

// countParams 获取路径参数数量
func countParams(path string) uint16 {
	var n uint16
	s := bytesconv.StringToBytes(path)
	n += uint16(bytes.Count(s, strColon))
	n += uint16(bytes.Count(s, strStar))
	return n
}

// checkPathValid 检查路由合法
func checkPathValid(path string) {
	if path == nilString {
		panic("empty path")
	}
	if path[0] != '/' {
		panic("path must begin with '/'")
	}
	for i, c := range []byte(path) {
		switch c {
		case ':':
			if (i < len(path)-1 && path[i+1] == '/') || i == (len(path)-1) {
				panic("wildcards must be named with a non-empty name in path '" + path + "'")
			}
			i++
			for ; i < len(path) && path[i] != '/'; i++ {
				if path[i] == ':' || path[i] == '*' {
					panic("only one wildcard per path segment is allowed, find multi in path '" + path + "'")
				}
			}
		case '*':
			if i == len(path)-1 {
				panic("wildcards must be named with a non-empty name in path '" + path + "'")
			}
			if i > 0 && path[i-1] != '/' {
				panic(" no / before wildcards in path " + path)
			}
			for ; i < len(path); i++ {
				if path[i] == '/' {
					panic("catch-all routes are only allowed at the end of the path in path '" + path + "'")
				}
			}
		}
	}
}

// addRoute 添加路径
func (r *Router) addRoute(path string, handlers HandlersChain) {
	checkPathValid(path)

	var (
		pnames []string // Param names
		ppath  = path   // Pristine path
	)

	if handlers == nil {
		panic(fmt.Sprintf("Adding route without handler function: %v", path))
	}

	// Add the front static route part of a non-static route
	for i, lcpIndex := 0, len(path); i < lcpIndex; i++ {
		// param route
		if path[i] == paramLabel {
			j := i + 1

			r.insert(path[:i], nil, staticKind, nilString, nil)
			for ; i < lcpIndex && path[i] != '/'; i++ {
			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, lcpIndex = j, len(path)

			if i == lcpIndex {
				// path node is last fragment of route path. ie. `/users/:id`
				r.insert(path[:i], handlers, paramKind, ppath, pnames)
				return
			} else {
				r.insert(path[:i], nil, paramKind, nilString, pnames)
			}
		} else if path[i] == anyLabel {
			r.insert(path[:i], nil, staticKind, nilString, nil)
			pnames = append(pnames, path[i+1:])
			r.insert(path[:i+1], handlers, allKind, ppath, pnames)
			return
		}
	}

	r.insert(path, handlers, staticKind, ppath, pnames)
}

// insert 插入路由
func (r *Router) insert(path string, handlers HandlersChain, t kind, ppath string, pnames []string) {
	currentNode := r.root
	if currentNode == nil {
		panic("hertz: invalid node")
	}
	search := path

	for {
		searchLen := len(search)
		prefixLen := len(currentNode.prefix)
		lcpLen := 0

		max := prefixLen
		if searchLen < max {
			max = searchLen
		}
		for ; lcpLen < max && search[lcpLen] == currentNode.prefix[lcpLen]; lcpLen++ {
		}

		if lcpLen == 0 {
			// At root node
			currentNode.label = search[0]
			currentNode.prefix = search
			if handlers != nil {
				currentNode.kind = t
				currentNode.handlers = handlers
				currentNode.ppath = ppath
				currentNode.pnames = pnames
			}
			currentNode.isLeaf = currentNode.children == nil && currentNode.paramChild == nil && currentNode.anyChild == nil
		} else if lcpLen < prefixLen {
			// Split node
			n := newNode(
				currentNode.kind,
				currentNode.prefix[lcpLen:],
				currentNode,
				currentNode.children,
				currentNode.handlers,
				currentNode.ppath,
				currentNode.pnames,
				currentNode.paramChild,
				currentNode.anyChild,
			)
			// Update parent path for all children to new node
			for _, child := range currentNode.children {
				child.parent = n
			}
			if currentNode.paramChild != nil {
				currentNode.paramChild.parent = n
			}
			if currentNode.anyChild != nil {
				currentNode.anyChild.parent = n
			}

			// Reset parent node
			currentNode.kind = staticKind
			currentNode.label = currentNode.prefix[0]
			currentNode.prefix = currentNode.prefix[:lcpLen]
			currentNode.children = nil
			currentNode.handlers = nil
			currentNode.ppath = nilString
			currentNode.pnames = nil
			currentNode.paramChild = nil
			currentNode.anyChild = nil
			currentNode.isLeaf = false

			// Only Static children could reach here
			currentNode.children = append(currentNode.children, n)

			if lcpLen == searchLen {
				// At parent node
				currentNode.kind = t
				currentNode.handlers = handlers
				currentNode.ppath = ppath
				currentNode.pnames = pnames
			} else {
				// Create child node
				n = newNode(t, search[lcpLen:], currentNode, nil, handlers, ppath, pnames, nil, nil)
				// Only Static children could reach here
				currentNode.children = append(currentNode.children, n)
			}
			currentNode.isLeaf = currentNode.children == nil && currentNode.paramChild == nil && currentNode.anyChild == nil
		} else if lcpLen < searchLen {
			search = search[lcpLen:]
			c := currentNode.findChildWithLabel(search[0])
			if c != nil {
				// Go deeper
				currentNode = c
				continue
			}
			// Create child node
			n := newNode(t, search, currentNode, nil, handlers, ppath, pnames, nil, nil)
			switch t {
			case staticKind:
				currentNode.children = append(currentNode.children, n)
			case paramKind:
				currentNode.paramChild = n
			case allKind:
				currentNode.anyChild = n
			}
			currentNode.isLeaf = currentNode.children == nil && currentNode.paramChild == nil && currentNode.anyChild == nil
		} else {
			// Node already exists
			if currentNode.handlers != nil && handlers != nil {
				panic("handlers are already registered for path '" + ppath + "'")
			}

			if handlers != nil {
				currentNode.handlers = handlers
				currentNode.ppath = ppath
				if len(currentNode.pnames) == 0 {
					currentNode.pnames = pnames
				}
			}
		}
		return
	}
}

// find 查找路由并将解析后的信息放入上下文
func (r *Router) find(path string, paramsPointer *Params, unescape bool) (res nodeValue) {
	var (
		cn          = r.root // current node
		search      = path   // current path
		searchIndex = 0
		buf         []byte
		paramIndex  int
	)

	backtrackToNextNodeKind := func(fromKind kind) (nextNodeKind kind, valid bool) {
		previous := cn
		cn = previous.parent
		valid = cn != nil

		// Next node type by priority
		if previous.kind == allKind {
			nextNodeKind = staticKind
		} else {
			nextNodeKind = previous.kind + 1
		}

		if fromKind == staticKind {
			// when backtracking is done from static kind block we did not change search so nothing to restore
			return
		}

		// restore search to value it was before we move to current node we are backtracking from.
		if previous.kind == staticKind {
			searchIndex -= len(previous.prefix)
		} else {
			paramIndex--
			// for param/any node.prefix value is always `:` so we can not deduce searchIndex from that and must use pValue
			// for that index as it would also contain part of path we cut off before moving into node we are backtracking from
			searchIndex -= len((*paramsPointer)[paramIndex].Value)
			(*paramsPointer) = (*paramsPointer)[:paramIndex]
		}
		search = path[searchIndex:]
		return
	}

	// search order: static > param > any
	for {
		if cn.kind == staticKind {
			if len(search) >= len(cn.prefix) && cn.prefix == search[:len(cn.prefix)] {
				// Continue search
				search = search[len(cn.prefix):]
				searchIndex = searchIndex + len(cn.prefix)
			} else {
				// not equal
				if (len(cn.prefix) == len(search)+1) &&
					(cn.prefix[len(search)]) == '/' && cn.prefix[:len(search)] == search && (cn.handlers != nil || cn.anyChild != nil) {
					res.tsr = true
				}
				// No matching prefix, let's backtrack to the first possible alternative node of the decision path
				nk, ok := backtrackToNextNodeKind(staticKind)
				if !ok {
					return // No other possibilities on the decision path
				} else if nk == paramKind {
					goto Param
				} else {
					// Not found (this should never be possible for static node we are looking currently)
					break
				}
			}
		}
		if search == nilString && len(cn.handlers) != 0 {
			res.handlers = cn.handlers
			break
		}

		// Static node
		if search != nilString {
			// If it can execute that logic, there is handler registered on the current node and search is `/`.
			if search == "/" && cn.handlers != nil {
				res.tsr = true
			}
			if child := cn.findChild(search[0]); child != nil {
				cn = child
				continue
			}
		}

		if search == nilString {
			if cd := cn.findChild('/'); cd != nil && (cd.handlers != nil || cd.anyChild != nil) {
				res.tsr = true
			}
		}

	Param:
		// Param node
		if child := cn.paramChild; search != nilString && child != nil {
			cn = child
			i := strings.Index(search, slash)
			if i == -1 {
				i = len(search)
			}
			(*paramsPointer) = (*paramsPointer)[:(paramIndex + 1)]
			val := search[:i]
			if unescape {
				if v, err := url.QueryUnescape(search[:i]); err == nil {
					val = v
				}
			}
			(*paramsPointer)[paramIndex].Value = val
			paramIndex++
			search = search[i:]
			searchIndex = searchIndex + i
			if search == nilString {
				if cd := cn.findChild('/'); cd != nil && (cd.handlers != nil || cd.anyChild != nil) {
					res.tsr = true
				}
			}
			continue
		}
	Any:
		// Any node
		if child := cn.anyChild; child != nil {
			// If any node is found, use remaining path for paramValues
			cn = child
			(*paramsPointer) = (*paramsPointer)[:(paramIndex + 1)]
			index := len(cn.pnames) - 1
			val := search
			if unescape {
				if v, err := url.QueryUnescape(search); err == nil {
					val = v
				}
			}

			(*paramsPointer)[index].Value = bytesconv.BytesToString(append(buf, val...))
			// update indexes/search in case we need to backtrack when no handler match is found
			paramIndex++
			searchIndex += len(search)
			search = nilString
			res.handlers = cn.handlers
			break
		}

		// Let's backtrack to the first possible alternative node of the decision path
		nk, ok := backtrackToNextNodeKind(allKind)
		if !ok {
			break // No other possibilities on the decision path
		} else if nk == paramKind {
			goto Param
		} else if nk == allKind {
			goto Any
		} else {
			// Not found
			break
		}
	}

	if cn != nil {
		res.pathTemplate = cn.ppath
		for i, name := range cn.pnames {
			(*paramsPointer)[i].Key = name
		}
	}

	return
}
