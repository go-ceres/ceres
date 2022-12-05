//go:build wireinject
// +build wireinject

package main

import (
    "github.com/go-ceres/ceres"
    "github.com/google/wire"
{{if .ImportsStr}}{{.ImportsStr}}{{end}}
)

// injectionApp build ceres
func injectionApp() (*ceres.App,func(),error) {
    panic(wire.Build(
        newApp,
	    controller.ProvideSet,
		domain.ProvideSet,
        infrastructure.ProvideSet,
        server.ProvideSet,
    ))
}
