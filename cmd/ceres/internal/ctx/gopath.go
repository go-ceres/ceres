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

package ctx

import (
	"errors"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

var errModuleCheck = errors.New("the work directory must be found in the go mod or the $GOPATH")

// projectFromGoPath 从工作路径中获取项目信息
func projectFromGoPath(workDir string) (*Project, error) {
	if len(workDir) == 0 {
		return nil, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return nil, err
	}

	workDir, err := pathx.ReadLink(workDir)
	if err != nil {
		return nil, err
	}

	buildContext := build.Default
	goPath := buildContext.GOPATH
	goPath, err = pathx.ReadLink(goPath)
	if err != nil {
		return nil, err
	}

	goSrc := filepath.Join(goPath, "src")
	if !pathx.FileExists(goSrc) {
		return nil, errModuleCheck
	}

	wd, err := filepath.Abs(workDir)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(wd, goSrc) {
		return nil, errModuleCheck
	}

	projectName := strings.TrimPrefix(wd, goSrc+string(filepath.Separator))
	return &Project{
		WorkDir: workDir,
		Name:    projectName,
		Path:    projectName,
		Dir:     filepath.Join(goSrc, projectName),
	}, nil
}
