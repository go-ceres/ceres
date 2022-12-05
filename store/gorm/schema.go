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

package gorm

import (
	"fmt"
	"strings"
)

// PaginationParam 分页查询参数
type PaginationParam struct {
	Pagination bool `form:"-"`                                     // 是否使用分页查询
	OnlyCount  bool `form:"-"`                                     // 是否仅查询count
	Page       int  `form:"page,default=1"`                        // 当前页
	PageSize   int  `form:"pageSize,default=10" binding:"max=100"` // 页大小
}

// PaginationResult 分页查询结果
type PaginationResult struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// QueryOptions 查询可选参数项
type QueryOptions struct {
	OrderFields  []*OrderField // 排序字段
	SelectFields []string      // 查询字段
}

// OrderDirection 排序方向
type OrderDirection int

// OrderFieldFunc 排序字段转换函数
type OrderFieldFunc func(string) string

const (
	// OrderByASC 升序排序
	OrderByASC OrderDirection = iota + 1
	// OrderByDESC 降序排序
	OrderByDESC
)

// NewOrderField 创建排序字段
func NewOrderField(key string, d OrderDirection) *OrderField {
	return &OrderField{
		Key:       key,
		Direction: d,
	}
}

// OrderField 排序字段
type OrderField struct {
	Key       string         // 字段名(字段名约束为小写蛇形)
	Direction OrderDirection // 排序方向
}

// GetPage 获取当前页
func (p PaginationParam) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 每页现实数量
func (p PaginationParam) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	return p.PageSize
}

// GetQueryOption 获取查询参数
func GetQueryOption(opts ...QueryOptions) QueryOptions {
	var opt QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// ParseOrder 解析排序字段
func ParseOrder(items []*OrderField, handle ...OrderFieldFunc) string {
	orders := make([]string, len(items))

	for i, item := range items {
		key := item.Key
		if len(handle) > 0 {
			key = handle[0](key)
		}

		direction := "ASC"
		if item.Direction == OrderByDESC {
			direction = "DESC"
		}
		orders[i] = fmt.Sprintf("%s %s", key, direction)
	}

	return strings.Join(orders, ",")
}
