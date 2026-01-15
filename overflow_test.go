package cache

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrementOverflow(t *testing.T) {
	t.Run("IncrementUint8 overflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", uint8(255), NoExpiration)

		_, err := c.IncrementUint8("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overflow would occur")

		// Value should remain unchanged
		val, _ := c.Get("counter")
		assert.Equal(t, uint8(255), val)
	})

	t.Run("IncrementUint16 overflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", uint16(65535), NoExpiration)

		_, err := c.IncrementUint16("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overflow would occur")
	})

	t.Run("IncrementUint32 overflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", uint32(math.MaxUint32), NoExpiration)

		_, err := c.IncrementUint32("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overflow would occur")
	})

	t.Run("IncrementUint64 overflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", uint64(math.MaxUint64), NoExpiration)

		_, err := c.IncrementUint64("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overflow would occur")
	})

	t.Run("IncrementInt8 overflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", int8(127), NoExpiration)

		_, err := c.IncrementInt8("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overflow would occur")
	})

	t.Run("IncrementInt8 underflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", int8(-128), NoExpiration)

		_, err := c.IncrementInt8("counter", -1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overflow would occur")
	})

	t.Run("IncrementInt16 overflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", int16(math.MaxInt16), NoExpiration)

		_, err := c.IncrementInt16("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overflow would occur")
	})

	t.Run("IncrementInt32 overflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", int32(math.MaxInt32), NoExpiration)

		_, err := c.IncrementInt32("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overflow would occur")
	})

	t.Run("IncrementInt64 overflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", int64(math.MaxInt64), NoExpiration)

		_, err := c.IncrementInt64("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overflow would occur")
	})

	t.Run("IncrementFloat32 infinity detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", float32(math.MaxFloat32), NoExpiration)

		_, err := c.IncrementFloat32("counter", float32(math.MaxFloat32))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overflow would occur")
	})

	t.Run("IncrementFloat64 infinity detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", float64(math.MaxFloat64), NoExpiration)

		_, err := c.IncrementFloat64("counter", float64(math.MaxFloat64))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overflow would occur")
	})
}

func TestDecrementUnderflow(t *testing.T) {
	t.Run("DecrementUint8 underflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", uint8(0), NoExpiration)

		_, err := c.DecrementUint8("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "underflow would occur")

		// Value should remain unchanged
		val, _ := c.Get("counter")
		assert.Equal(t, uint8(0), val)
	})

	t.Run("DecrementUint16 underflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", uint16(0), NoExpiration)

		_, err := c.DecrementUint16("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "underflow would occur")
	})

	t.Run("DecrementUint32 underflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", uint32(0), NoExpiration)

		_, err := c.DecrementUint32("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "underflow would occur")
	})

	t.Run("DecrementUint64 underflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", uint64(0), NoExpiration)

		_, err := c.DecrementUint64("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "underflow would occur")
	})

	t.Run("DecrementInt8 underflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", int8(-128), NoExpiration)

		_, err := c.DecrementInt8("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "underflow would occur")
	})

	t.Run("DecrementInt8 overflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", int8(127), NoExpiration)

		_, err := c.DecrementInt8("counter", -1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "underflow would occur")
	})

	t.Run("DecrementInt16 underflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", int16(math.MinInt16), NoExpiration)

		_, err := c.DecrementInt16("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "underflow would occur")
	})

	t.Run("DecrementInt32 underflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", int32(math.MinInt32), NoExpiration)

		_, err := c.DecrementInt32("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "underflow would occur")
	})

	t.Run("DecrementInt64 underflow detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", int64(math.MinInt64), NoExpiration)

		_, err := c.DecrementInt64("counter", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "underflow would occur")
	})

	t.Run("DecrementFloat32 infinity detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", float32(-math.MaxFloat32), NoExpiration)

		_, err := c.DecrementFloat32("counter", float32(math.MaxFloat32))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "underflow would occur")
	})

	t.Run("DecrementFloat64 infinity detection", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", float64(-math.MaxFloat64), NoExpiration)

		_, err := c.DecrementFloat64("counter", float64(math.MaxFloat64))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "underflow would occur")
	})
}

func TestIncrementDecrementWithinBounds(t *testing.T) {
	t.Run("IncrementUint8 within bounds succeeds", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", uint8(254), NoExpiration)

		result, err := c.IncrementUint8("counter", 1)
		assert.NoError(t, err)
		assert.Equal(t, uint8(255), result)
	})

	t.Run("DecrementUint8 within bounds succeeds", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", uint8(1), NoExpiration)

		result, err := c.DecrementUint8("counter", 1)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0), result)
	})

	t.Run("IncrementInt8 within bounds succeeds", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", int8(126), NoExpiration)

		result, err := c.IncrementInt8("counter", 1)
		assert.NoError(t, err)
		assert.Equal(t, int8(127), result)
	})

	t.Run("DecrementInt8 within bounds succeeds", func(t *testing.T) {
		c := New(NoExpiration, 0)
		c.Set("counter", int8(-127), NoExpiration)

		result, err := c.DecrementInt8("counter", 1)
		assert.NoError(t, err)
		assert.Equal(t, int8(-128), result)
	})
}
