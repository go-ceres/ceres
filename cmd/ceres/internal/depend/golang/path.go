package golang

import (
	"github.com/go-ceres/ceres/cmd/ceres/internal/ctx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"path/filepath"
	"strings"
)

func GetParentPackage(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	projectConf, err := ctx.PrepareProject(abs)
	if err != nil {
		return "", err
	}

	// fix https://github.com/zeromicro/go-zero/issues/1058
	wd := projectConf.WorkDir
	d := projectConf.Dir
	same, err := pathx.SameFile(wd, d)
	if err != nil {
		return "", err
	}

	trim := strings.TrimPrefix(projectConf.WorkDir, projectConf.Dir)
	if same {
		trim = strings.TrimPrefix(strings.ToLower(projectConf.WorkDir), strings.ToLower(projectConf.Dir))
	}

	return filepath.ToSlash(filepath.Join(projectConf.Path, trim)), nil
}
