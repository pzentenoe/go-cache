# Serialization Guide

Persist cache data to disk and restore it later using Go's Gob encoding.

## Overview

go-cache supports serialization of both standard Cache and ShardedCache using Go's encoding/gob package. This allows you
to:

- Persist cache state across application restarts
- Transfer cache data between instances
- Create cache backups
- Pre-populate caches from saved state

## Basic Usage

### Saving Cache to File

```go
c := cache.New(5*time.Minute, 10*time.Minute)

// Populate cache
c.Set("user:1", "Alice", cache.DefaultExpiration)
c.Set("user:2", "Bob", cache.DefaultExpiration)
c.Set("counter", int64(42), cache.NoExpiration)

// Save to file
err := c.SaveFile("cache.gob")
if err != nil {
log.Fatal(err)
}
```

### Loading Cache from File

```go
c := cache.New(5*time.Minute, 10*time.Minute)

// Load from file
err := c.LoadFile("cache.gob")
if err != nil {
log.Fatal(err)
}

// Data is now available
if val, found := c.Get("user:1"); found {
fmt.Println(val) // "Alice"
}
```

## Advanced Usage

### Using io.Writer/io.Reader

For more control, use `Save()` and `Load()` with io interfaces:

```go
// Save to custom writer
file, err := os.Create("cache.gob")
if err != nil {
return err
}
defer file.Close()

err = c.Save(file)

// Save to buffer
var buf bytes.Buffer
err := c.Save(&buf)

// Save to network
conn, _ := net.Dial("tcp", "remote:1234")
err := c.Save(conn)
```

```go
// Load from custom reader
file, err := os.Open("cache.gob")
if err != nil {
return err
}
defer file.Close()

err = c.Load(file)

// Load from buffer
var buf bytes.Buffer
err := c.Load(&buf)

// Load from network
conn, _ := net.Dial("tcp", "remote:1234")
err := c.Load(conn)
```

## Supported Types

Gob encoding supports most Go types:

### Basic Types ✅

```go
c.Set("string", "hello", cache.DefaultExpiration)
c.Set("int", 42, cache.DefaultExpiration)
c.Set("float", 3.14, cache.DefaultExpiration)
c.Set("bool", true, cache.DefaultExpiration)
```

### Composite Types ✅

```go
// Slices
c.Set("slice", []int{1, 2, 3}, cache.DefaultExpiration)

// Maps
c.Set("map", map[string]int{"a": 1, "b": 2}, cache.DefaultExpiration)

// Structs
type User struct {
Name  string
Email string
}
c.Set("user", User{Name: "Alice", Email: "alice@example.com"}, cache.DefaultExpiration)

// Pointers
user := &User{Name: "Bob"}
c.Set("userptr", user, cache.DefaultExpiration)
```

### Custom Types with Gob

Basic struct types work automatically with Gob:

```go
type CustomStruct struct {
    Field1 string
    Field2 int
}

c.Set("custom", CustomStruct{"value", 123}, cache.DefaultExpiration)
c.SaveFile("cache.gob") // Basic structs are handled automatically
```

**Note:** Simple structs with exported fields work without explicit registration. For complex types (interfaces, embedded interfaces, or types requiring special handling), use `gob.Register()` as shown in the Troubleshooting section.

### Unsupported Types ❌

Gob cannot encode:

- Channels
- Functions
- Unsafe pointers
- Complex recursive structures

## Sharded Cache Serialization

ShardedCache serialization works identically to standard Cache:

```go
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 16)

// Populate
for i := 0; i < 1000; i++ {
sc.Set(fmt.Sprintf("key%d", i), i, cache.DefaultExpiration)
}

// Save
err := sc.SaveFile("sharded.gob")

// Load (MUST use same number of shards!)
newSC := cache.NewSharded(5*time.Minute, 10*time.Minute, 16)
err = newSC.LoadFile("sharded.gob")
```

**Important:** ShardedCache must have the same number of shards when loading as when saving, or an error is returned:

```go
sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 16)
sc.SaveFile("cache.gob")

// This will fail!
sc2 := cache.NewSharded(5*time.Minute, 10*time.Minute, 32)
err := sc2.LoadFile("cache.gob")
// Error: shard count mismatch: saved with 16 shards, current cache has 32
```

## Load Behavior

### Merging with Existing Items

`Load()` adds items from the file but does NOT overwrite existing items:

```go
c := cache.New(5*time.Minute, 10*time.Minute)

// Existing item
c.Set("key1", "existing", cache.DefaultExpiration)

// File contains: key1="from_file", key2="value2"
c.LoadFile("cache.gob")

// Existing item is preserved
val, _ := c.Get("key1")
fmt.Println(val) // "existing" (NOT "from_file")

// New items are added
val, _ = c.Get("key2")
fmt.Println(val) // "value2"
```

### Loading into Empty Cache

To replace all data, flush before loading:

```go
c.Flush()
c.LoadFile("cache.gob")
```

## Complete Example

