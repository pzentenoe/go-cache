package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pzentenoe/go-cache"
)

func main() {
	fmt.Println("=== Concurrent Operations Example ===")

	// Standard cache concurrent operations
	demonstrateStandardCacheConcurrency()

	// Sharded cache for better concurrent performance
	demonstrateShardedCacheConcurrency()

	// Concurrent increment operations
	demonstrateConcurrentIncrements()
}

func demonstrateStandardCacheConcurrency() {
	fmt.Println("=== Standard Cache Concurrency ===")
	c := cache.New(5*time.Minute, 10*time.Minute)

	var wg sync.WaitGroup
	numReaders := 50
	numWriters := 50
	operations := 1000

	start := time.Now()

	// Writers
	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := fmt.Sprintf("key:%d:%d", id, j)
				c.Set(key, j, cache.DefaultExpiration)
			}
		}(i)
	}

	// Readers
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := fmt.Sprintf("key:%d:%d", id%numWriters, j)
				c.Get(key)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	totalOps := (numReaders + numWriters) * operations
	fmt.Printf("Completed %d operations in %v\n", totalOps, elapsed)
	fmt.Printf("Throughput: %.0f ops/sec\n", float64(totalOps)/elapsed.Seconds())
	fmt.Printf("Final cache size: %d items\n\n", c.ItemCount())
}

func demonstrateShardedCacheConcurrency() {
	fmt.Println("=== Sharded Cache Concurrency (Better Performance) ===")
	sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 32) // 32 shards

	var wg sync.WaitGroup
	numReaders := 50
	numWriters := 50
	operations := 1000

	start := time.Now()

	// Writers
	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := fmt.Sprintf("key:%d:%d", id, j)
				sc.Set(key, j, cache.DefaultExpiration)
			}
		}(i)
	}

	// Readers
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := fmt.Sprintf("key:%d:%d", id%numWriters, j)
				sc.Get(key)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	totalOps := (numReaders + numWriters) * operations
	fmt.Printf("Completed %d operations in %v\n", totalOps, elapsed)
	fmt.Printf("Throughput: %.0f ops/sec\n", float64(totalOps)/elapsed.Seconds())
	fmt.Printf("Final cache size: %d items\n\n", sc.ItemCount())
}

func demonstrateConcurrentIncrements() {
	fmt.Println("=== Concurrent Increment Operations ===")
	c := cache.New(cache.NoExpiration, 0)

	// Initialize counters
	c.Set("counter_int", int64(0), cache.NoExpiration)
	c.Set("counter_float", float64(0), cache.NoExpiration)

	var wg sync.WaitGroup
	numGoroutines := 100
	incrementsPerGoroutine := 1000

	start := time.Now()

	// Concurrent integer increments
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				c.Increment("counter_int", 1)
			}
		}()
	}

	// Concurrent float increments
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				c.IncrementFloat("counter_float", 0.5)
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	expectedInt := int64(numGoroutines * incrementsPerGoroutine)
	expectedFloat := float64(numGoroutines*incrementsPerGoroutine) * 0.5

	if val, found := c.Get("counter_int"); found {
		fmt.Printf("Integer counter: %v (expected: %d)\n", val, expectedInt)
		if val.(int64) == expectedInt {
			fmt.Println("✓ Integer increments are thread-safe!")
		}
	}

	if val, found := c.Get("counter_float"); found {
		fmt.Printf("Float counter: %.1f (expected: %.1f)\n", val, expectedFloat)
		if val.(float64) == expectedFloat {
			fmt.Println("✓ Float increments are thread-safe!")
		}
	}

	fmt.Printf("Completed in %v\n\n", elapsed)

	// Demonstrate concurrent Add operations (only one should succeed)
	fmt.Println("=== Concurrent Add Operations (Race Condition Demo) ===")
	c2 := cache.New(cache.NoExpiration, 0)

	var successCount int32
	numAdders := 100

	wg.Add(numAdders)
	for i := 0; i < numAdders; i++ {
		go func(id int) {
			defer wg.Done()
			err := c2.Add("exclusive_key", id, cache.NoExpiration)
			if err == nil {
				atomic.AddInt32(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("Out of %d concurrent Add attempts, %d succeeded\n", numAdders, successCount)
	if val, found := c2.Get("exclusive_key"); found {
		fmt.Printf("Winner value: %v\n", val)
	}

	if successCount == 1 {
		fmt.Println("✓ Add operation is thread-safe (only 1 succeeded)!")
	}

	fmt.Println("\n=== Example Complete ===")
}
