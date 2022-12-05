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

//go:embed tpl/entity.go.tpl
var entityTemplate string

// genEntity 生成领域层实体
func (g *Generator) genEntity(ctx DirContext, proto model.Proto) error {
	// 生成服务
	dir := ctx.GetEntity()
	for _, service := range proto.Service {
		entityFilename, err := formatx.FileNamingFormat(g.style.Name, service.Name+"_entity")
		if err != nil {
			return err
		}
		entityFile := filepath.Join(dir.Filename, entityFilename+".go")
		text, err := pathx.LoadTpl(category, entityTemplateFilename, entityTemplate)
		if err != nil {
			return err
		}

		err = templatex.With("entity").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"ServiceName": stringx.NewString(service.Name).ToCamel(),
		}, entityFile, false)
		if err != nil {
			return err
		}
	}
	return nil
}
