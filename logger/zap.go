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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var _ Logger = (*Zap)(nil)

type Zap struct {
	desugar       *zap.Logger
	lv            *zap.AtomicLevel
	core          zapcore.Core
	config        *ZapConfig
	sugar         *zap.SugaredLogger
	encoderConfig *zapcore.EncoderConfig
}

// SetLevel 设置等级
func (z *Zap) SetLevel(level string) {
	l := ParseLevel(level)
	if l.String() != "" {
		_ = z.lv.UnmarshalText([]byte(l.String()))
	}
}

func normalizeMessage(msg string) string {
	return fmt.Sprintf("%-40s", msg)
}

func (z *Zap) Log(level Level, msg string, fields ...LogField) {
	if z.config.Debug {
		msg = normalizeMessage(msg)
	}
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	switch level {
	case DebugLevel:
		z.desugar.Debug(msg, zapFields...)
	case InfoLevel:
		z.desugar.Info(msg, zapFields...)
	case WarnLevel:
		z.desugar.Warn(msg, zapFields...)
	case ErrorLevel:
		z.desugar.Error(msg, zapFields...)
	case PanicLevel:
		z.desugar.Panic(msg, zapFields...)
	case FatalLevel:
		z.desugar.Fatal(msg, zapFields...)
	}
}

func (z *Zap) Logf(level Level, msg string, args ...interface{}) {
	if z.config.Debug {
		msg = normalizeMessage(msg)
	}
	switch level {
	case DebugLevel:
		z.sugar.Debugf(msg, args...)
	case InfoLevel:
		z.sugar.Infof(msg, args...)
	case WarnLevel:
		z.sugar.Warnf(msg, args...)
	case ErrorLevel:
		z.sugar.Errorf(msg, args...)
	case PanicLevel:
		z.sugar.Panicf(msg, args...)
	case FatalLevel:
		z.sugar.Fatalf(msg, args...)
	}
}

func (z *Zap) Logw(level Level, msg string, keysAndValues ...interface{}) {
	if z.config.Debug {
		msg = normalizeMessage(msg)
	}
	switch level {
	case DebugLevel:
		z.sugar.Debugw(msg, keysAndValues...)
	case InfoLevel:
		z.sugar.Infow(msg, keysAndValues...)
	case WarnLevel:
		z.sugar.Warnw(msg, keysAndValues...)
	case ErrorLevel:
		z.sugar.Errorw(msg, keysAndValues...)
	case PanicLevel:
		z.sugar.Panicw(msg, keysAndValues...)
	case FatalLevel:
		z.sugar.Fatalw(msg, keysAndValues...)
	}
}

func (z *Zap) With(fields ...LogField) Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	desugarLogger := z.desugar.With(zapFields...)
	return &Zap{
		desugar: desugarLogger,
		lv:      z.lv,
		sugar:   desugarLogger.Sugar(),
		config:  z.config,
	}
}

func (z *Zap) Sync() error {
	if err := z.desugar.Sync(); err != nil {
		return err
	}
	return z.sugar.Sync()
}

func (z *Zap) Close() {

}

func NewZapLogger(configs ...*ZapConfig) *Zap {
	config := DefaultZapConfig()
	if len(configs) > 0 {
		config = configs[0]
	}
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddStacktrace(zap.DPanicLevel))
	if config.AddCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(config.CallerSkip))
	}
	if len(config.Fields) > 0 {
		zapOptions = append(zapOptions, zap.Fields(config.Fields...))
	}

	var ws zapcore.WriteSyncer
	if config.Debug {
		ws = os.Stdout
	} else {
		ws = zapcore.AddSync(newRotate(config))
	}

	if config.Async {
		ws, _ = Buffer(ws, defaultBufferSize, defaultFlushInterval)
	}

	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := lv.UnmarshalText([]byte(config.Level)); err != nil {
		panic(err)
	}

	encoderConfig := *config.EncoderConfig
	core := config.Core
	if core == nil {
		core = zapcore.NewCore(
			func() zapcore.Encoder {
				if config.Debug {
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
		)
	}

	zapLogger := zap.New(
		core,
		zapOptions...,
	)
	return &Zap{
		desugar: zapLogger,
		lv:      &lv,
		config:  config,
		sugar:   zapLogger.Sugar(),
	}
}
