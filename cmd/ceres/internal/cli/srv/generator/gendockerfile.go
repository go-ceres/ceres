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
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser/model"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"path/filepath"
)

//go:embed tpl/dockerfile.dockerfile.tpl
var dockerfileTemplate string

// genWire 生成wire依赖注入文件
func (g *Generator) genDockerfile(ctx DirContext, proto model.Proto) error {
	dir := ctx.GetDeployment()
	mainFile := filepath.Join(dir.Filename, "Dockerfile")
	text, err := pathx.LoadTpl(category, dockerfileTemplateFilename, dockerfileTemplate)
	if err != nil {
		return err
	}
	return templatex.With("dockerfile").Parse(text).SaveTo(map[string]interface{}{
		"Name": ctx.GetServiceName().Source(),
	}, mainFile, false)
}
