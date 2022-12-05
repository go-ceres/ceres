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
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/config"
	"github.com/go-ceres/ceres/cmd/ceres/internal/cli/srv/parser/model"
	"github.com/go-ceres/ceres/cmd/ceres/internal/ctx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/pathx"
	"github.com/go-ceres/ceres/cmd/ceres/internal/util/stringx"
	"github.com/go-ceres/cli/v2"
	"path/filepath"
	"strings"
)

var _ DirContext = (*defaultDirContext)(nil)

type dirKey string

const (
	wdKey             dirKey = "wd" // 工作
	binKey            dirKey = "bin"
	interfaceKey      dirKey = "interface" // 接口层
	aclKey            dirKey = "acl"
	bootstrapKey      dirKey = "bootstrap"      // 启动配置项
	configsKey        dirKey = "configs"        // 配置文件
	deploymentKey     dirKey = "deployment"     // 部属相关
	protoKey          dirKey = "proto"          // pb文件输出 (接口层下的子目录)
	internalKey       dirKey = "internal"       // 内部代码
	infrastructureKey dirKey = "infrastructure" // 底层代码 (属于内部代码下) 如常量，需要用到的底层包
	pkgKey            dirKey = "pkg"            // 组件依赖安装
	repositoryKey     dirKey = "repository"     // 存储仓
	serverKey         dirKey = "server"         // 服务注册层（接口层下的子目录）
	domainKey         dirKey = "domain"         // 领域层代码，主要实现逻辑 (属于内部代码下)
	businessKey       dirKey = "business"       // 业务处理逻辑层
	entityKey         dirKey = "entity"         // 领域实体
	irepositoryKey    dirKey = "irepository"    // 领域存储接口
	controllerKey     dirKey = "controller"     // 服务接入层 // 很潜入的一层，用于
)

type (
	// DirContext 文件夹上下文接口
	DirContext interface {
		GetMain() Dir      // 项目启动文件
		GetBootstrap() Dir // 获取启动配置目录 ------配置
		GetAcl() Dir
		GetDeployment() Dir     //获取部属相关
		GetConfigs() Dir        // 获取配置目录 -------配置文件
		GetInterface() Dir      // 接口层	-------对外提供服务
		GetProto() Dir          // proto输出目录 --------protoc输出目录
		GetInternalKey() Dir    // 内部代码层-----防止包调错
		GetServer() Dir         // 应用层服务启动层-------注册服务
		GetDomain() Dir         // 业务逻辑层-----------对应领域层
		GetInfrastructure() Dir // 基础设施层 ------基础设置，如数据，调用的额外包，常量
		GetPkg() Dir            // 包依赖层
		GetRepository() Dir     // 仓储
		GetController() Dir     // 服务层
		GetBusiness() Dir       // 业务逻辑层
		GetIRepository() Dir    // 领域存储层接口
		GetEntity() Dir         // 领域对象
		GetServiceName() stringx.String
		SetProtoDir(proto string)
	}
	// Dir 文件路径
	Dir struct {
		Package         string                                 // 文件完整包名
		Base            string                                 // 文件最后一级包名
		Filename        string                                 // 文件完整路径
		GetChildPackage func(childPath string) (string, error) // 文件夹完整路径
	}
	// defaultDirContext 文件夹管理上下文
	defaultDirContext struct {
		dirMap      map[dirKey]Dir // 文件夹集合
		serviceName stringx.String // 服务名，该类型方便对字符串进行处理
		project     *ctx.Project   // 项目上下文
		cliCtx      *cli.Context
	}
)

func (d defaultDirContext) GetAcl() Dir {
	return d.dirMap[aclKey]
}

func (d defaultDirContext) GetDeployment() Dir {
	return d.dirMap[deploymentKey]
}

func (d defaultDirContext) GetBusiness() Dir {
	return d.dirMap[businessKey]
}

func (d defaultDirContext) GetIRepository() Dir {
	return d.dirMap[irepositoryKey]
}

func (d defaultDirContext) GetEntity() Dir {
	return d.dirMap[entityKey]
}

func (d defaultDirContext) GetController() Dir {
	return d.dirMap[controllerKey]
}

func (d defaultDirContext) GetPkg() Dir {
	return d.dirMap[pkgKey]
}

func (d defaultDirContext) GetRepository() Dir {
	return d.dirMap[repositoryKey]
}

func (d defaultDirContext) GetMain() Dir {
	return d.dirMap[wdKey]
}

