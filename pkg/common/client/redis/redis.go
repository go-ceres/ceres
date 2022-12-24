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
	"github.com/go-redis/redis"
	"time"
)

type Client struct {
	options *Options
	client  redis.UniversalClient
}

func New(opts ...Option) *Client {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return NewWithOptions(options)
}

func NewWithOptions(options *Options) *Client {
	cli := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:              options.Addrs,
		MaxRedirects:       options.MaxRetries,
		ReadOnly:           options.ReadOnly,
		Password:           options.Password,
		MaxRetries:         options.MaxRetries,
		DialTimeout:        options.DialTimeout,
		ReadTimeout:        options.ReadTimeout,
		WriteTimeout:       options.WriteTimeout,
		PoolSize:           options.PoolSize,
		MinIdleConns:       options.MinIdleConns,
		IdleTimeout:        options.IdleTimeout,
		OnConnect:          options.OnConnect,
		MinRetryBackoff:    options.MinRetryBackoff,
		MaxRetryBackoff:    options.MaxRetryBackoff,
		MaxConnAge:         options.MaxConnAge,
		PoolTimeout:        options.PoolTimeout,
		IdleCheckFrequency: options.IdleCheckFrequency,
		RouteByLatency:     options.RouteByLatency,
		RouteRandomly:      options.RouteRandomly,
		MasterName:         options.MasterName,
		TLSConfig:          options.TLSConfig,
	})
	if err := cli.Ping().Err(); err != nil {
		panic(err)
	}
	return &Client{
		client:  cli,
		options: options,
	}
}

// Keys 查询指定前缀的key
func (r *Client) Keys(pattern string) []string {
	sliceObj, err := r.client.Keys(pattern).Result()
	if err != nil {
		sliceObj = []string{}
	}
	return sliceObj
}

// Close 关闭redis连接
func (r *Client) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// Get 从redis获取string
func (r *Client) Get(key string) string {
	strCmd := r.client.Get(key)
	if err := strCmd.Err(); err != nil {
		return ""
	}
	return strCmd.Val()
}

// GetBytes 从redis获取[]byte
func (r *Client) GetBytes(key string) []byte {
	b, err := r.client.Get(key).Bytes()
	if err != nil && err != redis.Nil {
		return []byte{}
	}
	return b
}

func (r *Client) MGet(keys ...string) ([]string, error) {
	sliceObj := r.client.MGet(keys...)
	if err := sliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	tmp := sliceObj.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice, nil
}

func (r *Client) MGets(keys []string) ([]interface{}, error) {
	ret, err := r.client.MGet(keys...).Result()
	if err != nil && err != redis.Nil {
		return []interface{}{}, err
	}
	return ret, nil
}

// Set 设置redis的string
func (r *Client) Set(key string, value interface{}, expire time.Duration) bool {
	err := r.client.Set(key, value, expire).Err()
	return err == nil
}

// HGetAll 从redis获取hash的所有键值对
func (r *Client) HGetAll(key string) map[string]string {
	hashObj := r.client.HGetAll(key)
	hash := hashObj.Val()
	return hash
}

// HGet 从redis获取hash单个值
func (r *Client) HGet(key string, fields string) (string, error) {
	strObj := r.client.HGet(key, fields)
	err := strObj.Err()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if err == redis.Nil {
		return "", nil
	}
	return strObj.Val(), nil
}

// HMGetMap 批量获取hash值，返回map
func (r *Client) HMGetMap(key string, fields []string) map[string]string {
	if len(fields) == 0 {
		return make(map[string]string)
	}
	sliceObj := r.client.HMGet(key, fields...)
	if err := sliceObj.Err(); err != nil && err != redis.Nil {
		return make(map[string]string)
	}

	tmp := sliceObj.Val()
	hashRet := make(map[string]string, len(tmp))

	var tmpTagID string

	for k, v := range tmp {
		tmpTagID = fields[k]
		if v != nil {
			hashRet[tmpTagID] = v.(string)
		} else {
			hashRet[tmpTagID] = ""
		}
	}
	return hashRet
}

// HMSet 设置redis的hash
func (r *Client) HMSet(key string, hash map[string]interface{}, expire time.Duration) bool {
	if len(hash) > 0 {
		err := r.client.HMSet(key, hash).Err()
		if err != nil {
			return false
		}
		if expire > 0 {
			r.client.Expire(key, expire)
		}
		return true
	}
	return false
}

// HSet hset
func (r *Client) HSet(key string, field string, value interface{}) bool {
	err := r.client.HSet(key, field, value).Err()
	return err == nil
}

// HDel ...
func (r *Client) HDel(key string, field ...string) bool {
	IntObj := r.client.HDel(key, field...)
	err := IntObj.Err()
	return err == nil
}

