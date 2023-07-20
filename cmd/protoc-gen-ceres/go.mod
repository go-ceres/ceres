module github.com/go-ceres/ceres/cmd/protoc-gen-ceres

go 1.19

require (
	github.com/go-ceres/ceres v0.0.8
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	golang.org/x/sys v0.4.0 // indirect
)

//replace github.com/go-ceres/ceres => ../../ // 开发时，发布时注释
