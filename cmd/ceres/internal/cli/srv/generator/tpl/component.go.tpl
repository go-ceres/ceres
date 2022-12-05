package pkg

import (
    {{range .ImportPackage}}{{.}}
{{end}}
)

func New{{.CamelName}}() {{.TypeName}} {
    return {{.InitStr}}
}
