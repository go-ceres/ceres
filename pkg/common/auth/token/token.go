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
	"github.com/go-ceres/ceres/internal/strings"
	"github.com/google/uuid"
	s "strings"
)

type Style string

const (
	StyleUuid       Style = "uuid"        // uuid样式
	StyleSimpleUuid Style = "simple-uuid" // uuid不带下划线
	StyleRandom32   Style = "random-32"   // 随机32位字符串
	StyleRandom64   Style = "random-64"   // 随机64位字符串
)

type Info struct {
	Name                 string `json:"tokenName"`            // token名称
	Value                string `json:"tokenValue"`           // token值
	IsLogin              bool   `json:"isLogin"`              // 是否登录
	LogicType            string `json:"logicType"`            // 当前登录类型
	TokenTimeout         int64  `json:"tokenTimeout"`         // 当前token失效时间
	SessionTimeout       int64  `json:"sessionTimeout"`       // 共享session的有效时间
	TokenSessionTimeout  int64  `json:"tokenSessionTimeout"`  // 当前token的session有效期
	TokenActivityTimeout int64  `json:"tokenActivityTimeout"` // 当前token无操作状态下的有效时间
	LoginDevice          string `json:"loginDevice"`          // 当前登录设备
}

// Sign 签名结构体
type Sign struct {
	Value  string `json:"value"`  // token值
	Device string `json:"device"` // 设备
}

// String 序列化成string
func (s *Sign) String() string {
	return "Sign [value=" + s.Value + ", device=" + s.Device + "]"
}

var _ Builder = (*defaultTokenBuilder)(nil)

// Builder token的生成接口
type Builder interface {
	Build(loginId string, logicType string, device string) string
}

type defaultTokenBuilder struct {
	style Style
}

func (d *defaultTokenBuilder) Build(loginId string, logicType string, device string) string {

	style := d.style
	switch style {
	case StyleUuid:
		return uuid.New().String()
	case StyleSimpleUuid:
		return s.ReplaceAll(uuid.New().String(), "_", "")
	case StyleRandom32:
		return strings.RandStr(32)
	case StyleRandom64:
		return strings.RandStr(64)
	default:
		return uuid.New().String()
	}
}
