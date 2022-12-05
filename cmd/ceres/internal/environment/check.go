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

package environment

import (
	"fmt"
	"github.com/go-ceres/ceres/cmd/ceres/internal/depend/protoc"
	"github.com/go-ceres/ceres/cmd/ceres/internal/depend/protocgenceres"
	"github.com/go-ceres/ceres/cmd/ceres/internal/depend/protocgengo"
	"github.com/go-ceres/ceres/cmd/ceres/internal/depend/protocgengogrpc"
	"github.com/go-ceres/ceres/cmd/ceres/internal/depend/wire"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/log"
	"strings"
	"time"
)

type bin struct {
	name   string
	exists bool
	get    func(cacheDir string) (string, error)
}

var bins = []bin{
	{
		name:   "protoc",
		exists: protoc.Exists(),
		get:    protoc.Install,
	},
	{
		name:   "wire",
		exists: wire.Exists(),
		get:    wire.Install,
	},
	{
		name:   "protoc-gen-go",
		exists: protocgengo.Exists(),
		get:    protocgengo.Install,
	},
	{
		name:   "protoc-gen-go-grpc",
		exists: protocgengogrpc.Exists(),
		get:    protocgengogrpc.Install,
	},
	{
		name:   "protoc-gen-ceres",
		exists: protocgenceres.Exists(),
		get:    protocgenceres.Install,
	},
}

func Prepare(install, force, verbose bool) error {
	log := log.NewLog(verbose)
	pending := true
	log.Info("[ceres-env]: preparing to check env")
	defer func() {
		if p := recover(); p != nil {
			log.Error("%+v", p)
			return
		}
		if pending {
			log.Info("\n[ceres-env]: congratulations! your ceres environment is ready!")
		} else {
			log.Error(`
[ceres-env]: check env finish, some dependencies is not found in PATH, you can execute
command 'ceres env check --install' to install it, for details, please execute command 
'ceres env check --help'`)
		}
	}()
	for _, e := range bins {
		time.Sleep(200 * time.Millisecond)
		log.Info("")
		log.Info("[ceres-env]: looking up %q", e.name)
		if e.exists {
			log.Info("[ceres-env]: %q is installed", e.name)
			continue
		}
		log.Warn("[ceres-env]: %q is not found in PATH", e.name)
		if install {
			install := func() {
				log.Info("[ceres-env]: preparing to install %q", e.name)
				path, err := e.get(Get(CeresCache))
				if err != nil {
					log.Error("[ceres-env]: an error interrupted the installation: %+v", err)
					pending = false
				} else {
					log.Info("[ceres-env]: %q is already installed in %q", e.name, path)
				}
			}
			if force {
				install()
				continue
			}
			log.Info("[ceres-env]: do you want to install %q [y: YES, n: No]", e.name)
			for {
				var in string
				_, _ = fmt.Scanln(&in)
				var brk bool
				switch {
				case strings.EqualFold(in, "y"):
					install()
					brk = true
				case strings.EqualFold(in, "n"):
					pending = false
					log.Info("[ceres-env]: %q installation is ignored", e.name)
					brk = true
				default:
					log.Error("[ceres-env]: invalid input, input 'y' for yes, 'n' for no")
				}
				if brk {
					break
				}
			}
		} else {
			pending = false
		}
	}
	return nil
}
