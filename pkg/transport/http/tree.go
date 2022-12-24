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
	"strings"
	"unicode"
)

type (
	kind uint8
	node struct {
		kind       kind
		label      byte
		prefix     string   // 前缀
		parent     *node    // 父节点
		children   children // 子节点
		ppath      string   //
		pnames     []string // 参数名称集合
		handlers   HandlersChain
		paramChild *node
		anyChild   *node
		isLeaf     bool
	}
	nodeValue struct {
		handlers     HandlersChain
		tsr          bool
		pathTemplate string
	}
	children []*node
)

// newNode 创建节点
func newNode(t kind, pre string, p *node, child children, mh HandlersChain, ppath string, pnames []string, paramChildren, anyChildren *node) *node {
	return &node{
		kind:       t,
		label:      pre[0],
		prefix:     pre,
		parent:     p,
		children:   child,
		ppath:      ppath,
		pnames:     pnames,
		handlers:   mh,
		paramChild: paramChildren,
		anyChild:   anyChildren,
		isLeaf:     child == nil && paramChildren == nil && anyChildren == nil,
	}
}

// findChild 查找子节点
func (n *node) findChild(l byte) *node {
	for _, c := range n.children {
		if c.label == l {
			return c
		}
	}
	return nil
}

// findChildWithLabel
func (n *node) findChildWithLabel(l byte) *node {
	for _, c := range n.children {
		if c.label == l {
			return c
		}
	}
	if l == paramLabel {
		return n.paramChild
	}
	if l == anyLabel {
		return n.anyChild
	}
	return nil
}

// findCaseInsensitivePath 根据路查找
func (n *node) findCaseInsensitivePath(path string, fixTrailingSlash bool) (ciPath []byte, found bool) {
	ciPath = make([]byte, 0, len(path)+1) // preallocate enough memory
	// Match paramKind.
	if n.label == paramLabel {
		end := 0
		for end < len(path) && path[end] != '/' {
			end++
		}
		ciPath = append(ciPath, path[:end]...)
		if end < len(path) {
			if len(n.children) > 0 {
				path = path[end:]

				goto loop
			}

			if fixTrailingSlash && len(path) == end+1 {
				return ciPath, true
			}
			return
		}

		if n.handlers != nil {
			return ciPath, true
		}
		if fixTrailingSlash && len(n.children) == 1 {
			// No handle found. Check if a handle for this path with(without) a trailing slash exists
			n = n.children[0]
			if n.prefix == "/" && n.handlers != nil {
				return append(ciPath, '/'), true
			}
		}
		return
	}

	// Match allKind.
	if n.label == anyLabel {
		return append(ciPath, path...), true
	}

	// Match static kind.
	if len(path) >= len(n.prefix) && strings.EqualFold(path[:len(n.prefix)], n.prefix) {
		path = path[len(n.prefix):]
		ciPath = append(ciPath, n.prefix...)

		if len(path) == 0 {
			if n.handlers != nil {
				return ciPath, true
			}
			// No handle found.
			// Try to fix the path by adding a trailing slash.
			if fixTrailingSlash {
				for i := 0; i < len(n.children); i++ {
					if n.children[i].label == '/' {
						n = n.children[i]
						if (len(n.prefix) == 1 && n.handlers != nil) ||
							(n.prefix == "*" && n.children[0].handlers != nil) {
							return append(ciPath, '/'), true
						}
						return
					}
				}
			}
			return
		}
	} else if fixTrailingSlash {
		// Nothing found.
		// Try to fix the path by adding / removing a trailing slash.
		if path == "/" {
			return ciPath, true
		}
		if len(path)+1 == len(n.prefix) && n.prefix[len(path)] == '/' &&
			strings.EqualFold(path, n.prefix[:len(path)]) &&
			n.handlers != nil {
			return append(ciPath, n.prefix...), true
		}
	}

loop:
	// First match static kind.
	for _, node := range n.children {
		if unicode.ToLower(rune(path[0])) == unicode.ToLower(rune(node.label)) {
			out, found := node.findCaseInsensitivePath(path, fixTrailingSlash)
			if found {
				return append(ciPath, out...), true
			}
		}
	}

	if n.paramChild != nil {
		out, found := n.paramChild.findCaseInsensitivePath(path, fixTrailingSlash)
		if found {
			return append(ciPath, out...), true
		}
	}

	if n.anyChild != nil {
		out, found := n.anyChild.findCaseInsensitivePath(path, fixTrailingSlash)
		if found {
			return append(ciPath, out...), true
		}
	}

	// Nothing found. We can recommend to redirect to the same URL
	// without a trailing slash if a leaf exists for that path
	found = fixTrailingSlash && path == "/" && n.handlers != nil
	return
}
