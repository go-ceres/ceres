//    Copyright 2022. Go-Ceres
//    Author https://github.com/go-ceres/go-ceres
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

package wire

import (
	"github.com/go-ceres/ceres/cmd/ceres/internal/depend/golang"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/env"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/execx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/installer"
	"strings"
)

const (
	Name = "wire"
	url  = "github.com/google/wire/cmd/wire@latest"
)

// Install 安装
func Install(cacheDir string) (string, error) {
	return installer.Install(cacheDir, Name, func(dest string) (string, error) {
		err := golang.Install(url)
		return dest, err
	})
}

// Exists 检测是否存在
func Exists() bool {
	_, err := env.LookUpWire()
	return err == nil
}

// Version is used to get the version of the protoc-gen-go-grpc plugin.
func Version() (string, error) {
	path, err := env.LookUpWire()
	if err != nil {
		return "", err
	}
	version, err := execx.Command(path+" commands", "")
	if err != nil {
		return "", err
	}
	fields := strings.Fields(version)
	if len(fields) > 1 {
		return fields[1], nil
	}
	return "", nil
}
