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
	"expvar"
	"fmt"
	strUtils "github.com/go-ceres/ceres/internal/strings"
	"github.com/valyala/fasthttp"
	"path"
	"strings"
	"time"
)

var _ IRouter = (*RouterGroup)(nil)

var (
	// Counter for total number of fs calls
	fsCalls = expvar.NewInt("fsCalls")

	// Counters for various response status codes
	fsOKResponses          = expvar.NewInt("fsOKResponses")
	fsNotModifiedResponses = expvar.NewInt("fsNotModifiedResponses")
	fsNotFoundResponses    = expvar.NewInt("fsNotFoundResponses")
	fsOtherResponses       = expvar.NewInt("fsOtherResponses")

	// Total size in bytes for OK response bodies served.
	fsResponseBodyBytes = expvar.NewInt("fsResponseBodyBytes")
)

// IRouter 路由接口定义
type IRouter interface {
	IRoutes
	Group(string, ...HandlerFunc) *RouterGroup
}

// IRoutes 路由接口定义
type IRoutes interface {
	Use(...HandlerFunc) IRoutes
	Handle(string, string, ...HandlerFunc) IRoutes
	Any(string, ...HandlerFunc) IRoutes
	GET(string, ...HandlerFunc) IRoutes
	POST(string, ...HandlerFunc) IRoutes
	DELETE(string, ...HandlerFunc) IRoutes
	PATCH(string, ...HandlerFunc) IRoutes
	PUT(string, ...HandlerFunc) IRoutes
	CONNECT(string, ...HandlerFunc) IRoutes
	OPTIONS(string, ...HandlerFunc) IRoutes
	TRACE(string, ...HandlerFunc) IRoutes
	HEAD(string, ...HandlerFunc) IRoutes
	StaticFile(string, string) IRoutes
	Static(string, string) IRoutes
	StaticFS(string, *fasthttp.FS) IRoutes
}

// RouterGroup 路由组定义
type RouterGroup struct {
	handlers HandlersChain // 方法集合
	basePath string        // 路径
	server   *Server       // 服务
	root     bool          // 是否是跟
}

// BasePath 返回基础路径
func (group *RouterGroup) BasePath() string {
	return group.basePath
}

// Use 使用中间件
func (group *RouterGroup) Use(handlerFunc ...HandlerFunc) IRoutes {
	group.handlers = append(group.handlers, handlerFunc...)
	return group.returnIRouter()
}

// Handle 注册方法
func (group *RouterGroup) Handle(httpMethod string, relativePath string, handlers ...HandlerFunc) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.mergeHandlers(handlers)
	group.server.addRoute(httpMethod, absolutePath, handlers)
	return group.returnIRouter()
}

func (group *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	group.GET(relativePath, handlers...)
	group.POST(relativePath, handlers...)
	group.DELETE(relativePath, handlers...)
	group.PATCH(relativePath, handlers...)
	group.PUT(relativePath, handlers...)
	group.CONNECT(relativePath, handlers...)
	group.OPTIONS(relativePath, handlers...)
	group.TRACE(relativePath, handlers...)
	group.HEAD(relativePath, handlers...)
	return group.returnIRouter()
}

func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.Handle(MethodGet, relativePath, handlers...)
}

func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.Handle(MethodPost, relativePath, handlers...)
}

func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.Handle(MethodDelete, relativePath, handlers...)
}

func (group *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.Handle(MethodPatch, relativePath, handlers...)
}

func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.Handle(MethodPut, relativePath, handlers...)
}

func (group *RouterGroup) CONNECT(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.Handle(MethodConnect, relativePath, handlers...)
}

func (group *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.Handle(MethodOptions, relativePath, handlers...)
}

func (group *RouterGroup) TRACE(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.Handle(MethodTrace, relativePath, handlers...)
}

func (group *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.Handle(MethodHead, relativePath, handlers...)
}

// Group 添加组
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		handlers: group.mergeHandlers(handlers),
		basePath: group.calculateAbsolutePath(relativePath),
		server:   group.server,
	}
}

// StaticFile 静态文件方法
func (group *RouterGroup) StaticFile(relativePath string, filepath string) IRoutes {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	handler := func(ctx *Context) error {
		return ctx.File(filepath)
	}
	group.GET(relativePath, handler)
	group.HEAD(relativePath, handler)
	return group.returnIRouter()
}

