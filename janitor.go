package cache

import (
	"runtime"
	"sync"
	"time"
)

// janitor periodically deletes expired items. A single implementation serves
// both Cache and shardedCache via the cleanup function.
type janitor struct {
	interval       time.Duration
	stop           chan struct{}
	updateInterval chan time.Duration
	cleanup        func()
	mu             sync.Mutex
	paused         bool
}

func newJanitor(ci time.Duration, cleanup func()) *janitor {
	return &janitor{
		interval:       ci,
		stop:           make(chan struct{}, 1),
		updateInterval: make(chan time.Duration, 1),
		cleanup:        cleanup,
	}
}

func (j *janitor) Run() {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			j.mu.Lock()
			paused := j.paused
			j.mu.Unlock()
			if !paused {
				j.cleanup()
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

// Pause temporarily pauses the janitor.
// Safe to call multiple times; subsequent calls are no-ops.
func (j *janitor) Pause() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.paused = true
}

// Resume resumes the janitor after a pause.
// Safe to call multiple times; subsequent calls are no-ops.
func (j *janitor) Resume() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.paused = false
}

func stopJanitor(c *Cache) {
	c.mu.RLock()
	j := c.janitor
	c.mu.RUnlock()
	if j != nil {
		j.stop <- struct{}{}
	}
}

func runJanitor(c *Cache, ci time.Duration) {
	j := newJanitor(ci, c.DeleteExpired)
	c.mu.Lock()
	c.janitor = j
	c.mu.Unlock()
	go j.Run()
	runtime.SetFinalizer(c, stopJanitor)
}
