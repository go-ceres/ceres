[ceres.application]
	Name = "{{.serviceName}}"
[ceres.application.server.grpc]
	Address="0.0.0.0:5201"{{if .HttpServer}}
[ceres.application.server.{{.HttpServer}}]
	Address="0.0.0.0:5200"{{end}}
{{range .components}}{{.ConfigStr}}{{end}}
