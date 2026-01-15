# go-cache Examples

This directory contains runnable examples demonstrating various features of go-cache.

## Running Examples

Each example is a standalone Go program with its own `main.go` file:

```bash
# Run basic example
cd basic && go run main.go

# Run sharded cache example
cd sharded && go run main.go

# Run serialization example
cd serialization && go run main.go

# Run concurrent operations example
cd concurrent && go run main.go

# Run overflow protection example
cd overflow && go run main.go

# Run janitor control example
cd janitor && go run main.go
```

## Examples Overview

### [basic/](basic/)
**Core operations and expiration**

Demonstrates:
- Creating a cache
- Set/Get/Delete operations
- Expiration (default, custom, none)
- Add and Replace operations
- GetWithExpiration
- Item counting
- Flush operation

**Best for:** Learning the fundamentals

### [sharded/](sharded/)
**High-concurrency patterns**

Demonstrates:
- Creating a sharded cache
- Performance comparison with standard cache
- Concurrent read/write operations
- Increment operations across shards
- OnEvicted callbacks
- Managing multiple shards

**Best for:** High-performance, concurrent applications

### [serialization/](serialization/)
**Save/Load cache data**

Demonstrates:
- Saving cache to file
- Loading cache from file
- Serializing custom types
- Sharded cache serialization
- Data persistence patterns

**Best for:** Applications requiring cache persistence

### [concurrent/](concurrent/)
**Thread-safe operations**

Demonstrates:
- Concurrent reads and writes
- Standard vs sharded cache performance
- Thread-safe increment operations
- Add operation race conditions
- Benchmarking concurrent operations

**Best for:** Understanding concurrency behavior

### [overflow/](overflow/)
**Numeric operation safety**

Demonstrates:
- Unsigned integer overflow protection
- Signed integer overflow/underflow protection
- Float infinity protection
- Safe operations within bounds
- All numeric types (uint8, uint16, uint32, uint64, int8, int16, int32, int64, float32, float64)

**Best for:** Applications with numeric counters

### [janitor/](janitor/)
**Runtime cleanup management**

Demonstrates:
- Automatic janitor cleanup
- Pausing and resuming janitor
- Dynamically changing cleanup interval
- Combined janitor controls
- Manual cleanup with DeleteExpired

**Best for:** Fine-tuning cache cleanup behavior

## Quick Start

If you're new to go-cache, start with these examples in order:

1. **basic/** - Learn the fundamentals
2. **serialization/** - Understand persistence
3. **concurrent/** - See thread-safety in action
4. **sharded/** - Learn when to use sharding

For specific use cases:
- Need high performance? → **sharded/**
- Need numeric counters? → **overflow/**
- Need cleanup control? → **janitor/**

## Example Code Structure

Each example follows this structure:

```
example-name/
└── main.go    # Standalone executable
```

All examples:
- Are self-contained
- Include comments explaining each step
- Can be run independently
- Use realistic scenarios

## Common Patterns

### Creating a Cache

```go
// Standard cache
c := cache.New(5*time.Minute, 10*time.Minute)

// Sharded cache (for high concurrency)
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 16)
```

### Basic Operations

```go
// Set
c.Set("key", "value", cache.DefaultExpiration)

// Get
if val, found := c.Get("key"); found {
    fmt.Println(val)
}

// Delete
c.Delete("key")
```

### Type Assertions

Since cache stores `interface{}`, you need type assertions:

```go
// String
if val, found := c.Get("name"); found {
    name := val.(string)
}

// Struct
type User struct { Name string }
if val, found := c.Get("user"); found {
    user := val.(User)
}

// Numbers
if val, found := c.Get("count"); found {
    count := val.(int64)
}
```

## Modifying Examples

Feel free to modify any example to experiment:

```go
// Change expiration times
c := cache.New(1*time.Minute, 30*time.Second)

// Try different shard counts
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 32)

// Experiment with concurrency
numGoroutines := 1000  // Increase concurrency
```

## Need Help?

- 📖 Read the [documentation](../docs/)
- 💬 Open an [issue](https://github.com/pzentenoe/go-cache/issues)
- ⭐ Star the [repository](https://github.com/pzentenoe/go-cache)

## Next Steps

After running the examples:

1. Read the [API Reference](../docs/api-reference.md) for complete method documentation
2. Check the [Sharded Cache Guide](../docs/sharded-cache.md) for performance tuning
3. Review [Best Practices](../docs/getting-started.md#best-practices)
4. Integrate go-cache into your application

---

Happy caching! 🚀