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
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/model/parser"
	"github.com/go-ceres/ceres/cmd/ceres/internal/ctx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/stringx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
)

//go:embed tpl/import.go.tpl
var importTemplate string

// genImports 生成包导入
func (g *Generator) genImports(table parser.Table, cache, time bool, projectCtx *ctx.Project) (string, error) {
	tplText, err := pathx.LoadTpl(category, importTemplateFile, importTemplate)
	if err != nil {
		return "", err
	}

	text, err := templatex.With("import").Parse(tplText).Execute(map[string]interface{}{
		"time":        time,
		"cache":       cache,
		"unTitleName": stringx.NewString(table.Name.ToCamel()).UnTitle(),
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
