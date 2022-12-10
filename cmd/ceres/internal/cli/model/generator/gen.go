//    Copyright 2022. Go-Ceres
//    Author https://github.com/go-ceres/go-ceres
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
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/model/config"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/model/parser"
	"github.com/go-ceres/ceres/cmd/ceres/internal/ctx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/formatx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CodeDescribe struct {
	FileName string // 文件名
	Update   bool   // 是否允许更新
	Content  string // 文件内容
}

// GeneratorFromDDl 生成模型根据DDl描述文件
func (g *Generator) GeneratorFromDDl(file string, conf *config.ModelConfig) error {
	// 1.检查输出路径
	abs, err := filepath.Abs(conf.Dist)
	if err != nil {
		return err
	}

	// 2.获取项目信息
	projectCtx, err := ctx.PrepareProject(abs)
	if err != nil {
		return err
	}

	// 3.创建输出文件夹
	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	// 4.解析文件
	tables, err := parser.Parse(file, conf.DataBase, conf.Strict)
	if err != nil {
		return err
	}
	// 5.生成代码
	codes, err := g.genCodeDescribe(tables, projectCtx, conf)
	if err != nil {
		return err
	}
	// 6.写出代码
	return g.createFile(codes)
}

// genFromDDl 生成代码描述
func (g *Generator) genCodeDescribe(tables []*parser.Table, projectCtx *ctx.Project, conf *config.ModelConfig) ([]*CodeDescribe, error) {
	var res = make([]*CodeDescribe, 0)
	for _, table := range tables {
		code, err := g.genModel(*table, projectCtx, conf)
		if err != nil {
			return nil, err
		}
		res = append(res, code...)
	}
	return res, nil
}

// genModel 生成模型
func (g *Generator) genModel(table parser.Table, project *ctx.Project, conf *config.ModelConfig) ([]*CodeDescribe, error) {
	var res = make([]*CodeDescribe, 0)
	// 说明没有主键
	if len(table.Primary.Name.Source()) == 0 {
		return nil, fmt.Errorf("table %s: missing primary key", table.Name.Source())
	}

	// 生成存储器
	repository, err := g.genRepository(table, project, conf)
	if err != nil {
		return nil, err
	}
	// 生成自定义存储器
	res = append(res, repository)

	//customRepository, err := g.genCustomRepository(table, project, conf)
	//if err != nil {
	//	return nil, err
	//}
	//res = append(res, customRepository)
	return res, nil
}

// genRepository 生成存储器
func (g *Generator) genRepository(table parser.Table, projectCtx *ctx.Project, conf *config.ModelConfig) (*CodeDescribe, error) {
	res := new(CodeDescribe)

	// 构建包导入代码
	importsCode, err := g.genImports(table, conf.Cache, table.Time, projectCtx)
	if err != nil {
		return nil, err
	}

	// 构建结构体代码
	structCode, err := g.genStruct(table, projectCtx, conf)
	if err != nil {
		return nil, err
	}

	// 生成表名
	tableNameCode, err := g.genTableName(table, conf.Prefix)
	if err != nil {
		return nil, err
	}

	newCode, err := g.genNew(table, conf.Cache)
	if err != nil {
		return nil, err
	}
	// 生成获取Db代码
	getDbCode, err := g.genGetDb(table)
	if err != nil {
		return nil, err
	}

	// 数据迁移代码
	autoMigrateCode, err := g.genAutoMigrate(table)
	if err != nil {
		return nil, err
	}

	// 新增代码
	createCode, err := g.GenCreate(table)
	if err != nil {
		return nil, err
	}

	// 删除代码
	deleteCode, err := g.genDelete(table, projectCtx, conf)
	if err != nil {
		return nil, err
	}

	// 修改代码
	updateCode, err := g.GenUpdate(table, projectCtx, conf)
	if err != nil {
		return nil, err
	}

	// 查询一条代码
	findCode, err := g.genFind(table)
	if err != nil {
		return nil, err
	}

	queryCode, err := g.genQueryListBytSql(table)
	if err != nil {
		return nil, err
	}

	content, err := g.genModelGenCode(map[string]interface{}{
		"pkg":         filepath.Base(projectCtx.WorkDir),
		"imports":     importsCode,
		"types":       structCode,
		"tablename":   tableNameCode,
		"db":          getDbCode,
		"automigrate": autoMigrateCode,
		"find":        findCode,
		"new":         newCode,
		"create":      createCode,
		"update":      updateCode,
		"delete":      deleteCode,
		"query":       queryCode,
	})
	if err != nil {
		return nil, err
	}
	modelFilename, err := formatx.FileNamingFormat(g.style.Name,
		fmt.Sprintf("%s_repository", table.Name.Source()))
	if err != nil {
		return nil, err
	}
	res.Content = content
	res.Update = true
	res.FileName = filepath.Join(projectCtx.WorkDir, modelFilename+"_gen.go")
	return res, nil
}

// createFile 创建文件
func (g *Generator) createFile(codes []*CodeDescribe) error {
	for _, code := range codes {
		exists := pathx.FileExists(code.FileName)
		// 如果文件存在并且是不允许更新的情况下提示信息
		if exists && !code.Update {
			g.log.Warning("%s already exists, ignored.", code.FileName)
			continue
		}
		// 写入文件
		if err := ioutil.WriteFile(code.FileName, []byte(code.Content), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