func (d defaultDirContext) GetBootstrap() Dir {
	return d.dirMap[bootstrapKey]
}

func (d defaultDirContext) GetConfigs() Dir {
	return d.dirMap[configsKey]
}

func (d defaultDirContext) GetInterface() Dir {
	return d.dirMap[interfaceKey]
}

func (d defaultDirContext) GetProto() Dir {
	return d.dirMap[protoKey]
}

func (d defaultDirContext) GetInternalKey() Dir {
	return d.dirMap[internalKey]
}

func (d defaultDirContext) GetServer() Dir {
	return d.dirMap[serverKey]
}

func (d defaultDirContext) GetDomain() Dir {
	return d.dirMap[domainKey]
}

func (d defaultDirContext) GetInfrastructure() Dir {
	return d.dirMap[infrastructureKey]
}

func (d defaultDirContext) GetServiceName() stringx.String {
	return d.serviceName
}

func (d defaultDirContext) SetProtoDir(protoPath string) {
	d.dirMap[protoKey] = Dir{
		Filename: protoPath,
		Package:  filepath.ToSlash(filepath.Join(d.project.Path, strings.TrimPrefix(protoPath, d.project.Dir))),
		Base:     filepath.Base(protoPath),
	}
}

func (d *Dir) Valid() bool {
	return len(d.Filename) > 0 && len(d.Package) > 0
}

