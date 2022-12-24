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

package file

import (
	"github.com/go-ceres/ceres/pkg/common/config"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type file struct {
	path string
}

// NewSource 创建文件资源管理器
func NewSource(path string) config.Source {
	return &file{path: path}
}

// Load 加载资源
func (f *file) Load() ([]*config.DataSet, error) {
	stat, err := os.Stat(f.path)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return f.loadDir(f.path, make([]*config.DataSet, 0))
	}
	dataSet, err := f.loadFile(f.path)
	if err != nil {
		return nil, err
	}
	return []*config.DataSet{dataSet}, nil
}

// Watch 启动文件监控
func (f *file) Watch() (config.Watcher, error) {
	return newWatcher(f)
}

// loadDir 加载文件夹文件
func (f *file) loadDir(path string, dataSets []*config.DataSet) ([]*config.DataSet, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range files {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		dataSet, err := f.loadFile(filepath.Join(path, entry.Name()))
		if err != nil {
			return nil, err
		}
		dataSets = append(dataSets, dataSet)
	}
	return dataSets, nil
}

// loadFile 加载文件
func (f *file) loadFile(path string) (*config.DataSet, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &config.DataSet{
		Key:       info.Name(),
		Format:    strings.TrimPrefix(filepath.Ext(info.Name()), "."),
		Data:      data,
		Timestamp: info.ModTime(),
	}, nil
}
