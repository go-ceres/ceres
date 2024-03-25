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

package token

import (
	"github.com/go-ceres/ceres/pkg/common/config"
	"github.com/go-ceres/ceres/pkg/common/logger"
)

const ModName = "auth.token"

type Option func(o *Options)

// Options 配置信息
type Options struct {
	LogicType       string         `json:"logicType"`       // 登录逻辑类型（可分前台、后台、等等）
	TokenName       string         `json:"tokenName"`       // token名称
	Timeout         int64          `json:"timeout"`         // 共享session的过期时间
	ActivityTimeout int64          `json:"activityTimeout"` // 临时过期时间，即在误操作情况下，该token过期时间
	IsConcurrent    bool           `json:"isConcurrent"`    // 是否支持多账号同时登录
	IsShare         bool           `json:"isShare"`         // 是否可以共享token
	TokenStyle      Style          `json:"tokenStyle"`      // token样式
	AutoRenew       bool           `json:"autoRenew"`       // 自动续签
	TokenPrefix     string         `json:"tokenPrefix"`     // token前缀
	IsLog           bool           `json:"isLog"`           // 是否打印日志
	CheckLogin      bool           `json:"checkLogin"`      // 是否检查登录
	tokenBuilder    Builder        // token生成类
	storage         Storage        // 数据存储接口
	permission      Permission     // 权限管理接口
	listener        Listener       // 监听器
	logger          *logger.Logger // 日志打印组件
}

func DefaultOptions() *Options {
	conf := &Options{
		TokenName:       "Authorization",
		Timeout:         2592000,
		ActivityTimeout: -1,
		IsConcurrent:    true,
		IsShare:         true,
		TokenStyle:      StyleUuid,
		AutoRenew:       true,
		TokenPrefix:     "Bearer ",
		IsLog:           true,
		CheckLogin:      true,
		logger:          logger.With(logger.FieldMod(ModName)),
	}
	return conf
}

func ScanRawConfig(key string) *Options {
	conf := DefaultOptions()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

func ScanConfig(name ...string) *Options {
	key := "application.auth.token"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

func WithLogicType(LogicType string) Option {
	return func(o *Options) {
		o.LogicType = LogicType
	}
}

func WithTokenName(TokenName string) Option {
	return func(o *Options) {
		o.TokenName = TokenName
	}
}

func WithTimeout(Timeout int64) Option {
	return func(o *Options) {
		o.Timeout = Timeout
	}
}

func WithActivityTimeout(ActivityTimeout int64) Option {
	return func(o *Options) {
		o.ActivityTimeout = ActivityTimeout
	}
}

func WithIsConcurrent(IsConcurrent bool) Option {
	return func(o *Options) {
		o.IsConcurrent = IsConcurrent
	}
}

func WithIsShare(IsShare bool) Option {
	return func(o *Options) {
		o.IsShare = IsShare
	}
}

func WithTokenStyle(TokenStyle Style) Option {
	return func(o *Options) {
		o.TokenStyle = TokenStyle
	}
}

func WithAutoRenew(AutoRenew bool) Option {
	return func(o *Options) {
		o.AutoRenew = AutoRenew
	}
}

func WithTokenPrefix(TokenPrefix string) Option {
	return func(o *Options) {
		o.TokenPrefix = TokenPrefix
	}
}

func WithIsLog(IsLog bool) Option {
	return func(o *Options) {
		o.IsLog = IsLog
	}
}

func WithCheckLogin(CheckLogin bool) Option {
	return func(o *Options) {
		o.CheckLogin = CheckLogin
	}
}

func WithStorage(storage Storage) Option {
	return func(o *Options) {
		o.storage = storage
	}
}

func WithPermission(permission Permission) Option {
	return func(o *Options) {
		o.permission = permission
	}
}

func WithTokenBuilder(tokenBuilder Builder) Option {
	return func(o *Options) {
		o.tokenBuilder = tokenBuilder
	}
}

func WithListener(listener Listener) Option {
	return func(o *Options) {
		o.listener = listener
	}
}

func WithLogger(logger *logger.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

func (o *Options) WithOptions(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *Options) Build() Logic {
	return NewLogic(o)
}

type loginOptions struct {
	device  string // 登录的客户端设备标识
	timeout int64  // 指定当前登录token的有效时间（如果没有指定则使用全局配置）
}
type LoginOption func(o *loginOptions)

// DefaultOption 默认的登录额外参数
func defaultLoginOptions(c *Options) *loginOptions {
	opts := &loginOptions{
		timeout: c.Timeout,
		device:  defaultLoginDevice,
	}
	return opts
}

// LoginDevice 客户端
func LoginDevice(device string) LoginOption {
	return func(o *loginOptions) {
		o.device = device
	}
}

// LoginTimeout 当前此次登录的过期时间
func LoginTimeout(timeout int64) LoginOption {
	return func(o *loginOptions) {
		o.timeout = timeout
	}
}
