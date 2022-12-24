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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	DB        = gorm.DB
	Dialector = gorm.Dialector
)

// New 创建gorm
func New(opts ...Option) *DB {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return NewWithOptions(options)
}

// NewWithOptions ...创建orm
func NewWithOptions(options *Options) *DB {
	// 连接驱动
	dialector := options.dialect
	if dialector == nil {
		dialector = mysql.Open(options.DNS)
	}
	// 构建日志
	gormConfig := buildGormConfig(options)
	// 创建gorm
	inner, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		panic(err)
	}
	sqlDb, err := inner.DB()
	if err != nil {
		panic(err)
	}
	// 设置最大可空闲连接数
	sqlDb.SetMaxIdleConns(options.MaxIdleConns)
	// 设置最大打开连接数
	sqlDb.SetMaxOpenConns(options.MaxOpenConns)
	// 设置连接最大存活时间
	if options.ConnMaxLifetime != 0 {
		sqlDb.SetConnMaxLifetime(options.ConnMaxLifetime)
	}
	// 设置连接最大可空闲时间
	if options.MaxIdleTime != 0 {
		sqlDb.SetConnMaxIdleTime(options.MaxIdleTime)
	}
	if err := sqlDb.Ping(); err != nil {
		panic(err)
	}
	return inner
}

// buildGormConfig 构建gorm配置文件
func buildGormConfig(options *Options) *gorm.Config {
	// gorm配置
	conf := &gorm.Config{
		SkipDefaultTransaction:                   options.SkipDefaultTransaction,
		FullSaveAssociations:                     options.FullSaveAssociations,
		DryRun:                                   options.DryRun,
		PrepareStmt:                              options.PrepareStmt,
		DisableAutomaticPing:                     options.DisableAutomaticPing,
		DisableForeignKeyConstraintWhenMigrating: options.DisableForeignKeyConstraintWhenMigrating,
		DisableNestedTransaction:                 options.DisableNestedTransaction,
		AllowGlobalUpdate:                        options.AllowGlobalUpdate,
		QueryFields:                              options.QueryFields,
		CreateBatchSize:                          options.CreateBatchSize,
		NowFunc:                                  options.NowFunc,
	}
	// 设置日志
	conf.Logger = newLog(options.logger, options.LogConfig)

	return conf
}
