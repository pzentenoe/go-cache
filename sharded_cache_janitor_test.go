package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShardedJanitor(t *testing.T) {
	t.Run("Run sharded janitor", func(t *testing.T) {
		sc := newShardedCache(2, DefaultExpiration)
		runShardedJanitor(sc, 1*time.Millisecond)

		assert.NotNil(t, sc.janitor)
		assert.Equal(t, 1*time.Millisecond, sc.janitor.interval)

		sc.Set("key1", "value1", 1*time.Millisecond)
		time.Sleep(2 * time.Millisecond)

		time.Sleep(5 * time.Millisecond)

		_, found := sc.Get("key1")
		assert.False(t, found)

		stopShardedJanitor(&unexportedShardedCache{sc})
	})

	t.Run("Stop sharded janitor", func(t *testing.T) {
		sc := newShardedCache(2, DefaultExpiration)
		runShardedJanitor(sc, 1*time.Millisecond)

		assert.NotNil(t, sc.janitor)
		assert.Equal(t, 1*time.Millisecond, sc.janitor.interval)

		stopShardedJanitor(&unexportedShardedCache{sc})

		// Wait enough for the goroutine to process the stop signal
		time.Sleep(10 * time.Millisecond)

		// Janitor goroutine should have exited; adding an item with
		// short expiration should NOT be cleaned up automatically
		sc.Set("key1", "value1", 1*time.Millisecond)
		time.Sleep(10 * time.Millisecond)
		// Item count includes expired items since janitor is stopped
		assert.Equal(t, 1, sc.ItemCount())
	})
}

// TestShardedCloseIdempotent verifies Close can be called multiple times and
// that janitor controls become no-ops afterwards.
func TestShardedCloseIdempotent(t *testing.T) {
	sc := NewSharded(NoExpiration, 10*time.Millisecond, 2)
	assert.NotPanics(t, func() {
		sc.Close()
		sc.Close()
	})
	assert.NotPanics(t, func() {
		sc.PauseJanitor()
		sc.ResumeJanitor()
		sc.SetJanitorInterval(time.Second)
	})
}

// TestShardedSetJanitorIntervalNonPositive is a regression test: a
// non-positive interval panicked the janitor goroutine via time.NewTicker.
// Non-positive intervals must be ignored.
func TestShardedSetJanitorIntervalNonPositive(t *testing.T) {
	sc := NewSharded(NoExpiration, 50*time.Millisecond, 2)
	defer sc.Close()
	assert.NotPanics(t, func() {
		sc.SetJanitorInterval(0)
		sc.SetJanitorInterval(-time.Second)
	})
}

// TestStopShardedJanitorAfterClose is a regression test: the GC finalizer
// calls stopShardedJanitor, which dereferenced the janitor field nil'ed by
// Close. It must be a no-op instead.
func TestStopShardedJanitorAfterClose(t *testing.T) {
	sc := NewSharded(NoExpiration, time.Millisecond, 2).(*unexportedShardedCache)
	sc.Close()
	assert.NotPanics(t, func() { stopShardedJanitor(sc) })
}

// TestShardedJanitorControlConcurrentWithClose exercises the janitor field
// synchronization under -race: controls racing with Close must not panic.
func TestShardedJanitorControlConcurrentWithClose(t *testing.T) {
	sc := NewSharded(NoExpiration, 10*time.Millisecond, 2)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 100; i++ {
			sc.PauseJanitor()
			sc.ResumeJanitor()
		}
	}()
	sc.Close()
	<-done
}

// TestShardedSetJanitorIntervalConcurrentWithClose is a regression test: an
// interval update racing Close could block forever once the janitor exited.
func TestShardedSetJanitorIntervalConcurrentWithClose(t *testing.T) {
	for i := 0; i < 100; i++ {
		sc := NewSharded(NoExpiration, time.Millisecond, 2)
		done := make(chan struct{})
		go func() {
			defer close(done)
			for j := 0; j < 50; j++ {
				sc.SetJanitorInterval(time.Duration(j+1) * time.Millisecond)
			}
		}()
		sc.Close()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("SetJanitorInterval blocked after the janitor stopped")
		}
	}
}
