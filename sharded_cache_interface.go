package cache

import (
	"io"
	"time"
)

// ShardedCache interface
type ShardedCache interface {
	Set(k string, x any, d time.Duration)
	SetDefault(k string, x any)
	Add(k string, x any, d time.Duration) error
	Replace(k string, x any, d time.Duration) error
	Get(k string) (any, bool)
	GetWithExpiration(k string) (any, time.Time, bool)
	Increment(k string, n int64) error
	IncrementFloat(k string, n float64) error
	Decrement(k string, n int64) error
	DecrementFloat(k string, n float64) error
	Delete(k string)
	DeleteExpired()
	Items() []map[string]Item
	ItemCount() int
	OnEvicted(f func(string, any))
	Flush()
	Save(w io.Writer) error
	SaveFile(fname string) error
	Load(r io.Reader) error
	LoadFile(fname string) error
}

type shardedCache struct {
	seed    uint32
	m       uint32
	cs      []*Cache
	janitor *shardedJanitor
}

type unexportedShardedCache struct {
	*shardedCache
}
