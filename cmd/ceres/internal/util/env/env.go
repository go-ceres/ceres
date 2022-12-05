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

package env

import (
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/vars"
	"os/exec"
	"runtime"
	"strings"
)

const (
	bin                = "bin"
	binGo              = "go"
	binWire            = "wire"
	binProtoc          = "protoc"
	binProtocGenGo     = "protoc-gen-go"
	binProtocGenGrpcGo = "protoc-gen-go-grpc"
	binProtocGenCeres  = "protoc-gen-ceres"
)

// LookPath 根据平台给所检测的软件加上后缀
func LookPath(binName string) (string, error) {
	suffix := getExeSuffix()
	if len(suffix) > 0 && !strings.HasSuffix(binName, suffix) {
		binName = binName + suffix
	}

	bin, err := exec.LookPath(binName)
	if err != nil {
		return "", err
	}
	return bin, nil
}

// CanExec 判断环境是否可运行
func CanExec() bool {
	switch runtime.GOOS {
	case vars.OsJs, vars.OsIOS:
		return false
	default:
		return true
	}
}

// LookUpProtoc 根据平台获取protoc的软件路径
func LookUpProtoc() (string, error) {
	suffix := getExeSuffix()
	xProtoc := binProtoc + suffix
	return LookPath(xProtoc)
}

// LookUpProtocGenGo 根据平台获取protoc-gen-go软件路径
func LookUpProtocGenGo() (string, error) {
	suffix := getExeSuffix()
	xProtocGenGo := binProtocGenGo + suffix
	return LookPath(xProtocGenGo)
}

// LookUpWire 根据平台获取wire的软件路径
func LookUpWire() (string, error) {
	suffix := getExeSuffix()
	xWire := binWire + suffix
	return LookPath(xWire)
}

// LookUpProtocGenGoGrpc 根据平台获取protoc-gen-go-grpc软件路径
func LookUpProtocGenGoGrpc() (string, error) {
	suffix := getExeSuffix()
	xProtocGenGoGrpc := binProtocGenGrpcGo + suffix
	return LookPath(xProtocGenGoGrpc)
}

// LookUpProtocGenCeres 根据平台获取protoc-gen-ceres软件路径
func LookUpProtocGenCeres() (string, error) {
	suffix := getExeSuffix()
	xProtocGenGoGrpc := binProtocGenCeres + suffix
	return LookPath(xProtocGenGoGrpc)
}

// getExeSuffix 获取平台软件的运行后缀
func getExeSuffix() string {
	if runtime.GOOS == vars.OsWindows {
		return ".exe"
	}
	return ""
}
