package {{.package}}

import (
{{if .notStream}} "context" {{end}}

{{.imports}}
)

type {{.service}}Controller struct {
    {{.unimplementedServer}}
    bus *{{.businessType}}
}

func New{{.service}}Controller(bus *{{.businessType}}) *{{.service}}Controller {
    return &{{.service}}Controller{
		bus:bus,
    }
}

{{.funcs}}
