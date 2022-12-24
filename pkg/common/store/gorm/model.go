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
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Model struct {
	ID        uint64 `gorm:"primaryKey;"`
	CreatedAt int64
	UpdatedAt int64
}

// GetDB 从上下文中获取DB
func GetDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	trans, ok := FromTrans(ctx)
	if ok && !FromNoTrans(ctx) {
		db, ok := trans.(*gorm.DB)
		if ok {
			if FromTransLock(ctx) {
				db = db.Clauses(clause.Locking{Strength: "UPDATE"})
			}
			return db
		}
	}

	return defDB
}

// GetDbWithModel 获取db
func GetDbWithModel(ctx context.Context, defDB *gorm.DB, m interface{}) *gorm.DB {
	return GetDB(ctx, defDB).Model(m)
}

// WrapPageQuery 包装带有分页的查询
func WrapPageQuery(ctx context.Context, db *gorm.DB, pp PaginationParam, out interface{}, maxPage ...int) (*PaginationResult, error) {
	if pp.OnlyCount {
		var count int64
		err := db.Count(&count).Error
		if err != nil {
			return nil, err
		}
		return &PaginationResult{Total: count}, nil
	} else if !pp.Pagination {
		err := db.Find(out).Error
		return nil, err
	}

	total, err := FindPage(ctx, db, pp, out, maxPage...)
	if err != nil {
		return nil, err
	}

	return &PaginationResult{
		Total: total,
		Page:  pp.GetPage(),
		Size:  pp.GetPageSize(),
	}, nil
}

// FindPage 查询分页数据
func FindPage(ctx context.Context, db *gorm.DB, pp PaginationParam, out interface{}, maxPage ...int) (int64, error) {
	var max int = 0
	if len(maxPage) > 0 {
		max = maxPage[0]
	}
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return 0, err
	} else if count == 0 {
		return count, nil
	}
	current, pageSize := pp.GetPage(), pp.GetPageSize()
	if current > int64(max) {
		current = int64(max)
	}
	if current > 0 && pageSize > 0 {
		db = db.Offset(int((current - 1) * pageSize)).Limit(int(pageSize))
	} else if pageSize > 0 {
		db = db.Limit(int(pageSize))
	}

	err = db.Find(out).Error
	return count, err
}

// FindOne 查询单条数据
func FindOne(db *gorm.DB, out interface{}) (bool, error) {
	result := db.First(out)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
