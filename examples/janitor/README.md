# Janitor Control Example

This example demonstrates runtime control over automatic cleanup (janitor).

## What You'll Learn

- How the janitor works
- Pausing automatic cleanup
- Resuming automatic cleanup
- Changing cleanup interval dynamically
- Manual cleanup with `DeleteExpired()`

## Running This Example

```bash
cd examples/janitor
go run main.go
```

## Operations Demonstrated

### 1. Automatic Janitor

The janitor runs in the background, removing expired items:

```go
// Janitor runs cleanup every 10 minutes
c := cache.New(5*time.Minute, 10*time.Minute)
```

### 2. Pause Janitor

Temporarily stop automatic cleanup:

```go
c.PauseJanitor()
// Expired items remain in cache
// Useful for temporary performance boost
```

### 3. Resume Janitor

Restart automatic cleanup:

```go
c.ResumeJanitor()
// Janitor resumes with original interval
// Immediately cleans up accumulated expired items
```

### 4. Change Interval

Dynamically adjust cleanup frequency:

```go
c.SetJanitorInterval(5 * time.Minute)
// Janitor now runs every 5 minutes
// Takes effect immediately
```

### 5. Manual Cleanup

Clean up expired items on-demand:

```go
c.DeleteExpired()
// Removes all expired items immediately
// Works even when janitor is paused
```

## When to Use Janitor Controls

### Pause Janitor

Use cases:
- **Bulk Operations**: Pause during large batch inserts
- **Performance Critical**: Pause during high-load periods
- **Testing**: Deterministic behavior in tests
- **Maintenance**: Prevent cleanup during migrations

### Change Interval

Use cases:
- **Off-Peak Cleanup**: Less frequent during business hours
- **Adaptive Cleanup**: Adjust based on load
- **Memory Pressure**: More frequent when memory is tight
- **Power Saving**: Less frequent on battery power

### Manual Cleanup

Use cases:
- **Controlled Timing**: Clean up at specific moments
- **Before Snapshots**: Clean before serialization
- **Memory Recovery**: Force cleanup when needed
- **Testing**: Verify expiration behavior

## Performance Considerations

### Janitor Overhead

- Minimal: Runs in separate goroutine
- Sleep-based: No busy waiting
- Lock-based: Brief write lock during cleanup

### Pausing Benefits

- Eliminates cleanup overhead
- Improves write performance ~5-10%
- Trades memory for speed

### Trade-offs

```
Frequent Cleanup (1min):
  ✅ Lower memory usage
  ❌ Slightly higher CPU

Infrequent Cleanup (60min):
  ✅ Lower CPU usage
  ❌ Higher memory usage

Paused:
  ✅ Maximum performance
  ❌ Memory grows unbounded
```

## Example Patterns

### Batch Operation Pattern

```go
c.PauseJanitor()
// Perform bulk operations
for i := 0; i < 10000; i++ {
    c.Set(fmt.Sprintf("key%d", i), i, cache.DefaultExpiration)
}
c.ResumeJanitor()
```

### Adaptive Interval Pattern

```go
if highLoad {
    c.SetJanitorInterval(60 * time.Minute)  // Less frequent
} else {
    c.SetJanitorInterval(10 * time.Minute)  // More frequent
}
```

## Related Documentation

- [Janitor Control Guide](../../docs/janitor-control.md)
- [API Reference](../../docs/api-reference.md)

## Next Steps

- [Basic Example](../basic/) - Understand core operations
- [Concurrent Example](../concurrent/) - See janitor thread-safety
