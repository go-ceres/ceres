type (
	default{{.camelName}}Repository struct {
		db *gorm.DB{{if .cache}}
		cache cache.Cache{{end}}
	}

	{{.camelName}} struct{
		{{.fields}}
	}

	{{.camelName}}List []*{{.camelName}}

	QueryParam struct {
        gorm.PaginationParam
        gorm.QueryOptions
    }

	{{.camelName}}QueryResult struct {
        PageResult *gorm.PaginationResult
        List                  {{.camelName}}List
    }
)
