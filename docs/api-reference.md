# API Reference

Complete reference for all go-cache methods and types.

## Cache Creation

### `New(defaultExpiration, cleanupInterval time.Duration) *Cache`

Creates a new cache with the specified default expiration and cleanup interval.

**Parameters:**

- `defaultExpiration`: Default expiration time for items (use `NoExpiration` for no expiration)
- `cleanupInterval`: How often to run the janitor cleanup (use `0` to disable automatic cleanup)

**Returns:** Pointer to a new Cache instance

**Example:**

```go
c := cache.New(5*time.Minute, 10*time.Minute)
```

### `NewFrom(defaultExpiration, cleanupInterval time.Duration, items map[string]Item) *Cache`

Creates a cache from an existing items map.

**Parameters:**

- `defaultExpiration`: Default expiration time
- `cleanupInterval`: Cleanup interval
- `items`: Pre-existing items map

**Returns:** Pointer to a new Cache instance

---

## Set Operations

### `Set(k string, x any, d time.Duration)`

Sets an item in the cache, replacing any existing item.

**Parameters:**

- `k`: Key
- `x`: Value (any type)
- `d`: Expiration duration (`DefaultExpiration`, `NoExpiration`, or custom duration)

**Example:**

```go
c.Set("user:1", "Alice", cache.DefaultExpiration)
c.Set("session", "abc123", 1*time.Hour)
c.Set("config", cfg, cache.NoExpiration)
```

### `SetDefault(k string, x any)`

Sets an item using the cache's default expiration time.

**Parameters:**

- `k`: Key
- `x`: Value

**Example:**

```go
c.SetDefault("key", "value")
// Equivalent to: c.Set("key", "value", cache.DefaultExpiration)
```

### `Add(k string, x any, d time.Duration) error`

Adds an item only if it doesn't already exist or has expired.

**Parameters:**

- `k`: Key
- `x`: Value
- `d`: Expiration duration

**Returns:** Error if item already exists

**Example:**

```go
err := c.Add("user:1", "Alice", cache.DefaultExpiration)
if err != nil {
// Item already exists
}
```

### `Replace(k string, x any, d time.Duration) error`

Replaces an existing item's value. Returns error if item doesn't exist.

**Parameters:**

- `k`: Key
- `x`: New value
- `d`: New expiration duration

**Returns:** Error if item doesn't exist

**Example:**

```go
err := c.Replace("user:1", "Bob", cache.DefaultExpiration)
if err != nil {
// Item doesn't exist
}
```

---

## Get Operations

### `Get(k string) (any, bool)`

Gets an item from the cache.

**Parameters:**

- `k`: Key

**Returns:**

- `any`: The value (or nil if not found)
- `bool`: Whether the key was found

**Example:**

```go
if val, found := c.Get("user:1"); found {
user := val.(string)
fmt.Println(user)
}
```

### `GetWithExpiration(k string) (any, time.Time, bool)`

Gets an item and its expiration time.

**Parameters:**

- `k`: Key

**Returns:**

- `any`: The value
- `time.Time`: Expiration time (zero if no expiration)
- `bool`: Whether the key was found

**Example:**

```go
if val, expTime, found := c.GetWithExpiration("session"); found {
if !expTime.IsZero() {
fmt.Println("Expires at:", expTime)
fmt.Println("TTL:", time.Until(expTime))
}
}
```

---

## Delete Operations

### `Delete(k string)`

Deletes an item from the cache. Triggers `OnEvicted` callback if set.

**Parameters:**

- `k`: Key

**Example:**

```go
c.Delete("user:1")
```

### `DeleteExpired()`

Manually deletes all expired items. Triggers `OnEvicted` for each deleted item.

**Example:**

```go
c.DeleteExpired()
```

### `Flush()`

Deletes all items from the cache. Does NOT trigger `OnEvicted` callbacks.

**Example:**

```go
c.Flush()
```

---

## Increment/Decrement Operations

### Integer Operations

#### `Increment(k string, n int64) error`

Increments a numeric item by n. Works with int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32,
float64.

**Returns:** Error if item is not numeric or doesn't exist

```go
c.Set("counter", int64(10), cache.NoExpiration)
c.Increment("counter", 5) // counter = 15
```

#### `Decrement(k string, n int64) error`

Decrements a numeric item by n.

```go
c.Decrement("counter", 3) // counter = 12
```

### Typed Integer Operations (with Overflow Protection)

All typed operations include overflow/underflow protection:

