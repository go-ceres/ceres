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

package token

import "github.com/go-ceres/ceres/pkg/common/auth/token"

// Option ...
type Option func(o *Options)

// Options token中间创建参数
type Options struct {
	options         *token.Options
	CheckPermission bool `json:"checkRole"`
	tokenOpts       []token.Option
}

func DefaultOptions() *Options {
	return &Options{
		CheckPermission: false,
		options:         token.DefaultOptions(),
		tokenOpts:       []token.Option{},
	}
}

// WithCheckPermission 检查角色权限
func WithCheckPermission(checkPermission bool) Option {
	return func(o *Options) {
		o.CheckPermission = checkPermission
	}
}

// WithTokenOptions token客户端配置
func WithTokenOptions(opts ...token.Option) Option {
	return func(o *Options) {
		o.tokenOpts = opts
	}
}
