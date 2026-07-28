package cache

import (
	"runtime"
	"sync"
	"time"
)

const (
	// NoExpiration For use with functions that take an expiration time.
	NoExpiration time.Duration = -1
	// DefaultExpiration For use with functions that take an expiration time. Equivalent to
	// passing in the same expiration duration as was given to New() or
	// NewFrom() when the cache was created (e.g. 5 minutes.)
	DefaultExpiration time.Duration = 0
)

// Cache struct for cache control
type Cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
	onEvicted         func(string, any)
	janitor           *janitor
}

// Set Add an item to the cache, replacing any existing item. If the duration is 0
// (DefaultExpiration), the cache's default expiration time is used. If it is -1
// (NoExpiration), the item never expires.
func (c *Cache) Set(k string, x any, d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(k, x, d)
}

// set is the unexported version of Set that assumes the caller holds the lock.
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

// SetDefault Add an item to the cache, replacing any existing item, using the default
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
		return keyErrorf(ErrAlreadyExists, errItemAlreadyExistsFormat, k)
	}
	c.set(k, x, d)
	return nil
}

// Replace Set a new value for the cache key only if it already exists, and the existing
// item hasn't expired. Returns an error otherwise.
func (c *Cache) Replace(k string, x any, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.get(k)
	if !found {
		return keyErrorf(ErrNotFound, errItemDoesNotExistFormat, k)
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
	if !found || item.Expired() {
		return nil, time.Time{}, false
	}

	if item.Expiration > 0 {
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
	v, onEvicted, evicted := c.delete(k)
	c.mu.Unlock()
	if evicted {
		onEvicted(k, v)
	}
}

// delete removes the key and captures the eviction callback, both while
// synchronized, so the callback can be safely invoked after the lock is
// released even if OnEvicted(nil) runs concurrently.
func (c *Cache) delete(k string) (any, func(string, any), bool) {
	v, found := c.items[k]
	if !found {
		return nil, nil, false
	}
	delete(c.items, k)
	if c.onEvicted != nil {
		return v.Object, c.onEvicted, true
	}
	return nil, nil, false
}

type keyAndValue struct {
	key   string
	value any
}

// DeleteExpired Delete all expired items from the cache.
func (c *Cache) DeleteExpired() {
	var evictedItems []keyAndValue
	c.mu.Lock()
	onEvicted := c.onEvicted
	for k, v := range c.items {
		if v.Expired() {
			delete(c.items, k)
			if onEvicted != nil {
				evictedItems = append(evictedItems, keyAndValue{k, v.Object})
			}
		}
	}
	c.mu.Unlock()
	for _, v := range evictedItems {
		onEvicted(v.key, v.value)
	}
}

// OnEvicted Sets an (optional) function that is called with the key and value when an
// item is evicted from the cache. (Including when it is deleted manually, but
// not when it is overwritten.) Set to nil to disable.
func (c *Cache) OnEvicted(f func(string, any)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onEvicted = f
}

// Flush Delete all items from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]Item)
}

// PauseJanitor temporarily pauses the automatic cleanup of expired items.
// The janitor will stop deleting expired items until ResumeJanitor is called.
// Safe to call multiple times. This method has no effect if the janitor is not running.
func (c *Cache) PauseJanitor() {
	c.mu.RLock()
	j := c.janitor
	c.mu.RUnlock()
	if j != nil {
		j.Pause()
	}
}

// ResumeJanitor resumes the automatic cleanup of expired items after it was paused.
// Safe to call multiple times. This method has no effect if the janitor is not running or not paused.
func (c *Cache) ResumeJanitor() {
	c.mu.RLock()
	j := c.janitor
	c.mu.RUnlock()
	if j != nil {
		j.Resume()
	}
}

// SetJanitorInterval dynamically updates the janitor's cleanup interval.
// The new interval will take effect immediately. Non-positive intervals are
// ignored. This method has no effect if the janitor is not running.
func (c *Cache) SetJanitorInterval(d time.Duration) {
	c.mu.RLock()
	j := c.janitor
	c.mu.RUnlock()
	if j != nil {
		j.setInterval(d)
	}
}

// Close stops the janitor goroutine and releases resources.
// After calling Close, the cache can still be used but expired items
// will no longer be cleaned up automatically. Safe to call multiple times.
func (c *Cache) Close() {
	c.mu.Lock()
	j := c.janitor
	c.janitor = nil
	c.mu.Unlock()
	if j != nil {
		// Remove the GC finalizer: the janitor is already stopped, and
		// stopJanitor would dereference the nil'ed janitor field.
		runtime.SetFinalizer(c, nil)
		j.stop <- struct{}{}
	}
}
