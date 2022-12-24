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

package http

import (
	"github.com/go-ceres/ceres/internal/httputil"
	"github.com/go-ceres/ceres/pkg/common/errors"
)

// HandlerFunc 方法定义
type HandlerFunc func(ctx *Context) error

// HandlersChain 方法组
type HandlersChain []HandlerFunc

// Last 获取最后一个方法
func (h HandlersChain) Last() HandlerFunc {
	if length := len(h); length > 0 {
		return h[length-1]
	}
	return nil
}

// DefaultErrorHandler 默认的错误处理方法
var DefaultErrorHandler = func(c *Context, err error) error {
	code := StatusInternalServerError
	var e *errors.Error
	body := err.Error()
	if errors.As(err, &e) {
		code = int(e.Code)
		subtype := httputil.ContentSubtype(c.GetRequestHeader(HeaderContentType))
		switch subtype {
		case "json":
			body = e.Json()
		default:
			body = e.Error()
		}
	}
	c.SetResponseHeader(HeaderContentType, MIMETextPlainCharsetUTF8)
	return c.SetStatusCode(code).SendString(body)
}
