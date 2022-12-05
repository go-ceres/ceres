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
	"github.com/go-ceres/ceres/config"
	"github.com/go-ceres/ceres/logger"
	"gorm.io/gorm"
	log "gorm.io/gorm/logger"
	"time"
)

// Config 配置信息
type Config struct {
	Drive           string        `json:"drive"`           // 驱动
	DNS             string        `json:"dns"`             // 连接字符串
	Debug           bool          `json:"debug"`           // 是否开启调试
	MaxIdleConns    int           `json:"maxIdleConns"`    // 最大空闲连接数
	MaxOpenConns    int           `json:"maxOpenConns"`    // 最大活动连接数
	ConnMaxLifetime time.Duration `json:"connMaxLifetime"` // 连接的最大存活时间

	*GormConfig                // gorm初始化配置
	*LogConfig                 // 日志配置
	dialect     Dialector      // 驱动连接器
	logger      *logger.Helper // 日志库
}

// LogConfig 日志配置
type LogConfig struct {
	SlowThreshold time.Duration // 日志时间阈值
	Colorful      bool          // 是否开启日志颜色区别
	LogLevel      string        // 日志等级
}

// DefaultLogConfig 默认的日志配置
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		SlowThreshold: time.Second,
		Colorful:      false,
		LogLevel:      "",
	}
}

type GormConfig gorm.Config

// DefaultConfig 默认gorm配置
func DefaultConfig() *Config {
	return &Config{
		Drive:           "mysql",
		DNS:             "",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		GormConfig:      &GormConfig{},
		LogConfig:       DefaultLogConfig(),
	}
}

// ScanRawConfig 扫描未包装的key
func ScanRawConfig(key string) *Config {
	conf := DefaultConfig()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 标准名称扫描
func ScanConfig(name ...string) *Config {
	key := "ceres.application.store.gorm"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

// initLogger 初始化日志
func (c *Config) initLogger() {
	// 默认日志配置
	logConf := log.Config{
		SlowThreshold: time.Second, // 慢 SQL 阈值
		LogLevel:      log.Silent,  // Log level
		Colorful:      false,       // 禁用彩色打印
	}
	// 转换等级
	if c.LogLevel != "" {
		logConf.LogLevel = ConvertLevel(c.LogLevel)
	}
	logConf.Colorful = c.Colorful
	logConf.SlowThreshold = c.SlowThreshold
	dbLog := newLog(c.logger, logConf)
	if c.Debug {
		dbLog = dbLog.LogMode(log.Info)
	}
	// gorm的配置信息
	c.GormConfig.Logger = dbLog
}

// WithDialector 单独设置Dialector
func (c *Config) WithDialector(dialect Dialector) *Config {
	c.dialect = dialect
	return c
}

// Build 构建gorm数据库链接
func (c *Config) Build() *DB {
	// 初始化日志
	c.initLogger()
	// 创建驱动
	if driver, ok := drivers[c.Drive]; !ok {
		logger.Panicf("%s driver is not set", c.Drive)
	} else {
		c.dialect = driver(c.DNS)
	}
	// 数据库
	db, err := Open(c.dialect, c)
	if err != nil {
		logger.Panicf("open gorm", "err", err, "value", c)
	}
	return db
}
