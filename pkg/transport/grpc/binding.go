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
	"github.com/go-ceres/ceres/pkg/common/binding"
	"net/http"
	"net/url"
	"sync"
)

var (
	defaultBinding  = binding.New()
	bindRequestPool = sync.Pool{
		New: func() interface{} {
			return &bindRequest{}
		},
	}
)

func acquireBindRequest(metadata *Metadata) *bindRequest {
	request := bindRequestPool.Get().(*bindRequest)
	request.metadata = metadata
	return request
}

func releaseBindRequest(req *bindRequest) {
	req.metadata = nil
	bindRequestPool.Put(req)
}

func warpRequest(req *Metadata) *bindRequest {
	return acquireBindRequest(req)
}

type bindRequest struct {
	metadata *Metadata
}

func (br *bindRequest) GetParams() url.Values {
	return url.Values{}
}

func (br *bindRequest) GetMethod() string {
	return "POST"
}

func (br *bindRequest) GetQuery() url.Values {
	return url.Values{}
}

func (br *bindRequest) GetContentType() string {
	return "application/x-protobuf"
}

func (br *bindRequest) GetHeader() http.Header {
	header := make(http.Header)
	for _, key := range br.metadata.requestHeader.Keys() {
		header.Set(key, br.metadata.requestHeader.Get(key))
	}
	return header
}

func (br *bindRequest) GetCookies() []*http.Cookie {
	return []*http.Cookie{}
}

func (br *bindRequest) GetBody() ([]byte, error) {
	return []byte(""), nil
}

func (br *bindRequest) GetPostForm() (url.Values, error) {
	postMap := make(url.Values)
	return postMap, nil
}

func (br *bindRequest) GetForm() (url.Values, error) {
	formMap := make(url.Values)
	return formMap, nil
}
