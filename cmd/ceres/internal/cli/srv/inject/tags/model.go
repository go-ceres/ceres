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

// Field 字段
type Field struct {
	Name       string // 字段名称
	Start      int64  // tag开始
	End        int64  // tag结束
	CurrentTag string // 当前的标签
	InjectTag  string // 要注入的tag
}

// Message 消息信息
type Message struct {
	Name   string // 名称,转为全小写
	Fields map[string]*Field
}
