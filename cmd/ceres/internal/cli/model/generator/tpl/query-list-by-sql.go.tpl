func (m *default{{.camelName}}Repository) QueryListBySql(ctx context.Context, params *QueryParam, sql string,args ...interface{} ) (*{{.camelName}}QueryResult,error) {
	db := Get{{.camelName}}Db(ctx, m.db)
	if len(sql) > 0 {
        db.Where(sql, args)
    }
    opt := gorm.GetQueryOption(params.QueryOptions)
	opt.OrderFields = append(opt.OrderFields, gorm.NewOrderField("{{.primary}}", gorm.OrderByDESC))
	db = db.Order(gorm.ParseOrder(opt.OrderFields))
	if len(opt.OrderFields) > 0 {
        db.Select(opt.SelectFields)
    }
	var po {{.camelName}}List
	pr, err := gorm.WrapPageQuery(ctx, db, params.PaginationParam, &po)
	if err != nil {
        return nil, err
    }
    qr := &{{.camelName}}QueryResult{
        List:       po,
        PageResult: pr,
    }
    return qr, nil
}
