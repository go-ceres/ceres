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
	"github.com/go-ceres/ceres/config"
	"github.com/go-ceres/ceres/logger"
)

const ModName = "auth.token"

// Config 配置信息
type Config struct {
	LogicType       string         `json:"logicType"`       // 登录逻辑类型（可分前台、后台、等等）
	TokenName       string         `json:"tokenName"`       // token名称
	Timeout         int64          `json:"timeout"`         // 共享session的过期时间
	ActivityTimeout int64          `json:"activityTimeout"` // 临时过期时间，即在误操作情况下，该token过期时间
	IsConcurrent    bool           `json:"isConcurrent"`    // 是否支持多账号同时登录
	IsShare         bool           `json:"isShare"`         // 是否可以共享token
	TokenStyle      TokenStyle     `json:"tokenStyle"`      // token样式
	AutoRenew       bool           `json:"autoRenew"`       // 自动续签
	TokenPrefix     string         `json:"tokenPrefix"`     // token前缀
	IsLog           bool           `json:"isLog"`           // 是否打印日志
	CheckLogin      bool           `json:"checkLogin"`      // 是否检查登录
	tokenBuilder    TokenBuilder   // token生成类
	storage         Storage        // 数据存储接口
	permission      Permission     // 权限管理接口
	listener        Listener       // 监听器
	logger          *logger.Helper // 日志打印组件
}

func DefaultConfig() *Config {
	conf := &Config{
		TokenName:       "ceres-token",
		Timeout:         2592000,
		ActivityTimeout: -1,
		IsConcurrent:    true,
		IsShare:         true,
		TokenStyle:      TokenStyleUuid,
		AutoRenew:       true,
		TokenPrefix:     "Bearer ",
		IsLog:           true,
		CheckLogin:      true,
		logger:          logger.With(logger.FieldMod(ModName)),
	}
	return conf
}

func ScanRawConfig(key string) *Config {
	conf := DefaultConfig()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

func ScanConfig(name ...string) *Config {
	key := "ceres.application.auth.token"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

func (c *Config) SetLogicType(LogicType string) *Config {
	c.LogicType = LogicType
	return c
}

func (c *Config) SetTokenName(TokenName string) *Config {
	c.TokenName = TokenName
	return c
}

func (c *Config) SetTimeout(Timeout int64) *Config {
	c.Timeout = Timeout
	return c
}

func (c *Config) SetActivityTimeout(ActivityTimeout int64) *Config {
	c.ActivityTimeout = ActivityTimeout
	return c
}

func (c *Config) SetIsConcurrent(IsConcurrent bool) *Config {
	c.IsConcurrent = IsConcurrent
	return c
}

func (c *Config) SetIsShare(IsShare bool) *Config {
	c.IsShare = IsShare
	return c
}

func (c *Config) SetTokenStyle(TokenStyle TokenStyle) *Config {
	c.TokenStyle = TokenStyle
	return c
}

func (c *Config) SetAutoRenew(AutoRenew bool) *Config {
	c.AutoRenew = AutoRenew
	return c
}

func (c *Config) SetTokenPrefix(TokenPrefix string) *Config {
	c.TokenPrefix = TokenPrefix
	return c
}

func (c *Config) SetIsLog(IsLog bool) *Config {
	c.IsLog = IsLog
	return c
}

func (c *Config) SetCheckLogin(CheckLogin bool) *Config {
	c.CheckLogin = CheckLogin
	return c
}

func (c *Config) SetStorage(storage Storage) *Config {
	c.storage = storage
	return c
}

func (c *Config) SetPermission(permission Permission) *Config {
	c.permission = permission
	return c
}

func (c *Config) SetTokenBuilder(tokenBuilder TokenBuilder) *Config {
	c.tokenBuilder = tokenBuilder
	return c
}

func (c *Config) SetListener(listener Listener) *Config {
	c.listener = listener
	return c
}

func (c *Config) SetLogger(logger *logger.Helper) *Config {
	c.logger = logger
	return c
}

func (c *Config) Build() *Logic {
	return NewLogic(c)
}
