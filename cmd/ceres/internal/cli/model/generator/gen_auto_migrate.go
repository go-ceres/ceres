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
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"strings"
)

//go:embed tpl/automigrate.go.tpl
var autoMigrateTemplate string

// genAutoMigrate 生成迁移代码
func (g *Generator) genAutoMigrate(table parser.Table) (string, error) {
	tplText, err := pathx.LoadTpl(category, autoMigrateTemplateFile, autoMigrateTemplate)
	if err != nil {
		return "", err
	}
	text, err := templatex.With("option").Parse(tplText).Execute(map[string]interface{}{
		"camelName": table.Name.ToCamel(),
		"options":   strings.Join(table.Options, " "),
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
