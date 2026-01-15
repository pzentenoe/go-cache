package cache

import (
	"fmt"
	"math"
)

// Increment an item of type int, int8, int16, int32, int64, uintptr, uint,
// uint8, uint32, or uint64, float32 or float64 by n. Returns an error if the
// item's value is not an integer, if it was not found, or if it is not
// possible to increment it by n. To retrieve the incremented value, use one
// of the specialized methods, e.g. IncrementInt64.
func (c *Cache) Increment(k string, n int64) error {
	return c.increment(k, n, func(val any) (any, error) {
		switch val := val.(type) {
		case int:
			return val + int(n), nil
		case int8:
			return val + int8(n), nil
		case int16:
			return val + int16(n), nil
		case int32:
			return val + int32(n), nil
		case int64:
			return val + n, nil
		case uint:
			return val + uint(n), nil
		case uintptr:
			return val + uintptr(n), nil
		case uint8:
			return val + uint8(n), nil
		case uint16:
			return val + uint16(n), nil
		case uint32:
			return val + uint32(n), nil
		case uint64:
			return val + uint64(n), nil
		case float32:
			return val + float32(n), nil
		case float64:
			return val + float64(n), nil
		default:
			return nil, fmt.Errorf("the value for %s is not an integer or float", k)
		}
	})
}

func (c *Cache) increment(k string, n any, incrementFunc func(any) (any, error)) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return fmt.Errorf("item %s not found", k)
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
	return c.increment(k, n, func(val any) (any, error) {
		switch val := val.(type) {
		case float32:
			return val + float32(n), nil
		case float64:
			return val + n, nil
		default:
			return nil, fmt.Errorf("the value for %s does not have type float32 or float64", k)
		}
	})
}

