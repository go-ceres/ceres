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

package form

import (
	"github.com/go-ceres/ceres/codec"
	"github.com/go-playground/form/v4"
	"google.golang.org/protobuf/proto"
	"net/url"
	"reflect"
)

// nullStr 空字符串
const nullStr = "null"

var (
	encoder = form.NewEncoder()
	decoder = form.NewDecoder()
)

const Name = "x-www-form-urlencoded"

func init() {
	decoder.SetTagName("json")
	encoder.SetTagName("json")
	codec.RegisterCodec(formCodec{encoder: encoder, decoder: decoder})
}

// Name 编解码器名称
func (formCodec) Name() string {
	return Name
}

// formCodec yaml编解码器
type formCodec struct {
	encoder *form.Encoder
	decoder *form.Decoder
}

// Marshal 编码
func (f formCodec) Marshal(v interface{}) ([]byte, error) {
	var vs url.Values
	var err error
	if m, ok := v.(proto.Message); ok {
		vs, err = EncodeValues(m)
		if err != nil {
			return nil, err
		}
	} else {
		vs, err = f.encoder.Encode(v)
		if err != nil {
			return nil, err
		}
	}
	for k, v := range vs {
		if len(v) == 0 {
			delete(vs, k)
		}
	}
	return []byte(vs.Encode()), nil
}

// Unmarshal 解码
func (f formCodec) Unmarshal(data []byte, v interface{}) error {
	vs, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	if m, ok := v.(proto.Message); ok {
		return DecodeValues(m, vs)
	}
	if m, ok := rv.Interface().(proto.Message); ok {
		return DecodeValues(m, vs)
	}

	return f.decoder.Decode(v, vs)
}
