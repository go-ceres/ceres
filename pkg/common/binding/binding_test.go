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

package binding

import (
	"testing"
)

type Values struct {
	Id      int64    `json:"id" path:"id"`
	Name    string   `json:"name" query:"name"`
	Cookie  string   `json:"cookie" cookie:"cookie"`
	Auth    string   `json:"auth" header:"auth"`
	Float   float64  `json:"float" header:"float"`
	ArrayT  []string `json:"arrayT" header:"arrayT"`
	Anytest struct {
		Ceshi string `json:"ceshi"`
	} `json:"anytest"`
}

func TestMarshal(t *testing.T) {
	bind := New()
	v := &Values{
		Id:     10,
		Name:   "liuqin",
		Auth:   "asdasdfasdasd",
		Cookie: "123456",
		Float:  3.1415926535,
		ArrayT: []string{"a", "b", "c"},
		Anytest: struct {
			Ceshi string `json:"ceshi"`
		}{Ceshi: "aaa"},
	}
	bind.marshal("POST", "/user/:id", v)
	//println(err)
}
