package cache

import "time"

func (sc *shardedCache) bucket(k string) *Cache {
	return sc.cs[djb33(sc.seed, k)%sc.m]
}

func (sc *shardedCache) Set(k string, x any, d time.Duration) {
	sc.bucket(k).Set(k, x, d)
}

func (sc *shardedCache) Add(k string, x any, d time.Duration) error {
	return sc.bucket(k).Add(k, x, d)
}

func (sc *shardedCache) Replace(k string, x any, d time.Duration) error {
	return sc.bucket(k).Replace(k, x, d)
}

func (sc *shardedCache) Get(k string) (any, bool) {
	return sc.bucket(k).Get(k)
}

func (sc *shardedCache) Increment(k string, n int64) error {
	return sc.bucket(k).Increment(k, n)
}

func (sc *shardedCache) IncrementFloat(k string, n float64) error {
	return sc.bucket(k).IncrementFloat(k, n)
}

func (sc *shardedCache) Decrement(k string, n int64) error {
	return sc.bucket(k).Decrement(k, n)
}

func (sc *shardedCache) Delete(k string) {
	sc.bucket(k).Delete(k)
}

func (sc *shardedCache) DeleteExpired() {
	for _, v := range sc.cs {
		v.DeleteExpired()
	}
}

func (sc *shardedCache) Items() []map[string]Item {
	res := make([]map[string]Item, len(sc.cs))
	for i, v := range sc.cs {
		res[i] = v.Items()
	}
	return res
}

func (sc *shardedCache) Flush() {
	for _, v := range sc.cs {
		v.Flush()
	}
}
