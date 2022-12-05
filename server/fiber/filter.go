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

import "github.com/gofiber/fiber/v2"

// FilterFunc 过滤器方法
type FilterFunc func(handler fiber.Handler) fiber.Handler

// FilterChain 执行通道
func FilterChain(filters ...interface{}) FilterFunc {
	return func(next fiber.Handler) fiber.Handler {
		for i := len(filters) - 1; i >= 0; i-- {
			switch filters[i].(type) {
			case FilterFunc:
				next = filters[i].(FilterFunc)(next)
			}
		}
		return next
	}
}
