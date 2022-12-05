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

var defaultHelper = NewHelper(NewZapLogger())

func SetLogger(log Logger) {
	defaultHelper = NewHelper(log)
}

func GetLogger() Logger {
	return defaultHelper.GetLogger()
}

// With 设置全局字段
func With(fields ...LogField) *Helper {
	return defaultHelper.With(fields...)
}

func Debug(msg string, fields ...LogField) {
	defaultHelper.Debug(msg, fields...)
}
func Info(msg string, fields ...LogField) {
	defaultHelper.Info(msg, fields...)
}
func Warn(msg string, fields ...LogField) {
	defaultHelper.Warn(msg, fields...)
}
func Error(msg string, fields ...LogField) {
	defaultHelper.Error(msg, fields...)
}
func Panic(msg string, fields ...LogField) {
	defaultHelper.Panic(msg, fields...)
}
func Fatal(msg string, fields ...LogField) {
	defaultHelper.Fatal(msg, fields...)
}

func Debugf(msg string, args ...interface{}) {
	defaultHelper.Debugf(msg, args...)
}
func Infof(msg string, args ...interface{}) {
	defaultHelper.Infof(msg, args...)
}
func Warnf(msg string, args ...interface{}) {
	defaultHelper.Warnf(msg, args...)
}
func Errorf(msg string, args ...interface{}) {
	defaultHelper.Errorf(msg, args...)
}
func Panicf(msg string, args ...interface{}) {
	defaultHelper.Panicf(msg, args...)
}
func Fatalf(msg string, args ...interface{}) {
	defaultHelper.Fatalf(msg, args...)
}

func Debugw(msg string, keyValues ...interface{}) {
	defaultHelper.Debugw(msg, keyValues...)
}
func Infow(msg string, keyValues ...interface{}) {
	defaultHelper.Infow(msg, keyValues...)
}
func Warnw(msg string, keyValues ...interface{}) {
	defaultHelper.Warnw(msg, keyValues...)
}
func Errorw(msg string, keyValues ...interface{}) {
	defaultHelper.Errorw(msg, keyValues...)
}
func Panicw(msg string, keyValues ...interface{}) {
	defaultHelper.Panicw(msg, keyValues...)
}
func Fatalw(msg string, keyValues ...interface{}) {
	defaultHelper.Fatalw(msg, keyValues...)
}
