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

package http

import (
	"context"
	"crypto/tls"
	"github.com/go-ceres/ceres/internal/matcher"
	"github.com/go-ceres/ceres/pkg/common/codec"
	"github.com/go-ceres/ceres/pkg/common/config"
	"github.com/go-ceres/ceres/pkg/common/logger"
	"github.com/go-ceres/ceres/pkg/transport"
	"time"
)

type TlsConfig struct {
	CertFile string `json:"certFile"` // tls配置
	KeyFile  string `json:"keyFile"`  // tls配置
}

// ServerOption http服务创建参数
type ServerOption func(o *ServerOptions)

// ServerOptions http服务创建参数结构体
type ServerOptions struct {
	Network                       string                              `json:"network"`                       // 网络类型
	Address                       string                              `json:"address"`                       // 连接地址
	Timeout                       time.Duration                       `json:"timeout"`                       // 超时时间
	TlsConf                       *TlsConfig                          `json:"tlsConf"`                       // tls配置
	AllowedMethods                []string                            `json:"allowed_methods"`               // 允许添加的方法
	DisablePrintRoute             bool                                `json:"disablePrintRoute"`             // 禁止打印路由
	Concurrency                   int                                 `json:"concurrency"`                   // 最大并发数.默认:256*1024
	DisableDefaultContentType     bool                                `json:"disableDefaultContentType"`     // 禁用默认的响应类型 默认值：false
	DisableHeaderNamesNormalizing bool                                `json:"disableHeaderNamesNormalizing"` // 禁用默认的标头名称：Content-Type，默认值为：false
	DisableKeepalive              bool                                `json:"disableKeepalive"`              // 禁用保持活动连接, 默认值：false
	MaxRequestBodySize            int                                 `json:"MaxRequestBodySize"`            // 服务接收最大的请求体大小，默认值：4*1024*1024
	ServerHeader                  string                              `json:"serverHeader"`                  // 服务器默认header
	CaseSensitive                 bool                                `json:"case_sensitive"`                // When set to true, enables case sensitive routing.
	ReadTimeout                   time.Duration                       `json:"readTimeout"`                   // 读取超时时间，默认值：unlimited
	WriteTimeout                  time.Duration                       `json:"writeTimeout"`                  // 写入超时时间，默认值：unlimited
	IdleTimeout                   time.Duration                       `json:"idleTimeout"`                   // 启用keep-alive时等待下一个请求的最长时间,默认值：unlimited
	ReadBufferSize                int                                 `json:"readBufferSize"`                // 读取的最大buff大小，默认值：4096
	WriteBufferSize               int                                 `json:"writeBufferSize"`               // 写入的最大buff大小，默认值：4096
	GetOnly                       bool                                `json:"getOnly"`                       // 如果为true，则拒绝所有非get的请求，默认值：false
	ReduceMemoryUsage             bool                                `json:"reduceMemoryUsage"`             // 使用更高的CPU换取内存使用率，默认值：false
	StreamRequestBody             bool                                `json:"streamRequestBody"`             // 启用请求体流，默认值为：false
	DisablePreParseMultipartForm  bool                                `json:"disablePreParseMultipartForm"`  // 如果设置为true，则不会预解析Multipart Form数据。默认值：false
	CompressedFileSuffix          string                              `json:"compressedFileSuffix"`          // 为原始文件添加后缀，默认值：ceres
	UseRawPath                    bool                                `json:"useRawPath"`                    // 启用后，则url。RawPath将用于查找参数。默认值为：false
	UnescapePathValues            bool                                `json:"unescapePathValues"`            // 启用后，则路径值将会被取消转义，默认值：true
	RemoveExtraSlash              bool                                `json:"removeExtraSlash"`              // 启用后，即使使用了额外的斜杠也能匹配到路由,默认值为：false
	RedirectTrailingSlash         bool                                `json:"redirectTrailingSlash"`         // 启用该配置，如果请求了/get/，但只有/get则会重定向到/get 默认值为：true
	RedirectFixedPath             bool                                `json:"redirectFixedPath"`             // 启用该配置，则会尝试修复路由 默认值为：false
	middleware                    matcher.Matcher                     // 中间件
	JSONCodec                     codec.Codec                         // json编解码器
	XmlCodec                      codec.Codec                         // xml编解码器
	ErrorHandler                  func(ctx *Context, err error) error // 错误处理函数
	EncodeResponseFunc            EncodeResponseFunc                  // 请求解码方法
	logger                        *logger.Logger                      //日志组件
}

// ClientOption 客户端创建参数
type ClientOption func(o *ClientOptions)

