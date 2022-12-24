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

package nacos

import (
	"github.com/go-ceres/ceres/pkg/common/config"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/file"
	"os"
	"time"
)

type Option func(o *Options)

// Options 配置信息
type Options struct {
	Address       []string       `json:"address"`       // 服务器地址
	Weight        float64        `json:"weight"`        // 初始化权重
	Cluster       string         `json:"cluster"`       // 集群
	Group         string         `json:"group"`         // 分组
	Kind          string         `json:"kind"`          // 默认协议
	ContentPath   string         `json:"contentPath"`   // nacos server contextpath
	ClientOptions *ClientOptions `json:"clientOptions"` // nacos客户端配置
}

// ClientOptions 客户端配置
type ClientOptions struct {
	TimeoutMs            uint64        `json:"timeoutMs"`            //timeout for requesting Nacos server, default value is 10000ms
	BeatInterval         int64         `json:"beatInterval"`         //the time interval for sending beat to server,default value is 5000ms
	NamespaceId          string        `json:"namespaceId"`          //the namespaceId of Nacos.When namespace is public, fill in the blank string here.
	AppName              string        `json:"appName"`              //the appName
	Endpoint             string        `json:"endpoint"`             //the endpoint for get Nacos server addresses
	RegionId             string        `json:"regionId"`             //the regionId for kms
	AccessKey            string        `json:"accessKey"`            //the AccessKey for kms
	SecretKey            string        `json:"secretKey"`            //the SecretKey for kms
	OpenKMS              bool          `json:"openKMS"`              //it's to open kms,default is false. https://help.aliyun.com/product/28933.html
	CacheDir             string        `json:"cacheDir"`             //the directory for persist nacos service info,default value is current path
	UpdateThreadNum      int           `json:"updateThreadNum"`      //the number of gorutine for update nacos service info,default value is 20
	NotLoadCacheAtStart  bool          `json:"notLoadCacheAtStart"`  //not to load persistent nacos service info in CacheDir at start time
	UpdateCacheWhenEmpty bool          `json:"updateCacheWhenEmpty"` //update cache when get empty service instance from server
	Username             string        `json:"username"`             //the username for nacos auth
	Password             string        `json:"password"`             //the password for nacos auth
	LogDir               string        `json:"logDir"`               //the directory for log, default is current path
	RotateTime           string        `json:"rotateTime"`           //the rotate time for log, eg: 30m, 1h, 24h, default is 24h
	MaxAge               int64         `json:"maxAge"`               //the max age of a log file, default value is 3
	LogLevel             string        `json:"logLevel"`             //the level of log, it's must be debug,info,warn,error, default value is info
	ContextPath          string        `json:"contextPath"`          //the nacos server contextpath
	Initial              int           `json:"initial"`              //the sampling initial of log
	Thereafter           int           `json:"thereafter"`           //the sampling thereafter of log
	Tick                 time.Duration `json:"tick"`                 //the sampling tick of log
}

// ToClientConfig 转换成客户端配置
func (c *ClientOptions) ToClientConfig() *constant.ClientConfig {
	return &constant.ClientConfig{
		TimeoutMs:            c.TimeoutMs,
		BeatInterval:         c.BeatInterval,
		NamespaceId:          c.NamespaceId,
		AppName:              c.AppName,
		Endpoint:             c.Endpoint,
		RegionId:             c.RegionId,
		AccessKey:            c.AccessKey,
		SecretKey:            c.SecretKey,
		OpenKMS:              c.OpenKMS,
		CacheDir:             c.CacheDir,
		UpdateThreadNum:      c.UpdateThreadNum,
		NotLoadCacheAtStart:  c.NotLoadCacheAtStart,
		UpdateCacheWhenEmpty: c.UpdateCacheWhenEmpty,
		Username:             c.Username,
		Password:             c.Password,
		LogDir:               c.LogDir,
		RotateTime:           c.RotateTime,
		MaxAge:               c.MaxAge,
		LogLevel:             c.LogLevel,
		LogSampling: &constant.ClientLogSamplingConfig{
			Initial:    c.Initial,
			Thereafter: c.Thereafter,
			Tick:       c.Tick,
		},
		ContextPath: c.ContextPath,
	}
}

// DefaultOptions 默认配置
func DefaultOptions() *Options {
	return &Options{
		Cluster: "DEFAULT",
		Group:   constant.DEFAULT_GROUP,
		Weight:  100,
		Kind:    "grpc",
		ClientOptions: &ClientOptions{
			TimeoutMs:            10 * 1000,
			BeatInterval:         5 * 1000,
			OpenKMS:              false,
			CacheDir:             file.GetCurrentPath() + string(os.PathSeparator) + "cache",
			UpdateThreadNum:      20,
			NotLoadCacheAtStart:  true,
			UpdateCacheWhenEmpty: false,
			LogDir:               file.GetCurrentPath() + string(os.PathSeparator) + "log",
			RotateTime:           "24h",
			MaxAge:               3,
			LogLevel:             "info",
		},
	}
}

