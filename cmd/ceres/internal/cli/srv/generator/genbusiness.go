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

//go:embed tpl/business.go.tpl
var businessTemplate string

// genRepository 生成领域层数据仓库接口
func (g *Generator) genBusiness(ctx DirContext, proto model.Proto) error {
	// 生成服务
	dir := ctx.GetBusiness()
	for _, service := range proto.Service {
		businessFilename, err := formatx.FileNamingFormat(g.style.Name, service.Name+"_business")
		if err != nil {
			return err
		}
		imports := make([]string, 0)
		irepositoryImport := fmt.Sprintf(`"%s"`, ctx.GetIRepository().Package)
		imports = append(imports, irepositoryImport)

		businessFile := filepath.Join(dir.Filename, businessFilename+".go")
		text, err := pathx.LoadTpl(category, businessTemplateFilename, businessTemplate)
		if err != nil {
			return err
		}

		err = templatex.With("business").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"imports":                strings.Join(imports, "\n"),
			"IRepositoryPackageName": ctx.GetIRepository().Base,
			"ServiceName":            stringx.NewString(service.Name).ToCamel(),
			"unTitleServiceName":     stringx.NewString(stringx.NewString(service.Name).ToCamel()).UnTitle(),
		}, businessFile, false)
		if err != nil {
			return err
		}
	}
	return g.genBusinessProvide(ctx, proto)
}

// genBusinessProvide 生成业务服务提供者
func (g *Generator) genBusinessProvide(ctx DirContext, proto model.Proto) error {
	dir := ctx.GetBusiness()
	provideFilename, err := formatx.FileNamingFormat(g.style.Name, "business")
	if err != nil {
		return err
	}
	provideFile := filepath.Join(dir.Filename, provideFilename+".go")
	provideSetList := make([]string, 0)
	for _, service := range proto.Service {
		provideSetList = append(provideSetList, fmt.Sprintf(`New%sBusiness`, stringx.NewString(service.Name).ToCamel()))
	}
	text, err := pathx.LoadTpl(category, provideFilename, provideTemplate)
	if err != nil {
		return err
	}
	return templatex.With("business-provide").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"PackageName":   dir.Base,
		"ProvideSetStr": strings.Join(provideSetList, ","),
	}, provideFile, true)
}
