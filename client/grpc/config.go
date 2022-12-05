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
	"github.com/go-ceres/ceres/logger"
	"github.com/go-ceres/ceres/middleware"
	"github.com/go-ceres/ceres/registry"
	"github.com/go-ceres/ceres/selector"
	"google.golang.org/grpc"
	"time"
)

const ModName = "client.grpc"

// Config 配置信息
type Config struct {
	Block        bool                          `json:"block"`       // 是否一直等待连接成功
	Insecure     bool                          `json:"insecure"`    // 是否忽略安全
	Debug        bool                          `json:"debug"`       // 是否调试模式
	Endpoint     string                        `json:"endpoint"`    // 连接地址
	TlsConfig    *tls.Config                   `json:"tlsConfig"`   // 安全认证
	Timeout      time.Duration                 `json:"timeout"`     // 超时
	DialTimeout  time.Duration                 `json:"dialTimeout"` // 调用超时
	OnDialError  string                        `json:"OnDialError"` // 构建错误处理 panic | error
	balancer     string                        // 负载均衡器名称
	discovery    registry.Registry             // 服务发现
	middleware   []middleware.Middleware       // 中间件
	interceptors []grpc.UnaryClientInterceptor // 拦截器
	dialOpts     []grpc.DialOption             // 调用参数
	filters      []selector.NodeFilter         // 节点选择器
	logger       *logger.Helper                // 日志
}

func DefaultConfig() *Config {
	return &Config{
		Endpoint:    "",
		Block:       true,
		Timeout:     3 * time.Second,
		DialTimeout: 3 * time.Second,
		balancer:    "selector",
		Insecure:    true,
		Debug:       false,
		logger:      logger.With(logger.FieldMod(ModName)),
	}
}

// ScanRawConfig 扫描原始
func ScanRawConfig(key string) *Config {
	conf := DefaultConfig()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 扫描配置
func ScanConfig(name ...string) *Config {
	key := "ceres.application.client.grpc"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

func (c *Config) SetBlock(Block bool) *Config {
	c.Block = Block
	return c
}

func (c *Config) SetInsecure(Insecure bool) *Config {
	c.Insecure = Insecure
	return c
}

func (c *Config) SetDebug(Debug bool) *Config {
	c.Debug = Debug
	return c
}

func (c *Config) SetEndpoint(Endpoint string) *Config {
	c.Endpoint = Endpoint
	return c
}

func (c *Config) SetTlsConfig(TlsConfig *tls.Config) *Config {
	c.TlsConfig = TlsConfig
	return c
}

func (c *Config) SetTimeout(Timeout time.Duration) *Config {
	c.Timeout = Timeout
	return c
}

// SetDialTimeout 设置请求超时时间
func (c *Config) SetDialTimeout(DialTimeout time.Duration) *Config {
	c.DialTimeout = DialTimeout
	return c
}

// SetOnDialError 错误处理级别 "panic" or "error"
func (c *Config) SetOnDialError(OnDialError string) *Config {
	c.OnDialError = OnDialError
	return c
}

// SetMiddleware 设置内部中间件
func (c *Config) SetMiddleware(m ...middleware.Middleware) *Config {
	c.middleware = m
	return c
}

// SetDiscovery 设置服务发现
func (c *Config) SetDiscovery(d registry.Registry) *Config {
	c.discovery = d
	return c
}

// SetTLSConfig 设置TLS认证
func (c *Config) SetTLSConfig(t *tls.Config) *Config {
	c.TlsConfig = t
	return c
}

// SetUnaryInterceptor 设置拦截器
func (c *Config) SetUnaryInterceptor(ins ...grpc.UnaryClientInterceptor) *Config {
	c.interceptors = ins
	return c
}

// SetDialOption 设置调用参数
func (c *Config) SetDialOption(opts ...grpc.DialOption) *Config {
	c.dialOpts = opts
	return c
}

// SetNodeFilter 设置拦截器
func (c *Config) SetNodeFilter(filters ...selector.NodeFilter) *Config {
	c.filters = filters
	return c
}

// Build 构建连接
func (c *Config) Build() *grpc.ClientConn {
	return newGrpcClient(c)
}
