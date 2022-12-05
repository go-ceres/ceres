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
	"fmt"
	"github.com/go-ceres/ceres/codec"
	"strings"
)

// Decoder 数据解码器
type Decoder func(*DataSet, map[string]interface{}) error

// defaultDecoder 默认的解码器
var defaultDecoder Decoder = func(set *DataSet, target map[string]interface{}) error {
	if set.Format == "" {
		keys := strings.Split(set.Key, ".")
		for i, k := range keys {
			if i == len(keys)-1 {
				target[k] = set.Data
			} else {
				sub := make(map[string]interface{})
				target[k] = sub
				target = sub
			}
		}
		return nil
	}
	if codec := codec.LoadCodec(set.Format); codec != nil {
		return codec.Unmarshal(set.Data, &target)
	}
	return fmt.Errorf("unsupported key: %s format: %s", set.Key, set.Format)
}
