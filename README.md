# go-cache

[![codecov](https://codecov.io/github/pzentenoe/go-cache/graph/badge.svg?token=3W164MZ18S)](https://codecov.io/github/pzentenoe/go-cache)
![CI](https://github.com/pzentenoe/go-cache/actions/workflows/actions.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/pzentenoe/go-cache)](https://goreportcard.com/report/github.com/pzentenoe/go-cache)
![License](https://img.shields.io/github/license/pzentenoe/go-cache.svg)
![GitHub release](https://img.shields.io/github/v/release/pzentenoe/go-cache.svg)

High-performance, thread-safe in-memory cache for Go with expiration support and advanced features.

## Features

- ⚡ **High Performance**: Optimized for concurrent access with optional sharding
- 🔒 **Thread-Safe**: All operations protected with RWMutex
- ⏰ **Flexible Expiration**: Per-item, default, or no expiration
- 🧹 **Janitor**: Background cleanup with runtime control
- 💾 **Serialization**: Persist cache to disk using Gob encoding
- 🛡️ **Overflow Protection**: Built-in protection for numeric operations
- 📊 **High Concurrency**: Sharded cache for reduced lock contention (2-4x faster)
- 🎯 **Simple API**: Intuitive interface, easy to integrate

## Quick Start

### Installation

```bash
go get github.com/pzentenoe/go-cache
```

### Basic Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/pzentenoe/go-cache"
)

func main() {
    // Create cache with 5-minute default expiration and 10-minute cleanup
    c := cache.New(5*time.Minute, 10*time.Minute)

    // Set a value
    c.Set("mykey", "myvalue", cache.DefaultExpiration)

    // Get a value
    if val, found := c.Get("mykey"); found {
        fmt.Println("Found:", val)
    }
}
```

### High-Concurrency Example

```go
// Use sharded cache for high-concurrency workloads
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 32)

// Same API as standard cache
sc.Set("key", "value", cache.DefaultExpiration)
val, found := sc.Get("key")
```

## Documentation

### 📚 Guides

- [Getting Started](docs/getting-started.md) - Installation, basic usage, core concepts
- [API Reference](docs/api-reference.md) - Complete method documentation
- [Sharded Cache](docs/sharded-cache.md) - High-concurrency usage guide
- [Serialization](docs/serialization.md) - Persist cache to disk
- [Janitor Control](docs/janitor-control.md) - Runtime cleanup management

### 💡 Examples

Runnable examples in [`examples/`](examples/):

- [Basic Usage](examples/basic/) - Core operations and expiration
- [Sharded Cache](examples/sharded/) - High-concurrency patterns
- [Serialization](examples/serialization/) - Save/Load cache data
- [Concurrent Operations](examples/concurrent/) - Thread-safe operations
- [Overflow Protection](examples/overflow/) - Numeric operation safety
- [Janitor Control](examples/janitor/) - Runtime cleanup management

Run any example:

```bash
cd examples/basic && go run main.go
```

## Core Operations

### Basic CRUD

```go
// Set operations
c.Set("key", "value", cache.DefaultExpiration)
c.SetDefault("key", "value")
c.Add("key", "value", 5*time.Minute)    // Only if not exists
c.Replace("key", "new", 5*time.Minute)  // Only if exists

// Get operations
val, found := c.Get("key")
val, expTime, found := c.GetWithExpiration("key")

// Delete operations
c.Delete("key")
c.DeleteExpired() // Remove expired items
c.Flush()         // Remove all items
```

### Numeric Operations

All numeric operations include overflow/underflow protection:

```go
// Increment/Decrement
c.Increment("counter", 1)
c.Decrement("counter", 1)
c.IncrementFloat("price", 5.50)

// Type-safe operations with error handling
result, err := c.IncrementUint64("views", 100)
if err != nil {
    // Overflow would occur
}
```

### Persistence

```go
// Save and load cache to/from disk
c.SaveFile("cache.gob")
c.LoadFile("cache.gob")
```

### Janitor Control

```go
// Runtime cleanup management
c.PauseJanitor()
c.ResumeJanitor()
c.SetJanitorInterval(5 * time.Minute)
```

## Recent Updates

The latest release adds overflow protection, janitor control, and complete ShardedCache API parity. See [CHANGELOG.md](CHANGELOG.md) for full version history and release notes.

## Performance

### Standard Cache

- Suitable for most applications
- Single lock for all operations
- ~500,000 ops/sec with 100 concurrent goroutines

### Sharded Cache

- Recommended for high-concurrency scenarios
- Multiple independent caches with separate locks
- ~2,000,000 ops/sec with 100 concurrent goroutines (4x improvement)
- Configurable shard count (8, 16, 32, 64)

**When to use sharded cache:**

- ✅ High concurrent read/write operations (100+ goroutines)
- ✅ Lock contention identified in profiling
- ✅ Maximum throughput required

See [Sharded Cache Guide](docs/sharded-cache.md) for benchmarks and best practices.

## Thread Safety

All operations are thread-safe and can be called from multiple goroutines:

```go
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        c.Set(fmt.Sprintf("key%d", id), id, cache.DefaultExpiration)
        c.Get(fmt.Sprintf("key%d", id))
    }(i)
}
wg.Wait()
```

## Working with Types

The cache uses `interface{}` internally, supporting any Go type. Use type assertions when retrieving values:

```go
// Storing values
c.Set("user", User{Name: "Alice"}, cache.DefaultExpiration)
c.Set("count", 42, cache.DefaultExpiration)

// Retrieving with type assertion
if val, found := c.Get("user"); found {
    user := val.(User)  // Type assertion
    fmt.Println(user.Name)
}

// Safe type assertion
if val, found := c.Get("count"); found {
    if count, ok := val.(int); ok {
        fmt.Println("Count:", count)
    }
}
```

For numeric operations, use the built-in increment/decrement methods which handle types safely.

## Contributing

We welcome contributions! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Add tests for new functionality
4. Ensure tests pass and coverage is maintained (currently **92.9%**):
   ```bash
   go test -race ./...        # Run with race detector
   go test -cover ./...       # Check coverage
   ```
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and release notes.

## Support

- 📖 [Documentation](docs/)
- 💬 [Issues](https://github.com/pzentenoe/go-cache/issues)
- 🌟 [Star on GitHub](https://github.com/pzentenoe/go-cache)

### Buy Me a Coffee

<a href="https://www.buymeacoffee.com/pzentenoe" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

Thank you for your support! ❤️

## Author

**Pablo Zenteno** - [pzentenoe](https://github.com/pzentenoe)
