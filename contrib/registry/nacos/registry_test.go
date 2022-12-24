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

package nacos

import (
	"context"
	"github.com/go-ceres/ceres/pkg/common/logger"
	"github.com/go-ceres/ceres/pkg/transport"
	"testing"
)

func TestTestRegistry(t *testing.T) {
	conf := DefaultOptions()
	conf.NamespaceId = "0799197e-ba01-4e9e-b6e5-815cbedef2d7"
	conf.Address = append(conf.Address, "http://123.57.16.239:8848/nacos")
	register := conf.Build()
	serverInfo := &transport.ServiceInfo{
		ID:        "123456789",
		Name:      "test",
		Version:   "v1.0.0",
		Endpoints: []string{"http://123.57.16.239:8848/nacos"},
		Metadata:  map[string]string{"ceshi": "1"},
	}
	err := register.Register(context.Background(), serverInfo)
	if err != nil {
		t.Error(err)
	}
	service, err := register.GetService(context.Background(), "test.http")
	if err != nil {
		t.Error(err)
	}
	err = register.Deregister(context.Background(), serverInfo)
	if err != nil {
		t.Error(err)
	}
	logger.Info("获取到了", logger.FieldAny("aaa", service))
}
