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

package matcher

import (
	"github.com/go-ceres/ceres/pkg/transport"
	"sort"
	"strings"
)

// Matcher 中间件匹配器
type Matcher interface {
	Use(mw ...transport.Middleware)
	Add(selector string, mw ...transport.Middleware)
	Match(operator string) []transport.Middleware
}

// New 创建
func New() Matcher {
	return &matcher{
		matchs: make(map[string][]transport.Middleware),
	}
}

type matcher struct {
	prefix []string
	data   []transport.Middleware
	matchs map[string][]transport.Middleware
}

func (m *matcher) Use(ms ...transport.Middleware) {
	m.data = ms
}

// Add 添加中间件到匹配器
func (m *matcher) Add(selector string, mw ...transport.Middleware) {
	if strings.HasPrefix(selector, "*") {
		selector = strings.TrimPrefix(selector, "*")
		m.prefix = append(m.prefix, selector)
		sort.Slice(m.prefix, func(i, j int) bool {
			return m.prefix[i] > m.prefix[j]
		})
	}
	m.matchs[selector] = mw
}

// Match 匹配
func (m *matcher) Match(operation string) []transport.Middleware {
	ms := make([]transport.Middleware, 0, len(m.data))
	if len(m.data) > 0 {
		ms = append(ms, m.data...)
	}
	if next, ok := m.matchs[operation]; ok {
		return append(ms, next...)
	}
	for _, prefix := range m.prefix {
		if strings.HasPrefix(operation, prefix) {
			return append(ms, m.matchs[prefix]...)
		}
	}
	return ms
}
