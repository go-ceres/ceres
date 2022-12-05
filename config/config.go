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

package config

import (
	"context"
	"errors"
	"github.com/go-ceres/ceres/logger"
	"golang.org/x/sync/errgroup"
	"reflect"
	"sync"
	"time"
)

var _ Config = (*config)(nil)

var ErrNotFound = errors.New("path not found")

// Observer 配置监听者
type Observer func(string, Value)

// Config 配置接口
type Config interface {
	// Load 加载数据
	Load(source ...Source) error
	// Scan 扫描数据到结构体
	Scan(v interface{}) error
	// Get 获取指定path数据
	Get(path string) Value
	// Watch 监听数据
	Watch(key string, ob Observer) error
	// Close 关闭
	Close() error
}

type config struct {
	opts      Options
	reader    Reader
	cached    *sync.Map
	observers *sync.Map
	watchers  []Watcher
}

// New 创建一个配置实例
func New(opts ...Option) (Config, error) {
	o := Options{
		decoder:  defaultDecoder,
		resolver: DefaultResolver,
	}
	for _, opt := range opts {
		opt(&o)
	}
	c := &config{
		opts:      o,
		reader:    newReader(o),
		cached:    &sync.Map{},
		observers: &sync.Map{},
		watchers:  make([]Watcher, 0),
	}
	return c, c.loadAndWatch(c.opts.sources...)
}

// Load 加载数据源
func (c *config) Load(sources ...Source) error {
	// 加入到资源池
	c.opts.sources = append(c.opts.sources, sources...)
	// 加载新加入的资源
	return c.loadAndWatch(sources...)
}

// run 运行 配置监听
func (c *config) loadAndWatch(sources ...Source) error {
	for _, source := range sources {
		dataSet, err := source.Load()
		if err != nil {
			return err
		}
		for _, v := range dataSet {
			logger.Debugf("config loaded: %s format: %s", v.Key, v.Format)
		}
		if err = c.reader.Merge(dataSet...); err != nil {
			logger.Errorf("failed to merge config source: %v", err)
			return err
		}
		w, err := source.Watch()
		if err != nil {
			logger.Errorf("failed to watch config source: %v", err)
			return err
		}
		c.watchers = append(c.watchers, w)
		go c.watch(w)
	}
	if err := c.reader.Resolve(); err != nil {
		logger.Errorf("failed to resolve config source: %v", err)
		return err
	}
	return nil
}

// Scan 扫描数据结构体
func (c *config) Scan(v interface{}) error {
	data, err := c.reader.Source()
	if err != nil {
		return err
	}
	return unmarshalJSON(data, v)
}

// Get 获取原子值
func (c *config) Get(path string) Value {
	if v, ok := c.cached.Load(path); ok {
		return v.(Value)
	}
	if v, ok := c.reader.Get(path); ok {
		c.cached.Store(path, v)
		return v
	}
	return &atomicValue{}
}

// Watch 配置一个监听者
func (c *config) Watch(key string, ob Observer) error {
	if v := c.Get(key); v.Load() == nil {
		return ErrNotFound
	}
	c.observers.Store(key, ob)
	return nil
}

// Close 关闭配置服务
func (c *config) Close() error {
	var ego errgroup.Group
	for _, watcher := range c.watchers {
		watcher := watcher
		ego.Go(watcher.Stop)
	}
	return ego.Wait()
}

// watch 监听数据
func (c *config) watch(w Watcher) {
	for {
		dataSets, err := w.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Infof("watcher's ctx cancel : %v", err)
				return
			}
			time.Sleep(time.Second)
			logger.Errorf("failed to watch next config: %v", err)
			continue
		}
		if err := c.reader.Merge(dataSets...); err != nil {
			logger.Errorf("failed to merge next config: %v", err)
			continue
		}
		if err := c.reader.Resolve(); err != nil {
			logger.Errorf("failed to resolve next config: %v", err)
			continue
		}
		c.cached.Range(func(key, value interface{}) bool {
			k := key.(string)
			v := value.(Value)
			if n, ok := c.reader.Get(k); ok && reflect.TypeOf(n.Load()) == reflect.TypeOf(v.Load()) && !reflect.DeepEqual(n.Load(), v.Load()) {
				v.Store(n.Load())
				if o, ok := c.observers.Load(k); ok {
					o.(Observer)(k, v)
				}
			}
			return true
		})
	}
}
