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
		assert.Equal(t, 1*time.Millisecond, sc.janitor.Interval)

		stopShardedJanitor(&unexportedShardedCache{sc})

		time.Sleep(10 * time.Millisecond)

		select {
		case <-sc.janitor.stop:
		default:
			t.Fatal("Janitor did not stop in time")
		}
	})
}
