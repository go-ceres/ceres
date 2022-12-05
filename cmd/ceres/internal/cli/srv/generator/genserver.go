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
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser/model"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/formatx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/stringx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"path/filepath"
	"strings"
)

//go:embed tpl/grpc.go.tpl
var grpcServerTemplate string

//go:embed tpl/http.go.tpl
var httpServerTemplate string

// genServer 生成服务注册代码
func (g *Generator) genServer(ctx DirContext, proto model.Proto, conf *config.Config) error {
	if err := g.genGRPCServer(ctx, proto); err != nil {
		return err
	}
	if len(conf.HttpServer) > 0 {
		if err := g.genHTTPServer(ctx, proto, conf); err != nil {
			return err
		}
	}
	return g.genServerProvide(ctx, proto, conf)
}

// genGRPCServer 生成GRPC注册服务
func (g *Generator) genGRPCServer(ctx DirContext, proto model.Proto) error {
	// 生成服务
	dir := ctx.GetServer()
	pbImport := fmt.Sprintf(`proto "%v"`, ctx.GetProto().Package)
	serviceImport := fmt.Sprintf(`"%s"`, ctx.GetController().Package)
	imports := []string{pbImport, serviceImport}
	filename, err := formatx.FileNamingFormat(g.style.Name, "grpc")
	if err != nil {
		return err
	}
	serverFile := filepath.Join(dir.Filename, filename+".go")
	var serverParamsList = make([]string, 0)
	var registerServerList = make([]string, 0)
	for _, service := range proto.Service {
		paramName := stringx.NewString(stringx.NewString(service.Name).ToCamel()).UnTitle()
		paramType := stringx.NewString(service.Name).ToCamel() + "Controller"
		serverParamsList = append(serverParamsList, fmt.Sprintf("%s *%s.%s", paramName, ctx.GetController().Base, paramType))
		registerServerList = append(
			registerServerList,
			fmt.Sprintf("%s.Register%s(srv,%s)", "proto", stringx.NewString(service.Name).ToCamel()+"Server", paramName),
		)
	}

	text, err := pathx.LoadTpl(category, grpcServerTemplateFilename, grpcServerTemplate)
	if err != nil {
		return err
	}
	return templatex.With("grpc-server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"Imports":         strings.Join(imports, "\n"),
		"serverParamsStr": strings.Join(serverParamsList, ","),
		"registerListStr": strings.Join(registerServerList, "\n"),
	}, serverFile, true)
}

// genHTTPServer 生成HTTP服务
func (g *Generator) genHTTPServer(ctx DirContext, proto model.Proto, conf *config.Config) error {
	// 生成服务
	dir := ctx.GetServer()
	pbImport := fmt.Sprintf(`proto "%v"`, ctx.GetProto().Package)
	serviceImport := fmt.Sprintf(`"%s"`, ctx.GetController().Package)
	serverImport := fmt.Sprintf(`"%s"`, "github.com/go-ceres/ceres/server/"+conf.HttpServer)
	imports := []string{pbImport, serviceImport, serverImport}
	filename, err := formatx.FileNamingFormat(g.style.Name, "http")
	if err != nil {
		return err
	}
	serverFile := filepath.Join(dir.Filename, filename+".go")
	var serverParamsList = make([]string, 0)
	var registerServerList = make([]string, 0)
	for _, service := range proto.Service {
		paramName := stringx.NewString(stringx.NewString(service.Name).ToCamel()).UnTitle()
		paramType := stringx.NewString(service.Name).ToCamel() + "Controller"
		serverParamsList = append(serverParamsList, fmt.Sprintf("%s *%s.%s", paramName, ctx.GetController().Base, paramType))
		registerServerList = append(
			registerServerList,
			fmt.Sprintf("%s.Register%s(srv,%s)", "proto", stringx.NewString(service.Name).ToCamel()+"HTTPServer", paramName),
		)
	}
	text, err := pathx.LoadTpl(category, httpServerTemplateFilename, httpServerTemplate)
	if err != nil {
		return err
	}
	return templatex.With("http-server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"Imports":         strings.Join(imports, "\n"),
		"ServerType":      conf.HttpServer,
		"serverParamsStr": strings.Join(serverParamsList, ","),
		"registerListStr": strings.Join(registerServerList, "\n"),
	}, serverFile, true)
}

// genServerProvide 生成服务提供者
func (g *Generator) genServerProvide(ctx DirContext, proto model.Proto, conf *config.Config) error {
	// 生成服务
	dir := ctx.GetServer()
	filename, err := formatx.FileNamingFormat(g.style.Name, "server")
	if err != nil {
		return err
	}
	serverFile := filepath.Join(dir.Filename, filename+".go")
	text, err := pathx.LoadTpl(category, provideTemplateFilename, provideTemplate)
	if err != nil {
		return err
	}
	ProvideSetStr := "NewGRPCServer"
	if len(conf.HttpServer) > 0 {
		ProvideSetStr = ProvideSetStr + ",NewHTTPServer"
	}
	return templatex.With("server-provide").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"PackageName":   "server",
		"ProvideSetStr": ProvideSetStr,
	}, serverFile, true)
}
