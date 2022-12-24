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

package app

import (
	"context"
	"github.com/go-ceres/ceres"
	"github.com/go-ceres/ceres/pkg/common/config"
	"github.com/go-ceres/ceres/pkg/common/logger"
	"github.com/go-ceres/ceres/pkg/transport"
	"net/url"
	"time"
)

// HookType 定义 hook类型
type HookType uint

// HookFunc 钩子回调方法
type HookFunc func(ctx context.Context)

const (
	ModName              = "app"
	BeforeStart HookType = iota
	BeforeStop
	AfterStart
	AfterStop
)

type Option func(o *Options)

// Options 配置信息
type Options struct {
	ctx              context.Context         // 应用上下文
	ID               string                  `json:"id"`               // 应用唯一标识
	Name             string                  `json:"name"`             // 应用名称
	Version          string                  `json:"version"`          // 应用版本
	Metadata         map[string]string       `json:"metadata"`         // 附加信息
	Endpoints        []*url.URL              `json:"endpoints"`        // 服务地址
	Region           string                  `json:"region"`           // 服务所属地域
	Zone             string                  `json:"zone"`             // 服务所属分区
	HideBanner       bool                    `json:"hideBanner"`       // 隐藏打印横幅
	MaxProc          int64                   `json:"maxProc"`          // 处理器内核优化
	RegistrarTimeout time.Duration           `json:"registrarTimeout"` // 服务注册超时时间
	StopTimeout      time.Duration           `json:"stopTimeout"`      // 注销服务超时时间
	hooks            map[HookType][]HookFunc // 启动钩子
	transports       []transport.Transport   // 服务集合
	registry         transport.Registry      // 注册中心
	logger           logger.Logger           // 日志组件
}

// DefaultOptions 默认配置
func DefaultOptions() *Options {
	return &Options{
		ctx:              context.Background(),
		ID:               ceres.AppId(),
		Name:             ceres.AppName(),
		Version:          ceres.AppVersion(),
		Region:           ceres.AppRegion(),
		Zone:             ceres.AppZone(),
		HideBanner:       false,
		RegistrarTimeout: 10 * time.Second,
		StopTimeout:      3 * time.Second,
		hooks:            map[HookType][]HookFunc{},
		transports:       []transport.Transport{},
	}
}

// ScanRawConfig 扫描原始键
func ScanRawConfig(key string) *Options {
	conf := DefaultOptions()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 标准扫描
func ScanConfig(name ...string) *Options {
	key := "application"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

// WithContext 设置应用上下文
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.ctx = ctx
	}
}

// WithId 设置应用id
func WithId(id string) Option {
	return func(o *Options) {
		o.ID = id
	}
}

// WithName 设置应用名称
func WithName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// WithVersion 设置应用版本
func WithVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

// WithMetadata 应用元信息
func WithMetadata(metadata map[string]string) Option {
	return func(o *Options) {
		o.Metadata = metadata
	}
}

// WithRegion 设置应用部属地域
func WithRegion(region string) Option {
	return func(o *Options) {
		o.Region = region
	}
}

// WithZone 设置应用部属分区
func WithZone(zone string) Option {
	return func(o *Options) {
		o.Zone = zone
	}
}

// HideBanner 隐藏应用横幅打印
func HideBanner() Option {
	return func(o *Options) {
		o.HideBanner = true
	}
}

// WithRegistrarTimeout 设置注册服务超时时间
func WithRegistrarTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.RegistrarTimeout = timeout
	}
}

// WithStopTimeout 设置停止服务超时时间
func WithStopTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.StopTimeout = timeout
	}
}

// WithHooks 设置钩子
func WithHooks(hookType HookType, hooks ...HookFunc) Option {
	return func(o *Options) {
		if o.hooks[hookType] == nil {
			o.hooks[hookType] = make([]HookFunc, 0)
		}
		o.hooks[hookType] = hooks
	}
}

// AddHooks 添加钩子
func AddHooks(hookType HookType, hooks ...HookFunc) Option {
	return func(o *Options) {
		if o.hooks[hookType] == nil {
			o.hooks[hookType] = make([]HookFunc, 0)
		}
		o.hooks[hookType] = append(o.hooks[hookType], hooks...)
	}
}

// WithRegistry 设置注册中心
func WithRegistry(reg transport.Registry) Option {
	return func(o *Options) {
		o.registry = reg
	}
}

// WithTransport 设置传输对象
func WithTransport(transports ...transport.Transport) Option {
	return func(o *Options) {
		o.transports = transports
	}
}

// WithLogger 设置日志
func WithLogger(log logger.Logger) Option {
	return func(o *Options) {
		o.logger = log
	}
}

// WithMaxProc 设置cpu配额
func WithMaxProc(maxProc int64) Option {
	return func(o *Options) {
		o.MaxProc = maxProc
	}
}

// WithOption 设置参数
func (o *Options) WithOption(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Build 构建应用
func (o *Options) Build() *Application {
	return NewWithOptions(o)
}
