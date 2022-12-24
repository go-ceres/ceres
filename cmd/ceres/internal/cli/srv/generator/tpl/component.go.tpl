package pkg

import (
    {{range .ImportPackage}}{{.}}
{{end}}
)

{{.ExtraFunc}}



func New{{.CamelName}}() {{.TypeName}} {
    return {{.InitStr}}
}
