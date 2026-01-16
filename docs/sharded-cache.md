# Sharded Cache

High-performance cache with reduced lock contention for concurrent workloads.

## Overview

ShardedCache splits the cache into multiple independent shards, each with its own lock. This significantly improves
performance in high-concurrency scenarios by reducing lock contention.

## When to Use Sharded Cache

Use ShardedCache when:

- ✅ You have high concurrent read/write operations
- ✅ Multiple goroutines access the cache simultaneously
- ✅ Performance benchmarks show lock contention
- ✅ You need maximum throughput

Use standard Cache when:

- ❌ Low concurrency (single goroutine or few goroutines)
- ❌ Simplicity is more important than maximum performance
- ❌ Memory overhead of multiple shards is a concern

## Creating a Sharded Cache

```go
import "github.com/pzentenoe/go-cache"

// Create with 16 shards (recommended starting point)
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 16)
```

### Choosing the Number of Shards

**General Guidelines:**

- **8-16 shards**: Good starting point for most applications
- **32 shards**: High concurrency (100+ goroutines)
- **64 shards**: Extremely high concurrency (1000+ goroutines)
- **Power of 2**: Use powers of 2 (8, 16, 32, 64) for optimal hashing

**Trade-offs:**

- More shards = Less contention, Higher memory overhead
- Fewer shards = More contention, Lower memory overhead

**Example:**

```go
// Low concurrency
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 8)

// High concurrency
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 32)

// Very high concurrency
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 64)
```

## API Compatibility

ShardedCache implements the same interface as standard Cache:

```go
// All standard cache operations work identically
sc.Set("key", "value", cache.DefaultExpiration)
sc.Get("key")
sc.Delete("key")
sc.Increment("counter", 1)
sc.IncrementFloat("price", 5.0)
sc.SaveFile("cache.gob")
```

## Performance Comparison

### Standard Cache

```go
c := cache.New(5*time.Minute, 10*time.Minute)

// 100 goroutines, 1000 operations each
// All operations compete for the same lock
// Throughput: ~500,000 ops/sec
```

### Sharded Cache (16 shards)

```go
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 16)

// 100 goroutines, 1000 operations each
// Operations distributed across 16 shards
// Throughput: ~2,000,000 ops/sec (4x improvement)
```

## How Sharding Works

### Hash Distribution

ShardedCache uses DJB33 hash function to distribute keys:

```go
// Pseudo-code
func bucket(key string) *Cache {
hash := djb33(seed, key)
shardIndex := hash % numShards
return shards[shardIndex]
}
```

Keys are deterministically assigned to shards:

- Same key always goes to same shard
- Keys are evenly distributed
- No rebalancing needed

### Lock Contention Reduction

**Standard Cache:**

```
All operations → Single Lock → Cache
```

**Sharded Cache (4 shards example):**

```
Operations → Shard 0 (Lock 0)
          ↘ Shard 1 (Lock 1)
          ↘ Shard 2 (Lock 2)
          ↘ Shard 3 (Lock 3)
```

## Complete Example

```go
package main

import (
	"fmt"
	"sync"
	"time"
	"github.com/pzentenoe/go-cache"
)

func main() {
	// Create sharded cache with 32 shards
	sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 32)

	var wg sync.WaitGroup
	numGoroutines := 100
	operationsPerGoroutine := 1000

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key:%d:%d", id, j)
				sc.Set(key, j, cache.DefaultExpiration)
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("Total items: %d\n", sc.ItemCount())
	// Output: Total items: 100000
}
```

## Serialization

ShardedCache supports full serialization:

```go
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 16)

// Populate
for i := 0; i < 1000; i++ {
sc.Set(fmt.Sprintf("key%d", i), i, cache.DefaultExpiration)
}

// Save
sc.SaveFile("sharded.gob")

// Load (must use same number of shards!)
newSC := cache.NewSharded(5*time.Minute, 10*time.Minute, 16)
newSC.LoadFile("sharded.gob")
```

**Important:** When loading, the cache must have the same number of shards as when saved, or an error is returned.

## Advanced Operations

### ItemCount Across Shards

```go
count := sc.ItemCount()
// Returns total items across all shards
```

### OnEvicted Callback

```go
sc.OnEvicted(func (key string, value any) {
log.Printf("Evicted from shard: %s = %v", key, value)
})
// Callback is registered for ALL shards
```

