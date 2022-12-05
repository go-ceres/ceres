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
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/execx"
	"path/filepath"
)

// Project 项目解析（go mod）
type Project struct {
	WorkDir string
	Name    string
	Dir     string
	Path    string
}

// PrepareProject 解析项目信息
func PrepareProject(workDir string) (*Project, error) {
	ctx, err := GetProjectInfo(workDir)
	if err == nil {
		return ctx, nil
	}
	name := filepath.Base(workDir)
	_, err = execx.Command("go mod init "+name, workDir)
	if err != nil {
		return nil, err
	}
	return GetProjectInfo(workDir)
}

// GetProjectInfo 获取项目信息
func GetProjectInfo(workDir string) (*Project, error) {
	isGoMod, err := IsGoMod(workDir)
	if err != nil {
		return nil, err
	}

	if isGoMod {
		return projectFromGoMod(workDir)
	}
	return projectFromGoPath(workDir)
}
