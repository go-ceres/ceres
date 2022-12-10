func new{{.camelName}}Repository(db *gorm.DB{{if .cache}}cache cache.Cache{{end}}) *default{{.camelName}}Repository {
	return &default{{.camelName}}Repository{{.extra.LeftBrackets}}
		db: db,{{if .cache}}
		cache: cache,{{end}}
	}
}