// ClientOptions 客户端创建参数结构体
type ClientOptions struct {
	Debug          bool                   `json:"debug"`     // 是否开启调试模式，默认值为：false
	Endpoint       string                 `json:"endpoint"`  // 请求地址：默认值为：""
	Block          bool                   `json:"block"`     // 是否
	UserAgent      string                 `json:"userAgent"` // user-agent 请求头，默认：""
	Timeout        time.Duration          `json:"timeout"`   // 请求超时时间，默认值：2s
	TlsConf        *tls.Config            `json:"tlsConf"`   // tls认证信息，默认值为：nil
	decodeResponse DecodeResponseFunc     // 响应信息解码器
	encodeRequest  EncodeRequestFunc      // 请求体编码器
	errorDecoder   DecodeErrorFunc        // 错误解码器
	middleware     []transport.Middleware // 中间件
	nodeFilters    []transport.NodeFilter // 节点过滤器
	discovery      transport.Discover     // 服务发现
	ctx            context.Context
	logger         *logger.Logger // 日志组件
}

type callHook int

const (
	before callHook = iota
	after
)

// CallOption 客户端创建参数
type CallOption func(info *callInfo, hook callHook, ctx interface{}) error

type callInfo struct {
	contentType  string // 文本类型
	operation    string // 对应的grpc方法
	pathTemplate string // 请求地址模板
}

// ========================= 服务端配置参数相关=====================

// DefaultServerOptions 默认配置信息
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		Network:               "tcp",
		Address:               "0.0.0.0:5200",
		AllowedMethods:        []string{MethodGet, MethodPost, MethodConnect, MethodDelete, MethodOptions, MethodHead, MethodPatch, MethodPut, MethodTrace},
		Concurrency:           256 * 1024,
		MaxRequestBodySize:    4 * 1024 * 1024,
		ReadBufferSize:        4096,
		WriteBufferSize:       4096,
		CompressedFileSuffix:  ".ceres.gz",
		UseRawPath:            false,
		UnescapePathValues:    true,
		RemoveExtraSlash:      false,
		RedirectTrailingSlash: true,
		RedirectFixedPath:     false,
		middleware:            matcher.New(),
		JSONCodec:             codec.LoadCodec("json"),
		XmlCodec:              codec.LoadCodec("xml"),
		ErrorHandler:          DefaultErrorHandler,
		EncodeResponseFunc:    DefaultResponseEncode,
		logger:                logger.With(logger.FieldMod("transport.http.server")),
	}
}

// ScanServerRawConfig 扫描配置
func ScanServerRawConfig(key string) *ServerOptions {
	o := DefaultServerOptions()
	if err := config.Get(key).Scan(o); err != nil {
		panic(err)
	}
	return o
}

// ScanServerConfig 标准配置扫描
func ScanServerConfig(name ...string) *ServerOptions {
	key := "application.transport.http.server"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanServerRawConfig(key)
}

func WithServerNetwork(Network string) ServerOption {
	return func(o *ServerOptions) {
		o.Network = Network
	}
}

func WithServerAddress(Address string) ServerOption {
	return func(o *ServerOptions) {
		o.Address = Address
	}
}

func WithServerTimeout(Timeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.Timeout = Timeout
	}
}

func WithServerTlsConf(TlsConf *TlsConfig) ServerOption {
	return func(o *ServerOptions) {
		o.TlsConf = TlsConf
	}
}

func WithServerDisablePrintRoute(DisablePrintRoute bool) ServerOption {
	return func(o *ServerOptions) {
		o.DisablePrintRoute = DisablePrintRoute
	}
}

func WithServerConcurrency(Concurrency int) ServerOption {
	return func(o *ServerOptions) {
		o.Concurrency = Concurrency
	}
}

func WithServerDisableDefaultContentType(DisableDefaultContentType bool) ServerOption {
	return func(o *ServerOptions) {
		o.DisableDefaultContentType = DisableDefaultContentType
	}
}

func WithServerDisableHeaderNamesNormalizing(DisableHeaderNamesNormalizing bool) ServerOption {
	return func(o *ServerOptions) {
		o.DisableHeaderNamesNormalizing = DisableHeaderNamesNormalizing
	}
}

func WithServerDisableKeepalive(DisableKeepalive bool) ServerOption {
	return func(o *ServerOptions) {
		o.DisableKeepalive = DisableKeepalive
	}
}

func WithServerMaxRequestBodySize(MaxRequestBodySize int) ServerOption {
	return func(o *ServerOptions) {
		o.MaxRequestBodySize = MaxRequestBodySize
	}
}

