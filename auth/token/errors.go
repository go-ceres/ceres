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
	"fmt"
	"github.com/go-ceres/ceres/errors"
)

// reason 列表
const (
	DisableLoginUser = "DISABLE_LOGIN_USER"
	NoToken          = "NO_TOKEN"
	NoUserId         = "NO_USER_ID"
	UserNotLogin     = "USER_NOT_LOGIN"
)

// key 过期常量
const (
	NeverExpire    int64 = -1 //常量，表示一个key永不过期 (在一个key被标注为永远不过期时返回此值)
	NotValueExpire int64 = -2 //常量，表示系统中不存在这个缓存 (在对不存在的key获取剩余存活时间时返回此值)
)

// AbnormalList 异常map
var AbnormalList = map[string]bool{
	NotToken:     true,
	InvalidToken: true,
	TokenTimeout: true,
	BeReplaced:   true,
	KickOut:      true,
}

// 异常与消息定义
const (
	NotToken            = "-1"
	NotTokenMessage     = "未提供Token"
	InvalidToken        = "-2"
	InvalidTokenMessage = "Token无效"
	TokenTimeout        = "-3"
	TokenTimeoutMessage = "Token已过期"
	BeReplaced          = "-4"
	BeReplacedMessage   = "Token已被顶下线"
	KickOut             = "-5"
	KickOutMessage      = "Token已被踢下线"
	DefaultMessage      = "当前会话未登录"
)

// ErrorDisableLogin 用户被封禁
func ErrorDisableLogin(format string, args ...interface{}) *errors.Error {
	return errors.Unauthorized(DisableLoginUser, fmt.Sprintf(format, args...))
}

func IsDisableLogin(err error) bool {
	return errors.Code(err) == 401 && errors.Reason(err) == DisableLoginUser
}

func ErrorNoToken(format string, args ...interface{}) *errors.Error {
	return errors.Unauthorized(NoToken, fmt.Sprintf(format, args...))
}

func IsNoToken(err error) bool {
	return errors.Code(err) == 401 && errors.Reason(err) == NoToken
}

func ErrorNoUserId(format string, args ...interface{}) *errors.Error {
	return errors.Unauthorized(NoUserId, fmt.Sprintf(format, args...))
}

func IsNoUserId(err error) bool {
	return errors.Code(err) == 401 && errors.Reason(err) == NoUserId
}

// ErrorNotLogin 用户未登录
func ErrorNotLogin(logicType string, code string, token string) *errors.Error {
	var msg = ""
	switch code {
	case NotToken:
		msg = NotTokenMessage
	case InvalidToken:
		msg = InvalidTokenMessage
	case TokenTimeout:
		msg = TokenTimeoutMessage
	case BeReplaced:
		msg = BeReplacedMessage
	case KickOut:
		msg = KickOutMessage
	default:
		msg = DefaultMessage
	}
	if len(token) == 0 {
		msg = msg + ":" + token
	}
	return errors.Unauthorized(UserNotLogin, msg).WithMetadata(map[string]string{
		"logicType": logicType,
		"code":      code,
		"message":   msg,
	})
}

func IsNotLogin(err error) bool {
	return errors.Code(err) == 401 && errors.Reason(err) == UserNotLogin
}
