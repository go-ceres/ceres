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

var defaultSelectorFactory = &wrapSelector{}

type wrapSelector struct{ ISelectorFactory }

// DefaultSelectorFactory 默认的选择器构建工厂
type DefaultSelectorFactory struct {
	WeightedNodeFactory IWeightedNodeFactory // 设置权重节点的创建工厂
	BalancerFactory     IBalancerFactory     // 负载均衡器创建工厂
}

func (b *DefaultSelectorFactory) Name() string {
	return b.BalancerFactory.Name()
}

func (b *DefaultSelectorFactory) Create() ISelector {
	return &DefaultSelector{
		WeightedNodeFactory: b.WeightedNodeFactory,
		Balancer:            b.BalancerFactory.Create(),
	}
}

// ISelectorFactory 节点选择器的创建工厂接口
type ISelectorFactory interface {
	Name() string      // 构建器名称
	Create() ISelector // 构建方法
}

// SetSelectorFactory 设置选择器构建工厂
func SetSelectorFactory(factory ISelectorFactory) {
	defaultSelectorFactory.ISelectorFactory = factory
}

// GetSelectorFactory 获取构建器
func GetSelectorFactory() ISelectorFactory {
	return defaultSelectorFactory.ISelectorFactory
}
