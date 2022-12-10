package {{.pkg}}

type {{.upperStartCamelObject}}Repository struct {
	*default{{.upperStartCamelObject}}Model
}

func New{{.upperStartCamelObject}}Repository() *{{.upperStartCamelObject}}Repository {
	return &{{.upperStartCamelObject}}Repository{
		default{{.upperStartCamelObject}}Repository: new{{.upperStartCamelObject}}Repository(),
	}
}
