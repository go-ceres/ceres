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
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/config"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/formatx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"path/filepath"
	"strings"
)

//go:embed tpl/bootstrap.go.tpl
var bootstrapTemplate string

// genWire 生成wire依赖注入文件
func (g *Generator) genBootstrap(ctx DirContext, conf *config.Config) error {
	dir := ctx.GetBootstrap()
	imports := make([]string, 0)
	bootstrapFilename, err := formatx.FileNamingFormat(g.style.Name, "main")
	if err != nil {
		return err
	}
	bootstrapFile := filepath.Join(dir.Filename, bootstrapFilename+".go")
	text, err := pathx.LoadTpl(category, bootstrapTemplateFilename, bootstrapTemplate)
	if err != nil {
		return err
	}
	return templatex.With("bootstrap").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"PackageImports": strings.Join(imports, "\n"),
		"hasRegistry":    conf.Registry,
		"HttpServer":     conf.HttpServer,
	}, bootstrapFile, false)
}
