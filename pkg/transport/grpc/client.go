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
	"fmt"
	"github.com/fatih/color"
	"github.com/go-ceres/ceres/pkg/common/logger"
	"github.com/go-ceres/ceres/pkg/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"time"
)

// NewClient 创建客户端
func NewClient(opts ...ClientOption) (*grpc.ClientConn, error) {
	options := DefaultClientOptions()
	for _, opt := range opts {
		opt(options)
	}
	return NewClientWithOptions(options)
}

// NewClientWithOptions 创建客户端
func NewClientWithOptions(options *ClientOptions) (*grpc.ClientConn, error) {
	var ctx = context.Background()
	var dialOptions = options.dialOpts
	if options.Block {
		if options.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, options.DialTimeout)
			defer cancel()
		}

		dialOptions = append(dialOptions, grpc.WithBlock())
	}
	ints := []grpc.UnaryClientInterceptor{
		unaryClientInterceptor(options.middleware, options.Timeout, options.filters),
	}
	if len(options.interceptors) > 0 {
		ints = append(ints, options.interceptors...)
	}
	if options.Debug {
		ints = append(ints, debugUnaryClientInterceptor(options.Endpoint))
	}

	dialOptions = append(dialOptions,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, options.Balancer)),
		grpc.WithChainUnaryInterceptor(ints...),
	)
	if options.discovery != nil {
		b := &discoveryResolverBuilder{
			discoverer:       options.discovery,
			timeout:          time.Second * 10,
			insecure:         options.Insecure,
			debugLogDisabled: false,
		}
		dialOptions = append(dialOptions, grpc.WithResolvers(b))
	}
	if options.Insecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if options.TlsConfig != nil {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(options.TlsConfig)))
	}
	cc, err := grpc.DialContext(ctx, options.Endpoint, dialOptions...)
	if err != nil {
		if options.OnDialError == "panic" {
			options.logger.Panic("dial grpc server panic err", logger.FieldError(err))
		} else {
			options.logger.Error("dial grpc server panic err", logger.FieldError(err))
		}
	}
	options.logger.Infof("grpc client started")
	return cc, nil
}

// debugUnaryClientInterceptor 日志拦截器
func debugUnaryClientInterceptor(addr string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var p peer.Peer
		prefix := fmt.Sprintf("[%s]", addr)
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			prefix = prefix + "(" + remote.Addr.String() + ")"
		}

		fmt.Printf("%-50s[%s] => %s\n", color.GreenString(prefix), time.Now().Format("04:05.000"), color.GreenString("Send: "+method+" | %v", req))
		err := invoker(ctx, method, req, reply, cc, append(opts, grpc.Peer(&p))...)
		if err != nil {
			fmt.Printf("%-50s[%s] => %s\n", color.RedString(prefix), time.Now().Format("04:05.000"), color.RedString("Erro: %v", err.Error()))
		} else {
			fmt.Printf("%-50s[%s] => %s\n", color.GreenString(prefix), time.Now().Format("04:05.000"), color.GreenString("Recv: %v", reply))
		}

		return err
	}
}

// unaryClientInterceptor 过滤器设置
func unaryClientInterceptor(ms []transport.Middleware, timeout time.Duration, filters []transport.NodeFilter) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		metadata := AcquireMetadata()
		defer ReleaseMetadata(metadata)
		metadata.endpoint = cc.Target()
		metadata.operation = method
		metadata.requestHeader = headerCarrier{}
		metadata.requestHeader = headerCarrier{}
		metadata.nodeFilters = filters
		ctx = transport.NewMetadataClientContext(ctx, metadata)
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.MetadataFromClientContext(ctx); ok {
				header := tr.RequestHeader()
				keys := header.Keys()
				KeyValues := make([]string, 0, len(keys))
				for _, key := range keys {
					KeyValues = append(KeyValues, key, header.Get(key))
				}
				ctx = grpcmd.AppendToOutgoingContext(ctx, KeyValues...)
			}
			return reply, invoker(ctx, method, req, reply, cc, opts...)
		}
		if len(ms) > 0 {
			h = transport.MiddlewareChain(ms...)(h)
		}
		var p transport.Peer
		ctx = transport.NewPeerContext(ctx, &p)
		_, err := h(ctx, req)
		return err
	}
}
