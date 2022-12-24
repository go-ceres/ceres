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

import _ "embed"

//go:embed tpl/provide.go.tpl
var provideTemplate string

const (
	category                       = "srv"
	configTemplateFilename         = "config.toml.tpl"
	componentTemplateFilename      = "component.go.tpl"
	infrastructureTemplateFileName = "infrastructure.go.tpl"
	serviceTemplateFilename        = "service.go.tpl"
	serviceFuncTemplateFilename    = "service-func.go.tpl"
	actionTemplateFilename         = "action.go.tpl"
	actionFuncTemplateFilename     = "action-func.go.tpl"
	controllerTemplateFilename     = "controller.go.tpl"
	provideTemplateFilename        = "provide.go.tpl"
	grpcServerTemplateFilename     = "grpc.go.tpl"
	httpServerTemplateFilename     = "http.go.tpl"
	irepositoryTemplateFilename    = "irepository.go.tpl"
	entityTemplateFilename         = "entity.tpl"
	businessTemplateFilename       = "business.go.tpl"
	repositoryTemplateFilename     = "repository.go.tpl"
	wireTemplateFilename           = "wire.go.tpl"
	bootstrapTemplateFilename      = "bootstrap.go.tpl"
	mainTemplateFilename           = "main.go.tpl"
	makefileTemplateFilename       = "makefile.makefile.tpl"
	dockerfileTemplateFilename     = "dockerfile.dockerfile.tpl"
	//logicTemplateFileFile          = "logic.tpl"
	//serviceTemplateFileFile        = "service.tpl"
	//dtoTemplateFileFile            = "dto.tpl"
	//serviceFuncTemplateFileFile    = "service-func.tpl"
	//serverTemplateFile             = "server.tpl"
	//serverFuncTemplateFile         = "server-func.tpl"
	//globalTemplateFile             = "global.tpl"
	//logicFuncTemplateFileFile      = "logic-func.tpl"
	//bootTemplateFileFile           = "boot.tpl"
	//serverBootTemplateFileFile     = "server-boot.tpl"
	//mainTemplateFileFile           = "main.tpl"
)
