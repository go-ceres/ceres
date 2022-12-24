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

var (
	defaultLogger = DefaultOptions().Build()
)

func SetLogger(log *Logger) {
	defaultLogger = log
}

func GetLogger() *Logger {
	return defaultLogger
}

func With(fields ...Field) *Logger {
	return defaultLogger.With(fields...)
}
func Debug(msg string, fields ...Field) {
	defaultLogger.Debug(msg, fields...)
}
func Info(msg string, fields ...Field) {
	defaultLogger.Info(msg, fields...)
}
func Warn(msg string, fields ...Field) {
	defaultLogger.Warn(msg, fields...)
}
func Error(msg string, fields ...Field) {
	defaultLogger.Error(msg, fields...)
}
func Panic(msg string, fields ...Field) {
	defaultLogger.Panic(msg, fields...)
}
func Fatal(msg string, fields ...Field) {
	defaultLogger.Fatal(msg, fields...)
}

func Debugf(msg string, args ...interface{}) {
	defaultLogger.Debugf(msg, args...)
}
func Infof(msg string, args ...interface{}) {
	defaultLogger.Infof(msg, args...)
}
func Warnf(msg string, args ...interface{}) {
	defaultLogger.Warnf(msg, args...)
}
func Errorf(msg string, args ...interface{}) {
	defaultLogger.Errorf(msg, args...)
}
func Panicf(msg string, args ...interface{}) {
	defaultLogger.Panicf(msg, args...)
}
func Fatalf(msg string, args ...interface{}) {
	defaultLogger.Fatalf(msg, args...)
}

func Debugw(msg string, keyValues ...interface{}) {
	defaultLogger.Debugw(msg, keyValues...)
}
func Infow(msg string, keyValues ...interface{}) {
	defaultLogger.Infow(msg, keyValues...)
}
func Warnw(msg string, keyValues ...interface{}) {
	defaultLogger.Warnw(msg, keyValues...)
}
func Errorw(msg string, keyValues ...interface{}) {
	defaultLogger.Errorw(msg, keyValues...)
}

func Panicw(msg string, keyValues ...interface{}) {
	defaultLogger.Panicw(msg, keyValues...)
}

func Fatalw(msg string, keyValues ...interface{}) {
	defaultLogger.Fatalw(msg, keyValues...)
}

func Sync() error {
	return defaultLogger.Sync()
}
