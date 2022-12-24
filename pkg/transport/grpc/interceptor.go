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
	"github.com/go-ceres/ceres/pkg/transport"
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
		ctx, cancel := ic.Merge(ctx, s.baseContext)
		defer cancel()
		md, _ := grpcMetadata.FromIncomingContext(ctx)
		replyHeader := grpcMetadata.MD{}
		metadata := AcquireMetadata()
		defer ReleaseMetadata(metadata)
		metadata.operation = info.FullMethod
		metadata.requestHeader = headerCarrier(md)
		metadata.replyHeader = headerCarrier(replyHeader)
		// 绑定header
		err := defaultBinding.Unmarshal(req, warpRequest(metadata), nil)
		if err != nil {
			return nil, err
		}
		if s.endpoint != nil {
			metadata.endpoint = s.endpoint.String()
		}
		ctx = transport.NewMetadataServerContext(ctx, metadata)
		if s.opts.Timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, s.opts.Timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		if next := s.opts.middleware.Match(metadata.Operation()); len(next) > 0 {
			h = transport.MiddlewareChain(next...)(h)
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
		ctx, cancel := ic.Merge(ss.Context(), s.baseContext)
		defer cancel()
		tr := AcquireMetadata()
		md, _ := grpcMetadata.FromIncomingContext(ctx)
		replyHeader := grpcMetadata.MD{}
		tr.operation = info.FullMethod
		tr.requestHeader = headerCarrier(md)
		tr.replyHeader = headerCarrier(replyHeader)
		if s.endpoint != nil {
			tr.endpoint = s.endpoint.String()
		}
		ctx = transport.NewMetadataServerContext(ctx, tr)
		ws := NewWrappedStream(ctx, ss)
		err := handler(srv, ws)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return err
	}
}
