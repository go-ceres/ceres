package business

import (
    {{.imports}}
)

type {{.ServiceName}}Business struct {
    repo {{.IRepositoryPackageName}}.I{{.ServiceName}}Repository
}

func New{{.ServiceName}}Business(repo {{.IRepositoryPackageName}}.I{{.ServiceName}}Repository) *{{.ServiceName}}Business {
    return &{{.ServiceName}}Business{
		repo:repo,
    }
}
