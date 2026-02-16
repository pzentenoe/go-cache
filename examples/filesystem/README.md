# Filesystem Persistence Example

This example demonstrates production-ready patterns for persisting cache data to the filesystem with crash safety and graceful shutdown.

## What You'll Learn

- Warm start (load cache from disk on startup)
- Atomic save (write to temp file, then rename)
- Periodic background save
- Graceful shutdown with final save
- Explicit `Close()` for janitor cleanup
- OS signal handling (Ctrl+C / SIGTERM)

## Running This Example

```bash
cd examples/filesystem
go run main.go
```

## Operations Demonstrated

1. **Warm Start**
   - Load existing cache from disk if available
   - Start fresh if no cache file exists
   - Graceful fallback on load errors

2. **Atomic Save**
   - Write to temporary file first
   - Rename to final path for crash safety
   - Clean up temp file on errors

3. **Periodic Background Save**
   - Save cache to disk at regular intervals
   - Runs in a separate goroutine
   - Stoppable via channel signal

4. **Graceful Shutdown**
   - Handle OS signals (SIGINT, SIGTERM)
   - Delete expired items before final save
   - Stop periodic saver and janitor with `Close()`

## Use Cases

- **Warm Cache on Restart**: Avoid cold start by loading previous cache data
- **Crash Recovery**: Atomic writes prevent corrupted cache files
- **Long-Running Services**: Periodic saves protect against data loss
- **Graceful Deploys**: Signal handling ensures clean shutdown

## Important Notes

- **Atomic Save Pattern**: Always write to a temp file and rename. This prevents partial writes from corrupting your cache file if the process crashes mid-save.
- **Signal Handling**: The example handles both `SIGINT` (Ctrl+C) and `SIGTERM` for clean shutdown in containerized environments.
- **Close()**: Always call `Close()` when you're done with the cache to stop the janitor goroutine and release resources.

## Related Documentation

- [Serialization Guide](../../docs/serialization.md)
- [Janitor Control Guide](../../docs/janitor-control.md)
- [Getting Started](../../docs/getting-started.md)

## Next Steps

- [Serialization Example](../serialization/) - Basic save/load operations
- [Janitor Example](../janitor/) - Runtime cleanup management