func WithServerServerHeader(ServerHeader string) ServerOption {
	return func(o *ServerOptions) {
		o.ServerHeader = ServerHeader
	}
}

func WithServerReadTimeout(ReadTimeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.ReadTimeout = ReadTimeout
	}
}

func WithServerWriteTimeout(WriteTimeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.WriteTimeout = WriteTimeout
	}
}

func WithServerIdleTimeout(IdleTimeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.IdleTimeout = IdleTimeout
	}
}

func WithServerReadBufferSize(ReadBufferSize int) ServerOption {
	return func(o *ServerOptions) {
		o.ReadBufferSize = ReadBufferSize
	}
}

func WithServerWriteBufferSize(WriteBufferSize int) ServerOption {
	return func(o *ServerOptions) {
		o.WriteBufferSize = WriteBufferSize
	}
}

func WithServerGetOnly(GetOnly bool) ServerOption {
	return func(o *ServerOptions) {
		o.GetOnly = GetOnly
	}
}

func WithServerReduceMemoryUsage(ReduceMemoryUsage bool) ServerOption {
	return func(o *ServerOptions) {
		o.ReduceMemoryUsage = ReduceMemoryUsage
	}
}

func WithServerStreamRequestBody(StreamRequestBody bool) ServerOption {
	return func(o *ServerOptions) {
		o.StreamRequestBody = StreamRequestBody
	}
}

func WithServerDisablePreParseMultipartForm(DisablePreParseMultipartForm bool) ServerOption {
	return func(o *ServerOptions) {
		o.DisablePreParseMultipartForm = DisablePreParseMultipartForm
	}
}

func WithServerCompressedFileSuffix(CompressedFileSuffix string) ServerOption {
	return func(o *ServerOptions) {
		o.CompressedFileSuffix = CompressedFileSuffix
	}
}

func WithServerUseRawPath(UseRawPath bool) ServerOption {
	return func(o *ServerOptions) {
		o.UseRawPath = UseRawPath
	}
}

func WithServerUnescapePathValues(UnescapePathValues bool) ServerOption {
	return func(o *ServerOptions) {
		o.UnescapePathValues = UnescapePathValues
	}
}

func WithServerRemoveExtraSlash(RemoveExtraSlash bool) ServerOption {
	return func(o *ServerOptions) {
		o.RemoveExtraSlash = RemoveExtraSlash
	}
}

func WithServerRedirectTrailingSlash(RedirectTrailingSlash bool) ServerOption {
	return func(o *ServerOptions) {
		o.RedirectTrailingSlash = RedirectTrailingSlash
	}
}

func WithServerRedirectFixedPath(RedirectFixedPath bool) ServerOption {
	return func(o *ServerOptions) {
		o.RedirectFixedPath = RedirectFixedPath
	}
}

func WithServerMiddleware(middlewares ...transport.Middleware) ServerOption {
	return func(o *ServerOptions) {
		o.middleware.Use(middlewares...)
	}
}

func AddServerMiddleware(selector string, middlewares ...transport.Middleware) ServerOption {
	return func(o *ServerOptions) {
		o.middleware.Add(selector, middlewares...)
	}
}

func WithServerJSONCodec(JSONCodec codec.Codec) ServerOption {
	return func(o *ServerOptions) {
		o.JSONCodec = JSONCodec
	}
}

func WithServerXmlCodec(XmlCodec codec.Codec) ServerOption {
	return func(o *ServerOptions) {
		o.XmlCodec = XmlCodec
	}
}

func WithServerErrorHandler(ErrorHandler func(ctx *Context, err error) error) ServerOption {
	return func(o *ServerOptions) {
		o.ErrorHandler = ErrorHandler
	}
}

func WithServerEncodeResponseFunc(EncodeResponseFunc EncodeResponseFunc) ServerOption {
	return func(o *ServerOptions) {
		o.EncodeResponseFunc = EncodeResponseFunc
	}
}

func WithServerLogger(logger *logger.Logger) ServerOption {
	return func(o *ServerOptions) {
		o.logger = logger
	}
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
	return NewWithOptions(o)
}

// =========================== 客户端配置参数相关 ======================

// DefaultClientOptions 默认的客户端配置参数
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Debug:          false,
		Timeout:        time.Second * 2,
		Block:          true,
		encodeRequest:  defaultRequestEncoder,
		decodeResponse: defaultResponseDecoder,
		errorDecoder:   defaultErrorDeCoder,
		ctx:            context.Background(),
		logger:         logger.With(logger.FieldMod("transport.http.client")),
	}
}

