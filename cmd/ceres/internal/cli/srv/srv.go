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

package srv

import (
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/action"
	"github.com/go-ceres/ceres/cmd/ceres/internal/common/flag"
	"github.com/go-ceres/cli/v2"
)

var (
	Flags    []cli.Flag
	Commands = []*cli.Command{
		{
			Name:        "protoc",
			Flags:       append(action.ProtocFlags, flag.CommonFlags...),
			Description: "Example：ceres srv protoc -f xx.proto --go_out=./ --go-grpc_out=./ ./",
			Action:      action.ProtocAction,
		},
	}
)
