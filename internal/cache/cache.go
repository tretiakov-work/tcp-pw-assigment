package cache

import (
	"sync"
	"time"
)

type Cache struct {
	data     map[string]cacheItem
	mutex    sync.Mutex
	stopChan chan struct{}
}

type cacheItem struct {
	value    []byte
	expireAt time.Time
}

func New() *Cache {
	c := &Cache{
		data:     make(map[string]cacheItem),
		stopChan: make(chan struct{}),
	}
	go c.startCleanupRoutine()
	return c
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.data[key]
	if !ok {
		return nil, false
	}

	if item.expireAt.Before(time.Now()) {
		delete(c.data, key)
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheItem{
		value:    value,
		expireAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

func (c *Cache) startCleanupRoutine() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanupExpiredItems()
		case <-c.stopChan:
			return
		}
	}
}

func (c *Cache) cleanupExpiredItems() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if item.expireAt.Before(now) {
			delete(c.data, key)
		}
	}
}

func (c *Cache) Stop() {
	close(c.stopChan)
}
