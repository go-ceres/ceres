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

package transport

import (
	"context"
	_ "github.com/go-ceres/ceres/pkg/common/codec/json"
	_ "github.com/go-ceres/ceres/pkg/common/codec/proto"
	_ "github.com/go-ceres/ceres/pkg/common/codec/toml"
	_ "github.com/go-ceres/ceres/pkg/common/codec/xml"
	_ "github.com/go-ceres/ceres/pkg/common/codec/yaml"
	"net/url"
)

// Kind 服务类型
type Kind string

// String 服务类型转字符串
func (k Kind) String() string {
	return string(k)
}

// Transport 定义接口
type Transport interface {
	Kind() Kind
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Endpointer 服务地址接口
type Endpointer interface {
	Endpoint() (*url.URL, error)
}
