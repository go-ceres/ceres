package server

import (
    "github.com/go-ceres/ceres/server/grpc"
    {{.Imports}}
)

func NewGRPCServer({{.serverParamsStr}}) *grpc.Server {
    srv := grpc.ScanConfig().Build()
    {{.registerListStr}}
    return srv
}
