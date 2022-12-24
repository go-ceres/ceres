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

package logger

import (
	"github.com/fatih/color"
	"github.com/go-ceres/ceres/pkg/common/config"
	"github.com/go-ceres/ceres/pkg/common/logger/rotate"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

const (
	ModeFile = "file"
	ModeStd  = "std"
)

type Mode string

type Option func(o *Options)

type Options struct {
	Mode           Mode        `json:"mode"`       // 是否是调试模式
	Level          Level       `json:"level"`      // Level 日志初始等级
	AddCaller      bool        `json:"addCaller"`  // 是否添加调用者信息
	CallerSkip     int         `json:"callerSkip"` // 调用者层级
	Async          bool        `json:"async"`      // 异步写入
	Fields         []zap.Field // 日志初始化字段
	EncoderConfig  *zapcore.EncoderConfig
	Core           []zapcore.Core // 日志写入者,可添加多个
	*rotate.Config                // 日志轮转配置
	configKey      string         // 日志等级配置路径
}

// DefaultOptions 默认参数
func DefaultOptions() *Options {
	return &Options{
		Config:        rotate.DefaultConfig(),
		Level:         DebugLevel,
		CallerSkip:    1,
		AddCaller:     false,
		Async:         false,
		Mode:          ModeStd,
		EncoderConfig: defaultEncoderConfig(),
	}
}

// ScanRawConfig 扫描配置
func ScanRawConfig(key string) *Options {
	options := DefaultOptions()
	if err := config.Get(key).Scan(options); err != nil {
		panic(err)
	}
	options.configKey = key
	return options
}

// ScanConfig 扫描配置
func ScanConfig(name ...string) *Options {
	key := "application.logger"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

// WithMod 设置日志模式
func WithMod(mode Mode) Option {
	return func(o *Options) {
		o.Mode = mode
	}
}

// WithLevel 日志等级
func WithLevel(level Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// WithAddCaller 添加调用者信息
func WithAddCaller(addCaller bool) Option {
	return func(o *Options) {
		o.AddCaller = addCaller
	}
}

// WithCallerSkip 调用者层级
func WithCallerSkip(skip int) Option {
	return func(o *Options) {
		o.CallerSkip = skip
	}
}

// WithFields 初始化字段
func WithFields(fields ...Field) Option {
	return func(o *Options) {
		o.Fields = fields
	}
}

// WithRotateConfig 日志轮转配置
func WithRotateConfig(conf *rotate.Config) Option {
	return func(o *Options) {
		o.Config = conf
	}
}

// WithEncoderConfig 编码配置
func WithEncoderConfig(encoderConfig *zapcore.EncoderConfig) Option {
	return func(o *Options) {
		o.EncoderConfig = encoderConfig
	}
}

// WithZapCore 添加额外的日志写入者
func WithZapCore(cores ...zapcore.Core) Option {
	return func(o *Options) {
		o.Core = cores
	}
}

// WithOptions 手动设置参数
func (o *Options) WithOptions(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Build 构建日志组件
func (o *Options) Build() *Logger {
	l := NewWithOptions(o)
	if o.configKey != "" {
		_ = l.AutoLevel(o.configKey + ".level")
	}
	return l
}

// defaultLevelEncodeLevel 默认等级编码配置
func defaultLevelEncodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorize = color.New()
	switch lv {
	case zapcore.DebugLevel:
		colorize.Add(color.FgBlue)
	case zapcore.InfoLevel:
		colorize.Add(color.FgGreen)
	case zapcore.WarnLevel:
		colorize.Add(color.FgYellow)
	case zapcore.ErrorLevel, zap.PanicLevel, zap.DPanicLevel, zap.FatalLevel:
		colorize.Add(color.FgRed)
	default:
	}
	enc.AppendString(colorize.Sprint(lv.String()))
}

// defaultEncoderConfig 默认的zap编码初始化配置
func defaultEncoderConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:          "ts",
		LevelKey:         "lv",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stack",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeTime:       timeEncoderUnix,
		ConsoleSeparator: "\t",
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
	}
}

// timeEncoder 默认的时间格式化函数
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func timeEncoderUnix(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}
