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
	ic "github.com/go-ceres/ceres/internal/context"
	"github.com/go-ceres/ceres/middleware"
	"github.com/go-ceres/ceres/transport"
	trGrpc "github.com/go-ceres/ceres/transport/grpc"
	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"
)

// wrappedStream is rewrite grpc stream's context
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func NewWrappedStream(ctx context.Context, stream grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{
		ServerStream: stream,
		ctx:          ctx,
	}
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// unaryServerInterceptor is a gRPC unary server interceptor
func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := ic.Merge(ctx, s.baseCtx)
		defer cancel()
		md, _ := grpcMetadata.FromIncomingContext(ctx)
		replyHeader := grpcMetadata.MD{}
		tr := new(trGrpc.Transport)
		tr.SetOperation(info.FullMethod)
		tr.SetRequestHeader(trGrpc.HeaderInstance(md))
		tr.SetReplyHeader(trGrpc.HeaderInstance(replyHeader))
		if s.endpoint != nil {
			tr.SetEndpoint(s.endpoint.String())
		}
		ctx = transport.NewServerContext(ctx, tr)
		if s.config.Timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, s.config.Timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		if next := s.config.middleware.Match(tr.Operation()); len(next) > 0 {
			h = middleware.Chain(next...)(h)
		}
		reply, err := h(ctx, req)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return reply, err
	}
}

// streamServerInterceptor is a gRPC stream server interceptor
func (s *Server) streamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, cancel := ic.Merge(ss.Context(), s.baseCtx)
		defer cancel()
		md, _ := grpcMetadata.FromIncomingContext(ctx)
		replyHeader := grpcMetadata.MD{}
		tr := new(trGrpc.Transport)
		tr.SetEndpoint(s.endpoint.String())
		tr.SetOperation(info.FullMethod)
		tr.SetRequestHeader(trGrpc.HeaderInstance(md))
		tr.SetReplyHeader(trGrpc.HeaderInstance(replyHeader))
		ctx = transport.NewServerContext(ctx, tr)
		ws := NewWrappedStream(ctx, ss)
		err := handler(srv, ws)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return err
	}
}
