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
	"errors"
	"github.com/go-ceres/ceres/internal/cycle"
	"github.com/go-ceres/ceres/internal/endpoint"
	"github.com/go-ceres/ceres/pkg/common/logger"
	"github.com/go-ceres/ceres/pkg/transport"
	"net/url"
	"strings"
	"time"
)

// Target 解析后的目标地址
type Target struct {
	Scheme    string
	Authority string
	Endpoint  string
}

// parseTarget 解析地址为目标地址
func parseTarget(endpoint string, insecure bool) (*Target, error) {
	if !strings.Contains(endpoint, "://") {
		if insecure {
			endpoint = "http://" + endpoint
		} else {
			endpoint = "https://" + endpoint
		}
	}
	parse, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	target := &Target{
		Scheme:    parse.Scheme,
		Authority: parse.Host,
	}
	if len(parse.Path) > 1 {
		target.Endpoint = parse.Path[1:]
	}
	return target, nil
}

// resolver 服务发现客户端
type resolver struct {
	discovery transport.Discover //
	selector  transport.Selector // 节点平衡器
	target    *Target            // 目标地址
	watcher   transport.Watcher  // 服务监听者
	insecure  bool               // 是否是不安全的
	logger    *logger.Logger     // 日志
	cycle     *cycle.Cycle
}

// newResolver 节点解析器
func newResolver(ctx context.Context, log *logger.Logger, discovery transport.Discover, target *Target, selector transport.Selector, block, insecure bool) (*resolver, error) {
	watcher, err := discovery.Watch(ctx, target.Endpoint)
	if err != nil {
		return nil, err
	}
	r := &resolver{
		logger:    log,
		target:    target,
		watcher:   watcher,
		selector:  selector,
		insecure:  insecure,
		discovery: discovery,
		cycle:     cycle.NewCycle(),
	}
	if block {
		r.cycle.Run(func() error {
			for {
				services, err := watcher.Next()
				if err != nil {
					return err
				}
				if r.update(services) {
					return nil
				}
			}
		})
		select {
		case err := <-r.cycle.Wait():
			if err != nil {
				stopErr := watcher.Stop()
				if stopErr != nil {
					log.Errorf("failed to http client watch stop: %v, error: %+v", target, stopErr)
				}
				return nil, err
			}
		case <-ctx.Done():
			log.Errorf("http client watch service %v reaching context deadline!", target)
			stopErr := watcher.Stop()
			if stopErr != nil {
				log.Errorf("failed to http client watch stop: %v, error: %+v", target, stopErr)
			}
			return nil, ctx.Err()
		}

	}
	go func() {
		_ = r.run()
	}()
	return r, nil
}

// run 运行服务发现
func (r *resolver) run() error {
	for {
		services, err := r.watcher.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			r.logger.Errorf("http client watch service %v got unexpected error:=%v", r.target, err)
			time.Sleep(time.Second)
			continue
		}
		r.update(services)
	}
}

// update 修改数据
func (r *resolver) update(services []*transport.ServiceInfo) bool {
	nodes := make([]transport.Node, 0)
	for _, ins := range services {
		ept, err := endpoint.ParseEndpoint(ins.Endpoints, endpoint.Scheme("http", !r.insecure))
		if err != nil {
			r.logger.Errorf("Failed to parse (%v) discovery endpoint: %v error %v", r.target, ins.Endpoints, err)
			continue
		}
		if ept == "" {
			continue
		}
		nodes = append(nodes, transport.NewNode("http", ept, ins))
	}
	if len(nodes) == 0 {
		r.logger.Warnf("[http resolver]Zero endpoint found,refused to write,set: %s ins: %v", r.target.Endpoint, nodes)
		return false
	}
	r.selector.Store(nodes)
	return true
}

// Close 停止服务监控
func (r *resolver) Close() error {
	return r.watcher.Stop()
}
