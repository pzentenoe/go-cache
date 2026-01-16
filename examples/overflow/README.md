# Overflow Protection Example

This example demonstrates built-in overflow/underflow protection for numeric operations.

## What You'll Learn

- Unsigned integer overflow detection
- Signed integer underflow detection
- Float infinity protection
- Safe increment/decrement operations
- Error handling for boundary conditions

## Running This Example

```bash
cd examples/overflow
go run main.go
```

## Operations Demonstrated

### 1. Unsigned Integer Overflow

Protected types:
- `uint8`, `uint16`, `uint32`, `uint64`
- `IncrementUint8()`, `IncrementUint16()`, etc.

```go
// Detect overflow before it happens
result, err := c.IncrementUint8("counter", 10)
if err != nil {
    // Would exceed math.MaxUint8 (255)
}
```

### 2. Signed Integer Overflow/Underflow

Protected types:
- `int8`, `int16`, `int32`, `int64`
- `IncrementInt8()`, `DecrementInt32()`, etc.

```go
// Detect underflow
result, err := c.DecrementInt8("counter", 200)
if err != nil {
    // Would go below math.MinInt8 (-128)
}
```

### 3. Float Infinity Detection

Protected types:
- `float32`, `float64`
- `IncrementFloat32()`, `IncrementFloat()` (float64)

```go
// Prevent infinity
result, err := c.IncrementFloat32("price", math.MaxFloat32)
if err != nil {
    // Would result in +Inf
}
```

## Why This Matters

### Without Protection

```go
var counter uint8 = 255
counter += 1  // Wraps to 0 (silent overflow!)
```

### With Protection

```go
c.Set("counter", uint8(255), cache.NoExpiration)
result, err := c.IncrementUint8("counter", 1)
if err != nil {
    fmt.Println("Overflow prevented:", err)
}
```

## All Protected Operations

| Type      | Increment              | Decrement              |
|-----------|------------------------|------------------------|
| `uint8`   | `IncrementUint8()`    | `DecrementUint8()`    |
| `uint16`  | `IncrementUint16()`   | `DecrementUint16()`   |
| `uint32`  | `IncrementUint32()`   | `DecrementUint32()`   |
| `uint64`  | `IncrementUint64()`   | `DecrementUint64()`   |
| `int8`    | `IncrementInt8()`     | `DecrementInt8()`     |
| `int16`   | `IncrementInt16()`    | `DecrementInt16()`    |
| `int32`   | `IncrementInt32()`    | `DecrementInt32()`    |
| `int64`   | `IncrementInt64()`    | `DecrementInt64()`    |
| `float32` | `IncrementFloat32()`  | `DecrementFloat32()`  |
| `float64` | `IncrementFloat()`    | `DecrementFloat()`    |

## Use Cases

- **Rate Limiters**: Prevent counter wrap-around
- **Financial Calculations**: Avoid silent precision loss
- **Statistics**: Maintain data integrity
- **Metrics**: Reliable counter increments

## Related Documentation

- [API Reference](../../docs/api-reference.md)
- [Getting Started](../../docs/getting-started.md)

## Next Steps

- [Concurrent Example](../concurrent/) - Thread-safe increments
- [Janitor Example](../janitor/) - Optimize cleanup for counters
- [Basic Example](../basic/) - Learn core operations
