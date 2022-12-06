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
	"github.com/go-ceres/ceres/client/grpc/discovery"
	"github.com/go-ceres/ceres/logger"
	"github.com/go-ceres/ceres/middleware"
	"github.com/go-ceres/ceres/selector"
	"github.com/go-ceres/ceres/selector/wrr"
	"github.com/go-ceres/ceres/transport"
	trGrpc "github.com/go-ceres/ceres/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	grpcmd "google.golang.org/grpc/metadata"
	"time"
)

func init() {

	if selector.GetSelectorFactory() == nil {
		selector.SetSelectorFactory(wrr.NewSelectorFactory())
	}
}

// newGrpcClient 创建客户端
func newGrpcClient(c *Config) *grpc.ClientConn {
	var ctx = context.Background()
	var dialOptions = c.dialOpts
	if c.Block {
		if c.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, c.DialTimeout)
			defer cancel()
		}

		dialOptions = append(dialOptions, grpc.WithBlock())
	}
	ints := []grpc.UnaryClientInterceptor{
		unaryClientInterceptor(c.middleware, c.Timeout, c.filters),
	}
	if len(c.interceptors) > 0 {
		ints = append(ints, c.interceptors...)
	}
	if c.Debug {
		ints = append(ints, debugUnaryClientInterceptor(c.Endpoint))
	}
	dialOptions = append(dialOptions,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, c.balancer)),
		grpc.WithChainUnaryInterceptor(ints...),
	)
	if c.discovery != nil {
		dialOptions = append(dialOptions, grpc.WithResolvers(
			discovery.NewBuilder(c.discovery, discovery.WithInsecure(c.Insecure)),
		))
	}
	if c.Insecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if c.TlsConfig != nil {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(c.TlsConfig)))
	}
	cc, err := grpc.DialContext(ctx, c.Endpoint, dialOptions...)
	if err != nil {
		if c.OnDialError == "panic" {
			c.logger.Panic("dial grpc server panic err", logger.FieldError(err))
		} else {
			c.logger.Error("dial grpc server panic err", logger.FieldError(err))
		}
	}
	logger.Infof("grpc client started")
	return cc
}

// unaryClientInterceptor 过滤器设置
func unaryClientInterceptor(ms []middleware.Middleware, timeout time.Duration, filters []selector.NodeFilter) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		tr := new(trGrpc.Transport)
		tr.SetEndpoint(cc.Target())
		tr.SetOperation(method)
		tr.SetRequestHeader(trGrpc.HeaderInstance{})
		tr.SetNodeFilters(filters)
		ctx = transport.NewClientContext(ctx, tr)
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				header := tr.RequestHeader()
				keys := header.Keys()
				keyvals := make([]string, 0, len(keys))
				for _, k := range keys {
					keyvals = append(keyvals, k, header.Get(k))
				}
				ctx = grpcmd.AppendToOutgoingContext(ctx, keyvals...)
			}
			return reply, invoker(ctx, method, req, reply, cc, opts...)
		}
		if len(ms) > 0 {
			h = middleware.Chain(ms...)(h)
		}
		var p selector.Peer
		ctx = selector.NewPeerContext(ctx, &p)
		_, err := h(ctx, req)
		return err
	}
}
