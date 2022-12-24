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

package grpc

import (
	"crypto/tls"
	"github.com/go-ceres/ceres/internal/matcher"
	"github.com/go-ceres/ceres/pkg/common/config"
	"github.com/go-ceres/ceres/pkg/common/logger"
	"github.com/go-ceres/ceres/pkg/transport"
	"google.golang.org/grpc"
	"time"
)

// ModName 模块名称
const ModName = "transport.grpc"

// ================= 服务端参数 ===================

// ServerOptions 参数信息
type ServerOptions struct {
	Network            string                         // net.listen network
	Address            string                         // 服务地址
	Timeout            time.Duration                  // 超时时间
	SlowQueryThreshold time.Duration                  // 在debug模式下慢查询阈值,如果请求到响应超过此值，则会打印日志
	TlsConf            *tls.Config                    // tls配置信息
	Reflection         bool                           // 是否反射服务
	Health             bool                           // 是否设置健康服务
	Debug              bool                           // 是否开启调试模式
	middleware         matcher.Matcher                // 中间件
	unaryInts          []grpc.UnaryServerInterceptor  // grpc服务拦截器
	streamInts         []grpc.StreamServerInterceptor // 数据流服务拦截器
	grpcOpts           []grpc.ServerOption            // grpc服务额外参数
	logger             *logger.Logger                 //日志组件
}

// ServerOption 参数信息
type ServerOption func(o *ServerOptions)

// DefaultServerOptions 默认的参数信息
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		Network:            "tcp",
		Address:            "127.0.0.1:5201",
		SlowQueryThreshold: 3 * time.Second,
		Debug:              false,
		middleware:         matcher.New(),
		logger:             logger.With(logger.FieldMod(ModName)),
	}
}

// ScanServerRawConfig 扫描配置
func ScanServerRawConfig(key string) *ServerOptions {
	conf := DefaultServerOptions()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanServerConfig 标准配置扫描
func ScanServerConfig(name ...string) *ServerOptions {
	key := "application.transport.grpc.server"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanServerRawConfig(key)
}

// WithOption 设置参数
func (o *ServerOptions) WithOption(opts ...ServerOption) *ServerOptions {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Build 构建grpc服务
func (o *ServerOptions) Build() *Server {
	return NewServerWithOptions(o)
}

// WithServerNetWork 设置网络
func WithServerNetWork(network string) ServerOption {
	return func(o *ServerOptions) {
		o.Network = network
	}
}

// WithServerAddress 设置地址
func WithServerAddress(address string) ServerOption {
	return func(o *ServerOptions) {
		o.Address = address
	}
}

// WithServerTimeout 设置超时时间
func WithServerTimeout(timeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.Timeout = timeout
	}
}

// WithServerSlowQueryThreshold 慢查询时间，超过此时间会打印慢查询日志
func WithServerSlowQueryThreshold(timeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.SlowQueryThreshold = timeout
	}
}

// WithServerTlsConfig 设置tls配置
func WithServerTlsConfig(tlsConfig *tls.Config) ServerOption {
	return func(o *ServerOptions) {
		o.TlsConf = tlsConfig
	}
}

// WithServerReflection 设置是否反射服务
func WithServerReflection(reflection bool) ServerOption {
	return func(o *ServerOptions) {
		o.Reflection = reflection
	}
}

// WithServerHealth 设置是否健康检查
func WithServerHealth(health bool) ServerOption {
	return func(o *ServerOptions) {
		o.Health = health
	}
}

// WithServerDebug 是否开启调试
func WithServerDebug(debug bool) ServerOption {
	return func(o *ServerOptions) {
		o.Debug = debug
	}
}

// WithServerMiddleware 添加中间件
func WithServerMiddleware(mws ...transport.Middleware) ServerOption {
	return func(o *ServerOptions) {
		o.middleware.Use(mws...)
	}
}

// AddServerMiddleware 添加中间件
func AddServerMiddleware(selector string, mw transport.Middleware) ServerOption {
	return func(o *ServerOptions) {
		o.middleware.Add(selector, mw)
	}
}

// WithServerUnaryServerInterceptor 设置grpc中间件
func WithServerUnaryServerInterceptor(ins ...grpc.UnaryServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.unaryInts = ins
	}
}

// AddServerUnaryServerInterceptor 添加grpc中间件
func AddServerUnaryServerInterceptor(ins ...grpc.UnaryServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.unaryInts = append(o.unaryInts, ins...)
	}
}

// WithServerStreamServerInterceptor 设置grpc流中间件
func WithServerStreamServerInterceptor(ins ...grpc.StreamServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.streamInts = ins
	}
}

// AddServerStreamServerInterceptor 添加grpc中间件
func AddServerStreamServerInterceptor(ins ...grpc.StreamServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.streamInts = append(o.streamInts, ins...)
	}
}

// WithServerGrpcOption 设置grpc额外参数
func WithServerGrpcOption(opts ...grpc.ServerOption) ServerOption {
	return func(o *ServerOptions) {
		o.grpcOpts = opts
	}
}

