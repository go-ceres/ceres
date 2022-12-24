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
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"
)

// Value 配置值接口
type Value interface {
	IsEmpty() bool
	Bool() (bool, error)
	Int() (int64, error)
	String() (string, error)
	Float() (float64, error)
	Duration() (time.Duration, error)
	Slice() ([]Value, error)
	Map() (map[string]Value, error)
	StringSlice() ([]string, error)
	StringMap() (map[string]string, error)
	Scan(val interface{}) error
	Bytes() ([]byte, error)
	Store(interface{})
	Load() interface{}
}

// newAtomicValue 创建原子值
func newAtomicValue(value interface{}) Value {
	v := &atomicValue{}
	v.Store(value)
	return v
}

type atomicValue struct {
	atomic.Value
}

// IsEmpty 判断是否为空
func (v *atomicValue) IsEmpty() bool {
	if v.Load() != nil {
		return false
	}
	return true
}

// Bool 获取bool值
func (v *atomicValue) Bool() (bool, error) {
	switch val := v.Load().(type) {
	case bool:
		return val, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
		return strconv.ParseBool(fmt.Sprint(val))
	}
	return false, v.typeAssertError()
}

// Int 获取int值
func (v *atomicValue) Int() (int64, error) {
	switch val := v.Load().(type) {
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
	case uint:
		return int64(val), nil
	case uint8:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint64:
		return int64(val), nil
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case json.Number:
		return val.Int64()
	case string:
		return strconv.ParseInt(val, 10, 64)
	}
	return 0, v.typeAssertError()
}

// String 获取string
func (v *atomicValue) String() (string, error) {
	switch val := v.Load().(type) {
	case string:
		return val, nil
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprint(val), nil
	case []byte:
		return string(val), nil
	case fmt.Stringer:
		return val.String(), nil
	}
	return "", v.typeAssertError()
}

// Float 获取float值
func (v *atomicValue) Float() (float64, error) {
	switch val := v.Load().(type) {
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case float32:
		return float64(val), nil
	case json.Number:
		return val.Float64()
	case float64:
		return val, nil
	case string:
		return strconv.ParseFloat(val, 64)
	}
	return 0.0, v.typeAssertError()
}

// Duration 获取成Duration
func (v *atomicValue) Duration() (time.Duration, error) {
	val, err := v.Int()
	if err != nil {
		return 0, err
	}
	return time.Duration(val), nil
}

// StringSlice 获取成字符串切片
func (v *atomicValue) StringSlice() ([]string, error) {
	vals, ok := v.Load().([]interface{})
	if !ok {
		return []string{}, v.typeAssertError()
	}
	slices := make([]string, 0, len(vals))
	for _, val := range vals {
		s, err := newAtomicValue(val).String()
		if err != nil {
			return nil, err
		}
		slices = append(slices, s)
	}
	return slices, nil
}

// StringMap 获取成stringMap
func (v *atomicValue) StringMap() (map[string]string, error) {
	m := make(map[string]string)
	vals, ok := v.Load().(map[string]interface{})
	if !ok {
		return m, v.typeAssertError()
	}
	for key, val := range vals {
		s, err := newAtomicValue(val).String()
		if err != nil {
			return m, err
		}
		m[key] = s
	}
	return m, nil
}

// Slice 获取Value切片
func (v *atomicValue) Slice() ([]Value, error) {
	vals, ok := v.Load().([]interface{})
	if !ok {
		return nil, v.typeAssertError()
	}
	slices := make([]Value, 0, len(vals))
	for _, val := range vals {
		slices = append(slices, newAtomicValue(val))
	}
	return slices, nil
}

// Map 获取Value Map
func (v *atomicValue) Map() (map[string]Value, error) {
	vals, ok := v.Load().(map[string]interface{})
	if !ok {
		return nil, v.typeAssertError()
	}
	m := make(map[string]Value, len(vals))
	for key, val := range vals {
		a := new(atomicValue)
		a.Store(val)
		m[key] = a
	}
	return m, nil
}

// Scan 扫描到指定对象
func (v *atomicValue) Scan(val interface{}) error {
	data, err := json.Marshal(v.Load())
	if err != nil {
		return err
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	return d.Decode(val)
}

// Bytes 获取字节码
func (v *atomicValue) Bytes() ([]byte, error) {
	switch val := v.Load().(type) {
	case []byte:
		return val, nil
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return []byte(fmt.Sprint(val)), nil
	case string:
		return []byte(val), nil
	case fmt.Stringer:
		return []byte(val.String()), nil
	}
	return []byte(""), v.typeAssertError()
}

func (v *atomicValue) typeAssertError() error {
	return fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.Load()))
}
