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

package fiber

import (
	"github.com/go-ceres/ceres/codec"
	"github.com/go-ceres/ceres/errors"
	"github.com/go-ceres/ceres/internal/httputil"
	"github.com/go-ceres/ceres/transport/http/binding"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"strings"
)

type Redirector interface {
	Redirect() (string, int)
}

type EncodeErrorFunc func(ctx *fiber.Ctx, err error)

// EncodeReplyFunc 定义响应体编码器方法
type EncodeReplyFunc func(ctx *fiber.Ctx, out interface{}) error

// DecodeRequestFunc 请求体解码器方法
type DecodeRequestFunc func(*fiber.Ctx, interface{}) error

// DefaultErrorFunc 默认错误回调
func DefaultErrorFunc(ctx *fiber.Ctx, err error) {
	se := errors.FromError(err)
	encoding, _ := CodecForCtx(ctx, "Accept")
	body, err := encoding.Marshal(se)
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.Set("Content-Type", httputil.ContentType(encoding.Name()))
	ctx.Status(int(se.Code))
	_ = ctx.Send(body)
}

// CodecForCtx 从上文中获取编码
func CodecForCtx(ctx *fiber.Ctx, name string) (codec.Codec, bool) {
	accepts := strings.Split(ctx.Get(name), ",")
	for _, accept := range accepts {
		encoding := codec.LoadCodec(httputil.ContentSubtype(accept))
		if encoding != nil {
			return encoding, true
		}
	}
	return codec.LoadCodec("json"), false
}

// DefaultRequestQuery 默认的query参数解码器
func DefaultRequestQuery(f *fiber.Ctx, v interface{}) error {
	b := f.Request().URI().QueryArgs().QueryString()
	return binding.BindQueryByte(b, v)
}

// DefaultRequestParams 默认请求变量解析
func DefaultRequestParams(f *fiber.Ctx, v interface{}) error {
	varRaws := f.AllParams()
	vars := make(url.Values, len(varRaws))
	for k, v := range varRaws {
		vars[k] = []string{v}
	}
	return binding.BindQuery(vars, v)
}

// DefaultRequestDecoder 默认响应体编码
func DefaultRequestDecoder(ctx *fiber.Ctx, v interface{}) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(Redirector); ok {
		url, code := rd.Redirect()
		return ctx.Redirect(url, code)
	}
	c, _ := CodecForCtx(ctx, "Accept")
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
