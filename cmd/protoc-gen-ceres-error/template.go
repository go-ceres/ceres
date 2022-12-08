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

package main

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed errors.go.tpl
var errorsTemplate string

type errorInfo struct {
	Name       string
	Value      string
	HTTPCode   int
	CamelValue string
	Comment    string
	HasComment bool
}

type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errors").Parse(errorsTemplate)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return buf.String()
}
