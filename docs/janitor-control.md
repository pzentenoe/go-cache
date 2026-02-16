# Janitor Control Guide

Runtime control over automatic cleanup of expired items.

## Overview

The janitor is a background goroutine that automatically deletes expired items at regular intervals. go-cache
provides methods to control janitor behavior at runtime.

## Default Behavior

When you create a cache with a cleanup interval, the janitor runs automatically:

```go
// Janitor runs every 10 minutes
c := cache.New(5*time.Minute, 10*time.Minute)
```

Without a janitor, you must manually clean up:

```go
// No automatic cleanup
c := cache.New(5*time.Minute, 0)

// Manual cleanup
c.DeleteExpired()
```

## Control Methods

### PauseJanitor()

Temporarily pauses automatic cleanup. Expired items remain in cache until janitor is resumed or `DeleteExpired()` is
called.

```go
c.PauseJanitor()
// Janitor is paused - no automatic cleanup
```

**Use cases:**

- Bulk operations where cleanup might interfere
- Temporary performance boost during high load
- Maintenance windows

### ResumeJanitor()

Resumes automatic cleanup after pause.

```go
c.ResumeJanitor()
// Janitor resumes - automatic cleanup continues
```

### SetJanitorInterval(d time.Duration)

Dynamically changes the cleanup interval.

```go
// Change to faster cleanup
c.SetJanitorInterval(1 * time.Minute)

// Change to slower cleanup
c.SetJanitorInterval(30 * time.Minute)
```

**Use cases:**

- Adjust cleanup frequency based on load
- Speed up cleanup during off-peak hours
- Slow down cleanup during peak hours

## Complete Example

```go
package main

import (
	"fmt"
	"time"
	"github.com/pzentenoe/go-cache"
)

func main() {
	// Start with 5-minute cleanup interval
	c := cache.New(1*time.Minute, 5*time.Minute)

	// Phase 1: Normal operation
	fmt.Println("Phase 1: Normal operation")
	c.Set("temp1", "value", 1*time.Minute)
	time.Sleep(2 * time.Minute)
	fmt.Printf("Items: %d (should be 0, cleaned by janitor)\n", c.ItemCount())

	// Phase 2: Pause for bulk operations
	fmt.Println("\nPhase 2: Pausing janitor for bulk operations")
	c.PauseJanitor()

	for i := 0; i < 100; i++ {
		c.Set(fmt.Sprintf("bulk:%d", i), i, 1*time.Minute)
	}
	fmt.Printf("Added 100 items\n")

	time.Sleep(2 * time.Minute)
	fmt.Printf("After expiration (paused): %d items\n", c.ItemCount())
	// Items are expired but not deleted

	// Phase 3: Speed up cleanup and resume
	fmt.Println("\nPhase 3: Speeding up cleanup")
	c.SetJanitorInterval(10 * time.Second)
	c.ResumeJanitor()

	time.Sleep(15 * time.Second)
	fmt.Printf("After fast cleanup: %d items\n", c.ItemCount())
}
```

## Pause/Resume Workflow

```go
c := cache.New(100*time.Millisecond, 50*time.Millisecond)

// Add items
for i := 0; i < 10; i++ {
c.Set(fmt.Sprintf("key%d", i), i, 100*time.Millisecond)
}
fmt.Printf("Added 10 items\n")

// Pause janitor
c.PauseJanitor()
fmt.Println("Janitor paused")

// Wait for items to expire
time.Sleep(150 * time.Millisecond)
fmt.Printf("After expiration (paused): %d items still in cache\n", c.ItemCount())
// Output: 10 items (expired but not cleaned)

// Resume janitor
c.ResumeJanitor()
fmt.Println("Janitor resumed")

// Wait for janitor to run
time.Sleep(100 * time.Millisecond)
fmt.Printf("After resume: %d items\n", c.ItemCount())
// Output: 0 items (cleaned up)
```

## Dynamic Interval Changes

### Adapting to Load

```go
func adjustCleanupInterval(c *cache.Cache, load int) {
switch {
case load > 1000:
// High load - slow down cleanup
c.SetJanitorInterval(30 * time.Minute)
case load > 100:
// Medium load - normal cleanup
c.SetJanitorInterval(10 * time.Minute)
default:
// Low load - aggressive cleanup
c.SetJanitorInterval(1 * time.Minute)
}
}
```

### Time-based Adjustment

```go
// Aggressive cleanup during off-peak hours
go func () {
ticker := time.NewTicker(1 * time.Hour)
for range ticker.C {
hour := time.Now().Hour()

if hour >= 2 && hour <= 6 {
// 2 AM - 6 AM: Aggressive cleanup
c.SetJanitorInterval(1 * time.Minute)
} else {
// Business hours: Conservative cleanup
c.SetJanitorInterval(15 * time.Minute)
}
}
}()
```

## Multiple Pause/Resume Cycles

```go
c := cache.New(50*time.Millisecond, 25*time.Millisecond)

// Cycle 1
c.PauseJanitor()
c.Set("temp1", "value", 50*time.Millisecond)
time.Sleep(100 * time.Millisecond)
c.ResumeJanitor()
time.Sleep(50 * time.Millisecond)

// Cycle 2
c.PauseJanitor()
c.Set("temp2", "value", 50*time.Millisecond)
time.Sleep(100 * time.Millisecond)
c.ResumeJanitor()
time.Sleep(50 * time.Millisecond)

// Both cycles work independently
```

