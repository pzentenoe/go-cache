package cache

import (
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

		foundKey1 := false
		foundKey2 := false

		for _, shardItems := range items {
			if item, ok := shardItems["key1"]; ok {
				assert.Equal(t, "value1", item.Object)
				foundKey1 = true
			}
			if item, ok := shardItems["key2"]; ok {
				assert.Equal(t, "value2", item.Object)
				foundKey2 = true
			}
		}

		assert.True(t, foundKey1)
		assert.True(t, foundKey2)
	})
}

func TestShardedCache_Flush(t *testing.T) {
	t.Run("Flush all items from the sharded cache", func(t *testing.T) {
		sc := setupShardedCache()
		sc.Set("key1", "value1", NoExpiration)
		sc.Set("key2", "value2", NoExpiration)
		sc.Flush()

		items := sc.Items()

		for _, shardItems := range items {
			assert.Empty(t, shardItems)
		}
	})
}
