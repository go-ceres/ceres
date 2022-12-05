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

package node

import (
	"context"
	"github.com/go-ceres/ceres/selector"
	"sync/atomic"
	"time"
)

const (
	defaultWeight = 100
)

var (
	_ selector.IWeightedNode        = (*DirectNode)(nil)
	_ selector.IWeightedNodeFactory = (*DirectNodeFactory)(nil)
)

// DirectNode is endpoint instance
type DirectNode struct {
	selector.INode
	lastPick int64 // 最后一次挑选时间时间
}

// DirectNodeFactory Direct节点的生产工厂
type DirectNodeFactory struct{}

// Create 创建节点
func (*DirectNodeFactory) Create(n selector.INode) selector.IWeightedNode {
	return &DirectNode{INode: n, lastPick: 0}
}

// Pick 选择节点
func (n *DirectNode) Pick() selector.DoneFunc {
	now := time.Now().UnixNano()
	atomic.StoreInt64(&n.lastPick, now)
	return func(ctx context.Context, di selector.InvokeDoneInfo) {}
}

// Weight 节点的有效权重
func (n *DirectNode) Weight() float64 {
	if n.InitialWeight() != nil {
		return float64(*n.InitialWeight())
	}
	return defaultWeight
}

// PickElapsed 计算节点挑选时间时长
func (n *DirectNode) PickElapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - atomic.LoadInt64(&n.lastPick))
}

// Raw 返回原始节点
func (n *DirectNode) Raw() selector.INode {
	return n.INode
}
