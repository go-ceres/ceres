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
	"github.com/go-ceres/ceres/internal/bytesconv"
	"github.com/go-ceres/ceres/internal/matcher"
	"github.com/go-ceres/ceres/pkg/transport"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"io"
	"mime/multipart"
	"net/http"
	"sync"
	"time"
)

var (
	contextPool = sync.Pool{
		New: func() interface{} {
			return new(Context)
		},
	}
)

// Context 定义上下文
type Context struct {
	index        int8
	params       Params
	fastCtx      *fasthttp.RequestCtx
	handlers     HandlersChain
	middleware   matcher.Matcher
	matched      bool
	method       string
	methodInt    int
	pathOriginal string
	server       *Server
	pathTemplate string
}

// acquireContext 借用
func (s *Server) acquireContext(fastCtx *fasthttp.RequestCtx) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.server = s
	ctx.fastCtx = fastCtx
	ctx.params = make(Params, s.maxParams)
	ctx.index = -1
	ctx.matched = false
	ctx.pathOriginal = bytesconv.BytesToString(fastCtx.URI().PathOriginal())
	ctx.method = bytesconv.BytesToString(fastCtx.Request.Header.Method())
	ctx.methodInt = s.trees.MethodInt(ctx.method)
	ctx.middleware = s.opts.middleware
	return ctx
}

// releaseContext 释放上下文
func (s *Server) releaseContext(c *Context) {
	c.params = c.params[0:0]
	c.handlers = nil
	c.fastCtx = nil
	c.server = nil
	c.pathTemplate = ""
	c.middleware = nil
	c.pathOriginal = ""
	contextPool.Put(c)
}

// URI 获取URI信息
func (ctx *Context) URI() *fasthttp.URI {
	return ctx.fastCtx.Request.URI()
}

// Path 获取路径
func (ctx *Context) Path() string {
	return bytesconv.BytesToString(ctx.URI().Path())
}

// Index 获取当前执行函数的下标
func (ctx *Context) Index() int8 {
	return ctx.index
}

// PathTemplate 获取路径模板,例如：/user/:id
func (ctx *Context) PathTemplate() string {
	return ctx.pathTemplate
}

// GetRequestHeader 获取请求体header
func (ctx *Context) GetRequestHeader(key string) string {
	return bytesconv.BytesToString(ctx.Request().Header.Peek(key))
}

// GetResponseHeader 获取响应header
func (ctx *Context) GetResponseHeader(key string) string {
	return bytesconv.BytesToString(ctx.Response().Header.Peek(key))
}

// AllRequestHeaders 获取所有的请求头
func (ctx *Context) AllRequestHeaders() map[string]string {
	headers := make(map[string]string)
	ctx.Request().Header.VisitAll(func(k, v []byte) {
		headers[bytesconv.BytesToString(k)] = bytesconv.BytesToString(v)
	})
	return headers
}

// AllResponseHeaders 获取所有的响应头
func (ctx *Context) AllResponseHeaders() map[string]string {
	headers := make(map[string]string)
	ctx.Response().Header.VisitAll(func(k, v []byte) {
		headers[bytesconv.BytesToString(k)] = bytesconv.BytesToString(v)
	})

	return headers
}

// RequestBodyStream 获取请求体流
func (ctx *Context) RequestBodyStream() io.Reader {
	return ctx.fastCtx.RequestBodyStream()
}

// MultipartForm 获取多文件
func (ctx *Context) MultipartForm() (*multipart.Form, error) {
	return ctx.fastCtx.MultipartForm()
}

// MultipartFormValue 获取指定键的multipartForm值
func (ctx *Context) MultipartFormValue(key string) (string, bool) {
	mf, err := ctx.MultipartForm()
	if err == nil && mf.Value != nil {
		vv := mf.Value[key]
		if len(vv) > 0 {
			return vv[0], true
		}
	}
	return "", false
}

// Middleware 执行grpc与http通用中间件
func (ctx *Context) Middleware(h transport.Handler) transport.Handler {
	tr, ok := transport.MetadataFromServerContext(ctx.UserContext())
	if ok {
		return transport.MiddlewareChain(ctx.middleware.Match(tr.Operation())...)(h)
	}
	return transport.MiddlewareChain(ctx.middleware.Match(ctx.Path())...)(h)
}

func (ctx *Context) SetUserContext(c context.Context) {
	ctx.fastCtx.SetUserValue(userContextKey{}, c)
}

// UserContext 获取上下文
func (ctx *Context) UserContext() context.Context {
	userCtx, ok := ctx.fastCtx.UserValue(userContextKey{}).(context.Context)
	if !ok {
		userCtx = context.Background()
		ctx.fastCtx.SetUserValue(userContextKey{}, userCtx)
	}
	return userCtx
}

// Operation grpc请求路径
func (ctx *Context) Operation() string {
	return ctx.pathTemplate
}

// GetFastCtx 获取fast的上下文
func (ctx *Context) GetFastCtx() *fasthttp.RequestCtx {
	return ctx.fastCtx
}

// Request 获取请求体
func (ctx *Context) Request() *fasthttp.Request {
	return &ctx.fastCtx.Request
}

// Response 返回响应体
func (ctx *Context) Response() *fasthttp.Response {
	return &ctx.fastCtx.Response
}

// SetStatusCode 设置状态码
func (ctx *Context) SetStatusCode(code int) *Context {
	ctx.Response().SetStatusCode(code)
	return ctx
}

// Redirect 重定向
func (ctx *Context) Redirect(statusCode int, uri []byte) error {
	ctx.redirect(uri, statusCode)
	return nil
}

