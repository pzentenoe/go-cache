# go-cache

## Overview

The go-cache library is a high-performance, in-memory key-value store written in Go, designed to provide fast, temporary storage for your applications. It is suitable for caching purposes where data needs to be accessed quickly and frequently, but persistence is not required.

![CI](https://github.com/pzentenoe/go-cache/actions/workflows/actions.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/pzentenoe/go-cache)](https://goreportcard.com/report/github.com/pzentenoe/go-cache)
![License](https://img.shields.io/github/license/pzentenoe/go-cache.svg)
![GitHub issues](https://img.shields.io/github/issues/pzentenoe/go-cache.svg)
![GitHub release](https://img.shields.io/github/v/release/pzentenoe/go-cache.svg)
![GitHub last commit](https://img.shields.io/github/last-commit/pzentenoe/go-cache.svg)
![Quality Gate](https://sonarqube.vikingcode.cl/api/project_badges/measure?project=go-cache&metric=alert_status&token=sqb_308a7e8ba0f56c9872f5db0aa7ec5ed3a491ab1e)
![Coverage](https://sonarqube.vikingcode.cl/api/project_badges/measure?project=go-cache&metric=coverage&token=sqb_308a7e8ba0f56c9872f5db0aa7ec5ed3a491ab1e)
![Bugs](https://sonarqube.vikingcode.cl/api/project_badges/measure?project=go-cache&metric=bugs&token=sqb_308a7e8ba0f56c9872f5db0aa7ec5ed3a491ab1e)

### Buy Me a Coffee

<a href="https://www.buymeacoffee.com/pzentenoe" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

Thank you for your support! ❤️


### Key Features
- Simple and Easy to Use: The library offers a straightforward interface for common cache operations like adding, retrieving, and deleting items.
- Flexible Expiration: Items can have custom expiration times, including no expiration. This is useful for controlling the lifespan of cached data.
- Automatic Expiration: Automatically removes expired items based on a configurable cleanup interval, ensuring efficient memory usage.
- Thread-Safe: The library uses synchronization primitives to ensure safe concurrent access to the cache.
- Support for Various Types: It supports caching items of various types including int, float32, float64, and more.
- Serialization: Provides methods to save and load cache data using Gob encoding.

### Architecture
The go-cache library is structured to offer both a standard cache and a sharded cache for higher concurrency needs.

### Standard Cache
- Cache: The main cache structure that stores items and handles operations like set, get, delete, etc.
- Item: Represents an individual cached item with its value and expiration time.
- Janitor: A background process that periodically cleans up expired items from the cache.

### Sharded Cache
For scenarios requiring high concurrency, the library provides a sharded cache implementation:
- **shardedCache:** Splits the cache into multiple shards, each managed by its own Cache instance to reduce lock contention.
- **shardedJanitor:** A janitor process specific to the sharded cache, responsible for cleaning up expired items in each shard.

### Serialization
The library supports saving and loading cache data to and from files using Gob encoding, allowing the cache state to be persisted and restored

## Installation

To use `go-cache` in your project, install it using the following Go command:

```bash
go get github.com/pzentenoe/go-cache
```
Import go-cache in your project:
```go
import "github.com/pzentenoe/go-cache"
```
## Usage Examples

### Creating and Using a Standard Cache
```go
package main

import (
	"fmt"
	"time"

	"github.com/pzentenoe/go-cache"
)

func main() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)
	// Set the value of the key "key1" to "value1", with the default expiration time
	c.Set("key1", "value1", cache.DefaultExpiration)

	val, found := c.Get("key1")
	if found {
		fmt.Println("Found value:", val)
	} else {
		fmt.Println("Item not found")
	}

	c.SaveFile("cache.data")
	c.LoadFile("cache.data")
}
```
### Creating and Using a Sharded Cache
```go
package main

import (
	"fmt"
	"time"

	"github.com/pzentenoe/go-cache"
)

func main() {
	// Create a sharded cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	sc := cache.NewSharded(5*time.Minute, 10*time.Minute, 10)
	sc.Set("key1", "value1", cache.DefaultExpiration)

	val, found := sc.Get("key1")
	if found {
		fmt.Println("Found value:", val)
	} else {
		fmt.Println("Item not found")
	}

	sc.SaveFile("sharded_cache.data")
	sc.LoadFile("sharded_cache.data")
}
```
### Incrementing and Decrementing Values
```go
package main

import (
	"fmt"
	"time"

	"github.com/pzentenoe/go-cache"
)

func main() {
	c := cache.New(5*time.Minute, 10*time.Minute)
	c.Set("counter", 1, cache.DefaultExpiration)

	c.Increment("counter", 1)
	val, _ := c.Get("counter")
	fmt.Println("Counter after increment:", val)

	c.Decrement("counter", 1)
	val, _ = c.Get("counter")
	fmt.Println("Counter after decrement:", val)
}
```


## Methods

#### Set
```go
Set(k string, x any, d time.Duration)
```
Adds an item to the cache, replacing any existing item. If the duration is DefaultExpiration, the cache’s default expiration time is used. If it is NoExpiration, the item never expires.

#### SetDefault
```go
SetDefault(k string, x any)
```
Adds an item to the cache using the default expiration time.

#### Add
```go
Add(k string, x any, d time.Duration) error
```
Adds an item to the cache only if an item doesn’t already exist for the given key, or if the existing item has expired. Returns an error otherwise.

#### Replace
```go
Replace(k string, x any, d time.Duration) error
```
Sets a new value for the cache key only if it already exists, and the existing item hasn’t expired. Returns an error otherwise.

#### Get
```go
Get(k string) (any, bool)
```
Gets an item from the cache. Returns the item or nil, and a boolean indicating whether the key was found.

#### GetWithExpiration
```go
GetWithExpiration(k string) (any, time.Time, bool)
```
Returns an item and its expiration time from the cache. If the item never expires, a zero value for time.Time is returned.

#### Delete
```go
Delete(k string)
```
Deletes an item from the cache.

#### DeleteExpired
```go
DeleteExpired()
```
Deletes all expired items from the cache.

### OnEvicted
```go
OnEvicted(f func(string, any))
```
Sets a function that is called with the key and value when an item is evicted from the cache. Set to nil to disable.

#### Flush
```go
Flush()
```
Deletes all items from the cache.

#### Increment
```go
Increment(k string, n int64) error
```
Increments an item of type int, int8, int16, int32, int64, uintptr, uint, uint8, uint32, or uint64, float32, or float64 by n. Returns an error if the item’s value is not an integer or if it was not found.

##### IncrementFloat
```go
IncrementFloat(k string, n float64) error
```
Increments an item of type float32 or float64 by n. Returns an error if the item’s value is not floating point, if it was not found, or if it is not possible to increment it by n.

#### Decrement
```go
Decrement(k string, n int64) error
```
Decrements an item of type int, int8, int16, int32, int64, uintptr, uint, uint8, uint32, or uint64, float32, or float64 by n. Returns an error if the item’s value is not an integer or if it was not found.

#### DecrementFloat
```go
DecrementFloat(k string, n float64) error
```
Decrements an item of type float32 or float64 by n. Returns an error if the item’s value is not floating point, if it was not found, or if it is not possible to decrement it by n.

#### Items
```go
Items() map[string]Item
```
Copies all unexpired items in the cache into a new map and returns it.

#### ItemCount
```go
ItemCount() int
```
Returns the number of items in the cache. This may include items that have expired but have not yet been cleaned up.

#### Save
```go
Save(w io.Writer) error
```
Writes the cache’s items (using Gob) to an io.Writer.

#### SaveFile
```go
SaveFile(fname string) error
```
Saves the cache’s items to the given filename, creating the file if it doesn’t exist and overwriting it if it does.

#### Load
```go
Load(r io.Reader) error
```
Adds (Gob-serialized) cache items from an io.Reader, excluding any items with keys that already exist (and haven’t expired) in the current cache.

#### LoadFile
```go
LoadFile(fname string) error
```
Loads and adds cache items from the given filename, excluding any items with keys that already exist in the current cache.


## Testing

Execute the tests with:

```bash
go test ./...
```

## Contributing
We welcome contributions! Please fork the project and submit pull requests to the `main` branch. Make sure to add tests
for new functionalities and document any significant changes.

## License
This project is released under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Changelog
For a detailed changelog, refer to [CHANGELOG.md](CHANGELOG.md).

## Author
- **Pablo Zenteno** - _Full Stack Developer_ - [pzentenoe](https://github.com/pzentenoe)
