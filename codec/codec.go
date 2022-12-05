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

package codec

import (
	"strings"
	"sync"
)

// Codec 编解码器接口
type Codec interface {
	// Name 编解码器名称
	Name() string
	// Marshal 编码
	Marshal(interface{}) ([]byte, error)
	// Unmarshal 解码
	Unmarshal([]byte, interface{}) error
}

var codecs = &sync.Map{}

// RegisterCodec 注册编解码器
func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register nil codec")
	}
	if codec.Name() == "" {
		panic("cannot register codec with codec name is empty")
	}
	codecs.Store(strings.ToLower(codec.Name()), codec)
}

// LoadCodec 加载编解码器
func LoadCodec(name string) Codec {
	if v, ok := codecs.Load(name); ok {
		return v.(Codec)
	}
	return nil
}
