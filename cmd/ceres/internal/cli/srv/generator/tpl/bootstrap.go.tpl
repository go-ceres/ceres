package main

import (
    "github.com/go-ceres/ceres"
    "github.com/go-ceres/ceres/pkg/app"
    "github.com/go-ceres/ceres/pkg/common/config"
    "github.com/go-ceres/ceres/pkg/common/config/file"
    "github.com/go-ceres/ceres/pkg/common/flag"
    "github.com/go-ceres/ceres/pkg/common/logger"
    "github.com/go-ceres/ceres/pkg/transport"
    "github.com/go-ceres/ceres/pkg/transport/grpc"{{if .HttpServer}}
    "github.com/go-ceres/ceres/pkg/transport/http"{{end}}



{{.PackageImports}}
)

var (
    confPath = flag.String("conf", "../configs/config.toml", "config path, eg: -conf config.yaml","f")
)

func newApp(gs *grpc.Server{{if .HttpServer}},hs *http.Server{{end}}{{if .hasRegistry}},registry transport.Registry{{end}}) *app.Application{
    return app.ScanConfig().WithOption({{if .hasRegistry}}
			app.WithRegistry(registry),{{end}}
			app.WithTransport(gs{{if .HttpServer}},hs{{end}}),
        ).Build()
}

func main()  {
	// parse flag
    flag.Parse()

    // init config
    if err := config.Load(file.NewSource(*confPath)); err!=nil {
		panic(err)
    }

	// init logger
	logger.SetLogger(
		logger.ScanConfig().WithOptions(
			logger.WithFields(
                logger.FieldAid(ceres.AppId()),
				logger.FieldName(ceres.AppName()),
				logger.FieldVersion(ceres.AppVersion()),
				logger.FieldHostName(ceres.HostName()),
            ),
        ).Build(),
    )

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
