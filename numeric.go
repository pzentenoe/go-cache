package cache

import (
	"fmt"
	"math"
)

// Numeric type sets supported by the typed increment/decrement methods.
type signed interface {
	int | int8 | int16 | int32 | int64
}

type unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

type floats interface {
	float32 | float64
}

// signedBounds returns the [min, max] range of T as int64. Widening is
// lossless: every signed type fits in int64.
func signedBounds[T signed]() (min, max int64) {
	switch any(*new(T)).(type) {
	case int:
		return math.MinInt, math.MaxInt
	case int8:
		return math.MinInt8, math.MaxInt8
	case int16:
		return math.MinInt16, math.MaxInt16
	case int32:
		return math.MinInt32, math.MaxInt32
	case int64:
		return math.MinInt64, math.MaxInt64
	}
	return 0, 0 // unreachable: T is constrained to signed integers
}

// unsignedMax returns the max value of T as uint64. Widening is lossless.
func unsignedMax[T unsigned]() uint64 {
	switch any(*new(T)).(type) {
	case uint, uintptr:
		return math.MaxUint
	case uint8:
		return math.MaxUint8
	case uint16:
		return math.MaxUint16
	case uint32:
		return math.MaxUint32
	case uint64:
		return math.MaxUint64
	}
	return 0 // unreachable: T is constrained to unsigned integers
}

// mutateTyped applies op to the cached value for k under lock and stores the
// result. The stored value must be exactly of type T; otherwise an error is
// returned instead of panicking on a failed type assertion.
func mutateTyped[T any](c *Cache, k string, n T, op func(k string, v, n T) (T, error)) (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var zero T
	v, found := c.items[k]
	if !found || v.Expired() {
		return zero, fmt.Errorf(errItemNotFoundFormat, k)
	}
	val, ok := v.Object.(T)
	if !ok {
		return zero, fmt.Errorf(errTypeMismatchFormat, k)
	}
	nv, err := op(k, val, n)
	if err != nil {
		return zero, err
	}
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

// Bounds checks run in the widened int64/uint64 domain, which preserves the
// exact semantics of the original per-type formulas: v and n promote without
// loss, and once the check passes the result is guaranteed to fit in T.
// Floats are the exception: they operate in T, since a float32 overflow to
// +Inf would not overflow in float64.

func incrementSigned[T signed](k string, v, n T) (T, error) {
	min, max := signedBounds[T]()
	v64, n64 := int64(v), int64(n)
	if (n64 > 0 && v64 > max-n64) || (n64 < 0 && v64 < min-n64) {
		return v, fmt.Errorf(errOverflowFormat, k)
	}
	return T(v64 + n64), nil
}

func decrementSigned[T signed](k string, v, n T) (T, error) {
	min, max := signedBounds[T]()
	v64, n64 := int64(v), int64(n)
	if (n64 > 0 && v64 < min+n64) || (n64 < 0 && v64 > max+n64) {
		return v, fmt.Errorf(errUnderflowFormat, k)
	}
	return T(v64 - n64), nil
}

func incrementUnsigned[T unsigned](k string, v, n T) (T, error) {
	max := unsignedMax[T]()
	v64, n64 := uint64(v), uint64(n)
	if v64 > max-n64 {
		return v, fmt.Errorf(errOverflowFormat, k)
	}
	return T(v64 + n64), nil
}

func decrementUnsigned[T unsigned](k string, v, n T) (T, error) {
	v64, n64 := uint64(v), uint64(n)
	if v64 < n64 {
		return v, fmt.Errorf(errUnderflowFormat, k)
	}
	return T(v64 - n64), nil
}

func incrementFloat[T floats](k string, v, n T) (T, error) {
	sum := v + n
	if math.IsInf(float64(sum), 0) {
		return v, fmt.Errorf(errOverflowFormat, k)
	}
	return sum, nil
}

func decrementFloat[T floats](k string, v, n T) (T, error) {
	diff := v - n
	if math.IsInf(float64(diff), 0) {
		return v, fmt.Errorf(errUnderflowFormat, k)
	}
	return diff, nil
}
