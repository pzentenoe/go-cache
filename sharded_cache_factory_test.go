package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShardedCache(t *testing.T) {
	t.Run("Create new sharded cache", func(t *testing.T) {
		sc := newShardedCache(2, DefaultExpiration)

		assert.NotNil(t, sc)
		assert.Equal(t, uint32(2), sc.m)
		assert.Len(t, sc.cs, 2)
		assert.NotNil(t, sc.cs[0])
		assert.NotNil(t, sc.cs[1])
	})
}

// TestNewShardedInvalidShards is a regression test: shards <= 0 used to panic
// later with a modulo-by-zero on the first cache operation. It must fail fast
// at construction instead.
func TestNewShardedInvalidShards(t *testing.T) {
	assert.Panics(t, func() { NewSharded(NoExpiration, 0, 0) })
	assert.Panics(t, func() { NewSharded(NoExpiration, 0, -1) })
	assert.NotPanics(t, func() { NewSharded(NoExpiration, 0, 1) })
}
