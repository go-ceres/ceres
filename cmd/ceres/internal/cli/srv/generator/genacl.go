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
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/config"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/formatx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"path/filepath"
)

func (g *Generator) genAcl(ctx DirContext, conf *config.Config) error {
	// 生成服务
	dir := ctx.GetAcl()
	filename, err := formatx.FileNamingFormat(g.style.Name, "acl")
	if err != nil {
		return err
	}
	serverFile := filepath.Join(dir.Filename, filename+".go")
	text, err := pathx.LoadTpl(category, provideTemplateFilename, provideTemplate)
	if err != nil {
		return err
	}
	ProvideSetStr := ""
	return templatex.With("acl-provide").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"PackageName":   "acl",
		"ProvideSetStr": ProvideSetStr,
	}, serverFile, false)
}
