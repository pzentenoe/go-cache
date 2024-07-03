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
		assert.Equal(t, 1*time.Millisecond, sc.janitor.Interval)

		// Allow time for the janitor to run and delete expired items
		time.Sleep(5 * time.Millisecond)

		// Ensure janitor is running by checking if DeleteExpired is called
		sc.Set("key1", "value1", 1*time.Millisecond)
		time.Sleep(2 * time.Millisecond)
		sc.janitor.Run(sc)
		_, found := sc.Get("key1")
		assert.False(t, found)

		stopShardedJanitor(&unexportedShardedCache{sc})
	})

	t.Run("Stop sharded janitor", func(t *testing.T) {
		sc := newShardedCache(2, DefaultExpiration)
		runShardedJanitor(sc, 1*time.Millisecond)

		assert.NotNil(t, sc.janitor)
		assert.Equal(t, 1*time.Millisecond, sc.janitor.Interval)

		// Ensure the janitor stops
		stopShardedJanitor(&unexportedShardedCache{sc})
		select {
		case <-sc.janitor.stop:
			// Success
		case <-time.After(1 * time.Second):
			t.Fatal("Janitor did not stop in time")
		}
	})
}
