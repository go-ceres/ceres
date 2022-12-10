//    Copyright 2022. Go-Ceres
//    Author https://github.com/go-ceres/go-ceres
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
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/model/config"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/model/parser"
	"github.com/go-ceres/ceres/cmd/ceres/internal/ctx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
)

//go:embed tpl/update.go.tpl
var updateTemplate string

// GenUpdate 生成修改代码
func (g *Generator) GenUpdate(table parser.Table, projectCtx *ctx.Project, args *config.ModelConfig) (string, error) {
	tplText, err := pathx.LoadTpl(category, updateTemplateFile, updateTemplate)
	if err != nil {
		return "", err
	}
	text, err := templatex.With("update").Parse(tplText).Execute(map[string]interface{}{
		"camelName":    table.Name.ToCamel(),
		"originalName": table.Primary.OriginalName,
		"fieldType":    table.Primary.Type,
		"fieldName":    table.Primary.Name.UnTitle(),
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
