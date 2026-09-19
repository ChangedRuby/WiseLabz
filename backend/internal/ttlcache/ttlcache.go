// Package ttlcache provides a tiny in-memory cache whose entries expire after
// a fixed TTL, for short-lived read-mostly API responses.
package ttlcache

import (
	"sync"
	"time"
)

type entry[V any] struct {
	value   V
	expires time.Time
}

// Cache is a concurrency-safe TTL cache.
type Cache[V any] struct {
	mu  sync.Mutex
	ttl time.Duration
	m   map[string]entry[V]
}

// New returns a cache whose entries live for ttl.
func New[V any](ttl time.Duration) *Cache[V] {
	return &Cache[V]{ttl: ttl, m: make(map[string]entry[V])}
}

// Get returns the cached value for key if present and unexpired.
func (c *Cache[V]) Get(key string) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.m[key]
	if !ok || time.Now().After(e.expires) {
		delete(c.m, key)
		var zero V
		return zero, false
	}
	return e.value, true
}

// Set stores value under key, sweeping expired entries so the map stays bounded.
func (c *Cache[V]) Set(key string, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, e := range c.m {
		if now.After(e.expires) {
			delete(c.m, k)
		}
	}
	c.m[key] = entry[V]{value: value, expires: now.Add(c.ttl)}
}
