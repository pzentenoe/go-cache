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
