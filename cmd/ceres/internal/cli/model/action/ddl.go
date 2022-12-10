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

package action

import (
	"errors"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/model/config"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/model/generator"
	"github.com/go-ceres/ceres/cmd/ceres/internal/style"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/cli/v2"
)

var DDlFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "src",
		Value:   "",
		Aliases: []string{"s"},
		Usage:   "The path or path globbing patterns of the ddl",
	},
	&cli.StringFlag{
		Name:    "dist",
		Value:   "",
		Aliases: []string{"d"},
		Usage:   "The target dir",
	},
	&cli.BoolFlag{
		Name:  "strict",
		Value: false,
		Usage: "The strict mode is enabled",
	},
	&cli.StringFlag{
		Name:  "database",
		Value: "",
		Usage: "database name",
	},
	&cli.StringFlag{
		Name:  "prefix",
		Value: "",
		Usage: "table name prefix",
	},
	&cli.StringFlag{
		Name:  "style",
		Value: "go_ceres",
		Usage: "The filename style",
	},
}
var DDlAction cli.ActionFunc = func(ctx *cli.Context) error {
	conf := config.DefaultPoConfig()
	conf.Src = ctx.String("src")
	conf.Dist = ctx.String("dist")
	newStyle, err := style.NewStyle(ctx.String("style"))
	if err != nil {
		return err
	}
	conf.Style = newStyle
	conf.DataBase = ctx.String("database")
	conf.Strict = ctx.Bool("strict")
	conf.Prefix = ctx.String("prefix")
	conf.Home = ctx.String("home")
	conf.Remote = ctx.String("remote")
	conf.Branch = ctx.String("branch")
	verbose := ctx.Bool("verbose")
	if len(conf.Src) == 0 {
		return errors.New("expected path or path globbing patterns, but nothing found")
	}
	// 模板相关
	if len(conf.Remote) > 0 {
		repo, _ := pathx.CloneIntoGitHome(conf.Remote, conf.Branch)
		if len(repo) > 0 {
			conf.Home = repo
		}
	}
	if len(conf.Home) > 0 {
		// 设置home目录
		pathx.RegisterCeresHome(conf.Home)
	}
	files, err := pathx.MatchFiles(conf.Src)
	if err != nil {
		return err
	}
	// 创建生成器
	g := generator.NewGenerator(ctx, conf.Style, verbose)
	for _, file := range files {
		if err := g.GeneratorFromDDl(file, conf); err != nil {
			return err
		}
	}
	return nil
}
