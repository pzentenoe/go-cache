package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/pzentenoe/go-cache"
)

func main() {
	fmt.Println("=== Sharded Cache Example ===")

	// Create a sharded cache with 16 shards for high concurrency
	// Default expiration: 5 minutes, cleanup interval: 10 minutes
	sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 16)

	fmt.Println("Created sharded cache with 16 shards")

	// Basic operations (same API as standard cache)
	sc.Set("user:1", "Alice", cache.DefaultExpiration)
	sc.Set("user:2", "Bob", cache.DefaultExpiration)
	sc.SetDefault("user:3", "Charlie")

	fmt.Println("Set 3 users in cache")

	// Get values
	if val, found := sc.Get("user:1"); found {
		fmt.Printf("user:1 = %v\n", val)
	}

	// Demonstrate high concurrency with sharded cache
	fmt.Println("\n=== Concurrent Operations Demo ===")
	fmt.Println("Starting 100 goroutines writing to cache...")

	var wg sync.WaitGroup
	numGoroutines := 100
	operationsPerGoroutine := 1000

	start := time.Now()

	// Writers
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key:%d:%d", id, j)
				sc.Set(key, j, cache.DefaultExpiration)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	totalOps := numGoroutines * operationsPerGoroutine
	fmt.Printf("Completed %d write operations in %v\n", totalOps, elapsed)
	fmt.Printf("Throughput: %.0f ops/sec\n", float64(totalOps)/elapsed.Seconds())
	fmt.Printf("Total items in cache: %d\n", sc.ItemCount())

	// Demonstrate increment operations
	fmt.Println("\n=== Increment Operations ===")
	sc.Set("counter", int64(0), cache.NoExpiration)

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				sc.Increment("counter", 1)
			}
		}()
	}

	wg.Wait()

	if val, found := sc.Get("counter"); found {
		fmt.Printf("Counter after 10,000 concurrent increments: %v\n", val)
	}

	// Float operations
	sc.Set("price", 100.50, cache.DefaultExpiration)
	sc.IncrementFloat("price", 9.50)
	if val, found := sc.Get("price"); found {
		fmt.Printf("Price after increment: %.2f\n", val)
	}

	sc.DecrementFloat("price", 10.00)
	if val, found := sc.Get("price"); found {
		fmt.Printf("Price after decrement: %.2f\n", val)
	}

	// GetWithExpiration
	fmt.Println("\n=== Expiration Info ===")
	if val, expTime, found := sc.GetWithExpiration("user:1"); found {
		fmt.Printf("Value: %v\n", val)
		if !expTime.IsZero() {
			fmt.Printf("Expires at: %v\n", expTime.Format(time.RFC3339))
			fmt.Printf("Time until expiration: %v\n", time.Until(expTime))
		}
	}

	// OnEvicted callback
	fmt.Println("\n=== Eviction Callback ===")
	sc.OnEvicted(func(key string, value any) {
		fmt.Printf("Evicted: %s = %v\n", key, value)
	})

	sc.Set("temp", "will be deleted", cache.DefaultExpiration)
	sc.Delete("temp") // This triggers OnEvicted callback

	// Flush all shards
	fmt.Printf("\n=== Cleanup ===\n")
	fmt.Printf("Items before flush: %d\n", sc.ItemCount())
	sc.Flush()
	fmt.Printf("Items after flush: %d\n", sc.ItemCount())

	fmt.Println("\n=== Example Complete ===")
}
