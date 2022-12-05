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

package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	//log := DefaultZapConfig().Build()
	//
	//log.sugar.Infow("", "msg", fmt.Sprint("aaa", "vvv", map[string]string{"a": "v"}))
	With(LogField{Key: "mod", Value: "测试"}).Debugf("aaaa %v", "aaa")
	With(LogField{Key: "mod", Value: "测试"}).Info("ceshi")
	Info("aaaa", FieldAny("any", map[string]string{"aaa": "vvv"}))
	Infow("测试", "mod", "abc", "ceshi", map[string]string{"aaa": "vvv"})
	Warn("[registry.nacos] nacos client init error：%v")
	Error("[registry.nacos] nacos client init error：%v")
}
