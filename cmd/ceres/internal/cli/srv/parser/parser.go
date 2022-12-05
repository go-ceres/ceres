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

package parser

import (
	"github.com/emicklei/proto"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser/model"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

type (
	// DefaultProtoParser types an empty struct
	DefaultProtoParser struct{}
)

// NewDefaultProtoParser 创建一个新的解析结构体
func NewDefaultProtoParser() *DefaultProtoParser {
	return &DefaultProtoParser{}
}

// Parse 将proto原始文件解析为结构体
func (p *DefaultProtoParser) Parse(src string, multiple ...bool) (model.Proto, error) {
	var ret model.Proto

	abs, err := filepath.Abs(src)
	if err != nil {
		return model.Proto{}, err
	}

	r, err := os.Open(abs)
	if err != nil {
		return ret, err
	}
	defer func(r *os.File) {
		_ = r.Close()
	}(r)

	parser := proto.NewParser(r)
	set, err := parser.Parse()
	if err != nil {
		return ret, err
	}

	var serviceList model.Services
	proto.Walk(
		set,
		proto.WithImport(func(i *proto.Import) {
			ret.Import = append(ret.Import, model.Import{Import: i})
		}),
		proto.WithMessage(func(message *proto.Message) {
			ret.Message = append(ret.Message, model.Message{Message: message})
		}),
		proto.WithPackage(func(p *proto.Package) {
			ret.Package = model.Package{Package: p}
		}),
		proto.WithService(func(service *proto.Service) {
			serv := model.Service{Service: service}
			elements := service.Elements
			for _, el := range elements {
				v, _ := el.(*proto.RPC)
				if v == nil {
					continue
				}
				serv.RPC = append(serv.RPC, &model.RPC{RPC: v})
			}

			serviceList = append(serviceList, serv)
		}),
		proto.WithOption(func(option *proto.Option) {
			if option.Name == "go_package" {
				ret.GoPackage = option.Constant.Source
			}
		}),
	)
	if err = serviceList.Validate(abs, multiple...); err != nil {
		return ret, err
	}

	if len(ret.GoPackage) == 0 {
		ret.GoPackage = ret.Package.Name
	}

	ret.PbPackage = GoSanitized(filepath.Base(ret.GoPackage))
	ret.Src = abs
	ret.Name = filepath.Base(abs)
	ret.Service = serviceList

	return ret, nil
}

// GoSanitized 对字符串进行格式化处理
func GoSanitized(s string) string {
	// Sanitize the input to the set of valid characters,
	// which must be '_' or be in the Unicode L or N categories.
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, s)

	// Prepend '_' in the event of a Go keyword conflict or if
	// the identifier is invalid (does not start in the Unicode L category).
	r, _ := utf8.DecodeRuneInString(s)
	if token.Lookup(s).IsKeyword() || !unicode.IsLetter(r) {
		return "_" + s
	}
	return s
}

// CamelCase 驼峰处理
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
