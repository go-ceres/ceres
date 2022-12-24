package {{.package}}

import (
{{if .notStream}} "context" {{end}}

{{.imports}}
)

type {{.actionName}} struct {
    bus *{{.businessType}}
}

func New{{.actionName}}(bus *{{.businessType}}) *{{.actionName}} {
    return &{{.actionName}}{
        bus:bus,
    }
}

{{.funcs}}
