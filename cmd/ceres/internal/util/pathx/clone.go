//    Copyright 2022. ceres
//    Author https://github.com/go-ceres/ceres
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package pathx

import (
	"fmt"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/env"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func CloneIntoGitHome(url, branch string) (dir string, err error) {
	gitHome, err := GetGitHome()
	if err != nil {
		return "", err
	}
	_ = os.RemoveAll(gitHome)
	ext := filepath.Ext(url)
	repo := strings.TrimSuffix(filepath.Base(url), ext)
	dir = filepath.Join(gitHome, repo)
	if FileExists(dir) {
		_ = os.RemoveAll(dir)
	}
	path, err := env.LookPath("git")
	if err != nil {
		return "", err
	}
	if !env.CanExec() {
		return "", fmt.Errorf("os %q can not call 'exec' command", runtime.GOOS)
	}
	args := []string{"clone"}
	if len(branch) > 0 {
		args = append(args, "-b", branch)
	}
	args = append(args, url, dir)
	cmd := exec.Command(path, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}
