package cache

import ()

// Increment an item of type int, int8, int16, int32, int64, uintptr, uint,
// uint8, uint32, or uint64, float32 or float64 by n. Returns an error if the
// item's value is not an integer, if it was not found, or if it is not
// possible to increment it by n. To retrieve the incremented value, use one
// of the specialized methods, e.g. IncrementInt64.
func (c *Cache) Increment(k string, n int64) error {
	return c.increment(k, func(val any) (any, error) {
		switch val := val.(type) {
		case int:
			return incrementSigned[int](k, val, int(n))
		case int8:
			return incrementSigned[int8](k, val, int8(n))
		case int16:
			return incrementSigned[int16](k, val, int16(n))
		case int32:
			return incrementSigned[int32](k, val, int32(n))
		case int64:
			return incrementSigned[int64](k, val, n)
		case uint:
			return incrementUnsigned[uint](k, val, uint(n))
		case uintptr:
			return incrementUnsigned[uintptr](k, val, uintptr(n))
		case uint8:
			return incrementUnsigned[uint8](k, val, uint8(n))
		case uint16:
			return incrementUnsigned[uint16](k, val, uint16(n))
		case uint32:
			return incrementUnsigned[uint32](k, val, uint32(n))
		case uint64:
			return incrementUnsigned[uint64](k, val, uint64(n))
		case float32:
			return incrementFloat[float32](k, val, float32(n))
		case float64:
			return incrementFloat[float64](k, val, float64(n))
		default:
			return nil, keyErrorf(ErrTypeMismatch, errNotIntegerOrFloatFormat, k)
		}
	})
}

func (c *Cache) increment(k string, incrementFunc func(any) (any, error)) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return keyErrorf(ErrNotFound, errItemNotFoundFormat, k)
	}
	newValue, err := incrementFunc(v.Object)
	if err != nil {
		return err
	}
	v.Object = newValue
	c.items[k] = v
	return nil
}

// IncrementFloat Increment an item of type float32 or float64 by n. Returns an error if the
// item's value is not floating point, if it was not found, or if it is not
// possible to increment it by n. Pass a negative number to decrement the
// value. To retrieve the incremented value, use one of the specialized methods,
// e.g. IncrementFloat64.
func (c *Cache) IncrementFloat(k string, n float64) error {
	return c.increment(k, func(val any) (any, error) {
		switch val := val.(type) {
		case float32:
			return incrementFloat[float32](k, val, float32(n))
		case float64:
			return incrementFloat[float64](k, val, n)
		default:
			return nil, keyErrorf(ErrTypeMismatch, errNotFloatTypeFormat, k)
		}
	})
}

// IncrementInt Increment an item of type int by n. Returns an error if the item's value is
// not an int, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementInt(k string, n int) (int, error) {
	return mutateTyped(c, k, n, incrementSigned[int])
}

// IncrementInt8 Increment an item of type int8 by n. Returns an error if the item's value is
// not an int8, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementInt8(k string, n int8) (int8, error) {
	return mutateTyped(c, k, n, incrementSigned[int8])
}

// IncrementInt16 Increment an item of type int16 by n. Returns an error if the item's value is
// not an int16, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementInt16(k string, n int16) (int16, error) {
	return mutateTyped(c, k, n, incrementSigned[int16])
}

// IncrementInt32 Increment an item of type int32 by n. Returns an error if the item's value is
// not an int32, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementInt32(k string, n int32) (int32, error) {
	return mutateTyped(c, k, n, incrementSigned[int32])
}

// IncrementInt64 Increment an item of type int64 by n. Returns an error if the item's value is
// not an int64, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementInt64(k string, n int64) (int64, error) {
	return mutateTyped(c, k, n, incrementSigned[int64])
}

// IncrementUint Increment an item of type uint by n. Returns an error if the item's value is
// not an uint, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementUint(k string, n uint) (uint, error) {
	return mutateTyped(c, k, n, incrementUnsigned[uint])
}

// IncrementUintptr Increment an item of type uintptr by n. Returns an error if the item's value
// is not an uintptr, or if it was not found. If there is no error, the
// incremented value is returned.
func (c *Cache) IncrementUintptr(k string, n uintptr) (uintptr, error) {
	return mutateTyped(c, k, n, incrementUnsigned[uintptr])
}

// IncrementUint8 Increment an item of type uint8 by n. Returns an error if the item's value is
// not an uint8, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementUint8(k string, n uint8) (uint8, error) {
	return mutateTyped(c, k, n, incrementUnsigned[uint8])
}

// IncrementUint16 Increment an item of type uint16 by n. Returns an error if the item's value is
// not an uint16, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementUint16(k string, n uint16) (uint16, error) {
	return mutateTyped(c, k, n, incrementUnsigned[uint16])
}

// IncrementUint32 Increment an item of type uint32 by n. Returns an error if the item's value is
// not an uint32, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementUint32(k string, n uint32) (uint32, error) {
	return mutateTyped(c, k, n, incrementUnsigned[uint32])
}

// IncrementUint64 Increment an item of type uint64 by n. Returns an error if the item's value is
// not an uint64, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementUint64(k string, n uint64) (uint64, error) {
	return mutateTyped(c, k, n, incrementUnsigned[uint64])
}

// IncrementFloat32 Increment an item of type float32 by n. Returns an error if the item's value
// is not a float32, or if it was not found. If there is no error, the
// incremented value is returned.
func (c *Cache) IncrementFloat32(k string, n float32) (float32, error) {
	return mutateTyped(c, k, n, incrementFloat[float32])
}

// IncrementFloat64 Increment an item of type float64 by n. Returns an error if the item's value
// is not a float64, or if it was not found. If there is no error, the
// incremented value is returned.
func (c *Cache) IncrementFloat64(k string, n float64) (float64, error) {
	return mutateTyped(c, k, n, incrementFloat[float64])
}
