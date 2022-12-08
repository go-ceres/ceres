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

import "time"

// Session session会话
type Session struct {
	Id         string                 `json:"id"`         // sessionId
	CreateTime int64                  `json:"createTime"` // 创建session的时间戳
	DataMap    map[string]interface{} `json:"dataMap"`    // 数据存储
	SignList   []*TokenSign           `json:"signList"`   // 签名列表
	storage    Storage                // 数据存储接口
}

// NewSession 新建一个session
func NewSession(id string, storage Storage) *Session {
	return &Session{
		Id:         id,
		CreateTime: time.Now().UnixMilli(),
		DataMap:    map[string]interface{}{},
		SignList:   make([]*TokenSign, 0),
		storage:    storage,
	}
}

// GetSign 根据token值获取签名
func (s *Session) GetSign(tokenValue string) *TokenSign {
	for i := 0; i < len(s.SignList); i++ {
		if s.SignList[i].Value == tokenValue {
			return s.SignList[i]
		}
	}
	return nil
}

// GetTimeout 获取此session的剩余存活时间
func (s *Session) GetTimeout() int64 {
	return s.storage.TTl(s.Id)
}

// AddTokenSign 在user-session上记录签名
func (s *Session) AddTokenSign(tokenValue string, device string) {
	// 如果已经存在于列表中，则无需再次添加
	for _, s2 := range s.SignList {
		if s2.Value == tokenValue {
			return
		}
	}
	// 添加并更新
	s.SignList = append(s.SignList, &TokenSign{
		Value:  tokenValue,
		Device: device,
	})
	s.Update()
}

// RemoveSign 移除指定token值得签名
func (s *Session) RemoveSign(tokenValue string) {
	for i := 0; i < len(s.SignList); i++ {
		if s.SignList[i].Value == tokenValue {
			s.SignList = append(s.SignList[:i], s.SignList[i+1:]...)
		}
	}
	s.Update()
}

// LogoutByTokenSignCountToZero 当session上面的token签名列表长度为0时，注销用户级session
func (s *Session) LogoutByTokenSignCountToZero() {
	if len(s.SignList) == 0 {
		s.LoginOut()
	}
}

// LoginOut 注销会话
func (s *Session) LoginOut() {
	// 删除storage
	s.storage.Del(s.Id)
}

// Update 更新持久库
func (s *Session) Update() {
	s.storage.UpdateObject(s.Id, s)
}

// UpdateMinTimeout 修改此Session的最小剩余存活时间 (只有在Session的过期时间低于指定的minTimeout时才会进行修改)
//形参:
//	minTimeout – 过期时间 (单位: 秒)
func (s *Session) UpdateMinTimeout(minTimeout int64) {
	if s.GetTimeout() < minTimeout {
		s.storage.UpdateObjectTTl(s.Id, minTimeout)
	}
}
