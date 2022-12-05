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
	"context"
	"github.com/go-ceres/ceres/middleware"
	"github.com/go-ceres/ceres/transport"
	"github.com/go-ceres/ceres/transport/http"
	"github.com/gofiber/fiber/v2"
	"io"
	"time"
)

var _ http.Context = (*wrapper)(nil)

type wrapper struct {
	srv *Server
	ctx *fiber.Ctx
}

func (c *wrapper) Bind(data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *wrapper) BindParams(out interface{}) error {
	return c.srv.config.decParams(c.ctx, out)
}

func (c *wrapper) BindQuery(out interface{}) error {
	err := c.srv.config.decQuery(c.ctx, out)
	if err != nil {
		return err
	}
	return nil
}

func (c *wrapper) BindBody(out interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *wrapper) Returns(code int, data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *wrapper) Result(code int, data interface{}) error {
	c.ctx.Status(code)
	return c.srv.config.enReply(c.ctx, data)
}

func (c *wrapper) XML(code int, data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *wrapper) JSON(code int, data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *wrapper) Stream(code int, contextType string, rd io.Reader) error {
	//TODO implement me
	panic("implement me")
}

// Middleware 执行中间件
func (c *wrapper) Middleware(h middleware.Handler) middleware.Handler {
	if tr, ok := transport.FromServerContext(c.ctx.UserContext()); ok {
		return middleware.Chain(c.srv.config.middleware.Match(tr.Operation())...)(h)
	}
	return middleware.Chain(c.srv.config.middleware.Match(c.ctx.Route().Path)...)(h)
}

// reset 重置数据
func (c *wrapper) reset(ctx *fiber.Ctx) {
	c.ctx = ctx
}
func (c *wrapper) Deadline() (time.Time, bool) {
	if c.ctx == nil {
		return time.Time{}, false
	}
	return c.ctx.UserContext().Deadline()
}

func (c *wrapper) Done() <-chan struct{} {
	if c.ctx == nil {
		return nil
	}
	return c.ctx.UserContext().Done()
}

func (c *wrapper) Err() error {
	if c.ctx == nil {
		return context.Canceled
	}
	return c.ctx.UserContext().Err()
}

func (c *wrapper) Value(key interface{}) interface{} {
	if c.ctx == nil {
		return nil
	}
	return c.ctx.UserContext().Value(key)
}
