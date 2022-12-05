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

package wrr

import (
	"context"
	"github.com/go-ceres/ceres/selector"
	"github.com/go-ceres/ceres/selector/node"
	"sync"
)

const (
	// Name is wrr balancer name
	Name = "wrr"
)

var _ selector.IBalancer = (*Balancer)(nil) // Name is balancer name

func init() {
	selector.SetSelectorFactory(NewSelectorFactory())
}

// BalancerFactory is wrr builder
type BalancerFactory struct{}

// Create c创建负载均衡器函数
func (b *BalancerFactory) Create() selector.IBalancer {
	return &Balancer{currentWeight: make(map[string]float64)}
}
func (b *BalancerFactory) Name() string {
	return Name
}

// Balancer is a random balancer.
type Balancer struct {
	mu            sync.Mutex
	currentWeight map[string]float64
}

// Pick is pick a weighted node.
func (p *Balancer) Pick(_ context.Context, nodes []selector.IWeightedNode) (selector.IWeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	var totalWeight float64
	var selected selector.IWeightedNode
	var selectWeight float64

	// nginx wrr load balancing algorithm: http://blog.csdn.net/zhangskd/article/details/50194069
	p.mu.Lock()
	for _, node := range nodes {
		totalWeight += node.Weight()
		cwt := p.currentWeight[node.Address()]
		// current += effectiveWeight
		cwt += node.Weight()
		p.currentWeight[node.Address()] = cwt
		if selected == nil || selectWeight < cwt {
			selectWeight = cwt
			selected = node
		}
	}
	p.currentWeight[selected.Address()] = selectWeight - totalWeight
	p.mu.Unlock()

	d := selected.Pick()
	return selected, d, nil
}

// NewSelectorFactory 返回选择器的创建工厂
func NewSelectorFactory() selector.ISelectorFactory {
	return &selector.DefaultSelectorFactory{
		BalancerFactory:     &BalancerFactory{},
		WeightedNodeFactory: &node.DirectNodeFactory{},
	}
}