// SetWithErr ...
func (r *Client) SetWithErr(key string, value interface{}, expire time.Duration) error {
	err := r.client.Set(key, value, expire).Err()
	return err
}

// SetNx 设置redis的string 如果键已存在
func (r *Client) SetNx(key string, value interface{}, expiration time.Duration) bool {

	result, err := r.client.SetNX(key, value, expiration).Result()

	if err != nil {
		return false
	}

	return result
}

// SetNxWithErr 设置redis的string 如果键已存在
func (r *Client) SetNxWithErr(key string, value interface{}, expiration time.Duration) (bool, error) {
	result, err := r.client.SetNX(key, value, expiration).Result()
	return result, err
}

// Incr redis自增
func (r *Client) Incr(key string) bool {
	err := r.client.Incr(key).Err()
	return err == nil
}

// IncrWithErr ...
func (r *Client) IncrWithErr(key string) (int64, error) {
	ret, err := r.client.Incr(key).Result()
	return ret, err
}

// IncrBy 将 key 所储存的值加上增量 increment 。
func (r *Client) IncrBy(key string, increment int64) (int64, error) {
	intObj := r.client.IncrBy(key, increment)
	if err := intObj.Err(); err != nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// Decr redis自减
func (r *Client) Decr(key string) bool {
	err := r.client.Decr(key).Err()
	return err == nil
}

// Scan ...
func (r *Client) Scan(cursor uint64, match string, count int64) ([]string, error) {
	result, _, err := r.client.Scan(cursor, match, count).Result()
	return result, err
}

// Type ...
func (r *Client) Type(key string) (string, error) {
	statusObj := r.client.Type(key)
	if err := statusObj.Err(); err != nil {
		return "", err
	}

	return statusObj.Val(), nil
}

// ZRevRange 倒序获取有序集合的部分数据
func (r *Client) ZRevRange(key string, start, stop int64) ([]string, error) {
	strSliceObj := r.client.ZRevRange(key, start, stop)
	if err := strSliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// ZRevRangeWithScores ...
func (r *Client) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	zSliceObj := r.client.ZRevRangeWithScores(key, start, stop)
	if err := zSliceObj.Err(); err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}
	return zSliceObj.Val(), nil
}

// ZRange ...
func (r *Client) ZRange(key string, start, stop int64) ([]string, error) {
	strSliceObj := r.client.ZRange(key, start, stop)
	if err := strSliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// ZRevRank ...
func (r *Client) ZRevRank(key string, member string) (int64, error) {
	intObj := r.client.ZRevRank(key, member)
	if err := intObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// ZRevRangeByScore ...
func (r *Client) ZRevRangeByScore(key string, opt redis.ZRangeBy) ([]string, error) {
	res, err := r.client.ZRevRangeByScore(key, opt).Result()
	if err != nil && err != redis.Nil {
		return []string{}, err
	}

	return res, nil
}

// ZRevRangeByScoreWithScores ...
func (r *Client) ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) ([]redis.Z, error) {
	res, err := r.client.ZRevRangeByScoreWithScores(key, opt).Result()
	if err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}

	return res, nil
}

// HMGet 批量获取hash值
func (r *Client) HMGet(key string, fileds []string) []string {
	sliceObj := r.client.HMGet(key, fileds...)
	if err := sliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	tmp := sliceObj.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice
}

// ZCard 获取有序集合的基数
func (r *Client) ZCard(key string) (int64, error) {
	IntObj := r.client.ZCard(key)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}
	return IntObj.Val(), nil
}

// ZScore 获取有序集合成员 member 的 score 值
func (r *Client) ZScore(key string, member string) (float64, error) {
	FloatObj := r.client.ZScore(key, member)
	err := FloatObj.Err()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	return FloatObj.Val(), err
}

