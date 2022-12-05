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

package zipx

import (
	"archive/zip"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"io"
	"os"
	"path/filepath"
)

// Unpacking 解压
func Unpacking(name, destPath string, mapper func(f *zip.File) bool) error {
	r, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		ok := mapper(file)
		if ok {
			err = fileCopy(file, destPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// fileCopy 文件复制
func fileCopy(file *zip.File, destPath string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	filename := filepath.Join(destPath, filepath.Base(file.Name))
	dir := filepath.Dir(filename)
	err = pathx.MkdirIfNotExist(dir)
	if err != nil {
		return err
	}

	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, rc)
	return err
}
