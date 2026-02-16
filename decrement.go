package cache

import (
	"fmt"
	"math"
)

// Decrement an item of type int, int8, int16, int32, int64, uintptr, uint,
// uint8, uint32, or uint64, float32 or float64 by n. Returns an error if the
// item's value is not an integer, if it was not found, or if it is not
// possible to decrement it by n. To retrieve the decremented value, use one
// of the specialized methods, e.g. DecrementInt64.
func (c *Cache) Decrement(k string, n int64) error {
	return c.decrement(k, func(val any) (any, error) {
		switch val := val.(type) {
		case int:
			nv := int(n)
			if (nv > 0 && val < math.MinInt+nv) || (nv < 0 && val > math.MaxInt+nv) {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return val - nv, nil
		case int8:
			nv := int8(n)
			if (nv > 0 && val < math.MinInt8+nv) || (nv < 0 && val > math.MaxInt8+nv) {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return val - nv, nil
		case int16:
			nv := int16(n)
			if (nv > 0 && val < math.MinInt16+nv) || (nv < 0 && val > math.MaxInt16+nv) {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return val - nv, nil
		case int32:
			nv := int32(n)
			if (nv > 0 && val < math.MinInt32+nv) || (nv < 0 && val > math.MaxInt32+nv) {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return val - nv, nil
		case int64:
			if (n > 0 && val < math.MinInt64+n) || (n < 0 && val > math.MaxInt64+n) {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return val - n, nil
		case uint:
			nv := uint(n)
			if val < nv {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return val - nv, nil
		case uintptr:
			nv := uintptr(n)
			if val < nv {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return val - nv, nil
		case uint8:
			nv := uint8(n)
			if val < nv {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return val - nv, nil
		case uint16:
			nv := uint16(n)
			if val < nv {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return val - nv, nil
		case uint32:
			nv := uint32(n)
			if val < nv {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return val - nv, nil
		case uint64:
			nv := uint64(n)
			if val < nv {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return val - nv, nil
		case float32:
			result := val - float32(n)
			if math.IsInf(float64(result), 0) {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return result, nil
		case float64:
			result := val - float64(n)
			if math.IsInf(result, 0) {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return result, nil
		default:
			return nil, fmt.Errorf(errNotIntegerFormat, k)
		}
	})
}

func (c *Cache) decrement(k string, decrementFunc func(any) (any, error)) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return fmt.Errorf(errItemNotFoundFormat, k)
	}
	newValue, err := decrementFunc(v.Object)
	if err != nil {
		return err
	}
	v.Object = newValue
	return nil
}

// DecrementFloat Decrement an item of type float32 or float64 by n. Returns an error if the
// item's value is not floating point, if it was not found, or if it is not
// possible to decrement it by n. Pass a negative number to decrement the
// value. To retrieve the decremented value, use one of the specialized methods,
// e.g. DecrementFloat64.
func (c *Cache) DecrementFloat(k string, n float64) error {
	return c.decrement(k, func(val any) (any, error) {
		switch val := val.(type) {
		case float32:
			result := val - float32(n)
			if math.IsInf(float64(result), 0) {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return result, nil
		case float64:
			result := val - n
			if math.IsInf(result, 0) {
				return nil, fmt.Errorf(errUnderflowFormat, k)
			}
			return result, nil
		default:
			return nil, fmt.Errorf(errNotFloatTypeFormat, k)
		}
	})
}

// DecrementInt Decrement an item of type int by n. Returns an error if the item's value is
// not an int, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt(k string, n int) (int, error) {
	result := c.decrementTyped(k, n, int(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(int), nil
}

// DecrementInt8 Decrement an item of type int8 by n. Returns an error if the item's value is
// not an int8, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt8(k string, n int8) (int8, error) {
	result := c.decrementTyped(k, n, int8(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(int8), nil
}

// DecrementInt16 Decrement an item of type int16 by n. Returns an error if the item's value is
// not an int16, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt16(k string, n int16) (int16, error) {
	result := c.decrementTyped(k, n, int16(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(int16), nil
}

// DecrementInt32 Decrement an item of type int32 by n. Returns an error if the item's value is
// not an int32, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt32(k string, n int32) (int32, error) {
	result := c.decrementTyped(k, n, int32(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(int32), nil
}

// DecrementInt64 Decrement an item of type int64 by n. Returns an error if the item's value is
// not an int64, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt64(k string, n int64) (int64, error) {
	result := c.decrementTyped(k, n, int64(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(int64), nil
}

// DecrementUint Decrement an item of type uint by n. Returns an error if the item's value is
// not an uint, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementUint(k string, n uint) (uint, error) {
	result := c.decrementTyped(k, n, uint(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(uint), nil
}

// DecrementUintptr Decrement an item of type uintptr by n. Returns an error if the item's value
// is not an uintptr, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementUintptr(k string, n uintptr) (uintptr, error) {
	result := c.decrementTyped(k, n, uintptr(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(uintptr), nil
}

// DecrementUint8 Decrement an item of type uint8 by n. Returns an error if the item's value
// is not an uint8, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementUint8(k string, n uint8) (uint8, error) {
	result := c.decrementTyped(k, n, uint8(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(uint8), nil
}

// DecrementUint16 Decrement an item of type uint16 by n. Returns an error if the item's value
// is not an uint16, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementUint16(k string, n uint16) (uint16, error) {
	result := c.decrementTyped(k, n, uint16(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(uint16), nil
}

// DecrementUint32 Decrement an item of type uint32 by n. Returns an error if the item's value
// is not an uint32, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementUint32(k string, n uint32) (uint32, error) {
	result := c.decrementTyped(k, n, uint32(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(uint32), nil
}

// DecrementUint64 Decrement an item of type uint64 by n. Returns an error if the item's value
// is not an uint64, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementUint64(k string, n uint64) (uint64, error) {
	result := c.decrementTyped(k, n, uint64(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(uint64), nil
}

// DecrementFloat32 Decrement an item of type float32 by n. Returns an error if the item's value
// is not a float32, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementFloat32(k string, n float32) (float32, error) {
	result := c.decrementTyped(k, n, float32(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(float32), nil
}

// DecrementFloat64 Decrement an item of type float64 by n. Returns an error if the item's value
// is not a float64, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementFloat64(k string, n float64) (float64, error) {
	result := c.decrementTyped(k, n, float64(0))
	if result.err != nil {
		return 0, result.err
	}
	return result.value.(float64), nil
}

func (c *Cache) decrementTyped(k string, n any, zero any) operationResult {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return operationResult{zero, fmt.Errorf(errItemNotFoundFormat, k)}
	}
	switch val := v.Object.(type) {
	case int:
		nv := n.(int)
		if (nv > 0 && val < math.MinInt+nv) || (nv < 0 && val > math.MaxInt+nv) {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = val - nv
	case int8:
		nv := n.(int8)
		if (nv > 0 && val < math.MinInt8+nv) || (nv < 0 && val > math.MaxInt8+nv) {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = val - nv
	case int16:
		nv := n.(int16)
		if (nv > 0 && val < math.MinInt16+nv) || (nv < 0 && val > math.MaxInt16+nv) {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = val - nv
	case int32:
		nv := n.(int32)
		if (nv > 0 && val < math.MinInt32+nv) || (nv < 0 && val > math.MaxInt32+nv) {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = val - nv
	case int64:
		nv := n.(int64)
		if (nv > 0 && val < math.MinInt64+nv) || (nv < 0 && val > math.MaxInt64+nv) {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = val - nv
	case uint:
		nv := n.(uint)
		if val < nv {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = val - nv
	case uintptr:
		nv := n.(uintptr)
		if val < nv {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = val - nv
	case uint8:
		nv := n.(uint8)
		if val < nv {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = val - nv
	case uint16:
		nv := n.(uint16)
		if val < nv {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = val - nv
	case uint32:
		nv := n.(uint32)
		if val < nv {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = val - nv
	case uint64:
		nv := n.(uint64)
		if val < nv {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = val - nv
	case float32:
		nv := n.(float32)
		result := val - nv
		if math.IsInf(float64(result), 0) {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = result
	case float64:
		nv := n.(float64)
		result := val - nv
		if math.IsInf(result, 0) {
			return operationResult{zero, fmt.Errorf(errUnderflowFormat, k)}
		}
		v.Object = result
	default:
		return operationResult{zero, fmt.Errorf(errUnsupportedTypeFormat, k)}
	}
	return operationResult{v.Object, nil}
}
