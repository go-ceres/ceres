//go:build wireinject
// +build wireinject

package main

import (
    "github.com/go-ceres/ceres/pkg/app"
    "github.com/google/wire"
{{if .ImportsStr}}{{.ImportsStr}}{{end}}
)

// injectionApp build ceres
func injectionApp() (*app.Application,func(),error) {
    panic(wire.Build(
        newApp,
		action.ProvideSet,
		domain.ProvideSet,
        infrastructure.ProvideSet,
        server.ProvideSet,
    ))
}
