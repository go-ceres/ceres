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
	"github.com/go-ceres/ceres/logger"
	"github.com/go-ceres/ceres/transport/http"
	"github.com/gofiber/fiber/v2"
	"testing"
)

func TestFiberServer(t *testing.T) {
	var filter FilterFunc = func(next fiber.Handler) fiber.Handler {
		return func(ctx *fiber.Ctx) error {
			logger.Infof("进来了")
			return next(ctx)
		}
	}
	srv := New()
	srv.Get("/member/user/:id", func(context http.Context) error {
		var in map[string]interface{}
		if err := context.BindQuery(&in); err != nil {
			return err
		}
		return nil
	}, filter)
	err := srv.Start(context.Background())
	if err != nil {
		logger.Errorw("启动失败", logger.FieldError(err))
	}

}
