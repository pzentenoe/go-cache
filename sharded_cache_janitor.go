package cache

import (
	"runtime"
	"time"
)

type shardedJanitor struct {
	Interval time.Duration
	stop     chan struct{}
}

func (j *shardedJanitor) Run(sc *shardedCache) {
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			sc.DeleteExpired()
		case <-j.stop:
			return
		}
	}
}

// Stop sends a signal to stop the janitor's Run loop and waits for it to finish
func (j *shardedJanitor) Stop() {
	close(j.stop)
}

func stopShardedJanitor(sc *unexportedShardedCache) {
	sc.janitor.Stop()
}

func runShardedJanitor(sc *shardedCache, ci time.Duration) {
	j := &shardedJanitor{
		Interval: ci,
		stop:     make(chan struct{}),
	}
	sc.janitor = j
	go j.Run(sc)
	runtime.SetFinalizer(sc, func(sc *shardedCache) {
		sc.janitor.Stop()
	})
}
