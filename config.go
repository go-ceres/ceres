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

package ceres

import (
	"github.com/go-ceres/ceres/config"
	"github.com/go-ceres/ceres/logger"
	"github.com/go-ceres/ceres/registry"
	"github.com/go-ceres/ceres/server"
	"github.com/go-ceres/ceres/version"
	"net/url"
	"time"
)

// HookType 定义 hook类型
type HookType uint

// HookFunc 钩子回调方法
type HookFunc func()

const (
	ModName              = "app"
	BeforeStart HookType = iota
	BeforeStop
	AfterStart
	AfterStop
)

// Config 配置信息
type Config struct {
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
	servers          []server.Server         // 服务集合
	registry         registry.Registry       // 注册中心
	logger           logger.Logger           // 日志组件
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		ID:               version.AppId(),
		Name:             version.AppName(),
		Version:          version.AppVersion(),
		Region:           version.AppRegion(),
		Zone:             version.AppZone(),
		HideBanner:       false,
		RegistrarTimeout: 10 * time.Second,
		StopTimeout:      3 * time.Second,
		hooks:            map[HookType][]HookFunc{},
		servers:          []server.Server{},
	}
}

// ScanRawConfig 扫描原始键
func ScanRawConfig(key string) *Config {
	conf := DefaultConfig()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 标准扫描
func ScanConfig(name ...string) *Config {
	key := "ceres.application"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

// SetID 设置应用id
func (c *Config) SetID(id string) *Config {
	c.ID = id
	return c
}

// SetName 设置应用名称
func (c *Config) SetName(name string) *Config {
	c.Name = name
	return c
}

// SetVersion 设置应用版本
func (c *Config) SetVersion(version string) *Config {
	c.Version = version
	return c
}

// SetMetadata 设置应用版本
func (c *Config) SetMetadata(metadata map[string]string) *Config {
	c.Metadata = metadata
	return c
}

// SetRegion 设置应用部属地域
func (c *Config) SetRegion(region string) *Config {
	c.Region = region
	return c
}

// SetZone 设置应用部属分区
func (c *Config) SetZone(zone string) *Config {
	c.Zone = zone
	return c
}

// SetHideBanner 隐藏应用横幅
func (c *Config) SetHideBanner() *Config {
	c.HideBanner = true
	return c
}

// SetRegistrarTimeout 设置注册服务超时时间
func (c *Config) SetRegistrarTimeout(timeout time.Duration) *Config {
	c.RegistrarTimeout = timeout
	return c
}

// SetStopTimeout 设置停止服务超时时间
func (c *Config) SetStopTimeout(timeout time.Duration) *Config {
	c.StopTimeout = timeout
	return c
}

// AddHooks 添加钩子
func (c *Config) AddHooks(hookType HookType, hooks ...HookFunc) *Config {
	if c.hooks[hookType] == nil {
		c.hooks[hookType] = make([]HookFunc, 0)
	}
	c.hooks[hookType] = append(c.hooks[hookType], hooks...)
	return c
}

// SetRegistry 设置注册中心
func (c *Config) SetRegistry(reg registry.Registry) *Config {
	c.registry = reg
	return c
}

// AddServers 添加服务
func (c *Config) AddServers(servers ...server.Server) *Config {
	c.servers = append(c.servers, servers...)
	return c
}

// Build 构建应用
func (c *Config) Build() *App {
	return New(c)
}
