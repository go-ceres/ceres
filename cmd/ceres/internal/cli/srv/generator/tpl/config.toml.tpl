[application]
	Name = "{{.serviceName}}"
[application.logger]
	mode = "std"
	level = "debug"
[application.transport.grpc.server]
	Address="0.0.0.0:5201"{{if .HttpServer}}
[application.transport.http.server]
	Address="0.0.0.0:5200"{{end}}
{{range .components}}{{.ConfigStr}}{{end}}
