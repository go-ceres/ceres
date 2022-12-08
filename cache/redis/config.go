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

package redis

import (
	"github.com/go-ceres/ceres/client/redis"
	"github.com/go-ceres/ceres/config"
	"github.com/go-ceres/ceres/logger"
)

var ModName = "cache.redis"

type Config struct {
	Prefix string `json:"prefix"`
	client *redis.Client
	logger *logger.Helper
}

func DefaultConfig() *Config {
	return &Config{
		Prefix: "ceres:ceres",
		logger: logger.With(logger.FieldMod("cache.redis")),
	}
}

// ScanRawConfig 扫描无包装key
func ScanRawConfig(key string) *Config {
	conf := DefaultConfig()
	if err := config.Get(key).Scan(conf); err != nil {
		logger.Fatalf("parse config error: %V", err)
	}
	return conf
}

// ScanConfig 标准扫描配置
func ScanConfig(name ...string) *Config {
	key := "ceres.application.cache.redis"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

// SetClient 设置redis客户端
func (c *Config) SetClient(client *redis.Client) *Config {
	c.client = client
	return c
}

// SetLogger 设置日志
func (c *Config) SetLogger(log logger.Logger) *Config {
	c.logger = logger.NewHelper(log, logger.FieldMod(ModName))
	return c
}

// Build 构建缓存
func (c *Config) Build() *Cache {
	return NewCache(c)
}
