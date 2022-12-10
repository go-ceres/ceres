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

package fiber

import (
	"github.com/go-ceres/ceres/config"
	"github.com/go-ceres/ceres/internal/matcher"
	"github.com/go-ceres/ceres/logger"
	"time"
)

const ModName = "server.fiber"

type TlsConfig struct {
	CertFile string `json:"certFile"` // tls配置
	KeyFile  string `json:"keyFile"`  // tls配置
}

type Config struct {
	Network    string            `json:"network"` // 网络类型
	Address    string            `json:"address"` // 连接地址
	Timeout    time.Duration     `json:"timeout"` // 超时时间
	TlsConf    *TlsConfig        `json:"tlsConf"` // tls配置
	middleware matcher.Matcher   // 中间件
	ene        EncodeErrorFunc   // 编码错误回调
	enReply    EncodeReplyFunc   // 响应体编码器
	decQuery   DecodeRequestFunc // 获取值
	decParams  DecodeRequestFunc
	logger     *logger.Helper //日志组件
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	conf := &Config{
		Network:    "tcp",
		Address:    "0.0.0.0:5200",
		ene:        DefaultErrorFunc,
		enReply:    DefaultRequestDecoder,
		decQuery:   DefaultRequestQuery,
		decParams:  DefaultRequestParams,
		middleware: matcher.New(),
		logger:     logger.With(logger.FieldMod(ModName)),
	}
	return conf
}

// ScanRawConfig 原始扫描
func ScanRawConfig(key string) *Config {
	conf := DefaultConfig()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 标准扫描
func ScanConfig(name ...string) *Config {
	key := "ceres.application.server.fiber"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

// SetNetwork 设置网络类型
func (c *Config) SetNetwork(Network string) *Config {
	c.Network = Network
	return c
}

// SetAddress 设置服务器地址
func (c *Config) SetAddress(Address string) *Config {
	c.Address = Address
	return c
}

// SetTimeout 设置超时
func (c *Config) SetTimeout(Timeout time.Duration) *Config {
	c.Timeout = Timeout
	return c
}

// SetTlsConf 设置tls认证
func (c *Config) SetTlsConf(TlsConf *TlsConfig) *Config {
	c.TlsConf = TlsConf
	return c
}

// SetEne 设置错误处理函数
func (c *Config) SetEne(ene EncodeErrorFunc) *Config {
	c.ene = ene
	return c
}

// SetEnc 设置响应编码器
func (c *Config) SetEnc(enReply EncodeReplyFunc) *Config {
	c.enReply = enReply
	return c
}

// SetDecQuery 设置请求解码器
func (c *Config) SetDecQuery(decQuery DecodeRequestFunc) *Config {
	c.decQuery = decQuery
	return c
}

// SetDecParams 设置参数解码器
func (c *Config) SetDecParams(decParams DecodeRequestFunc) *Config {
	c.decParams = decParams
	return c
}

// SetLogger 设置日志
func (c *Config) SetLogger(logger *logger.Helper) *Config {
	c.logger = logger
	return c
}

// Build 构建fiber服务
func (c *Config) Build() *Server {
	return New(c)
}
