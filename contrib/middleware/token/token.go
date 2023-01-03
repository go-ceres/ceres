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

package token

import (
	"context"
	"github.com/go-ceres/ceres/pkg/common/auth/token"
	"github.com/go-ceres/ceres/pkg/common/errors"
	"github.com/go-ceres/ceres/pkg/transport"
	"strings"
)

type UserInfo struct {
	UserId string // 用户编号
}

// Server 创建服务器中间件
func Server(opts ...Option) transport.Middleware {
	o := DefaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	options := o.options
	for _, opt := range o.tokenOpts {
		opt(options)
	}
	logic := token.NewLogic(options)
	return func(handler transport.Handler) transport.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			md, ok := transport.MetadataFromServerContext(ctx)
			if !ok {
				return nil, errors.Unauthorized("NOT_HEADER", "no header")
			}
			headerToken := md.RequestHeader().Get(options.TokenName)
			if len(options.TokenPrefix) > 0 {
				headerToken = strings.TrimPrefix(headerToken, options.TokenPrefix)
			}
			// 如果没有id
			loginId, err := logic.GetLoginId(headerToken)
			if err != nil {
				return nil, err
			}
			info := &UserInfo{
				UserId: loginId,
			}
			// 检查权限
			if o.CheckPermission {
				permission := logic.HasPathPermission(loginId, md.Operation(), "*")
				if !permission {
					return nil, errors.Unauthorized("NO_PERMISSION", "no permission")
				}
			}
			// 设置header
			md.RequestHeader().Set("userId", loginId)
			// 设置上下文用户信息
			ctx = NewUserContext(ctx, info)
			return handler(ctx, req)
		}
	}
}

type userKey struct{}

// NewUserContext 设置用户id
func NewUserContext(ctx context.Context, info *UserInfo) context.Context {
	return context.WithValue(ctx, userKey{}, info)
}

// UserFromContext 从上下文中获取用户
func UserFromContext(ctx context.Context) (info *UserInfo, ok bool) {
	info, ok = ctx.Value(userKey{}).(*UserInfo)
	return
}
