// Package localcache provides a simple in-memory cache with a TTL.
package localcache

import (
	"sync"
	"time"
)

// Cache is an interface for a simple in-memory cache with Get and Set methods.
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

// cache is an implementation of the Cache interface with a map, ttl and mutex lock.
type localCache struct {
	data map[string]interface{}
	ttl  time.Duration
	mu   sync.Mutex
}

// Get returns the value for a given key and a boolean indicating if the key exists.
func (c *localCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.data[key]
	return value, ok
}

// Set sets the value for a given key and deletes the key after the TTL.
func (c *localCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	f := func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		delete(c.data, key)
	}
	time.AfterFunc(c.ttl * time.Second, f)
}

// New returns a new cache instance with a given TTL.
func New(ttl time.Duration) Cache {
	return &localCache{
		data: make(map[string]interface{}),
		ttl:  ttl,
	}
}
