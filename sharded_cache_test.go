package cache

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupShardedCache() *shardedCache {
	sc := newShardedCache(2, DefaultExpiration)
	return sc
}

func TestShardedCache_Set(t *testing.T) {
	t.Run("Set an item in the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", "value1", NoExpiration)

		val, found := sc.Get("key1")
		assert.True(t, found)
		assert.Equal(t, "value1", val)
	})
}

func TestShardedCache_Add(t *testing.T) {
	t.Run("Add an item to the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		err := sc.Add("key1", "value1", NoExpiration)

		assert.NoError(t, err)

		val, found := sc.Get("key1")
		assert.True(t, found)
		assert.Equal(t, "value1", val)
	})

	t.Run("Add an existing item to the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", "value1", NoExpiration)
		err := sc.Add("key1", "value2", NoExpiration)

		assert.Error(t, err)

		val, found := sc.Get("key1")
		assert.True(t, found)
		assert.Equal(t, "value1", val)
	})
}

func TestShardedCache_Replace(t *testing.T) {
	t.Run("Replace an existing item in the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", "value1", NoExpiration)
		err := sc.Replace("key1", "value2", NoExpiration)

		assert.NoError(t, err)

		val, found := sc.Get("key1")
		assert.True(t, found)
		assert.Equal(t, "value2", val)
	})

	t.Run("Replace a non-existing item in the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		err := sc.Replace("key1", "value1", NoExpiration)

		assert.Error(t, err)
	})
}

func TestShardedCache_Increment(t *testing.T) {
	t.Run("Increment an integer item in the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", 10, NoExpiration)
		err := sc.Increment("key1", 5)

		assert.NoError(t, err)

		val, found := sc.Get("key1")
		assert.True(t, found)
		assert.Equal(t, 15, val)
	})
}

func TestShardedCache_IncrementFloat(t *testing.T) {
	t.Run("Increment a float item in the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", 10.5, NoExpiration)
		err := sc.IncrementFloat("key1", 2.5)

		assert.NoError(t, err)

		val, found := sc.Get("key1")
		assert.True(t, found)
		assert.Equal(t, 13.0, val)
	})
}

func TestShardedCache_Decrement(t *testing.T) {
	t.Run("Decrement an integer item in the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", 10, NoExpiration)
		err := sc.Decrement("key1", 5)

		assert.NoError(t, err)

		val, found := sc.Get("key1")
		assert.True(t, found)
		assert.Equal(t, 5, val)
	})
}

func TestShardedCache_Delete(t *testing.T) {
	t.Run("Delete an item from the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", "value1", NoExpiration)
		sc.Delete("key1")

		_, found := sc.Get("key1")
		assert.False(t, found)
	})
}

func TestShardedCache_DeleteExpired(t *testing.T) {
	t.Run("Delete expired items from the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", "value1", 1*time.Millisecond)
		time.Sleep(2 * time.Millisecond)
		sc.DeleteExpired()

		_, found := sc.Get("key1")
		assert.False(t, found)
	})
}

func TestShardedCache_Items(t *testing.T) {
	t.Run("Get all items from the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", "value1", NoExpiration)
		sc.Set("key2", "value2", NoExpiration)

		items := sc.Items()

		assert.Equal(t, 2, len(items))

		item1, ok := items["key1"]
		assert.True(t, ok)
		assert.Equal(t, "value1", item1.Object)

		item2, ok := items["key2"]
		assert.True(t, ok)
		assert.Equal(t, "value2", item2.Object)
	})
}

func TestShardedCache_Flush(t *testing.T) {
	t.Run("Flush all items from the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", "value1", NoExpiration)
		sc.Set("key2", "value2", NoExpiration)
		sc.Flush()

		items := sc.Items()
		assert.Empty(t, items)
	})
}

func TestShardedCache_SetDefault(t *testing.T) {
	t.Run("SetDefault with default expiration", func(t *testing.T) {
		sc := setupShardedCache()
		sc.SetDefault("key1", "value1")

		val, found := sc.Get("key1")
		assert.True(t, found)
		assert.Equal(t, "value1", val)
	})
}

func TestShardedCache_GetWithExpiration(t *testing.T) {
	t.Run("GetWithExpiration non-existing item", func(t *testing.T) {
		sc := setupShardedCache()
		_, expTime, found := sc.GetWithExpiration("nonexistent")
		assert.False(t, found)
		assert.True(t, expTime.IsZero())
	})

	t.Run("GetWithExpiration existing item with no expiration", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", "value1", NoExpiration)

		val, expTime, found := sc.GetWithExpiration("key1")
		assert.True(t, found)
		assert.Equal(t, "value1", val)
		assert.True(t, expTime.IsZero())
	})

	t.Run("GetWithExpiration existing item with expiration", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", "value1", 5*time.Second)

		val, expTime, found := sc.GetWithExpiration("key1")
		assert.True(t, found)
		assert.Equal(t, "value1", val)
		assert.False(t, expTime.IsZero())
		assert.True(t, expTime.After(time.Now()))
	})
}