// ZAdd 将一个或多个 member 元素及其 score 值加入到有序集 key 当中
func (r *Client) ZAdd(key string, members ...redis.Z) (int64, error) {
	IntObj := r.client.ZAdd(key, members...)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (r *Client) ZCount(key string, min, max string) (int64, error) {
	IntObj := r.client.ZCount(key, min, max)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// Del redis删除
func (r *Client) Del(key ...string) int64 {
	result, err := r.client.Del(key...).Result()
	if err != nil {
		return 0
	}
	return result
}

// DelWithErr ...
func (r *Client) DelWithErr(key string) (int64, error) {
	result, err := r.client.Del(key).Result()
	return result, err
}

// HIncrBy 哈希field自增
func (r *Client) HIncrBy(key string, field string, incr int) int64 {
	result, err := r.client.HIncrBy(key, field, int64(incr)).Result()
	if err != nil {
		return 0
	}
	return result
}

// HIncrByWithErr 哈希field自增并且返回错误
func (r *Client) HIncrByWithErr(key string, field string, incr int) (int64, error) {
	return r.client.HIncrBy(key, field, int64(incr)).Result()
}

// Exists 键是否存在
func (r *Client) Exists(key string) bool {
	result, err := r.client.Exists(key).Result()
	if err != nil {
		return false
	}
	return result == 1
}

// ExistsWithErr ...
func (r *Client) ExistsWithErr(key string) (bool, error) {
	result, err := r.client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// LPush 将一个或多个值 value 插入到列表 key 的表头
func (r *Client) LPush(key string, values ...interface{}) (int64, error) {
	IntObj := r.client.LPush(key, values...)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// RPush 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
func (r *Client) RPush(key string, values ...interface{}) (int64, error) {
	IntObj := r.client.RPush(key, values...)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// RPop 移除并返回列表 key 的尾元素。
func (r *Client) RPop(key string) (string, error) {
	strObj := r.client.RPop(key)
	if err := strObj.Err(); err != nil {
		return "", err
	}

	return strObj.Val(), nil
}

// LRange 获取列表指定范围内的元素
func (r *Client) LRange(key string, start, stop int64) ([]string, error) {
	result, err := r.client.LRange(key, start, stop).Result()
	if err != nil {
		return []string{}, err
	}

	return result, nil
}

// LLen ...
func (r *Client) LLen(key string) int64 {
	IntObj := r.client.LLen(key)
	if err := IntObj.Err(); err != nil {
		return 0
	}

	return IntObj.Val()
}

// LLenWithErr ...
func (r *Client) LLenWithErr(key string) (int64, error) {
	ret, err := r.client.LLen(key).Result()
	return ret, err
}

// LRem ...
func (r *Client) LRem(key string, count int64, value interface{}) int64 {
	IntObj := r.client.LRem(key, count, value)
	if err := IntObj.Err(); err != nil {
		return 0
	}

	return IntObj.Val()
}

// LIndex ...
func (r *Client) LIndex(key string, idx int64) (string, error) {
	ret, err := r.client.LIndex(key, idx).Result()
	return ret, err
}

// LTrim ...
func (r *Client) LTrim(key string, start, stop int64) (string, error) {
	ret, err := r.client.LTrim(key, start, stop).Result()
	return ret, err
}

// ZRemRangeByRank 移除有序集合中给定的排名区间的所有成员
func (r *Client) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	result, err := r.client.ZRemRangeByRank(key, start, stop).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

// Expire 设置过期时间
func (r *Client) Expire(key string, expiration time.Duration) (bool, error) {
	result, err := r.client.Expire(key, expiration).Result()
	if err != nil {
		return false, err
	}

	return result, err
}

// ZRem 从zset中移除变量
func (r *Client) ZRem(key string, members ...interface{}) (int64, error) {
	result, err := r.client.ZRem(key, members...).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// SAdd 向set中添加成员
func (r *Client) SAdd(key string, member ...interface{}) (int64, error) {
	intObj := r.client.SAdd(key, member...)
	if err := intObj.Err(); err != nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// SMembers 返回set的全部成员
func (r *Client) SMembers(key string) ([]string, error) {
	strSliceObj := r.client.SMembers(key)
	if err := strSliceObj.Err(); err != nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// SIsMember ...
func (r *Client) SIsMember(key string, member interface{}) (bool, error) {
	boolObj := r.client.SIsMember(key, member)
	if err := boolObj.Err(); err != nil {
		return false, err
	}
	return boolObj.Val(), nil
}

// HKeys 获取hash的所有域
func (r *Client) HKeys(key string) []string {
	strObj := r.client.HKeys(key)
	if err := strObj.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	return strObj.Val()
}

// HLen 获取hash的长度
func (r *Client) HLen(key string) int64 {
	intObj := r.client.HLen(key)
	if err := intObj.Err(); err != nil && err != redis.Nil {
		return 0
	}
	return intObj.Val()
}

// GeoAdd 写入地理位置
func (r *Client) GeoAdd(key string, location *redis.GeoLocation) (int64, error) {
	res, err := r.client.GeoAdd(key, location).Result()
	if err != nil {
		return 0, err
	}

	return res, nil
}

// GeoRadius 根据经纬度查询列表
func (r *Client) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	res, err := r.client.GeoRadius(key, longitude, latitude, query).Result()
	if err != nil {
		return []redis.GeoLocation{}, err
	}

	return res, nil
}

// TTL 查询过期时间
func (r *Client) TTL(key string) (int64, error) {
	if result, err := r.client.TTL(key).Result(); err != nil {
		return -2, err
	} else {
		return int64(result.Seconds()), nil
	}
}

// Subscribe 监听器
func (r *Client) Subscribe(keys ...string) *redis.PubSub {
	return r.client.Subscribe(keys...)
}
