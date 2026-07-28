package cache

import (
	"time"
)

func stopShardedJanitor(sc *unexportedShardedCache) {
	sc.mu.Lock()
	j := sc.janitor
	sc.mu.Unlock()
	if j != nil {
		j.stop <- struct{}{}
	}
}

func runShardedJanitor(sc *shardedCache, ci time.Duration) {
	j := newJanitor(ci, sc.DeleteExpired)
	sc.mu.Lock()
	sc.janitor = j
	sc.mu.Unlock()
	go j.Run()
}
