//    Copyright 2022. ceres
//    Author https://github.com/go-ceres/ceres
//
//    Licensed under the Apache License, CeresVersion 2.0 (the "License");
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

package environment

import (
	"fmt"
	"github.com/go-ceres/ceres/cmd/ceres/internal/depend/protoc"
	"github.com/go-ceres/ceres/cmd/ceres/internal/depend/protocgengo"
	"github.com/go-ceres/ceres/cmd/ceres/internal/depend/protocgengogrpc"
	sortedmap "github.com/go-ceres/ceres/cmd/ceres/internal/util/collection"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/version"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var ceresEnv *sortedmap.SortedMap

const (
	CeresOS                = "CERES_OS"
	CeresArch              = "CERES_ARCH"
	CeresHome              = "CERES_HOME"
	CeresDebug             = "CERES_DEBUG"
	CeresCache             = "CERES_CACHE"
	CeresVersion           = "CERES_VERSION"
	ProtocVersion          = "PROTOC_VERSION"
	ProtocGenGoVersion     = "PROTOC_GEN_GO_VERSION"
	ProtocGenGoGRPCVersion = "PROTO_GEN_GO_GRPC_VERSION"

	envFileDir = "env"
)

// init initializes the ceres environment variables, the environment variables of the function are set in order,
// please do not change the logic order of the code.
func init() {
	defaultCeresHome, err := pathx.GetDefaultCeresHome()
	if err != nil {
		log.Fatalln(err)
	}
	ceresEnv = sortedmap.New()
	ceresEnv.SetKV(CeresOS, runtime.GOOS)
	ceresEnv.SetKV(CeresArch, runtime.GOARCH)
	existsEnv := readEnv(defaultCeresHome)
	if existsEnv != nil {
		ceresHome, ok := existsEnv.GetString(CeresHome)
		if ok && len(ceresHome) > 0 {
			ceresEnv.SetKV(CeresHome, ceresHome)
		}
		if debug := existsEnv.GetOr(CeresDebug, "").(string); debug != "" {
			if strings.EqualFold(debug, "true") || strings.EqualFold(debug, "false") {
				ceresEnv.SetKV(CeresDebug, debug)
			}
		}
		if value := existsEnv.GetStringOr(CeresCache, ""); value != "" {
			ceresEnv.SetKV(CeresCache, value)
		}
	}
	if !ceresEnv.HasKey(CeresHome) {
		ceresEnv.SetKV(CeresHome, defaultCeresHome)
	}
	if !ceresEnv.HasKey(CeresDebug) {
		ceresEnv.SetKV(CeresDebug, "False")
	}

	if !ceresEnv.HasKey(CeresCache) {
		cacheDir, _ := pathx.GetCacheDir()
		ceresEnv.SetKV(CeresCache, cacheDir)
	}

	ceresEnv.SetKV(CeresVersion, version.Version)
	protocVer, _ := protoc.Version()
	ceresEnv.SetKV(ProtocVersion, protocVer)

	protocGenGoVer, _ := protocgengo.Version()
	ceresEnv.SetKV(ProtocGenGoVersion, protocGenGoVer)

	protocGenGoGrpcVer, _ := protocgengogrpc.Version()
	ceresEnv.SetKV(ProtocGenGoGRPCVersion, protocGenGoGrpcVer)
}

func Print() string {
	return strings.Join(ceresEnv.Format(), "\n")
}

func Get(key string) string {
	return GetOr(key, "")
}

func GetOr(key, def string) string {
	return ceresEnv.GetStringOr(key, def)
}

func readEnv(ceresHome string) *sortedmap.SortedMap {
	envFile := filepath.Join(ceresHome, envFileDir)
	data, err := os.ReadFile(envFile)
	if err != nil {
		return nil
	}
	dataStr := string(data)
	lines := strings.Split(dataStr, "\n")
	sm := sortedmap.New()
	for _, line := range lines {
		_, _, err = sm.SetExpression(line)
		if err != nil {
			continue
		}
	}
	return sm
}

func WriteEnv(kv []string) error {
	defaultCeresHome, err := pathx.GetDefaultCeresHome()
	if err != nil {
		log.Fatalln(err)
	}
	data := sortedmap.New()
	for _, e := range kv {
		_, _, err := data.SetExpression(e)
		if err != nil {
			return err
		}
	}
	data.RangeIf(func(key, value interface{}) bool {
		switch key.(string) {
		case CeresHome, CeresCache:
			path := value.(string)
			if !pathx.FileExists(path) {
				err = fmt.Errorf("[writeEnv]: path %q is not exists", path)
				return false
			}
		}
		if ceresEnv.HasKey(key) {
			ceresEnv.SetKV(key, value)
			return true
		} else {
			err = fmt.Errorf("[writeEnv]: invalid key: %v", key)
			return false
		}
	})
	if err != nil {
		return err
	}
	envFile := filepath.Join(defaultCeresHome, envFileDir)
	return os.WriteFile(envFile, []byte(strings.Join(ceresEnv.Format(), "\n")), 0o777)
}
