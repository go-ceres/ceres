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
	"github.com/go-ceres/ceres/pkg/common/config"
	"github.com/go-ceres/ceres/pkg/common/logger/rotate"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"strings"
)

// Logger 日志结构体
type Logger struct {
	desugar       *zap.Logger
	lv            *zap.AtomicLevel
	core          zapcore.Core
	options       *Options
	sugar         *zap.SugaredLogger
	encoderConfig *zapcore.EncoderConfig
}

// New 创建日志
func New(opts ...Option) *Logger {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return NewWithOptions(options)
}

// NewWithOptions 创建日志
func NewWithOptions(options *Options) *Logger {
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddStacktrace(zap.DPanicLevel))
	if options.AddCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(options.CallerSkip))
	}
	if len(options.Fields) > 0 {
		zapOptions = append(zapOptions, zap.Fields(options.Fields...))
	}

	var ws zapcore.WriteSyncer
	if options.Mode == ModeStd {
		ws = os.Stdout
	} else {
		ws = zapcore.AddSync(rotate.NewRotate(options.Config))
	}

	if options.Async {
		ws, _ = Buffer(ws, defaultBufferSize, defaultFlushInterval)
	}

	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := lv.UnmarshalText([]byte(options.Level)); err != nil {
		panic(err)
	}
	encoderConfig := *options.EncoderConfig
	cores := options.Core
	cores = append(cores, zapcore.NewCore(
		func() zapcore.Encoder {
			if options.Mode == ModeStd {
				encoderConfig.EncodeLevel = defaultLevelEncodeLevel
				encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
				encoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, arrayEncoder zapcore.PrimitiveArrayEncoder) {
					arrayEncoder.AppendString(caller.FullPath())
				}
				return zapcore.NewConsoleEncoder(encoderConfig)
			}
			return zapcore.NewJSONEncoder(encoderConfig)
		}(),
		ws,
		lv,
	))
	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(
		core,
		zapOptions...,
	)
	return &Logger{
		desugar: zapLogger,
		lv:      &lv,
		options: options,
		sugar:   zapLogger.Sugar(),
	}
}

// AutoLevel 自动配置日志等级
func (l *Logger) AutoLevel(key string) error {
	return config.Watch(key, func(s string, value config.Value) {
		levelStr, err := config.Get(key).String()
		if err == nil {
			lvText := strings.ToLower(levelStr)
			if lvText != "" {
				l.Info("update level", FieldString("level", lvText))
				_ = l.lv.UnmarshalText([]byte(lvText))
			}
		}
	})
}

// Sync 同步日志缓存
func (l *Logger) Sync() error {
	return l.desugar.Sync()
}

// StdLog ...
func (l *Logger) StdLog() *log.Logger {
	return zap.NewStdLog(l.desugar)
}

// With ...
func (l *Logger) With(fields ...Field) *Logger {
	desugarLogger := l.desugar.With(fields...)
	return &Logger{
		desugar: desugarLogger,
		lv:      l.lv,
		sugar:   desugarLogger.Sugar(),
		options: l.options,
	}
}

// Debug ...
func (l *Logger) Debug(msg string, fields ...Field) {
	if l.options.Mode == ModeStd {
		msg = normalizeMessage(msg)
	}
	l.desugar.Debug(msg, fields...)
}

// Info ...
func (l *Logger) Info(msg string, fields ...Field) {
	l.desugar.Info(msg, fields...)
}

// Warn ...
func (l *Logger) Warn(msg string, fields ...Field) {
	if l.options.Mode == ModeStd {
		msg = normalizeMessage(msg)
	}
	l.desugar.Warn(msg, fields...)
}

// Error ...
func (l *Logger) Error(msg string, fields ...Field) {
	if l.options.Mode == ModeStd {
		msg = normalizeMessage(msg)
	}
	l.desugar.Error(msg, fields...)
}

// Fatal ...
func (l *Logger) Fatal(msg string, fields ...Field) {
	if l.options.Mode == ModeStd {
		msg = normalizeMessage(msg)
	}
	l.desugar.Fatal(msg, fields...)
}

// Panic ...
func (l *Logger) Panic(msg string, fields ...Field) {
	if l.options.Mode == ModeStd {
		msg = normalizeMessage(msg)
	}
	l.desugar.Panic(msg, fields...)
}

// Debugf ...
func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.sugar.Debugf(sprintf(msg, args...))
}

// Infof ...
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.sugar.Infof(sprintf(msg, args...))
}

// Warnf ...
func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.sugar.Warnf(sprintf(msg, args...))
}

// Errorf ...
func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.sugar.Errorf(sprintf(msg, args...))
}

// Fatalf ...
func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.sugar.Fatalf(sprintf(msg, args...))
}

// Panicf ...
func (l *Logger) Panicf(msg string, args ...interface{}) {
	l.sugar.Panicf(sprintf(msg, args...))
}

func sprintf(template string, args ...interface{}) string {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	return msg
}

// Debugw ...
func (l *Logger) Debugw(msg string, args ...interface{}) {
	if l.options.Mode == ModeStd {
		msg = normalizeMessage(msg)
	}
	l.sugar.Debugw(msg, args...)
}

// Infow ...
func (l *Logger) Infow(msg string, keyAndValues ...interface{}) {
	if l.options.Mode == ModeStd {
		msg = normalizeMessage(msg)
	}
	l.sugar.Infow(msg, keyAndValues...)
}

// Warnw ...
func (l *Logger) Warnw(msg string, keyAndValues ...interface{}) {
	if l.options.Mode == ModeStd {
		msg = normalizeMessage(msg)
	}
	l.sugar.Warnw(msg, keyAndValues...)
}

// Errorw ...
func (l *Logger) Errorw(msg string, keyAndValues ...interface{}) {
	if l.options.Mode == ModeStd {
		msg = normalizeMessage(msg)
	}
	l.sugar.Errorw(msg, keyAndValues...)
}

// Fatalw ...
func (l *Logger) Fatalw(msg string, keyAndValues ...interface{}) {
	if l.options.Mode == ModeStd {
		msg = normalizeMessage(msg)
	}
	l.sugar.Fatalw(msg, keyAndValues...)
}

// Panicw ...
func (l *Logger) Panicw(msg string, keyAndValues ...interface{}) {
	if l.options.Mode == ModeStd {
		msg = normalizeMessage(msg)
	}
	l.sugar.Panicw(msg, keyAndValues...)
}

func normalizeMessage(msg string) string {
	return fmt.Sprintf("%-40s", msg)
}