```go
package main

import (
	"fmt"
	"time"
	"github.com/pzentenoe/go-cache"
)

type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

func main() {
	// Create and populate cache
	c := cache.New(5*time.Minute, 10*time.Minute)

	users := []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()},
		{ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Now()},
	}

	for _, user := range users {
		c.Set(fmt.Sprintf("user:%d", user.ID), user, cache.DefaultExpiration)
	}

	c.Set("config", map[string]interface{}{
		"max_connections": 100,
		"timeout":         30 * time.Second,
	}, cache.NoExpiration)

	// Save to file
	fmt.Println("Saving cache...")
	err := c.SaveFile("app_cache.gob")
	if err != nil {
		panic(err)
	}

	// Simulate app restart - create new cache
	fmt.Println("\nSimulating app restart...")
	newCache := cache.New(5*time.Minute, 10*time.Minute)

	// Load from file
	fmt.Println("Loading cache from disk...")
	err = newCache.LoadFile("app_cache.gob")
	if err != nil {
		panic(err)
	}

	// Verify data
	if val, found := newCache.Get("user:1"); found {
		user := val.(User)
		fmt.Printf("Loaded user: %+v\n", user)
	}

	if val, found := newCache.Get("config"); found {
		config := val.(map[string]interface{})
		fmt.Printf("Loaded config: %v\n", config)
	}
}
```

## Error Handling

### Save Errors

```go
err := c.SaveFile("cache.gob")
if err != nil {
// Possible errors:
// - File permission denied
// - Disk full
// - Type registration error
log.Printf("Failed to save cache: %v", err)
}
```

### Load Errors

```go
err := c.LoadFile("cache.gob")
if err != nil {
// Possible errors:
// - File not found
// - Corrupt file
// - Incompatible Gob version
// - Shard count mismatch (ShardedCache)
log.Printf("Failed to load cache: %v", err)
}
```

## Best Practices

### 1. Handle Load Errors Gracefully

```go
c := cache.New(5*time.Minute, 10*time.Minute)

err := c.LoadFile("cache.gob")
if err != nil {
log.Printf("Warning: Could not load cache: %v", err)
// Continue with empty cache
} else {
log.Println("Cache loaded successfully")
}
```

### 2. Save Periodically

```go
// Save cache every 5 minutes
ticker := time.NewTicker(5 * time.Minute)
go func () {
for range ticker.C {
err := c.SaveFile("cache.gob")
if err != nil {
log.Printf("Failed to save cache: %v", err)
}
}
}()
```

### 3. Save on Shutdown

```go
func main() {
c := cache.New(5*time.Minute, 10*time.Minute)

// Handle graceful shutdown
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

go func () {
<-sigChan
fmt.Println("\nSaving cache before shutdown...")
c.SaveFile("cache.gob")
os.Exit(0)
}()

// Your application logic...
}
```

### 4. Use Atomic Writes

```go
// Write to temp file first, then rename
tmpFile := "cache.gob.tmp"
finalFile := "cache.gob"

err := c.SaveFile(tmpFile)
if err != nil {
return err
}

// Atomic rename
err = os.Rename(tmpFile, finalFile)
if err != nil {
os.Remove(tmpFile)
return err
}
```

### 5. Compress Large Caches

```go
import (
"compress/gzip"
"os"
)

// Save with compression
file, _ := os.Create("cache.gob.gz")
defer file.Close()

gzWriter := gzip.NewWriter(file)
defer gzWriter.Close()

c.Save(gzWriter)

// Load with decompression
file, _ := os.Open("cache.gob.gz")
defer file.Close()

gzReader, _ := gzip.NewReader(file)
defer gzReader.Close()

c.Load(gzReader)
```

## Performance Considerations

### Save Performance

Save time depends on:

- Number of items
- Complexity of stored types
- Disk I/O speed

Approximate times (SSD):

- 1,000 items: ~5ms
- 10,000 items: ~50ms
- 100,000 items: ~500ms

### Load Performance

Load time is similar to save time but may be slightly faster.

### Memory Usage

During serialization:

- Save: No additional memory (streams to disk)
- Load: Temporary memory for deserialization

## Troubleshooting

### "gob: type not registered" or encoding errors

For complex types (interfaces, function types, or types with unexported fields), register them explicitly:

```go
type MyComplexType struct { ... }

func init() {
    gob.Register(MyComplexType{})
    gob.Register(&MyComplexType{}) // If storing pointers
}
```

**When to use `gob.Register()`:**
- ✅ Storing interface{} values containing custom types
- ✅ Types with embedded interfaces
- ✅ Pointer types
- ✅ When getting "gob: type not registered" errors

**Not needed for:**
- ❌ Simple structs with exported fields
- ❌ Built-in types (string, int, bool, etc.)

### "shard count mismatch"

```go
// Ensure same shard count
const SHARD_COUNT = 16

sc1 := cache.NewSharded(exp, cleanup, SHARD_COUNT)
sc1.SaveFile("cache.gob")

sc2 := cache.NewSharded(exp, cleanup, SHARD_COUNT) // Use constant
sc2.LoadFile("cache.gob")
```

### Corrupt File

```go
// Validate before loading
fileInfo, err := os.Stat("cache.gob")
if err != nil {
log.Println("Cache file not found, starting fresh")
return
}

if fileInfo.Size() == 0 {
log.Println("Cache file is empty, starting fresh")
return
}

c.LoadFile("cache.gob")
```

## See Also

- [API Reference](api-reference.md)
- [Serialization Example](../examples/serialization/) - Hands-on demonstration
- [Getting Started Guide](getting-started.md) - Core concepts
- [Sharded Cache Guide](sharded-cache.md) - Serializing sharded caches
- [Go Gob Documentation](https://pkg.go.dev/encoding/gob)
