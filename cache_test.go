package cache

import (
	"testing"
	"time"
)

func TestCache_Set(t *testing.T) {
	cache := New(DefaultExpiration, 0)

	t.Run("Set item with default expiration", func(t *testing.T) {
		cache.Set("key1", "value1", DefaultExpiration)
		if v, found := cache.Get("key1"); !found || v != "value1" {
			t.Errorf("Expected to find key1 with value 'value1', got %v", v)
		}
	})

	t.Run("Set item with no expiration", func(t *testing.T) {
		cache.Set("key2", "value2", NoExpiration)
		if v, found := cache.Get("key2"); !found || v != "value2" {
			t.Errorf("Expected to find key2 with value 'value2', got %v", v)
		}
	})

	t.Run("Set item with custom expiration", func(t *testing.T) {
		cache.Set("key3", "value3", 100*time.Millisecond)
		if v, found := cache.Get("key3"); !found || v != "value3" {
			t.Errorf("Expected to find key3 with value 'value3', got %v", v)
		}

		// Wait for key3 to expire
		time.Sleep(200 * time.Millisecond)
		if _, found := cache.Get("key3"); found {
			t.Error("Expected key3 to be expired")
		}
	})

	t.Run("Overwrite existing item", func(t *testing.T) {
		cache.Set("key1", "newValue1", DefaultExpiration)
		if v, found := cache.Get("key1"); !found || v != "newValue1" {
			t.Errorf("Expected to find key1 with value 'newValue1', got %v", v)
		}
	})
}

func TestCache_Add(t *testing.T) {
	cache := New(DefaultExpiration, 0)

	t.Run("Add item that doesn't exist", func(t *testing.T) {
		err := cache.Add("key1", "value1", DefaultExpiration)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if v, found := cache.Get("key1"); !found || v != "value1" {
			t.Errorf("Expected to find key1 with value 'value1', got %v", v)
		}
	})

	t.Run("Add item that already exists", func(t *testing.T) {
		err := cache.Add("key1", "newValue1", DefaultExpiration)
		if err == nil {
			t.Error("Expected an error, got nil")
		}
		if v, found := cache.Get("key1"); !found || v != "value1" {
			t.Errorf("Expected to find key1 with value 'value1', got %v", v)
		}
	})

	t.Run("Add item with no expiration", func(t *testing.T) {
		err := cache.Add("key2", "value2", NoExpiration)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if v, found := cache.Get("key2"); !found || v != "value2" {
			t.Errorf("Expected to find key2 with value 'value2', got %v", v)
		}
	})

	t.Run("Add item with custom expiration", func(t *testing.T) {
		err := cache.Add("key3", "value3", 100*time.Millisecond)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if v, found := cache.Get("key3"); !found || v != "value3" {
			t.Errorf("Expected to find key3 with value 'value3', got %v", v)
		}

		// Wait for key3 to expire
		time.Sleep(200 * time.Millisecond)
		if _, found := cache.Get("key3"); found {
			t.Error("Expected key3 to be expired")
		}
	})
}

func TestCache_Replace(t *testing.T) {
	cache := New(DefaultExpiration, 0)

	t.Run("Replace non-existing item", func(t *testing.T) {
		err := cache.Replace("key1", "value1", DefaultExpiration)
		if err == nil {
			t.Error("Expected error when replacing non-existing item, got nil")
		}
	})

	t.Run("Replace existing item", func(t *testing.T) {
		cache.Set("key2", "value2", DefaultExpiration)
		err := cache.Replace("key2", "newValue2", DefaultExpiration)
		if err != nil {
			t.Errorf("Expected no error when replacing existing item, got %v", err)
		}
		if v, found := cache.Get("key2"); !found || v != "newValue2" {
			t.Errorf("Expected to find key2 with value 'newValue2', got %v", v)
		}
	})

	t.Run("Replace expired item", func(t *testing.T) {
		cache.Set("key3", "value3", 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond)
		err := cache.Replace("key3", "newValue3", DefaultExpiration)
		if err == nil {
			t.Error("Expected error when replacing expired item, got nil")
		}
	})
}

