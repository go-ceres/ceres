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
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/stringx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"strings"
)

//go:embed tpl/field.go.tpl
var fieldTemplate string

// genFields 生成所有字段代码
func (g *Generator) genFields(fields []*parser.Field, po bool) (string, error) {
	var fieldList []string
	for _, field := range fields {
		res, err := genField(field, po)
		if err != nil {
			return "", err
		}
		fieldList = append(fieldList, res)
	}
	return strings.Join(fieldList, "\n"), nil
}

// genQueryFields 生成查询字段
func (g *Generator) genQueryFields(fields []*parser.Field) (string, error) {
	var queryFieldList []string
	for _, field := range fields {
		this := new(parser.Field)
		this.OriginalName = field.OriginalName
		// 主键配置，用于多个主键查询
		if field.Primary {
			this.Name = stringx.NewString(field.Name.ToCamel() + "s")
			this.OriginalName = field.OriginalName + "s"
			this.Type = "[]" + field.Type
		} else if field.Unique || field.Fulltext {
			this.Name = field.Name
			switch field.Type {
			case "int64", "int32", "int", "uint64", "uint32", "uint", "float64", "float32":
				this.Type = "*" + field.Type
			default:
				this.Type = field.Type
			}
		} else {
			this.Name = field.Name
			this.Type = field.Type
		}
		res, err := genQueryField(this)
		if err != nil {
			return "", err
		}
		queryFieldList = append(queryFieldList, strings.ReplaceAll(res, "\n", ""))
	}
	return strings.Join(queryFieldList, "\n"), nil
}

// genQueryField 生成查询字段
func genQueryField(field *parser.Field) (string, error) {
	tag, err := genQueryTag(field)
	if err != nil {
		return "", err
	}
	tplText, err := pathx.LoadTpl(category, fieldTemplateFile, fieldTemplate)
	if err != nil {
		return "", err
	}
	text, err := templatex.With("field").Parse(tplText).Execute(map[string]interface{}{
		"name": field.Name.ToCamel(),
		"type": field.Type,
		"tag":  tag,
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}

// genField 生成带个字段代码
func genField(field *parser.Field, po bool) (string, error) {
	tplText, err := pathx.LoadTpl(category, fieldTemplateFile, fieldTemplate)
	if err != nil {
		return "", err
	}
	data := map[string]interface{}{
		"name": field.Name.ToCamel(),
		"type": field.Type,
	}
	if po {
		tag, err := genTag(field.Tag)
		if err != nil {
			return "", err
		}
		data["tag"] = tag
	}

	text, err := templatex.With("field").Parse(strings.ReplaceAll(tplText, "\n", "")).Execute(data)
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
