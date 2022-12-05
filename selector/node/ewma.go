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

package node

import (
	"container/list"
	"context"
	"github.com/go-ceres/ceres/errors"
	"github.com/go-ceres/ceres/selector"
	"math"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// The mean lifetime of `cost`, it reaches its half-life after Tau*ln(2).
	tau = int64(time.Millisecond * 600)
	// if statistic not collected,we add a big lag penalty to endpoint
	penalty = uint64(time.Second * 10)
)

var (
	_ selector.IWeightedNode        = (*EWMANode)(nil)
	_ selector.IWeightedNodeFactory = (*EWMAWeightedNodeFactory)(nil)
)

// EWMANode is endpoint instance
type EWMANode struct {
	selector.INode
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

// EWMAWeightedNodeFactory ewma 节点的创建工厂
type EWMAWeightedNodeFactory struct {
	ErrHandler func(err error) (isErr bool)
}

// Create 创建权重节点的函数
func (b *EWMAWeightedNodeFactory) Create(n selector.INode) selector.IWeightedNode {
	s := &EWMANode{
		INode:      n,
		lag:        0,
		success:    1000,
		inflight:   1,
		inflights:  list.New(),
		errHandler: b.ErrHandler,
	}
	return s
}

func (n *EWMANode) health() uint64 {
	return atomic.LoadUint64(&n.success)
}

func (n *EWMANode) load() (load uint64) {
	now := time.Now().UnixNano()
	avgLag := atomic.LoadInt64(&n.lag)
	lastPredictTs := atomic.LoadInt64(&n.predictTs)
	predictInterval := avgLag / 5
	if predictInterval < int64(time.Millisecond*5) {
		predictInterval = int64(time.Millisecond * 5)
	}
	if predictInterval > int64(time.Millisecond*200) {
		predictInterval = int64(time.Millisecond * 200)
	}
	if now-lastPredictTs > predictInterval && atomic.CompareAndSwapInt64(&n.predictTs, lastPredictTs, now) {
		var (
			total   int64
			count   int
			predict int64
		)
		n.lk.RLock()
		first := n.inflights.Front()
		for first != nil {
			lag := now - first.Value.(int64)
			if lag > avgLag {
				count++
				total += lag
			}
			first = first.Next()
		}
		if count > (n.inflights.Len()/2 + 1) {
			predict = total / int64(count)
		}
		n.lk.RUnlock()
		atomic.StoreInt64(&n.predict, predict)
	}

	if avgLag == 0 {
		// penalty is the penalty value when there is no data when the node is just started.
		// The default value is 1e9 * 10
		load = penalty * uint64(atomic.LoadInt64(&n.inflight))
		return
	}
	predict := atomic.LoadInt64(&n.predict)
	if predict > avgLag {
		avgLag = predict
	}
	load = uint64(avgLag) * uint64(atomic.LoadInt64(&n.inflight))
	return
}

// Pick 选择节点
func (n *EWMANode) Pick() selector.DoneFunc {
	now := time.Now().UnixNano()
	atomic.StoreInt64(&n.lastPick, now)
	atomic.AddInt64(&n.inflight, 1)
	atomic.AddInt64(&n.reqs, 1)
	n.lk.Lock()
	e := n.inflights.PushBack(now)
	n.lk.Unlock()
	return func(ctx context.Context, idf selector.InvokeDoneInfo) {
		n.lk.Lock()
		n.inflights.Remove(e)
		n.lk.Unlock()
		atomic.AddInt64(&n.inflight, -1)

		now := time.Now().UnixNano()
		// get moving average ratio w
		stamp := atomic.SwapInt64(&n.stamp, now)
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
		oldLag := atomic.LoadInt64(&n.lag)
		if oldLag == 0 {
			w = 0.0
		}
		lag = int64(float64(oldLag)*w + float64(lag)*(1.0-w))
		atomic.StoreInt64(&n.lag, lag)

		success := uint64(1000) // error value ,if error set 1
		if idf.Err != nil {
			if n.errHandler != nil && n.errHandler(idf.Err) {
				success = 0
			}
			var netErr net.Error
			if errors.Is(context.DeadlineExceeded, idf.Err) || errors.Is(context.Canceled, idf.Err) ||
				errors.IsServiceUnavailable(idf.Err) || errors.IsGatewayTimeout(idf.Err) || errors.As(idf.Err, &netErr) {
				success = 0
			}
		}
		oldSuc := atomic.LoadUint64(&n.success)
		success = uint64(float64(oldSuc)*w + float64(success)*(1.0-w))
		atomic.StoreUint64(&n.success, success)
	}
}

// Weight is node effective weight.
func (n *EWMANode) Weight() (weight float64) {
	weight = float64(n.health()*uint64(time.Second)) / float64(n.load())
	return
}

func (n *EWMANode) PickElapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - atomic.LoadInt64(&n.lastPick))
}

func (n *EWMANode) Raw() selector.INode {
	return n.INode
}
