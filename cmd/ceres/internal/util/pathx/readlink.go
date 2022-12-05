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

package pathx

import (
	"os"
	"path/filepath"
)

// ReadLink 递归返回命名符号链接的目标。
func ReadLink(name string) (string, error) {
	name, err := filepath.Abs(name)
	if err != nil {
		return "", err
	}

	if _, err := os.Lstat(name); err != nil {
		return name, nil
	}

	// uncheck condition: ignore file path /var, maybe be temporary file path
	if name == "/" || name == "/var" {
		return name, nil
	}

	isLink, err := isLink(name)
	if err != nil {
		return "", err
	}

	if !isLink {
		dir, base := filepath.Split(name)
		dir = filepath.Clean(dir)
		dir, err := ReadLink(dir)
		if err != nil {
			return "", err
		}

		return filepath.Join(dir, base), nil
	}

	link, err := os.Readlink(name)
	if err != nil {
		return "", err
	}

	dir, base := filepath.Split(link)
	dir = filepath.Dir(dir)
	dir, err = ReadLink(dir)
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, base), nil
}
