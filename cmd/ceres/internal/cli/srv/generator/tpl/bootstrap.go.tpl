package main

import (
    "github.com/go-ceres/ceres"
    "github.com/go-ceres/ceres/config"
    "github.com/go-ceres/ceres/config/file"
    "github.com/go-ceres/ceres/flag"
    "github.com/go-ceres/ceres/logger"
    "github.com/go-ceres/ceres/server/grpc"


{{.PackageImports}}
)

var (
    confPath = flag.String("conf", "../configs/config.toml", "config path, eg: -conf config.yaml","f")
)

func newApp(gs *grpc.Server{{if .HttpServer}},hs *{{.HttpServer}}.Server{{end}}{{if .hasRegistry}},registry registry.Registry{{end}}) *ceres.App{
    return ceres.ScanConfig().{{if .hasRegistry}}
        SetRegistry(registry).{{end}}
        AddServers(gs{{if .HttpServer}},hs{{end}}).
        Build()
}

func main()  {
	// parse flag
    flag.Parse()
	// init logger
    logger.SetLogger(logger.NewZapLogger(logger.DefaultZapConfig()))
    // load configuration Source
    if err := config.Load(file.NewSource(*confPath)); err!=nil {
		panic(err)
    }
	// wire dependency injection
    app , clear , err := injectionApp()
    if err != nil {
        panic(err)
    }
	defer clear()
	// start
	err = app.Run()
    if err!=nil {
        panic(err)
    }
}
