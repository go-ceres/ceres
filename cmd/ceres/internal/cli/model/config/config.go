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

package config

import "github.com/go-ceres/ceres/cmd/ceres/internal/style"

type ModelConfig struct {
	Src      string // 输入的文件
	Dist     string // 输出文件目录
	DataBase string // 数据库名
	Home     string // home文件夹
	Remote   string // 远程地址
	Branch   string // 远程分支
	Prefix   string // 前缀
	Strict   bool   // 严格模式
	Cache    bool   // 是否使用缓存
	Style    *style.Style
}

func DefaultPoConfig() *ModelConfig {
	return &ModelConfig{}
}
