func AutoMigrateGorm{{.camelName}}(db *gorm.DB) error {
	return db{{if .options}}.Set("gorm:table_options","{{.options}}"){{end}}.AutoMigrate(new({{.camelName}}))
}
