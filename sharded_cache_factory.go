package cache

import (
	"crypto/rand"
	"math"
	"math/big"
	insecurerand "math/rand"
	"os"
	"runtime"
	"time"
)

func newShardedCache(n int, de time.Duration) *shardedCache {
	max := big.NewInt(0).SetUint64(uint64(math.MaxUint32))
	rnd, err := rand.Int(rand.Reader, max)
	var seed uint32
	if err != nil {
		os.Stderr.Write([]byte("WARNING: go-cache's newShardedCache failed to read from the system CSPRNG (/dev/urandom or equivalent.) Your system's security may be compromised. Continuing with an insecure seed.\n"))
		seed = insecurerand.Uint32()
	} else {
		seed = uint32(rnd.Uint64())
	}
	sc := &shardedCache{
		seed: seed,
		m:    uint32(n),
		cs:   make([]*Cache, n),
	}
	for i := 0; i < n; i++ {
		c := &Cache{
			defaultExpiration: de,
			items:             make(map[string]Item),
		}
		sc.cs[i] = c
	}
	return sc
}

// NewSharded returns a new sharded cache with a given default expiration
// duration and cleanup interval. The shards parameter determines the number
// of internal shards to use for the cache. More shards reduce lock contention
// in high-concurrency scenarios but use more memory. A good starting point is
// between 8-32 shards depending on your use case.
//
// If the expiration duration is less than one (or NoExpiration), the items
// in the cache never expire (by default), and must be deleted manually. If
// the cleanup interval is less than one, expired items are not deleted from
// the cache before calling DeleteExpired().
func NewSharded(defaultExpiration, cleanupInterval time.Duration, shards int) ShardedCache {
	return unexportedNewSharded(defaultExpiration, cleanupInterval, shards)
}

func unexportedNewSharded(defaultExpiration, cleanupInterval time.Duration, shards int) *unexportedShardedCache {
	if defaultExpiration == 0 {
		defaultExpiration = NoExpiration
	}
	sc := newShardedCache(shards, defaultExpiration)
	SC := &unexportedShardedCache{sc}
	if cleanupInterval > 0 {
		runShardedJanitor(sc, cleanupInterval)
		runtime.SetFinalizer(SC, stopShardedJanitor)
	}
	return SC
}
