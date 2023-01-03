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
	"encoding/json"
	"errors"
	"github.com/go-ceres/ceres/internal/endpoint"
	"github.com/go-ceres/ceres/pkg/common/logger"
	"github.com/go-ceres/ceres/pkg/transport"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"strings"
	"time"
)

var (
	_ resolver.Builder  = (*directResolverBuilder)(nil)
	_ resolver.Builder  = (*discoveryResolverBuilder)(nil)
	_ resolver.Resolver = (*directResolver)(nil)
	_ resolver.Resolver = (*discoveryResolver)(nil)
)

func init() {
	resolver.Register(&directResolverBuilder{})
}

// directResolverBuilder 支持direct://<authority>/127.0.0.1:9000,127.0.0.2:9000 格式的
type directResolverBuilder struct{}

func (d *directResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	addrs := make([]resolver.Address, 0)
	for _, addr := range strings.Split(strings.TrimPrefix(target.URL.Path, "/"), ",") {
		addrs = append(addrs, resolver.Address{Addr: addr})
	}
	err := cc.UpdateState(resolver.State{
		Addresses: addrs,
	})
	if err != nil {
		return nil, err
	}
	return newDirectResolver(), nil
}

func (d *directResolverBuilder) Scheme() string {
	return "direct"
}

type directResolver struct{}

func newDirectResolver() resolver.Resolver {
	return &directResolver{}
}

func (r *directResolver) Close() {
}

func (r *directResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

// discoveryResolverBuilder 服务发现构建器
type discoveryResolverBuilder struct {
	discoverer       transport.Discover
	timeout          time.Duration
	logger           *logger.Logger
	insecure         bool
	debugLogDisabled bool
}

func (d *discoveryResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	watchRes := &struct {
		err error
		w   transport.Watcher
	}{}

	done := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		w, err := d.discoverer.Watch(ctx, strings.TrimPrefix(target.URL.Path, "/"))
		watchRes.w = w
		watchRes.err = err
		close(done)
	}()

	var err error
	select {
	case <-done:
		err = watchRes.err
	case <-time.After(d.timeout):
		err = errors.New("discovery create watcher overtime")
	}
	if err != nil {
		cancel()
		return nil, err
	}
	r := &discoveryResolver{
		w:                watchRes.w,
		cc:               cc,
		ctx:              ctx,
		cancel:           cancel,
		insecure:         d.insecure,
		logger:           d.logger,
		debugLogDisabled: d.debugLogDisabled,
	}
	go r.watch()
	return r, nil
}

func (d *discoveryResolverBuilder) Scheme() string {
	return "discovery"
}

// discoveryResolver
type discoveryResolver struct {
	w                transport.Watcher
	cc               resolver.ClientConn
	ctx              context.Context
	cancel           context.CancelFunc
	logger           *logger.Logger
	insecure         bool
	debugLogDisabled bool
}

func (r *discoveryResolver) watch() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}
		ins, err := r.w.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			r.logger.Errorf("[resolver] Failed to watch discovery endpoint: %v", err)
			time.Sleep(time.Second)
			continue
		}
		r.update(ins)
	}
}

func (r *discoveryResolver) update(ins []*transport.ServiceInfo) {
	addrs := make([]resolver.Address, 0)
	endpoints := make(map[string]struct{})
	for _, in := range ins {
		ept, err := endpoint.ParseEndpoint(in.Endpoints, endpoint.Scheme("grpc", !r.insecure))
		if err != nil {
			logger.Errorf("[resolver] Failed to parse discovery endpoint: %v", err)
			continue
		}
		if ept == "" {
			continue
		}
		// filter redundant endpoints
		if _, ok := endpoints[ept]; ok {
			continue
		}
		endpoints[ept] = struct{}{}
		addr := resolver.Address{
			ServerName: in.Name,
			Attributes: parseAttributes(in.Metadata),
			Addr:       ept,
		}
		addr.Attributes = addr.Attributes.WithValue("rawServiceInstance", in)
		addrs = append(addrs, addr)
	}
	if len(addrs) == 0 {
		r.logger.Warnf("[resolver] not found endpoint list,refused to write, instances: %v", ins)
		return
	}
	err := r.cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		r.logger.Errorf("[resolver] failed to update state: %s", err)
	}

	if !r.debugLogDisabled {
		b, _ := json.Marshal(ins)
		r.logger.Infof("[resolver] update instances: %s", b)
	}
}

func (r *discoveryResolver) Close() {
	r.cancel()
	err := r.w.Stop()
	if err != nil {
		r.logger.Errorf("[resolver] failed to watch top: %s", err)
	}
}

func (r *discoveryResolver) ResolveNow(options resolver.ResolveNowOptions) {}

func parseAttributes(md map[string]string) *attributes.Attributes {
	var a *attributes.Attributes
	for k, v := range md {
		if a == nil {
			a = attributes.New(k, v)
		} else {
			a = a.WithValue(k, v)
		}
	}
	return a
}
