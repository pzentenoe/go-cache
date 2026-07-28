package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkCacheGetExpiring(b *testing.B) {
	benchmarkCacheGet(b, 5*time.Minute)
}

func BenchmarkCacheGetNotExpiring(b *testing.B) {
	benchmarkCacheGet(b, NoExpiration)
}

func benchmarkCacheGet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	c := New(exp, 0)
	c.Set("foobarba", "zquux", DefaultExpiration)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Get("foobarba")
	}
}

func BenchmarkCacheGetManyConcurrentNotExpiring(b *testing.B) {
	b.StopTimer()
	n := 10000
	c := New(NoExpiration, 0)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := "foo" + strconv.Itoa(i)
		keys[i] = k
		c.Set(k, "bar", NoExpiration)
	}
	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for _, v := range keys {
		go func(k string) {
			for j := 0; j < each; j++ {
				c.Get(k)
			}
			wg.Done()
		}(v)
	}
	b.StartTimer()
	wg.Wait()
}

func BenchmarkCacheSet(b *testing.B) {
	c := New(NoExpiration, 0)
	for i := 0; i < b.N; i++ {
		c.Set("foobarba", "zquux", NoExpiration)
	}
}

func BenchmarkShardedCacheSet(b *testing.B) {
	tc := unexportedNewSharded(NoExpiration, 0, 10)
	for i := 0; i < b.N; i++ {
		tc.Set("foobarba", "zquux", NoExpiration)
	}
}

func BenchmarkCacheIncrementInt64(b *testing.B) {
	c := New(NoExpiration, 0)
	c.Set("counter", int64(0), NoExpiration)
	for i := 0; i < b.N; i++ {
		if _, err := c.IncrementInt64("counter", 1); err != nil {
			b.Fatal(err)
		}
	}
}

// ShardedCache exposes the generic Increment/Decrement, not the typed
// variants, so the sharded benchmarks use Increment for an even comparison.
func BenchmarkShardedCacheIncrement(b *testing.B) {
	tc := unexportedNewSharded(NoExpiration, 0, 10)
	tc.Set("counter", int64(0), NoExpiration)
	for i := 0; i < b.N; i++ {
		if err := tc.Increment("counter", 1); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCacheIncrementManyConcurrent measures lock contention on the
// single-lock cache: every goroutine increments its own key.
func BenchmarkCacheIncrementManyConcurrent(b *testing.B) {
	b.StopTimer()
	n := 100
	c := New(NoExpiration, 0)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := "counter" + strconv.Itoa(i)
		keys[i] = k
		c.Set(k, int64(0), NoExpiration)
	}
	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for _, v := range keys {
		go func(k string) {
			for j := 0; j < each; j++ {
				// Reaching MaxInt64 takes ~2^63 iterations; safe to ignore err.
				_ = c.Increment(k, 1)
			}
			wg.Done()
		}(v)
	}
	b.StartTimer()
	wg.Wait()
}

// BenchmarkShardedCacheIncrementManyConcurrent is the sharded counterpart:
// per-shard locks let increments on different keys proceed in parallel.
func BenchmarkShardedCacheIncrementManyConcurrent(b *testing.B) {
	b.StopTimer()
	n := 100
	tc := unexportedNewSharded(NoExpiration, 0, 10)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := "counter" + strconv.Itoa(i)
		keys[i] = k
		tc.Set(k, int64(0), NoExpiration)
	}
	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for _, v := range keys {
		go func(k string) {
			for j := 0; j < each; j++ {
				_ = tc.Increment(k, 1)
			}
			wg.Done()
		}(v)
	}
	b.StartTimer()
	wg.Wait()
}
