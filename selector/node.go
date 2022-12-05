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
	"github.com/go-ceres/ceres/registry"
	"strconv"
	"time"
)

// IWeightedNode 权重节点
type IWeightedNode interface {
	INode
	Raw() INode                 //返回原始节点
	Weight() float64            // 运行时计算出的权重
	Pick() DoneFunc             // 获取节点
	PickElapsed() time.Duration //最近一次获取节点到现在的时间差
}

// IWeightedNodeFactory 权重节点的生产工厂接口
type IWeightedNodeFactory interface {
	Create(node INode) IWeightedNode
}

// INode 节点接口
type INode interface {
	Scheme() string                     // 节点协议
	Address() string                    // 服务地址
	InitialWeight() *int64              // 初始化权重
	ServiceInfo() *registry.ServiceInfo // 服务注册信息
}

// DefaultNode 默认节点
type DefaultNode struct {
	scheme      string                // 协议
	address     string                // 地址
	weight      *int64                // 权重
	serviceInfo *registry.ServiceInfo // 服务信息
}

// InitialWeight 权重
func (d DefaultNode) InitialWeight() *int64 {
	return d.weight
}

// Scheme 获取协议
func (d DefaultNode) Scheme() string {
	return d.scheme
}

// Address 获取服务地址
func (d DefaultNode) Address() string {
	return d.address
}

// ServiceInfo 获取服务信息
func (d DefaultNode) ServiceInfo() *registry.ServiceInfo {
	return d.serviceInfo
}

// NewNode 创建一个节点
func NewNode(scheme string, address string, info *registry.ServiceInfo) INode {
	node := &DefaultNode{
		scheme:      scheme,
		address:     address,
		serviceInfo: info,
	}
	if info != nil {
		if str, ok := info.Metadata["weight"]; ok {
			if weight, err := strconv.ParseInt(str, 10, 64); err == nil {
				node.weight = &weight
			}
		}
	}
	return node
}
