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
	"github.com/go-ceres/ceres/transport"
	"net/http"
)

var _ Transporter = (*Transport)(nil)

// Transporter http传输协议接口
type Transporter interface {
	transport.Transporter
	Request() *http.Request
	PathTemplate() string
}

// Transport http传输协议结构
type Transport struct {
	endpoint      string
	operation     string
	requestHeader headerInstance
	replyHeader   headerInstance
	request       *http.Request
	pathTemplate  string
}

func (tp *Transport) SetOperation(op string) {
	tp.operation = op
}

// Kind 返回当前传输协议类型
func (tp *Transport) Kind() transport.Kind {
	return "http"
}

// Operation 返回当前请求对应protoc 所生成的全链路地址
func (tp *Transport) Operation() string {
	return tp.operation
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

// Request 获取请求对象
func (tp *Transport) Request() *http.Request {
	return tp.request
}

// PathTemplate 获取路径模板
func (tp *Transport) PathTemplate() string {
	return tp.pathTemplate
}

// headerInstance 响应头实例
type headerInstance http.Header

// Get 获取头信息
func (hc headerInstance) Get(key string) string {
	return http.Header(hc).Get(key)
}

// Set 设置信息到头信息实例
func (hc headerInstance) Set(key string, value string) {
	http.Header(hc).Set(key, value)
}

// Keys 获取所有的头信息键
func (hc headerInstance) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}
