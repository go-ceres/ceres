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

package http

import (
	"context"
	"crypto/tls"
	"github.com/go-ceres/ceres/config"
	"github.com/go-ceres/ceres/logger"
	"github.com/go-ceres/ceres/middleware"
	"github.com/go-ceres/ceres/registry"
	"github.com/go-ceres/ceres/selector"
	"net/http"
	"time"
)

const ModName = "client.http"

// Config 客户端参数
type Config struct {
	Endpoint     string                  `json:"endpoint"`  // 请求地址
	UserAgent    string                  `json:"userAgent"` // 用户代理
	TlsConf      *tls.Config             `json:"tlsConf"`   // tls认证
	Timeout      time.Duration           `json:"timeout"`   // 超时时间
	Balancer     string                  `json:"balancer"`  // 负载聚合名
	Block        bool                    `json:"block"`     // 是否阻塞
	ctx          context.Context         // 上下文
	encoder      EncodeRequestFunc       // 请求编码器
	decoder      DecodeResponseFunc      // 响应解码器
	errorDecoder DecodeErrorFunc         // 错误解码器
	transport    http.RoundTripper       // 传输协议接口
	nodeFilters  []selector.NodeFilter   // 节点选择
	discovery    registry.Registry       // 服务发现
	middleware   []middleware.Middleware // 中间件
	logger       *logger.Helper          // 日志
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Endpoint:     "",
		UserAgent:    "",
		Block:        false,
		Timeout:      2000 * time.Millisecond,
		ctx:          context.Background(),
		encoder:      DefaultRequestEncoder,
		decoder:      DefaultResponseDecoder,
		errorDecoder: DefaultErrorDecoder,
		transport:    http.DefaultTransport,
		logger:       logger.With(logger.FieldMod(ModName)),
	}
}

func ScanRawConfig(key string) *Config {
	conf := DefaultConfig()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 扫描配置文件
func ScanConfig(name ...string) *Config {
	key := "ceres.application.client.http"
	if len(name) > 0 {
		key = key + name[0]
	}
	return ScanRawConfig(key)
}

// Build 根据配置文件构建客户端
func (c *Config) Build() *Client {
	client, err := New(c)
	if err != nil {
		c.logger.Panic("new client failed: ", logger.FieldError(err))
	}
	return client
}

// WithTransport 设置传输协议
func (c *Config) WithTransport(tripper http.RoundTripper) *Config {
	c.transport = tripper
	return c
}

// WithRequestEncoder 设置请求编码器
func (c *Config) WithRequestEncoder(encode EncodeRequestFunc) *Config {
	c.encoder = encode
	return c
}

// WithResponseDecoder 设置响应解码器
func (c *Config) WithResponseDecoder(decoder DecodeResponseFunc) *Config {
	c.decoder = decoder
	return c
}

// WithErrorDecoder 设置错误解码器
func (c *Config) WithErrorDecoder(errorFunc DecodeErrorFunc) *Config {
	c.errorDecoder = errorFunc
	return c
}

// WithNodeFilter 设置节点过滤器
func (c *Config) WithNodeFilter(filter ...selector.NodeFilter) *Config {
	c.nodeFilters = filter
	return c
}

// WithDiscovery 设置客户端服务发现
func (c *Config) WithDiscovery(d registry.Registry) *Config {
	c.discovery = d
	return c
}

// WithContext 设置上下文
func (c *Config) WithContext(ctx context.Context) *Config {
	c.ctx = ctx
	return c
}

// WithTLSConfig 设置tls信息
func (c *Config) WithTLSConfig(tls *tls.Config) *Config {
	c.TlsConf = tls
	return c
}