// IncrementInt Increment an item of type int by n. Returns an error if the item's value is
// not an int, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementInt(k string, n int) (int, error) {
	result := c.incrementTyped(k, n, int(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(int), nil
}

// IncrementInt8 Increment an item of type int8 by n. Returns an error if the item's value is
// not an int8, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementInt8(k string, n int8) (int8, error) {
	result := c.incrementTyped(k, n, int8(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(int8), nil
}

// IncrementInt16 Increment an item of type int16 by n. Returns an error if the item's value is
// not an int16, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementInt16(k string, n int16) (int16, error) {
	result := c.incrementTyped(k, n, int16(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(int16), nil
}

// IncrementInt32 Increment an item of type int32 by n. Returns an error if the item's value is
// not an int32, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementInt32(k string, n int32) (int32, error) {
	result := c.incrementTyped(k, n, int32(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(int32), nil
}

// IncrementInt64 Increment an item of type int64 by n. Returns an error if the item's value is
// not an int64, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementInt64(k string, n int64) (int64, error) {
	result := c.incrementTyped(k, n, int64(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(int64), nil
}

// IncrementUint Increment an item of type uint by n. Returns an error if the item's value is
// not an uint, or if it was not found. If there is no error, the incremented
// value is returned.
func (c *Cache) IncrementUint(k string, n uint) (uint, error) {
	result := c.incrementTyped(k, n, uint(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(uint), nil
}

// IncrementUintptr Increment an item of type uintptr by n. Returns an error if the item's value
// is not an uintptr, or if it was not found. If there is no error, the
// incremented value is returned.
func (c *Cache) IncrementUintptr(k string, n uintptr) (uintptr, error) {
	result := c.incrementTyped(k, n, uintptr(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(uintptr), nil
}

// IncrementUint8 Increment an item of type uint8 by n. Returns an error if the item's value
// is not an uint8, or if it was not found. If there is no error, the
// incremented value is returned.
func (c *Cache) IncrementUint8(k string, n uint8) (uint8, error) {
	result := c.incrementTyped(k, n, uint8(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(uint8), nil
}

// IncrementUint16 Increment an item of type uint16 by n. Returns an error if the item's value
// is not an uint16, or if it was not found. If there is no error, the
// incremented value is returned.
func (c *Cache) IncrementUint16(k string, n uint16) (uint16, error) {
	result := c.incrementTyped(k, n, uint16(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(uint16), nil
}

// IncrementUint32 Increment an item of type uint32 by n. Returns an error if the item's value
// is not an uint32, or if it was not found. If there is no error, the
// incremented value is returned.
func (c *Cache) IncrementUint32(k string, n uint32) (uint32, error) {
	result := c.incrementTyped(k, n, uint32(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(uint32), nil
}

// IncrementUint64 Increment an item of type uint64 by n. Returns an error if the item's value
// is not an uint64, or if it was not found. If there is no error, the
// incremented value is returned.
func (c *Cache) IncrementUint64(k string, n uint64) (uint64, error) {
	result := c.incrementTyped(k, n, uint64(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(uint64), nil
}

// IncrementFloat32 Increment an item of type float32 by n. Returns an error if the item's value
// is not a float32, or if it was not found. If there is no error, the
// incremented value is returned.
func (c *Cache) IncrementFloat32(k string, n float32) (float32, error) {
	result := c.incrementTyped(k, n, float32(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(float32), nil
}

// IncrementFloat64 Increment an item of type float64 by n. Returns an error if the item's value
// is not a float64, or if it was not found. If there is no error, the
// incremented value is returned.
func (c *Cache) IncrementFloat64(k string, n float64) (float64, error) {
	result := c.incrementTyped(k, n, float64(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(float64), nil
}

type incrementResult struct {
	value any
	err   error
}

func (c *Cache) incrementTyped(k string, n any, zero any) incrementResult {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return incrementResult{zero, fmt.Errorf("item %s not found", k)}
	}
	switch val := v.Object.(type) {
	case int:
		nv := n.(int)
		if (nv > 0 && val > math.MaxInt-nv) || (nv < 0 && val < math.MinInt-nv) {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = val + nv
	case int8:
		nv := n.(int8)
		if (nv > 0 && val > math.MaxInt8-nv) || (nv < 0 && val < math.MinInt8-nv) {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = val + nv
	case int16:
		nv := n.(int16)
		if (nv > 0 && val > math.MaxInt16-nv) || (nv < 0 && val < math.MinInt16-nv) {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = val + nv
	case int32:
		nv := n.(int32)
		if (nv > 0 && val > math.MaxInt32-nv) || (nv < 0 && val < math.MinInt32-nv) {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = val + nv
	case int64:
		nv := n.(int64)
		if (nv > 0 && val > math.MaxInt64-nv) || (nv < 0 && val < math.MinInt64-nv) {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = val + nv
	case uint:
		nv := n.(uint)
		if val > math.MaxUint-nv {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = val + nv
	case uintptr:
		nv := n.(uintptr)
		if val > math.MaxUint-uintptr(nv) {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = val + nv
	case uint8:
		nv := n.(uint8)
		if val > math.MaxUint8-nv {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = val + nv
	case uint16:
		nv := n.(uint16)
		if val > math.MaxUint16-nv {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = val + nv
	case uint32:
		nv := n.(uint32)
		if val > math.MaxUint32-nv {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = val + nv
	case uint64:
		nv := n.(uint64)
		if val > math.MaxUint64-nv {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = val + nv
	case float32:
		nv := n.(float32)
		result := val + nv
		if math.IsInf(float64(result), 0) {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = result
	case float64:
		nv := n.(float64)
		result := val + nv
		if math.IsInf(result, 0) {
			return incrementResult{zero, fmt.Errorf("overflow would occur for %s", k)}
		}
		v.Object = result
	default:
		return incrementResult{zero, fmt.Errorf("the value for %s is not a supported type", k)}
	}
	c.items[k] = v
	return incrementResult{v.Object, nil}
}