- `IncrementInt(k string, n int) (int, error)`
- `IncrementInt8(k string, n int8) (int8, error)`
- `IncrementInt16(k string, n int16) (int16, error)`
- `IncrementInt32(k string, n int32) (int32, error)`
- `IncrementInt64(k string, n int64) (int64, error)`
- `IncrementUint(k string, n uint) (uint, error)`
- `IncrementUintptr(k string, n uintptr) (uintptr, error)`
- `IncrementUint8(k string, n uint8) (uint8, error)`
- `IncrementUint16(k string, n uint16) (uint16, error)`
- `IncrementUint32(k string, n uint32) (uint32, error)`
- `IncrementUint64(k string, n uint64) (uint64, error)`

Corresponding `Decrement*` methods exist for each type.

**Example:**

```go
c.Set("views", uint64(100), cache.NoExpiration)
result, err := c.IncrementUint64("views", 50)
if err != nil {
fmt.Println("Overflow would occur:", err)
} else {
fmt.Println("New value:", result) // 150
}
```

### Float Operations

#### `IncrementFloat(k string, n float64) error`

Increments a float32 or float64 value.

```go
c.Set("price", 19.99, cache.NoExpiration)
c.IncrementFloat("price", 5.00) // 24.99
```

#### `DecrementFloat(k string, n float64) error`

Decrements a float32 or float64 value.

```go
c.DecrementFloat("price", 2.50) // 22.49
```

#### Typed Float Operations (with Infinity Protection)

- `IncrementFloat32(k string, n float32) (float32, error)`
- `IncrementFloat64(k string, n float64) (float64, error)`
- `DecrementFloat32(k string, n float32) (float32, error)`
- `DecrementFloat64(k string, n float64) (float64, error)`

---

## Utility Methods

### `Items() map[string]Item`

Returns a copy of all unexpired items.

**Returns:** Map of all items in cache

**Example:**

```go
items := c.Items()
for key, item := range items {
fmt.Printf("%s: %v\n", key, item.Object)
}
```

### `ItemCount() int`

Returns the number of items in the cache (may include expired items not yet cleaned).

**Returns:** Item count

**Example:**

```go
count := c.ItemCount()
fmt.Printf("Cache has %d items\n", count)
```

### `OnEvicted(f func(string, any))`

Sets a callback function called when items are evicted. Set to `nil` to disable.

**Parameters:**

- `f`: Callback function receiving key and value

**Triggered by:**

- `Delete()`
- `DeleteExpired()` (for expired items)
- Manual eviction

**NOT triggered by:**

- `Set()` overwrites
- `Flush()`

**Example:**

```go
c.OnEvicted(func (key string, value any) {
log.Printf("Evicted: %s = %v", key, value)
})
```

---

## Janitor Control Methods

### `PauseJanitor()`

Temporarily pauses automatic cleanup of expired items.

**Example:**

```go
c.PauseJanitor()
// Perform bulk operations
c.ResumeJanitor()
```

### `ResumeJanitor()`

Resumes automatic cleanup after pause.

### `SetJanitorInterval(d time.Duration)`

Dynamically changes the cleanup interval.

**Parameters:**

- `d`: New cleanup interval

**Example:**

```go
c.SetJanitorInterval(5 * time.Minute)
```

---

## Serialization Methods

### `Save(w io.Writer) error`

Writes cache items to an io.Writer using Gob encoding.

**Parameters:**

- `w`: io.Writer to write to

**Returns:** Error if write fails

**Example:**

```go
file, _ := os.Create("cache.gob")
defer file.Close()
err := c.Save(file)
```

### `SaveFile(fname string) error`

Saves cache to a file.

**Parameters:**

- `fname`: File path

**Returns:** Error if save fails

**Example:**

```go
err := c.SaveFile("cache.gob")
```

### `Load(r io.Reader) error`

Loads cache items from an io.Reader (Gob encoded).

**Parameters:**

- `r`: io.Reader to read from

**Returns:** Error if read fails

**Example:**

```go
file, _ := os.Open("cache.gob")
defer file.Close()
err := c.Load(file)
```

### `LoadFile(fname string) error`

Loads cache from a file.

**Parameters:**

- `fname`: File path

**Returns:** Error if load fails

**Example:**

```go
err := c.LoadFile("cache.gob")
```

---

## Constants

### `NoExpiration`

Used to specify that an item should never expire.

```go
c.Set("config", cfg, cache.NoExpiration)
```

### `DefaultExpiration`

Uses the cache's default expiration time.

```go
c.Set("key", "value", cache.DefaultExpiration)
```

---

## Types

### `Item`

Represents a cached item.

```go
type Item struct {
Object     any   // The cached value
Expiration int64 // Expiration time in UnixNano
}
```

#### `(i Item) Expired() bool`

Returns true if the item has expired.

```go
if item.Expired() {
fmt.Println("Item has expired")
}
```

---

## ShardedCache

See [Sharded Cache Documentation](sharded-cache.md) for ShardedCache-specific API.