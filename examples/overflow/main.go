package main

import (
	"fmt"
	"math"

	"github.com/pzentenoe/go-cache"
)

func main() {
	fmt.Println("=== Overflow/Underflow Protection Example ===")
	c := cache.New(cache.NoExpiration, 0)

	// Unsigned integer overflow protection
	demonstrateUintOverflow(c)

	// Signed integer overflow/underflow protection
	demonstrateIntOverflow(c)

	// Float infinity protection
	demonstrateFloatOverflow(c)

	// Safe operations within bounds
	demonstrateSafeOperations(c)
}

func demonstrateUintOverflow(c *cache.Cache) {
	fmt.Println("=== Unsigned Integer Overflow Protection ===")

	// uint8 overflow
	c.Set("uint8_counter", uint8(255), cache.NoExpiration)
	fmt.Println("uint8_counter = 255 (max value)")

	_, err := c.IncrementUint8("uint8_counter", 1)
	if err != nil {
		fmt.Printf("✓ Overflow prevented: %v\n", err)
	}

	// Value should remain unchanged
	if val, _ := c.Get("uint8_counter"); val == uint8(255) {
		fmt.Println("✓ Value unchanged: 255")
	}

	// uint16 overflow
	c.Set("uint16_counter", uint16(65535), cache.NoExpiration)
	fmt.Println("\nuint16_counter = 65535 (max value)")

	_, err = c.IncrementUint16("uint16_counter", 1)
	if err != nil {
		fmt.Printf("✓ Overflow prevented: %v\n", err)
	}

	// uint32 overflow
	c.Set("uint32_counter", uint32(math.MaxUint32), cache.NoExpiration)
	fmt.Printf("\nuint32_counter = %d (max value)\n", math.MaxUint32)

	_, err = c.IncrementUint32("uint32_counter", 1)
	if err != nil {
		fmt.Printf("✓ Overflow prevented: %v\n", err)
	}

	// uint64 overflow
	c.Set("uint64_counter", uint64(math.MaxUint64), cache.NoExpiration)
	fmt.Printf("\nuint64_counter = %d (max value)\n", uint64(math.MaxUint64))

	_, err = c.IncrementUint64("uint64_counter", 1)
	if err != nil {
		fmt.Printf("✓ Overflow prevented: %v\n\n", err)
	}
}

func demonstrateIntOverflow(c *cache.Cache) {
	fmt.Println("=== Signed Integer Overflow/Underflow Protection ===")

	// int8 overflow
	c.Set("int8_counter", int8(127), cache.NoExpiration)
	fmt.Println("int8_counter = 127 (max value)")

	_, err := c.IncrementInt8("int8_counter", 1)
	if err != nil {
		fmt.Printf("✓ Overflow prevented: %v\n", err)
	}

	// int8 underflow
	c.Set("int8_counter", int8(-128), cache.NoExpiration)
	fmt.Println("\nint8_counter = -128 (min value)")

	_, err = c.DecrementInt8("int8_counter", 1)
	if err != nil {
		fmt.Printf("✓ Underflow prevented: %v\n", err)
	}

	// int16 overflow
	c.Set("int16_counter", int16(math.MaxInt16), cache.NoExpiration)
	fmt.Printf("\nint16_counter = %d (max value)\n", math.MaxInt16)

	_, err = c.IncrementInt16("int16_counter", 1)
	if err != nil {
		fmt.Printf("✓ Overflow prevented: %v\n", err)
	}

	// int32 overflow
	c.Set("int32_counter", int32(math.MaxInt32), cache.NoExpiration)
	fmt.Printf("\nint32_counter = %d (max value)\n", math.MaxInt32)

	_, err = c.IncrementInt32("int32_counter", 1)
	if err != nil {
		fmt.Printf("✓ Overflow prevented: %v\n", err)
	}

	// int64 overflow
	c.Set("int64_counter", int64(math.MaxInt64), cache.NoExpiration)
	fmt.Printf("\nint64_counter = %d (max value)\n", math.MaxInt64)

	_, err = c.IncrementInt64("int64_counter", 1)
	if err != nil {
		fmt.Printf("✓ Overflow prevented: %v\n\n", err)
	}
}

func demonstrateFloatOverflow(c *cache.Cache) {
	fmt.Println("=== Float Infinity Protection ===")

	// float32 overflow to infinity
	c.Set("float32_counter", float32(math.MaxFloat32), cache.NoExpiration)
	fmt.Printf("float32_counter = %.2e (max value)\n", float32(math.MaxFloat32))

	_, err := c.IncrementFloat32("float32_counter", float32(math.MaxFloat32))
	if err != nil {
		fmt.Printf("✓ Infinity prevented: %v\n", err)
	}

	// float64 overflow to infinity
	c.Set("float64_counter", float64(math.MaxFloat64), cache.NoExpiration)
	fmt.Printf("\nfloat64_counter = %.2e (max value)\n", math.MaxFloat64)

	_, err = c.IncrementFloat64("float64_counter", math.MaxFloat64)
	if err != nil {
		fmt.Printf("✓ Infinity prevented: %v\n", err)
	}

	// float32 underflow to -infinity
	c.Set("float32_counter", float32(-math.MaxFloat32), cache.NoExpiration)
	fmt.Printf("\nfloat32_counter = %.2e (min value)\n", float32(-math.MaxFloat32))

	_, err = c.DecrementFloat32("float32_counter", float32(math.MaxFloat32))
	if err != nil {
		fmt.Printf("✓ -Infinity prevented: %v\n\n", err)
	}
}

func demonstrateSafeOperations(c *cache.Cache) {
	fmt.Println("=== Safe Operations Within Bounds ===")

	// uint8: 254 + 1 = 255 (within bounds)
	c.Set("safe_uint8", uint8(254), cache.NoExpiration)
	result, err := c.IncrementUint8("safe_uint8", 1)
	if err == nil {
		fmt.Printf("✓ uint8: 254 + 1 = %d (success)\n", result)
	}

	// int8: 126 + 1 = 127 (within bounds)
	c.Set("safe_int8", int8(126), cache.NoExpiration)
	result8, err := c.IncrementInt8("safe_int8", 1)
	if err == nil {
		fmt.Printf("✓ int8: 126 + 1 = %d (success)\n", result8)
	}

	// uint8: 1 - 1 = 0 (within bounds)
	c.Set("safe_uint8_dec", uint8(1), cache.NoExpiration)
	resultDec, err := c.DecrementUint8("safe_uint8_dec", 1)
	if err == nil {
		fmt.Printf("✓ uint8: 1 - 1 = %d (success)\n", resultDec)
	}

	// int8: -127 - 1 = -128 (within bounds)
	c.Set("safe_int8_dec", int8(-127), cache.NoExpiration)
	result8Dec, err := c.DecrementInt8("safe_int8_dec", 1)
	if err == nil {
		fmt.Printf("✓ int8: -127 - 1 = %d (success)\n", result8Dec)
	}

	// float64: safe addition
	c.Set("safe_float64", 100.5, cache.NoExpiration)
	resultFloat, err := c.IncrementFloat64("safe_float64", 50.25)
	if err == nil {
		fmt.Printf("✓ float64: 100.5 + 50.25 = %.2f (success)\n", resultFloat)
	}

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\nKey Takeaway: All numeric operations are protected against")
	fmt.Println("overflow/underflow, preventing silent data corruption!")
}
