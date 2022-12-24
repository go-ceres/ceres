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

package tags

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	str := "ame     string `protobuf:\"bytes,2,opt,name=name,proto3\" json:\"name,omitempty\"`"
	findString := rInject.ReplaceAll([]byte(str), []byte(fmt.Sprintf("`%s`", "我进来了")))
	println(findString)
}