// mkdir 创建文件夹
func (g *Generator) mkdir(project *ctx.Project, proto model.Proto, conf *config.Config) (DirContext, error) {
	dirMap := make(map[dirKey]Dir)
	binDir := filepath.Join(project.WorkDir, "bin")                   // 执行文件生成目录
	interfaceDir := filepath.Join(project.WorkDir, "interface")       // 接口层
	aclDir := filepath.Join(interfaceDir, "acl")                      // 接口防腐层
	bootstrapDir := filepath.Join(project.WorkDir, "bootstrap")       // 服务启动目录
	configsDir := filepath.Join(project.WorkDir, "configs")           // 配置文件目录
	deploymentDir := filepath.Join(project.WorkDir, "deployment")     // 部属相关目录
	protoDir := filepath.Join(interfaceDir, "proto")                  // proto生成文件
	internalDir := filepath.Join(project.WorkDir, "internal")         // 内部包隔离
	serverDir := filepath.Join(internalDir, "server")                 // 服务注册
	controllerDir := filepath.Join(internalDir, "controller")         // 控制器，进行DTO-BO对象转换
	domainDir := filepath.Join(internalDir, "domain")                 // 业务层，主要处理业务罗，实现BO与DO数据转换
	businessDir := filepath.Join(domainDir, "business")               //业务逻辑处理
	entityDir := filepath.Join(domainDir, "entity")                   //领域实体
	irepositoryDir := filepath.Join(domainDir, "irepository")         // 领域存储接口
	infrastructureDir := filepath.Join(internalDir, "infrastructure") // 基础设置层
	pkgDir := filepath.Join(infrastructureDir, "pkg")                 // 依赖
	repositoryDir := filepath.Join(infrastructureDir, "repository")   // 仓储
	getChildPackage := func(parent, childPath string) (string, error) {
		child := strings.TrimPrefix(childPath, parent)
		abs := filepath.Join(parent, strings.ToLower(child))
		childPath = strings.TrimPrefix(abs, project.Dir)
		pkg := filepath.Join(project.Path, childPath)
		return filepath.ToSlash(pkg), nil
	}
	conf.ProtocOut = protoDir
	// 组装protoc参数
	protoArgs := wrapProtocCmd(conf, g.cliCtx.Args().Slice())
	// 设置protoc命令
	conf.ProtocCmd = strings.Join(protoArgs, " ")
	dirMap[businessKey] = Dir{
		Filename: businessDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(businessDir, project.Dir))),
		Base:     filepath.Base(businessDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(businessDir, childPath)
		},
	}
	dirMap[aclKey] = Dir{
		Filename: aclDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(aclDir, project.Dir))),
		Base:     filepath.Base(aclDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(aclDir, childPath)
		},
	}
	dirMap[entityKey] = Dir{
		Filename: entityDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(entityDir, project.Dir))),
		Base:     filepath.Base(entityDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(entityDir, childPath)
		},
	}
	dirMap[irepositoryKey] = Dir{
		Filename: irepositoryDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(irepositoryDir, project.Dir))),
		Base:     filepath.Base(irepositoryDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(irepositoryDir, childPath)
		},
	}
	dirMap[wdKey] = Dir{
		Filename: project.WorkDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(project.WorkDir, project.Dir))),
		Base:     filepath.Base(project.WorkDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(project.WorkDir, childPath)
		},
	}
	dirMap[interfaceKey] = Dir{
		Filename: interfaceDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(interfaceDir, project.Dir))),
		Base:     filepath.Base(interfaceDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(interfaceDir, childPath)
		},
	}
	dirMap[bootstrapKey] = Dir{
		Filename: bootstrapDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(bootstrapDir, project.Dir))),
		Base:     filepath.Base(bootstrapDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(bootstrapDir, childPath)
		},
	}
	dirMap[configsKey] = Dir{
		Filename: configsDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(configsDir, project.Dir))),
		Base:     filepath.Base(configsDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(configsDir, childPath)
		},
	}
	dirMap[deploymentKey] = Dir{
		Filename: deploymentDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(deploymentDir, project.Dir))),
		Base:     filepath.Base(deploymentDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(deploymentDir, childPath)
		},
	}
	dirMap[protoKey] = Dir{
		Filename: protoDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(protoDir, project.Dir))),
		Base:     filepath.Base(protoDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(protoDir, childPath)
		},
	}
	dirMap[internalKey] = Dir{
		Filename: internalDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(internalDir, project.Dir))),
		Base:     filepath.Base(internalDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(internalDir, childPath)
		},
	}
	dirMap[serverKey] = Dir{
		Filename: serverDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(serverDir, project.Dir))),
		Base:     filepath.Base(serverDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(serverDir, childPath)
		},
	}
	dirMap[controllerKey] = Dir{
		Filename: controllerDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(controllerDir, project.Dir))),
		Base:     filepath.Base(controllerDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(controllerDir, childPath)
		},
	}
	dirMap[domainKey] = Dir{
		Filename: domainDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(domainDir, project.Dir))),
		Base:     filepath.Base(domainDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(domainDir, childPath)
		},
	}
	dirMap[infrastructureKey] = Dir{
		Filename: infrastructureDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(infrastructureDir, project.Dir))),
		Base:     filepath.Base(infrastructureDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(infrastructureDir, childPath)
		},
	}
	dirMap[pkgKey] = Dir{
		Filename: pkgDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(pkgDir, project.Dir))),
		Base:     filepath.Base(pkgDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(pkgDir, childPath)
		},
	}
	dirMap[repositoryKey] = Dir{
		Filename: repositoryDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(repositoryDir, project.Dir))),
		Base:     filepath.Base(repositoryDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(repositoryDir, childPath)
		},
	}
	dirMap[binKey] = Dir{
		Filename: binDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(binDir, project.Dir))),
		Base:     filepath.Base(binDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(binDir, childPath)
		},
	}

	for _, dir := range dirMap {
		if err := pathx.MkdirIfNotExist(dir.Filename); err != nil {
			return nil, err
		}
	}

	serviceName := strings.TrimSuffix(proto.Name, filepath.Ext(proto.Name))
	return &defaultDirContext{
		dirMap:      dirMap,
		project:     project,
		serviceName: stringx.NewString(strings.ReplaceAll(serviceName, "-", "")),
	}, nil
}

// wrapProtocCmd 包装protoc命令
func wrapProtocCmd(conf *config.Config, args []string) []string {
	res := append([]string{"protoc"}, args...)
	// 支持多个文件路径
	if len(conf.ProtoPath) > 0 {
		for _, s := range conf.ProtoPath {
			res = append(res, "--proto_path", s)
		}
	}
	// go插件参数
	for _, goOpt := range conf.GoOpt {
		res = append(res, "--go_opt", goOpt)
	}
	// grpc插件参数
	for _, goGrpcOpt := range conf.GoGrpcOpt {
		res = append(res, "--go-grpc_opt", goGrpcOpt)
	}
	// go数据目录
	res = append(res, "--go_out", conf.ProtocOut)
	// grpc输出目录
	res = append(res, "--go-grpc_out", conf.ProtocOut)
	if len(conf.HttpServer) > 0 {
		// http输出目录
		res = append(res, "--ceres_out", conf.ProtocOut)
	}
	// 插件
	for _, plugin := range conf.Plugins {
		res = append(res, "--plugin="+plugin)
		// 查看有没有插件参数
	}
	return res
}
