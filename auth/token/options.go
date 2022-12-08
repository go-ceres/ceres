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

type loginOptions struct {
	device  string // 登录的客户端设备标识
	timeout int64  // 指定当前登录token的有效时间（如果没有指定则使用全局配置）
}
type LoginOption func(o *loginOptions)

// DefaultOption 默认的登录额外参数
func defaultLoginOptions(c *Config) *loginOptions {
	opts := &loginOptions{
		timeout: c.Timeout,
		device:  defaultLoginDevice,
	}
	return opts
}

// LoginDevice 客户端
func LoginDevice(device string) LoginOption {
	return func(o *loginOptions) {
		o.device = device
	}
}

// LoginTimeout 当前此次登录的过期时间
func LoginTimeout(timeout int64) LoginOption {
	return func(o *loginOptions) {
		o.timeout = timeout
	}
}
