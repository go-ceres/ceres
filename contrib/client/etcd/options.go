// Copyright 2023. ceres
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

package etcd

import (
	"github.com/go-ceres/ceres/pkg/common/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Options clientv3.Config

type Option func(o *Options)

func WithAutoSyncInterval(autoSyncInterval time.Duration) Option {
	return func(o *Options) {
		o.AutoSyncInterval = autoSyncInterval
	}
}

func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(o *Options) {
		o.DialTimeout = dialTimeout
	}
}

func WithDialKeepAliveTime(dialKeepAliveTime time.Duration) Option {
	return func(o *Options) {
		o.DialKeepAliveTime = dialKeepAliveTime
	}
}

func WithDialKeepAliveTimeout(dialKeepAliveTimeout time.Duration) Option {
	return func(o *Options) {
		o.DialKeepAliveTimeout = dialKeepAliveTimeout
	}
}

func DefaultOption() *Options {
	return &Options{
		Endpoints: []string{"http://127.0.0.1:12379"},
	}
}

func ScanRawConfig(key string) *Options {
	opts := DefaultOption()
	if err := config.Get(key).Scan(opts); err != nil {
		panic(err)
	}
	return opts
}

func ScanConfig(name ...string) *Options {
	key := "application.transport.client.etcd"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

func (o *Options) WithOptions(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *Options) Build() *Client {
	return NewWithOptions(o)
}
