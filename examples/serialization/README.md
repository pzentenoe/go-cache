# Serialization Example

This example demonstrates how to persist cache data to disk and restore it.

## What You'll Learn

- Saving cache to file
- Loading cache from file
- Serializing custom types
- Sharded cache serialization
- Data persistence patterns

## Running This Example

```bash
cd examples/serialization
go run main.go
```

## Operations Demonstrated

1. **Saving to Disk**
   - `SaveFile()` - write cache to file
   - Gob encoding format
   - Preserves expiration times

2. **Loading from Disk**
   - `LoadFile()` - restore cache from file
   - Validates data integrity
   - Maintains cache structure

3. **Custom Types**
   - Structs and complex types
   - Type registration for Gob
   - Safe serialization

4. **Sharded Cache**
   - Shard-aware serialization
   - Correct key redistribution
   - Shard count validation

## Use Cases

- **Application Restart**: Preserve cache across restarts
- **Backup**: Create cache snapshots
- **Migration**: Transfer cache between instances
- **Pre-population**: Load initial data from file

## Important Notes

⚠️ **Gob Encoding Requirements**:
- Custom types must be registered with `gob.Register()`
- Fields must be exported (capitalized)
- Circular references not supported

⚠️ **Sharded Cache**:
- Shard count must match on load
- Changing shard count requires regenerating the file

## Related Documentation

- [Serialization Guide](../../docs/serialization.md)
- [Getting Started](../../docs/getting-started.md)

## Next Steps

- [Basic Example](../basic/) - Learn core operations first
- [Sharded Example](../sharded/) - Understand sharding before serializing
