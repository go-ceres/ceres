package golang

import (
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"go/build"
	"os"
	"path/filepath"
)

// GoBin returns a path of GOBIN.
func GoBin() string {
	def := build.Default
	goroot := os.Getenv("GOPATH")
	bin := filepath.Join(goroot, "bin")
	if !pathx.FileExists(bin) {
		gopath := os.Getenv("GOROOT")
		bin = filepath.Join(gopath, "bin")
	}
	if !pathx.FileExists(bin) {
		bin = os.Getenv("GOBIN")
	}
	if !pathx.FileExists(bin) {
		bin = filepath.Join(def.GOPATH, "bin")
	}
	return bin
}
