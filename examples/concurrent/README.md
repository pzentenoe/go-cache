# Concurrent Operations Example

This example demonstrates thread-safe operations and performance under high concurrency.

## What You'll Learn

- Thread-safe read/write operations
- Concurrent increment operations
- Race condition prevention
- Performance comparison: standard vs sharded
- Best practices for concurrent access

## Running This Example

```bash
cd examples/concurrent
go run main.go

# Run with race detector
go run -race main.go
```

## Operations Demonstrated

1. **Concurrent Reads/Writes**
   - Multiple goroutines accessing cache
   - No data races
   - Automatic synchronization

2. **Atomic Increment**
   - Thread-safe counters
   - `Increment()` operations
   - Race-free increments

3. **Add Race Conditions**
   - Using `Add()` correctly
   - Understanding TOCTTOU issues
   - Proper concurrent patterns

4. **Performance Benchmarking**
   - Standard cache under load
   - Sharded cache under load
   - Real-world comparison

## Thread Safety Guarantees

All cache operations are thread-safe:
- ✅ `Get()` / `Set()` / `Delete()`
- ✅ `Increment()` / `Decrement()`
- ✅ `Add()` / `Replace()`
- ✅ `Flush()` / `ItemCount()`

## Performance Results

Example output with 100 concurrent goroutines:

```
Standard Cache:  ~500,000 ops/sec
Sharded Cache:   ~2,000,000 ops/sec
Speedup:         4x
```

## Common Patterns

### Safe Counter Pattern

```go
// Good: Thread-safe increment
c.Increment("counter", 1)

// Bad: Race condition
if val, found := c.Get("counter"); found {
    c.Set("counter", val.(int) + 1, cache.NoExpiration)
}
```

### Concurrent Access Pattern

```go
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        c.Set(fmt.Sprintf("key%d", id), id, cache.DefaultExpiration)
    }(i)
}
wg.Wait()
```

## Related Documentation

- [Sharded Cache Guide](../../docs/sharded-cache.md)
- [API Reference](../../docs/api-reference.md)

## Next Steps

- [Sharded Example](../sharded/) - Optimize for high concurrency
- [Overflow Example](../overflow/) - Safe numeric operations
