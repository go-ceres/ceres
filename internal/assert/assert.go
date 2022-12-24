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

package assert

import (
	"fmt"
	"github.com/go-ceres/ceres/pkg/common/logger"
)

func Debugf(guard bool, text string, args ...interface{}) {
	if !guard {
		logger.Debugf(text, args...)
	}
}

func Infof(guard bool, text string, args ...interface{}) {
	if !guard {
		logger.Infof(text, args...)
	}
}

func Warnf(guard bool, text string, args ...interface{}) {
	if !guard {
		logger.Warnf(text, args...)
	}
}

func Errorf(guard bool, text string, args ...interface{}) {
	if !guard {
		logger.Errorf(text, args...)
	}
}

func Panic(guard bool, text string, args ...interface{}) {
	if !guard {
		panic(fmt.Sprintf(text, args...))
	}
}
