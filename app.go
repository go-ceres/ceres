// Copyright 2022. ceres
// Author https://github.com/go-ceres/ceres
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ceres

import (
	"context"
	"github.com/fatih/color"
	"github.com/go-ceres/ceres/internal/cycle"
	"github.com/go-ceres/ceres/internal/signals"
	"github.com/go-ceres/ceres/logger"
	"github.com/go-ceres/ceres/registry"
	"github.com/go-ceres/ceres/server"
	"github.com/go-ceres/ceres/version"
	"go.uber.org/automaxprocs/maxprocs"
	"runtime"
	"sync"
	"time"
)

// AppInfo 应用信息
type AppInfo interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoint() []string
}

func (app *App) ID() string {
	return app.serviceInfo.ID
}

func (app *App) Name() string {
	return app.serviceInfo.Name
}

func (app *App) Version() string {
	return app.serviceInfo.Version
}

func (app *App) Metadata() map[string]string {
	return app.serviceInfo.Metadata
}

func (app *App) Endpoint() []string {
	return app.serviceInfo.Endpoints
}

// App 应用生命周期管理
type App struct {
	conf        *Config               // 应用配置信息
	locker      *sync.RWMutex         // 读写锁
	cycle       *cycle.Cycle          // 生命周期管理
	startupOnce sync.Once             // 启动执行函数
	stopOnce    sync.Once             // 停止执行函数
	serviceInfo *registry.ServiceInfo // 服务信息
	logger      *logger.Helper        // 日志
}

// New 新建应用
func New(c ...*Config) *App {
	conf := DefaultConfig()
	if len(c) > 0 {
		conf = c[0]
	}
	if conf.logger != nil {
		logger.SetLogger(conf.logger)
	}
	return &App{
		conf:   conf,
		cycle:  cycle.NewCycle(),
		locker: &sync.RWMutex{},
		logger: logger.With(logger.FieldMod(ModName)),
	}
}

// initialize 初始化
func (app *App) startup() (err error) {
	app.startupOnce.Do(func() {
		err = app.serialRunner(
			app.printBanner,
			app.initMaxProcs,
		)
	})
	return
}

// printBanner 打印banner
func (app *App) printBanner() error {
	if app.conf.HideBanner {
		return nil
	}
	const banner = `
			  ____ _____ ____  _____ ____
			 / ___| ____|  _ \| ____/ ___|
			| |   |  _| | |_) |  _| \___ \
			| |___| |___|  _ <| |___ ___) |
			 \____|_____|_| \_|_____|____/
			
			ceres@` + version.CeresVersion + `    http://go-ceres.com/


`
	color.Green(banner)
	return nil
}

// InitMaxProcs 设置处理器调用
func (app *App) initMaxProcs() error {
	if maxProcs := app.conf.MaxProc; maxProcs != 0 {
		runtime.GOMAXPROCS(int(maxProcs))
	} else {
		if _, err := maxprocs.Set(); err != nil {
			logger.Panic("auto max procs", logger.FieldMod("application"), logger.FieldError(err))
		}
	}
	return nil
}

