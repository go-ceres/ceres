package server

import (
    "github.com/go-ceres/ceres/pkg/transport/grpc"
    {{.Imports}}
)

func NewGRPCServer({{.serverParamsStr}}) *grpc.Server {
    srv := grpc.ScanServerConfig().Build()
    {{.registerListStr}}
    return srv
}
