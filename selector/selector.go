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

package selector

import (
	"context"
	"github.com/go-ceres/ceres/errors"
	"sync/atomic"
)

// ErrNoAvailable 没有可用节点
var ErrNoAvailable = errors.ServiceUnavailable("no_available_node", "no_available_node")

type DoneFunc func(ctx context.Context, di InvokeDoneInfo)

// ISelector 节点选择器接口
type ISelector interface {
	ReBalancer
	Select(ctx context.Context, opts ...Option) (selected INode, done DoneFunc, err error)
}

// ReBalancer 节点平衡器
type ReBalancer interface {
	Apply(nodes []INode)
}

// InvokeDoneInfo 调用RPC完成之后传入回调函数的信息
type InvokeDoneInfo struct {
	ResponseMetadata       // 响应的metadata数据
	Err              error // 响应错误信息
	BytesSent        bool  // 是否有字节发送过服务器
	BytesReceived    bool  // 标识是否有从服务器接受到过数据
}

// ResponseMetadata 响应信息
type ResponseMetadata interface {
	Get(key string) string
}

// InvokeDoneCallBack 调用RPC完成后的回调
type InvokeDoneCallBack func(ctx context.Context, info InvokeDoneInfo)

// DefaultSelector 默认选择器
type DefaultSelector struct {
	WeightedNodeFactory IWeightedNodeFactory
	Balancer            IBalancer
	nodes               atomic.Value
}

// Select 选择节点函数
func (d *DefaultSelector) Select(ctx context.Context, opts ...Option) (INode, DoneFunc, error) {
	var (
		options    Options
		candidates []IWeightedNode
	)
	nodes, ok := d.nodes.Load().([]IWeightedNode)
	if !ok {
		return nil, nil, ErrNoAvailable
	}
	for _, o := range opts {
		o(&options)
	}
	if len(options.NodeFilters) > 0 {
		newNodes := make([]INode, len(nodes))
		for i, wc := range nodes {
			newNodes[i] = wc
		}
		for _, filter := range options.NodeFilters {
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
	wn, done, err := d.Balancer.Pick(ctx, candidates)
	if err != nil {
		return nil, nil, err
	}
	p, ok := FromPeerContext(ctx)
	if ok {
		p.Node = wn.Raw()
	}
	return wn.Raw(), done, nil
}

// Apply 修改节点数据
func (d *DefaultSelector) Apply(nodes []INode) {
	weightedNodes := make([]IWeightedNode, 0, len(nodes))
	for _, n := range nodes {
		weightedNodes = append(weightedNodes, d.WeightedNodeFactory.Create(n))
	}
	// TODO: Do not delete unchanged nodes
	d.nodes.Store(weightedNodes)
}
