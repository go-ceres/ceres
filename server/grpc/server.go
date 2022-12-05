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
	"github.com/go-ceres/ceres/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"net/url"
)

var _ server.Server = (*Server)(nil)

// Server Grpc服务
type Server struct {
	*grpc.Server
	endpoint *url.URL
	config   *Config         // 配置信息
	baseCtx  context.Context // 原始的context
	listener net.Listener    // 服务监听器
	health   *health.Server  // 健康服务
}

// New 新建
func New(conf ...*Config) *Server {
	var c = DefaultConfig()
	if len(conf) > 0 {
		c = conf[0]
	}
	srv := &Server{
		baseCtx: context.Background(),
		config:  c,
		health:  health.NewServer(),
	}
	unaryInts := []grpc.UnaryServerInterceptor{
		srv.unaryServerInterceptor(),
	}
	streamInts := []grpc.StreamServerInterceptor{
		srv.streamServerInterceptor(),
	}
	if len(srv.config.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.config.unaryInts...)
	}
	if len(srv.config.streamInts) > 0 {
		streamInts = append(streamInts, srv.config.streamInts...)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
		grpc.ChainStreamInterceptor(streamInts...),
	}
	if srv.config.TlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.config.TlsConf)))
	}
	if len(srv.config.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.config.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	if srv.config.Health {
		grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	}
	if srv.config.Reflection {
		reflection.Register(srv.Server)
	}
	return srv
}

func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return err
	}
	s.baseCtx = ctx
	s.config.logger.Infof("[gRPC] server listening on: %s", s.listener.Addr().String())
	s.health.Resume()
	err := s.Serve(s.listener)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) GracefulStop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

// listenAndEndpoint 启动grpc服务并且获取运行的ip地址
func (s *Server) listenAndEndpoint() error {
	if s.listener == nil {
		listener, err := net.Listen(s.config.Network, s.config.Address)
		if err != nil {
			return err
		}
		s.listener = listener
	}
	if s.endpoint == nil {
		addr, err := host.Extract(s.config.Address, s.listener)
		if err != nil {
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("grpc", s.config.TlsConf != nil), addr)
	}
	return nil
}
