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
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/templatex"
	"path/filepath"
)

//go:embed tpl/config.toml.tpl
var configTemplate string

// GenConfig 生成配置文件
func (g *Generator) GenConfig(ctx DirContext, _ model.Proto, conf *config.Config) error {
	dir := ctx.GetConfigs()
	configFilename, err := formatx.FileNamingFormat(g.style.Name, "config")
	if err != nil {
		return err
	}
	err = pathx.MkdirIfNotExist(dir.Filename)
	if err != nil {
		return err
	}
	// 文件名
	fileName := filepath.Join(dir.Filename, fmt.Sprintf("%v.toml", configFilename))
	// 获取模板内容
	context, err := pathx.LoadTpl(category, configTemplateFilename, configTemplate)
	if err != nil {
		return err
	}
	return templatex.With("config").Parse(context).SaveTo(map[string]interface{}{
		"serviceName": ctx.GetServiceName().UnTitle(),
		"HttpServer":  conf.HttpServer,
		"components":  conf.Components,
	}, fileName, false)
}
