package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestJanitorPauseResume tests pausing and resuming the janitor
func TestJanitorPauseResume(t *testing.T) {
	t.Run("Pause and resume janitor", func(t *testing.T) {
		// Create cache with short expiration and cleanup interval
		c := New(50*time.Millisecond, 25*time.Millisecond)

		// Add items that will expire
		for i := 0; i < 10; i++ {
			c.Set("key", i, 50*time.Millisecond)
		}

		// Pause janitor
		c.PauseJanitor()

		// Wait for items to expire
		time.Sleep(100 * time.Millisecond)

		// Items should still be in cache (janitor is paused)
		// Note: They are expired but not deleted
		count := c.ItemCount()
		assert.Greater(t, count, 0, "Items should still be in cache while janitor is paused")

		// Resume janitor
		c.ResumeJanitor()

		// Wait for janitor to run
		time.Sleep(50 * time.Millisecond)

		// Now expired items should be cleaned up
		count = c.ItemCount()
		assert.Equal(t, 0, count, "Expired items should be cleaned up after janitor resumes")
	})

	t.Run("Pause janitor on cache without janitor", func(t *testing.T) {
		c := New(NoExpiration, 0) // No janitor
		// Should not panic
		assert.NotPanics(t, func() {
			c.PauseJanitor()
			c.ResumeJanitor()
		})
	})
}

// TestSetJanitorInterval tests dynamically changing the janitor interval
func TestSetJanitorInterval(t *testing.T) {
	t.Run("Update janitor interval", func(t *testing.T) {
		// Start with 100ms interval
		c := New(50*time.Millisecond, 100*time.Millisecond)

		// Add an item that expires in 50ms
		c.Set("key1", "value1", 50*time.Millisecond)

		// Change interval to 10ms (faster cleanup)
		c.SetJanitorInterval(10 * time.Millisecond)

		// Wait for item to expire and be cleaned
		time.Sleep(100 * time.Millisecond)

		// Item should be deleted by now
		_, found := c.Get("key1")
		assert.False(t, found, "Item should be deleted after interval change")
	})

	t.Run("Set interval on cache without janitor", func(t *testing.T) {
		c := New(NoExpiration, 0) // No janitor
		// Should not panic
		assert.NotPanics(t, func() {
			c.SetJanitorInterval(100 * time.Millisecond)
		})
	})
}

// TestJanitorControlCombinations tests various combinations of janitor control
func TestJanitorControlCombinations(t *testing.T) {
	t.Run("Pause, change interval, resume", func(t *testing.T) {
		c := New(50*time.Millisecond, 100*time.Millisecond)

		// Add items
		c.Set("key1", "value1", 50*time.Millisecond)
		c.Set("key2", "value2", 50*time.Millisecond)

		// Pause janitor
		c.PauseJanitor()

		// Wait for expiration
		time.Sleep(100 * time.Millisecond)

		// Change interval to faster cleanup
		c.SetJanitorInterval(10 * time.Millisecond)

		// Resume janitor
		c.ResumeJanitor()

		// Wait for cleanup with new interval
		time.Sleep(50 * time.Millisecond)

		// Items should be cleaned up
		count := c.ItemCount()
		assert.Equal(t, 0, count, "Items should be cleaned up after resume")
	})

	t.Run("Multiple pause and resume cycles", func(t *testing.T) {
		c := New(20*time.Millisecond, 10*time.Millisecond)

		// Add persistent items
		c.Set("persistent1", "value1", NoExpiration)
		c.Set("persistent2", "value2", NoExpiration)

		// First pause-resume cycle
		c.PauseJanitor()
		c.Set("temp1", "value", 20*time.Millisecond)
		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, 3, c.ItemCount(), "All items present while paused")

		c.ResumeJanitor()
		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, 2, c.ItemCount(), "Only persistent items remain")

		// Second pause-resume cycle
		c.PauseJanitor()
		c.Set("temp2", "value", 20*time.Millisecond)
		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, 3, c.ItemCount(), "All items present while paused")

		c.ResumeJanitor()
		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, 2, c.ItemCount(), "Only persistent items remain")
	})
}

