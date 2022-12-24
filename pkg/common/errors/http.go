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

// Unauthorized 权限认证失败错误
func Unauthorized(reason, message string) *Error {
	return New(401, reason, message)
}

// IsUnauthorized 判断是否是权限错误
func IsUnauthorized(err error) bool {
	return Code(err) == 401
}

// Forbidden new Forbidden error that is mapped to a 403 response.
func Forbidden(reason, message string) *Error {
	return New(403, reason, message)
}

// IsForbidden determines if err is an error which indicates a Forbidden error.
// It supports wrapped errors.
func IsForbidden(err error) bool {
	return Code(err) == 403
}

// NotFound new NotFound error that is mapped to a 404 response.
func NotFound(reason, message string) *Error {
	return New(404, reason, message)
}

// IsNotFound determines if err is an error which indicates an NotFound error.
// It supports wrapped errors.
func IsNotFound(err error) bool {
	return Code(err) == 404
}

// MethodNotAllowed 请求超时
func MethodNotAllowed(reason, message string) *Error {
	return New(405, reason, message)
}

// IsMethodNotAllowed 是否是请求超时
func IsMethodNotAllowed(err error) bool {
	return Code(err) == 405
}

// RequestTimeout 请求超时
func RequestTimeout(reason, message string) *Error {
	return New(408, reason, message)
}

// IsRequestTimeout 是否是请求超时
func IsRequestTimeout(err error) bool {
	return Code(err) == 408
}

// Conflict 409 冲突错误
func Conflict(reason, message string) *Error {
	return New(409, reason, message)
}

// IsConflict 是否是409冲突错误
func IsConflict(err error) bool {
	return Code(err) == 409
}

// RequestHeaderFieldsTooLarge 请求头字段太大
func RequestHeaderFieldsTooLarge(reason, message string) *Error {
	return New(409, reason, message)
}

// IsRequestHeaderFieldsTooLarge 是否是请求头字段太大
func IsRequestHeaderFieldsTooLarge(err error) bool {
	return Code(err) == 409
}

// RequestEntityTooLarge 请求体太大
func RequestEntityTooLarge(reason, message string) *Error {
	return New(413, reason, message)
}

// IsRequestEntityTooLarge 是否是请求体太大
func IsRequestEntityTooLarge(err error) bool {
	return Code(err) == 413
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
