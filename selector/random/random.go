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

package random

import (
	"context"
	"github.com/go-ceres/ceres/selector"
	"github.com/go-ceres/ceres/selector/node"
	"math/rand"
)

func init() {
	selector.SetSelectorFactory(NewRandomFactory())
}

// Name 负载均衡器名称
const Name = "random"

// Balancer 负载均衡器
type Balancer struct{}

// Pick 选择节点的方法
func (b *Balancer) Pick(ctx context.Context, nodes []selector.IWeightedNode) (selector.IWeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	cur := rand.Intn(len(nodes))
	selected := nodes[cur]
	d := selected.Pick()
	return selected, d, nil
}

// BalancerFactory 负载均衡器创建工厂
type BalancerFactory struct{}

// Name 工厂名称
func (b *BalancerFactory) Name() string {
	return Name
}

// Create 创建函数
func (b *BalancerFactory) Create() selector.IBalancer {
	return &Balancer{}
}

// NewRandomFactory 创建选择器工厂
func NewRandomFactory() selector.ISelectorFactory {
	return &selector.DefaultSelectorFactory{
		BalancerFactory:     &BalancerFactory{},
		WeightedNodeFactory: &node.DirectNodeFactory{},
	}
}
