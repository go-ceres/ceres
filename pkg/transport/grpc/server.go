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

package grpc

import (
	"context"
	"github.com/go-ceres/ceres/internal/endpoint"
	"github.com/go-ceres/ceres/internal/host"
	"github.com/go-ceres/ceres/pkg/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"net/url"
)

const (
	KindGrpc transport.Kind = "grpc"
)

var (
	_ transport.Transport = (*Server)(nil)
)

// Server Grpc服务
type Server struct {
	*grpc.Server
	baseContext context.Context // 应用上下文
	endpoint    *url.URL        // 入口地址
	opts        *ServerOptions  // 配置信息
	listener    net.Listener    // 服务监听器
	health      *health.Server  // 健康服务
}

// NewServer 新建
func NewServer(opts ...ServerOption) *Server {
	options := DefaultServerOptions()
	for _, opt := range opts {
		opt(options)
	}
	return NewServerWithOptions(options)
}

// NewServerWithOptions 创建根据参数信息
func NewServerWithOptions(options *ServerOptions) *Server {
	srv := &Server{
		baseContext: context.Background(),
		opts:        options,
		health:      health.NewServer(),
	}
	unaryInts := []grpc.UnaryServerInterceptor{
		srv.unaryServerInterceptor(),
	}
	streamInts := []grpc.StreamServerInterceptor{
		srv.streamServerInterceptor(),
	}
	if len(srv.opts.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.opts.unaryInts...)
	}
	if len(srv.opts.streamInts) > 0 {
		streamInts = append(streamInts, srv.opts.streamInts...)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
		grpc.ChainStreamInterceptor(streamInts...),
	}
	if srv.opts.TlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.opts.TlsConf)))
	}
	if len(srv.opts.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.opts.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	if srv.opts.Health {
		grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	}
	if srv.opts.Reflection {
		reflection.Register(srv.Server)
	}
	return srv
}

// Kind 传输协议类型
func (s *Server) Kind() transport.Kind {
	return KindGrpc
}

// Endpoint 获取入口地址
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

// Start 启动服务
func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return err
	}
	s.baseContext = ctx
	s.opts.logger.Infof("[GRPC] server listening on: %s", s.listener.Addr().String())
	if s.opts.Health {
		s.health.Resume()
	}
	return s.Serve(s.listener)
}

// Stop 关闭
func (s *Server) Stop(ctx context.Context) error {
	s.opts.logger.Info("[GRPC] server stopping")
	s.GracefulStop()
	return nil
}

// listenAndEndpoint 启动grpc服务并且获取运行的ip地址
func (s *Server) listenAndEndpoint() error {
	if s.listener == nil {
		listener, err := net.Listen(s.opts.Network, s.opts.Address)
		if err != nil {
			return err
		}
		s.listener = listener
	}
	if s.endpoint == nil {
		addr, err := host.Extract(s.opts.Address, s.listener)
		if err != nil {
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("grpc", s.opts.TlsConf != nil), addr)
	}
	return nil
}
