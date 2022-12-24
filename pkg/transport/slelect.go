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

package transport

import (
	"context"
	"github.com/go-ceres/ceres/pkg/common/errors"
	"sync/atomic"
)

const selectorName = "p2c"

var (
	selectorBuilder SelectorBuilder
	_               Selector = (*DefaultSelector)(nil)
	// ErrNoAvailable 没有可用节点
	ErrNoAvailable = errors.ServiceUnavailable("NO_AVAILABLE_NODE", "no available node")
)

func init() {
	selectorBuilder = &DefaultSelectorBuilder{
		Name:             selectorName,
		BalancerBuilder:  &DefaultBalancerBuilder{},
		WightNodeBuilder: &DefaultWeightNodeBuilder{},
	}
}

// NodeFilter 过滤器
type NodeFilter func(ctx context.Context, nodes []Node) []Node

// SelectOption 选择接节点参数
type SelectOption func(o *SelectOptions)

// WithNodeFilter 设置节点过滤器
func WithNodeFilter(filters ...NodeFilter) SelectOption {
	return func(o *SelectOptions) {
		o.filters = filters
	}
}

// SelectOptions 选择节点参数信息
type SelectOptions struct {
	filters []NodeFilter
}

// SelectorBuilder 选择器构建器接口
type SelectorBuilder interface {
	Build() Selector
}

// DefaultSelectorBuilder 默认的选择器的构建器
type DefaultSelectorBuilder struct {
	Name             string
	WightNodeBuilder WeightNodeBuilder
	BalancerBuilder  BalancerBuilder
}

// Build 构建选择器方法
func (s *DefaultSelectorBuilder) Build() Selector {
	return &DefaultSelector{
		name:              s.Name,
		weightNodeBuilder: s.WightNodeBuilder,
		balancer:          s.BalancerBuilder.Build(),
	}
}

// Selector 选择器接口
type Selector interface {
	// Name 选择器名称
	Name() string
	// Store 平衡方法
	Store(nodes []Node)
	// Select 选择节点
	Select(ctx context.Context, opts ...SelectOption) (selected Node, doneFunc DoneFunc, err error)
}

// DefaultSelector 默认的选择器实现
type DefaultSelector struct {
	weightNodeBuilder WeightNodeBuilder // 权重节点构建器
	balancer          Balancer          // 负载均衡构建器
	name              string            // 选择器名称
	nodes             atomic.Value
}

// Name 选择器名称
func (s *DefaultSelector) Name() string {
	return s.name
}

// Store 存储节点
func (s *DefaultSelector) Store(nodes []Node) {
	weightedNodes := make([]IWeightedNode, 0, len(nodes))
	for _, n := range nodes {
		weightedNodes = append(weightedNodes, s.weightNodeBuilder.Build(n))
	}
	s.nodes.Store(weightedNodes)
}

// Select 选择节点
func (s *DefaultSelector) Select(ctx context.Context, opts ...SelectOption) (selected Node, doneFunc DoneFunc, err error) {
	var (
		options    SelectOptions
		candidates []IWeightedNode
	)
	nodes, ok := s.nodes.Load().([]IWeightedNode)
	if !ok {
		return nil, nil, ErrNoAvailable
	}
	for _, o := range opts {
		o(&options)
	}
	if len(options.filters) > 0 {
		newNodes := make([]Node, len(nodes))
		for i, wc := range nodes {
			newNodes[i] = wc
		}
		for _, filter := range options.filters {
			newNodes = filter(ctx, newNodes)
		}
		candidates = make([]IWeightedNode, len(newNodes))
		for i, n := range newNodes {
			candidates[i] = n.(IWeightedNode)
		}
	} else {
		candidates = nodes
	}

	if len(candidates) == 0 {
		return nil, nil, ErrNoAvailable
	}
	wn, done, err := s.balancer.Pick(ctx, candidates)
	if err != nil {
		return nil, nil, err
	}
	p, ok := FromPeerContext(ctx)
	if ok {
		p.Node = wn.Raw()
	}
	return wn.Raw(), done, nil
}

// SetSelectorBuilder 注册节点选择器
func SetSelectorBuilder(builder SelectorBuilder) {
	selectorBuilder = builder
}

// GetSelectorBuilder 获取
func GetSelectorBuilder() SelectorBuilder {
	return selectorBuilder
}
