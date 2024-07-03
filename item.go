package cache

import (
	"time"
)

type Item struct {
	Object     any
	Expiration int64
}

// Returns true if the item has expired.
func (item Item) Expired() bool {
	return item.Expiration > 0 && time.Now().UnixNano() > item.Expiration
}

// Copies all unexpired items in the cache into a new map and returns it.
func (c *Cache) Items() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := make(map[string]Item, len(c.items))
	for k, v := range c.items {
		if !v.Expired() {
			m[k] = v
		}
	}
	return m
}

// Returns the number of items in the cache. This may include items that have
// expired, but have not yet been cleaned up.
func (c *Cache) ItemCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
