package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	data map[string]cacheEntry
	mu   *sync.Mutex
}

func (c Cache) Add(key string, val []byte) {
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mu.Lock()
	c.data[key] = entry
	c.mu.Unlock()
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elt, ok := c.data[key]
	if ok {
		return elt.val, ok
	}
	return nil, ok
}

func (c Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		for key, entry := range c.data {
			if entry.createdAt.Add(interval).Before(now) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) Cache {
	var mu sync.Mutex
	cache := Cache{
		data: map[string]cacheEntry{},
		mu:   &mu,
	}
	go cache.reapLoop(interval)
	return cache
}
