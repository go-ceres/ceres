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

package app

import (
	"context"
	"github.com/fatih/color"
	"github.com/go-ceres/ceres"
	"github.com/go-ceres/ceres/internal/cycle"
	"github.com/go-ceres/ceres/internal/signals"
	"github.com/go-ceres/ceres/pkg/common/logger"
	"github.com/go-ceres/ceres/pkg/transport"
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

// Application 应用生命周期管理
type Application struct {
	ctx         context.Context        // 应用上下文
	options     *Options               // 应用配置信息
	locker      *sync.RWMutex          // 读写锁
	cycle       *cycle.Cycle           // 生命周期管理
	startupOnce sync.Once              // 启动执行函数
	stopOnce    sync.Once              // 停止执行函数
	serviceInfo *transport.ServiceInfo // 服务信息
	logger      *logger.Logger         // 日志
}

// New 新建应用
func New(opts ...Option) *Application {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return NewWithOptions(options)
}

// NewWithOptions 根据全参数创建应用
func NewWithOptions(opts ...*Options) *Application {
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}
	return &Application{
		options: options,
		ctx:     options.ctx,
		cycle:   cycle.NewCycle(),
		locker:  &sync.RWMutex{},
		logger:  logger.With(logger.FieldMod(ModName)),
	}
}

func (app *Application) ID() string {
	return app.serviceInfo.ID
}

func (app *Application) Name() string {
	return app.serviceInfo.Name
}

func (app *Application) Version() string {
	return app.serviceInfo.Version
}

func (app *Application) Metadata() map[string]string {
	return app.serviceInfo.Metadata
}

func (app *Application) Endpoint() []string {
	return app.serviceInfo.Endpoints
}

// initialize 初始化
func (app *Application) startup() (err error) {
	app.startupOnce.Do(func() {
		err = app.serialRunner(
			app.printBanner,
			app.initMaxProcs,
		)
	})
	return
}

// printBanner 打印banner
func (app *Application) printBanner() error {
	if app.options.HideBanner {
		return nil
	}
	const banner = `
			  ____ _____ ____  _____ ____
			 / ___| ____|  _ \| ____/ ___|
			| |   |  _| | |_) |  _| \___ \
			| |___| |___|  _ <| |___ ___) |
			 \____|_____|_| \_|_____|____/
			
			ceres@` + ceres.Version + `    http://go-ceres.com/


`
	color.Green(banner)
	return nil
}

// InitMaxProcs 设置处理器调用
func (app *Application) initMaxProcs() error {
	if maxProcs := app.options.MaxProc; maxProcs != 0 {
		runtime.GOMAXPROCS(int(maxProcs))
	} else {
		if _, err := maxprocs.Set(); err != nil {
			logger.Panic("auto max procs", logger.FieldMod("application"), logger.FieldError(err))
		}
	}
	return nil
}

// serialRunner 串行执行器
func (app *Application) serialRunner(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// Run 启动
func (app *Application) Run() error {
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
func (app *Application) waitSignals() {
	signals.Shutdown(func(grace bool) {
		_ = app.Stop()
	})
}

// startServers 启动服务
func (app *Application) startServers() error {
	info, err := app.buildServerInfo()
	if err != nil {
		return err
	}
	app.locker.Lock()
	app.serviceInfo = info
	app.locker.Unlock()

	appCtx := NewContext(app.ctx, app)
	wg := sync.WaitGroup{}
	// 启动前钩子
	app.runHook(BeforeStart)
	for _, srv := range app.options.transports {
		srv := srv
		wg.Add(1)
		app.cycle.Run(func() error {
			wg.Done()
			return srv.Start(appCtx)
		})
	}
	wg.Wait()

	// 注册服务
	if app.options.registry != nil {
		ctx, cancel := context.WithTimeout(appCtx, 3*time.Second)
		defer cancel()
		if err = app.options.registry.Register(ctx, info); err != nil {
			return err
		}
	}

	// 启动后钩子
	app.runHook(AfterStart)

	return nil
}

// buildServerInfo 构建服务信息
func (app *Application) buildServerInfo() (*transport.ServiceInfo, error) {
	endpoints := make([]string, 0)
	if len(endpoints) == 0 {
		for _, srv := range app.options.transports {
			if r, ok := srv.(transport.Endpointer); ok {
				e, err := r.Endpoint()
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, e.String())
			}
		}
	}
	info := &transport.ServiceInfo{
		ID:      app.options.ID,
		Name:    app.options.Name,
		Version: app.options.Version,
		Metadata: map[string]string{
			"region": app.options.Region,
			"zone":   app.options.Zone,
		},
		Endpoints: endpoints,
	}
	for key, value := range app.options.Metadata {
		info.Metadata[key] = value
	}
	return info, nil
}

// Stop 停止应用
func (app *Application) Stop() (err error) {
	app.stopOnce.Do(func() {
		// 执行钩子
		app.runHook(BeforeStop)
		// 服务信息
		serverInfo := app.serviceInfo
		// 注销服务
		stopCtx, cancel := context.WithTimeout(NewContext(app.ctx, app), app.options.StopTimeout)
		defer cancel()
		if app.options.registry != nil && serverInfo != nil {
			if err = app.options.registry.Deregister(stopCtx, serverInfo); err != nil {
				app.logger.Error("stop server error", logger.FieldError(err))
			}
		}
		// 停止服务
		app.locker.RLock()
		for _, s := range app.options.transports {
			s := s
			app.cycle.Run(func() error {
				return s.Stop(stopCtx)
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
func (app *Application) clear() {
	_ = logger.GetLogger().Sync()
}

// runHook 运行钩子
func (app *Application) runHook(k HookType) {
	hooks, ok := app.options.hooks[k]
	if ok {
		ctx := NewContext(app.ctx, app)
		for _, hook := range hooks {
			hook(ctx)
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
