package cache

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSentinelErrors verifies that returned errors match the exported
// sentinels via errors.Is while keeping the historical message text.
func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name     string
		run      func(c *Cache) error
		sentinel error
		legacy   string
	}{
		{
			name:     "Increment on missing key wraps ErrNotFound",
			run:      func(c *Cache) error { return c.Increment("k", 1) },
			sentinel: ErrNotFound,
			legacy:   "item k not found",
		},
		{
			name: "Add on existing key wraps ErrAlreadyExists",
			run: func(c *Cache) error {
				c.Set("k", 1, NoExpiration)
				return c.Add("k", 2, NoExpiration)
			},
			sentinel: ErrAlreadyExists,
			legacy:   "item k already exists",
		},
		{
			name:     "Replace on missing key wraps ErrNotFound",
			run:      func(c *Cache) error { return c.Replace("k", 1, NoExpiration) },
			sentinel: ErrNotFound,
			legacy:   "item k doesn't exist",
		},
		{
			name: "IncrementInt8 overflow wraps ErrOverflow",
			run: func(c *Cache) error {
				c.Set("k", int8(100), NoExpiration)
				_, err := c.IncrementInt8("k", 100)
				return err
			},
			sentinel: ErrOverflow,
			legacy:   "overflow would occur for k",
		},
		{
			name: "DecrementUint8 underflow wraps ErrUnderflow",
			run: func(c *Cache) error {
				c.Set("k", uint8(5), NoExpiration)
				_, err := c.DecrementUint8("k", 10)
				return err
			},
			sentinel: ErrUnderflow,
			legacy:   "underflow would occur for k",
		},
		{
			name: "Typed op on mismatched type wraps ErrTypeMismatch",
			run: func(c *Cache) error {
				c.Set("k", int32(5), NoExpiration)
				_, err := c.IncrementInt64("k", 1)
				return err
			},
			sentinel: ErrTypeMismatch,
			legacy:   "the value for k does not match the type expected by this operation",
		},
		{
			name: "Increment on non-numeric value wraps ErrTypeMismatch",
			run: func(c *Cache) error {
				c.Set("k", "nope", NoExpiration)
				return c.Increment("k", 1)
			},
			sentinel: ErrTypeMismatch,
			legacy:   "the value for k is not an integer or float",
		},
		{
			name: "IncrementFloat on int value wraps ErrTypeMismatch",
			run: func(c *Cache) error {
				c.Set("k", 5, NoExpiration)
				return c.IncrementFloat("k", 1.5)
			},
			sentinel: ErrTypeMismatch,
			legacy:   "the value for k does not have type float32 or float64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(NoExpiration, 0)
			err := tt.run(c)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, tt.sentinel), "errors.Is(%v, %v)", err, tt.sentinel)
			assert.Equal(t, tt.legacy, err.Error(), "legacy message must not change")
		})
	}
}
