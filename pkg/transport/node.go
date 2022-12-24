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

package transport

import (
	"container/list"
	"context"
	"github.com/go-ceres/ceres/pkg/common/errors"
	"math"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	_ WeightNodeBuilder = (*DefaultWeightNodeBuilder)(nil)
	_ Node              = (*defaultNode)(nil)
	// The mean lifetime of `cost`, it reaches its half-life after Tau*ln(2).
	tau = int64(time.Millisecond * 600)
	// if statistic not collected,we add a big lag penalty to endpoint
	penalty = uint64(time.Second * 10)
)

type ReplyMetadata interface {
	Get(key string) string
}

type DoneInfo struct {
	ReplyMetadata       // 响应的metadata数据
	Err           error // 响应错误信息
	BytesSent     bool  // 是否有字节发送过服务器
	BytesReceived bool  // 标识是否有从服务器接受到过数据
}

type DoneFunc func(ctx context.Context, di DoneInfo)

// WeightNodeBuilder 权重节点构建器接口
type WeightNodeBuilder interface {
	Build(node Node) IWeightedNode
}

// DefaultWeightNodeBuilder 默认权重节点构建器实现
type DefaultWeightNodeBuilder struct {
	ErrHandler func(err error) (isErr bool)
}

// Build 构建权重节点
func (w *DefaultWeightNodeBuilder) Build(node Node) IWeightedNode {
	s := &DefaultWeightedNode{
		Node:       node,
		lag:        0,
		success:    1000,
		inflight:   1,
		inflights:  list.New(),
		errHandler: w.ErrHandler,
	}
	return s
}

// IWeightedNode 权重节点
type IWeightedNode interface {
	Node
	Raw() Node                  //返回原始节点
	Weight() float64            // 运行时计算出的权重
	Pick() DoneFunc             // 获取节点
	PickElapsed() time.Duration //最近一次获取节点到现在的时间差
}

// DefaultWeightedNode 默认权重节点
type DefaultWeightedNode struct {
	Node

	// client statistic data
	lag       int64
	success   uint64
	inflight  int64
	inflights *list.List
	// last collected timestamp
	stamp     int64
	predictTs int64
	predict   int64
	// request number in a period time
	reqs int64
	// last lastPick timestamp
	lastPick int64

	errHandler func(err error) (isErr bool)
	lk         sync.RWMutex
}

func (wn *DefaultWeightedNode) health() uint64 {
	return atomic.LoadUint64(&wn.success)
}

func (wn *DefaultWeightedNode) load() (load uint64) {
	now := time.Now().UnixNano()
	avgLag := atomic.LoadInt64(&wn.lag)
	lastPredictTs := atomic.LoadInt64(&wn.predictTs)
	predictInterval := avgLag / 5
	if predictInterval < int64(time.Millisecond*5) {
		predictInterval = int64(time.Millisecond * 5)
	}
	if predictInterval > int64(time.Millisecond*200) {
		predictInterval = int64(time.Millisecond * 200)
	}
	if now-lastPredictTs > predictInterval && atomic.CompareAndSwapInt64(&wn.predictTs, lastPredictTs, now) {
		var (
			total   int64
			count   int
			predict int64
		)
		wn.lk.RLock()
		first := wn.inflights.Front()
		for first != nil {
			lag := now - first.Value.(int64)
			if lag > avgLag {
				count++
				total += lag
			}
			first = first.Next()
		}
		if count > (wn.inflights.Len()/2 + 1) {
			predict = total / int64(count)
		}
		wn.lk.RUnlock()
		atomic.StoreInt64(&wn.predict, predict)
	}

	if avgLag == 0 {
		// penalty is the penalty value when there is no data when the node is just started.
		// The default value is 1e9 * 10
		load = penalty * uint64(atomic.LoadInt64(&wn.inflight))
		return
	}
	predict := atomic.LoadInt64(&wn.predict)
	if predict > avgLag {
		avgLag = predict
	}
	load = uint64(avgLag) * uint64(atomic.LoadInt64(&wn.inflight))
	return
}

// Pick pick a node.
func (wn *DefaultWeightedNode) Pick() DoneFunc {
	now := time.Now().UnixNano()
	atomic.StoreInt64(&wn.lastPick, now)
	atomic.AddInt64(&wn.inflight, 1)
	atomic.AddInt64(&wn.reqs, 1)
	wn.lk.Lock()
	e := wn.inflights.PushBack(now)
	wn.lk.Unlock()
	return func(ctx context.Context, di DoneInfo) {
		wn.lk.Lock()
		wn.inflights.Remove(e)
		wn.lk.Unlock()
		atomic.AddInt64(&wn.inflight, -1)

		now := time.Now().UnixNano()
		// get moving average ratio w
		stamp := atomic.SwapInt64(&wn.stamp, now)
		td := now - stamp
		if td < 0 {
			td = 0
		}
		w := math.Exp(float64(-td) / float64(tau))

		start := e.Value.(int64)
		lag := now - start
		if lag < 0 {
			lag = 0
		}
		oldLag := atomic.LoadInt64(&wn.lag)
		if oldLag == 0 {
			w = 0.0
		}
		lag = int64(float64(oldLag)*w + float64(lag)*(1.0-w))
		atomic.StoreInt64(&wn.lag, lag)

		success := uint64(1000) // error value ,if error set 1
		if di.Err != nil {
			if wn.errHandler != nil && wn.errHandler(di.Err) {
				success = 0
			}
			var netErr net.Error
			if errors.Is(context.DeadlineExceeded, di.Err) || errors.Is(context.Canceled, di.Err) ||
				errors.IsServiceUnavailable(di.Err) || errors.IsGatewayTimeout(di.Err) || errors.As(di.Err, &netErr) {
				success = 0
			}
		}
		oldSuc := atomic.LoadUint64(&wn.success)
		success = uint64(float64(oldSuc)*w + float64(success)*(1.0-w))
		atomic.StoreUint64(&wn.success, success)
	}
}

// Weight is node effective weight.
func (wn *DefaultWeightedNode) Weight() (weight float64) {
	weight = float64(wn.health()*uint64(time.Second)) / float64(wn.load())
	return
}

func (wn *DefaultWeightedNode) PickElapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - atomic.LoadInt64(&wn.lastPick))
}

func (wn *DefaultWeightedNode) Raw() Node {
	return wn.Node
}

// Node 服务节点接口
type Node interface {
	Scheme() string            // 节点协议
	Address() string           // 服务地址
	InitialWeight() *int64     // 初始化权重
	ServiceInfo() *ServiceInfo // 服务注册信息
}

// defaultNode 默认服务器节点实现
type defaultNode struct {
	scheme      string       // 协议
	address     string       // 地址
	weight      *int64       // 权重
	serviceInfo *ServiceInfo // 服务信息
}

// Scheme 节点协议
func (d *defaultNode) Scheme() string {
	return d.scheme
}

// Address 节点地址
func (d *defaultNode) Address() string {
	return d.address
}

// InitialWeight 初始化权重
func (d *defaultNode) InitialWeight() *int64 {
	return d.weight
}

// ServiceInfo 服务信息
func (d *defaultNode) ServiceInfo() *ServiceInfo {
	return d.serviceInfo
}

// NewNode 创建节点
func NewNode(scheme, address string, info *ServiceInfo) Node {
	node := &defaultNode{
		scheme:      scheme,
		address:     address,
		serviceInfo: info,
	}
	if info != nil {
		if str, ok := info.Metadata["weight"]; ok {
			if weight, err := strconv.ParseInt(str, 10, 64); err == nil {
				node.weight = &weight
			}
		}
	}
	return node
}
