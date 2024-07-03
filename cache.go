package cache

import (
	"fmt"
	"sync"
	"time"
)

const (
	// For use with functions that take an expiration time.
	NoExpiration time.Duration = -1
	// For use with functions that take an expiration time. Equivalent to
	// passing in the same expiration duration as was given to New() or
	// NewFrom() when the cache was created (e.g. 5 minutes.)
	DefaultExpiration time.Duration = 0
)

type Cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
	onEvicted         func(string, any)
	janitor           *janitor
}

// Add an item to the cache, replacing any existing item. If the duration is 0
// (DefaultExpiration), the cache's default expiration time is used. If it is -1
// (NoExpiration), the item never expires.
func (c *Cache) Set(k string, x any, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
}

func (c *Cache) set(k string, x any, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
}

// Add an item to the cache, replacing any existing item, using the default
// expiration.
func (c *Cache) SetDefault(k string, x any) {
	c.Set(k, x, DefaultExpiration)
}

// Add an item to the cache only if an item doesn't already exist for the given
// key, or if the existing item has expired. Returns an error otherwise.
func (c *Cache) Add(k string, x any, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.get(k)
	if found {
		return fmt.Errorf("Item %s already exists", k)
	}
	c.set(k, x, d)
	return nil
}

// Set a new value for the cache key only if it already exists, and the existing
// item hasn't expired. Returns an error otherwise.
func (c *Cache) Replace(k string, x any, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.get(k)
	if !found {
		return fmt.Errorf("Item %s doesn't exist", k)
	}
	c.set(k, x, d)
	return nil
}

// Get an item from the cache. Returns the item or nil, and a bool indicating
// whether the key was found.
func (c *Cache) Get(k string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[k]
	if !found || item.Expired() {
		return nil, false
	}
	return item.Object, true
}

// GetWithExpiration returns an item and its expiration time from the cache.
// It returns the item or nil, the expiration time if one is set (if the item
// never expires a zero value for time.Time is returned), and a bool indicating
// whether the key was found.
func (c *Cache) GetWithExpiration(k string) (any, time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[k]
	if !found {
		return nil, time.Time{}, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, time.Time{}, false
		}
		return item.Object, time.Unix(0, item.Expiration), true
	}

	// If expiration <= 0 (i.e. no expiration time set), return the item and a zeroed time.Time
	return item.Object, time.Time{}, true
}

func (c *Cache) get(k string) (any, bool) {
	item, found := c.items[k]
	if !found || item.Expired() {
		return nil, false
	}
	return item.Object, true
}

// Delete an item from the cache. Does nothing if the key is not in the cache.
func (c *Cache) Delete(k string) {
	c.mu.Lock()
	v, evicted := c.delete(k)
	c.mu.Unlock()
	if evicted {
		c.onEvicted(k, v)
	}
}

func (c *Cache) delete(k string) (any, bool) {
	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Object, true
		}
	}
	delete(c.items, k)
	return nil, false
}

type keyAndValue struct {
	key   string
	value any
}

// Delete all expired items from the cache.
func (c *Cache) DeleteExpired() {
	var evictedItems []keyAndValue
	c.mu.Lock()
	for k, v := range c.items {
		if v.Expired() {
			ov, evicted := c.delete(k)
			if evicted {
				evictedItems = append(evictedItems, keyAndValue{k, ov})
			}
		}
	}
	c.mu.Unlock()
	for _, v := range evictedItems {
		c.onEvicted(v.key, v.value)
	}
}

// Sets an (optional) function that is called with the key and value when an
// item is evicted from the cache. (Including when it is deleted manually, but
// not when it is overwritten.) Set to nil to disable.
func (c *Cache) OnEvicted(f func(string, any)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onEvicted = f
}

// Delete all items from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]Item)
}
