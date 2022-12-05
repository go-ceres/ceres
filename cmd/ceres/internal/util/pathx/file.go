//    Copyright 2022. ceres
//    Author https://github.com/go-ceres/ceres
//
//    Licensed under the Apache License, CeresVersion 2.0 (the "License");
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

package pathx

import (
	"github.com/go-ceres/ceres/cmd/ceres/internal/version"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	ceresDir = ".ceres"
	gitDir   = ".git"
	cacheDir = "cache"
)

var ceresHome string

// RegisterCeresHome 注册home目录
func RegisterCeresHome(home string) {
	ceresHome = home
}

// MatchFiles 搜索匹配文件
func MatchFiles(src string) ([]string, error) {
	dir, pattern := filepath.Split(src)
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(abs)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		match, err := filepath.Match(pattern, name)
		if err != nil {
			return nil, err
		}

		if !match {
			continue
		}

		res = append(res, filepath.Join(abs, name))
	}
	return res, nil
}

// MkdirIfNotFound 如果文件夹不存在则创建文件夹
func MkdirIfNotFound(dir string) error {
	if len(dir) == 0 {
		return nil
	}
	// 创建文件夹
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// FileExists 判断文件是否存在
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// GetGitHome 获取ceres的git home目录
func GetGitHome() (string, error) {
	ceresHome, err := GetCeresHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(ceresHome, gitDir), nil
}

// GetCacheDir 获取缓存目录
func GetCacheDir() (string, error) {
	ceresHome, err := GetCeresHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(ceresHome, cacheDir), nil
}

func SameFile(path1, path2 string) (bool, error) {
	stat1, err := os.Stat(path1)
	if err != nil {
		return false, err
	}

	stat2, err := os.Stat(path2)
	if err != nil {
		return false, err
	}

	return os.SameFile(stat1, stat2), nil
}

// Copy 复制
func Copy(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	dir := filepath.Dir(dest)
	err = MkdirIfNotExist(dir)
	if err != nil {
		return err
	}
	w, err := os.Create(dest)
	if err != nil {
		return err
	}
	_ = w.Chmod(os.ModePerm)
	defer func(w *os.File) {
		_ = w.Close()
	}(w)
	_, err = io.Copy(w, f)
	return err
}

// GetCeresHome 获取ceres的home目录
func GetCeresHome() (home string, err error) {
	defer func() {
		if err != nil {
			return
		}
		info, err := os.Stat(home)
		if err == nil && !info.IsDir() {
			_ = os.Rename(home, home+".old")
			_ = MkdirIfNotExist(home)
		}
	}()
	if len(ceresHome) != 0 {
		home = ceresHome
		return
	}
	home, err = GetDefaultCeresHome()
	return
}

// GetDefaultCeresHome 获取默认的ceres的home目录
func GetDefaultCeresHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ceresDir), nil
}

// LoadTpl 加载模板
func LoadTpl(category, filename, def string) (string, error) {
	dir, err := GetTplDir(category)
	if err != nil {
		return "", err
	}

	filename = filepath.Join(dir, filename)
	if !FileExists(filename) {
		return def, nil
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// GetTplDir 获取指定模板路径
func GetTplDir(category string) (string, error) {
	home, err := GetCeresHome()
	if err != nil {
		return "", err
	}
	if home == ceresHome {
		beforeTplDir := filepath.Join(home, version.Version, category)
		fs, _ := ioutil.ReadDir(beforeTplDir)
		var hasContent bool
		for _, f := range fs {
			if f.Size() > 0 {
				hasContent = true
				break
			}
		}
		if hasContent {
			return beforeTplDir, nil
		}
		return filepath.Join(home, category), nil
	}
	return filepath.Join(home, version.Version, category), nil
}
