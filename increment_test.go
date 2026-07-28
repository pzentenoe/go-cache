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

	t.Run("Increment non-float type", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", "not a float", DefaultExpiration)
		err := c.IncrementFloat("key", 5.0)
		assert.Error(t, err)
	})

	t.Run("Increment non-float int type", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", 10, DefaultExpiration)
		err := c.IncrementFloat("key", 5.0)
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

func TestCache_IncrementTyped_Overflow(t *testing.T) {
	t.Run("IncrementInt overflow with positive n", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", 9223372036854775807, DefaultExpiration) // math.MaxInt64
		_, err := c.IncrementInt("key", 1)
		assert.Error(t, err)
	})

	t.Run("IncrementInt underflow with negative n", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", -9223372036854775808, DefaultExpiration) // math.MinInt64
		_, err := c.IncrementInt("key", -1)
		assert.Error(t, err)
	})

	t.Run("IncrementInt8 overflow with positive n", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int8(127), DefaultExpiration) // math.MaxInt8
		_, err := c.IncrementInt8("key", 1)
		assert.Error(t, err)
	})

	t.Run("IncrementInt8 underflow with negative n", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int8(-128), DefaultExpiration) // math.MinInt8
		_, err := c.IncrementInt8("key", -1)
		assert.Error(t, err)
	})

	t.Run("IncrementInt16 overflow with positive n", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int16(32767), DefaultExpiration) // math.MaxInt16
		_, err := c.IncrementInt16("key", 1)
		assert.Error(t, err)
	})

	t.Run("IncrementInt16 underflow with negative n", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int16(-32768), DefaultExpiration) // math.MinInt16
		_, err := c.IncrementInt16("key", -1)
		assert.Error(t, err)
	})

	t.Run("IncrementInt32 overflow with positive n", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int32(2147483647), DefaultExpiration) // math.MaxInt32
		_, err := c.IncrementInt32("key", 1)
		assert.Error(t, err)
	})

	t.Run("IncrementInt32 underflow with negative n", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int32(-2147483648), DefaultExpiration) // math.MinInt32
		_, err := c.IncrementInt32("key", -1)
		assert.Error(t, err)
	})

	t.Run("IncrementInt64 overflow with positive n", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int64(9223372036854775807), DefaultExpiration) // math.MaxInt64
		_, err := c.IncrementInt64("key", 1)
		assert.Error(t, err)
	})

	t.Run("IncrementInt64 underflow with negative n", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", int64(-9223372036854775808), DefaultExpiration) // math.MinInt64
		_, err := c.IncrementInt64("key", -1)
		assert.Error(t, err)
	})

	t.Run("IncrementUint overflow", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint(18446744073709551615), DefaultExpiration) // math.MaxUint64
		_, err := c.IncrementUint("key", 1)
		assert.Error(t, err)
	})

	t.Run("IncrementUintptr overflow", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uintptr(18446744073709551615), DefaultExpiration) // math.MaxUint
		_, err := c.IncrementUintptr("key", 1)
		assert.Error(t, err)
	})

	t.Run("IncrementUint8 overflow", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint8(255), DefaultExpiration) // math.MaxUint8
		_, err := c.IncrementUint8("key", 1)
		assert.Error(t, err)
	})

	t.Run("IncrementUint16 overflow", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint16(65535), DefaultExpiration) // math.MaxUint16
		_, err := c.IncrementUint16("key", 1)
		assert.Error(t, err)
	})

	t.Run("IncrementUint32 overflow", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint32(4294967295), DefaultExpiration) // math.MaxUint32
		_, err := c.IncrementUint32("key", 1)
		assert.Error(t, err)
	})

	t.Run("IncrementUint64 overflow", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", uint64(18446744073709551615), DefaultExpiration) // math.MaxUint64
		_, err := c.IncrementUint64("key", 1)
		assert.Error(t, err)
	})

	t.Run("IncrementFloat32 overflow", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", float32(3.4e38), DefaultExpiration) // Near math.MaxFloat32
		_, err := c.IncrementFloat32("key", 1e38)
		assert.Error(t, err)
	})

	t.Run("IncrementFloat64 overflow", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key", float64(1.7e308), DefaultExpiration) // Near math.MaxFloat64
		_, err := c.IncrementFloat64("key", 1e308)
		assert.Error(t, err)
	})
}

// TestIncrementTypeMismatch is a regression test: typed increment methods
// asserted n to the stored value's type, panicking on mismatch (e.g.
// IncrementInt64 on an int32 value). They must return an error instead.
func TestIncrementTypeMismatch(t *testing.T) {
	t.Run("IncrementInt64 on int32 value", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("key", int32(5), NoExpiration)
		_, err := c.IncrementInt64("key", 1)
		assert.Error(t, err)

		val, found := c.Get("key")
		assert.True(t, found)
		assert.Equal(t, int32(5), val, "value must be unchanged on type mismatch")
	})

	t.Run("IncrementInt on string value", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("key", "not-a-number", NoExpiration)
		_, err := c.IncrementInt("key", 1)
		assert.Error(t, err)
	})

	t.Run("Matching type still works", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("key", int32(5), NoExpiration)
		v, err := c.IncrementInt32("key", 2)
		assert.NoError(t, err)
		assert.Equal(t, int32(7), v)
	})
}
