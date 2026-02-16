package cache

import (
	"time"
)

type shardedJanitor struct {
	Interval       time.Duration
	stop           chan struct{}
	pause          chan struct{}
	resume         chan struct{}
	updateInterval chan time.Duration
}

func (j *shardedJanitor) Run(sc *shardedCache) {
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()

	paused := false

	for {
		select {
		case <-ticker.C:
			if !paused {
				sc.DeleteExpired()
			}
		case <-j.pause:
			paused = true
		case <-j.resume:
			paused = false
		case newInterval := <-j.updateInterval:
			ticker.Stop()
			j.Interval = newInterval
			ticker = time.NewTicker(newInterval)
		case <-j.stop:
			return
		}
	}
}

// Stop sends a signal to stop the janitor's Run loop
func (j *shardedJanitor) Stop() {
	close(j.stop)
}

func stopShardedJanitor(sc *unexportedShardedCache) {
	sc.janitor.Stop()
}

func runShardedJanitor(sc *shardedCache, ci time.Duration) {
	j := &shardedJanitor{
		Interval:       ci,
		stop:           make(chan struct{}),
		pause:          make(chan struct{}),
		resume:         make(chan struct{}),
		updateInterval: make(chan time.Duration),
	}
	sc.janitor = j
	go j.Run(sc)
}
