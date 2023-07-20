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
	"github.com/go-ceres/ceres/pkg/transport"
	"testing"
)

func Cors() transport.Middleware {
	return func(handler transport.Handler) transport.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			c, ok := ctx.(*Context)
			method := ""
			if ok {
				method = string(c.Request().Header.Method())
				c.Response().Header.Set("Access-Control-Allow-Origin", "*")
				c.Response().Header.Set("Access-Control-Allow-Headers", "*")
				c.Response().Header.Set("Access-Control-Allow-Methods", "*")
				c.Response().Header.Set("Access-Control-Expose-Headers", "*")
				c.Response().Header.Set("Access-Control-Allow-Credentials", "true")
				if method == "OPTIONS" {
					c.SetStatusCode(StatusOK)
					return nil, nil
				}
			}
			return handler(ctx, req)
		}
	}
}

func TestServer(t *testing.T) {
	srv := NewServer(
		WithServerMiddleware(
			Cors(),
		),
	)
	srv.Use(func(ctx *Context) error {
		method := string(ctx.Request().Header.Method())
		ctx.Response().Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response().Header.Set("Access-Control-Allow-Headers", "*")
		ctx.Response().Header.Set("Access-Control-Allow-Methods", "*")
		ctx.Response().Header.Set("Access-Control-Expose-Headers", "*")
		ctx.Response().Header.Set("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			ctx.SetStatusCode(StatusOK)
			return nil
		}
		return ctx.Next()
	})
	srv.Static("/files/uploads/*filepath", "./static/uploads")
	//srv.GET("/user/:id", func(ctx *Context) error {
	//	println("进来了")
	//	return errors.New(405, "错误", "错误信息")
	//})
	err := srv.Start(context.Background())
	if err != nil {
		t.Error(err)
	}
}
