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
	"crypto/tls"
	"github.com/go-ceres/ceres/pkg/common/config"
	"github.com/go-ceres/ceres/pkg/common/logger"
	"github.com/go-redis/redis"
	"time"
)

const ModName = "client.redis"

type Mode string

type Option func(o *Options)

// Options 配置参数
type Options struct {
	Addrs              []string      `json:"addrs"`        // 连接地址
	Password           string        `json:"password"`     // 密码
	DB                 int           `json:"db"`           // DB，默认为0, 一般应用不推荐使用DB分片
	PoolSize           int           `json:"pool_size"`    // 集群内每个节点的最大连接池限制 默认每个CPU10个连接
	MaxRetries         int           `json:"maxRetries"`   //网络相关的错误最大重试次数 默认5次
	MinIdleConns       int           `json:"minIdleConns"` // 最小空闲连接数,默认100
	DialTimeout        time.Duration `json:"dialTimeout"`  // 连接超时
	ReadTimeout        time.Duration `json:"readTimeout"`  //读取超时 默认3s
	WriteTimeout       time.Duration `json:"writeTimeout"` // 写入超时 默认3s
	IdleTimeout        time.Duration `json:"idleTimeout"`  // 连接最大空闲时间，默认60s, 超过该时间，连接会被主动关闭
	Debug              bool          `json:"debug"`        // 是否开启debug模式
	ReadOnly           bool          `json:"readOnly"`     // 集群模式中在从属节点上启用读模式
	MinRetryBackoff    time.Duration `json:"minRetryBackoff"`
	MaxRetryBackoff    time.Duration `json:"maxRetryBackoff"`
	MaxConnAge         time.Duration `json:"maxConnAge"`
	PoolTimeout        time.Duration `json:"poolTimeout"`
	IdleCheckFrequency time.Duration `json:"idleCheckFrequency"`
	MaxRedirects       int           `json:"maxRedirects"`
	RouteByLatency     bool          `json:"routeByLatency"`
	RouteRandomly      bool          `json:"routeRandomly"`
	MasterName         string        `json:"masterName"`
	OnConnect          func(*redis.Conn) error
	TLSConfig          *tls.Config
	logger             *logger.Logger // 日志组件
}

func WithAddrs(addrs []string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}
func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}
func WithDB(db int) Option {
	return func(o *Options) {
		o.DB = db
	}
}
func WithPoolSize(poolSize int) Option {
	return func(o *Options) {
		o.PoolSize = poolSize
	}
}
func WithMaxRetries(maxRetries int) Option {
	return func(o *Options) {
		o.MaxRetries = maxRetries
	}
}
func WithMinIdleConns(minIdleConns int) Option {
	return func(o *Options) {
		o.MinIdleConns = minIdleConns
	}
}
func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(o *Options) {
		o.DialTimeout = dialTimeout
	}
}
func WithReadTimeout(readTimeout time.Duration) Option {
	return func(o *Options) {
		o.ReadTimeout = readTimeout
	}
}
func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(o *Options) {
		o.WriteTimeout = writeTimeout
	}
}
func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(o *Options) {
		o.IdleTimeout = idleTimeout
	}
}
func WithDebug(debug bool) Option {
	return func(o *Options) {
		o.Debug = debug
	}
}
func WithReadOnly(readOnly bool) Option {
	return func(o *Options) {
		o.ReadOnly = readOnly
	}
}
func WithMinRetryBackoff(minRetryBackoff time.Duration) Option {
	return func(o *Options) {
		o.MinRetryBackoff = minRetryBackoff
	}
}
func WithMaxRetryBackoff(maxRetryBackoff time.Duration) Option {
	return func(o *Options) {
		o.MaxRetryBackoff = maxRetryBackoff
	}
}
func WithMaxConnAge(maxConnAge time.Duration) Option {
	return func(o *Options) {
		o.MaxConnAge = maxConnAge
	}
}
func WithPoolTimeout(poolTimeout time.Duration) Option {
	return func(o *Options) {
		o.PoolTimeout = poolTimeout
	}
}
func WithIdleCheckFrequency(idleCheckFrequency time.Duration) Option {
	return func(o *Options) {
		o.IdleCheckFrequency = idleCheckFrequency
	}
}
func WithMaxRedirects(maxRedirects int) Option {
	return func(o *Options) {
		o.MaxRedirects = maxRedirects
	}
}
func WithRouteByLatency(routeByLatency bool) Option {
	return func(o *Options) {
		o.RouteByLatency = routeByLatency
	}
}
func WithRouteRandomly(routeRandomly bool) Option {
	return func(o *Options) {
		o.RouteRandomly = routeRandomly
	}
}
func WithMasterName(masterName string) Option {
	return func(o *Options) {
		o.MasterName = masterName

	}
}
func WithOnConnect(onConnect func(*redis.Conn) error) Option {
	return func(o *Options) {
		o.OnConnect = onConnect
	}
}
func WithTLSConfig(tLSConfig *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = tLSConfig
	}
}
func WithLogger(logger *logger.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

// DefaultOptions 默认配置参数
func DefaultOptions() *Options {
	return &Options{
		Addrs:        []string{"127.0.0.1:6379"},
		DB:           0,
		PoolSize:     10,
		MaxRetries:   5,
		MinIdleConns: 100,
		DialTimeout:  time.Second * 3,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
		IdleTimeout:  time.Second * 60,
		Debug:        false,
		ReadOnly:     false,
		logger:       logger.With(logger.FieldMod(ModName)),
	}
}

// ScanRawConfig 扫描配置文件
func ScanRawConfig(key string) *Options {
	conf := DefaultOptions()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 标准扫描
func ScanConfig(name ...string) *Options {
	key := "application.client.redis"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

// WithOptions 配置参数
func (o *Options) WithOptions(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Build 构建redis
func (o *Options) Build() *Client {
	return NewWithOptions(o)
}
