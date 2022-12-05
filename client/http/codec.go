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
	"github.com/go-ceres/ceres/codec"
	"github.com/go-ceres/ceres/errors"
	"github.com/go-ceres/ceres/internal/httputil"
	"io"
	"net/http"
)

// DecodeErrorFunc 解码错误回调方法
type DecodeErrorFunc func(ctx context.Context, res *http.Response) error

// EncodeRequestFunc 请求编码器
type EncodeRequestFunc func(ctx context.Context, contentType string, in interface{}) (body []byte, err error)

// DecodeResponseFunc 响应解码器
type DecodeResponseFunc func(ctx context.Context, res *http.Response, out interface{}) error

// DefaultRequestEncoder is an HTTP request encoder.
func DefaultRequestEncoder(ctx context.Context, contentType string, in interface{}) ([]byte, error) {
	name := httputil.ContentSubtype(contentType)
	body, err := codec.LoadCodec(name).Marshal(in)
	if err != nil {
		return nil, err
	}
	return body, err
}

// DefaultResponseDecoder is an HTTP response decoder.
func DefaultResponseDecoder(ctx context.Context, res *http.Response, v interface{}) error {
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return CodecForResponse(res).Unmarshal(data, v)
}

// DefaultErrorDecoder 默认的错误解码器
func DefaultErrorDecoder(ctx context.Context, res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err == nil {
		e := new(errors.Error)
		if err = CodecForResponse(res).Unmarshal(data, e); err == nil {
			e.Code = int32(res.StatusCode)
			return e
		}
	}
	return errors.Newf(int32(res.StatusCode), errors.UnknownReason, "").WithCause(err)
}

// CodecForResponse 获取编解码器
func CodecForResponse(r *http.Response) codec.Codec {
	loadCodec := codec.LoadCodec(httputil.ContentSubtype(r.Header.Get("Content-Type")))
	if loadCodec != nil {
		return loadCodec
	}
	return codec.LoadCodec("json")
}
