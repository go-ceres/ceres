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
	"github.com/go-ceres/ceres/pkg/common/errors"
	"testing"
)

func TestServer(t *testing.T) {
	srv := NewServer()
	srv.Use(func(ctx *Context) error {
		println("我也进来了")

		return ctx.Next()
	}).GET("/user/:id", func(ctx *Context) error {
		println("进来了")
		return errors.New(405, "错误", "错误信息")
	})
	err := srv.Start(context.Background())
	if err != nil {
		t.Error(err)
	}
}
