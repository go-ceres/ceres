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
	"github.com/valyala/fasthttp"
	"path"
	"strings"
)

var _ IRouter = (*RouterGroup)(nil)

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
	//TODO implement me
	panic("implement me")
}

func (group *RouterGroup) StaticFS(relativePath string, fs *fasthttp.FS) IRoutes {
	//TODO implement me
	panic("implement me")
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
