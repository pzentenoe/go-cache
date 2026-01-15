package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pzentenoe/go-cache"
)

type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

func main() {
	fmt.Println("=== Serialization Example ===\n")

	// Create cache and populate with data
	c := cache.New(5*time.Minute, 10*time.Minute)

	fmt.Println("Populating cache with data...")

	// Add various types of data
	c.Set("string_key", "Hello, World!", cache.DefaultExpiration)
	c.Set("int_key", 42, cache.DefaultExpiration)
	c.Set("float_key", 3.14159, cache.DefaultExpiration)
	c.Set("bool_key", true, cache.DefaultExpiration)

	users := []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()},
		{ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Now()},
		{ID: 3, Name: "Charlie", Email: "charlie@example.com", CreatedAt: time.Now()},
	}

	for _, user := range users {
		key := fmt.Sprintf("user:%d", user.ID)
		c.Set(key, user, cache.DefaultExpiration)
	}

	// Add a map
	config := map[string]interface{}{
		"max_connections": 100,
		"timeout":         30 * time.Second,
		"enabled":         true,
	}
	c.Set("config", config, cache.DefaultExpiration)

	fmt.Printf("Cache populated with %d items\n\n", c.ItemCount())

	// Save to file
	filename := "cache_data.gob"
	fmt.Printf("Saving cache to %s...\n", filename)
	err := c.SaveFile(filename)
	if err != nil {
		fmt.Printf("Error saving cache: %v\n", err)
		return
	}
	fmt.Println("Cache saved successfully")

	// Create a new empty cache
	fmt.Println("\nCreating new empty cache...")
	newCache := cache.New(5*time.Minute, 10*time.Minute)
	fmt.Printf("New cache has %d items\n", newCache.ItemCount())

	// Load from file
	fmt.Printf("\nLoading cache from %s...\n", filename)
	err = newCache.LoadFile(filename)
	if err != nil {
		fmt.Printf("Error loading cache: %v\n", err)
		return
	}
	fmt.Printf("Cache loaded successfully. Now has %d items\n\n", newCache.ItemCount())

	// Verify loaded data
	fmt.Println("=== Verifying Loaded Data ===")

	if val, found := newCache.Get("string_key"); found {
		fmt.Printf("string_key: %v\n", val)
	}

	if val, found := newCache.Get("user:1"); found {
		user := val.(User)
		fmt.Printf("user:1: %+v\n", user)
	}

	if val, found := newCache.Get("config"); found {
		fmt.Printf("config: %v\n", val)
	}

	// Clean up
	fmt.Println("\n=== Cleanup ===")
	err = os.Remove(filename)
	if err != nil {
		fmt.Printf("Warning: Could not delete %s: %v\n", filename, err)
	} else {
		fmt.Printf("Deleted %s\n", filename)
	}

	// Demonstrate sharded cache serialization
	fmt.Println("\n=== Sharded Cache Serialization ===")

	sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 8)

	// Populate sharded cache
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("item:%d", i)
		sc.Set(key, i*10, cache.DefaultExpiration)
	}

	fmt.Printf("Sharded cache populated with %d items\n", sc.ItemCount())

	// Save sharded cache
	shardedFilename := "sharded_cache.gob"
	err = sc.SaveFile(shardedFilename)
	if err != nil {
		fmt.Printf("Error saving sharded cache: %v\n", err)
		return
	}
	fmt.Printf("Sharded cache saved to %s\n", shardedFilename)

	// Load into new sharded cache (must have same number of shards)
	newSC := cache.NewSharded(5*time.Minute, 10*time.Minute, 8)
	err = newSC.LoadFile(shardedFilename)
	if err != nil {
		fmt.Printf("Error loading sharded cache: %v\n", err)
		return
	}
	fmt.Printf("Sharded cache loaded successfully. Now has %d items\n", newSC.ItemCount())

	// Verify some items
	if val, found := newSC.Get("item:42"); found {
		fmt.Printf("item:42: %v\n", val)
	}

	// Clean up sharded cache file
	os.Remove(shardedFilename)

	fmt.Println("\n=== Example Complete ===")
}
