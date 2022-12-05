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

package errorx

import (
	"fmt"
	"github.com/go-ceres/ceres/cmd/ceres/internal/environment"
	"strings"
)

var errorFormat = `ceres error: %+v
ceres env:
%s
%s`

// CeresError represents a ceres error.
type CeresError struct {
	message []string
	err     error
}

func (e *CeresError) Error() string {
	detail := wrapMessage(e.message...)
	return fmt.Sprintf(errorFormat, e.err, environment.Print(), detail)
}

// Wrap wraps an error with ceres version and message.
func Wrap(err error, message ...string) error {
	e, ok := err.(*CeresError)
	if ok {
		return e
	}

	return &CeresError{
		message: message,
		err:     err,
	}
}

func wrapMessage(message ...string) string {
	if len(message) == 0 {
		return ""
	}
	return fmt.Sprintf(`message: %s`, strings.Join(message, "\n"))
}
