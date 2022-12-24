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
	"go.uber.org/zap"
)

type Field = zap.Field

var (
	FieldInt8       = zap.Int8
	FieldInt16      = zap.Int16
	FieldInt32      = zap.Int32
	FieldInt64      = zap.Int64
	FieldDuration   = zap.Duration
	FieldError      = zap.Error
	FieldBool       = zap.Bool
	FieldUint       = zap.Uint
	FieldUint8      = zap.Uint8
	FieldUint16     = zap.Uint16
	FieldUint32     = zap.Uint32
	FieldUint64     = zap.Uint64
	FieldString     = zap.String
	FieldByteString = zap.ByteString
	FieldNamespace  = zap.Namespace
	FieldReflect    = zap.Reflect
	FieldSkip       = zap.Skip
	FieldAny        = zap.Any
	FieldObject     = zap.Object
)

// FieldAid 应用id
func FieldAid(aid string) Field {
	return FieldString("aid", aid)
}

// FieldName 应用名称
func FieldName(name string) Field {
	return FieldString("name", name)
}

// FieldMod 模块
func FieldMod(mod string) Field {
	return FieldString("mod", mod)
}

// FieldRegion 地域
func FieldRegion(region string) Field {
	return FieldString("region", region)
}

// FieldZone 分区
func FieldZone(zone string) Field {
	return FieldString("zone", zone)
}

// FieldVersion 版本
func FieldVersion(version string) Field {
	return FieldString("version", version)
}

// FieldHostName 主机名
func FieldHostName(hostName string) Field {
	return FieldString("hostName", hostName)
}

// FieldEndpoint 入口地址
func FieldEndpoint(endpoint string) Field {
	return FieldString("endpoint", endpoint)
}
