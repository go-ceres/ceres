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

package fiber

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-ceres/ceres/errors"
	"github.com/go-ceres/ceres/internal/endpoint"
	"github.com/go-ceres/ceres/internal/host"
	"github.com/go-ceres/ceres/middleware"
	"github.com/go-ceres/ceres/transport"
	tHttp "github.com/go-ceres/ceres/transport/http"
	"github.com/gofiber/fiber/v2"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var _ tHttp.Server = (*Server)(nil)

type Server struct {
	endpoint *url.URL     // 服务地址
	App      *fiber.App   // fiber
	ctxPool  sync.Pool    // 上下文对象池
	listener net.Listener // 网络监听
	config   *Config      // 配置信息
}

func New(config ...*Config) *Server {
	conf := DefaultConfig()
	if len(config) > 0 {
		conf = config[0]
	}
	srv := &Server{
		App: fiber.New(fiber.Config{
			DisableStartupMessage: true,
		}),
		config: conf,
	}
	srv.ctxPool.New = func() any {
		return &wrapper{
			srv: srv,
			ctx: nil,
		}
	}
	//srv.App.Use(srv.filter)
	return srv
}

// GetFiber 返回fiber
func (s *Server) GetFiber() *fiber.App {
	return s.App
}

// Use 添加中间件
func (s *Server) Use(selector string, m ...middleware.Middleware) {
	s.config.middleware.Add(selector, m...)
}

// filter 设置过滤
func (s *Server) filter() FilterFunc {
	return func(next fiber.Handler) fiber.Handler {
		return func(fiberCtx *fiber.Ctx) error {
			var (
				ctx    context.Context
				cancel context.CancelFunc
			)
			if s.config.Timeout > 0 {
				ctx, cancel = context.WithTimeout(fiberCtx.Context(), s.config.Timeout)
			} else {
				ctx, cancel = context.WithCancel(fiberCtx.Context())
			}
			defer cancel()
			pathTemplate := fiberCtx.Path()
			if route := s.getCurrentRouter(fiberCtx); route != nil {
				pathTemplate = route.Path
			}

			tr := &Transport{
				context:      fiberCtx,
				operation:    pathTemplate,
				pathTemplate: pathTemplate,
				requestHeader: &RequestHeaderCarrier{
					ctx: fiberCtx,
				},
				replyHeader: &ReplyHeaderCarrier{
					ctx: fiberCtx,
				},
			}
			if s.endpoint != nil {
				tr.endpoint = s.endpoint.String()
			}
			fiberCtx.SetUserContext(transport.NewServerContext(ctx, tr))
			return next(fiberCtx)
		}
	}
}

// getCurrentRouter 匹配当前路由
func (s *Server) getCurrentRouter(ctx *fiber.Ctx) *fiber.Route {
	params := ctx.AllParams()
	currenParams := make(map[string]interface{})
	for _, key := range ctx.Route().Params {
		currenParams[key] = params[key]
	}

	currentUrl, _ := ctx.GetRouteURL(ctx.Route().Name, currenParams)
	if currentUrl == ctx.Path() {
		return ctx.Route()
	}
	return nil
}

func (s *Server) Handler(method, relativePath string, h tHttp.HandleFunc, filters ...interface{}) {
	name := strings.Join(strings.Split(strings.TrimPrefix(relativePath, "/"), "/"), "_")
	var next = func(ctx *fiber.Ctx) error {
		// 设置内部上下文
		iCtx := s.ctxPool.Get().(*wrapper)
		iCtx.ctx = ctx
		// 处理请求函数
		if err := h(iCtx); err != nil {
			s.config.ene(ctx, err)
		}
		iCtx.ctx = nil
		s.ctxPool.Put(iCtx)
		return nil
	}
	next = FilterChain(filters...)(next)
	next = FilterChain(s.filter())(next)
	s.App.Add(method, relativePath, next).Name(name)
}

func (s *Server) GET(path string, h tHttp.HandleFunc, filters ...interface{}) {
	s.Handler(fiber.MethodGet, path, h, filters...)
}

func (s *Server) POST(path string, h tHttp.HandleFunc, filters ...interface{}) {
	s.Handler(http.MethodPost, path, h, filters...)
}

func (s *Server) HEAD(path string, h tHttp.HandleFunc, filters ...interface{}) {
	s.Handler(http.MethodHead, path, h, filters...)
}

func (s *Server) PUT(path string, h tHttp.HandleFunc, filters ...interface{}) {
	s.Handler(http.MethodPut, path, h, filters...)
}

func (s *Server) PATCH(path string, h tHttp.HandleFunc, filters ...interface{}) {
	s.Handler(http.MethodPatch, path, h, filters...)
}

func (s *Server) DELETE(path string, h tHttp.HandleFunc, filters ...interface{}) {
	s.Handler(http.MethodDelete, path, h, filters...)
}

func (s *Server) CONNECT(path string, h tHttp.HandleFunc, filters ...interface{}) {
	s.Handler(http.MethodConnect, path, h, filters...)
}

func (s *Server) OPTIONS(path string, h tHttp.HandleFunc, filters ...interface{}) {
	s.Handler(http.MethodOptions, path, h, filters...)
}

func (s *Server) TRACE(path string, h tHttp.HandleFunc, filters ...interface{}) {
	s.Handler(http.MethodTrace, path, h, filters...)
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return err
	}
	s.config.logger.Infof("[HTTP] server listening on: %s", s.listener.Addr().String())
	listener := s.listener
	if s.config.TlsConf != nil {
		// Set TLS config with handler
		cert, err := tls.LoadX509KeyPair(s.config.TlsConf.CertFile, s.config.TlsConf.KeyFile)
		if err != nil {
			return fmt.Errorf("tls: cannot load TLS key pair from certFile=%q and keyFile=%q: %s", s.config.TlsConf.CertFile, s.config.TlsConf.KeyFile, err)
		}

		tlsHandler := &fiber.TLSHandler{}
		config := &tls.Config{
			MinVersion: tls.VersionTLS12,
			Certificates: []tls.Certificate{
				cert,
			},
			GetCertificate: tlsHandler.GetClientInfo,
		}
		// Setup listener
		listener = tls.NewListener(listener, config)
		if err != nil {
			return err
		}
		s.App.SetTLSHandler(tlsHandler)
	}
	err := s.App.Listener(listener)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) GracefulStop(ctx context.Context) error {
	s.config.logger.Info("server stopping")
	return s.GracefulStop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	s.config.logger.Info("server stopping")
	return s.Stop(ctx)
}

// listenAndEndpoint 启动网络监听，并获取入口地址
func (s *Server) listenAndEndpoint() error {
	if s.listener == nil {
		lis, err := net.Listen(s.config.Network, s.config.Address)
		if err != nil {
			return err
		}
		s.listener = lis
	}
	if s.endpoint == nil {
		addr, err := host.Extract(s.config.Address, s.listener)
		if err != nil {
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("http", s.config.TlsConf != nil), addr)
	}
	return nil
}

func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}
