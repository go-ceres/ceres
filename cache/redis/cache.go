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

package redis

import (
	"encoding/json"
	"time"
)

type Cache struct {
	conf *Config
}

func NewCache(c ...*Config) *Cache {
	conf := DefaultConfig()
	if len(c) > 0 {
		conf = c[0]
	}
	if conf.client == nil {
		conf.logger.Fatalf("redis client not set")
	}
	return &Cache{
		conf: conf,
	}
}

// getSaveKey 获取实际存储key
func (c *Cache) getSaveKey(key string) string {
	if len(c.conf.Prefix) > 0 {
		return c.conf.Prefix + ":" + key
	}
	return key
}

// Has 查询是否包含缓存
func (c *Cache) Has(key string) bool {
	key = c.getSaveKey(key)
	for _, s := range c.conf.client.Keys(key) {
		if key == s {
			return true
		}
	}
	return false
}

// Get 获取缓存，可以带默认值
func (c *Cache) Get(key string, def ...string) string {
	key = c.getSaveKey(key)
	ret := ""
	if len(def) > 0 {
		ret = def[0]
	}
	if r := c.conf.client.Get(key); r != "" {
		ret = r
	}
	return ret
}

// Set 设置缓存
func (c *Cache) Set(key string, value string, timeout int64) bool {
	key = c.getSaveKey(key)
	return c.conf.client.Set(key, value, time.Second*time.Duration(timeout))
}

// SetObject 设置对象
func (c *Cache) SetObject(key string, value interface{}, timeout int64) bool {
	key = c.getSaveKey(key)
	marshal, err := json.Marshal(value)
	if err != nil {
		return false
	}
	return c.conf.client.Set(key, marshal, time.Second*time.Duration(timeout))
}

// GetObject 获取obj
func (c *Cache) GetObject(key string, obj interface{}) bool {
	key = c.getSaveKey(key)
	str := c.conf.client.Get(key)
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		return false
	}
	return true
}

// Update 修改数据,并且不修改过期时间
func (c *Cache) Update(key string, value string) {
	expire, err := c.conf.client.TTL(c.getSaveKey(key))
	if err != nil {
		return
	}
	c.Set(key, value, expire)
}

// UpdateObject 修改持久化数据
func (c *Cache) UpdateObject(key string, value interface{}) bool {
	expire, err := c.conf.client.TTL(c.getSaveKey(key))
	if err != nil {
		return false
	}
	return c.SetObject(key, value, expire)
}

// UpdateObjectTTl 修改持久化时间
func (c *Cache) UpdateObjectTTl(key string, timeout int64) {
	key = c.getSaveKey(key)
	_, _ = c.conf.client.Expire(key, time.Duration(timeout)*time.Second)
}

// Del 删除缓存
func (c *Cache) Del(key string) bool {
	key = c.getSaveKey(key)
	return c.conf.client.Del(key) == 0
}

// TTl 获取剩余过期时间
func (c *Cache) TTl(key string) int64 {
	key = c.getSaveKey(key)
	ttl, _ := c.conf.client.TTL(key)
	return ttl
}

// Clear 清除缓存
func (c *Cache) Clear() bool {
	keys := c.conf.client.Keys(c.conf.Prefix + "*")
	return c.conf.client.Del(keys...) == 0
}
