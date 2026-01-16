# Sharded Cache Example

This example demonstrates how to use ShardedCache for high-concurrency scenarios.

## What You'll Learn

- Creating a sharded cache
- Understanding shard distribution
- Performance benefits of sharding
- Using the same API as standard cache
- Concurrent operations across shards

## Running This Example

```bash
cd examples/sharded
go run main.go
```

## Key Concepts

### What is Sharding?

Sharding splits the cache into multiple independent segments (shards), each with its own lock. This dramatically reduces lock contention under high concurrency.

### Performance Comparison

- **Standard Cache**: ~500,000 ops/sec (100 goroutines)
- **Sharded Cache**: ~2,000,000 ops/sec (100 goroutines)
- **Speedup**: 4x improvement

## Operations Demonstrated

1. **Creating Sharded Cache**
   - Configuring shard count (8, 16, 32, 64)
   - Default: 16 shards

2. **Concurrent Access**
   - Multiple goroutines reading/writing
   - Automatic key distribution

3. **Same API**
   - All standard Cache methods work
   - Transparent sharding

## When to Use Sharded Cache

Use sharded cache when:
- ✅ You have 100+ concurrent goroutines
- ✅ Profiling shows lock contention
- ✅ Maximum throughput is critical

Stick with standard cache when:
- ❌ Low concurrency (< 10 goroutines)
- ❌ Small dataset (< 1000 items)
- ❌ Simplicity preferred over performance

## Related Documentation

- [Sharded Cache Guide](../../docs/sharded-cache.md)
- [Performance Benchmarks](../../docs/sharded-cache.md#benchmarking)

## Next Steps

- [Concurrent Example](../concurrent/) - Compare standard vs sharded performance
- [Serialization Example](../serialization/) - Persist sharded caches
