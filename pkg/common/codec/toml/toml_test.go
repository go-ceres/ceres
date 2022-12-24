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

import "testing"

type A struct {
	Name string
	B
}

type B struct {
	Name string
	C
}

type C struct {
	Name string
	D
}

type D struct {
	Name string
	F    string
}

func TestBuffWire(t *testing.T) {
	toml := tomlCodec{}
	data := A{
		Name: "aaa",
		B: B{
			Name: Name,
			C: C{
				Name: Name,
				D: D{
					F: Name,
				},
			},
		},
	}
	marshal, err := toml.Marshal(data)
	print(marshal, err)
}
