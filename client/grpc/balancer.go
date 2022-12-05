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
	"github.com/go-ceres/ceres/registry"
	"github.com/go-ceres/ceres/selector"
	"github.com/go-ceres/ceres/transport"
	trGrpc "github.com/go-ceres/ceres/transport/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
)

const (
	balancerName = "selector"
)

var (
	_ base.PickerBuilder = (*balancerBuilder)(nil)
	_ balancer.Picker    = (*balancerPicker)(nil)
)

func init() {
	b := base.NewBalancerBuilder(
		balancerName,
		&balancerBuilder{
			builder: selector.GetSelectorFactory(),
		},
		base.Config{HealthCheck: true},
	)
	balancer.Register(b)
}

type balancerBuilder struct {
	builder selector.ISelectorFactory
}

// Build creates a grpc Picker.
func (b *balancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		// Block the RPC until a new picker is available via UpdateState().
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	nodes := make([]selector.INode, 0, len(info.ReadySCs))
	for conn, info := range info.ReadySCs {
		ins, _ := info.Address.Attributes.Value("rawServiceInstance").(*registry.ServiceInfo)
		nodes = append(nodes, &grpcNode{
			INode:   selector.NewNode("grpc", info.Address.Addr, ins),
			subConn: conn,
		})
	}
	p := &balancerPicker{
		selector: b.builder.Create(),
	}
	p.selector.Apply(nodes)
	return p
}

// balancerPicker is a grpc picker.
type balancerPicker struct {
	selector selector.ISelector
}

// Pick pick instances.
func (p *balancerPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var filters []selector.NodeFilter
	if tr, ok := transport.FromClientContext(info.Ctx); ok {
		if gtr, ok := tr.(*trGrpc.Transport); ok {
			filters = gtr.NodeFilters()
		}
	}

	n, done, err := p.selector.Select(info.Ctx, selector.WithNodeFilter(filters...))
	if err != nil {
		return balancer.PickResult{}, err
	}

	return balancer.PickResult{
		SubConn: n.(*grpcNode).subConn,
		Done: func(di balancer.DoneInfo) {
			done(info.Ctx, selector.InvokeDoneInfo{
				Err:              di.Err,
				BytesSent:        di.BytesSent,
				BytesReceived:    di.BytesReceived,
				ResponseMetadata: Trailer(di.Trailer),
			})
		},
	}, nil
}

// Trailer is a grpc trailer MD.
type Trailer metadata.MD

// Get get a grpc trailer value.
func (t Trailer) Get(k string) string {
	v := metadata.MD(t).Get(k)
	if len(v) > 0 {
		return v[0]
	}
	return ""
}

type grpcNode struct {
	selector.INode
	subConn balancer.SubConn
}
