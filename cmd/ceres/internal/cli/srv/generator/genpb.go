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

package generator

import (
	"fmt"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/config"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/execx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (g *Generator) GenPb(ctx DirContext, c *config.Config) error {
	return g.genPbDirect(ctx, c)
}

func (g *Generator) genPbDirect(ctx DirContext, c *config.Config) error {
	g.log.Debug("[command]: %s", c.ProtocCmd)
	err := pathx.MkdirIfNotExist(c.ProtocOut)
	if err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = execx.Command(c.ProtocCmd, pwd)
	if err != nil {
		return err
	}
	return g.setPbDir(ctx, c)
}

func (g *Generator) setPbDir(ctx DirContext, c *config.Config) error {
	pbDir, err := findPbFile(c.ProtocOut, false)
	if err != nil {
		return err
	}
	if len(pbDir) == 0 {
		return fmt.Errorf("pg.go is not found under %q", c.ProtocOut)
	}
	grpcDir, err := findPbFile(c.ProtocOut, true)
	if err != nil {
		return err
	}
	if len(grpcDir) == 0 {
		return fmt.Errorf("_grpc.pb.go is not found in %q", c.ProtocOut)
	}
	if pbDir == c.Dist {
		return fmt.Errorf("the output of pb.go and _grpc.pb.go must not be the same "+
			"with --dist:\npb output: %s\ngrpc out: %s", pbDir, c.Dist)
	}
	ctx.SetProtoDir(pbDir)
	return nil
}

const (
	pbSuffix   = "pb.go"
	grpcSuffix = "_grpc.pb.go"
)

func findPbFile(current string, grpc bool) (string, error) {
	fileSystem := os.DirFS(current)
	var ret string
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, pbSuffix) {
			if grpc {
				if strings.HasSuffix(path, grpcSuffix) {
					ret = path
					return os.ErrExist
				}
			} else if !strings.HasSuffix(path, grpcSuffix) {
				ret = path
				return os.ErrExist
			}
		}
		return nil
	})
	if err == os.ErrExist {
		return filepath.Dir(filepath.Join(current, ret)), nil
	}
	return "", err
}
