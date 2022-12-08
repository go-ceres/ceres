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

type TokenInfo struct {
	TokenName            string `json:"tokenName"`            // token名称
	TokenValue           string `json:"tokenValue"`           // token值
	IsLogin              bool   `json:"isLogin"`              // 是否登录
	LogicType            string `json:"logicType"`            // 当前登录类型
	TokenTimeout         int64  `json:"tokenTimeout"`         // 当前token失效时间
	SessionTimeout       int64  `json:"sessionTimeout"`       // 共享session的有效时间
	TokenSessionTimeout  int64  `json:"tokenSessionTimeout"`  // 当前token的session有效期
	TokenActivityTimeout int64  `json:"tokenActivityTimeout"` // 当前token无操作状态下的有效时间
	LoginDevice          string `json:"loginDevice"`          // 当前登录设备
}
