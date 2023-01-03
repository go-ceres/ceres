module github.com/go-ceres/ceres/cmd/protoc-gen-ceres-error

go 1.18

require (
	github.com/go-ceres/ceres v0.0.4
	golang.org/x/text v0.5.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	golang.org/x/sys v0.2.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.51.0 // indirect
)

//replace github.com/go-ceres/ceres => ../../  // 开发时，发布时注释
