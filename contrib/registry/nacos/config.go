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

package nacos

import (
	"github.com/go-ceres/ceres/config"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/file"
	"os"
)

// Config 配置信息
type Config struct {
	Address []string `json:"address"` // 服务器地址
	Weight  float64  `json:"weight"`  // 初始化权重
	Cluster string   `json:"cluster"` // 集群
	Group   string   `json:"group"`   // 分组
	Kind    string   `json:"kind"`    // 协议
	// 客户端配置
	*constant.ClientConfig
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Cluster: "DEFAULT",
		Group:   constant.DEFAULT_GROUP,
		Weight:  100,
		Kind:    "grpc",
		ClientConfig: &constant.ClientConfig{
			TimeoutMs:            10 * 1000,
			BeatInterval:         5 * 1000,
			OpenKMS:              false,
			CacheDir:             file.GetCurrentPath() + string(os.PathSeparator) + "cache",
			UpdateThreadNum:      20,
			NotLoadCacheAtStart:  true,
			UpdateCacheWhenEmpty: false,
			LogDir:               file.GetCurrentPath() + string(os.PathSeparator) + "log",
			RotateTime:           "24h",
			MaxAge:               3,
			LogLevel:             "info",
		},
	}
}

// ScanRawConfig 扫描无封装key
func ScanRawConfig(key string) *Config {
	conf := DefaultConfig()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 扫描配置
func ScanConfig(name ...string) *Config {
	key := "ceres.application.registry.nacos"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

// Build 构建注册中心
func (c *Config) Build() *Registry {
	return newRegistry(c)
}
