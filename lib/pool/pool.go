// Copyright 2024 openGemini Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pool

import (
	"sync"
)

type CachePool[T any] struct {
	pool         sync.Pool
	capacityChan chan struct{}
	newFunc      func() T
}

func NewCachePool[T any](newFunc func() T, maxSize int) *CachePool[T] {
	return &CachePool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return newFunc()
			},
		},
		capacityChan: make(chan struct{}, maxSize),
		newFunc:      newFunc,
	}
}

func (c *CachePool[T]) Get() T {
	select {
	case <-c.capacityChan:
		item := c.pool.Get()

		return item.(T)
	default:
		return c.newFunc()
	}
}

func (c *CachePool[T]) Put(x T) {
	select {
	case c.capacityChan <- struct{}{}:
		c.pool.Put(x)
	default:
		// Pool is full, discard the item
	}
}
