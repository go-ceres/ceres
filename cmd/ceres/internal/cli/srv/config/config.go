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

import "github.com/go-ceres/ceres/cmd/ceres/internal/util/stringx"

type ComponentType int8

const (
	Registry ComponentType = iota
	Orm
)

// Config 新建项目时的配置文件
type Config struct {
	Dist         string       // 项目输出路径，例如: .
	ProtocOut    string       // proto文件输出目录
	ConfigSource string       // 配置组件的类型，例如：file
	Registry     bool         // 注册中心，例如：etcd
	HttpServer   bool         // http服务实现
	ProtoPath    []string     // proto_path 参数
	GoOpt        []string     // protoc 的 opt参数
	GoGrpcOpt    []string     // protoc 的go-grpc_opt 参数
	Plugins      []string     // protoc的插件
	ProtoFile    string       // proto文件
	ProtocCmd    string       // protoc的命令
	Components   []*Component // 选择的组件
}

// Component 选择的额外组件
type Component struct {
	Type          ComponentType  // 组件类别
	ExtraFunc     string         // 额外的字符串
	CamelName     string         // 应用于 方法的名称
	Name          stringx.String // 组件名称
	ImportPackage []string       // 需要导入的包
	InitStr       string         // 初始化方法字符串
	ConfigStr     string         // 配置文件信息
	TypeName      string         // 类型
}

// DefaultConfig 默认配置信息
func DefaultConfig() *Config {
	return &Config{
		Components: make([]*Component, 0),
	}
}
