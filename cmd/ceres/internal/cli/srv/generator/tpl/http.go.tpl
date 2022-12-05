package server

import (
	{{.Imports}}
)

func NewHTTPServer({{.serverParamsStr}}) *{{.ServerType}}.Server {
    srv := {{.ServerType}}.ScanConfig().Build()
    {{.registerListStr}}
    return srv
}
