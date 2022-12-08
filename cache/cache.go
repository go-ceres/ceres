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

package cache

// Cache 缓存接口
type Cache interface {
	// Has 查询是否包含缓存
	Has(key string) bool
	// Get 获取缓存
	Get(key string, def ...string) string
	// Set 写入缓存
	Set(key string, value string, timeout int64) bool
	// SetObject 设置对象
	SetObject(key string, value interface{}, timeout int64) bool
	// GetObject 获取obj
	GetObject(key string, obj interface{}) bool
	// Update 修改缓存，不修改时间
	Update(key string, value string)
	// UpdateObject 修改持久化数据
	UpdateObject(key string, value interface{}) bool
	// UpdateObjectTTl 修改持久化时间
	UpdateObjectTTl(key string, timeout int64)
	// Del 删除缓存
	Del(key string) bool
	// TTl 获取缓存剩余过期时间
	TTl(key string) int64
	// Clear 清空缓存
	Clear() bool
}