// redirect 重定向
func (ctx *Context) redirect(uri []byte, statusCode int) {
	ctx.Response().Header.SetCanonical(bytesconv.StringToBytes(HeaderLocation), uri)
	statusCode = getRedirectStatusCode(statusCode)
	ctx.Response().SetStatusCode(statusCode)
}

// getRedirectStatusCode 获取重定向状态码
func getRedirectStatusCode(statusCode int) int {
	if statusCode == StatusMovedPermanently || statusCode == StatusFound ||
		statusCode == StatusSeeOther || statusCode == StatusTemporaryRedirect ||
		statusCode == StatusPermanentRedirect {
		return statusCode
	}
	return StatusFound
}

// SetResponseHeader 设置响应体header
func (ctx *Context) SetResponseHeader(key, value string) {
	if value == "" {
		ctx.Response().Header.Del(key)
		return
	}
	ctx.Response().Header.Set(key, value)
}

// SetRequestHeader 设置请求体header
func (ctx *Context) SetRequestHeader(key, value string) {
	if value == "" {
		ctx.Request().Header.Del(key)
		return
	}
	ctx.Request().Header.Set(key, value)
}

// GetParam 获取路径参数
func (ctx *Context) GetParam(key string) string {
	return ctx.params.MustGet(key)
}

// AllParam 返回所有的路径参数
func (ctx *Context) AllParam() Params {
	return ctx.params
}

// GetContextType 返回响应体文本类型
func (ctx *Context) GetContextType() string {
	return bytesconv.BytesToString(ctx.Response().Header.ContentType())
}

// GetCookie 获取cookie
func (ctx *Context) GetCookie(key string) string {
	return bytesconv.BytesToString(ctx.Request().Header.Cookie(key))
}

// SetCookie 设置cookie
func (ctx *Context) SetCookie(name, value string, opts ...CookieOption) *Context {
	ck := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(ck)
	for _, opt := range opts {
		opt(ck)
	}
	ck.SetKey(name)
	ck.SetValue(value)
	ctx.Response().Header.SetCookie(ck)
	return ctx
}

// ShouldBind 绑定数据
func (ctx *Context) ShouldBind(out interface{}) error {
	req := warpRequest(ctx.Request())
	defer releaseBindRequest(req)
	return defaultBinding.Unmarshal(out, req, &ctx.params)
}

// Result 返回响应结果
func (ctx *Context) Result(code int, v interface{}) error {
	ctx.SetStatusCode(code)
	return ctx.server.opts.EncodeResponseFunc(ctx, v)
}

// SendStatus 发送状态码
func (ctx *Context) SendStatus(code int) error {
	ctx.SetStatusCode(code)
	if len(ctx.fastCtx.Response.Body()) == 0 {
		return ctx.SendString(statusMessage[code])
	}
	return nil
}

// SendString 发送字符串
func (ctx *Context) SendString(body string) error {
	ctx.fastCtx.SetBodyString(body)
	return nil
}

// Write 添加body
func (ctx *Context) Write(p []byte) (int, error) {
	ctx.fastCtx.Response.AppendBody(p)
	return len(p), nil
}

// IsPost 是否是post请求
func (ctx *Context) IsPost() bool {
	return ctx.fastCtx.IsPost()
}

// IsGet 是否是get请求
func (ctx *Context) IsGet() bool {
	return ctx.fastCtx.IsGet()
}

// IsHead 是否是head请求
func (ctx *Context) IsHead() bool {
	return ctx.fastCtx.IsHead()
}

// IsPut 是否是put请求
func (ctx *Context) IsPut() bool {
	return ctx.fastCtx.IsPut()
}

// IsPatch 是否是patch请求
func (ctx *Context) IsPatch() bool {
	return ctx.fastCtx.IsPatch()
}

// IsOptions 是否是Options请求
func (ctx *Context) IsOptions() bool {
	return ctx.fastCtx.IsOptions()
}

// IsConnect 是否是Connect请求
func (ctx *Context) IsConnect() bool {
	return ctx.fastCtx.IsConnect()
}

// IsDelete 是否是Delete请求
func (ctx *Context) IsDelete() bool {
	return ctx.fastCtx.IsDelete()
}

// IsTrace 是否是Trace请求
func (ctx *Context) IsTrace() bool {
	return ctx.fastCtx.IsTrace()
}

// IsTLS 是否是tls链接
func (ctx *Context) IsTLS() bool {
	return ctx.fastCtx.IsTLS()
}

// Next 执行堆栈中与当前匹配路由的下一个方法
func (ctx *Context) Next() (err error) {
	ctx.index++
	if ctx.index < int8(len(ctx.handlers)) {
		err = ctx.handlers[ctx.index](ctx)
	}
	return err
}

// File 发送文件
func (ctx *Context) File(path string) error {
	ctx.fastCtx.SendFile(path)
	ctx.SetStatusCode(http.StatusOK)
	return nil
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	if ctx.fastCtx == nil {
		return time.Time{}, false
	}
	return ctx.fastCtx.Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	if ctx.fastCtx == nil {
		return nil
	}
	return ctx.fastCtx.Done()
}

func (ctx *Context) Err() error {
	if ctx.fastCtx == nil {
		return context.Canceled
	}
	return ctx.fastCtx.Err()
}

// WithValue 设置值
func (ctx *Context) WithValue(key, value any) {
	parant := ctx.UserContext()
	child := context.WithValue(parant, key, value)
	ctx.SetUserContext(child)
}

func (ctx *Context) Value(key any) any {
	if ctx.fastCtx == nil {
		return nil
	}
	return ctx.UserContext().Value(key)
}