// ScanClientRawConfig 扫描客户端参数
func ScanClientRawConfig(key string) *ClientOptions {
	o := DefaultClientOptions()
	if err := config.Get(key).Scan(o); err != nil {
		panic(err)
	}
	return o
}

// ScanClientConfig 扫描客户端配置参数
func ScanClientConfig(name ...string) *ClientOptions {
	key := "application.transport.http.client"
	if len(name) > 0 {
		key = key + "." + name[0]
	}
	return ScanClientRawConfig(key)
}

// WithClientDebug 设置客户端调试信息
func WithClientDebug(debug bool) ClientOption {
	return func(o *ClientOptions) {
		o.Debug = debug
	}
}

// WithClientEndpoint 设置入口地址
func WithClientEndpoint(endpoint string) ClientOption {
	return func(o *ClientOptions) {
		o.Endpoint = endpoint
	}
}

// WithClientUserAgent 设置user-agent
func WithClientUserAgent(userAgent string) ClientOption {
	return func(o *ClientOptions) {
		o.UserAgent = userAgent
	}
}

// WithClientTimeout 设置客户端请求超时时间
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(o *ClientOptions) {
		o.Timeout = timeout
	}
}

// WithClientTlsConfig 设置客户端tls认证信息
func WithClientTlsConfig(conf *tls.Config) ClientOption {
	return func(o *ClientOptions) {
		o.TlsConf = conf
	}
}

// WitClientMiddleware 设置中间件
func WitClientMiddleware(middlewares []transport.Middleware) ClientOption {
	return func(o *ClientOptions) {
		o.middleware = middlewares
	}
}

// WithClientNodeFilters 设置节点过滤器
func WithClientNodeFilters(filters []transport.NodeFilter) ClientOption {
	return func(o *ClientOptions) {
		o.nodeFilters = filters
	}
}

// WithClientDecodeResponseFunc 设置响应解码方法
func WithClientDecodeResponseFunc(fn DecodeResponseFunc) ClientOption {
	return func(o *ClientOptions) {
		o.decodeResponse = fn
	}
}

// WithClientDecodeErrorFunc 错误解码方法
func WithClientDecodeErrorFunc(fn DecodeErrorFunc) ClientOption {
	return func(o *ClientOptions) {
		o.errorDecoder = fn
	}
}

// WithClientDiscovery 设置服务发现
func WithClientDiscovery(discovery transport.Discover) ClientOption {
	return func(o *ClientOptions) {
		o.discovery = discovery
	}
}

// WithClientContext 设置请求上下文
func WithClientContext(ctx context.Context) ClientOption {
	return func(o *ClientOptions) {
		o.ctx = ctx
	}
}

// WithClientLogger 设置日志组件
func WithClientLogger(log *logger.Logger) ClientOption {
	return func(o *ClientOptions) {
		o.logger = log
	}
}

// WithOption 设置参数
func (c *ClientOptions) WithOption(opts ...ClientOption) *ClientOptions {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Build 根据配置参数构建客户端
func (c *ClientOptions) Build() (*Client, error) {
	return NewClientWithOptions(c)
}

// ===================== 请求参数 ======================

func defaultCallInfo(path string) *callInfo {
	return &callInfo{
		contentType:  "application/json",
		operation:    path,
		pathTemplate: path,
	}
}

func WithCallOperation(op string) CallOption {
	return func(info *callInfo, hook callHook, ctx interface{}) error {
		if hook == before {
			info.operation = op
		}
		return nil
	}
}

func WithCallContentType(contentType string) CallOption {
	return func(info *callInfo, hook callHook, ctx interface{}) error {
		if hook == before {
			info.contentType = contentType
		}
		return nil
	}
}

func WithCallPathTemplate(pathTemplate string) CallOption {
	return func(info *callInfo, hook callHook, ctx interface{}) error {
		if hook == before {
			info.pathTemplate = pathTemplate
		}
		return nil
	}
}

func WithCallBeforeHeader(key, value string) CallOption {
	return func(info *callInfo, hook callHook, ctx interface{}) error {
		if hook == before {
			req, ok := ctx.(*Request)
			if ok {
				req.Header.Set(key, value)
			}
		}
		return nil
	}
}

func WithCallAfterHeader(key, value string) CallOption {
	return func(info *callInfo, hook callHook, ctx interface{}) error {
		if hook == after {
			resp, ok := ctx.(*Response)
			if ok {
				resp.Header.Set(key, value)
			}
		}
		return nil
	}
}

type disableLogger struct{}

func (d *disableLogger) Printf(_ string, _ ...interface{}) {}
