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
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser/model"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/formatx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/stringx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"path/filepath"
)

//go:embed tpl/irepository.go.tpl
var irepositoryTemplate string

// genRepository 生成领域层数据仓库接口
func (g *Generator) genIRepository(ctx DirContext, proto model.Proto) error {
	// 生成服务
	dir := ctx.GetIRepository()
	for _, service := range proto.Service {
		irepositoryFilename, err := formatx.FileNamingFormat(g.style.Name, "i"+service.Name+"_repository")
		if err != nil {
			return err
		}
		irepositoryFile := filepath.Join(dir.Filename, irepositoryFilename+".go")
		text, err := pathx.LoadTpl(category, irepositoryTemplateFilename, irepositoryTemplate)
		if err != nil {
			return err
		}

		err = templatex.With("irepository").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"ServiceName": stringx.NewString(service.Name).ToCamel(),
		}, irepositoryFile, false)
		if err != nil {
			return err
		}
	}
	return nil
}
