package cache

import "sync"

// Cache
type cacheItem[T any] struct {
	value     T
	timestamp int64 // unix seconds
}

type Cache[T any] struct {
	data map[string]cacheItem[T]
	mu   sync.RWMutex
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{data: make(map[string]cacheItem[T])}
}

func (c *Cache[T]) Set(key string, value T, now int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheItem[T]{value: value, timestamp: now}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.data[key]
	if !ok {
		 var zero T
		 return zero, false
	}
	return item.value, true
}

// GetWithTTL returns cached value if not expired, otherwise calls fetch and updates cache
func (c *Cache[T]) GetWithTTL(key string, ttlSeconds int64, now int64, fetch func() T) T {
       c.mu.RLock()
       item, ok := c.data[key]
       c.mu.RUnlock()
       if ok && now-item.timestamp <= ttlSeconds {
	       return item.value
       }
       // fetch fresh data (simulate DB)
       val := fetch()
       c.Set(key, val, now)
       return val
}
