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
	"fmt"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser/model"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/formatx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"path/filepath"
	"strings"
)

func (g *Generator) genDomainProvide(ctx DirContext, proto model.Proto) error {
	dir := ctx.GetDomain()
	provideFilename, err := formatx.FileNamingFormat(g.style.Name, "domain")
	if err != nil {
		return err
	}
	repositoryProvideFile := filepath.Join(dir.Filename, provideFilename+".go")
	businessImport := fmt.Sprintf(`"%s"`, ctx.GetBusiness().Package)
	imports := []string{businessImport}
	text, err := pathx.LoadTpl(category, provideTemplateFilename, provideTemplate)
	if err != nil {
		return err
	}
	return templatex.With("business-provide").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"PackageName":   dir.Base,
		"ImportsStr":    strings.Join(imports, "\n"),
		"ProvideSetStr": strings.Join([]string{ctx.GetBusiness().Base + ".ProvideSet"}, ","),
	}, repositoryProvideFile, true)
}
