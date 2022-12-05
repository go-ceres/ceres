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

// LogField 日志字段
type LogField struct {
	Key   string
	Value interface{}
}

// FieldMod 模块
func FieldMod(mod string) LogField {
	return LogField{
		Key:   "mod",
		Value: mod,
	}
}

// FieldId 应用id
func FieldId(id string) LogField {
	return LogField{
		Key:   "appId",
		Value: id,
	}
}

// FieldName 应用name
func FieldName(name string) LogField {
	return LogField{
		Key:   "appName",
		Value: name,
	}
}

// FieldVersion 应用name
func FieldVersion(version string) LogField {
	return LogField{
		Key:   "appVersion",
		Value: version,
	}
}

// FieldTraceId trace id
func FieldTraceId(traceId string) LogField {
	return LogField{
		Key:   "traceId",
		Value: traceId,
	}
}

// FieldTraceSpan trace span
func FieldTraceSpan(traceSpan string) LogField {
	return LogField{
		Key:   "traceId",
		Value: traceSpan,
	}
}

// FieldError 错误字段
func FieldError(err error) LogField {
	return LogField{
		Key:   "err",
		Value: err,
	}
}

// FieldAny 任意类型
func FieldAny(key string, value interface{}) LogField {
	return LogField{
		Key:   key,
		Value: value,
	}
}
