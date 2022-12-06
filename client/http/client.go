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
	"bytes"
	"context"
	"fmt"
	"github.com/go-ceres/ceres/errors"
	"github.com/go-ceres/ceres/internal/host"
	"github.com/go-ceres/ceres/middleware"
	"github.com/go-ceres/ceres/selector"
	"github.com/go-ceres/ceres/selector/wrr"
	"github.com/go-ceres/ceres/transport"
	"io"
	"net/http"
)

const SupportPackageIsVersion1 = true

func init() {
	if selector.GetSelectorFactory() == nil {
		selector.SetSelectorFactory(wrr.NewSelectorFactory())
	}
}

// Client http客户端
type Client struct {
	config   *Config
	insecure bool               // 是否是不安全的
	target   *Target            // 目标地址
	resolver *resolver          // 服务发现
	cc       *http.Client       // http客户端
	selector selector.ISelector // 选择器
}

// New 新建一个客户端
func New(c *Config) (*Client, error) {
	if c.TlsConf != nil {
		if tr, ok := c.transport.(*http.Transport); ok {
			tr.TLSClientConfig = c.TlsConf
		}
	}
	insecure := c.TlsConf == nil
	// 解析目标地址
	target, err := ParseTarget(c.Endpoint, insecure)
	if err != nil {
		return nil, err
	}
	// 获取选择器
	selectorBuilder := selector.GetSelectorFactory().Create()
	if selectorBuilder == nil {
		return nil, fmt.Errorf("not found selectorBuild, selectorBuild name is %v", c.Balancer)
	}
	sl := selectorBuilder
	var r *resolver
	if c.discovery != nil {
		if target.Scheme == "discovery" {
			r, err = newResolver(c.ctx, c.logger, c.discovery, target, sl, c.Block, insecure)
			if err != nil {
				return nil, fmt.Errorf("new resolver failed! err：%v", err)
			}
		} else if _, _, err := host.ExtractHostPort(c.Endpoint); err != nil {
			return nil, fmt.Errorf("invalid endpoint format：%v", c.Endpoint)
		}
	}
	return &Client{
		target:   target,
		config:   c,
		insecure: insecure,
		resolver: r,
		cc: &http.Client{
			Timeout:   c.Timeout,
			Transport: c.transport,
		},
		selector: sl,
	}, nil
}

// Invoke 提供给外部调用
func (c *Client) Invoke(ctx context.Context, method, path string, args interface{}, reply interface{}, opts ...CallOption) error {
	var (
		contentType string
		body        io.Reader
	)
	info := newDefaultCallInfo(path)
	for _, opt := range opts {
		if err := opt.before(&info); err != nil {
			return err
		}
	}
	if args != nil {
		data, err := c.config.encoder(ctx, info.contentType, args)
		if err != nil {
			return err
		}
		contentType = info.contentType
		body = bytes.NewReader(data)
	}
	url := fmt.Sprintf("%s://%s%s", c.target.Scheme, c.target.Authority, path)
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	if contentType != "" {
		request.Header.Set("Content-Type", info.contentType)
	}
	if c.config.UserAgent != "" {
		request.Header.Set("User-Agent", c.config.UserAgent)
	}
	ctx = transport.NewClientContext(ctx, &Transport{
		endpoint:      c.config.Endpoint,
		requestHeader: headerInstance(request.Header),
		operation:     info.operation,
		request:       request,
		pathTemplate:  info.pathTemplate,
	})
	return c.invoke(ctx, request, args, reply, info, opts...)
}

// invoke 调用http客户端并处理数据
func (c *Client) invoke(ctx context.Context, req *http.Request, args interface{}, reply interface{}, info callInfo, opts ...CallOption) error {
	handler := func(ctx context.Context, in interface{}) (interface{}, error) {
		response, err := c.do(req.WithContext(ctx))
		if err != nil {
			return nil, err
		}
		if response != nil {
			cs := csAttempt{response: response}
			for _, opt := range opts {
				opt.after(&info, &cs)
			}
		}
		defer response.Body.Close()
		if err := c.config.decoder(ctx, response, reply); err != nil {
			return nil, err
		}
		return reply, nil
	}
	var p selector.Peer
	ctx = selector.NewPeerContext(ctx, &p)
	if len(c.config.middleware) > 0 {
		handler = middleware.Chain(c.config.middleware...)(handler)
	}
	_, err := handler(ctx, args)
	return err
}

// Do 提供给外部调用
func (c *Client) Do(req *http.Request, opts ...CallOption) (*http.Response, error) {
	info := newDefaultCallInfo(req.URL.Path)
	for _, o := range opts {
		if err := o.before(&info); err != nil {
			return nil, err
		}
	}
	return c.do(req)
}

// do 发起请求
func (c *Client) do(req *http.Request) (*http.Response, error) {
	var done func(context.Context, selector.InvokeDoneInfo)
	if c.resolver != nil {
		var (
			err  error
			node selector.INode
		)
		if node, done, err = c.selector.Select(req.Context(), selector.WithNodeFilter(c.config.nodeFilters...)); err != nil {
			return nil, errors.ServiceUnavailable("NODE_NOT_FOUND", err.Error())
		}
		if c.insecure {
			req.URL.Scheme = "http"
		} else {
			req.URL.Scheme = "https"
		}
		req.URL.Host = node.Address()
		req.Host = node.Address()
	}
	resp, err := c.cc.Do(req)
	if err == nil {
		err = c.config.errorDecoder(req.Context(), resp)
	}
	if done != nil {
		done(req.Context(), selector.InvokeDoneInfo{Err: err})
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Close 关闭客户端
func (c *Client) Close() error {
	if c.resolver != nil {
		return c.resolver.Close()
	}
	return nil
}
