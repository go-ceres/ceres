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

import "time"

// DataSet 数据源结构
type DataSet struct {
	Key       string    // 原始来源
	Data      []byte    // 原始数据
	Format    string    // 格式化
	Timestamp time.Time // 数据时间
}

// Watcher 数据监听器
type Watcher interface {
	Next() ([]*DataSet, error)
	Stop() error
}

// Source 数据源接口
type Source interface {
	Load() ([]*DataSet, error) // 加载数据
	Watch() (Watcher, error)   // 数据监听者
}
