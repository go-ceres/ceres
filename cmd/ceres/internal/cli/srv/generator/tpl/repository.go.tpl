package repository

import (
    "github.com/go-ceres/ceres/pkg/common/store/gorm"
	{{.Imports}}
)

var _ {{.IRepositoryPackageName}}.I{{.ServiceName}}Repository = (*{{.ServiceName}}Repository)(nil)

type {{.ServiceName}}Repository struct {
    db *gorm.DB
}

func New{{.ServiceName}}Repository(db *gorm.DB) {{.IRepositoryPackageName}}.I{{.ServiceName}}Repository {
    return &{{.ServiceName}}Repository{
	    db:db,
    }
}
