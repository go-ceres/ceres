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
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/config"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser"
	"github.com/go-ceres/ceres/cmd/ceres/internal/ctx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"path/filepath"
)

func (g *Generator) Generate(conf *config.Config) error {
	// 1.检查输出路径
	abs, err := filepath.Abs(conf.Dist)
	if err != nil {
		return err
	}

	// 2.创建输出文件夹
	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	// 3.检查环境是否安装
	err = g.Prepare()
	if err != nil {
		return err
	}

	// 4.获取项目信息
	projectCtx, err := ctx.PrepareProject(abs)
	if err != nil {
		return err
	}

	// 5.翻译proto原始文件为结构体
	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse(conf.ProtoFile, true)
	if err != nil {
		return err
	}

	// 6.创建文件夹
	dirCtx, err := g.mkdir(projectCtx, proto, conf)
	if err != nil {
		return err
	}

	// 7.生成pb文件
	if err := g.GenPb(dirCtx, conf); err != nil {
		return err
	}

	// 8.生成配置文件
	if err := g.GenConfig(dirCtx, proto, conf); err != nil {
		return err
	}

	// 9.生成数据存储接口
	if err := g.genIRepository(dirCtx, proto); err != nil {
		return err
	}

	// 10.生成实体对象
	if err := g.genEntity(dirCtx, proto); err != nil {
		return err
	}
	// 11.生成业务business业务代码
	if err := g.genBusiness(dirCtx, proto); err != nil {
		return err
	}
	// 12.生成基础设施层存储层
	if err := g.genRepository(dirCtx, proto); err != nil {
		return err
	}

	// 13.生成domain的provide
	if err := g.genDomainProvide(dirCtx, proto); err != nil {
		return err
	}

	// 14.生成基础设施层的包依赖初始化
	if err := g.genPkg(dirCtx, conf); err != nil {
		return err
	}

	// 15.生成基础设施的provide
	if err := g.genInfrastructure(dirCtx, conf); err != nil {
		return err
	}

	// 16.生成服务
	if err := g.GenController(dirCtx, proto, conf); err != nil {
		return err
	}

	// 17.生成服务注册
	if err := g.genServer(dirCtx, proto, conf); err != nil {
		return err
	}

	// 18.生成acl，接口防腐层依赖注入
	if err := g.genAcl(dirCtx, conf); err != nil {
		return err
	}

	// 19.生成依赖注入文件
	if err := g.genWire(dirCtx, proto); err != nil {
		return err
	}

	// 20.生成bootstrap文件
	if err = g.genBootstrap(dirCtx, conf); err != nil {
		return err
	}

	// 21.生成makefile文件
	if err := g.genMakefile(dirCtx, proto); err != nil {
		return err
	}
	// 22.生成dockerfile文件
	if err := g.genDockerfile(dirCtx, proto); err != nil {
		return err
	}
	return nil
}
