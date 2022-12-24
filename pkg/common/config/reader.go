//    Copyright 2022. ceres
//    Author https://github.com/go-ceres/ceres
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package config

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/go-ceres/ceres/pkg/common/codec"
	"github.com/imdario/mergo"
	"strings"
	"sync"
)

// 数据读取器
type reader struct {
	opts Options
	data map[string]interface{} // 数据容器
	lock sync.Mutex             // 数据锁
}

type Reader interface {
	Merge(...*DataSet) error
	Get(string) (Value, bool)
	Source() ([]byte, error)
	String() (string, error)
	Resolve() error
}

// newReader 创建数据读取器
func newReader(opts Options) Reader {
	return &reader{
		opts: opts,
		data: make(map[string]interface{}),
		lock: sync.Mutex{},
	}
}

// Merge 合并数据集
func (r *reader) Merge(dataSets ...*DataSet) error {
	merged, err := r.cloneMap()
	if err != nil {
		return err
	}
	for _, dataSet := range dataSets {
		next := make(map[string]interface{})
		if err := r.opts.decoder(dataSet, next); err != nil {
			return err
		}
		if err := mergo.Map(&merged, convertMap(next), mergo.WithOverride); err != nil {
			return err
		}
	}
	r.lock.Lock()
	r.data = merged
	r.lock.Unlock()
	return nil
}

// Get 读取值
func (r *reader) Get(path string) (Value, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return readValue(r.data, path)
}

// Source 获取byte数据
func (r *reader) Source() ([]byte, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return marshalJSON(convertMap(r.data))
}

// String 返回字符串
func (r *reader) String() (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	i, err := marshalJSON(convertMap(r.data))
	if err != nil {
		return "", err
	}
	return string(i), nil
}

// Resolve 重置数据
func (r *reader) Resolve() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.opts.resolver(r.data)
}

func (r *reader) cloneMap() (map[string]interface{}, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return cloneMap(r.data)
}

func cloneMap(src map[string]interface{}) (map[string]interface{}, error) {
	// https://gist.github.com/soroushjp/0ec92102641ddfc3ad5515ca76405f4d
	var buf bytes.Buffer
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(src)
	if err != nil {
		return nil, err
	}
	var clone map[string]interface{}
	err = dec.Decode(&clone)
	if err != nil {
		return nil, err
	}
	return clone, nil
}

func convertMap(src interface{}) interface{} {
	switch m := src.(type) {
	case map[string]interface{}:
		dst := make(map[string]interface{}, len(m))
		for k, v := range m {
			dst[k] = convertMap(v)
		}
		return dst
	case map[interface{}]interface{}:
		dst := make(map[string]interface{}, len(m))
		for k, v := range m {
			dst[fmt.Sprint(k)] = convertMap(v)
		}
		return dst
	case []interface{}:
		dst := make([]interface{}, len(m))
		for k, v := range m {
			dst[k] = convertMap(v)
		}
		return dst
	case []byte:
		// there will be no binary data in the config data
		return string(m)
	default:
		return src
	}
}

func readValue(values map[string]interface{}, path string) (Value, bool) {
	var (
		next = values
		keys = strings.Split(path, ".")
		last = len(keys) - 1
	)
	for idx, key := range keys {
		value, ok := next[key]
		if !ok {
			return nil, false
		}
		if idx == last {
			av := newAtomicValue(value)
			return av, true
		}
		switch vm := value.(type) {
		case map[string]interface{}:
			next = vm
		default:
			return nil, false
		}
	}
	return nil, false
}

func marshalJSON(v interface{}) ([]byte, error) {
	return codec.LoadCodec("json").Marshal(v)
}

func unmarshalJSON(data []byte, v interface{}) error {
	return codec.LoadCodec("json").Unmarshal(data, v)
}
