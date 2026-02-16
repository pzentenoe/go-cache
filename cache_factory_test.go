package cache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	t.Run("New cache with default expiration", func(t *testing.T) {
		cache := newCache(DefaultExpiration, make(map[string]Item))
		if cache.defaultExpiration != DefaultExpiration {
			t.Errorf("Expected default expiration to be %v, got %v", DefaultExpiration, cache.defaultExpiration)
		}
		if cache.items == nil {
			t.Error("Expected items map to be initialized")
		}
	})

	t.Run("New cache with NoExpiration", func(t *testing.T) {
		cache := newCache(NoExpiration, make(map[string]Item))
		if cache.defaultExpiration != NoExpiration {
			t.Errorf("Expected default expiration to be %v, got %v", NoExpiration, cache.defaultExpiration)
		}
		if cache.items == nil {
			t.Error("Expected items map to be initialized")
		}
	})
}

func TestNewCacheWithJanitor(t *testing.T) {
	t.Run("New cache with janitor", func(t *testing.T) {
		cache := newCacheWithJanitor(DefaultExpiration, 1*time.Minute, make(map[string]Item))
		if cache.defaultExpiration != DefaultExpiration {
			t.Errorf("Expected default expiration to be %v, got %v", DefaultExpiration, cache.defaultExpiration)
		}
		if cache.items == nil {
			t.Error("Expected items map to be initialized")
		}
		if cache.janitor == nil {
			t.Error("Expected janitor to be initialized")
		}
	})

	t.Run("New cache without janitor", func(t *testing.T) {
		cache := newCacheWithJanitor(DefaultExpiration, 0, make(map[string]Item))
		if cache.defaultExpiration != DefaultExpiration {
			t.Errorf("Expected default expiration to be %v, got %v", DefaultExpiration, cache.defaultExpiration)
		}
		if cache.items == nil {
			t.Error("Expected items map to be initialized")
		}
		if cache.janitor != nil {
			t.Error("Expected janitor not to be initialized")
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("New cache with default expiration and cleanup interval", func(t *testing.T) {
		cache := New(DefaultExpiration, 1*time.Minute)
		if cache.defaultExpiration != DefaultExpiration {
			t.Errorf("Expected default expiration to be %v, got %v", DefaultExpiration, cache.defaultExpiration)
		}
		if cache.items == nil {
			t.Error("Expected items map to be initialized")
		}
		if cache.janitor == nil {
			t.Error("Expected janitor to be initialized")
		}
	})

	t.Run("New cache with NoExpiration and no cleanup interval", func(t *testing.T) {
		cache := New(NoExpiration, 0)
		if cache.defaultExpiration != NoExpiration {
			t.Errorf("Expected default expiration to be %v, got %v", NoExpiration, cache.defaultExpiration)
		}
		if cache.items == nil {
			t.Error("Expected items map to be initialized")
		}
		if cache.janitor != nil {
			t.Error("Expected janitor not to be initialized")
		}
	})
}

func TestNewFrom(t *testing.T) {
	t.Run("NewFrom with default expiration and cleanup interval", func(t *testing.T) {
		items := make(map[string]Item)
		cache := NewFrom(DefaultExpiration, 1*time.Minute, items)
		if cache.defaultExpiration != DefaultExpiration {
			t.Errorf("Expected default expiration to be %v, got %v", DefaultExpiration, cache.defaultExpiration)
		}
		if cache.items == nil || len(cache.items) != 0 {
			t.Error("Expected items map to be initialized and empty")
		}
		if cache.janitor == nil {
			t.Error("Expected janitor to be initialized")
		}
	})

	t.Run("NewFrom with NoExpiration and no cleanup interval", func(t *testing.T) {
		items := make(map[string]Item)
		cache := NewFrom(NoExpiration, 0, items)
		if cache.defaultExpiration != NoExpiration {
			t.Errorf("Expected default expiration to be %v, got %v", NoExpiration, cache.defaultExpiration)
		}
		if cache.items == nil || len(cache.items) != 0 {
			t.Error("Expected items map to be initialized and empty")
		}
		if cache.janitor != nil {
			t.Error("Expected janitor not to be initialized")
		}
	})

	t.Run("NewFrom with existing items map", func(t *testing.T) {
		items := map[string]Item{
			"key1": {Object: "value1", Expiration: 0},
		}
		cache := NewFrom(DefaultExpiration, 1*time.Minute, items)
		if cache.defaultExpiration != DefaultExpiration {
			t.Errorf("Expected default expiration to be %v, got %v", DefaultExpiration, cache.defaultExpiration)
		}
		if cache.items == nil || len(cache.items) != 1 {
			t.Error("Expected items map to be initialized with one item")
		}
		if cache.janitor == nil {
			t.Error("Expected janitor to be initialized")
		}
		if cache.items["key1"].Object != "value1" {
			t.Errorf("Expected item with key 'key1' to have value 'value1', got %v", cache.items["key1"].Object)
		}
	})
}
