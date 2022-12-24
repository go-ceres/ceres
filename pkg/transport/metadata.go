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

import "context"

// Metadata 定义上下文接口
type Metadata interface {
	// Kind 元数据类型
	Kind() Kind
	// Endpoint 请求地址
	Endpoint() string
	// Operation 获取grpc请求全路径
	Operation() string
	// RequestHeader 请求头
	RequestHeader() Header
	// ReplyHeader 响应头
	ReplyHeader() Header
}

type (
	serverContextKey struct{}
	clientContextKey struct{}
)

type Header interface {
	Get(key string) string
	Keys() []string
	Set(key, value string)
}

// NewMetadataServerContext 创建一个上下文，存储应用上下
func NewMetadataServerContext(parent context.Context, md Metadata) context.Context {
	return context.WithValue(parent, serverContextKey{}, md)
}

// MetadataFromServerContext 从上下文中获取应用上下文
func MetadataFromServerContext(parent context.Context) (md Metadata, ok bool) {
	md, ok = parent.Value(serverContextKey{}).(Metadata)
	return
}

// NewMetadataClientContext 从客户上下文中获取元数据
func NewMetadataClientContext(parent context.Context, md Metadata) context.Context {
	return context.WithValue(parent, clientContextKey{}, md)
}

// MetadataFromClientContext 获取客户端上下文
func MetadataFromClientContext(parent context.Context) (md Metadata, ok bool) {
	md, ok = parent.Value(clientContextKey{}).(Metadata)
	return
}
