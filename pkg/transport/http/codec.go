//// Copyright 2022. ceres
//// Author https://github.com/go-ceres/ceres
////
//// Licensed under the Apache License, Version 2.0 (the "License");
//// you may not use this file except in compliance with the License.
//// You may obtain a copy of the License at
////
//// http://www.apache.org/licenses/LICENSE-2.0
////
//// Unless required by applicable law or agreed to in writing, software
//// distributed under the License is distributed on an "AS IS" BASIS,
//// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//// See the License for the specific language governing permissions and
//// limitations under the License.
//

package http

import (
	"context"
	"github.com/go-ceres/ceres/internal/bytesconv"
	"github.com/go-ceres/ceres/internal/httputil"
	"github.com/go-ceres/ceres/pkg/common/codec"
	"github.com/go-ceres/ceres/pkg/common/errors"
	"strings"
)

type Redirector interface {
	Redirect() (string, int)
}

// EncodeRequestFunc 客户端请求体解码器定义
type EncodeRequestFunc = func(ctx context.Context, req *Request, v interface{}) ([]byte, error)

// DecodeResponseFunc 客户端响应体解码器定义
type DecodeResponseFunc = func(ctx context.Context, resp *Response, v interface{}) error

// DecodeErrorFunc 客户端请求错误回调方法定义
type DecodeErrorFunc func(ctx context.Context, res *Response) error

// EncodeResponseFunc 服务端响应体编码器定义
type EncodeResponseFunc func(ctx *Context, v interface{}) error

// CodecForResponse 根据响应信息获取解码器
func CodecForResponse(response *Response) codec.Codec {
	c := codec.LoadCodec(httputil.ContentSubtype(bytesconv.BytesToString(response.Header.Peek(HeaderContentType))))
	if c != nil {
		return c
	}
	return codec.LoadCodec("json")
}

// CodecForRequest 从上文中获取编码
func CodecForRequest(req *Request, name string) (codec.Codec, bool) {
	accepts := strings.Split(bytesconv.BytesToString(req.Header.Peek(name)), ",")
	for _, accept := range accepts {
		encoding := codec.LoadCodec(httputil.ContentSubtype(accept))
		if encoding != nil {
			return encoding, true
		}
	}
	return codec.LoadCodec("json"), false
}

// defaultRequestEncoder 默认的客户端请求体编码器
func defaultRequestEncoder(ctx context.Context, req *Request, in interface{}) ([]byte, error) {
	request, _ := CodecForRequest(req, HeaderContentType)
	body, err := request.Marshal(in)
	if err != nil {
		return nil, err
	}
	return body, err
}

// defaultResponseDecoder 默认的客户端响应体解码器
func defaultResponseDecoder(ctx context.Context, resp *Response, out interface{}) error {
	data := resp.Body()
	return CodecForResponse(resp).Unmarshal(data, out)
}

// defaultErrorDeCoder 默认的客户端请求错误回调方法
func defaultErrorDeCoder(ctx context.Context, resp *Response) error {
	if resp.StatusCode() >= 200 && resp.StatusCode() <= 299 {
		return nil
	}
	data := resp.Body()
	e := new(errors.Error)
	err := CodecForResponse(resp).Unmarshal(data, e)
	if err == nil {
		e.Code = int32(resp.StatusCode())
		return e
	}
	return errors.Newf(int32(resp.StatusCode()), errors.UnknownReason, "").WithCause(err)
}

// DefaultResponseEncode 默认的服务端响应体编码器
func DefaultResponseEncode(ctx *Context, v interface{}) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(Redirector); ok {
		url, code := rd.Redirect()
		return ctx.Redirect(code, bytesconv.StringToBytes(url))
	}
	c, _ := CodecForRequest(ctx.Request(), "Accept")
	data, err := c.Marshal(v)
	if err != nil {
		return err
	}
	ctx.Response().Header.Set("Content-Type", httputil.ContentType(c.Name()))
	_, err = ctx.Write(data)
	if err != nil {
		return err
	}
	return nil
}
