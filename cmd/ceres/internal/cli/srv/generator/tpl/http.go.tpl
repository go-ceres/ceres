package server

import (
	{{.Imports}}
    "github.com/go-ceres/ceres/pkg/transport/http"
)

func NewHTTPServer({{.serverParamsStr}}) *http.Server {
    srv := http.ScanServerConfig().Build()
    {{.registerListStr}}
    return srv
}
