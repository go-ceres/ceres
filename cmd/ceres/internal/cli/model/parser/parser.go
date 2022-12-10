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

package parser

import (
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/stringx"
	"io"
	"io/ioutil"
	"strings"
	"vitess.io/vitess/go/vt/sqlparser"
)

// Field 字段信息
type Field struct {
	Name         stringx.String // 字段名称
	OriginalName string         // 原始字段
	Type         string         // 字段类型（对应)
	Fulltext     bool           // 全文索引
	Spatial      bool           // 空间索引
	Unique       bool           // 唯一索引
	Primary      bool           // 是否时主键
	Number       bool           // 是否时数字类型
	Tag          string         // 字段对应的标签
}

// Statement 所有的表信息
type Statement struct {
	Tables []*Table // 表信息
}

// Table 表信息
type Table struct {
	Name     stringx.String // 表名
	DataBase stringx.String // 数据库名
	Primary  *Field         // 主键
	Options  []string       // 创建表结构的参数
	Time     bool           // 是否需要引入time组件
	Fields   []*Field       // 字段属性
}

// Index 索引信息
type Index sqlparser.IndexInfo

// Parse 解析sql原始文件
func Parse(filename, database string, strict bool) ([]*Table, error) {
	var tables []*Table
	sql, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	tokenizer := sqlparser.NewStringTokenizer(string(sql))
	for {
		next, err := sqlparser.ParseNext(tokenizer)
		if err == io.EOF {
			break
		}
		createTable, ok := next.(*sqlparser.CreateTable)
		if ok {
			tables = append(tables, ParseTable(createTable, database, strict))
		}
	}
	return tables, nil
}

// ParseTable 转换CreateTable
func ParseTable(table *sqlparser.CreateTable, database string, strict bool) *Table {
	//
	res := &Table{
		DataBase: stringx.NewString(database),
		Name:     stringx.NewString(table.Table.Name.String()),
		Options:  []string{},
		Fields:   []*Field{},
	}
	indexMap := make(map[string][]*Index)
	// 获取创建表结构的参数
	for _, option := range table.TableSpec.Options {
		res.Options = append(res.Options, option.Name+"="+option.String)
	}
	// 处理索引
	for _, index := range table.TableSpec.Indexes {
		for _, column := range index.Columns {
			is := func() []*Index {
				tempIs, ok := indexMap[column.Column.String()]
				if !ok {
					tempIs = make([]*Index, 0)
				}
				return tempIs
			}()
			is = append(is, &Index{
				Type:           index.Info.Type,
				Name:           index.Info.Name,
				ConstraintName: index.Info.ConstraintName,
				Primary:        index.Info.Primary,
				Spatial:        index.Info.Spatial,
				Fulltext:       index.Info.Fulltext,
				Unique:         index.Info.Unique,
			})
			indexMap[column.Column.String()] = is
		}
	}
	// 获取创建表时的字段
	for _, column := range table.TableSpec.Columns {
		field := new(Field)
		field.Name = stringx.NewString(column.Name.String())
		field.Number = false
		field.OriginalName = column.Name.String()
		field.Type = func(columnType string) string {
			switch columnType {
			case "bigint", "int":
				field.Number = true
				return "int64"
			case "varchar", "text", "char":
				return "string"
			case "tinyint":
				field.Number = true
				return "int"
			case "timestamp", "datetime", "date":
				res.Time = true
				return "time.Time"
			case "float":
				field.Number = true
				return "float32"
			case "double":
				field.Number = true
				return "float64"
			default:
				return "string"
			}
		}(strings.ToLower(column.Type.Type))
		// 标签
		tags := []string{"column:" + column.Name.String()}

		// 处理索引tag
		primary := false
		if indexs, ok := indexMap[column.Name.String()]; ok {
			for _, index := range indexs {
				if index.Primary {
					tags = append(tags, "primaryKey")
					primary = true
					field.Primary = true
				} else if index.Fulltext {
					if !index.Name.IsEmpty() && strings.ToLower(index.Name.String()) != "fulltext" {
						tags = append(tags, "index:"+index.Name.String()+",class:FULLTEXT")
					} else {
						tags = append(tags, "index:,class:FULLTEXT")
					}
					field.Fulltext = true
				} else if index.Unique {
					field.Unique = true
					if !index.Name.IsEmpty() && strings.ToLower(index.Name.String()) != "primary" {
						tags = append(tags, "uniqueIndex:"+index.Name.String())
					} else {
						tags = append(tags, "unique")
					}
				} else if index.Spatial {

				} else {
					if !index.Name.IsEmpty() && strings.ToLower(index.Name.String()) != "index" {
						tags = append(tags, "index:"+index.Name.String())
					} else {
						tags = append(tags, "index")
					}
				}
			}
		}
		// 处理类型
		switch strings.ToLower(column.Type.Type) {
		case "enum":
			tags = append(tags, "type:"+column.Type.Type+"("+strings.Join(column.Type.EnumValues, ",")+")")
		case "timestamp", "datetime", "date", "text", "longtext", "tinytext", "mediumtext":
			tags = append(tags, "type:"+column.Type.Type)
		default:
			// 如果是主键，不设置类型
			if !field.Primary {
				tags = append(tags, "type:"+column.Type.Type+"("+column.Type.Length.Val+")")
			}
		}
		// 自增列tag
		if column.Type.Options.Autoincrement {
			tags = append(tags, "autoIncrement")
		}
		// 是否为空
		if column.Type.Options.Null != nil && !*column.Type.Options.Null {
			tags = append(tags, "not null")
		}
		// 默认值
		defaultVal := getDefault(column)
		if column.Type.Options.Default != nil && len(defaultVal) > 0 {
			tags = append(tags, "default:"+defaultVal)
		}
		// 备注
		if column.Type.Options.Comment != nil {
			tags = append(tags, "comment:"+column.Type.Options.Comment.Val)
		}
		field.Tag = strings.Join(tags, ";")
		if primary {
			res.Primary = field
		}
		res.Fields = append(res.Fields, field)
	}

	return res
}

// getDefault 获取默认值
func getDefault(column *sqlparser.ColumnDefinition) string {
	defaultVal := ""
	if column.Type.Options.Default != nil {
		switch stmt := column.Type.Options.Default.(type) {
		case *sqlparser.Literal:
			switch stmt.Type {
			case sqlparser.StrVal:
				defaultVal = "'" + stmt.Val + "'"
			case sqlparser.IntVal, sqlparser.FloatVal, sqlparser.DOUBLE:
				defaultVal = stmt.Val
			}
		case *sqlparser.FuncExpr:
			if stmt.Name.String() == "current_timestamp" {
				defaultVal = "CURRENT_TIMESTAMP"
			}

		}

	}
	if column.Type.Options.OnUpdate != nil {
		switch onUpdate := column.Type.Options.OnUpdate.(type) {
		case *sqlparser.FuncExpr:
			if onUpdate.Name.String() == "current_timestamp" {
				defaultVal = "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"
			}
		}
	}

	return defaultVal
}
