{{.head}}

package {{.package}}

import (
{{if .notStream}} "context" {{end}}

{{.imports}}
)

type {{.service}}Service struct {
    {{.unimplementedServer}}
    {{ range .descList }}{{.UnTitleName}}Action *{{.ActionPackage}}.{{.ActionName}}
    {{ end }}
}

func New{{.service}}Service(
    {{ range .descList }}{{.UnTitleName}}Action *{{.ActionPackage}}.{{.ActionName}},
    {{ end }}
	) *{{.service}}Service {
    return &{{.service}}Service{
        {{ range .descList }}{{.UnTitleName}}Action:{{.UnTitleName}}Action,
        {{ end }}
    }
}

{{.funcs}}
