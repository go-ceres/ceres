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

package cycle

import (
	"sync"
	"sync/atomic"
)

// Cycle 生命周期管理
type Cycle struct {
	locker  *sync.Mutex
	wg      *sync.WaitGroup
	done    chan struct{}
	quit    chan error
	closing uint32
	waiting uint32
}

// NewCycle 创建
func NewCycle() *Cycle {
	return &Cycle{
		locker:  &sync.Mutex{},
		wg:      &sync.WaitGroup{},
		done:    make(chan struct{}),
		quit:    make(chan error),
		closing: 0,
		waiting: 0,
	}
}

// Run 运行
func (c *Cycle) Run(fn func() error) {
	c.locker.Lock()
	//todo add check options panic before waiting
	defer c.locker.Unlock()
	c.wg.Add(1)
	go func(c *Cycle) {
		defer c.wg.Done()
		if err := fn(); err != nil {
			c.quit <- err
		}
	}(c)
}

// Done 结束通道
func (c *Cycle) Done() <-chan struct{} {
	if atomic.CompareAndSwapUint32(&c.waiting, 0, 1) {
		go func(c *Cycle) {
			c.locker.Lock()
			defer c.locker.Unlock()
			c.wg.Wait()
			close(c.done)
		}(c)
	}
	return c.done
}

// Close 手动关闭
func (c *Cycle) Close() {
	c.locker.Lock()
	defer c.locker.Unlock()
	if atomic.CompareAndSwapUint32(&c.closing, 0, 1) {
		close(c.quit)
	}
}

// DoneAndClose 结束并关闭
func (c *Cycle) DoneAndClose() {
	<-c.Done()
	c.Close()
}

// Wait 一直等待，一直到有错误
func (c *Cycle) Wait() <-chan error {
	return c.quit
}
