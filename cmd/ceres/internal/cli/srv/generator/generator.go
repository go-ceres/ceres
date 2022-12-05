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

package generator

import (
	"github.com/go-ceres/ceres/cmd/ceres/internal/environment"
	"github.com/go-ceres/ceres/cmd/ceres/internal/style"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/log"
	"github.com/go-ceres/ceres/logger"
	"github.com/go-ceres/cli/v2"
)

// Generator 代码生成器
type Generator struct {
	style   *style.Style
	log     *log.Log
	cliCtx  *cli.Context
	verbose bool
}

// NewGenerator 创建代码生成器
func NewGenerator(ctx *cli.Context, styleStr string, verbose bool) *Generator {
	newStyle, err := style.NewStyle(styleStr)
	if err != nil {
		logger.Panic("", logger.FieldError(err))
	}
	return &Generator{
		style:   newStyle,
		log:     log.NewLog(verbose),
		cliCtx:  ctx,
		verbose: verbose,
	}
}

// Prepare 前置步骤安装没有安装的依赖
func (g *Generator) Prepare() error {
	return environment.Prepare(true, true, g.verbose)
}
