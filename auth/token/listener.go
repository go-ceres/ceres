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
	"github.com/go-ceres/ceres/logger"
)

var _ Listener = (*defaultListener)(nil)

// Listener 监听器
type Listener interface {
	// DoLogin 每次登录触发
	// 形参:
	//	loginType – 账号类别
	//	loginId – 账号id
	//	loginOptions – 登录参数
	DoLogin(loginType string, loginId string, options loginOptions)
	// DoLogout 每次注销时触发
	//形参:
	//	loginType – 账号类别
	//	loginId – 账号id
	//	tokenValue – token值
	DoLogout(loginType string, loginId string, tokenValue string)
	// DoKickout 每次被踢下线时触发
	//形参:
	//	loginType – 账号类别
	//	loginId – 账号id
	//	tokenValue – token值
	DoKickout(loginType string, loginId string, tokenValue string)
	// DoReplaced 每次被顶下线时触发
	//形参:
	//	loginType – 账号类别
	//	loginId – 账号id
	//	tokenValue – token值
	DoReplaced(loginType string, loginId string, tokenValue string)
	// DoDisable 账号被封禁时触发
	//形参:
	//	loginType – 账号类别
	//	loginId – 账号id
	//	disableTime – 封禁的时长
	DoDisable(loginType string, loginId string, disableTime int64)
	// DoUntieDisable 每次账号被解封时触发
	//形参:
	//	loginType – 账号类别
	//	loginId – 账号id
	DoUntieDisable(loginType string, loginId string)
	// DoCreateSession 每次创建Session时触发
	//形参:
	//	id – SessionId
	DoCreateSession(id string)
	//DoLogoutSession 每次注销Session时触发
	//形参:
	//	id – SessionId
	DoLogoutSession(id string)
}

// defaultListener 默认的监听器
type defaultListener struct {
	logger *logger.Helper
}

func (d defaultListener) DoLogin(loginType string, loginId string, options loginOptions) {
	d.logger.Infof("登录逻辑：%s的用户：%s在设备：%s上登录成功", loginType, loginId, options.device)
}

func (d defaultListener) DoLogout(loginType string, loginId string, tokenValue string) {
	d.logger.Infof("登录逻辑：%s的用户：%s使用的token：%s退出成功", loginType, loginId, tokenValue)
}

func (d defaultListener) DoKickout(loginType string, loginId string, tokenValue string) {
	d.logger.Infof("登录逻辑：%s的用户：%s使用的token：%s被踢下线", loginType, loginId, tokenValue)
}

func (d defaultListener) DoReplaced(loginType string, loginId string, tokenValue string) {
	d.logger.Infof("登录逻辑：%s的用户：%s使用的token：%s被顶下线", loginType, loginId, tokenValue)
}

func (d defaultListener) DoDisable(loginType string, loginId string, disableTime int64) {
	d.logger.Infof("登录逻辑：%s的用户：%s被禁止登录，禁止时间：%d秒", loginType, loginId, disableTime)
}

func (d defaultListener) DoUntieDisable(loginType string, loginId string) {
	d.logger.Infof("登录逻辑：%s的用户：%s解除封禁", loginType, loginId)
}

func (d defaultListener) DoCreateSession(id string) {
	d.logger.Infof("用户：%s创建session成功", id)
}

func (d defaultListener) DoLogoutSession(id string) {
	d.logger.Infof("用户：%s退出成功", id)
}
