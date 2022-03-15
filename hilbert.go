package hilbert

import (
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
}

type cache struct {
	items map[string]interface{}
	mu    sync.RWMutex
}

var _ Cache = &cache{}

func New() Cache {
	return &cache{
		items: make(map[string]interface{}),
	}
}

func (c *cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = value
	go c.remove(key, ttl)
}

func (c *cache) remove(key string, ttl time.Duration) {
	time.Sleep(ttl)

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.items[key]
	return v, ok
}
