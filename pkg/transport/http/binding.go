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
	"github.com/go-ceres/ceres/internal/bytesconv"
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

func acquireBindRequest(req *Request) *bindRequest {
	request := bindRequestPool.Get().(*bindRequest)
	request.req = req
	return request
}

func releaseBindRequest(req *bindRequest) {
	req.req = nil
	bindRequestPool.Put(req)
}

func warpRequest(req *Request) *bindRequest {
	return acquireBindRequest(req)
}

type bindRequest struct {
	req *Request
}

func (br *bindRequest) GetParams() url.Values {
	return nil
}

func (br *bindRequest) GetMethod() string {
	return bytesconv.BytesToString(br.req.Header.Method())
}

func (br *bindRequest) GetQuery() url.Values {
	queryMap := make(url.Values)
	br.req.URI().QueryArgs().VisitAll(func(key, value []byte) {
		keyStr := bytesconv.BytesToString(key)
		values := queryMap[keyStr]
		values = append(values, bytesconv.BytesToString(value))
		queryMap[keyStr] = values
	})

	return queryMap
}

func (br *bindRequest) GetContentType() string {
	return bytesconv.BytesToString(br.req.Header.ContentType())
}

func (br *bindRequest) GetHeader() http.Header {
	header := make(http.Header)
	br.req.Header.VisitAll(func(key, value []byte) {
		keyStr := bytesconv.BytesToString(key)
		values := header[keyStr]
		values = append(values, bytesconv.BytesToString(value))
		header[keyStr] = values
	})
	return header
}

func (br *bindRequest) GetCookies() []*http.Cookie {
	var cookies []*http.Cookie
	br.req.Header.VisitAllCookie(func(key, value []byte) {
		cookies = append(cookies, &http.Cookie{
			Name:  string(key),
			Value: string(value),
		})
	})
	return cookies
}

func (br *bindRequest) GetBody() ([]byte, error) {
	return br.req.Body(), nil
}

func (br *bindRequest) GetPostForm() (url.Values, error) {
	postMap := make(url.Values)
	br.req.PostArgs().VisitAll(func(key, value []byte) {
		keyStr := bytesconv.BytesToString(key)
		values := postMap[keyStr]
		values = append(values, bytesconv.BytesToString(value))
		postMap[keyStr] = values
	})
	mf, err := br.req.MultipartForm()
	if err == nil {
		for k, v := range mf.Value {
			if len(v) > 0 {
				postMap[k] = v
			}
		}
	}

	return postMap, nil
}

func (br *bindRequest) GetForm() (url.Values, error) {
	formMap := make(url.Values)
	br.req.URI().QueryArgs().VisitAll(func(key, value []byte) {
		keyStr := bytesconv.BytesToString(key)
		values := formMap[keyStr]
		values = append(values, bytesconv.BytesToString(value))
		formMap[keyStr] = values
	})
	br.req.PostArgs().VisitAll(func(key, value []byte) {
		keyStr := bytesconv.BytesToString(key)
		values := formMap[keyStr]
		values = append(values, bytesconv.BytesToString(value))
		formMap[keyStr] = values
	})

	return formMap, nil
}
