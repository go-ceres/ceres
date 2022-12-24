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

package tags

import (
	"fmt"
	"github.com/emicklei/proto"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser/model"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/stringx"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"regexp"
)

var (
	tagType = regexp.MustCompile(`^\(api\.([^)]*)\)`)
	rInject = regexp.MustCompile("`.+`")
	rTags   = regexp.MustCompile(`[\w_]+:"[^"]+"`)
)

func ParseFile(inputPath string, src interface{}) (res map[string]*Message, err error) {
	res = make(map[string]*Message)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, inputPath, src, parser.ParseComments)
	if err != nil {
		return
	}
	for _, decl := range f.Decls {
		// check if is generic declaration
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		var typeSpec *ast.TypeSpec
		for _, spec := range genDecl.Specs {
			if ts, tsOK := spec.(*ast.TypeSpec); tsOK {
				typeSpec = ts
				break
			}
		}

		// skip if can't get type spec
		if typeSpec == nil {
			continue
		}

		// not a struct, skip
		structDecl, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}
		message := new(Message)
		message.Name = typeSpec.Name.Name
		message.Fields = make(map[string]*Field)
		for _, field := range structDecl.Fields.List {
			if len(field.Names) > 0 && field.Tag != nil {
				thisField := new(Field)
				thisField.Name = field.Names[0].Name
				thisField.CurrentTag = field.Tag.Value
				thisField.Start = int64(field.Pos())
				thisField.End = int64(field.End())
				message.Fields[thisField.Name] = thisField
			}

		}
		res[message.Name] = message
	}
	return
}

func WireFile(inputPath string, goMessages map[string]*Message, pb model.Proto) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	contents, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}

	// 查找需要替换的,然后替换掉
	for _, message := range pb.Message {
		for _, element := range message.Elements {
			normalField, ok := element.(*proto.NormalField)
			fieldName := stringx.NewString(normalField.Name).Title() // 字段名称
			// 必须存在在结构体里的数据
			if ok && goMessages[message.Name] != nil && goMessages[message.Name].Fields[fieldName] != nil {
				// 获取当前字段的信息
				repMsg := goMessages[message.Name].Fields[fieldName]
				// 修改json字段
				repMsg.InjectTag += fmt.Sprintf(`json:"%s"`, normalField.Name)
				for _, option := range normalField.Options {
					t := tagType.FindStringSubmatch(option.Name)
					if len(t) > 1 {
						repMsg.InjectTag += fmt.Sprintf(` %s:"%s"`, t[1], option.Constant.Source)
					}
				}
				contents = injectTag(contents, repMsg)
			}
			goMessages, err = ParseFile(inputPath, contents)
			if err != nil {
				return err
			}
		}
	}

	// 输出文件
	if err = os.WriteFile(inputPath, contents, 0o644); err != nil {
		return err
	}
	return nil
}
