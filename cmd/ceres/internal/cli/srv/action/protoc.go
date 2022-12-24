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

package action

import (
	"errors"
	"fmt"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/config"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/generator"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/stringx"
	"github.com/go-ceres/cli/v2"
	"github.com/gookit/gcli/v3/interact"
	"os"
	"path/filepath"
)

var ProtocFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "proto_path",
		Usage: "proto path",
	},
	&cli.StringFlag{
		Name:    "dist",
		Aliases: []string{"d"},
		Usage:   "proto filename",
	},
	//&cli.StringSliceFlag{
	//	Name:  "plugin_opt",
	//	Usage: `proto plugin option, Example: plugin is protoc-go-markdown, plugin_opt value "--plugin_opt=go-markdown_opt=s" `,
	//},
	&cli.StringFlag{
		Name:  "style",
		Value: "go_ceres",
		Usage: "filename format style",
	},
	&cli.StringSliceFlag{
		Name:  "go_opt",
		Value: nil,
		Usage: "protoc args for go_opt",
	},
	&cli.StringSliceFlag{
		Name:  "go-grpc_opt",
		Value: nil,
		Usage: "protoc args for go-grpc_opt",
	},
	&cli.StringSliceFlag{
		Name:  "plugin",
		Value: nil,
		Usage: "protoc args for plugin",
	},
}

var (
	errInvalidDistOutput = errors.New("ceres: missing --dist")
)

var ProtocAction cli.ActionFunc = func(ctx *cli.Context) error {
	// 创建默认配置
	conf := config.DefaultConfig()
	// 额外proto文件的路径
	protoPath := ctx.StringSlice("proto_path")
	// 如果有则添加到数组
	if len(protoPath) > 0 {
		conf.ProtoPath = append(conf.ProtoPath, protoPath...)
	}
	// go_opt go插件参数
	conf.GoOpt = ctx.StringSlice("go_opt")
	// go-grpc_opt grpc插件参数
	conf.GoGrpcOpt = ctx.StringSlice("go-grpc_opt")
	// proto_file proto文件路径
	conf.ProtoFile = ctx.Args().First()
	// 输出目录
	conf.Dist = ctx.String("dist")
	if len(conf.Dist) == 0 {
		return errInvalidDistOutput
	}
	// 检查goOut目录是否有效
	var err error
	// 检查goGrpcOut目录是否有效
	conf.ProtocOut, err = filepath.Abs(conf.ProtocOut)
	if err != nil {
		return err
	}
	// 创建protoc编译输出目录
	if err := pathx.MkdirIfNotFound(conf.ProtocOut); err != nil {
		return err
	}
	// 获取当前工作目录
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	// 模板相关
	home := ctx.String("home")
	remote := ctx.String("remote")
	branch := ctx.String("branch")
	if len(remote) > 0 {
		repo, _ := pathx.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		// 设置home目录
		pathx.RegisterCeresHome(home)
	}
	// 设置输出目录
	if !filepath.IsAbs(conf.Dist) {
		conf.Dist = filepath.Join(pwd, conf.Dist)
	}
	// 配置输出路径
	conf.Dist, err = filepath.Abs(conf.Dist)
	if err != nil {
		return err
	}
	// 是否打印日志
	verbose := ctx.Bool("verbose")
	// 注册中心
	registry := interact.SelectOne(
		"please select registration Center!",
		[]string{"none", "nacos"},
		"0",
	)
	if registry != "none" {
		registryComponent := &config.Component{
			Type: config.Registry,
			ExtraFunc: `func NewRegistry(registry *` + registry + `.Registry) transport.Registry {
    return registry
}

func NewDiscovery(registry *` + registry + `.Registry) transport.Discover  {
    return registry
}`,
			CamelName: stringx.NewString(registry).ToCamel(),
			Name:      stringx.NewString(registry),
			ImportPackage: []string{
				fmt.Sprintf(`"github.com/go-ceres/ceres/contrib/registry/%s"`, registry),
				`"github.com/go-ceres/ceres/pkg/transport"`,
			},
			InitStr: registry + `.ScanConfig().Build()`,
			ConfigStr: func() string {
				if registry == "nacos" {
					return `[application.transport.registry.` + registry + `]
    Address=["http://127.0.0.1:8488"]
`
				}
				return ""
			}(),
			TypeName: "*" + registry + ".Registry",
		}
		conf.Components = append(conf.Components, registryComponent)
		conf.Registry = true
	}
	// orm
	orm := interact.SelectOne(
		"please select orm!",
		[]string{"gorm"},
		"0",
	)
	if orm != "none" {
		conf.Components = append(conf.Components, &config.Component{
			Type:      config.Orm,
			CamelName: stringx.NewString(orm).ToCamel(),
			Name:      stringx.NewString(orm),
			ImportPackage: []string{
				fmt.Sprintf(`"github.com/go-ceres/ceres/pkg/common/store/%s"`, orm),
			},
			InitStr: orm + `.ScanConfig().Build()`,
			ConfigStr: `[application.store.` + orm + `]
    dns="user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
`,
			TypeName: "*" + orm + ".DB",
		})
	}

	// 是否增加http服务
	conf.HttpServer = interact.Confirm("add http server?", false)

	// 创建生成器
	g := generator.NewGenerator(ctx, "go_ceres", verbose)
	return g.Generate(conf)
}
