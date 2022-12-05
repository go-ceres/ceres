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

package style

import (
	"errors"
	"strings"
)

const DefaultFormat = "go_ceres"

// Style 样式
type Style struct {
	Name string // 生成样式
}

func NewStyle(format string) (*Style, error) {
	if len(format) == 0 {
		format = DefaultFormat
	}
	cfg := &Style{
		Name: format,
	}
	err := validate(cfg)
	return cfg, err
}

func validate(cfg *Style) error {
	if len(strings.TrimSpace(cfg.Name)) == 0 {
		return errors.New("missing namingFormat")
	}
	return nil
}