func TestCache_Get(t *testing.T) {
	cache := New(DefaultExpiration, 0)

	t.Run("Get non-existing item", func(t *testing.T) {
		_, found := cache.Get("key1")
		if found {
			t.Error("Expected not to find a non-existing item, but found one")
		}
	})

	t.Run("Get existing item", func(t *testing.T) {
		cache.Set("key2", "value2", DefaultExpiration)
		v, found := cache.Get("key2")
		if !found || v != "value2" {
			t.Errorf("Expected to find key2 with value 'value2', got %v", v)
		}
	})

	t.Run("Get expired item", func(t *testing.T) {
		cache.Set("key3", "value3", 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond)
		_, found := cache.Get("key3")
		if found {
			t.Error("Expected not to find an expired item, but found one")
		}
	})
}

func TestCache_GetWithExpiration(t *testing.T) {
	cache := New(DefaultExpiration, 0)

	t.Run("GetWithExpiration non-existing item", func(t *testing.T) {
		_, _, found := cache.GetWithExpiration("key1")
		if found {
			t.Error("Expected not to find a non-existing item, but found one")
		}
	})

	t.Run("GetWithExpiration existing item", func(t *testing.T) {
		cache.Set("key2", "value2", DefaultExpiration)
		v, exp, found := cache.GetWithExpiration("key2")
		if !found || v != "value2" || !exp.IsZero() {
			t.Errorf("Expected to find key2 with value 'value2' and zero expiration, got %v, %v", v, exp)
		}
	})

	t.Run("GetWithExpiration expired item", func(t *testing.T) {
		cache.Set("key3", "value3", 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond)
		_, _, found := cache.GetWithExpiration("key3")
		if found {
			t.Error("Expected not to find an expired item, but found one")
		}
	})

	t.Run("GetWithExpiration item with expiration", func(t *testing.T) {
		cache.Set("key4", "value4", 5*time.Second)
		v, exp, found := cache.GetWithExpiration("key4")
		if !found || v != "value4" || exp.IsZero() {
			t.Errorf("Expected to find key4 with value 'value4' and non-zero expiration, got %v, %v", v, exp)
		}
	})
}

func TestCache_Delete(t *testing.T) {
	cache := New(DefaultExpiration, 0)

	t.Run("Delete non-existing item", func(t *testing.T) {
		cache.Delete("key1")
		_, found := cache.Get("key1")
		if found {
			t.Error("Expected not to find a non-existing item, but found one")
		}
	})

	t.Run("Delete existing item", func(t *testing.T) {
		cache.Set("key2", "value2", DefaultExpiration)
		cache.Delete("key2")
		_, found := cache.Get("key2")
		if found {
			t.Error("Expected not to find the deleted item, but found one")
		}
	})
}

func TestCache_DeleteExpired(t *testing.T) {
	cache := New(DefaultExpiration, 0)

	t.Run("DeleteExpired with no expired items", func(t *testing.T) {
		cache.Set("key1", "value1", DefaultExpiration)
		cache.DeleteExpired()
		_, found := cache.Get("key1")
		if !found {
			t.Error("Expected to find key1, but did not find it")
		}
	})

	t.Run("DeleteExpired with expired items", func(t *testing.T) {
		cache.Set("key2", "value2", 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond)
		cache.DeleteExpired()
		_, found := cache.Get("key2")
		if found {
			t.Error("Expected not to find the expired item, but found one")
		}
	})
}

func TestCache_OnEvicted(t *testing.T) {
	cache := New(DefaultExpiration, 0)

	t.Run("Set OnEvicted handler", func(t *testing.T) {
		var evictedKey string
		var evictedValue any
		cache.OnEvicted(func(k string, v any) {
			evictedKey = k
			evictedValue = v
		})

		cache.Set("key1", "value1", DefaultExpiration)
		cache.Delete("key1")

		if evictedKey != "key1" || evictedValue != "value1" {
			t.Errorf("Expected evicted key to be 'key1' and value 'value1', got %v and %v", evictedKey, evictedValue)
		}
	})
}

func TestCache_Flush(t *testing.T) {
	cache := New(DefaultExpiration, 0)

	t.Run("Flush cache", func(t *testing.T) {
		cache.Set("key1", "value1", DefaultExpiration)
		cache.Set("key2", "value2", DefaultExpiration)
		cache.Flush()
		if n := cache.ItemCount(); n != 0 {
			t.Errorf("Expected item count to be 0 after flush, got %d", n)
		}
	})
}
