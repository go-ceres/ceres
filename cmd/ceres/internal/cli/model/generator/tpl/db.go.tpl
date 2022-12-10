func Get{{.camelName}}Db(ctx context.Context,def *gorm.DB) *gorm.DB {
	return gorm.GetDbWithModel(ctx,def,new({{.camelName}}))
}
