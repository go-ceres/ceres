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
	"fmt"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser/model"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/formatx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/stringx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"path/filepath"
	"strings"
)

//go:embed tpl/repository.go.tpl
var repositoryTemplate string

// genRepository 生成数据仓库
func (g *Generator) genRepository(ctx DirContext, proto model.Proto) error {
	// 生成服务
	dir := ctx.GetRepository()
	irepositoryDir := ctx.GetIRepository()
	irepositoryImport := fmt.Sprintf(`"%s"`, irepositoryDir.Package)
	imports := []string{irepositoryImport}
	for _, service := range proto.Service {
		irepositoryFilename, err := formatx.FileNamingFormat(g.style.Name, service.Name+"_repository")
		if err != nil {
			return err
		}
		repositoryFile := filepath.Join(dir.Filename, irepositoryFilename+".go")
		text, err := pathx.LoadTpl(category, repositoryTemplateFilename, repositoryTemplate)
		if err != nil {
			return err
		}

		err = templatex.With("repository").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"Imports":                strings.Join(imports, ","),
			"IRepositoryPackageName": irepositoryDir.Base,
			"PackageName":            dir.Base,
			"ServiceName":            stringx.NewString(service.Name).ToCamel(),
		}, repositoryFile, false)
		if err != nil {
			return err
		}
	}
	return g.genRepositoryProvide(ctx, proto)
}

func (g *Generator) genRepositoryProvide(ctx DirContext, proto model.Proto) error {
	dir := ctx.GetRepository()
	provideFilename, err := formatx.FileNamingFormat(g.style.Name, "repository")
	if err != nil {
		return err
	}
	repositoryProvideFile := filepath.Join(dir.Filename, provideFilename+".go")
	provideSetList := make([]string, 0)
	for _, service := range proto.Service {
		provideSetList = append(provideSetList, fmt.Sprintf(`New%sRepository`, stringx.NewString(service.Name).ToCamel()))
	}
	text, err := pathx.LoadTpl(category, provideTemplateFilename, provideTemplate)
	if err != nil {
		return err
	}
	return templatex.With("repository-provide").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"PackageName":   dir.Base,
		"ProvideSetStr": strings.Join(provideSetList, ","),
	}, repositoryProvideFile, true)
}
