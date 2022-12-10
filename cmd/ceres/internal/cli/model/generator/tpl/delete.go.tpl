func (m *default{{.camelName}}Repository) Delete(ctx context.Context,{{.fieldName}}s []{{.fieldType}}) error {
	result:=Get{{.camelName}}Db(ctx,m.db).Where("{{.originalName}} IN (?)",{{.fieldName}}s).Delete({{.camelName}}{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
