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

package errors

// BadRequest 错误请求
func BadRequest(reason, message string) *Error {
	return New(400, reason, message)
}

// IsBadRequest 判断是否时请求错误
func IsBadRequest(err error) bool {
	return Code(err) == 400
}

// InternalServer 内部服务器错误
func InternalServer(reason, message string) *Error {
	return New(500, reason, message)
}

// IsInternalServer 是否是客户端管理
func IsInternalServer(err error) bool {
	return Code(err) == 500
}

// ServiceUnavailable 没有找到服务，映射到http的503
func ServiceUnavailable(reason, message string) *Error {
	return New(503, reason, message)
}

// IsServiceUnavailable 判断错误码是否是503
func IsServiceUnavailable(err error) bool {
	return Code(err) == 503
}

// GatewayTimeout 网关超时，对应到http的504
func GatewayTimeout(reason, message string) *Error {
	return New(504, reason, message)
}

// IsGatewayTimeout 判断是否是网关超时
func IsGatewayTimeout(err error) bool {
	return Code(err) == 504
}

// ClientClosed 客户端关闭对应http的499错误
func ClientClosed(reason, message string) *Error {
	return New(499, reason, message)
}

// IsClientClosed 判断是否是客户端关闭
func IsClientClosed(err error) bool {
	return Code(err) == 499
}
