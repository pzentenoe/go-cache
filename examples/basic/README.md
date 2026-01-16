# Basic Usage Example

This example demonstrates the fundamental operations of go-cache.

## What You'll Learn

- Creating a cache instance
- Set/Get/Delete operations
- Working with different expiration times
- Using Add and Replace methods
- Getting items with expiration info
- Counting items in the cache

## Running This Example

```bash
cd examples/basic
go run main.go
```

## Key Operations Demonstrated

1. **Creating a Cache**
   - Default expiration time
   - Cleanup interval configuration

2. **Setting Values**
   - `Set()` with default expiration
   - `Set()` with custom expiration
   - `SetDefault()` shorthand

3. **Getting Values**
   - `Get()` basic retrieval
   - `GetWithExpiration()` with timestamp

4. **Conditional Operations**
   - `Add()` - only if key doesn't exist
   - `Replace()` - only if key exists

5. **Cleanup Operations**
   - `Delete()` individual items
   - `Flush()` clear all items
   - `DeleteExpired()` manual cleanup

## Related Documentation

- [Getting Started Guide](../../docs/getting-started.md)
- [API Reference](../../docs/api-reference.md)

## Next Steps

After mastering basic operations, explore:
- [Concurrent Example](../concurrent/) - Thread-safe operations and patterns
- [Sharded Cache Example](../sharded/) - High-concurrency scenarios
- [Serialization Example](../serialization/) - Persistence
- [Overflow Example](../overflow/) - Safe numeric operations
- [Janitor Example](../janitor/) - Cleanup control
