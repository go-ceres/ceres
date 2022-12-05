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

package json

import (
	"bytes"
	"encoding/json"
	"github.com/go-ceres/ceres/codec"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"reflect"
)

const Name = "json"

var (
	// ProtoMarshal 编码器
	ProtoMarshal = protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	// ProtoUnmarshal 解码器.
	ProtoUnmarshal = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

type jsonCodec struct{}

func init() {
	codec.RegisterCodec(jsonCodec{})
}

// Name 当前编码器名称
func (jsonCodec) Name() string {
	return Name
}

// Marshal 编码
func (jsonCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return ProtoMarshal.Marshal(m)
	default:
		return json.Marshal(v)
	}
}

// Unmarshal 解码
func (jsonCodec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return ProtoUnmarshal.Unmarshal(data, m)
	default:
		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return ProtoUnmarshal.Unmarshal(data, m)
		}
		d := json.NewDecoder(bytes.NewReader(data))
		d.UseNumber()
		return d.Decode(m)
	}
}
