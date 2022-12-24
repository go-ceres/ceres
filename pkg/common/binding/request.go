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

package binding

import (
	"net/http"
	"net/url"
)

type Request interface {
	// GetMethod 获取请求方法
	GetMethod() string
	// GetParams 路径参数
	GetParams() url.Values
	// GetQuery 获取get请求参数
	GetQuery() url.Values
	// GetContentType 获取请求方式
	GetContentType() string
	// GetHeader 获取请求头
	GetHeader() http.Header
	// GetCookies 获取cookie
	GetCookies() []*http.Cookie
	// GetBody 获取请求body
	GetBody() ([]byte, error)
	// GetPostForm 获取post请求参数
	GetPostForm() (url.Values, error)
	// GetForm 获取get请求参数
	GetForm() (url.Values, error)
}

type MarshalRequest interface {
	Request
	// GetPath 获取完整路径
	GetPath() string
	// HasBody 是否有body
	HasBody() bool
	//Body 需要请求的body体
	Body() map[string]interface{}
}

type bindRequest struct {
	path    string
	hasBody bool
	params  url.Values
	query   url.Values
	header  http.Header
	cookie  []*http.Cookie
	body    map[string]interface{}
}

func (b bindRequest) GetPath() string {
	return b.path
}

func (b bindRequest) GetMethod() string {
	return ""
}

func (b bindRequest) GetParams() url.Values {
	return b.params
}

func (b bindRequest) GetQuery() url.Values {
	return b.query
}

func (b bindRequest) GetContentType() string {
	return ""
}

func (b bindRequest) GetHeader() http.Header {
	return b.header
}

func (b bindRequest) GetCookies() []*http.Cookie {
	return b.cookie
}

func (b *bindRequest) HasBody() bool {
	return b.hasBody
}

// Body 返回要编码的body
func (b *bindRequest) Body() map[string]interface{} {
	return b.body
}

func (b bindRequest) GetBody() ([]byte, error) {
	return nil, nil
}

func (b bindRequest) GetPostForm() (url.Values, error) {
	return nil, nil
}

func (b bindRequest) GetForm() (url.Values, error) {
	return nil, nil
}
