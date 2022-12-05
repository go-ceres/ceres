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

package ctx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/execx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const goModuleWithoutGoFiles = "command-line-arguments"

var errInvalidGoMod = errors.New("invalid go module")

// Module 模块的配置信息
type Module struct {
	Path      string
	Main      bool
	Dir       string
	GoMod     string
	GoVersion string
}

func (m *Module) validate() error {
	if m.Path == goModuleWithoutGoFiles || m.Dir == "" {
		return errInvalidGoMod
	}
	return nil
}

// IsGoMod 判断是否是mod
func IsGoMod(workDir string) (bool, error) {
	if len(workDir) == 0 {
		return false, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return false, err
	}

	data, err := execx.Command("go list -m -f '{{.GoMod}}'", workDir)
	if err != nil || len(data) == 0 {
		return false, nil
	}

	return true, nil
}

// GetRealModule 获取指定工作目录下的模块信息
func GetRealModule(workDir string) (*Module, error) {
	data, err := execx.Command("go list -json -m", workDir)
	if err != nil {
		return nil, err
	}
	modules, err := decodePackages(strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	for _, m := range modules {
		if strings.HasPrefix(workDir, m.Dir) {
			return &m, nil
		}
	}
	return nil, errors.New("no matched module")
}

// decodePackages 解析包
func decodePackages(rc io.Reader) ([]Module, error) {
	var modules []Module
	decoder := json.NewDecoder(rc)
	for decoder.More() {
		var m Module
		if err := decoder.Decode(&m); err != nil {
			return nil, fmt.Errorf("invalid module: %v", err)
		}
		modules = append(modules, m)
	}

	return modules, nil
}

// projectFromGoMod 从mod文件中解析项目信息
func projectFromGoMod(workDir string) (*Project, error) {
	if len(workDir) == 0 {
		return nil, errors.New("the work directory is not found")
	}
	workDir, err := pathx.ReadLink(workDir)
	if err != nil {
		return nil, err
	}

	m, err := GetRealModule(workDir)
	if err != nil {
		return nil, err
	}

	if err := m.validate(); err != nil {
		return nil, err
	}

	var res Project
	res.WorkDir = workDir
	res.Name = filepath.Base(m.Dir)
	dir, err := pathx.ReadLink(m.Dir)
	if err != nil {
		return nil, err
	}

	res.Dir = dir
	res.Path = m.Path
	return &res, nil
}
