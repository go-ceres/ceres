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
	"path/filepath"
	"strings"
)

const serviceFunctionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.service}}Service) {{.method}} ({{if .notStream}}ctx context.Context,{{if .hasReq}} req {{.request}}{{end}}{{else}}{{if .hasReq}} req {{.request}},{{end}}stream {{.streamBody}}{{end}}) ({{if .notStream}}{{.response}},{{end}}error) {
	return s.{{.actionDesc.UnTitleName}}Action.{{.method}}({{if .notStream}}ctx,{{if .hasReq}}req{{end}}{{else}}{{if .hasReq}} req,{{end}}stream{{end}})
}
`

//go:embed tpl/service.go.tpl
var serviceTemplate string

//go:embed tpl/upload.go.tpl
var uploadServiceTemplate string

type ActionDesc struct {
	ActionName    string
	UnTitleName   string
	ActionPackage string
}

// GenService 生成grpc实现类,业务入口
func (g *Generator) GenService(ctx DirContext, proto model.Proto, conf *config.Config) error {
	// 生成服务
	dir := ctx.GetService()
	pbImport := fmt.Sprintf(`"%v"`, ctx.GetProto().Package)
	for _, service := range proto.Service {
		serviceFilename, err := formatx.FileNamingFormat(g.style.Name, service.Name+"_service")
		imports := []string{pbImport}
		if err != nil {
			return err
		}
		// 如果是多服务
		if len(proto.Service) > 1 {
			actionChildPkg, err := ctx.GetAction().GetChildPackage(service.Name)
			if err != nil {
				return err
			}
			actionImport := fmt.Sprintf(`%s "%s"`, service.Name+"Actions", actionChildPkg)
			imports = append(imports, actionImport)

		} else {
			// 导入action包
			actionImport := fmt.Sprintf(`"%s"`, ctx.GetAction().Package)
			imports = append(imports, actionImport)
		}
		serviceFile := filepath.Join(dir.Filename, serviceFilename+".go")
		funcList, descList, err := g.genServiceFunctions(ctx, proto.PbPackage, service, len(proto.Service) > 1)
		if err != nil {
			return err
		}

		text, err := pathx.LoadTpl(category, serviceTemplateFilename, serviceTemplate)
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
		head := templatex.GetHead(proto.Name)
		err = templatex.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"head":      head,
			"pbPackage": proto.PbPackage,
			"package":   "service",
			"unimplementedServer": fmt.Sprintf("%s.Unimplemented%sServer", proto.PbPackage,
				stringx.NewString(service.Name).ToCamel()),
			"service":   stringx.NewString(service.Name).ToCamel(),
			"imports":   strings.Join(imports, "\n"),
			"funcs":     strings.Join(funcList, "\n"),
			"descList":  descList,
			"notStream": notStream,
		}, serviceFile, true)
		if err != nil {
			return err
		}
	}
	// 生成上传服务
	if err := g.GenUploadService(ctx, proto, conf); err != nil {
		return err
	}
	// 生成服务提供者
	return g.genServiceProvide(ctx, proto)
}

// GenUploadService 生成上传服务文件
func (g *Generator) GenUploadService(ctx DirContext, proto model.Proto, conf *config.Config) error {
	// 生成服务
	dir := ctx.GetService()
	uploadServiceFilename := filepath.Join(dir.Filename, "upload_service.go")
	text, err := pathx.LoadTpl(category, serviceTemplateFilename, uploadServiceTemplate)
	if err != nil {
		return err
	}
	return templatex.With("upload_service").
		GoFmt(true).
		Parse(text).
		SaveTo(map[string]interface{}{}, uploadServiceFilename, false)
}

// genServiceProvide 生成service的服务提供者
func (g *Generator) genServiceProvide(ctx DirContext, proto model.Proto) error {
	dir := ctx.GetService()
	provideFilename, err := formatx.FileNamingFormat(g.style.Name, "service")
	if err != nil {
		return err
	}
	serviceProvideFile := filepath.Join(dir.Filename, provideFilename+".go")
	provideSetList := make([]string, 0)
	for _, service := range proto.Service {
		provideSetList = append(provideSetList, fmt.Sprintf(`New%sService`, stringx.NewString(service.Name).ToCamel()))
	}
	// 增加upload_service
	provideSetList = append(provideSetList, `NewUploadService`)
	text, err := pathx.LoadTpl(category, provideTemplateFilename, provideTemplate)
	if err != nil {
		return err
	}
	return templatex.With("service-provide").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"PackageName":   dir.Base,
		"ProvideSetStr": strings.Join(provideSetList, ","),
	}, serviceProvideFile, true)
}

// genServiceFunctions 生成server方法
func (g *Generator) genServiceFunctions(ctx DirContext, goPackage string, service model.Service, multi bool) ([]string, []*ActionDesc, error) {
	var (
		functionList   []string
		actionDescList []*ActionDesc
	)
	for _, rpc := range service.RPC {
		text, err := pathx.LoadTpl(category, serviceFuncTemplateFilename, serviceFunctionTemplate)
		if err != nil {
			return nil, nil, err
		}
		comment := parser.GetComment(rpc.Doc())
		streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(service.Name),
			parser.CamelCase(rpc.Name), "Server")
		actionPackage := ""
		if multi {
			actionPackage = service.Name + "Actions"
		} else {
			actionPackage = "action"
		}
		desc := &ActionDesc{
			ActionName:    stringx.NewString(rpc.Name).ToCamel() + "Action",
			UnTitleName:   stringx.NewString(stringx.NewString(rpc.Name).ToCamel()).UnTitle(),
			ActionPackage: actionPackage,
		}
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
			"actionDesc":   desc,
		})
		if err != nil {
			return nil, nil, err
		}

		functionList = append(functionList, buffer.String())
		actionDescList = append(actionDescList, desc)
	}
	return functionList, actionDescList, nil
}
