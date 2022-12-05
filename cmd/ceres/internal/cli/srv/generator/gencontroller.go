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
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/config"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser/model"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/formatx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/stringx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"path/filepath"
	"strings"
)

const functionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.service}}Controller) {{.method}} ({{if .notStream}}ctx context.Context,{{if .hasReq}} req {{.request}}{{end}}{{else}}{{if .hasReq}} req {{.request}},{{end}}stream {{.streamBody}}{{end}}) ({{if .notStream}}{{.response}},{{end}}error) {
	return {{if .notStream}}&{{.responseType}}{},{{end}}nil
}
`

//go:embed tpl/controller.go.tpl
var controllerTemplate string

type LogicDesc struct {
	LogicName    string
	UnTitleName  string
	LogicPackage string
}

// GenController 生成grpc实现类,业务入口
func (g *Generator) GenController(ctx DirContext, proto model.Proto, conf *config.Config) error {
	// 生成服务
	dir := ctx.GetController()
	pbImport := fmt.Sprintf(`"%v"`, ctx.GetProto().Package)
	imports := []string{pbImport}

	for _, service := range proto.Service {
		controllerFilename, err := formatx.FileNamingFormat(g.style.Name, service.Name+"_controller")
		if err != nil {
			return err
		}

		// 导入business包
		businessImport := fmt.Sprintf(`"%s"`, ctx.GetBusiness().Package)
		imports = append(imports, businessImport)

		controllerFile := filepath.Join(dir.Filename, controllerFilename+".go")
		funcList, err := g.genFunctions(proto.PbPackage, service)
		if err != nil {
			return err
		}

		text, err := pathx.LoadTpl(category, controllerTemplateFilename, controllerTemplate)
		if err != nil {
			return err
		}

		notStream := false
		for _, rpc := range service.RPC {
			if !rpc.StreamsRequest && !rpc.StreamsReturns {
				notStream = true
				break
			}
		}
		err = templatex.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"pbPackage":    proto.PbPackage,
			"package":      "controller",
			"businessType": ctx.GetBusiness().Base + "." + stringx.NewString(service.Name).ToCamel() + "Business",
			"unimplementedServer": fmt.Sprintf("%s.Unimplemented%sServer", proto.PbPackage,
				stringx.NewString(service.Name).ToCamel()),
			"service":   stringx.NewString(service.Name).ToCamel(),
			"imports":   strings.Join(imports, "\n"),
			"funcs":     strings.Join(funcList, "\n"),
			"notStream": notStream,
		}, controllerFile, false)
		if err != nil {
			return err
		}
	}
	// 生成服务提供者
	return g.genControllerProvide(ctx, proto)
}

// genControllerProvide 生成service的服务提供者
func (g *Generator) genControllerProvide(ctx DirContext, proto model.Proto) error {
	dir := ctx.GetController()
	provideFilename, err := formatx.FileNamingFormat(g.style.Name, "controller")
	if err != nil {
		return err
	}
	serviceProvideFile := filepath.Join(dir.Filename, provideFilename+".go")
	provideSetList := make([]string, 0)
	for _, service := range proto.Service {
		provideSetList = append(provideSetList, fmt.Sprintf(`New%sController`, stringx.NewString(service.Name).ToCamel()))
	}
	text, err := pathx.LoadTpl(category, provideTemplateFilename, provideTemplate)
	if err != nil {
		return err
	}
	return templatex.With("service-provide").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"PackageName":   dir.Base,
		"ProvideSetStr": strings.Join(provideSetList, ","),
	}, serviceProvideFile, true)
}

// genFunctions 生成server方法
func (g *Generator) genFunctions(goPackage string, service model.Service) ([]string, error) {
	var (
		functionList []string
	)
	for _, rpc := range service.RPC {
		text, err := pathx.LoadTpl(category, serviceFuncTemplateFilename, functionTemplate)
		if err != nil {
			return nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(service.Name),
			parser.CamelCase(rpc.Name), "Server")
		buffer, err := templatex.With("func").Parse(text).Execute(map[string]interface{}{
			"service":      stringx.NewString(service.Name).ToCamel(),
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
			return nil, err
		}

		functionList = append(functionList, buffer.String())
	}
	return functionList, nil
}
