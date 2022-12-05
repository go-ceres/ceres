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
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser/model"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/formatx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"path/filepath"
	"strings"
)

//go:embed tpl/wire.go.tpl
var wireTemplate string

// genWire 生成wire依赖注入文件
func (g *Generator) genWire(ctx DirContext, proto model.Proto) error {
	// 生成服务
	dir := ctx.GetBootstrap()
	controllerImport := fmt.Sprintf(`"%s"`, ctx.GetController().Package)
	serverImport := fmt.Sprintf(`"%s"`, ctx.GetServer().Package)
	domainImport := fmt.Sprintf(`"%s"`, ctx.GetDomain().Package)
	infrastructureImport := fmt.Sprintf(`"%s"`, ctx.GetInfrastructure().Package)
	imports := []string{controllerImport, serverImport, domainImport, infrastructureImport}
	wireFilename, err := formatx.FileNamingFormat(g.style.Name, "wire")
	if err != nil {
		return err
	}
	wireFile := filepath.Join(dir.Filename, wireFilename+".go")
	text, err := pathx.LoadTpl(category, wireTemplateFilename, wireTemplate)
	if err != nil {
		return err
	}
	return templatex.With("wire").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"ImportsStr": strings.Join(imports, "\n"),
	}, wireFile, false)
}