// TestJanitorControlConcurrency tests janitor control under concurrent operations
func TestJanitorControlConcurrency(t *testing.T) {
	t.Run("Concurrent pause/resume calls", func(t *testing.T) {
		c := New(50*time.Millisecond, 25*time.Millisecond)

		// Should not panic or deadlock
		assert.NotPanics(t, func() {
			go c.PauseJanitor()
			go c.ResumeJanitor()
			go c.PauseJanitor()
			go c.ResumeJanitor()

			time.Sleep(100 * time.Millisecond)
		})
	})

	t.Run("Concurrent interval updates", func(t *testing.T) {
		c := New(50*time.Millisecond, 100*time.Millisecond)

		// Should not panic or deadlock
		assert.NotPanics(t, func() {
			go c.SetJanitorInterval(10 * time.Millisecond)
			go c.SetJanitorInterval(20 * time.Millisecond)
			go c.SetJanitorInterval(30 * time.Millisecond)

			time.Sleep(100 * time.Millisecond)
		})
	})
}

// TestJanitorIntervalPersistence tests that interval changes persist
func TestJanitorIntervalPersistence(t *testing.T) {
	t.Run("Interval change persists across cleanups", func(t *testing.T) {
		c := New(100*time.Millisecond, 200*time.Millisecond)

		// Change to faster interval
		newInterval := 20 * time.Millisecond
		c.SetJanitorInterval(newInterval)

		// Add multiple expiring items
		for i := 0; i < 3; i++ {
			c.Set("key", i, 30*time.Millisecond)
			// Wait for expiration + enough time for janitor to run
			time.Sleep(80 * time.Millisecond)

			// Each iteration should be cleaned up by janitor with new interval
			count := c.ItemCount()
			assert.Equal(t, 0, count, "Item should be cleaned up with new interval")
		}

		// Verify janitor interval was updated
		assert.NotNil(t, c.janitor)
		assert.Equal(t, newInterval, c.janitor.interval)
	})
}

// TestSetJanitorIntervalNonPositive is a regression test: a non-positive
// interval made the janitor goroutine panic in time.NewTicker, crashing the
// whole process. Non-positive intervals must be ignored.
func TestSetJanitorIntervalNonPositive(t *testing.T) {
	c := New(DefaultExpiration, 50*time.Millisecond)
	defer c.Close()

	assert.NotPanics(t, func() {
		c.SetJanitorInterval(0)
		c.SetJanitorInterval(-time.Second)
	})
	assert.Equal(t, 50*time.Millisecond, c.janitor.interval, "interval must be unchanged")
}

// TestCloseIdempotent verifies Close can be called multiple times and that
// janitor controls become no-ops afterwards.
func TestCloseIdempotent(t *testing.T) {
	c := New(DefaultExpiration, 10*time.Millisecond)
	assert.NotPanics(t, func() {
		c.Close()
		c.Close()
	})
	assert.NotPanics(t, func() {
		c.PauseJanitor()
		c.ResumeJanitor()
		c.SetJanitorInterval(time.Second)
	})
}

// TestStopJanitorAfterClose is a regression test: the GC finalizer calls
// stopJanitor, which dereferenced the janitor field nil'ed by Close, crashing
// the process. It must be a no-op instead.
func TestStopJanitorAfterClose(t *testing.T) {
	c := New(DefaultExpiration, time.Millisecond)
	c.Close()
	assert.NotPanics(t, func() { stopJanitor(c) })
}

// TestJanitorControlConcurrentWithClose exercises the janitor field
// synchronization under -race: controls racing with Close must not panic.
func TestJanitorControlConcurrentWithClose(t *testing.T) {
	c := New(DefaultExpiration, 10*time.Millisecond)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 100; i++ {
			c.PauseJanitor()
			c.ResumeJanitor()
		}
	}()
	c.Close()
	<-done
}
