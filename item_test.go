package cache

import (
	"testing"
	"time"
)

func TestItem_Expired(t *testing.T) {
	t.Run("Not expired", func(t *testing.T) {
		item := Item{
			Object:     "test",
			Expiration: time.Now().Add(1 * time.Hour).UnixNano(),
		}
		if item.Expired() {
			t.Errorf("Expected item to not be expired")
		}
	})

	t.Run("Expired", func(t *testing.T) {
		item := Item{
			Object:     "test",
			Expiration: time.Now().Add(-1 * time.Hour).UnixNano(),
		}
		if !item.Expired() {
			t.Errorf("Expected item to be expired")
		}
	})

	t.Run("No expiration", func(t *testing.T) {
		item := Item{
			Object:     "test",
			Expiration: 0,
		}
		if item.Expired() {
			t.Errorf("Expected item to not be expired")
		}
	})
}

func TestCache_Items(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	cache.Set("key1", "value1", DefaultExpiration)
	cache.Set("key2", "value2", DefaultExpiration)
	cache.Set("key3", "value3", 1*time.Second)

	// Esperar a que caduque key3
	time.Sleep(2 * time.Second)

	items := cache.Items()
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	if items["key1"].Object != "value1" {
		t.Errorf("Expected value1 for key1, got %v", items["key1"].Object)
	}
	if items["key2"].Object != "value2" {
		t.Errorf("Expected value2 for key2, got %v", items["key2"].Object)
	}
	if _, found := items["key3"]; found {
		t.Errorf("Did not expect to find key3")
	}
}

func TestCache_ItemCount(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	cache.Set("key1", "value1", DefaultExpiration)
	cache.Set("key2", "value2", DefaultExpiration)
	cache.Set("key3", "value3", 1*time.Second)

	// Esperar a que caduque key3
	time.Sleep(2 * time.Second)

	count := cache.ItemCount()
	if count != 3 {
		t.Errorf("Expected item count to be 3, got %d", count)
	}

	// Limpiar ítems expirados
	cache.DeleteExpired()
	count = cache.ItemCount()
	if count != 2 {
		t.Errorf("Expected item count to be 2 after deleting expired items, got %d", count)
	}
}
