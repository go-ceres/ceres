func (m *default{{.camelName}}Repository) FindOne(ctx context.Context,params {{.camelName}}) (*{{.camelName}},error) {
    var en = new({{.camelName}})
	db := Get{{.camelName}}Db(ctx,m.db).Where(&params)
	_,err := gorm.FindOne(db,en)
	if err != nil {
		return nil,err
	}
	return en,nil
}
