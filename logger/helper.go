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

// Helper 帮助类
type Helper struct {
	log    Logger
	fields []LogField
}

// NewHelper 创建日志帮助
func NewHelper(logger Logger, fields ...LogField) *Helper {
	var log Logger
	if len(fields) > 0 {
		log = logger.With(fields...)
	} else {
		log = logger
	}
	return &Helper{
		log: log,
	}
}

func (h *Helper) GetLogger() Logger {
	return h.log
}

// With 设置全局字段
func (h *Helper) With(fields ...LogField) *Helper {
	log := h.log.With(fields...)
	return &Helper{
		log: log,
	}
}

func (h *Helper) Debug(msg string, fields ...LogField) {
	h.log.Log(DebugLevel, msg, fields...)
}
func (h *Helper) Info(msg string, fields ...LogField) {
	h.log.Log(InfoLevel, msg, fields...)
}
func (h *Helper) Warn(msg string, fields ...LogField) {
	h.log.Log(WarnLevel, msg, fields...)
}
func (h *Helper) Error(msg string, fields ...LogField) {
	h.log.Log(ErrorLevel, msg, fields...)
}
func (h *Helper) Panic(msg string, fields ...LogField) {
	h.log.Log(PanicLevel, msg, fields...)
}
func (h *Helper) Fatal(msg string, fields ...LogField) {
	h.log.Log(FatalLevel, msg, fields...)
}

func (h *Helper) Debugf(msg string, args ...interface{}) {
	h.log.Logf(DebugLevel, msg, args...)
}
func (h *Helper) Infof(msg string, args ...interface{}) {
	h.log.Logf(InfoLevel, msg, args...)
}
func (h *Helper) Warnf(msg string, args ...interface{}) {
	h.log.Logf(WarnLevel, msg, args...)
}
func (h *Helper) Errorf(msg string, args ...interface{}) {
	h.log.Logf(ErrorLevel, msg, args...)
}
func (h *Helper) Panicf(msg string, args ...interface{}) {
	h.log.Logf(PanicLevel, msg, args...)
}
func (h *Helper) Fatalf(msg string, args ...interface{}) {
	h.log.Logf(FatalLevel, msg, args...)
}

func (h *Helper) Debugw(msg string, keyValues ...interface{}) {
	h.log.Logw(DebugLevel, msg, keyValues...)
}
func (h *Helper) Infow(msg string, keyValues ...interface{}) {
	h.log.Logw(InfoLevel, msg, keyValues...)
}
func (h *Helper) Warnw(msg string, keyValues ...interface{}) {
	h.log.Logw(WarnLevel, msg, keyValues...)
}
func (h *Helper) Errorw(msg string, keyValues ...interface{}) {
	h.log.Logw(ErrorLevel, msg, keyValues...)
}
func (h *Helper) Panicw(msg string, keyValues ...interface{}) {
	h.log.Logw(PanicLevel, msg, keyValues...)
}
func (h *Helper) Fatalw(msg string, keyValues ...interface{}) {
	h.log.Logw(FatalLevel, msg, keyValues...)
}

func (h *Helper) Sync() error {
	return h.log.Sync()
}

func (h *Helper) Close() {
	h.log.Close()
}