func TestShardedCache_DecrementFloat(t *testing.T) {
	t.Run("DecrementFloat a float item", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", 10.5, NoExpiration)
		err := sc.DecrementFloat("key1", 2.5)

		assert.NoError(t, err)

		val, found := sc.Get("key1")
		assert.True(t, found)
		assert.Equal(t, 8.0, val)
	})

	t.Run("DecrementFloat non-existing item", func(t *testing.T) {
		sc := setupShardedCache()
		err := sc.DecrementFloat("nonexistent", 1.0)
		assert.Error(t, err)
	})
}

func TestShardedCache_ItemCount(t *testing.T) {
	t.Run("ItemCount with no items", func(t *testing.T) {
		sc := setupShardedCache()
		count := sc.ItemCount()
		assert.Equal(t, 0, count)
	})

	t.Run("ItemCount with items", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", "value1", NoExpiration)
		sc.Set("key2", "value2", NoExpiration)
		sc.Set("key3", "value3", NoExpiration)

		count := sc.ItemCount()
		assert.Equal(t, 3, count)
	})
}

func TestShardedCache_OnEvicted(t *testing.T) {
	t.Run("OnEvicted callback is called", func(t *testing.T) {
		sc := setupShardedCache()
		evictedKeys := []string{}
		sc.OnEvicted(func(key string, value any) {
			evictedKeys = append(evictedKeys, key)
		})

		sc.Set("key1", "value1", NoExpiration)
		sc.Delete("key1")

		assert.Contains(t, evictedKeys, "key1")
	})
}

func TestShardedCache_SaveLoad(t *testing.T) {
	t.Run("Save and Load sharded cache", func(t *testing.T) {
		sc1 := newShardedCache(2, NoExpiration)
		sc1.Set("key1", "value1", NoExpiration)
		sc1.Set("key2", 123, NoExpiration)
		sc1.Set("key3", 45.6, NoExpiration)

		// Save to file
		err := sc1.SaveFile("test_sharded_cache.gob")
		assert.NoError(t, err)
		defer func() {
			// Clean up
			_ = os.Remove("test_sharded_cache.gob")
		}()

		// Create new cache and load
		sc2 := newShardedCache(2, NoExpiration)
		err = sc2.LoadFile("test_sharded_cache.gob")
		assert.NoError(t, err)

		// Verify data
		val1, found1 := sc2.Get("key1")
		assert.True(t, found1)
		assert.Equal(t, "value1", val1)

		val2, found2 := sc2.Get("key2")
		assert.True(t, found2)
		assert.Equal(t, 123, val2)

		val3, found3 := sc2.Get("key3")
		assert.True(t, found3)
		assert.Equal(t, 45.6, val3)
	})

	t.Run("Load with different shard count fails", func(t *testing.T) {
		sc1 := newShardedCache(2, NoExpiration)
		sc1.Set("key1", "value1", NoExpiration)

		err := sc1.SaveFile("test_sharded_mismatch.gob")
		assert.NoError(t, err)
		defer func() {
			_ = os.Remove("test_sharded_mismatch.gob")
		}()

		// Try to load into cache with different shard count
		sc2 := newShardedCache(4, NoExpiration)
		err = sc2.LoadFile("test_sharded_mismatch.gob")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "shard count mismatch")
	})

	t.Run("Load non-existent file fails", func(t *testing.T) {
		sc := newShardedCache(2, NoExpiration)
		err := sc.LoadFile("nonexistent_file.gob")
		assert.Error(t, err)
	})
}

// TestShardedCache_JanitorControlWithoutJanitor covers the nil-janitor paths:
// controls on a sharded cache created without cleanup interval are no-ops.
func TestShardedCache_JanitorControlWithoutJanitor(t *testing.T) {
	sc := NewSharded(NoExpiration, 0, 2)
	assert.NotPanics(t, func() {
		sc.PauseJanitor()
		sc.ResumeJanitor()
		sc.SetJanitorInterval(time.Second)
		sc.Close()
	})
}

// TestShardedCache_LoadShardMismatch verifies loading data saved with a
// different shard count fails instead of silently misplacing items.
func TestShardedCache_LoadShardMismatch(t *testing.T) {
	var buf bytes.Buffer
	saved := newShardedCache(2, NoExpiration)
	saved.Set("key", "value", NoExpiration)
	assert.NoError(t, saved.Save(&buf))

	loaded := newShardedCache(4, NoExpiration)
	err := loaded.Load(&buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "shard count mismatch")
}

// TestShardedCache_SaveFileError covers the file creation error path.
func TestShardedCache_SaveFileError(t *testing.T) {
	sc := newShardedCache(2, NoExpiration)
	err := sc.SaveFile("/invalid/path/to/file.gob")
	assert.Error(t, err)
}