### GetWithExpiration

```go
if val, expTime, found := sc.GetWithExpiration("session:123"); found {
fmt.Println("Value:", val)
if !expTime.IsZero() {
fmt.Println("TTL:", time.Until(expTime))
}
}
```

### Flush All Shards

```go
sc.Flush()
// Flushes ALL shards atomically
```

## Benchmarking

Example benchmark comparing standard vs sharded cache:

```go
package main

import (
	"fmt"
	"sync"
	"time"
	"github.com/pzentenoe/go-cache"
)

func benchmarkCache(name string, operations int, concurrent int, useSharded bool) {
	var c interface {
		Set(string, any, time.Duration)
		Get(string) (any, bool)
	}

	if useSharded {
		c = cache.NewSharded(5*time.Minute, 0, 32)
	} else {
		c = cache.New(5*time.Minute, 0)
	}

	start := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations/concurrent; j++ {
				key := fmt.Sprintf("key:%d:%d", id, j)
				c.Set(key, j, cache.DefaultExpiration)
				c.Get(key)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("%s: %d ops in %v (%.0f ops/sec)\n",
		name, operations, elapsed,
		float64(operations)/elapsed.Seconds())
}

func main() {
	ops := 1000000
	concurrent := 100

	benchmarkCache("Standard Cache", ops, concurrent, false)
	benchmarkCache("Sharded Cache (32 shards)", ops, concurrent, true)
}
```

Expected output (actual numbers may vary by system):

```
Standard Cache: 1000000 ops in 2.5s (400000 ops/sec)
Sharded Cache (32 shards): 1000000 ops in 600ms (1666666 ops/sec)
```

**Note:** Throughout this documentation, we use rounded reference numbers (~500,000 and ~2,000,000 ops/sec) for consistency and simplicity. Actual performance will vary based on hardware, Go version, workload characteristics, and system load. The key insight is the relative improvement: sharded cache typically delivers 3-4x better performance under high concurrency.

## Internal Structure

```
ShardedCache
├── Shard 0 (Cache instance)
│   ├── items map[string]Item
│   ├── mu sync.RWMutex
│   └── janitor
├── Shard 1 (Cache instance)
│   ├── items map[string]Item
│   ├── mu sync.RWMutex
│   └── janitor
...
└── Shard N (Cache instance)
    ├── items map[string]Item
    ├── mu sync.RWMutex
    └── janitor
```

Each shard:

- Has its own independent Cache instance
- Has its own lock (no contention with other shards)
- Has its own janitor for cleanup
- Operates completely independently

## Best Practices

1. **Choose appropriate shard count:**
   ```go
   // Start with 16, adjust based on profiling
   sc := cache.NewSharded(expiration, cleanup, 16)
   ```

2. **Use for high concurrency:**
   ```go
   // If you have 100+ concurrent goroutines
   sc := cache.NewSharded(expiration, cleanup, 32)
   ```

3. **Keep shard count consistent for serialization:**
   ```go
   const numShards = 16
   sc1 := cache.NewSharded(expiration, cleanup, numShards)
   sc1.SaveFile("cache.gob")

   sc2 := cache.NewSharded(expiration, cleanup, numShards) // Same count!
   sc2.LoadFile("cache.gob")
   ```

4. **Profile before optimizing:**
   ```go
   // Use pprof to identify if lock contention is an issue
   // Only switch to sharded if profiling shows contention
   ```

## Limitations

1. **Memory Overhead:**
    - Each shard has overhead for locks, maps, janitors
    - More shards = more memory

2. **Serialization Compatibility:**
    - Must use same shard count when loading
    - Cannot change shard count after data is saved

3. **Items() Method:**
    - Returns slice of maps (one per shard)
    - Not a single unified map like standard Cache

## Migration from Standard Cache

Easy drop-in replacement:

```go
// Before
c := cache.New(5*time.Minute, 10*time.Minute)

// After
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 16)

// All operations work the same!
sc.Set("key", "value", cache.DefaultExpiration)
sc.Get("key")
```

## See Also

- [API Reference](api-reference.md)
- [Getting Started Guide](getting-started.md)
- [Concurrent Operations Example](../examples/concurrent/)
- [Sharded Cache Example](../examples/sharded/)
