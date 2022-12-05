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

package model

import (
	"errors"
	"fmt"
	"github.com/emicklei/proto"
	"path/filepath"
	"strings"
)

type (
	// Services proto文件服务集合.
	Services []Service

	// Service proto文件单个服务的描述信息
	Service struct {
		*proto.Service
		RPC []*RPC
	}
)

func (s Services) Validate(filename string, multipleOpt ...bool) error {
	if len(s) == 0 {
		return errors.New("rpc service not found")
	}

	var multiple bool
	for _, c := range multipleOpt {
		multiple = c
	}

	if !multiple && len(s) > 1 {
		return errors.New("only one service expected")
	}

	name := filepath.Base(filename)
	for _, service := range s {
		for _, rpc := range service.RPC {
			if strings.Contains(rpc.RequestType, ".") {
				return fmt.Errorf("line %v:%v, request type must defined in %s",
					rpc.Position.Line,
					rpc.Position.Column, name)
			}
			if strings.Contains(rpc.ReturnsType, ".") {
				return fmt.Errorf("line %v:%v, returns type must defined in %s",
					rpc.Position.Line,
					rpc.Position.Column, name)
			}
		}
	}
	return nil
}
