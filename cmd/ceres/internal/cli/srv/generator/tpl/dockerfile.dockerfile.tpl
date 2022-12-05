FROM alpine:3.10
WORKDIR /{{.Name}}
RUN mkdir "etc"
RUN mkdir "configs"
COPY ./configs/ /{{.Name}}/configs
COPY ./bin/ /{{.Name}}/bin
EXPOSE 5200
EXPOSE 5201
ENTRYPOINT [ "./bin/{{.Name}}","-f","./configs/config.yaml" ]
