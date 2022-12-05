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

package fiber

import (
	"github.com/go-ceres/ceres/transport"
	"github.com/go-ceres/ceres/transport/http"
	"github.com/gofiber/fiber/v2"
)

var (
	KindFiber transport.Kind = "fiber"
	_         Transporter    = (*Transport)(nil)
)

type Transporter interface {
	http.Transporter
	Context() *fiber.Ctx
}

type Transport struct {
	endpoint      string
	operation     string
	requestHeader *RequestHeaderCarrier
	replyHeader   *ReplyHeaderCarrier
	context       *fiber.Ctx
	pathTemplate  string
}

func (t *Transport) SetOperation(op string) {
	t.operation = op
}

// Context 获取fiber的上下文
func (t *Transport) Context() *fiber.Ctx {
	return t.context
}

func (t *Transport) PathTemplate() string {
	return t.pathTemplate
}

func (t *Transport) Kind() transport.Kind {
	return KindFiber
}

func (t *Transport) Endpoint() string {
	return t.endpoint
}

func (t *Transport) Operation() string {
	return t.operation
}

func (t *Transport) RequestHeader() transport.Header {
	return t.requestHeader
}

func (t *Transport) ReplyHeader() transport.Header {
	return t.replyHeader
}

type RequestHeaderCarrier struct {
	ctx *fiber.Ctx
}

// Get 获取
func (r RequestHeaderCarrier) Get(key string) string {
	return string(r.ctx.Request().Header.Peek(key))
}

// Set 设置
func (r RequestHeaderCarrier) Set(key string, value string) {

	r.ctx.Request().Header.Set(key, value)
}

// Keys 获取全部header的键
func (r RequestHeaderCarrier) Keys() []string {
	headers := r.ctx.GetReqHeaders()
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	return keys
}

type ReplyHeaderCarrier struct {
	ctx *fiber.Ctx
}

// Get 获取header
func (r *ReplyHeaderCarrier) Get(key string) string {
	return r.ctx.GetRespHeader(key)
}

// Set 设置header
func (r *ReplyHeaderCarrier) Set(key string, value string) {
	r.ctx.Set(key, value)
}

func (r *ReplyHeaderCarrier) Keys() []string {
	headers := r.ctx.GetRespHeaders()
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	return keys
}
