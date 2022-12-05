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

package formatx

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	flagGo    = "GO"
	flagCeres = "CERES"

	unknown style = iota
	title
	lower
	upper
)

// ErrNamingFormat defines an error for unknown format
var ErrNamingFormat = errors.New("unsupported format")

type (
	styleFormat struct {
		before     string
		through    string
		after      string
		goStyle    style
		ceresStyle style
	}

	style int
)

// FileNamingFormat is used to format the file name. You can define the format style
// through the go and ceres formatting characters. For example, you can define the snake
// format as go_ceres, and the camel case format as goCeres. You can even specify the split
// character, such as go#Ceres, theoretically any combination can be used, but the prerequisite
// must meet the naming conventions of each operating system file name.
// Note: Formatting is based on snake or camel string
func FileNamingFormat(format, content string) (string, error) {
	upperFormat := strings.ToUpper(format)
	indexGo := strings.Index(upperFormat, flagGo)
	indexCeres := strings.Index(upperFormat, flagCeres)
	if indexGo < 0 || indexCeres < 0 || indexGo > indexCeres {
		return "", ErrNamingFormat
	}
	var (
		before, through, after string
		flagGo, flagCeres      string
		goStyle, ceresStyle    style
		err                    error
	)
	before = format[:indexGo]
	flagGo = format[indexGo : indexGo+2]
	through = format[indexGo+2 : indexCeres]
	flagCeres = format[indexCeres : indexCeres+5]
	after = format[indexCeres+5:]

	goStyle, err = getStyle(flagGo)
	if err != nil {
		return "", err
	}

	ceresStyle, err = getStyle(flagCeres)
	if err != nil {
		return "", err
	}
	var formatStyle styleFormat
	formatStyle.goStyle = goStyle
	formatStyle.ceresStyle = ceresStyle
	formatStyle.before = before
	formatStyle.through = through
	formatStyle.after = after
	return doFormat(formatStyle, content)
}

func doFormat(f styleFormat, content string) (string, error) {
	splits, err := split(content)
	if err != nil {
		return "", err
	}
	var join []string
	for index, split := range splits {
		if index == 0 {
			join = append(join, transferTo(split, f.goStyle))
			continue
		}
		join = append(join, transferTo(split, f.ceresStyle))
	}
	joined := strings.Join(join, f.through)
	return f.before + joined + f.after, nil
}

func transferTo(in string, style style) string {
	switch style {
	case upper:
		return strings.ToUpper(in)
	case lower:
		return strings.ToLower(in)
	case title:
		return strings.Title(in)
	default:
		return in
	}
}

func split(content string) ([]string, error) {
	var (
		list   []string
		reader = strings.NewReader(content)
		buffer = bytes.NewBuffer(nil)
	)
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				if buffer.Len() > 0 {
					list = append(list, buffer.String())
				}
				return list, nil
			}
			return nil, err
		}
		if r == '_' {
			if buffer.Len() > 0 {
				list = append(list, buffer.String())
			}
			buffer.Reset()
			continue
		}

		if r >= 'A' && r <= 'Z' {
			if buffer.Len() > 0 {
				list = append(list, buffer.String())
			}
			buffer.Reset()
		}
		buffer.WriteRune(r)
	}
}

func getStyle(flag string) (style, error) {
	compare := strings.ToLower(flag)
	switch flag {
	case strings.ToLower(compare):
		return lower, nil
	case strings.ToUpper(compare):
		return upper, nil
	case strings.Title(compare):
		return title, nil
	default:
		return unknown, fmt.Errorf("unexpected format: %s", flag)
	}
}
