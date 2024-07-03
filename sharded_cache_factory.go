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
