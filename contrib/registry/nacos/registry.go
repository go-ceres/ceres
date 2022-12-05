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
	"context"
	"errors"
	"fmt"
	"github.com/go-ceres/ceres/logger"
	"github.com/go-ceres/ceres/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net"
	"net/url"
	"strconv"
	"strings"
)

var _ registry.Registry = (*Registry)(nil)

// Registry nacos注册中心实现
type Registry struct {
	config       *Config
	client       naming_client.INamingClient
	srvConfig    []constant.ServerConfig
	clientConfig *constant.ClientConfig
}

func newRegistry(c *Config) *Registry {
	reg := &Registry{
		config:    c,
		srvConfig: make([]constant.ServerConfig, 0),
	}
	// 初始化配置信息
	for _, addr := range c.Address {
		parse, err := url.Parse(addr)
		if err != nil {
			logger.Panicf("[registry.nacos] url parse panic url：%v", addr)
			return nil
		}
		portStr := parse.Port()
		if len(portStr) == 0 {
			portStr = "80"
		}
		parseInt, err := strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			logger.Panicf("[registry.nacos] url parse panic url：%v", addr)
			return nil
		}
		ips := strings.Split(parse.Host, ":"+parse.Port())
		if len(ips) == 0 {
			logger.Panicf("[registry.nacos] url parse panic url：%v", addr)
			return nil
		}
		reg.srvConfig = append(reg.srvConfig, constant.ServerConfig{
			Scheme:      parse.Scheme,
			IpAddr:      ips[0],
			Port:        uint64(parseInt),
			ContextPath: parse.Path,
		})
	}

	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  c.ClientConfig,
		ServerConfigs: reg.srvConfig,
	})
	if err != nil {
		logger.Panicf("[registry.nacos] nacos client init error：%v")
	}
	reg.client = client
	return reg
}

// GetService 获取服务
func (r *Registry) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInfo, error) {
	res, err := r.client.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   r.config.Group,
		HealthyOnly: true,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*registry.ServiceInfo, 0, len(res))
	for _, in := range res {
		kind := r.config.Kind
		if k, ok := in.Metadata["kind"]; ok {
			kind = k
		}
		items = append(items, &registry.ServiceInfo{
			ID:        in.InstanceId,
			Name:      in.ServiceName,
			Version:   in.Metadata["version"],
			Metadata:  in.Metadata,
			Endpoints: []string{fmt.Sprintf("%s://%s:%d", kind, in.Ip, in.Port)},
		})
	}
	return items, nil
}

// Watch 监听服务
func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return newWatcher(ctx, r.client, serviceName, r.config.Group, r.config.Kind, []string{r.config.Cluster})
}

// Register 注册服务
func (r *Registry) Register(ctx context.Context, service *registry.ServiceInfo) error {
	if service.Name == "" {
		return errors.New("ServiceInfo.name is empty")
	}
	for _, endpoint := range service.Endpoints {
		u, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return err
		}
		p, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		var rmd map[string]string
		if service.Metadata == nil {
			rmd = map[string]string{
				"kind":    u.Scheme,
				"version": service.Version,
			}
		} else {
			rmd = make(map[string]string, len(service.Metadata)+2)
			for k, v := range service.Metadata {
				rmd[k] = v
			}
			rmd["kind"] = u.Scheme
			rmd["version"] = service.Version
		}
		_, e := r.client.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          host,
			Port:        uint64(p),
			ServiceName: service.Name + "." + u.Scheme,
			Weight:      r.config.Weight,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
			Metadata:    rmd,
			ClusterName: r.config.Cluster,
			GroupName:   r.config.Group,
		})
		if e != nil {
			return fmt.Errorf("register serviceInfo err %v,%v", e, endpoint)
		}
	}
	return nil
}

// Deregister 卸载服务
func (r *Registry) Deregister(ctx context.Context, service *registry.ServiceInfo) error {
	for _, endpoint := range service.Endpoints {
		u, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return err
		}
		p, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		if _, err = r.client.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          host,
			Port:        uint64(p),
			ServiceName: service.Name + "." + u.Scheme,
			GroupName:   r.config.Group,
			Cluster:     r.config.Cluster,
			Ephemeral:   true,
		}); err != nil {
			return err
		}
	}
	return nil
}
