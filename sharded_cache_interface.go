package cache

import "time"

type ShardedCache interface {
	Set(k string, x any, d time.Duration)
	Add(k string, x any, d time.Duration) error
	Replace(k string, x any, d time.Duration) error
	Get(k string) (any, bool)
	Increment(k string, n int64) error
	IncrementFloat(k string, n float64) error
	Decrement(k string, n int64) error
	Delete(k string)
	DeleteExpired()
	Items() []map[string]Item
	Flush()
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
