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

package transport

import (
	"context"
	_ "github.com/go-ceres/ceres/codec/form"
	_ "github.com/go-ceres/ceres/codec/json"
	_ "github.com/go-ceres/ceres/codec/proto"
	_ "github.com/go-ceres/ceres/codec/toml"
	_ "github.com/go-ceres/ceres/codec/xml"
	_ "github.com/go-ceres/ceres/codec/yaml"
)

// Kind 服务类型
type Kind string

func (k Kind) String() string {
	return string(k)
}

// Transporter 协议
type Transporter interface {
	// Kind 服务传输类别
	Kind() Kind
	// Endpoint 服务地址
	Endpoint() string
	// Operation proto生成的操作方法
	Operation() string
	// SetOperation 设置生成的操作方法
	SetOperation(op string)
	// RequestHeader 请求头
	RequestHeader() Header
	// ReplyHeader 响应头
	ReplyHeader() Header
}

// Header 请求头信息
type Header interface {
	Get(key string) string
	Set(key string, value string)
	Keys() []string
}

type (
	serverTransportKey struct{}
	clientTransportKey struct{}
)

// NewServerContext returns a new Context that carries value.
func NewServerContext(ctx context.Context, tr Transporter) context.Context {
	return context.WithValue(ctx, serverTransportKey{}, tr)
}

// FromServerContext returns the Transport value stored in ctx, if any.
func FromServerContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(serverTransportKey{}).(Transporter)
	return
}

// NewClientContext returns a new Context that carries value.
func NewClientContext(ctx context.Context, tr Transporter) context.Context {
	return context.WithValue(ctx, clientTransportKey{}, tr)
}

// FromClientContext returns the Transport value stored in ctx, if any.
func FromClientContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(clientTransportKey{}).(Transporter)
	return
}
