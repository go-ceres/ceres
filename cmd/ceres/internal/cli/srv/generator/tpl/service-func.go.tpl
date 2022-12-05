{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.service}}Service) {{.method}} ({{if .notStream}}ctx context.Context,{{if .hasReq}} req {{.request}}{{end}}{{else}}{{if .hasReq}} req {{.request}},{{end}}stream {{.streamBody}}{{end}}) ({{if .notStream}}{{.response}},{{end}}error) {
	return {{if .notStream}}&{{.responseType}}{},{{end}}nil
}
