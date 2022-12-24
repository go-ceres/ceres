// Copyright 2022. ceres
// Author https://github.com/go-ceres/ceres
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	_ "embed"
	"fmt"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/config"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser/model"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/formatx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/stringx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"gorm.io/gorm/utils"
	"path/filepath"
	"strings"
)

const actionFunctionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.actionName}}) {{.method}} ({{if .notStream}}ctx context.Context,{{if .hasReq}} req {{.request}}{{end}}{{else}}{{if .hasReq}} req {{.request}},{{end}}stream {{.streamBody}}{{end}}) ({{if .notStream}}{{.response}},{{end}}error) {
	// todo: data conversion between dto and bo, and delete this line
	return {{if .notStream}}&{{.responseType}}{},{{end}}nil
}
`

//go:embed tpl/action.go.tpl
var actionTemplate string

// genAction 生成操作实例
func (g *Generator) genAction(ctx DirContext, proto model.Proto, conf *config.Config) error {
	dir := ctx.GetAction()
	for _, service := range proto.Service {
		serviceName := service.Name
		// 导入business包
		businessImport := fmt.Sprintf(`"%s"`, ctx.GetBusiness().Package)

		for _, rpc := range service.RPC {
			var (
				err            error
				filename       string
				actionName     string
				actionFilename string
				packageName    string
			)

			actionName = fmt.Sprintf("%sAction", stringx.NewString(rpc.Name).ToCamel())
			actionFilename, err = formatx.FileNamingFormat(g.style.Name, rpc.Name+"_action")
			if err != nil {
				return err
			}
			if len(proto.Service) > 1 {
				childPkg, err := dir.GetChildPackage(serviceName)
				if err != nil {
					return err
				}
				packageName = filepath.Base(childPkg)
				if err := pathx.MkdirIfNotExist(filepath.Join(dir.Filename, filepath.Base(childPkg))); err != nil {
					return err
				}
				filename = filepath.Join(dir.Filename, packageName, actionFilename+".go")
			} else {
				packageName = "action"
				filename = filepath.Join(dir.Filename, actionFilename+".go")
			}

			functions, err := g.genActionFunction(serviceName, proto.PbPackage, actionName, rpc)
			if err != nil {
				return err
			}
			pbImport := fmt.Sprintf(`"%v"`, ctx.GetProto().Package)
			imports := []string{pbImport, businessImport}
			text, err := pathx.LoadTpl(category, actionTemplateFilename, actionTemplate)
			if err != nil {
				return err
			}

			if err = templatex.With("action").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
				"actionName":   actionName,
				"funcs":        functions,
				"notStream":    !rpc.StreamsRequest,
				"businessType": ctx.GetBusiness().Base + "." + stringx.NewString(service.Name).ToCamel() + "Business",
				"package":      packageName,
				"imports":      strings.Join(imports, "\n"),
			}, filename, false); err != nil {
				return err
			}
		}
	}
	return g.genActionProvide(ctx, proto)
}

// genActionFunction 生成操作实例方法
func (g *Generator) genActionFunction(serviceName, goPackage, actionName string,
	rpc *model.RPC) (string,
	error) {
	functions := make([]string, 0)
	text, err := pathx.LoadTpl(category, actionFuncTemplateFilename, actionFunctionTemplate)
	if err != nil {
		return "", err
	}

	comment := parser.GetComment(rpc.Doc())
	streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(serviceName),
		parser.CamelCase(rpc.Name), "Server")
	buffer, err := templatex.With("fun").Parse(text).Execute(map[string]interface{}{
		"actionName":   actionName,
		"method":       parser.CamelCase(rpc.Name),
		"request":      fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.RequestType)),
		"response":     fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.ReturnsType)),
		"responseType": fmt.Sprintf("%s.%s", goPackage, parser.CamelCase(rpc.ReturnsType)),
		"hasComment":   len(comment) > 0,
		"comment":      comment,
		"hasReq":       !rpc.StreamsRequest,
		"stream":       rpc.StreamsRequest || rpc.StreamsReturns,
		"notStream":    !rpc.StreamsRequest && !rpc.StreamsReturns,
		"streamBody":   streamServer,
	})
	if err != nil {
		return "", err
	}

	functions = append(functions, buffer.String())
	return strings.Join(functions, "\n"), nil
}

// genActionProvide 生成操作实例服务提供者
func (g *Generator) genActionProvide(ctx DirContext, proto model.Proto) error {
	dir := ctx.GetAction()
	provideFilename, err := formatx.FileNamingFormat(g.style.Name, "action")
	if err != nil {
		return err
	}
	actionProvideFile := filepath.Join(dir.Filename, provideFilename+".go")
	provideSetList := make([]string, 0)
	multiple := len(proto.Service) > 1
	var imports []string
	for _, service := range proto.Service {
		var actionPackage string
		var err error
		if multiple {
			actionPackage, err = dir.GetChildPackage(service.Name)
			if err != nil {
				return err
			}
			actionImport := fmt.Sprintf(`"%s"`, actionPackage)
			if !utils.Contains(imports, actionImport) {
				imports = append(imports, actionImport)
			}
		}
		for _, rpc := range service.RPC {
			provideSet := ""
			if multiple {
				actionDir := filepath.Base(actionPackage)
				provideSet = fmt.Sprintf(`%s.New%sAction`, actionDir, stringx.NewString(rpc.Name).ToCamel())
			} else {
				provideSet = fmt.Sprintf(`New%sAction`, stringx.NewString(rpc.Name).ToCamel())
			}
			provideSetList = append(provideSetList, provideSet)
		}
	}
	text, err := pathx.LoadTpl(category, provideTemplateFilename, provideTemplate)
	if err != nil {
		return err
	}
	return templatex.With("action-provide").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"PackageName":   dir.Base,
		"ImportsStr":    strings.Join(imports, "\n"),
		"ProvideSetStr": strings.Join(provideSetList, ","),
	}, actionProvideFile, true)
}
