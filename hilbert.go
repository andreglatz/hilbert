package hilbert

import (
	"sync"
	"time"
)

type item[T any] struct {
	value T
	timer *time.Timer
}

type cache[T any] struct {
	items map[string]item[T]
	ttl   time.Duration
	mu    sync.RWMutex
}

var _ Cache[any] = &cache[any]{}

func New[T any](ttl ...time.Duration) Cache[T] {
	if len(ttl) <= 0 {
		ttl = []time.Duration{DefaultTTL}
	}

	return &cache[T]{
		items: make(map[string]item[T]),
		ttl:   ttl[0],
	}
}

func (c *cache[T]) Set(key string, value T, ttl ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(ttl) <= 0 {
		ttl = []time.Duration{c.ttl}
	}

	if item, ok := c.items[key]; ok {
		item.timer.Stop()
	}

	timer := time.NewTimer(ttl[0])

	c.items[key] = item[T]{value, timer}

	go c.remove(key, timer)
}

func (c *cache[T]) remove(key string, timer *time.Timer) {
	<-timer.C

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.items[key]
	return v.value, ok
}

func (c *cache[T]) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		item.timer.Stop()
		delete(c.items, key)
	}
}

func (c *cache[T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		item.timer.Stop()
		delete(c.items, key)
	}

}
