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

package templatex

import (
	"bytes"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/errorx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	goformat "go/format"
	"io/ioutil"
	"text/template"
)

const regularPerm = 0o666

const (
	// DoNotEditHead added to the beginning of a file to prompt the user not to edit
	DoNotEditHead = "// Code generated by ceres. DO NOT EDIT."

	headTemplate = DoNotEditHead + `
// Source: {{.source}}`
)

// Template is a tool to provides the text/template operations
type Template struct {
	name  string
	text  string
	goFmt bool
}

// With returns an instance of Template
func With(name string) *Template {
	return &Template{
		name: name,
	}
}

// Parse accepts a source template and returns Template
func (t *Template) Parse(text string) *Template {
	t.text = text
	return t
}

// GoFmt sets the value to goFmt and marks the generated codes will be formatted or not
func (t *Template) GoFmt(format bool) *Template {
	t.goFmt = format
	return t
}

// SaveTo writes the codes to the target path
func (t *Template) SaveTo(data interface{}, path string, forceUpdate bool) error {
	if pathx.FileExists(path) && !forceUpdate {
		return nil
	}

	output, err := t.Execute(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, output.Bytes(), regularPerm)
}

// Execute returns the codes after the template executed
func (t *Template) Execute(data interface{}) (*bytes.Buffer, error) {
	tem, err := template.New(t.name).Parse(t.text)
	if err != nil {
		return nil, errorx.Wrap(err, "template parse error:", t.text)
	}

	buf := new(bytes.Buffer)
	if err = tem.Execute(buf, data); err != nil {
		return nil, errorx.Wrap(err, "template execute error:", t.text)
	}

	if !t.goFmt {
		return buf, nil
	}

	formatOutput, err := goformat.Source(buf.Bytes())
	if err != nil {
		return nil, errorx.Wrap(err, "go format error:", buf.String())
	}

	buf.Reset()
	buf.Write(formatOutput)
	return buf, nil
}

// GetHead 生成文件head
func GetHead(source string) string {
	buffer, _ := With("head").Parse(headTemplate).Execute(map[string]interface{}{
		"source": source,
	})
	return buffer.String()
}