// Static 静态目录
func (group *RouterGroup) Static(relativePath string, path string) IRoutes {
	if path == "" {
		path = "."
	}
	if relativePath == "" {
		relativePath = "/"
	}
	// Prefix always start with a '/'
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	// in case-sensitive routing, all to lowercase
	if group.server.opts.CaseSensitive {
		relativePath = strUtils.ToLower(relativePath)
	}
	// Strip trailing slashes from the root path
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	// Is prefix a partial wildcard?
	if strings.Contains(relativePath, "*") {
		// /john* -> /john
		relativePath = strings.Split(relativePath, "*")[0]
		// Fix this later
	}
	prefixLen := len(relativePath)
	if prefixLen > 1 && relativePath[prefixLen-1:] == "/" {
		// /john/ -> /john
		prefixLen--
		relativePath = relativePath[:prefixLen]
	}
	const cacheDuration = 10 * time.Second
	fs := &fasthttp.FS{
		Root:                 path,
		AllowEmptyRoot:       false,
		GenerateIndexPages:   false,
		AcceptByteRange:      false,
		Compress:             false,
		CompressedFileSuffix: group.server.opts.CompressedFileSuffix,
		CacheDuration:        cacheDuration,
		IndexNames:           []string{"index.html"},
		PathRewrite: func(fctx *fasthttp.RequestCtx) []byte {
			path := fctx.Path()
			if len(path) >= prefixLen {
				path = path[prefixLen:]
				if len(path) == 0 || path[len(path)-1] != '/' {
					path = append(path, '/')
				}
			}
			if len(path) > 0 && path[0] != '/' {
				path = append([]byte("/"), path...)
			}
			return path
		},
		PathNotFound: func(fctx *fasthttp.RequestCtx) {
			fctx.Response.SetBodyString("not found")
			fctx.Response.SetStatusCode(StatusNotFound)
		},
	}
	fileHandler := fs.NewRequestHandler()
	handler := func(ctx *Context) error {
		// Serve file
		fileHandler(ctx.fastCtx)
		// Return request if found and not forbidden
		status := ctx.fastCtx.Response.StatusCode()
		if status != StatusNotFound && status != StatusForbidden {
			group.updateFSCounters(ctx.fastCtx)
			return nil
		}
		// Reset response to default
		content := fmt.Sprintf(`file "%s" not found`, ctx.Path())
		ctx.fastCtx.SetContentType("")
		ctx.fastCtx.Response.SetStatusCode(StatusOK)
		ctx.fastCtx.Response.SetBodyString(content)
		return ctx.Next()
	}
	group.GET(relativePath+"/*filepath", handler)
	group.HEAD(relativePath+"/*filepath", handler)
	return group.returnIRouter()
}

func (group *RouterGroup) updateFSCounters(ctx *fasthttp.RequestCtx) {
	// Increment the number of fsHandler calls.
	fsCalls.Add(1)

	// Update other stats counters
	resp := &ctx.Response
	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		fsOKResponses.Add(1)
		fsResponseBodyBytes.Add(int64(resp.Header.ContentLength()))
	case fasthttp.StatusNotModified:
		fsNotModifiedResponses.Add(1)
	case fasthttp.StatusNotFound:
		fsNotFoundResponses.Add(1)
	default:
		fsOtherResponses.Add(1)
	}
}

func (group *RouterGroup) StaticFS(relativePath string, fs *fasthttp.FS) IRoutes {
	if relativePath == "" {
		relativePath = "/"
	}
	// Prefix always start with a '/'
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	// in case-sensitive routing, all to lowercase
	if group.server.opts.CaseSensitive {
		relativePath = strUtils.ToLower(relativePath)
	}
	// Is prefix a partial wildcard?
	if strings.Contains(relativePath, "*") {
		// /john* -> /john
		relativePath = strings.Split(relativePath, "*")[0]
		// Fix this later
	}
	prefixLen := len(relativePath)
	if prefixLen > 1 && relativePath[prefixLen-1:] == "/" {
		// /john/ -> /john
		prefixLen--
		relativePath = relativePath[:prefixLen]
	}
	fileHandler := fs.NewRequestHandler()
	handler := func(ctx *Context) error {
		// Serve file
		fileHandler(ctx.fastCtx)
		// Return request if found and not forbidden
		status := ctx.fastCtx.Response.StatusCode()
		if status != StatusNotFound && status != StatusForbidden {
			return nil
		}
		// Reset response to default
		ctx.fastCtx.SetContentType("")
		ctx.fastCtx.Response.SetStatusCode(StatusOK)
		ctx.fastCtx.Response.SetBodyString("")
		return ctx.Next()
	}
	group.GET(relativePath+"/*filepath", handler)
	group.HEAD(relativePath+"/*filepath", handler)
	return group.returnIRouter()
}

// mergeHandlers 合并
func (group *RouterGroup) mergeHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.handlers) + len(handlers)
	if finalSize >= int(AbortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.handlers)
	copy(mergedHandlers[len(group.handlers):], handlers)
	return mergedHandlers
}

// returnIRouter 返回router
func (group *RouterGroup) returnIRouter() IRouter {
	if group.root {
		return group.server
	}
	return group
}

// calculateAbsolutePath 计算路径
func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.basePath, relativePath)
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}
