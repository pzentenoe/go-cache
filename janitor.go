package cache

import (
	"runtime"
	"sync"
	"time"
)

type janitor struct {
	interval       time.Duration
	stop           chan struct{}
	updateInterval chan time.Duration
	mu             sync.Mutex
	paused         bool
}

func (j *janitor) Run(c *Cache) {
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
				c.DeleteExpired()
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
	j := &janitor{
		interval:       ci,
		stop:           make(chan struct{}, 1),
		updateInterval: make(chan time.Duration, 1),
	}
	c.mu.Lock()
	c.janitor = j
	c.mu.Unlock()
	go j.Run(c)
	runtime.SetFinalizer(c, stopJanitor)
}
