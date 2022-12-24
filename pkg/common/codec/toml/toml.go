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

package toml

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/go-ceres/ceres/pkg/common/codec"
)

const Name = "toml"

// tomlCodec protobuf编解码器
type tomlCodec struct{}

func init() {
	codec.RegisterCodec(tomlCodec{})
}

// Name 编码器名称
func (tomlCodec) Name() string {
	return Name
}

// Marshal 编码
func (tomlCodec) Marshal(v interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte(""))
	err := toml.NewEncoder(buf).Encode(v)
	return buf.Bytes(), err
}

// Unmarshal 解码
func (tomlCodec) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}
