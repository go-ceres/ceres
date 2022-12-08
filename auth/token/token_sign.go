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

// TokenSign 签名结构体
type TokenSign struct {
	Value  string `json:"value"`  // token值
	Device string `json:"device"` // 设备
}

// String 序列化成string
func (s *TokenSign) String() string {
	return "TokenSign [value=" + s.Value + ", device=" + s.Device + "]"
}
