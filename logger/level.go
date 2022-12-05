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

import "strings"

type Level uint8

const (
	LevelKey = "level"
	// DebugLevel 调试
	DebugLevel Level = iota
	// InfoLevel 正常
	InfoLevel
	// WarnLevel 警告
	WarnLevel
	// ErrorLevel 错误
	ErrorLevel
	// PanicLevel 致命，意外终止
	PanicLevel
	// FatalLevel 致命，但正常退出
	FatalLevel
)

// ParseLevel 解析等级
func ParseLevel(text string) Level {
	switch strings.ToUpper(text) {
	case "debug", "DEBUG":
		return DebugLevel
	case "info", "INFO", "": // make the zero value useful
		return InfoLevel
	case "warn", "WARN":
		return WarnLevel
	case "error", "ERROR":
		return ErrorLevel
	case "panic", "PANIC":
		return PanicLevel
	case "fatal", "FATAL":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// Key ...
func (l Level) Key() string {
	return LevelKey
}

// String 日志等级对应的字符串
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return ""
	}
}