## Manual Cleanup Alternative

If you don't need automatic cleanup, disable the janitor and use manual cleanup:

```go
// No janitor
c := cache.New(5*time.Minute, 0)

// Manual cleanup when needed
go func () {
ticker := time.NewTicker(10 * time.Minute)
for range ticker.C {
c.DeleteExpired()
}
}()

// Or cleanup on-demand
func cleanupIfNeeded(c *cache.Cache, maxItems int) {
if c.ItemCount() > maxItems {
c.DeleteExpired()
}
}
```

## Combining Controls

```go
c := cache.New(100*time.Millisecond, 100*time.Millisecond)

// Start with normal interval
fmt.Println("Starting with 100ms interval")

// Pause for bulk operation
c.PauseJanitor()
for i := 0; i < 1000; i++ {
c.Set(fmt.Sprintf("key%d", i), i, 100*time.Millisecond)
}

// Change interval while paused
c.SetJanitorInterval(10 * time.Millisecond)

// Resume with new interval
c.ResumeJanitor()

// Fast cleanup begins
time.Sleep(50 * time.Millisecond)
fmt.Printf("Items after fast cleanup: %d\n", c.ItemCount())
```

## Safe Usage with Nil Janitor

All control methods are safe to call even when janitor is not running:

```go
// Cache without janitor
c := cache.New(5*time.Minute, 0)

// These are safe (no-op)
c.PauseJanitor()
c.ResumeJanitor()
c.SetJanitorInterval(1 * time.Minute)
```

## Performance Considerations

### Pause/Resume Overhead

Pausing and resuming have minimal overhead:

- Non-blocking, mutex-protected operation
- Idempotent: safe to call multiple times (no deadlock risk)
- Immediate effect

### Interval Change Overhead

Changing interval requires:

- Stopping old ticker
- Creating new ticker
- Small memory allocation

### Best Practices

1. **Don't pause unnecessarily:**

```go
// Bad - unnecessary pause for single operation
c.PauseJanitor()
c.Set("key", "value", cache.DefaultExpiration)
c.ResumeJanitor()

// Good - pause only for bulk operations
c.PauseJanitor()
for i := 0; i < 10000; i++ {
c.Set(fmt.Sprintf("key%d", i), i, cache.DefaultExpiration)
}
c.ResumeJanitor()
```

2. **Balance cleanup frequency:**

```go
// Too aggressive - wastes CPU
c.SetJanitorInterval(100 * time.Millisecond)

// Too slow - memory builds up
c.SetJanitorInterval(24 * time.Hour)

// Balanced - adjust based on your use case
c.SetJanitorInterval(5 * time.Minute)
```

3. **Use pause during critical sections:**

```go
// Pause during critical operation
c.PauseJanitor()
err := performCriticalOperation()
c.ResumeJanitor()

if err != nil {
// Handle error
}
```

## Monitoring Janitor Activity

Track cleanup activity with eviction callbacks:

```go
var cleanupCount int64

c.OnEvicted(func (key string, value any) {
atomic.AddInt64(&cleanupCount, 1)
})

// Monitor cleanup rate
go func () {
ticker := time.NewTicker(1 * time.Minute)
for range ticker.C {
count := atomic.SwapInt64(&cleanupCount, 0)
fmt.Printf("Items cleaned in last minute: %d\n", count)

// Adjust interval based on cleanup rate
if count > 1000 {
c.SetJanitorInterval(30 * time.Second)
} else {
c.SetJanitorInterval(5 * time.Minute)
}
}
}()
```

## Common Patterns

### Bulk Loading

```go
func bulkLoad(c *cache.Cache, items map[string]any) {
c.PauseJanitor()
defer c.ResumeJanitor()

for k, v := range items {
c.Set(k, v, cache.DefaultExpiration)
}
}
```

### Maintenance Window

```go
func performMaintenance(c *cache.Cache) {
c.PauseJanitor()
defer c.ResumeJanitor()

// Perform maintenance tasks
c.DeleteExpired()

// Compact or reorganize
items := c.Items()
c.Flush()
for k, item := range items {
if !item.Expired() {
c.Set(k, item.Object, time.Until(time.Unix(0, item.Expiration)))
}
}
}
```

### Graceful Shutdown

```go
func shutdown(c *cache.Cache) {
// Clean up expired items one last time
c.DeleteExpired()

// Save to disk
c.SaveFile("cache_backup.gob")

// Stop janitor goroutine explicitly
c.Close()
}
```

## Sharded Cache Janitor Control

The sharded cache supports the same janitor control methods as the standard cache:

```go
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 16)

// Pause automatic cleanup across all shards
sc.PauseJanitor()

// Resume automatic cleanup
sc.ResumeJanitor()

// Change cleanup interval
sc.SetJanitorInterval(1 * time.Minute)
```

The sharded cache uses a single janitor that iterates over all shards when cleaning up expired items.

## See Also

- [API Reference](api-reference.md)
- [Sharded Cache](sharded-cache.md)
- [Getting Started](getting-started.md)
- [Janitor Control Example](../examples/janitor/) - Hands-on demonstration
- [Best Practices](getting-started.md#best-practices) - Memory management tips
