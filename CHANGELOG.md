# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2026-02-16

### Changes

- **`map[string]*Item` changed to `map[string]Item`**: All internal maps and public APIs now use value semantics instead of pointers. This affects:
    - `NewFrom()` parameter: `map[string]*Item` → `map[string]Item`
    - `Cache.Items()` return type: `map[string]*Item` → `map[string]Item`
    - `ShardedCache.Items()` return type: `[]map[string]*Item` → `map[string]Item` (now returns a flat unified map instead of per-shard slices)
    - `Item.Expired()` receiver: `*Item` → `Item` (value receiver)

### Added

- **`Close()` method** for `Cache` and `ShardedCache`: Explicitly stops the janitor goroutine and releases resources, eliminating reliance on `runtime.SetFinalizer` for goroutine cleanup
- **`Close()` added to `ShardedCache` interface**: Enables proper resource management through the interface

### Fixed

- **Janitor deadlock risk eliminated**: Replaced unbuffered `pause`/`resume` channels with `sync.Mutex` + boolean flag. `PauseJanitor()` and `ResumeJanitor()` are now idempotent and safe to call multiple times without blocking
- **`Items()` no longer leaks internal state**: With value semantics, callers receive copies of items, preventing external mutation of cache data without locks
- **`GetWithExpiration` duplicated expiration logic**: Now uses `item.Expired()` consistently instead of manually checking `time.Now().UnixNano() > item.Expiration`
- **`delete()` asymmetric control flow**: Simplified to check existence first, then delete, then decide eviction notification

### Refactored

- **`Set`/`set` deduplication**: Exported `Set()` now delegates to unexported `set()` with proper lock management, eliminating duplicated expiration calculation logic
- **`SaveFile` simplified return**: Removed redundant `if err != nil { return err }; return nil` pattern in favor of direct `return c.Save(fp)`
- **Janitor channels buffered**: `stop` and `updateInterval` channels now use buffer of 1 to prevent potential blocking
- **Factory stderr replaced with `log.Println`**: `newShardedCache` no longer writes directly to `os.Stderr`; uses standard `log` package instead

### Performance

- **Reduced GC pressure**: Value-based `Item` storage eliminates per-item heap allocations, significantly reducing garbage collector overhead for caches with many entries
- **Better CPU cache locality**: Items stored as values in the map are more contiguous in memory, improving iteration performance

### Documentation

- Updated `api-reference.md`: Corrected `Items()`, `NewFrom()` signatures, `Expired()` receiver, added `Close()` documentation, noted idempotent behavior of `PauseJanitor`/`ResumeJanitor`
- Updated `sharded-cache.md`: Corrected internal structure diagram (`map[string]Item`, single shared janitor), fixed `Items()` limitation description (now returns flat map)
- Updated `janitor-control.md`: Corrected pause overhead description (mutex-based, not channel-based), updated graceful shutdown pattern to use `Close()`
- Updated `serialization.md`: No API changes needed (uses `Save`/`Load` methods, not internal types)

## [1.2.0] - 2026-02-16

### Refactored

- **Error message constants**: Extracted all repeated error format strings (`errOverflowFormat`, `errUnderflowFormat`, `errItemNotFoundFormat`, etc.) into `constants.go` to eliminate duplication across `increment.go`, `decrement.go`, `cache.go`, `serialization.go`, and `sharded_cache.go`
- **Removed dead parameter**: Eliminated unused `n any` parameter from private `increment()` and `decrement()` helper methods — the value was already captured by the closure and never read inside the function body
- **Consistent item referencing**: Updated cache to use `map[string]*Item` for consistent item referencing and fixed potential underflow/overflow bugs in increment and decrement operations
- **Janitor control enhancements**: Added janitor control improvements for sharded cache
- **Idiomatic error creation**: Replaced `fmt.Errorf()` with `errors.New()` for error messages without format verbs, fixing `staticcheck SA1006` warnings

### Added

- **Pre-commit end-of-file fixer**: Updated pre-commit configuration to include end-of-file fixer hook

### Fixed

- **Underflow/overflow bugs**: Fixed potential underflow and overflow bugs in increment and decrement operations

### Documentation

- Fixed formatting in API reference for ShardedCache section
- Added FOSSA status badges for license and security in README

### CI/CD

- Added FOSSA scan action to workflow for license compliance
- Updated FOSSA action to version 1.7.0

### Tests

- Added tests for cache decrement and increment operations, including underflow and overflow cases
- Updated test coverage command to exclude examples directory

## [1.1.0] - 2025-01-14

### Added

- **Public API `NewSharded()`**: Exported function to create sharded caches as documented in README
- **Complete ShardedCache API**: Added missing methods to ShardedCache interface and implementation:
    - `SetDefault()` - Set items with default expiration
    - `GetWithExpiration()` - Get items with expiration time
    - `DecrementFloat()` - Decrement float values
    - `ItemCount()` - Count total items across all shards
    - `OnEvicted()` - Eviction callbacks for all shards
- **ShardedCache Serialization**: Full support for Save/SaveFile/Load/LoadFile operations
    - Validates shard count on load to prevent data corruption
    - Properly re-hashes items to correct shards during deserialization
- **Overflow/Underflow Protection**: Built-in validation for all increment/decrement operations
    - Prevents silent overflow for uint8, uint16, uint32, uint64, int8, int16, int32, int64
    - Detects infinity for float32 and float64 operations
    - Returns descriptive errors when boundaries would be exceeded
- **Janitor Control Methods**: Added runtime control over automatic cleanup
    - `PauseJanitor()` - Temporarily pause automatic cleanup
    - `ResumeJanitor()` - Resume automatic cleanup after pause
    - `SetJanitorInterval(d time.Duration)` - Dynamically change cleanup interval
- **Comprehensive Test Suite**: Added 700+ lines of new tests
    - Tests for all new ShardedCache methods
    - Overflow/underflow boundary condition tests for all numeric types
    - Concurrency and stress tests for thread safety validation
    - Janitor control methods tests with timing verification
    - Test coverage increased from 80.8% to 92.9%

### Changed

- **Go Version**: Updated from Go 1.22.5 to Go 1.25
- **Dependencies**: Updated `github.com/stretchr/testify` from v1.9.0 to v1.11.1
- **CI/CD Pipeline**: Updated GitHub Actions to latest versions:
    - `actions/setup-go@v3` → `@v5`
    - `actions/cache@v2` → `@v4`
    - `codecov/codecov-action@v4.0.1` → `@v4.6.0`
- **Thread Safety Improvement**: Janitor now uses `chan struct{}` instead of `chan bool` for better memory efficiency
  and Go idioms

### Fixed

- **ShardedCache Load Bug**: Fixed issue where items were not properly re-hashed to correct shards during
  deserialization
- **Error Message Consistency**: Standardized all error messages to use lowercase following Go conventions
    - Updated increment.go, decrement.go, cache.go, and serialization.go
    - Changed "Item %s not found" → "item %s not found"
    - Changed "The value for..." → "the value for..."
    - Changed "Error registering..." → "error registering..."

### Documentation

- Updated README.md with comprehensive "What's New" section
- Added overflow protection examples
- Documented all new ShardedCache methods
- Updated architecture documentation
- Added detailed testing instructions

## [1.0.0] - 2024-07-03

### Added

- Initial release of the `go-cache`.
- Release golang cache library
- Testing.
- golangcilint.yml
- README.md
- Examples
