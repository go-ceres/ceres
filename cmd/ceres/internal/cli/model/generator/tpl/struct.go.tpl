type (
	default{{.camelName}}Repository struct {
		db *gorm.DB{{if .cache}}
		cache cache.Cache{{end}}
	}

	{{.camelName}} struct{
		{{.fields}}
	}

	{{.camelName}}List []*{{.camelName}}

	{{.camelName}}ListQueryParam struct {
        gorm.PaginationParam
        gorm.QueryOptions
    }

	{{.camelName}}ListQueryResult struct {
        PageResult *gorm.PaginationResult
        List                  {{.camelName}}List
    }
)
