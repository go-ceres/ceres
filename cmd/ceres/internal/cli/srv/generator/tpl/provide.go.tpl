package {{.PackageName}}

import (
    "github.com/google/wire"
{{if .ImportsStr}}{{.ImportsStr}}{{end}}
)


var ProvideSet = wire.NewSet({{.ProvideSetStr}})
