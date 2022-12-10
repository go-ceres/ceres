func (m *default{{.camelName}}Repository) Update(ctx context.Context,{{.fieldName}} {{.fieldType}},param *{{.camelName}}) error {
    result:=Get{{.camelName}}Db(ctx,m.db).Where("{{.originalName}} = ?",{{.fieldName}}).Updates(param)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
