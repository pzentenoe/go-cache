package cache

import "fmt"

// Decrement an item of type int, int8, int16, int32, int64, uintptr, uint,
// uint8, uint32, or uint64, float32 or float64 by n. Returns an error if the
// item's value is not an integer, if it was not found, or if it is not
// possible to decrement it by n. To retrieve the decremented value, use one
// of the specialized methods, e.g. DecrementInt64.
func (c *Cache) Decrement(k string, n int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return fmt.Errorf("Item not found")
	}
	switch val := v.Object.(type) {
	case int:
		v.Object = val - int(n)
	case int8:
		v.Object = val - int8(n)
	case int16:
		v.Object = val - int16(n)
	case int32:
		v.Object = val - int32(n)
	case int64:
		v.Object = val - n
	case uint:
		v.Object = val - uint(n)
	case uintptr:
		v.Object = val - uintptr(n)
	case uint8:
		v.Object = val - uint8(n)
	case uint16:
		v.Object = val - uint16(n)
	case uint32:
		v.Object = val - uint32(n)
	case uint64:
		v.Object = val - uint64(n)
	case float32:
		v.Object = val - float32(n)
	case float64:
		v.Object = val - float64(n)
	default:
		return fmt.Errorf("The value for %s is not an integer", k)
	}
	c.items[k] = v
	return nil
}

// Decrement an item of type float32 or float64 by n. Returns an error if the
// item's value is not floating point, if it was not found, or if it is not
// possible to decrement it by n. Pass a negative number to decrement the
// value. To retrieve the decremented value, use one of the specialized methods,
// e.g. DecrementFloat64.
func (c *Cache) DecrementFloat(k string, n float64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return fmt.Errorf("Item %s not found", k)
	}
	switch val := v.Object.(type) {
	case float32:
		v.Object = val - float32(n)
	case float64:
		v.Object = val - n
	default:
		return fmt.Errorf("The value for %s does not have type float32 or float64", k)
	}
	c.items[k] = v
	return nil
}

// Decrement an item of type int by n. Returns an error if the item's value is
// not an int, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt(k string, n int) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(int)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}

// Decrement an item of type int8 by n. Returns an error if the item's value is
// not an int8, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt8(k string, n int8) (int8, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(int8)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int8", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}

// Decrement an item of type int16 by n. Returns an error if the item's value is
// not an int16, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt16(k string, n int16) (int16, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(int16)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int16", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}

// Decrement an item of type int32 by n. Returns an error if the item's value is
// not an int32, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt32(k string, n int32) (int32, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(int32)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int32", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}

// Decrement an item of type int64 by n. Returns an error if the item's value is
// not an int64, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementInt64(k string, n int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(int64)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int64", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}

// Decrement an item of type uint by n. Returns an error if the item's value is
// not an uint, or if it was not found. If there is no error, the decremented
// value is returned.
func (c *Cache) DecrementUint(k string, n uint) (uint, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(uint)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an uint", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}

// Decrement an item of type uintptr by n. Returns an error if the item's value
// is not an uintptr, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementUintptr(k string, n uintptr) (uintptr, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(uintptr)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an uintptr", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}

// Decrement an item of type uint8 by n. Returns an error if the item's value
// is not an uint8, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementUint8(k string, n uint8) (uint8, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(uint8)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an uint8", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}

// Decrement an item of type uint16 by n. Returns an error if the item's value
// is not an uint16, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementUint16(k string, n uint16) (uint16, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(uint16)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an uint16", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}

// Decrement an item of type uint32 by n. Returns an error if the item's value
// is not an uint32, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementUint32(k string, n uint32) (uint32, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(uint32)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an uint32", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}

// Decrement an item of type uint64 by n. Returns an error if the item's value
// is not an uint64, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementUint64(k string, n uint64) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(uint64)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an uint64", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}

// Decrement an item of type float32 by n. Returns an error if the item's value
// is not an float32, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementFloat32(k string, n float32) (float32, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(float32)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an float32", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}

// Decrement an item of type float64 by n. Returns an error if the item's value
// is not an float64, or if it was not found. If there is no error, the
// decremented value is returned.
func (c *Cache) DecrementFloat64(k string, n float64) (float64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	val, ok := v.Object.(float64)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an float64", k)
	}
	val -= n
	v.Object = val
	c.items[k] = v
	return val, nil
}
