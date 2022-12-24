//    Copyright 2022. ceres
//    Author https://github.com/go-ceres/ceres
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package config

import (
	_ "github.com/go-ceres/ceres/pkg/common/codec/json"
	_ "github.com/go-ceres/ceres/pkg/common/codec/toml"
	_ "github.com/go-ceres/ceres/pkg/common/codec/xml"
	_ "github.com/go-ceres/ceres/pkg/common/codec/yaml"
)

var DefaultConfig, _ = New()

// Load 加载数据源
func Load(sources ...Source) error {
	return DefaultConfig.Load(sources...)
}

// Scan 扫描数据结构体
func Scan(v interface{}) error {
	return DefaultConfig.Scan(v)
}

// Get 获取原子值
func Get(path string) Value {
	return DefaultConfig.Get(path)
}

// Watch 配置一个监听者
func Watch(key string, ob Observer) error {
	return DefaultConfig.Watch(key, ob)
}

// Close 关闭配置服务
func Close() error {
	return DefaultConfig.Close()
}
