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
	"github.com/go-ceres/ceres/pkg/transport"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
)

const (
	balanceName = "selector"
)

func init() {
	b := base.NewBalancerBuilder(
		balanceName,
		&pickerBuilder{
			selectorBuilder: transport.GetSelectorBuilder(),
		},
		base.Config{
			HealthCheck: true,
		},
	)
	balancer.Register(b)
}

type pickerBuilder struct {
	selectorBuilder transport.SelectorBuilder
}

// Build 构建匹配器
func (pb *pickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	// 获取selector
	if len(info.ReadySCs) == 0 {
		// Block the RPC until a new picker is available via UpdateState().
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	nodes := make([]transport.Node, 0, len(info.ReadySCs))
	for conn, info := range info.ReadySCs {
		ins, _ := info.Address.Attributes.Value("rawServiceInstance").(*transport.ServiceInfo)
		nodes = append(nodes, &grpcNode{
			Node:    transport.NewNode("grpc", info.Address.Addr, ins),
			subConn: conn,
		})
	}
	p := &picker{
		selector: pb.selectorBuilder.Build(),
	}
	p.selector.Store(nodes)
	return p
}

type picker struct {
	selector transport.Selector
}

func (p *picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var filters []transport.NodeFilter
	if md, ok := transport.MetadataFromClientContext(info.Ctx); ok {
		if gmd, ok := md.(*Metadata); ok {
			filters = gmd.NodeFilters()
		}
	}

	n, done, err := p.selector.Select(info.Ctx, transport.WithNodeFilter(filters...))
	if err != nil {
		return balancer.PickResult{}, err
	}

	return balancer.PickResult{
		SubConn: n.(*grpcNode).subConn,
		Done: func(di balancer.DoneInfo) {
			done(info.Ctx, transport.DoneInfo{
				Err:           di.Err,
				BytesSent:     di.BytesSent,
				BytesReceived: di.BytesReceived,
				ReplyMetadata: Trailer(di.Trailer),
			})
		},
	}, nil
}

type Trailer metadata.MD

func (md Trailer) Get(key string) string {
	v := metadata.MD(md).Get(key)
	if len(v) > 0 {
		return v[0]
	}
	return ""
}

type grpcNode struct {
	transport.Node
	subConn balancer.SubConn
}
