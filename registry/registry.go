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

package registry

import (
	"context"
)

// Registry 注册中心接口
type Registry interface {
	// Register 注册服务
	Register(ctx context.Context, service *ServiceInfo) error
	// Deregister 注销服务
	Deregister(ctx context.Context, service *ServiceInfo) error
	// GetService 获取指定服务名的所有服务
	GetService(ctx context.Context, serviceName string) ([]*ServiceInfo, error)
	// Watch 创建服务监听器
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

// Watcher 服务监听者
type Watcher interface {
	// Next 监听服务变化后返回
	Next() ([]*ServiceInfo, error)
	// Stop 关闭监听器
	Stop() error
}

// ServiceInfo 服务实例信息
type ServiceInfo struct {
	ID        string            `json:"id"`        // 服务id
	Name      string            `json:"name"`      // 服务名称
	Version   string            `json:"version"`   // 当前版本
	Endpoints []string          `json:"endpoints"` //服务地址
	Metadata  map[string]string `json:"metadata"`  //附加信息
}
