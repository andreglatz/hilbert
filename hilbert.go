package hilbert

import (
	"sync"
	"time"
)

type cache[T any] struct {
	items map[string]T
	ttl   time.Duration
	mu    sync.RWMutex
}

var _ Cache[any] = &cache[any]{}

func New[T any](ttl ...time.Duration) Cache[T] {
	if len(ttl) <= 0 {
		ttl = []time.Duration{DefaultTTL}
	}

	return &cache[T]{
		items: make(map[string]T),
		ttl:   ttl[0],
	}
}

func (c *cache[T]) Set(key string, value T, ttl ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = value

	if len(ttl) <= 0 {
		ttl = []time.Duration{c.ttl}
	}

	if ttl[0] != -1 {
		go c.remove(key, ttl[0])
	}
}

func (c *cache[T]) remove(key string, ttl time.Duration) {
	time.Sleep(ttl)

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.items[key]
	return v, ok
}
