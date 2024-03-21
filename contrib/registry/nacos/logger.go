// Copyright 2023. ceres
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

package nacos

import (
	cLogger "github.com/go-ceres/ceres/pkg/common/logger"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
)

var _ logger.Logger = (*Logger)(nil)

type Logger struct {
	logger *cLogger.Logger
}

func NewLogger(log *cLogger.Logger) *Logger {
	return &Logger{
		logger: log,
	}
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Infow("", args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warnw("", args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Errorw("", args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debugw("", args...)
}

func (l *Logger) Infof(fmt string, args ...interface{}) {
	l.logger.Infof(fmt, args...)
}

func (l *Logger) Warnf(fmt string, args ...interface{}) {
	l.logger.Warnf(fmt, args...)
}

func (l *Logger) Errorf(fmt string, args ...interface{}) {
	l.logger.Errorf(fmt, args...)
}

func (l *Logger) Debugf(fmt string, args ...interface{}) {
	l.logger.Debugf(fmt, args...)
}