// serialRunner 串行执行器
func (app *App) serialRunner(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// Run 启动
func (app *App) Run() error {
	// 启动前初始化
	if err := app.startup(); err != nil {
		return err
	}
	app.waitSignals()
	defer app.clear()

	// 启动服务
	app.cycle.Run(app.startServers)

	if err := <-app.cycle.Wait(); err != nil {
		return err
	}
	app.logger.Info("shutdown ceres, bye!")
	return nil
}

// waitSignals 等待退出命令
func (app *App) waitSignals() {
	signals.Shutdown(func(grace bool) {
		if grace {
			_ = app.GracefulStop(context.TODO())
		} else {
			_ = app.Stop()
		}
	})
}

// startServers 启动服务
func (app *App) startServers() error {
	info, err := app.buildServerInfo()
	if err != nil {
		return err
	}
	app.locker.Lock()
	app.serviceInfo = info
	app.locker.Unlock()

	appCtx := NewContext(context.Background(), app)
	wg := sync.WaitGroup{}
	// 启动前钩子
	app.runHook(BeforeStart)

	for _, srv := range app.conf.servers {
		srv := srv
		wg.Add(1)
		app.cycle.Run(func() error {
			wg.Done()
			return srv.Start(appCtx)
		})
	}
	wg.Wait()

	// 注册服务
	if app.conf.registry != nil {
		ctx, cancel := context.WithTimeout(appCtx, 3*time.Second)
		defer cancel()
		if err = app.conf.registry.Register(ctx, info); err != nil {
			return err
		}
	}

	// 启动后钩子
	app.runHook(AfterStart)

	return nil
}

// buildServerInfo 构建服务信息
func (app *App) buildServerInfo() (*registry.ServiceInfo, error) {
	endpoints := make([]string, 0)
	if len(endpoints) == 0 {
		for _, srv := range app.conf.servers {
			if r, ok := srv.(server.Endpointer); ok {
				e, err := r.Endpoint()
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, e.String())
			}
		}
	}
	info := &registry.ServiceInfo{
		ID:      app.conf.ID,
		Name:    app.conf.Name,
		Version: app.conf.Version,
		Metadata: map[string]string{
			"region": app.conf.Region,
			"zone":   app.conf.Zone,
		},
		Endpoints: endpoints,
	}
	for key, value := range app.conf.Metadata {
		info.Metadata[key] = value
	}
	return info, nil
}

// GracefulStop 优雅停止应用
func (app *App) GracefulStop(ctx context.Context) (err error) {
	app.stopOnce.Do(func() {
		app.runHook(BeforeStop)
		serverInfo := app.serviceInfo
		// 注销服务
		if app.conf.registry != nil && serverInfo != nil {
			ctx, cancel := context.WithTimeout(NewContext(context.Background(), app), app.conf.StopTimeout)
			defer cancel()
			if err = app.conf.registry.Deregister(ctx, serverInfo); err != nil {
				app.logger.Error("stop server error", logger.FieldError(err))
			}
		}
		// 停止服务
		app.locker.RLock()
		for _, s := range app.conf.servers {
			func(s server.Server) {
				app.cycle.Run(func() error {
					return s.GracefulStop(ctx)
				})
			}(s)
		}
		app.locker.RUnlock()

		<-app.cycle.Done()
		app.runHook(AfterStop)
		app.cycle.Close()
	})
	return err
}

// Stop 停止应用
func (app *App) Stop() (err error) {
	app.stopOnce.Do(func() {
		app.runHook(BeforeStop)
		serverInfo := app.serviceInfo
		// 注销服务
		if app.conf.registry != nil && serverInfo != nil {
			ctx, cancel := context.WithTimeout(NewContext(context.Background(), app), app.conf.StopTimeout)
			defer cancel()
			if err = app.conf.registry.Deregister(ctx, serverInfo); err != nil {
				app.logger.Error("stop server error", logger.FieldError(err))
			}
		}
		// 停止服务
		app.locker.RLock()
		for _, s := range app.conf.servers {
			s := s
			app.cycle.Run(func() error {
				return s.Stop(context.Background())
			})
		}
		app.locker.RUnlock()
		<-app.cycle.Done()
		app.runHook(AfterStop)
		app.cycle.Close()
	})
	return err
}

// clear 清除
func (app *App) clear() {
	_ = logger.GetLogger().Sync()
	logger.GetLogger().Close()
}

// runHook 运行钩子
func (app *App) runHook(k HookType) {
	hooks, ok := app.conf.hooks[k]
	if ok {
		for _, hook := range hooks {
			hook()
		}
	}
}

type appKey struct{}

// NewContext 创建附带服务信息的上下文
func NewContext(ctx context.Context, info AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, info)
}

// FromContext 从上下文中获取服务信息
func FromContext(ctx context.Context) (info AppInfo, ok bool) {
	info, ok = ctx.Value(appKey{}).(AppInfo)
	return
}
