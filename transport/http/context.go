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
	"context"
	"github.com/go-ceres/ceres/middleware"
	"io"
)

// Context 定义上下文
type Context interface {
	context.Context
	Bind(data interface{}) error
	BindParams(out interface{}) error
	BindQuery(out interface{}) error
	BindBody(out interface{}) error
	Returns(code int, data interface{}) error
	Result(code int, data interface{}) error
	XML(code int, data interface{}) error
	JSON(code int, data interface{}) error
	Stream(code int, contextType string, rd io.Reader) error
	Middleware(h middleware.Handler) middleware.Handler
}
