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
	"github.com/go-ceres/ceres/pkg/common/client/redis"
	"github.com/go-ceres/ceres/pkg/common/config"
	"github.com/go-ceres/ceres/pkg/common/logger"
)

var ModName = "cache.redis"

type Option func(o *Options)

type Options struct {
	Prefix string `json:"prefix"`
	client *redis.Client
	logger *logger.Logger
}

func DefaultOptions() *Options {
	return &Options{
		Prefix: "ceres:ceres",
		logger: logger.With(logger.FieldMod(ModName)),
	}
}

// ScanRawConfig 扫描无包装key
func ScanRawConfig(key string) *Options {
	conf := DefaultOptions()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 标准扫描配置
func ScanConfig(name ...string) *Options {
	key := "application.cache.redis"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

// WithClient 设置redis客户端
func WithClient(client *redis.Client) Option {
	return func(o *Options) {
		o.client = client
	}
}

// WithLogger 设置日志
func WithLogger(logger *logger.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

// WithOptions 手动设置参数
func (o *Options) WithOptions(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Build 构建缓存
func (o *Options) Build() *Cache {
	return NewWithOptions(o)
}
