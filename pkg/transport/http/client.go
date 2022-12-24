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
	"fmt"
	"github.com/go-ceres/ceres/internal/bytesconv"
	"github.com/go-ceres/ceres/pkg/common/errors"
	"github.com/go-ceres/ceres/pkg/transport"
	"github.com/valyala/fasthttp"
	"net"
	"strconv"
	"strings"
	"sync"
)

var (
	// 请求对象池
	requestPool = sync.Pool{
		New: func() interface{} {
			return &Request{}
		},
	}
	// 响应对象池
	responsePool = sync.Pool{
		New: func() interface{} {
			return &Response{}
		},
	}
)

type RequestHeader = fasthttp.RequestHeader

type ResponseHeader = fasthttp.ResponseHeader

type Request = fasthttp.Request

type Response = fasthttp.Response

type Args = fasthttp.Args

type RetryIfFunc = fasthttp.RetryIfFunc

// AcquireRequest 借用请求对象
func AcquireRequest() *Request {
	return requestPool.Get().(*Request)
}

// ReleaseRequest 归还请求对象
func ReleaseRequest(req *Request) {
	req.Reset()
	requestPool.Put(req)
}

// AcquireResponse 借用响应对象
func AcquireResponse() *Response {
	return responsePool.Get().(*Response)
}

// ReleaseResponse 归还响应对象
func ReleaseResponse(resp *Response) {
	resp.Reset()
	responsePool.Put(resp)
}

// Client 客户端结构体
type Client struct {
	options  *ClientOptions
	target   *Target
	insecure bool
	resolver *resolver
	selector transport.Selector
	cc       *fasthttp.HostClient
}

// NewClient 创建客户端
func NewClient(opts ...ClientOption) (*Client, error) {
	options := DefaultClientOptions()
	for _, opt := range opts {
		opt(options)
	}
	return NewClientWithOptions(options)
}

// NewClientWithOptions 创建客户端根据配置参数
func NewClientWithOptions(options *ClientOptions) (*Client, error) {
	insecure := options.TlsConf == nil
	target, err := parseTarget(options.Endpoint, insecure)
	if err != nil {
		return nil, err
	}
	var resolver *resolver
	var selector transport.Selector
	// 如果有复制均衡
	if options.discovery != nil {
		selector = transport.GetSelectorBuilder().Build()
		resolver, err = newResolver(options.ctx, options.logger, options.discovery, target, selector, options.Block, insecure)
	}
	return &Client{
		target:   target,
		insecure: insecure,
		resolver: resolver,
		selector: selector,
		options:  options,
		cc: &fasthttp.HostClient{
			TLSConfig: options.TlsConf,
			IsTLS:     !insecure,
			Name:      options.UserAgent,
		},
	}, nil
}

// Invoke ...
func (c *Client) Invoke(ctx context.Context, method, path string, args interface{}, reply interface{}, opts ...CallOption) (err error) {
	// 解码出request信息
	marshalReq, err := defaultBinding.Marshal(path, args)
	if err != nil {
		return err
	}
	var (
		contentType string
		body        []byte
	)
	req := AcquireRequest()
	defer ReleaseRequest(req)
	resp := AcquireResponse()
	defer ReleaseResponse(resp)
	info := defaultCallInfo(path)
	for _, opt := range opts {
		if err := opt(info, before, req); err != nil {
			return err
		}
	}
	// 处理路径参数
	path = marshalReq.GetPath()
	// 处理header
	headerParams := marshalReq.GetHeader()
	if len(headerParams) > 0 {
		for key, values := range headerParams {
			v := ""
			if len(values) > 0 {
				v = values[0]
			}
			req.Header.Set(key, v)
		}
	}

	// 设置cookie
	cookieParams := marshalReq.GetCookies()
	if len(cookieParams) > 0 {
		for _, cookie := range cookieParams {
			req.Header.SetCookie(cookie.Name, cookie.Value)
		}
	}

	// 处理query
	queryParams := marshalReq.GetQuery()
	if len(queryParams) > 0 {
		if query := queryParams.Encode(); query != "" {
			path = path + "?" + query
		}
	}
	// 设置body
	if marshalReq.HasBody() {
		// 设置请求体类型
		if contentType != "" {
			req.Header.Set(HeaderContentType, info.contentType)
		}
		// 获取并设置body数据
		body, err = c.options.encodeRequest(ctx, req, marshalReq.Body())
		if err != nil {
			return err
		}
	}
	// 设置url
	url := fmt.Sprintf("%s://%s%s", c.target.Scheme, c.target.Authority, path)
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.SetBody(body)
	ctx = transport.NewMetadataClientContext(ctx, &Metadata{
		endpoint:     c.options.Endpoint,
		request:      req,
		operation:    info.operation,
		pathTemplate: info.pathTemplate,
		response:     resp,
	})
	return c.invoke(ctx, req, resp, args, reply, info, opts...)
}

func (c *Client) invoke(ctx context.Context, req *Request, resp *Response, args interface{}, reply interface{}, info *callInfo, opts ...CallOption) error {
	h := func(ctx context.Context, in interface{}) (interface{}, error) {
		err := c.do(ctx, req, resp)
		if err != nil {
			return nil, err
		}
		for _, opt := range opts {
			err := opt(info, after, resp)
			if err != nil {
				return nil, err
			}
		}
		if err := c.options.decodeResponse(ctx, resp, reply); err != nil {
			return nil, err
		}
		return reply, nil
	}
	var p transport.Peer
	ctx = transport.NewPeerContext(ctx, &p)
	if len(c.options.middleware) > 0 {
		h = transport.MiddlewareChain(c.options.middleware...)(h)
	}
	_, err := h(ctx, args)
	return err
}

// Do 发送请求
func (c *Client) Do(ctx context.Context, req *Request, resp *Response, opts ...CallOption) error {
	info := defaultCallInfo(bytesconv.BytesToString(req.URI().Path()))
	for _, opt := range opts {
		if err := opt(info, before, req); err != nil {
			return err
		}
	}
	return c.do(ctx, req, resp)
}

func (c *Client) do(ctx context.Context, req *Request, resp *Response) error {
	var done func(context.Context, transport.DoneInfo)
	if c.resolver != nil {
		var (
			err  error
			node transport.Node
		)
		if node, done, err = c.selector.Select(ctx, transport.WithNodeFilter(c.options.nodeFilters...)); err != nil {
			return errors.ServiceUnavailable("NODE_NOT_FOUND", err.Error())
		}
		if c.insecure {
			req.URI().SetScheme("http")
		} else {
			req.URI().SetScheme("https")
		}
		req.URI().SetHost(node.Address())
		req.SetHost(node.Address())
	}
	c.cc.Addr = addMissingPort(bytesconv.BytesToString(req.URI().Host()), !c.insecure)
	err := c.cc.Do(req, resp)
	if err == nil {
		err = c.options.errorDecoder(ctx, resp)
	}
	if done != nil {
		done(ctx, transport.DoneInfo{Err: err})
	}
	if err != nil {
		return err
	}
	return nil
}

func addMissingPort(addr string, isTLS bool) string {
	n := strings.Index(addr, ":")
	if n >= 0 {
		return addr
	}
	port := 80
	if isTLS {
		port = 443
	}
	return net.JoinHostPort(addr, strconv.Itoa(port))
}
