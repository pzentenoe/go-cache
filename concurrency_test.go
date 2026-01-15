package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCacheConcurrentReadWrite tests concurrent read and write operations
func TestCacheConcurrentReadWrite(t *testing.T) {
	t.Run("Concurrent reads and writes to same cache", func(t *testing.T) {
		c := New(NoExpiration, 0)
		numGoroutines := 100
		numOperations := 1000

		var wg sync.WaitGroup
		wg.Add(numGoroutines * 2) // readers and writers

		// Writers
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					key := fmt.Sprintf("key-%d-%d", id, j)
					c.Set(key, j, NoExpiration)
				}
			}(i)
		}

		// Readers
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					key := fmt.Sprintf("key-%d-%d", id, j)
					c.Get(key)
				}
			}(i)
		}

		wg.Wait()

		// Verify cache is in consistent state
		assert.NotNil(t, c)
		count := c.ItemCount()
		assert.GreaterOrEqual(t, count, 0)
	})

	t.Run("Concurrent writes to same key", func(t *testing.T) {
		c := New(NoExpiration, 0)
		numGoroutines := 100
		key := "shared-key"

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(val int) {
				defer wg.Done()
				c.Set(key, val, NoExpiration)
			}(i)
		}

		wg.Wait()

		// Should have exactly one value (last write wins)
		val, found := c.Get(key)
		assert.True(t, found)
		assert.NotNil(t, val)
	})

	t.Run("Concurrent increment operations", func(t *testing.T) {
		c := New(NoExpiration, 0)
		key := "counter"
		c.Set(key, int64(0), NoExpiration)

		numGoroutines := 100
		incrementsPerGoroutine := 100
		expectedTotal := int64(numGoroutines * incrementsPerGoroutine)

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < incrementsPerGoroutine; j++ {
					c.Increment(key, 1)
				}
			}()
		}

		wg.Wait()

		val, found := c.Get(key)
		assert.True(t, found)
		assert.Equal(t, expectedTotal, val.(int64))
	})
}

// TestCacheDeleteDuringRead tests race conditions between Delete and Get
func TestCacheDeleteDuringRead(t *testing.T) {
	t.Run("Delete during concurrent reads", func(t *testing.T) {
		c := New(NoExpiration, 0)

		// Populate cache
		for i := 0; i < 1000; i++ {
			c.Set(fmt.Sprintf("key-%d", i), i, NoExpiration)
		}

		numReaders := 50
		numDeleters := 10

		var wg sync.WaitGroup
		wg.Add(numReaders + numDeleters)

		// Readers
		for i := 0; i < numReaders; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < 1000; j++ {
					c.Get(fmt.Sprintf("key-%d", j))
				}
			}()
		}

		// Deleters
		for i := 0; i < numDeleters; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < 1000; j++ {
					c.Delete(fmt.Sprintf("key-%d", j))
				}
			}()
		}

		wg.Wait()

		// Cache should be in consistent state
		assert.NotNil(t, c)
	})
}

// TestShardedCacheConcurrency tests sharded cache under concurrent load
func TestShardedCacheConcurrency(t *testing.T) {
	t.Run("Concurrent operations across shards", func(t *testing.T) {
		sc := NewSharded(NoExpiration, 0, 16)
		numGoroutines := 100
		numOperations := 1000

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					key := fmt.Sprintf("key-%d-%d", id, j)
					sc.Set(key, j, NoExpiration)
					sc.Get(key)
				}
			}(i)
		}

		wg.Wait()

		count := sc.ItemCount()
		assert.GreaterOrEqual(t, count, 0)
	})

	t.Run("Sharded cache under high load", func(t *testing.T) {
		sc := NewSharded(NoExpiration, 0, 32)
		numGoroutines := 200
		operationsPerGoroutine := 500

		var wg sync.WaitGroup
		wg.Add(numGoroutines * 3) // Set, Get, Delete operations

		// Writers
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < operationsPerGoroutine; j++ {
					key := fmt.Sprintf("key-%d-%d", id, j)
					sc.Set(key, j, NoExpiration)
				}
			}(i)
		}

		// Readers
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < operationsPerGoroutine; j++ {
					key := fmt.Sprintf("key-%d-%d", id, j)
					sc.Get(key)
				}
			}(i)
		}

		// Deleters
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < operationsPerGoroutine; j++ {
					if j%2 == 0 { // Delete every other key
						key := fmt.Sprintf("key-%d-%d", id, j)
						sc.Delete(key)
					}
				}
			}(i)
		}

		wg.Wait()

		// Verify consistent state
		count := sc.ItemCount()
		assert.GreaterOrEqual(t, count, 0)
	})
}

