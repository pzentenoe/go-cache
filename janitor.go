package cache

import (
	"runtime"
	"time"
)

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *Cache) {
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			return
		}
	}
}

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

func runJanitor(c *Cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
	runtime.SetFinalizer(c, stopJanitor)
}
