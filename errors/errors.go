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

import (
	"errors"
	"fmt"
	"github.com/go-ceres/ceres/transport"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

const (
	// SuccessCode 正确错误码
	SuccessCode = 200
	// UnknownCode 未知错误，服务器内部错误
	UnknownCode = 500
	// UnknownReason 未知的错误信息
	UnknownReason = ""
)

type Error struct {
	Status
	cause error
}

// New 新建错误
func New(code int32, reason string, message string) *Error {
	return &Error{
		Status: Status{
			Code:    code,
			Reason:  reason,
			Message: message,
		},
	}
}

// Newf 带格式化的新建
func Newf(code int32, reason, format string, a ...interface{}) *Error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

// Errorf 带格式化的错误，并返回错误接口类型
func Errorf(code int32, reason, format string, a ...interface{}) error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

// Error 实现错误接口
func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d mod = %s message = %s metadata = %v cause = %v", e.Code, e.Reason, e.Message, e.Metadata, e.cause)
}

// Raw 原始错误
func (e *Error) Raw() error {
	return e.cause
}

// Is 判断是否是指定错误
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Code == e.Code && se.Reason == e.Reason
	}
	return false
}

// WithMetadata 添加附加消息
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := Clone(e)
	err.Metadata = md
	return err
}

// AddMetadata 添加单个附加信息
func (e *Error) AddMetadata(key, value string) *Error {
	e.Metadata[key] = value
	return e
}

// WithCause 添加错误
func (e *Error) WithCause(cause error) *Error {
	err := Clone(e)
	err.cause = cause
	return err
}

// GRPCStatus 获取grpc状态
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(transport.ToGRPCCode(e.Code), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.Reason,
			Metadata: e.Metadata,
		})
	return s
}

// FromError 从原始错误中转换到框架内部错误
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if !ok {
		return New(UnknownCode, UnknownReason, err.Error())
	}
	ret := New(
		transport.FromGRPCCode(gs.Code()),
		UnknownReason,
		gs.Message(),
	)
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			ret.Reason = d.Reason
			return ret.WithMetadata(d.Metadata)
		}
	}
	return ret
}

// Code 获取错误码
func Code(err error) int32 {
	if err == nil {
		return SuccessCode
	}
	return FromError(err).Code
}

// Reason 获取错误原因
func Reason(err error) string {
	if err == nil {
		return UnknownReason
	}
	return FromError(err).Reason
}

// Clone 深度克隆
func Clone(err *Error) *Error {
	if err == nil {
		return nil
	}
	metadata := make(map[string]string, len(err.Metadata))
	for k, v := range err.Metadata {
		metadata[k] = v
	}
	return &Error{
		cause: err.cause,
		Status: Status{
			Code:     err.Code,
			Reason:   err.Reason,
			Message:  err.Message,
			Metadata: metadata,
		},
	}
}
