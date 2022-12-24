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
	_ "embed"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/config"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/formatx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"path/filepath"
	"strings"
)

//go:embed tpl/component.go.tpl
var componentTemplate string

// genPkg 生成基础设施层的包依赖
func (g *Generator) genPkg(ctx DirContext, conf *config.Config) error {
	dir := ctx.GetPkg()
	err2 := pathx.MkdirIfNotExist(dir.Filename)
	if err2 != nil {
		return err2
	}
	packageFileName, err := formatx.FileNamingFormat(g.style.Name, "pkg")
	if err != nil {
		return err
	}
	newFuncNames := make([]string, 0)
	for _, component := range conf.Components {
		newFuncNames = append(newFuncNames, "New"+component.CamelName)
		if component.Type == config.Registry {
			newFuncNames = append(newFuncNames, "NewRegistry,NewDiscovery")
		}
	}
	fileName := filepath.Join(dir.Filename, packageFileName+".go")
	text, err := pathx.LoadTpl(category, provideTemplateFilename, provideTemplate)
	if err != nil {
		return err
	}
	// 生成子依赖
	for _, component := range conf.Components {
		err := g.genComponent(ctx, component)
		if err != nil {
			return err
		}
	}
	return templatex.With("pkg-provide").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"PackageName":   "pkg",
		"ProvideSetStr": strings.Join(newFuncNames, ","),
	}, fileName, false)
}

// genSubPkg 初始化组件
func (g *Generator) genComponent(ctx DirContext, component *config.Component) error {
	dir := ctx.GetPkg()
	subPackageFileName, err := formatx.FileNamingFormat(g.style.Name, component.Name.Source())
	if err != nil {
		return err
	}
	fileName := filepath.Join(dir.Filename, subPackageFileName+".go")
	text, err := pathx.LoadTpl(category, componentTemplateFilename, componentTemplate)
	if err != nil {
		return err
	}
	return templatex.With("component").GoFmt(true).Parse(text).SaveTo(component, fileName, false)
}
