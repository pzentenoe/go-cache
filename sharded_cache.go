package cache

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"
)

func (sc *shardedCache) bucket(k string) *Cache {
	return sc.cs[djb33(sc.seed, k)%sc.m]
}

func (sc *shardedCache) Set(k string, x any, d time.Duration) {
	sc.bucket(k).Set(k, x, d)
}

func (sc *shardedCache) SetDefault(k string, x any) {
	sc.bucket(k).SetDefault(k, x)
}

func (sc *shardedCache) Add(k string, x any, d time.Duration) error {
	return sc.bucket(k).Add(k, x, d)
}

func (sc *shardedCache) Replace(k string, x any, d time.Duration) error {
	return sc.bucket(k).Replace(k, x, d)
}

func (sc *shardedCache) Get(k string) (any, bool) {
	return sc.bucket(k).Get(k)
}

func (sc *shardedCache) GetWithExpiration(k string) (any, time.Time, bool) {
	return sc.bucket(k).GetWithExpiration(k)
}

func (sc *shardedCache) Increment(k string, n int64) error {
	return sc.bucket(k).Increment(k, n)
}

func (sc *shardedCache) IncrementFloat(k string, n float64) error {
	return sc.bucket(k).IncrementFloat(k, n)
}

func (sc *shardedCache) Decrement(k string, n int64) error {
	return sc.bucket(k).Decrement(k, n)
}

func (sc *shardedCache) DecrementFloat(k string, n float64) error {
	return sc.bucket(k).DecrementFloat(k, n)
}

func (sc *shardedCache) Delete(k string) {
	sc.bucket(k).Delete(k)
}

func (sc *shardedCache) DeleteExpired() {
	for _, v := range sc.cs {
		v.DeleteExpired()
	}
}

func (sc *shardedCache) Items() map[string]Item {
	// ItemCount is an upper bound: it includes expired items not yet cleaned up.
	res := make(map[string]Item, sc.ItemCount())
	for _, v := range sc.cs {
		for k, item := range v.Items() {
			res[k] = item
		}
	}
	return res
}

func (sc *shardedCache) ItemCount() int {
	count := 0
	for _, c := range sc.cs {
		count += c.ItemCount()
	}
	return count
}

func (sc *shardedCache) OnEvicted(f func(string, any)) {
	for _, c := range sc.cs {
		c.OnEvicted(f)
	}
}

func (sc *shardedCache) Flush() {
	for _, v := range sc.cs {
		v.Flush()
	}
}

// PauseJanitor temporarily pauses the automatic cleanup of expired items.
// Safe to call multiple times; subsequent calls are no-ops.
func (sc *shardedCache) PauseJanitor() {
	sc.mu.Lock()
	j := sc.janitor
	sc.mu.Unlock()
	if j != nil {
		j.Pause()
	}
}

// ResumeJanitor resumes the automatic cleanup of expired items after it was paused.
// Safe to call multiple times; subsequent calls are no-ops.
func (sc *shardedCache) ResumeJanitor() {
	sc.mu.Lock()
	j := sc.janitor
	sc.mu.Unlock()
	if j != nil {
		j.Resume()
	}
}

// SetJanitorInterval dynamically updates the janitor's cleanup interval.
// Non-positive intervals are ignored.
func (sc *shardedCache) SetJanitorInterval(d time.Duration) {
	sc.mu.Lock()
	j := sc.janitor
	sc.mu.Unlock()
	if j != nil {
		j.setInterval(d)
	}
}

// Close stops the janitor goroutine and releases resources.
// After calling Close, the cache can still be used but expired items
// will no longer be cleaned up automatically. Safe to call multiple times.
func (sc *shardedCache) Close() {
	sc.mu.Lock()
	j := sc.janitor
	sc.janitor = nil
	sc.mu.Unlock()
	if j != nil {
		j.stop <- struct{}{}
	}
}

// Save writes the sharded cache's items (using Gob) to an io.Writer.
// All shards are serialized sequentially.
func (sc *shardedCache) Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = errors.New(errGobRegistration)
		}
	}()

	// Save the number of shards first
	if err := enc.Encode(len(sc.cs)); err != nil {
		return err
	}

	// Save each shard's items; value types are registered once across shards.
	registered := make(map[reflect.Type]struct{})
	for _, c := range sc.cs {
		c.mu.RLock()
		registerGobTypes(c.items, registered)
		if err := enc.Encode(&c.items); err != nil {
			c.mu.RUnlock()
			return err
		}
		c.mu.RUnlock()
	}
	return nil
}

// SaveFile saves the sharded cache's items to the given filename, creating
// the file if it doesn't exist, and overwriting it if it does.
func (sc *shardedCache) SaveFile(fname string) error {
	fp, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer fp.Close()
	return sc.Save(fp)
}

// Load adds (Gob-serialized) cache items from an io.Reader, excluding any
// items with keys that already exist (and haven't expired) in the current cache.
func (sc *shardedCache) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)

	// Read number of shards
	var numShards int
	if err := dec.Decode(&numShards); err != nil {
		return err
	}

	if numShards != len(sc.cs) {
		return fmt.Errorf("shard count mismatch: saved with %d shards, current cache has %d", numShards, len(sc.cs))
	}

	// Load items for each shard
	for i := 0; i < numShards; i++ {
		items := map[string]Item{}
		if err := dec.Decode(&items); err != nil {
			return err
		}

		// Re-hash each item to the correct shard
		for k, v := range items {
			c := sc.bucket(k)
			c.mu.Lock()
			if ov, found := c.items[k]; !found || ov.Expired() {
				c.items[k] = v
			}
			c.mu.Unlock()
		}
	}
	return nil
}

// LoadFile loads and adds cache items from the given filename, excluding any
// items with keys that already exist in the current cache.
func (sc *shardedCache) LoadFile(fname string) error {
	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fp.Close()
	return sc.Load(fp)
}
