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

package i18n

import (
	"context"
	"github.com/go-ceres/ceres/middleware"
)

var global I18n

// Server 服务端中间件
func Server(opts ...Option) middleware.Middleware {
	global = newI18n(opts...)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 设置当前上下文
			global.setCurrentContext(ctx)
			return handler(ctx, req)
		}
	}
}

// GetMessage 获取消息
func GetMessage(param interface{}) (string, error) {
	return global.getMessage(param)
}

// GetMustMessage 忽略错误的数据
func GetMustMessage(param interface{}) string {
	return global.mustGetMessage(param)
}