// TestJanitorUnderLoad tests janitor behavior under concurrent operations
func TestJanitorUnderLoad(t *testing.T) {
	t.Run("Janitor with concurrent expiring items", func(t *testing.T) {
		// Use short expiration and cleanup interval
		c := New(50*time.Millisecond, 25*time.Millisecond)
		numGoroutines := 50
		numOperations := 100

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					key := fmt.Sprintf("key-%d-%d", id, j)
					// Some with expiration, some without
					if j%2 == 0 {
						c.Set(key, j, 50*time.Millisecond)
					} else {
						c.Set(key, j, NoExpiration)
					}
				}
			}(i)
		}

		wg.Wait()

		// Wait for janitor to run multiple times
		time.Sleep(200 * time.Millisecond)

		// Verify cache is in consistent state
		count := c.ItemCount()
		assert.GreaterOrEqual(t, count, 0)

		// Most expired items should be cleaned up
		// We can't guarantee exact count due to timing, but it should be less than total
		assert.Less(t, count, numGoroutines*numOperations)
	})

	t.Run("Janitor with concurrent DeleteExpired calls", func(t *testing.T) {
		c := New(10*time.Millisecond, 0) // No automatic janitor

		// Add items with short expiration
		for i := 0; i < 1000; i++ {
			c.Set(fmt.Sprintf("key-%d", i), i, 10*time.Millisecond)
		}

		// Wait for items to expire
		time.Sleep(50 * time.Millisecond)

		// Multiple goroutines calling DeleteExpired simultaneously
		numGoroutines := 10
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				c.DeleteExpired()
			}()
		}

		wg.Wait()

		// All expired items should be deleted
		count := c.ItemCount()
		assert.Equal(t, 0, count)
	})
}

// TestConcurrentReplace tests concurrent Replace operations
func TestConcurrentReplace(t *testing.T) {
	t.Run("Concurrent Replace operations", func(t *testing.T) {
		c := New(NoExpiration, 0)
		key := "replace-key"
		c.Set(key, 0, NoExpiration)

		numGoroutines := 50
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(val int) {
				defer wg.Done()
				c.Replace(key, val, NoExpiration)
			}(i)
		}

		wg.Wait()

		val, found := c.Get(key)
		assert.True(t, found)
		assert.NotNil(t, val)
	})
}

// TestConcurrentAdd tests concurrent Add operations
func TestConcurrentAdd(t *testing.T) {
	t.Run("Concurrent Add operations on different keys", func(t *testing.T) {
		c := New(NoExpiration, 0)
		numGoroutines := 100

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		successCount := int32(0)
		var mu sync.Mutex

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				key := fmt.Sprintf("add-key-%d", id)
				err := c.Add(key, id, NoExpiration)
				if err == nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}(i)
		}

		wg.Wait()

		// All adds should succeed (different keys)
		assert.Equal(t, int32(numGoroutines), successCount)
		assert.Equal(t, numGoroutines, c.ItemCount())
	})

	t.Run("Concurrent Add operations on same key", func(t *testing.T) {
		c := New(NoExpiration, 0)
		key := "shared-add-key"
		numGoroutines := 100

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		successCount := int32(0)
		var mu sync.Mutex

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				err := c.Add(key, id, NoExpiration)
				if err == nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}(i)
		}

		wg.Wait()

		// Only one Add should succeed
		assert.Equal(t, int32(1), successCount)
		assert.Equal(t, 1, c.ItemCount())
	})
}

// TestConcurrentFlush tests Flush under concurrent operations
func TestConcurrentFlush(t *testing.T) {
	t.Run("Flush during concurrent operations", func(t *testing.T) {
		c := New(NoExpiration, 0)

		// Populate cache
		for i := 0; i < 1000; i++ {
			c.Set(fmt.Sprintf("key-%d", i), i, NoExpiration)
		}

		numWriters := 50
		numReaders := 50
		numFlushers := 5

		var wg sync.WaitGroup
		wg.Add(numWriters + numReaders + numFlushers)

		// Writers
		for i := 0; i < numWriters; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					c.Set(fmt.Sprintf("new-key-%d-%d", id, j), j, NoExpiration)
				}
			}(i)
		}

		// Readers
		for i := 0; i < numReaders; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					c.Get(fmt.Sprintf("key-%d", j))
				}
			}()
		}

		// Flushers
		for i := 0; i < numFlushers; i++ {
			go func() {
				defer wg.Done()
				time.Sleep(10 * time.Millisecond)
				c.Flush()
			}()
		}

		wg.Wait()

		// Cache should be in consistent state
		assert.NotNil(t, c)
	})
}

// TestConcurrentOnEvicted tests OnEvicted callback under concurrent operations
func TestConcurrentOnEvicted(t *testing.T) {
	t.Run("OnEvicted callback with concurrent deletions", func(t *testing.T) {
		c := New(NoExpiration, 0)

		var evictedCount int32
		var mu sync.Mutex

		c.OnEvicted(func(key string, value any) {
			mu.Lock()
			evictedCount++
			mu.Unlock()
		})

		// Add items
		numItems := 1000
		for i := 0; i < numItems; i++ {
			c.Set(fmt.Sprintf("key-%d", i), i, NoExpiration)
		}

		// Delete concurrently
		numGoroutines := 50
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < numItems/numGoroutines; j++ {
					idx := id*(numItems/numGoroutines) + j
					c.Delete(fmt.Sprintf("key-%d", idx))
				}
			}(i)
		}

		wg.Wait()

		// All items should have been evicted
		assert.Equal(t, int32(numItems), evictedCount)
		assert.Equal(t, 0, c.ItemCount())
	})
}
