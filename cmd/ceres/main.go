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

package main

import (
	"fmt"
	"github.com/go-ceres/ceres"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/model"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv"
	"github.com/go-ceres/cli/v2"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

var rootCommand = []*cli.Command{
	{
		Name:        "srv",
		Subcommands: srv.Commands,
		Usage:       "generate a service for ceres",
	},
	{
		Name:        "model",
		Subcommands: model.Commands,
		Usage:       "generate model for ceres",
	},
	{
		Name: "build",
	},
	{
		Name: "upgrade",
	},
	{
		Name: "run",
	},
}

func main() {
	app := cli.NewApp()
	app.Usage = "a cli tool for ceres"
	app.Description = "a cli tool for ceres"
	app.UseShortOptionHandling = true
	app.Version = fmt.Sprintf("%s %s/%s", ceres.Version, runtime.GOOS, runtime.GOARCH)
	app.Commands = rootCommand
	app.Flags = append(app.Flags)
	app.ExitErrHandler = func(context *cli.Context, err error) {
		_ = cli.ShowCommandHelp(context, context.Command.Name)
		logrus.Info(err)
	}
	_ = app.Run(os.Args)
}
