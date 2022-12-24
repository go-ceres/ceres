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
	"github.com/go-ceres/ceres/internal/bytesconv"
	"github.com/go-ceres/ceres/pkg/transport"
)

var _ transport.Metadata = (*Metadata)(nil)

// Metadata 元数据
type Metadata struct {
	endpoint     string    // 入口地址
	operation    string    // 操作地址
	pathTemplate string    // 地址模板
	request      *Request  // 请求信息
	response     *Response // 响应信息
}

func (md *Metadata) Kind() transport.Kind {
	return KindHttp
}

func (md *Metadata) Endpoint() string {
	return md.endpoint
}

func (md *Metadata) Operation() string {
	return md.operation
}

func (md *Metadata) Request() *Request {
	return md.request
}

func (md *Metadata) RequestHeader() transport.Header {
	return newHeaderCarrier(&md.request.Header)
}

func (md *Metadata) ReplyHeader() transport.Header {
	return newHeaderCarrier(&md.response.Header)
}

type header interface {
	Peek(key string) []byte
	PeekKeys() [][]byte
	Set(key, value string)
}

type headerCarrier struct {
	hd header
}

func newHeaderCarrier(hd header) *headerCarrier {
	return &headerCarrier{
		hd: hd,
	}
}

func (r *headerCarrier) Get(key string) string {
	return bytesconv.BytesToString(r.hd.Peek(key))
}

func (r *headerCarrier) Set(key, value string) {
	r.hd.Set(key, value)
}

func (r *headerCarrier) Keys() []string {
	keys := r.hd.PeekKeys()
	res := make([]string, 0, len(keys))
	for _, key := range keys {
		res = append(res, bytesconv.BytesToString(key))
	}
	return res
}

func SetOperation(ctx context.Context, op string) {
	md, ok := transport.MetadataFromServerContext(ctx)
	if ok {
		metadata, ok := md.(*Metadata)
		if ok {
			metadata.operation = op
		}
	}
}
