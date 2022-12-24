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
	"google.golang.org/grpc/metadata"
	"sync"
)

var (
	_            transport.Metadata = (*Metadata)(nil)
	metadataPool                    = sync.Pool{
		New: func() interface{} {
			return &Metadata{}
		},
	}
)

func AcquireMetadata() *Metadata {
	return metadataPool.Get().(*Metadata)
}

func ReleaseMetadata(md *Metadata) {
	md.reset()
	metadataPool.Put(md)
}

type Metadata struct {
	endpoint      string
	operation     string
	requestHeader headerCarrier
	replyHeader   headerCarrier
	nodeFilters   []transport.NodeFilter
}

func (md *Metadata) reset() {
	md.endpoint = ""
	md.nodeFilters = nil
	md.replyHeader = nil
	md.requestHeader = nil
}

func (md *Metadata) Kind() transport.Kind {
	return KindGrpc
}

func (md *Metadata) Endpoint() string {
	return md.endpoint
}

func (md *Metadata) Operation() string {
	return md.operation
}

func (md *Metadata) RequestHeader() transport.Header {
	return md.requestHeader
}

func (md *Metadata) ReplyHeader() transport.Header {
	return md.replyHeader
}

func (md *Metadata) NodeFilters() []transport.NodeFilter {
	return md.nodeFilters
}

type headerCarrier metadata.MD

func (hc headerCarrier) Get(key string) string {
	vals := metadata.MD(hc).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range metadata.MD(hc) {
		keys = append(keys, k)
	}
	return keys
}

func (hc headerCarrier) Set(key, value string) {
	metadata.MD(hc).Set(key, value)
}
