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
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var (
	forcePick                 = time.Second * 3
	_         Balancer        = (*DefaultBalancer)(nil)
	_         BalancerBuilder = (*DefaultBalancerBuilder)(nil)
)

// BalancerBuilder 负载均衡器构建接口
type BalancerBuilder interface {
	Build() Balancer
}

// Balancer 负载均衡器接口
type Balancer interface {
	Pick(ctx context.Context, nodes []IWeightedNode) (selected IWeightedNode, done DoneFunc, err error)
}

// DefaultBalancerBuilder 默认实现
type DefaultBalancerBuilder struct{}

// Build 构建默认负载均衡器
func (b *DefaultBalancerBuilder) Build() Balancer {
	return &DefaultBalancer{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// DefaultBalancer 默认实现
type DefaultBalancer struct {
	mu     sync.Mutex
	r      *rand.Rand
	picked int64
}

// prePick 选择两个不同的节点
func (bc *DefaultBalancer) prePick(nodes []IWeightedNode) (nodeA IWeightedNode, nodeB IWeightedNode) {
	bc.mu.Lock()
	a := bc.r.Intn(len(nodes))
	b := bc.r.Intn(len(nodes) - 1)
	bc.mu.Unlock()
	if b >= a {
		b = b + 1
	}
	nodeA, nodeB = nodes[a], nodes[b]
	return
}

// Pick 匹配节点
func (bc *DefaultBalancer) Pick(_ context.Context, nodes []IWeightedNode) (IWeightedNode, DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, ErrNoAvailable
	}
	if len(nodes) == 1 {
		done := nodes[0].Pick()
		return nodes[0], done, nil
	}

	var pc, upc IWeightedNode
	nodeA, nodeB := bc.prePick(nodes)
	// meta.Weight is the weight set by the service publisher in discovery
	if nodeB.Weight() > nodeA.Weight() {
		pc, upc = nodeB, nodeA
	} else {
		pc, upc = nodeA, nodeB
	}

	// If the failed node has never been selected once during forceGap, it is forced to be selected once
	// Take advantage of forced opportunities to trigger updates of success rate and delay
	if upc.PickElapsed() > forcePick && atomic.CompareAndSwapInt64(&bc.picked, 0, 1) {
		pc = upc
		atomic.StoreInt64(&bc.picked, 0)
	}
	done := pc.Pick()
	return pc, done, nil
}
