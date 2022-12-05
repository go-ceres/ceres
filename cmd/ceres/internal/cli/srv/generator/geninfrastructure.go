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
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/formatx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"path/filepath"
	"strings"
)

//go:embed tpl/provide.go.tpl
var infrastructureTemplate string

// genInfrastructure 生成基础设施层
func (g *Generator) genInfrastructure(ctx DirContext, conf *config.Config) error {
	dir := ctx.GetInfrastructure()
	err := pathx.MkdirIfNotExist(dir.Filename)
	if err != nil {
		return err
	}
	pkgPackage, err := dir.GetChildPackage("pkg")
	if err != nil {
		return err
	}
	pkgImport := fmt.Sprintf(`"%s"`, pkgPackage)
	repositoryPackage, err := dir.GetChildPackage("repository")
	if err != nil {
		return err
	}
	repositoryImport := fmt.Sprintf(`"%s"`, repositoryPackage)
	imports := []string{pkgImport, repositoryImport}

	infrastructureFileName, err := formatx.FileNamingFormat(g.style.Name, "infrastructure")
	if err != nil {
		return err
	}
	pkgSet := fmt.Sprintf("%s.ProvideSet", filepath.Base(pkgPackage))
	repositorySet := fmt.Sprintf("%s.ProvideSet", filepath.Base(repositoryPackage))
	provideSet := []string{pkgSet, repositorySet}

	fileName := filepath.Join(dir.Filename, infrastructureFileName+".go")
	text, err := pathx.LoadTpl(category, infrastructureTemplateFileName, infrastructureTemplate)
	if err != nil {
		return err
	}
	return templatex.With("infrastructure").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"PackageName":   "infrastructure",
		"ImportsStr":    strings.Join(imports, "\n"),
		"ProvideSetStr": strings.Join(provideSet, ","),
	}, fileName, false)
}
