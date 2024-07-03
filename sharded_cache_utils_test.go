package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDjb33(t *testing.T) {
	t.Run("Hash function is consistent", func(t *testing.T) {
		result1 := djb33(0, "test")
		result2 := djb33(0, "test")
		assert.Equal(t, result1, result2)
	})

	t.Run("Hash function with different seeds", func(t *testing.T) {
		result1 := djb33(1, "test")
		result2 := djb33(2, "test")
		assert.NotEqual(t, result1, result2)
	})

	t.Run("Hash function with different strings", func(t *testing.T) {
		result1 := djb33(0, "test")
		result2 := djb33(0, "another_test")
		assert.NotEqual(t, result1, result2)
	})

	t.Run("Hash function with same strings and different seeds", func(t *testing.T) {
		result1 := djb33(1, "same_string")
		result2 := djb33(2, "same_string")
		assert.NotEqual(t, result1, result2)
	})

	t.Run("Hash function is different for different inputs", func(t *testing.T) {
		result1 := djb33(0, "input1")
		result2 := djb33(0, "input2")
		assert.NotEqual(t, result1, result2)
	})
}
