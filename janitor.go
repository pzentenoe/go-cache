package cache

import (
	"runtime"
	"time"
)

type janitor struct {
	Interval       time.Duration
	stop           chan struct{}
	pause          chan struct{}
	resume         chan struct{}
	updateInterval chan time.Duration
}

func (j *janitor) Run(c *Cache) {
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()

	paused := false

	for {
		select {
		case <-ticker.C:
			if !paused {
				c.DeleteExpired()
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

func stopJanitor(c *Cache) {
	close(c.janitor.stop)
}

func runJanitor(c *Cache, ci time.Duration) {
	j := &janitor{
		Interval:       ci,
		stop:           make(chan struct{}),
		pause:          make(chan struct{}),
		resume:         make(chan struct{}),
		updateInterval: make(chan time.Duration),
	}
	c.janitor = j
	go j.Run(c)
	runtime.SetFinalizer(c, stopJanitor)
}
