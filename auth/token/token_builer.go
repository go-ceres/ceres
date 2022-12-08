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

package token

import (
	"github.com/go-ceres/ceres/internal/strings"
	"github.com/google/uuid"
	s "strings"
)

var _ TokenBuilder = (*defaultTokenBuilder)(nil)

// TokenBuilder token的生成接口
type TokenBuilder interface {
	Build(loginId string, logicType string, device string) string
}

type defaultTokenBuilder struct {
	style TokenStyle
}

func (d *defaultTokenBuilder) Build(loginId string, logicType string, device string) string {
	style := d.style
	switch style {
	case TokenStyleUuid:
		return uuid.New().String()
	case TokenStyleSimpleUuid:
		return s.ReplaceAll(uuid.New().String(), "_", "")
	case TokenStyleRandom32:
		return strings.RandStr(32)
	case TokenStyleRandom64:
		return strings.RandStr(64)
	default:
		return uuid.New().String()
	}
}