// AddServerGrpcOption 添加grpc额外参数
func AddServerGrpcOption(opts ...grpc.ServerOption) ServerOption {
	return func(o *ServerOptions) {
		o.grpcOpts = append(o.grpcOpts, opts...)
	}
}

// WithServerLogger 设置日志组件
func WithServerLogger(log *logger.Logger) ServerOption {
	return func(o *ServerOptions) {
		o.logger = log
	}
}

// ====================== 客户端参数 =================

// ClientOption 创建客户端的参数
type ClientOption func(o *ClientOptions)

// ClientOptions 客户端参数信息
type ClientOptions struct {
	Debug        bool                          `json:"debug"`       // 是否调试模式
	Endpoint     string                        `json:"endpoint"`    // 连接地址
	Block        bool                          `json:"block"`       // 是否一直等待连接成功
	Insecure     bool                          `json:"insecure"`    // 是否忽略安全
	TlsConfig    *tls.Config                   `json:"tlsConfig"`   // 安全认证
	Timeout      time.Duration                 `json:"timeout"`     // 超时
	DialTimeout  time.Duration                 `json:"dialTimeout"` // 调用超时
	OnDialError  string                        `json:"OnDialError"` // 构建错误处理 panic | error
	Balancer     string                        // 负载均衡器名称
	Selector     string                        `json:"selector"` // 选择器名称
	discovery    transport.Discover            // 服务发现
	middleware   []transport.Middleware        // 中间件
	interceptors []grpc.UnaryClientInterceptor // 拦截器
	dialOpts     []grpc.DialOption             // 调用参数
	filters      []transport.NodeFilter        // 节点选择器
	logger       *logger.Logger                // 日志
}

func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Endpoint:    "",
		Block:       true,
		Timeout:     3 * time.Second,
		DialTimeout: 3 * time.Second,
		Selector:    "p2c",
		Balancer:    balanceName,
		Insecure:    true,
		Debug:       false,
		logger:      logger.With(logger.FieldMod(ModName)),
	}
}

// ScanClientRawConfig 扫描配置文件
func ScanClientRawConfig(Key string) *ClientOptions {
	options := DefaultClientOptions()
	if err := config.Get(Key).Scan(options); err != nil {
		panic(err)
	}
	return options
}

// ScanClientConfig 标准扫描
func ScanClientConfig(name ...string) *ClientOptions {
	key := "application.transport.grpc.client"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanClientRawConfig(key)
}

func WithClientClientDebug(debug bool) ClientOption {
	return func(o *ClientOptions) {
		o.Debug = debug
	}
}

func WithClientEndpoint(Endpoint string) ClientOption {
	return func(o *ClientOptions) {
		o.Endpoint = Endpoint
	}
}

func WithClientBlock(Block bool) ClientOption {
	return func(o *ClientOptions) {
		o.Block = Block
	}
}

func WithClientInsecure(Insecure bool) ClientOption {
	return func(o *ClientOptions) {
		o.Insecure = Insecure
	}
}

func WithClientTlsConfig(TlsConfig *tls.Config) ClientOption {
	return func(o *ClientOptions) {
		o.TlsConfig = TlsConfig
	}
}

func WithClientTimeout(Timeout time.Duration) ClientOption {
	return func(o *ClientOptions) {
		o.Timeout = Timeout
	}
}

func WithClientDialTimeout(DialTimeout time.Duration) ClientOption {
	return func(o *ClientOptions) {
		o.DialTimeout = DialTimeout
	}
}

func WithClientOnDialError(OnDialError string) ClientOption {
	return func(o *ClientOptions) {
		o.OnDialError = OnDialError
	}
}

func WithClientBalancer(Balancer string) ClientOption {
	return func(o *ClientOptions) {
		o.Balancer = Balancer
	}
}

func WithClientSelector(Selector string) ClientOption {
	return func(o *ClientOptions) {
		o.Selector = Selector
	}
}

func WithClientDiscovery(discovery transport.Discover) ClientOption {
	return func(o *ClientOptions) {
		o.discovery = discovery
	}
}

func WithClientMiddleware(middleware []transport.Middleware) ClientOption {
	return func(o *ClientOptions) {
		o.middleware = middleware
	}
}

func WithClientInterceptors(interceptors []grpc.UnaryClientInterceptor) ClientOption {
	return func(o *ClientOptions) {
		o.interceptors = interceptors
	}
}

func WithClientDialOpts(dialOpts []grpc.DialOption) ClientOption {
	return func(o *ClientOptions) {
		o.dialOpts = dialOpts
	}
}

func WithClientFilters(filters []transport.NodeFilter) ClientOption {
	return func(o *ClientOptions) {
		o.filters = filters
	}
}

func WithClientLogger(logger *logger.Logger) ClientOption {
	return func(o *ClientOptions) {
		o.logger = logger
	}
}

// WithOptions 手动设置参数
func (o *ClientOptions) WithOptions(opts ...ClientOption) *ClientOptions {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *ClientOptions) Build() (*grpc.ClientConn, error) {
	return NewClientWithOptions(o)
}
