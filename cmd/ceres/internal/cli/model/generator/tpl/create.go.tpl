func (m *default{{.camelName}}Repository) Create(ctx context.Context, param *{{.camelName}}) error {
	result := Get{{.camelName}}Db(ctx,m.db).Create(param)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
