package hilbert

import (
	"sync"
	"time"
)

type Cache[T any] interface {
	Set(key string, value T, ttl time.Duration)
	Get(key string) (T, bool)
}

type cache[T any] struct {
	items map[string]T
	mu    sync.RWMutex
}

var _ Cache[any] = &cache[any]{}

func New[T any]() Cache[T] {
	return &cache[T]{
		items: make(map[string]T),
	}
}

func (c *cache[T]) Set(key string, value T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = value
	go c.remove(key, ttl)
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
