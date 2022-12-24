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

package gorm

import (
	"github.com/go-ceres/ceres/pkg/common/config"
	"github.com/go-ceres/ceres/pkg/common/logger"
	"gorm.io/gorm"
	"time"
)

const ModName = "store.gorm"

type Option func(o *Options)

type Options struct {
	Debug           bool           `json:"debug"`           // 是否开启调试
	DNS             string         `json:"dns"`             // 连接字符串
	MaxIdleConns    int            `json:"maxIdleConns"`    // 最大空闲连接数
	MaxOpenConns    int            `json:"maxOpenConns"`    // 最大活动连接数
	ConnMaxLifetime time.Duration  `json:"connMaxLifetime"` // 连接的最大存活时间
	MaxIdleTime     time.Duration  `json:"maxIdleTime"`     // 连接可最大空闲时间
	LogConfig       *LogConfig     `json:"logConfig"`       // 日志配置
	gormConfig      *gorm.Config   // gorm配置
	dialect         Dialector      // 驱动连接器
	logger          *logger.Logger // 日志库
	// 下面配置来自于gorm的配置，详情可查看gorm官方文档
	SkipDefaultTransaction                   bool `json:"skipDefaultTransaction"`
	FullSaveAssociations                     bool `json:"fullSaveAssociations"`
	DryRun                                   bool `json:"dryRun"`
	PrepareStmt                              bool `json:"prepareStmt"`
	DisableAutomaticPing                     bool `json:"disableAutomaticPing"`
	DisableForeignKeyConstraintWhenMigrating bool `json:"disableForeignKeyConstraintWhenMigrating"`
	DisableNestedTransaction                 bool `json:"disableNestedTransaction"`
	AllowGlobalUpdate                        bool `json:"allowGlobalUpdate"`
	QueryFields                              bool `json:"queryFields"`
	CreateBatchSize                          int  `json:"createBatchSize"`
	NowFunc                                  func() time.Time
}

func WithDebug(debug bool) Option {
	return func(o *Options) {
		o.Debug = debug
	}
}
func WithDns(dns string) Option {
	return func(o *Options) {
		o.DNS = dns
	}
}
func WithMaxIdleConns(mxIdleConns int) Option {
	return func(o *Options) {
		o.MaxIdleConns = mxIdleConns
	}
}
func WithMaxOpenConns(maxOpenConns int) Option {
	return func(o *Options) {
		o.MaxOpenConns = maxOpenConns
	}
}
func WithConnMaxIdleTime(timeout time.Duration) Option {
	return func(o *Options) {
		o.MaxIdleTime = timeout
	}
}
func WithConnMaxLifetime(timeout time.Duration) Option {
	return func(o *Options) {
		o.ConnMaxLifetime = timeout
	}
}
func Withdialect(dialect gorm.Dialector) Option {
	return func(o *Options) {
		o.dialect = dialect
	}
}
func WithSkipDefaultTransaction(disableNestedTransaction bool) Option {
	return func(o *Options) {
		o.DisableNestedTransaction = disableNestedTransaction
	}
}
func WithFullSaveAssociations(fullSaveAssociations bool) Option {
	return func(o *Options) {
		o.FullSaveAssociations = fullSaveAssociations
	}
}
func WithDryRun(dryRun bool) Option {
	return func(o *Options) {
		o.DryRun = dryRun
	}
}
func WithPrepareStmt(prepareStmt bool) Option {
	return func(o *Options) {
		o.PrepareStmt = prepareStmt
	}
}
func WithDisableAutomaticPing(disableAutomaticPing bool) Option {
	return func(o *Options) {
		o.DisableAutomaticPing = disableAutomaticPing
	}
}
func WithDisableForeignKeyConstraintWhenMigrating(disableForeignKeyConstraintWhenMigrating bool) Option {
	return func(o *Options) {
		o.DisableForeignKeyConstraintWhenMigrating = disableForeignKeyConstraintWhenMigrating
	}
}
func WithDisableNestedTransaction(disableNestedTransaction bool) Option {
	return func(o *Options) {
		o.DisableNestedTransaction = disableNestedTransaction
	}
}
func WithAllowGlobalUpdate(allowGlobalUpdate bool) Option {
	return func(o *Options) {
		o.AllowGlobalUpdate = allowGlobalUpdate
	}
}
func WithQueryFields(queryFields bool) Option {
	return func(o *Options) {
		o.QueryFields = queryFields
	}
}
func WithCreateBatchSize(size int) Option {
	return func(o *Options) {
		o.CreateBatchSize = size
	}
}

// WithNowFunc 设置当前时间生成方法
func WithNowFunc(nowFunc func() time.Time) Option {
	return func(o *Options) {
		o.NowFunc = nowFunc
	}
}

// WithLogger 设置日志
func WithLogger(log *logger.Logger) Option {
	return func(o *Options) {
		o.logger = log
	}
}

// DefaultOptions 默认参数
func DefaultOptions() *Options {
	return &Options{
		DNS:             "",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		LogConfig:       defaultLogConfig(),
		logger:          logger.With(logger.FieldMod(ModName)),
	}
}

// ScanRawConfig 扫描配置文件
func ScanRawConfig(key string) *Options {
	options := DefaultOptions()
	if err := config.Get(key).Scan(options); err != nil {
		panic(err)
	}
	return options
}

// ScanConfig 标准名称扫描
func ScanConfig(name ...string) *Options {
	key := "application.store.gorm"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

// WithOptions 设置额外参数
func (o *Options) WithOptions(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Build 构建DB
func (o *Options) Build() *DB {
	return NewWithOptions(o)
}
