package cache

import "time"

type shardedJanitor struct {
	Interval time.Duration
	stop     chan bool
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

func stopShardedJanitor(sc *unexportedShardedCache) {
	sc.janitor.stop <- true
}

func runShardedJanitor(sc *shardedCache, ci time.Duration) {
	j := &shardedJanitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	sc.janitor = j
	go j.Run(sc)
}
