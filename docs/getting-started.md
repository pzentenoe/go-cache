# Getting Started with go-cache

## Installation

Install go-cache using Go modules:

```bash
go get github.com/pzentenoe/go-cache
```

## Quick Start

### Basic Usage

```go
package main

import (
	"fmt"
	"time"
	"github.com/pzentenoe/go-cache"
)

func main() {
	// Create a cache with 5-minute default expiration
	// and 10-minute cleanup interval
	c := cache.New(5*time.Minute, 10*time.Minute)

	// Set a value
	c.Set("mykey", "myvalue", cache.DefaultExpiration)

	// Get a value
	if val, found := c.Get("mykey"); found {
		fmt.Println("Found:", val)
	}
}
```

## Core Concepts

### Expiration Times

go-cache supports three types of expiration:

1. **Default Expiration**: Uses the cache's default expiration time

```go
c.Set("key", "value", cache.DefaultExpiration)
c.SetDefault("key", "value") // Same as above
```

2. **Custom Expiration**: Specify a custom duration

```go
c.Set("key", "value", 30*time.Second)
```

3. **No Expiration**: Item never expires

```go
c.Set("key", "value", cache.NoExpiration)
```

### Automatic Cleanup (Janitor)

The janitor automatically removes expired items at regular intervals:

```go
// Create cache with cleanup every 10 minutes
c := cache.New(5*time.Minute, 10*time.Minute)

// Create cache without automatic cleanup
c := cache.New(5*time.Minute, 0)
```

When automatic cleanup is disabled, you can manually trigger it:

```go
c.DeleteExpired()
```

## Common Operations

### Set Operations

```go
// Set with default expiration
c.Set("user:1", "Alice", cache.DefaultExpiration)

// Set with custom expiration
c.Set("session", "abc123", 1*time.Hour)

// Set without expiration
c.Set("config", appConfig, cache.NoExpiration)

// Add (only if not exists)
err := c.Add("user:1", "Bob", cache.DefaultExpiration)
if err != nil {
// Key already exists
}

// Replace (only if exists)
err = c.Replace("user:1", "Charlie", cache.DefaultExpiration)
```

### Get Operations

```go
// Simple get
if val, found := c.Get("user:1"); found {
user := val.(string)
fmt.Println(user)
}

// Get with expiration info
if val, expTime, found := c.GetWithExpiration("user:1"); found {
user := val.(string)
if !expTime.IsZero() {
fmt.Println("Expires at:", expTime)
}
}
```

### Delete Operations

```go
// Delete single item
c.Delete("user:1")

// Delete all expired items
c.DeleteExpired()

// Delete all items
c.Flush()
```

### Increment/Decrement

```go
// Integer operations
c.Set("counter", int64(0), cache.NoExpiration)
c.Increment("counter", 1)
c.Decrement("counter", 1)

// Typed operations (with overflow protection)
c.Set("views", uint64(100), cache.NoExpiration)
result, err := c.IncrementUint64("views", 50)
if err != nil {
// Overflow would occur
}

// Float operations
c.Set("price", 19.99, cache.NoExpiration)
c.IncrementFloat("price", 5.00) // 24.99
c.DecrementFloat("price", 2.50) // 22.49
```

## Type Safety

go-cache stores values as `any` (interface{}). You need to type assert when retrieving:

```go
type User struct {
Name  string
Email string
}

user := User{Name: "Alice", Email: "alice@example.com"}
c.Set("user:1", user, cache.DefaultExpiration)

// Type assertion
if val, found := c.Get("user:1"); found {
user := val.(User)
fmt.Println(user.Name)
}
```

## Eviction Callbacks

Execute custom logic when items are evicted:

```go
c.OnEvicted(func (key string, value any) {
fmt.Printf("Evicted: %s = %v\n", key, value)
})

// Callback is triggered on:
// - Delete()
// - DeleteExpired() (for expired items)
// - Manual eviction
// NOT triggered on Set() overwrites
```

## Thread Safety

All operations are thread-safe. The cache uses `sync.RWMutex` for concurrent access:

```go
var wg sync.WaitGroup

// Safe concurrent writes
for i := 0; i < 100; i++ {
wg.Add(1)
go func (id int) {
defer wg.Done()
c.Set(fmt.Sprintf("key%d", id), id, cache.DefaultExpiration)
}(i)
}

wg.Wait()
```

## Best Practices

### Choosing Expiration Times

- **Short-lived data (1-5 minutes)**: API responses, session tokens
- **Medium-lived data (5-30 minutes)**: User preferences, configuration
- **Long-lived data (1+ hours)**: Static content, rarely-changing data
- **No expiration**: Use `cache.NoExpiration` sparingly (requires manual cleanup)

### When to Use Sharded Cache

Use `NewSharded()` when:
- ✅ You have 100+ concurrent goroutines accessing the cache
- ✅ Profiling shows lock contention on cache operations
- ✅ Maximum throughput is critical

Stick with standard `New()` when:
- ❌ Low concurrency (< 10 goroutines)
- ❌ Small dataset (< 1,000 items)
- ❌ Simplicity is more important than performance

### Memory Management

- Set appropriate expiration times to prevent unbounded growth
- Use `Flush()` to clear cache when needed (e.g., during maintenance)
- Monitor memory usage with `ItemCount()`
- Consider pausing the janitor during bulk operations

### Error Handling

Always check return values:
```go
// Check if item exists
if val, found := c.Get("key"); found {
    // Safe to use val
}

// Handle overflow errors
result, err := c.IncrementUint64("counter", 1)
if err != nil {
    log.Printf("Increment failed: %v", err)
}
```

### Type Safety

Use safe type assertions:
```go
// Safe: checks type before using
if val, found := c.Get("user"); found {
    if user, ok := val.(User); ok {
        fmt.Println(user.Name)
    }
}

// Unsafe: can panic
user := val.(User)  // Only use if you're certain of the type
```

## Next Steps

- [API Reference](api-reference.md) - Complete method documentation
- [Sharded Cache](sharded-cache.md) - High-concurrency usage
- [Serialization](serialization.md) - Persist cache to disk
- [Examples](../examples/) - Runnable code examples
