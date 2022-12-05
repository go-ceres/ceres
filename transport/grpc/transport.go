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

package grpc

import (
	"github.com/go-ceres/ceres/selector"
	"github.com/go-ceres/ceres/transport"
	"google.golang.org/grpc/metadata"
)

var KindGrpc transport.Kind = "grpc"

var _ transport.Transporter = (*Transport)(nil)

type Transport struct {
	endpoint      string
	operation     string
	requestHeader HeaderInstance
	replyHeader   HeaderInstance
	nodeFilters   []selector.NodeFilter
}

// SetEndpoint 设置地址
func (tp *Transport) SetEndpoint(endpoint string) {
	tp.endpoint = endpoint
}

// SetRequestHeader 设置请求头
func (tp *Transport) SetRequestHeader(requestHeader HeaderInstance) {
	tp.requestHeader = requestHeader
}

// SetReplyHeader 设置响应头
func (tp *Transport) SetReplyHeader(replyHeader HeaderInstance) {
	tp.replyHeader = replyHeader
}

// SetNodeFilters 设置节点过滤器
func (tp *Transport) SetNodeFilters(nodeFilters []selector.NodeFilter) {
	tp.nodeFilters = nodeFilters
}

// SetOperation 设置当前操作
func (tp *Transport) SetOperation(op string) {
	tp.operation = op
}

// NodeFilters 获取node节点过滤器
func (tp *Transport) NodeFilters() []selector.NodeFilter {
	return tp.nodeFilters
}

// Operation 返回当前请求对应protoc 所生成的全链路地址
func (tp *Transport) Operation() string {
	return tp.operation
}

// Kind 返回当前传输协议类型
func (tp *Transport) Kind() transport.Kind {
	return KindGrpc
}

// Endpoint 返回服务地址
func (tp *Transport) Endpoint() string {
	return tp.endpoint
}

// RequestHeader 返回请求头实例
func (tp *Transport) RequestHeader() transport.Header {
	return tp.requestHeader
}

// ReplyHeader 返回响应头实例
func (tp *Transport) ReplyHeader() transport.Header {
	return tp.replyHeader
}

// HeaderInstance 响应头实例
type HeaderInstance metadata.MD

// Get 获取头信息
func (h HeaderInstance) Get(key string) string {
	values := metadata.MD(h).Get(key)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

// Set 设置信息到头信息实例
func (h HeaderInstance) Set(key string, value string) {
	metadata.MD(h).Set(key, value)
}

// Keys 获取所有的头信息键
func (h HeaderInstance) Keys() []string {
	keys := make([]string, 0, len(h))
	for key := range metadata.MD(h) {
		keys = append(keys, key)
	}
	return keys
}
