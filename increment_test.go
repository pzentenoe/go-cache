package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache_Increment(t *testing.T) {
	t.Run("Increment int", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", 10, DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, 15, val)
	})

	t.Run("Increment int8", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int8(10), DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, int8(15), val)
	})

	t.Run("Increment int16", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int16(10), DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, int16(15), val)
	})

	t.Run("Increment int32", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int32(10), DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, int32(15), val)
	})

	t.Run("Increment int64", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int64(10), DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, int64(15), val)
	})

	t.Run("Increment uint", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint(10), DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, uint(15), val)
	})

	t.Run("Increment uintptr", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uintptr(10), DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, uintptr(15), val)
	})

	t.Run("Increment uint8", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint8(10), DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, uint8(15), val)
	})

	t.Run("Increment uint16", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint16(10), DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, uint16(15), val)
	})

	t.Run("Increment uint32", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint32(10), DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, uint32(15), val)
	})

	t.Run("Increment uint64", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint64(10), DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, uint64(15), val)
	})

	t.Run("Increment float32", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", float32(10.5), DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, float32(15.5), val)
	})

	t.Run("Increment float64", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", float64(10.5), DefaultExpiration)
		err := c.Increment("key", 5)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, float64(15.5), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		err := c.Increment("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementFloat(t *testing.T) {
	t.Run("Increment float32", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", float32(10.5), DefaultExpiration)
		err := c.IncrementFloat("key", 5.0)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, float32(15.5), val)
	})

	t.Run("Increment float64", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", float64(10.5), DefaultExpiration)
		err := c.IncrementFloat("key", 5.0)
		assert.NoError(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, float64(15.5), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		err := c.IncrementFloat("nonexistent", 5.0)
		assert.Error(t, err)
	})
}

func TestCache_IncrementInt(t *testing.T) {
	t.Run("Increment int", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", 10, DefaultExpiration)
		newVal, err := c.IncrementInt("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, 15, newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, 15, val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementInt("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementInt8(t *testing.T) {
	t.Run("Increment int8", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int8(10), DefaultExpiration)
		newVal, err := c.IncrementInt8("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, int8(15), newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, int8(15), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementInt8("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementInt16(t *testing.T) {
	t.Run("Increment int16", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int16(10), DefaultExpiration)
		newVal, err := c.IncrementInt16("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, int16(15), newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, int16(15), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementInt16("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementInt32(t *testing.T) {
	t.Run("Increment int32", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int32(10), DefaultExpiration)
		newVal, err := c.IncrementInt32("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, int32(15), newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, int32(15), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementInt32("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementInt64(t *testing.T) {
	t.Run("Increment int64", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int64(10), DefaultExpiration)
		newVal, err := c.IncrementInt64("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, int64(15), newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, int64(15), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementInt64("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementUint(t *testing.T) {
	t.Run("Increment uint", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint(10), DefaultExpiration)
		newVal, err := c.IncrementUint("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, uint(15), newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, uint(15), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementUint("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementUintptr(t *testing.T) {
	t.Run("Increment uintptr", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uintptr(10), DefaultExpiration)
		newVal, err := c.IncrementUintptr("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, uintptr(15), newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, uintptr(15), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementUintptr("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementUint8(t *testing.T) {
	t.Run("Increment uint8", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint8(10), DefaultExpiration)
		newVal, err := c.IncrementUint8("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, uint8(15), newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, uint8(15), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementUint8("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementUint16(t *testing.T) {
	t.Run("Increment uint16", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint16(10), DefaultExpiration)
		newVal, err := c.IncrementUint16("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, uint16(15), newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, uint16(15), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementUint16("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementUint32(t *testing.T) {
	t.Run("Increment uint32", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint32(10), DefaultExpiration)
		newVal, err := c.IncrementUint32("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, uint32(15), newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, uint32(15), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementUint32("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementUint64(t *testing.T) {
	t.Run("Increment uint64", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint64(10), DefaultExpiration)
		newVal, err := c.IncrementUint64("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, uint64(15), newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, uint64(15), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementUint64("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementFloat32(t *testing.T) {
	t.Run("Increment float32", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", float32(10.5), DefaultExpiration)
		newVal, err := c.IncrementFloat32("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, float32(15.5), newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, float32(15.5), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementFloat32("nonexistent", 5)
		assert.Error(t, err)
	})
}

func TestCache_IncrementFloat64(t *testing.T) {
	t.Run("Increment float64", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", float64(10.5), DefaultExpiration)
		newVal, err := c.IncrementFloat64("key", 5)
		assert.NoError(t, err)
		assert.Equal(t, float64(15.5), newVal)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, float64(15.5), val)
	})

	t.Run("Increment non-existent key", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		_, err := c.IncrementFloat64("nonexistent", 5)
		assert.Error(t, err)
	})
}
