package cache

import (
	"sync"
	"time"
)

type shardedJanitor struct {
	interval       time.Duration
	stop           chan struct{}
	updateInterval chan time.Duration
	mu             sync.Mutex
	paused         bool
}

func (j *shardedJanitor) Run(sc *shardedCache) {
	j.mu.Lock()
	ticker := time.NewTicker(j.interval)
	j.mu.Unlock()
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			j.mu.Lock()
			paused := j.paused
			j.mu.Unlock()
			if !paused {
				sc.DeleteExpired()
			}
		case newInterval := <-j.updateInterval:
			ticker.Stop()
			j.mu.Lock()
			j.interval = newInterval
			j.mu.Unlock()
			ticker = time.NewTicker(newInterval)
		case <-j.stop:
			return
		}
	}
}

// Pause temporarily pauses the sharded janitor.
// Safe to call multiple times; subsequent calls are no-ops.
func (j *shardedJanitor) Pause() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.paused = true
}

// Resume resumes the sharded janitor after a pause.
// Safe to call multiple times; subsequent calls are no-ops.
func (j *shardedJanitor) Resume() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.paused = false
}

func stopShardedJanitor(sc *unexportedShardedCache) {
	sc.mu.Lock()
	j := sc.janitor
	sc.mu.Unlock()
	if j != nil {
		j.stop <- struct{}{}
	}
}

func runShardedJanitor(sc *shardedCache, ci time.Duration) {
	j := &shardedJanitor{
		interval:       ci,
		stop:           make(chan struct{}, 1),
		updateInterval: make(chan time.Duration, 1),
	}
	sc.mu.Lock()
	sc.janitor = j
	sc.mu.Unlock()
	go j.Run(sc)
}
