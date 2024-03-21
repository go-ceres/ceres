FROM golang:1.22 AS build
WORKDIR /project
COPY . /project
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go env -w GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/{{.Name}} ./bootstrap


FROM alpine
ENV TZ Asia/Shanghai
COPY --from=build /project/bin /{{.Name}}/bin
COPY --from=build /project/configs /{{.Name}}/configs
WORKDIR /examine
EXPOSE 5200
EXPOSE 5201
ENTRYPOINT [ "./bin/{{.Name}}","-f","./configs/config.yaml" ]
