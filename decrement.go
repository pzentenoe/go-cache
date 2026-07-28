package cache

import ()

// Decrement an item of type int, int8, int16, int32, int64, uintptr, uint,
// uint8, uint32, or uint64, float32 or float64 by n. Returns an error if the
// item's value is not an integer, if it was not found, or if it is not
// possible to decrement it by n. To retrieve the decremented value, use one
// of the specialized methods, e.g. DecrementInt64.
func (c *Cache) Decrement(k string, n int64) error {
	return c.decrement(k, func(val any) (any, error) {
		switch val := val.(type) {
		case int:
			return decrementSigned[int](k, val, int(n))
		case int8:
			return decrementSigned[int8](k, val, int8(n))
		case int16:
			return decrementSigned[int16](k, val, int16(n))
		case int32:
			return decrementSigned[int32](k, val, int32(n))
		case int64:
			return decrementSigned[int64](k, val, n)
		case uint:
			return decrementUnsigned[uint](k, val, uint(n))
		case uintptr:
			return decrementUnsigned[uintptr](k, val, uintptr(n))
		case uint8:
			return decrementUnsigned[uint8](k, val, uint8(n))
		case uint16:
			return decrementUnsigned[uint16](k, val, uint16(n))
		case uint32:
			return decrementUnsigned[uint32](k, val, uint32(n))
		case uint64:
			return decrementUnsigned[uint64](k, val, uint64(n))
		case float32:
			return decrementFloat[float32](k, val, float32(n))
		case float64:
			return decrementFloat[float64](k, val, float64(n))
		default:
			return nil, keyErrorf(ErrTypeMismatch, errNotIntegerOrFloatFormat, k)
		}
	})
}

func (c *Cache) decrement(k string, decrementFunc func(any) (any, error)) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return keyErrorf(ErrNotFound, errItemNotFoundFormat, k)
	}
	newValue, err := decrementFunc(v.Object)
	if err != nil {
		return err
	}
	v.Object = newValue
	c.items[k] = v
	return nil
}

// DecrementFloat Decrement an item of type float32 or float64 by n. Returns an error if the
// item's value is not floating point, if it was not found, or if it is not
// possible to decrement it by n. Pass a negative number to increment the
// value. To retrieve the decremented value, use one of the specialized methods,
// e.g. DecrementFloat64.
func (c *Cache) DecrementFloat(k string, n float64) error {
	return c.decrement(k, func(val any) (any, error) {
		switch val := val.(type) {
		case float32:
			return decrementFloat[float32](k, val, float32(n))
		case float64:
			return decrementFloat[float64](k, val, n)
		default:
			return nil, keyErrorf(ErrTypeMismatch, errNotFloatTypeFormat, k)
		}
	})
}

// DecrementInt Decrement an item of type int by n. Returns an error if the item's value is
// not an int, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt(k string, n int) (int, error) {
	return mutateTyped(c, k, n, decrementSigned[int])
}

// DecrementInt8 Decrement an item of type int8 by n. Returns an error if the item's value is
// not an int8, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt8(k string, n int8) (int8, error) {
	return mutateTyped(c, k, n, decrementSigned[int8])
}

// DecrementInt16 Decrement an item of type int16 by n. Returns an error if the item's value is
// not an int16, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt16(k string, n int16) (int16, error) {
	return mutateTyped(c, k, n, decrementSigned[int16])
}

// DecrementInt32 Decrement an item of type int32 by n. Returns an error if the item's value is
// not an int32, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt32(k string, n int32) (int32, error) {
	return mutateTyped(c, k, n, decrementSigned[int32])
}

// DecrementInt64 Decrement an item of type int64 by n. Returns an error if the item's value is
// not an int64, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt64(k string, n int64) (int64, error) {
	return mutateTyped(c, k, n, decrementSigned[int64])
}

// DecrementUint Decrement an item of type uint by n. Returns an error if the item's value is
// not an uint, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementUint(k string, n uint) (uint, error) {
	return mutateTyped(c, k, n, decrementUnsigned[uint])
}

// DecrementUintptr Decrement an item of type uintptr by n. Returns an error if the item's value
// is not an uintptr, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementUintptr(k string, n uintptr) (uintptr, error) {
	return mutateTyped(c, k, n, decrementUnsigned[uintptr])
}

// DecrementUint8 Decrement an item of type uint8 by n. Returns an error if the item's value is
// not an uint8, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementUint8(k string, n uint8) (uint8, error) {
	return mutateTyped(c, k, n, decrementUnsigned[uint8])
}

// DecrementUint16 Decrement an item of type uint16 by n. Returns an error if the item's value is
// not an uint16, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementUint16(k string, n uint16) (uint16, error) {
	return mutateTyped(c, k, n, decrementUnsigned[uint16])
}

// DecrementUint32 Decrement an item of type uint32 by n. Returns an error if the item's value is
// not an uint32, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementUint32(k string, n uint32) (uint32, error) {
	return mutateTyped(c, k, n, decrementUnsigned[uint32])
}

// DecrementUint64 Decrement an item of type uint64 by n. Returns an error if the item's value is
// not an uint64, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementUint64(k string, n uint64) (uint64, error) {
	return mutateTyped(c, k, n, decrementUnsigned[uint64])
}

// DecrementFloat32 Decrement an item of type float32 by n. Returns an error if the item's value
// is not a float32, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementFloat32(k string, n float32) (float32, error) {
	return mutateTyped(c, k, n, decrementFloat[float32])
}

// DecrementFloat64 Decrement an item of type float64 by n. Returns an error if the item's value
// is not a float64, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementFloat64(k string, n float64) (float64, error) {
	return mutateTyped(c, k, n, decrementFloat[float64])
}
