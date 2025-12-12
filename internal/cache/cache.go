package cache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache: make(map[string]cacheEntry),
		mu:    sync.RWMutex{},
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			c.ReapLoop(interval)
		}
	}()

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if key == "" {
		return
	}
	c.cache[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if key == "" {
		return []byte(""), false
	}
	ce, ok := c.cache[key]
	return ce.val, ok
}

func (c *Cache) ReapLoop(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for key, val := range c.cache {
		if val.createdAt.Add(interval).Before(now) {
			delete(c.cache, key)
		}
	}
}
