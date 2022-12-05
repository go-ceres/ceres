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
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

// ZapConfig zap配置
type ZapConfig struct {
	Dir           string        // Dir 日志输出目录
	Debug         bool          // 是否是调试模式
	Name          string        // Name 日志文件名称
	Level         string        // Level 日志初始等级
	Fields        []zap.Field   // 日志初始化字段
	AddCaller     bool          // 是否添加调用者信息
	Prefix        string        // 日志前缀
	MaxSize       int           // 日志输出文件最大长度，超过改值则截断
	MaxAge        int           //日志最大存活时长
	MaxBackup     int           // 日志最大备份数量
	Interval      time.Duration // 日志磁盘刷盘间隔
	CallerSkip    int           // 调用者层级
	Async         bool          // 异步写入
	Queue         bool
	QueueSleep    time.Duration
	Core          zapcore.Core
	EncoderConfig *zapcore.EncoderConfig
	configKey     string // 日志等级配置路径
}

// Filename 文件名
func (z *ZapConfig) Filename() string {
	return fmt.Sprintf("%s/%s", z.Dir, z.Name)
}

// DefaultZapConfig 默认的zap日志的配置
func DefaultZapConfig() *ZapConfig {
	return &ZapConfig{
		Name:          "default.log",
		Dir:           ".",
		Level:         "debug",
		MaxSize:       500, // 500M
		MaxAge:        1,   // 1 day
		MaxBackup:     10,  // 10 backup
		Interval:      24 * time.Hour,
		CallerSkip:    1,
		AddCaller:     false,
		Async:         false,
		Queue:         false,
		Debug:         true,
		QueueSleep:    100 * time.Millisecond,
		EncoderConfig: defaultZapEncoderConfig(),
	}
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

// defaultZapEncoderConfig 默认的zap编码初始化配置
func defaultZapEncoderConfig() *zapcore.EncoderConfig {
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
