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

package grpc

import (
	"crypto/tls"
	"github.com/go-ceres/ceres/config"
	"github.com/go-ceres/ceres/internal/matcher"
	"github.com/go-ceres/ceres/logger"
	"google.golang.org/grpc"
	"time"
)

const ModName = "server.grpc"

// Config 配置信息
type Config struct {
	Network            string                         // net.listen network
	Address            string                         // 服务地址
	Timeout            time.Duration                  // 超时时间
	SlowQueryThreshold time.Duration                  // 在debug模式下慢查询阈值,如果请求到响应超过此值，则会打印日志
	TlsConf            *tls.Config                    // tls配置信息
	Reflection         bool                           // 是否反射服务
	Health             bool                           // 是否设置健康服务
	Debug              bool                           // 是否开启调试模式
	middleware         matcher.Matcher                // 中间件
	unaryInts          []grpc.UnaryServerInterceptor  // grpc服务拦截器
	streamInts         []grpc.StreamServerInterceptor // 数据流服务拦截器
	grpcOpts           []grpc.ServerOption            // grpc服务额外参数
	logger             *logger.Helper                 //日志组件
}

// Build 构建
func (c *Config) Build() *Server {
	return New(c)
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Network:            "tcp",
		Address:            "127.0.0.1:5201",
		SlowQueryThreshold: 3 * time.Second,
		Debug:              false,
		middleware:         matcher.New(),
		logger:             logger.With(logger.FieldMod(ModName)),
	}
}

// ScanRawConfig 扫描配置
func ScanRawConfig(key string) *Config {
	conf := DefaultConfig()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 标准配置扫描
func ScanConfig(name ...string) *Config {
	key := "ceres.application.server.grpc"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

func (c *Config) SetHealth(health bool) *Config {
	c.Health = health
	return c
}

func (c *Config) SetNetwork(Network string) *Config {
	c.Network = Network
	return c
}

func (c *Config) SetAddress(Address string) *Config {
	c.Address = Address
	return c
}

func (c *Config) SetSlowQueryThreshold(SlowQueryThreshold time.Duration) *Config {
	c.SlowQueryThreshold = SlowQueryThreshold
	return c
}

func (c *Config) SetTlsConf(TlsConf *tls.Config) *Config {
	c.TlsConf = TlsConf
	return c
}

func (c *Config) SetUnaryInts(unaryInts ...grpc.UnaryServerInterceptor) *Config {
	c.unaryInts = unaryInts
	return c
}

func (c *Config) SetStreamInts(streamInts ...grpc.StreamServerInterceptor) *Config {
	c.streamInts = streamInts
	return c
}

func (c *Config) SetGrpcOpts(grpcOpts ...grpc.ServerOption) *Config {
	c.grpcOpts = grpcOpts
	return c
}
