# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
