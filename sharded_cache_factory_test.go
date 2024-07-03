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
