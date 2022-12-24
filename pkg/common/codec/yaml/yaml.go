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

package yaml

import (
	"github.com/go-ceres/ceres/pkg/common/codec"
	"gopkg.in/yaml.v3"
)

func init() {
	codec.RegisterCodec(yamlCodec{})
}

const Name = "yaml"

// Name 编解码器名称
func (yamlCodec) Name() string {
	return "yaml"
}

// yamlCodec yaml编解码器
type yamlCodec struct{}

// Marshal 编码
func (yamlCodec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// Unmarshal 解码
func (yamlCodec) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
