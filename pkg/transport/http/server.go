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
	"context"
	"crypto/tls"
	"github.com/go-ceres/ceres/internal/assert"
	"github.com/go-ceres/ceres/internal/bytesconv"
	"github.com/go-ceres/ceres/internal/endpoint"
	"github.com/go-ceres/ceres/internal/host"
	"github.com/go-ceres/ceres/internal/path"
	"github.com/go-ceres/ceres/pkg/common/errors"
	"github.com/go-ceres/ceres/pkg/transport"
	"github.com/valyala/fasthttp"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	KindHttp                 transport.Kind = "grpc"
	ModName                                 = "server.http"
	SupportPackageIsVersion1                = true
)

var (
	_ transport.Transport  = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

// Server 服务定义
type Server struct {
	baseContext context.Context // 应用上下文
	listener    net.Listener    // 网络监听器
	tlsHandler  TLSHandler
	initOnce    sync.Once
	ctxPool     sync.Pool // 上下文对象池
	endpoint    *url.URL
	server      *fasthttp.Server
	maxParams   uint16
	opts        *ServerOptions
	trees       Routers
	RouterGroup // 路由
}

// NewServer 新建
func NewServer(opts ...ServerOption) *Server {
	o := DefaultServerOptions()
	for _, opt := range opts {
		opt(o)
	}
	return NewWithOptions(o)
}

// NewWithOptions 根据全参数创建
func NewWithOptions(options *ServerOptions) *Server {
	srv := &Server{
		opts:        options,
		baseContext: context.Background(),
		trees:       make(Routers, 0, 9),
		RouterGroup: RouterGroup{
			handlers: nil,
			basePath: "/",
			root:     true,
		},
	}
	srv.RouterGroup.server = srv
	// 初始化服务
	srv.init()
	return srv
}

// Kind Transport类型
func (s *Server) Kind() transport.Kind {
	return KindHttp
}

// Endpoint 获取入口地址
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

// Start 启动
func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return err
	}
	s.baseContext = ctx
	s.opts.logger.Infof("[HTTP] server listening on: %s", s.listener.Addr().String())
	var err error
	listener := s.listener
	if s.opts.TlsConf != nil {
		err = s.server.ServeTLS(listener, s.opts.TlsConf.CertFile, s.opts.TlsConf.CertFile)
	} else {
		err = s.server.Serve(listener)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// init 初始化
func (s *Server) init() {
	s.initOnce.Do(func() {
		s.server = &fasthttp.Server{
			Logger:                        &disableLogger{},
			LogAllErrors:                  false,
			ErrorHandler:                  s.serverErrorHandler,
			Handler:                       s.handler,
			Name:                          s.opts.ServerHeader,
			Concurrency:                   s.opts.Concurrency,
			NoDefaultContentType:          s.opts.DisableDefaultContentType,
			DisableHeaderNamesNormalizing: s.opts.DisableHeaderNamesNormalizing,
			DisableKeepalive:              s.opts.DisableKeepalive,
			MaxRequestBodySize:            s.opts.MaxRequestBodySize,
			NoDefaultServerHeader:         s.opts.ServerHeader == "",
			ReadTimeout:                   s.opts.ReadTimeout,
			WriteTimeout:                  s.opts.ReadTimeout,
			IdleTimeout:                   s.opts.IdleTimeout,
			ReadBufferSize:                s.opts.ReadBufferSize,
			WriteBufferSize:               s.opts.WriteBufferSize,
			GetOnly:                       s.opts.GetOnly,
			ReduceMemoryUsage:             s.opts.ReduceMemoryUsage,
			StreamRequestBody:             s.opts.StreamRequestBody,
			DisablePreParseMultipartForm:  s.opts.DisablePreParseMultipartForm,
		}
	})
}

func (s *Server) serverErrorHandler(requestCtx *fasthttp.RequestCtx, err error) {
	c := s.acquireContext(requestCtx)
	if _, ok := err.(*fasthttp.ErrSmallBuffer); ok {
		err = errors.RequestHeaderFieldsTooLarge("REQUEST_HEADER_FIELDS_TOO_LARGE", statusMessage[StatusRequestHeaderFieldsTooLarge])
	} else if netErr, ok := err.(*net.OpError); ok && netErr.Timeout() {
		err = errors.RequestTimeout("REQUEST_TIMEOUT", statusMessage[StatusRequestTimeout])
	} else if err == fasthttp.ErrBodyTooLarge {
		err = errors.RequestEntityTooLarge("BODY_TOO_LARGE", statusMessage[StatusRequestEntityTooLarge])
	} else if err == fasthttp.ErrGetOnly {
		err = errors.MethodNotAllowed("METHOD_NOT_ALLOWED", statusMessage[StatusMethodNotAllowed])
	} else if strings.Contains(err.Error(), "timeout") {
		err = errors.RequestTimeout("REQUEST_TIMEOUT", statusMessage[StatusRequestTimeout])
	} else {
		err = errors.BadRequest("BAD_REQUEST", statusMessage[StatusBadRequest])
	}
	if catch := s.opts.ErrorHandler(c, err); catch != nil {
		_ = c.SendStatus(StatusInternalServerError)
	}

	s.releaseContext(c)
}

func (s *Server) next(ctx *Context) (err error) {
	rPath := ctx.Path()
	unescape := false
	if s.opts.UseRawPath {
		rPath = ctx.pathOriginal
		unescape = s.opts.UnescapePathValues
	}
	if s.opts.RemoveExtraSlash {
		rPath = path.CleanPath(rPath)
	}
	tree := s.trees[ctx.methodInt]
	paramsPointer := &ctx.params
	value := tree.find(rPath, paramsPointer, unescape)
	if value.handlers != nil {
		metadata := &Metadata{
			operation:    value.pathTemplate,
			pathTemplate: value.pathTemplate,
			request:      ctx.Request(),
			response:     ctx.Response(),
		}
		if s.endpoint != nil {
			metadata.endpoint = s.endpoint.String()
		}
		ctx.SetUserContext(transport.NewMetadataServerContext(s.baseContext, metadata))
		ctx.handlers = value.handlers
		ctx.pathTemplate = value.pathTemplate
		err = ctx.Next()
		return
	}
	if ctx.method != MethodConnect && rPath != "/" {
		if value.tsr && s.opts.RedirectTrailingSlash {
			redirectTrailingSlash(ctx)
			return
		}
		if s.opts.RedirectFixedPath && redirectFixedPath(ctx, tree.root, s.opts.RedirectFixedPath) {
			return
		}
	}
	err = errors.NotFound("NOT_FOUND", "Cannot "+ctx.method+" "+ctx.pathOriginal)
	return
}

// handler 入口
func (s *Server) handler(fastCtx *fasthttp.RequestCtx) {
	// 分配上下文
	ctx := s.acquireContext(fastCtx)
	// 释放ctx
	defer s.releaseContext(ctx)
	// 没有找到method
	if ctx.methodInt == -1 {
		_ = ctx.SetStatusCode(StatusBadRequest).SendString(default405Body)
		return
	}
	err := s.next(ctx)
	if err != nil {
		if catch := s.opts.ErrorHandler(ctx, err); catch != nil {
			_ = ctx.SendStatus(StatusInternalServerError)
		}
	}
}

// addRoute 添加路由
func (s *Server) addRoute(method, path string, handlers HandlersChain) {
	if len(path) == 0 {
		panic("path should not be ''")
	}
	assert.Panic(path[0] == '/', "path must begin with '/'")
	assert.Panic(method != "", "HTTP method can not be empty")
	assert.Panic(len(handlers) > 0, "there must be at least one handler")

	if !s.opts.DisablePrintRoute {
		//debugPrintRoute(method, path, handlers)
	}

	methodRouter := s.trees.get(method)
	if methodRouter == nil {
		methodRouter = &Router{method: method, root: &node{}}
		s.trees = append(s.trees, methodRouter)
	}
	methodRouter.addRoute(path, handlers)

	// Update maxParams
	if paramsCount := countParams(path); paramsCount > s.maxParams {
		s.maxParams = paramsCount
	}
}

// Stop 关闭
func (s *Server) Stop(ctx context.Context) error {
	s.opts.logger.Info("[HTTP] server stopping")
	return s.server.ShutdownWithContext(ctx)
}

// listenAndEndpoint 启动网络监听，并获取入口地址
func (s *Server) listenAndEndpoint() error {
	if s.listener == nil {
		lis, err := net.Listen(s.opts.Network, s.opts.Address)
		if err != nil {
			return err
		}
		s.listener = lis
	}
	if s.endpoint == nil {
		addr, err := host.Extract(s.opts.Address, s.listener)
		if err != nil {
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("http", s.opts.TlsConf != nil), addr)
	}
	return nil
}

func redirectTrailingSlash(c *Context) {
	p := bytesconv.BytesToString(c.fastCtx.Request.URI().Path())
	if prefix := path.CleanPath(bytesconv.BytesToString(c.fastCtx.Request.Header.Peek("X-Forwarded-Prefix"))); prefix != "." {
		p = prefix + "/" + p
	}

	tmpURI := trailingSlashURL(p)

	query := c.fastCtx.Request.URI().QueryString()

	if len(query) > 0 {
		tmpURI = tmpURI + "?" + bytesconv.BytesToString(query)
	}

	c.fastCtx.Request.SetRequestURI(tmpURI)
	redirectRequest(c)
}

func redirectFixedPath(c *Context, root *node, trailingSlash bool) bool {
	rPath := bytesconv.BytesToString(c.fastCtx.Request.URI().Path())
	if fixedPath, ok := root.findCaseInsensitivePath(path.CleanPath(rPath), trailingSlash); ok {
		c.fastCtx.Request.SetRequestURI(bytesconv.BytesToString(fixedPath))
		redirectRequest(c)
		return true
	}
	return false
}

// redirectRequest 重定向
func redirectRequest(c *Context) {
	code := StatusMovedPermanently // Permanent redirect, request with GET method
	if bytesconv.BytesToString(c.fastCtx.Request.Header.Method()) != MethodGet {
		code = StatusTemporaryRedirect
	}

	c.fastCtx.Redirect(bytesconv.BytesToString(c.fastCtx.Request.URI().RequestURI()), code)
}

func trailingSlashURL(ts string) string {
	tmpURI := ts + "/"
	if length := len(ts); length > 1 && ts[length-1] == '/' {
		tmpURI = ts[:length-1]
	}
	return tmpURI
}

type TLSHandler struct {
	clientHelloInfo *tls.ClientHelloInfo
}

// GetClientInfo Callback function to set CHI
func (t *TLSHandler) GetClientInfo(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	t.clientHelloInfo = info
	return nil, nil
}
