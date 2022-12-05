GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
NAME={{.Name}}
GOOS=linux
GOARCH=amd64
TEMP=${shell git diff-index --quiet HEAD; echo $$?}
STATUS=
ifeq ($(TEMP),0)
    STATUS=Clean
else
    STATUS=Modified
endif


.PHONY: install
# init env
install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-ceres/ceres@latest
	go install github.com/go-ceres/ceres/cmd/protoc-gen-ceres@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: build
# build
build:
	mkdir -p ../bin/ && GOOS=${GOOS} GOARCH=${GOARCH} go build \
    -ldflags "\
    -X github.com/go-ceres/ceres/version.appName='${NAME}' \
    -X github.com/go-ceres/ceres/version.appVersion='$(VERSION)' \
    -X github.com/go-ceres/ceres/version.buildTime='${shell date '+%Y-%m-%d--%T'}' \
    -X github.com/go-ceres/ceres/version.buildUser='${shell whoami}' \
    -X github.com/go-ceres/ceres/version.buildStatus='${STATUS}' \
    -X github.com/go-ceres/ceres/version.buildHost='${shell hostName -f}'" \
    -o ../bin/{{.Name}} ../bootstrap/

.PHONY: docker-build
docker-build:
	docker build ../ --platform ${GOOS}/${GOARCH} -t ${NAME}:${VERSION} -f Dockerfile

.PHONY: all
# generate all
all:
	make build;