// ScanRawConfig 扫描无封装key
func ScanRawConfig(key string) *Options {
	conf := DefaultOptions()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 扫描配置
func ScanConfig(name ...string) *Options {
	key := "application.transport.registry.nacos"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanRawConfig(key)
}

// WithAddress nacos链接
func WithAddress(address ...string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

// WithWeight 初始化权重
func WithWeight(weight float64) Option {
	return func(o *Options) {
		o.Weight = weight
	}
}

// WithCluster 集群名称
func WithCluster(cluster string) Option {
	return func(o *Options) {
		o.Cluster = cluster
	}
}

// WithGroup 分组
func WithGroup(group string) Option {
	return func(o *Options) {
		o.Group = group
	}
}

func WithClientTimeoutMs(TimeoutMs uint64) Option {
	return func(o *Options) {
		o.ClientOptions.TimeoutMs = TimeoutMs
	}
}

func WithClientBeatInterval(BeatInterval int64) Option {
	return func(o *Options) {
		o.ClientOptions.BeatInterval = BeatInterval
	}
}

func WithClientNamespaceId(NamespaceId string) Option {
	return func(o *Options) {
		o.ClientOptions.NamespaceId = NamespaceId
	}
}

func WithClientAppName(AppName string) Option {
	return func(o *Options) {
		o.ClientOptions.AppName = AppName
	}
}

func WithClientEndpoint(Endpoint string) Option {
	return func(o *Options) {
		o.ClientOptions.Endpoint = Endpoint
	}
}

func WithClientRegionId(RegionId string) Option {
	return func(o *Options) {
		o.ClientOptions.RegionId = RegionId
	}
}

func WithClientAccessKey(AccessKey string) Option {
	return func(o *Options) {
		o.ClientOptions.AccessKey = AccessKey
	}
}

func WithClientSecretKey(SecretKey string) Option {
	return func(o *Options) {
		o.ClientOptions.SecretKey = SecretKey
	}
}

func WithClientOpenKMS(OpenKMS bool) Option {
	return func(o *Options) {
		o.ClientOptions.OpenKMS = OpenKMS
	}
}

func WithClientCacheDir(CacheDir string) Option {
	return func(o *Options) {
		o.ClientOptions.CacheDir = CacheDir
	}
}

func WithClientUpdateThreadNum(UpdateThreadNum int) Option {
	return func(o *Options) {
		o.ClientOptions.UpdateThreadNum = UpdateThreadNum
	}
}

func WithClientNotLoadCacheAtStart(NotLoadCacheAtStart bool) Option {
	return func(o *Options) {
		o.ClientOptions.NotLoadCacheAtStart = NotLoadCacheAtStart
	}
}

func WithClientUpdateCacheWhenEmpty(UpdateCacheWhenEmpty bool) Option {
	return func(o *Options) {
		o.ClientOptions.UpdateCacheWhenEmpty = UpdateCacheWhenEmpty
	}
}

func WithClientUsername(Username string) Option {
	return func(o *Options) {
		o.ClientOptions.Username = Username
	}
}

func WithClientPassword(Password string) Option {
	return func(o *Options) {
		o.ClientOptions.Password = Password
	}
}

func WithClientLogDir(LogDir string) Option {
	return func(o *Options) {
		o.ClientOptions.LogDir = LogDir
	}
}

func WithClientRotateTime(RotateTime string) Option {
	return func(o *Options) {
		o.ClientOptions.RotateTime = RotateTime
	}
}

func WithClientMaxAge(MaxAge int64) Option {
	return func(o *Options) {
		o.ClientOptions.MaxAge = MaxAge
	}
}

func WithClientLogLevel(LogLevel string) Option {
	return func(o *Options) {
		o.ClientOptions.LogLevel = LogLevel
	}
}

func WithClientContextPath(ContextPath string) Option {
	return func(o *Options) {
		o.ClientOptions.ContextPath = ContextPath
	}
}

func WithClientInitial(Initial int) Option {
	return func(o *Options) {
		o.ClientOptions.Initial = Initial
	}
}

func WithClientThereafter(Thereafter int) Option {
	return func(o *Options) {
		o.ClientOptions.Thereafter = Thereafter
	}
}

func WithClientTick(Tick time.Duration) Option {
	return func(o *Options) {
		o.ClientOptions.Tick = Tick
	}
}

// WithOptions 手动设置参数
func (o *Options) WithOptions(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Build 构建注册中心
func (o *Options) Build() *Registry {
	return NewWithOptions(o)
}
