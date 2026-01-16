package cache

import (
	"testing"
	"time"
)

func TestCache_Decrement(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	// Caso: Decrement int
	t.Run("Decrement int", func(t *testing.T) {
		cache.Set("intKey", 10, DefaultExpiration)
		err := cache.Decrement("intKey", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("intKey")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != 5 {
			t.Errorf("Expected value to be 5, got %v", value)
		}
	})

	// Caso: Decrement int8
	t.Run("Decrement int8", func(t *testing.T) {
		cache.Set("int8Key", int8(10), DefaultExpiration)
		err := cache.Decrement("int8Key", 3)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("int8Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != int8(7) {
			t.Errorf("Expected value to be 7, got %v", value)
		}
	})

	// Caso: Decrement int16
	t.Run("Decrement int16", func(t *testing.T) {
		cache.Set("int16Key", int16(100), DefaultExpiration)
		err := cache.Decrement("int16Key", 25)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("int16Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != int16(75) {
			t.Errorf("Expected value to be 75, got %v", value)
		}
	})

	// Caso: Decrement int32
	t.Run("Decrement int32", func(t *testing.T) {
		cache.Set("int32Key", int32(1000), DefaultExpiration)
		err := cache.Decrement("int32Key", 250)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("int32Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != int32(750) {
			t.Errorf("Expected value to be 750, got %v", value)
		}
	})

	// Caso: Decrement int64
	t.Run("Decrement int64", func(t *testing.T) {
		cache.Set("int64Key", int64(10000), DefaultExpiration)
		err := cache.Decrement("int64Key", 2500)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("int64Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != int64(7500) {
			t.Errorf("Expected value to be 7500, got %v", value)
		}
	})

	// Caso: Decrement uint
	t.Run("Decrement uint", func(t *testing.T) {
		cache.Set("uintKey", uint(50), DefaultExpiration)
		err := cache.Decrement("uintKey", 10)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("uintKey")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != uint(40) {
			t.Errorf("Expected value to be 40, got %v", value)
		}
	})

	// Caso: Decrement uintptr
	t.Run("Decrement uintptr", func(t *testing.T) {
		cache.Set("uintptrKey", uintptr(100), DefaultExpiration)
		err := cache.Decrement("uintptrKey", 30)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("uintptrKey")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != uintptr(70) {
			t.Errorf("Expected value to be 70, got %v", value)
		}
	})

	// Caso: Decrement uint8
	t.Run("Decrement uint8", func(t *testing.T) {
		cache.Set("uint8Key", uint8(50), DefaultExpiration)
		err := cache.Decrement("uint8Key", 15)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("uint8Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != uint8(35) {
			t.Errorf("Expected value to be 35, got %v", value)
		}
	})

	// Caso: Decrement uint16
	t.Run("Decrement uint16", func(t *testing.T) {
		cache.Set("uint16Key", uint16(500), DefaultExpiration)
		err := cache.Decrement("uint16Key", 100)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("uint16Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != uint16(400) {
			t.Errorf("Expected value to be 400, got %v", value)
		}
	})

	// Caso: Decrement uint32
	t.Run("Decrement uint32", func(t *testing.T) {
		cache.Set("uint32Key", uint32(5000), DefaultExpiration)
		err := cache.Decrement("uint32Key", 1000)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("uint32Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != uint32(4000) {
			t.Errorf("Expected value to be 4000, got %v", value)
		}
	})

	// Caso: Decrement uint64
	t.Run("Decrement uint64", func(t *testing.T) {
		cache.Set("uint64Key", uint64(50000), DefaultExpiration)
		err := cache.Decrement("uint64Key", 10000)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("uint64Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != uint64(40000) {
			t.Errorf("Expected value to be 40000, got %v", value)
		}
	})

	// Caso: Decrement float32
	t.Run("Decrement float32", func(t *testing.T) {
		cache.Set("float32Key", float32(10.5), DefaultExpiration)
		err := cache.Decrement("float32Key", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("float32Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != float32(10.5-5) {
			t.Errorf("Expected value to be %v, got %v", float32(10.5-5), value)
		}
	})

	// Caso: Decrement float64
	t.Run("Decrement float64", func(t *testing.T) {
		cache.Set("float64Key", float64(100.75), DefaultExpiration)
		err := cache.Decrement("float64Key", 25)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("float64Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != float64(75.75) {
			t.Errorf("Expected value to be %v, got %v", float64(75.75), value)
		}
	})

	// Caso: Decrement int not found
	t.Run("Decrement int not found", func(t *testing.T) {
		err := cache.Decrement("notFoundKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// Caso: Decrement int with expiration
	t.Run("Decrement int with expiration", func(t *testing.T) {
		cache.Set("intExpKey", 10, 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		err := cache.Decrement("intExpKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// Caso: Decrement non-integer type
	t.Run("Decrement non-integer type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		err := cache.Decrement("stringKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementFloat(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	t.Run("Decrement float32", func(t *testing.T) {
		cache.Set("float32Key", float32(10.5), DefaultExpiration)
		err := cache.DecrementFloat("float32Key", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("float32Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != float32(10.5-5) {
			t.Errorf("Expected value to be %v, got %v", float32(10.5-5), value)
		}
	})

	t.Run("Decrement float64", func(t *testing.T) {
		cache.Set("float64Key", float64(20.5), DefaultExpiration)
		err := cache.DecrementFloat("float64Key", 10.5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		value, found := cache.Get("float64Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != float64(20.5-10.5) {
			t.Errorf("Expected value to be %v, got %v", float64(20.5-10.5), value)
		}
	})

	t.Run("Decrement float not found", func(t *testing.T) {
		err := cache.DecrementFloat("notFoundKey", 5.0)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement float with expiration", func(t *testing.T) {
		cache.Set("floatExpKey", float64(20.5), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		err := cache.DecrementFloat("floatExpKey", 10.5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement non-float type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		err := cache.DecrementFloat("stringKey", 5.0)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementInt(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	t.Run("Decrement int", func(t *testing.T) {
		cache.Set("intKey", 10, DefaultExpiration)
		val, err := cache.DecrementInt("intKey", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected value to be %v, got %v", 5, val)
		}
		value, found := cache.Get("intKey")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != 5 {
			t.Errorf("Expected value to be %v, got %v", 5, value)
		}
	})

	t.Run("Decrement int not found", func(t *testing.T) {
		_, err := cache.DecrementInt("notFoundKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement int with expiration", func(t *testing.T) {
		cache.Set("intExpKey", 20, 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementInt("intExpKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement non-int type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementInt("stringKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementInt8(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	t.Run("Decrement int8", func(t *testing.T) {
		cache.Set("int8Key", int8(10), DefaultExpiration)
		val, err := cache.DecrementInt8("int8Key", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected value to be %v, got %v", 5, val)
		}
		value, found := cache.Get("int8Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != int8(5) {
			t.Errorf("Expected value to be %v, got %v", int8(5), value)
		}
	})

	t.Run("Decrement int8 not found", func(t *testing.T) {
		_, err := cache.DecrementInt8("notFoundKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement int8 with expiration", func(t *testing.T) {
		cache.Set("int8ExpKey", int8(20), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementInt8("int8ExpKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement non-int8 type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementInt8("stringKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementInt16(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	t.Run("Decrement int16", func(t *testing.T) {
		cache.Set("int16Key", int16(10), DefaultExpiration)
		val, err := cache.DecrementInt16("int16Key", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected value to be %v, got %v", 5, val)
		}
		value, found := cache.Get("int16Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != int16(5) {
			t.Errorf("Expected value to be %v, got %v", int16(5), value)
		}
	})

	t.Run("Decrement int16 not found", func(t *testing.T) {
		_, err := cache.DecrementInt16("notFoundKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement int16 with expiration", func(t *testing.T) {
		cache.Set("int16ExpKey", int16(20), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementInt16("int16ExpKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement non-int16 type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementInt16("stringKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementInt32(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	t.Run("Decrement int32", func(t *testing.T) {
		cache.Set("int32Key", int32(10), DefaultExpiration)
		val, err := cache.DecrementInt32("int32Key", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected value to be %v, got %v", 5, val)
		}
		value, found := cache.Get("int32Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != int32(5) {
			t.Errorf("Expected value to be %v, got %v", int32(5), value)
		}
	})

	t.Run("Decrement int32 not found", func(t *testing.T) {
		_, err := cache.DecrementInt32("notFoundKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement int32 with expiration", func(t *testing.T) {
		cache.Set("int32ExpKey", int32(20), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementInt32("int32ExpKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement non-int32 type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementInt32("stringKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementInt64(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	t.Run("Decrement int64", func(t *testing.T) {
		cache.Set("int64Key", int64(10), DefaultExpiration)
		val, err := cache.DecrementInt64("int64Key", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected value to be %v, got %v", 5, val)
		}
		value, found := cache.Get("int64Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != int64(5) {
			t.Errorf("Expected value to be %v, got %v", int64(5), value)
		}
	})

	t.Run("Decrement int64 not found", func(t *testing.T) {
		_, err := cache.DecrementInt64("notFoundKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement int64 with expiration", func(t *testing.T) {
		cache.Set("int64ExpKey", int64(20), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementInt64("int64ExpKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement non-int64 type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementInt64("stringKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementUint(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	// Caso: Decrement uint
	t.Run("Decrement uint", func(t *testing.T) {
		cache.Set("uintKey", uint(10), DefaultExpiration)
		val, err := cache.DecrementUint("uintKey", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected value to be %v, got %v", 5, val)
		}
		value, found := cache.Get("uintKey")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != uint(5) {
			t.Errorf("Expected value to be %v, got %v", uint(5), value)
		}
	})

	// Caso: Decrement uint not found
	t.Run("Decrement uint not found", func(t *testing.T) {
		_, err := cache.DecrementUint("notFoundKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// Caso: Decrement uint with expiration
	t.Run("Decrement uint with expiration", func(t *testing.T) {
		cache.Set("uintExpKey", uint(20), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementUint("uintExpKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// Caso: Decrement non-uint type
	t.Run("Decrement non-uint type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementUint("stringKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementUintptr(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	// Caso: Decrement uintptr
	t.Run("Decrement uintptr", func(t *testing.T) {
		cache.Set("uintptrKey", uintptr(10), DefaultExpiration)
		val, err := cache.DecrementUintptr("uintptrKey", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected value to be %v, got %v", 5, val)
		}
		value, found := cache.Get("uintptrKey")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != uintptr(5) {
			t.Errorf("Expected value to be %v, got %v", uintptr(5), value)
		}
	})

	// Caso: Decrement uintptr not found
	t.Run("Decrement uintptr not found", func(t *testing.T) {
		_, err := cache.DecrementUintptr("notFoundKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// Caso: Decrement uintptr with expiration
	t.Run("Decrement uintptr with expiration", func(t *testing.T) {
		cache.Set("uintptrExpKey", uintptr(20), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementUintptr("uintptrExpKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// Caso: Decrement non-uintptr type
	t.Run("Decrement non-uintptr type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementUintptr("stringKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementUint8(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	// Caso: Decrement uint8
	t.Run("Decrement uint8", func(t *testing.T) {
		cache.Set("uint8Key", uint8(10), DefaultExpiration)
		val, err := cache.DecrementUint8("uint8Key", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected value to be %v, got %v", 5, val)
		}
		value, found := cache.Get("uint8Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != uint8(5) {
			t.Errorf("Expected value to be %v, got %v", uint8(5), value)
		}
	})

	t.Run("Decrement uint8 not found", func(t *testing.T) {
		_, err := cache.DecrementUint8("notFoundKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement uint8 with expiration", func(t *testing.T) {
		cache.Set("uint8ExpKey", uint8(20), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementUint8("uint8ExpKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement non-uint8 type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementUint8("stringKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementUint16(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	t.Run("Decrement uint16", func(t *testing.T) {
		cache.Set("uint16Key", uint16(10), DefaultExpiration)
		val, err := cache.DecrementUint16("uint16Key", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected value to be %v, got %v", 5, val)
		}
		value, found := cache.Get("uint16Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != uint16(5) {
			t.Errorf("Expected value to be %v, got %v", uint16(5), value)
		}
	})

	t.Run("Decrement uint16 not found", func(t *testing.T) {
		_, err := cache.DecrementUint16("notFoundKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement uint16 with expiration", func(t *testing.T) {
		cache.Set("uint16ExpKey", uint16(20), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementUint16("uint16ExpKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement non-uint16 type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementUint16("stringKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementUint32(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	t.Run("Decrement uint32", func(t *testing.T) {
		cache.Set("uint32Key", uint32(10), DefaultExpiration)
		val, err := cache.DecrementUint32("uint32Key", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected value to be %v, got %v", 5, val)
		}
		value, found := cache.Get("uint32Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != uint32(5) {
			t.Errorf("Expected value to be %v, got %v", uint32(5), value)
		}
	})

	t.Run("Decrement uint32 not found", func(t *testing.T) {
		_, err := cache.DecrementUint32("notFoundKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement uint32 with expiration", func(t *testing.T) {
		cache.Set("uint32ExpKey", uint32(20), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementUint32("uint32ExpKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement non-uint32 type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementUint32("stringKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementUint64(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	t.Run("Decrement uint64", func(t *testing.T) {
		cache.Set("uint64Key", uint64(10), DefaultExpiration)
		val, err := cache.DecrementUint64("uint64Key", 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected value to be %v, got %v", 5, val)
		}
		value, found := cache.Get("uint64Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != uint64(5) {
			t.Errorf("Expected value to be %v, got %v", uint64(5), value)
		}
	})

	t.Run("Decrement uint64 not found", func(t *testing.T) {
		_, err := cache.DecrementUint64("notFoundKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement uint64 with expiration", func(t *testing.T) {
		cache.Set("uint64ExpKey", uint64(20), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementUint64("uint64ExpKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement non-uint64 type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementUint64("stringKey", 5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementFloat32(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	t.Run("Decrement float32", func(t *testing.T) {
		cache.Set("float32Key", float32(10.5), DefaultExpiration)
		val, err := cache.DecrementFloat32("float32Key", 2.5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != float32(8.0) {
			t.Errorf("Expected value to be %v, got %v", float32(8.0), val)
		}
		value, found := cache.Get("float32Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != float32(8.0) {
			t.Errorf("Expected value to be %v, got %v", float32(8.0), value)
		}
	})

	t.Run("Decrement float32 not found", func(t *testing.T) {
		_, err := cache.DecrementFloat32("notFoundKey", 2.5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement float32 with expiration", func(t *testing.T) {
		cache.Set("float32ExpKey", float32(20.5), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementFloat32("float32ExpKey", 2.5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement non-float32 type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementFloat32("stringKey", 2.5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementFloat64(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	t.Run("Decrement float64", func(t *testing.T) {
		cache.Set("float64Key", float64(10.5), DefaultExpiration)
		val, err := cache.DecrementFloat64("float64Key", 2.5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if val != float64(8.0) {
			t.Errorf("Expected value to be %v, got %v", float64(8.0), val)
		}
		value, found := cache.Get("float64Key")
		if !found {
			t.Fatalf("Expected to find key")
		}
		if value != float64(8.0) {
			t.Errorf("Expected value to be %v, got %v", float64(8.0), value)
		}
	})

	// Caso: Decrement float64 not found
	t.Run("Decrement float64 not found", func(t *testing.T) {
		_, err := cache.DecrementFloat64("notFoundKey", 2.5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement float64 with expiration", func(t *testing.T) {
		cache.Set("float64ExpKey", float64(20.5), 100*time.Millisecond)
		time.Sleep(200 * time.Millisecond) // Espera para que expire
		_, err := cache.DecrementFloat64("float64ExpKey", 2.5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Decrement non-float64 type", func(t *testing.T) {
		cache.Set("stringKey", "value", DefaultExpiration)
		_, err := cache.DecrementFloat64("stringKey", 2.5)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func TestCache_DecrementTyped_Underflow(t *testing.T) {
	cache := New(DefaultExpiration, 1*time.Minute)

	t.Run("DecrementInt underflow with positive n", func(t *testing.T) {
		cache.Set("intKey", -9223372036854775808, DefaultExpiration) // math.MinInt64
		_, err := cache.DecrementInt("intKey", 1)
		if err == nil {
			t.Fatalf("Expected underflow error, got nil")
		}
	})

	t.Run("DecrementInt8 underflow with positive n", func(t *testing.T) {
		cache.Set("int8Key", int8(-128), DefaultExpiration) // math.MinInt8
		_, err := cache.DecrementInt8("int8Key", 1)
		if err == nil {
			t.Fatalf("Expected underflow error, got nil")
		}
	})

	t.Run("DecrementInt16 underflow with positive n", func(t *testing.T) {
		cache.Set("int16Key", int16(-32768), DefaultExpiration) // math.MinInt16
		_, err := cache.DecrementInt16("int16Key", 1)
		if err == nil {
			t.Fatalf("Expected underflow error, got nil")
		}
	})

	t.Run("DecrementInt32 underflow with positive n", func(t *testing.T) {
		cache.Set("int32Key", int32(-2147483648), DefaultExpiration) // math.MinInt32
		_, err := cache.DecrementInt32("int32Key", 1)
		if err == nil {
			t.Fatalf("Expected underflow error, got nil")
		}
	})

	t.Run("DecrementInt64 underflow with positive n", func(t *testing.T) {
		cache.Set("int64Key", int64(-9223372036854775808), DefaultExpiration) // math.MinInt64
		_, err := cache.DecrementInt64("int64Key", 1)
		if err == nil {
			t.Fatalf("Expected underflow error, got nil")
		}
	})

	t.Run("DecrementInt overflow with negative n", func(t *testing.T) {
		cache.Set("intKey", 9223372036854775807, DefaultExpiration) // math.MaxInt64
		_, err := cache.DecrementInt("intKey", -1)
		if err == nil {
			t.Fatalf("Expected overflow error, got nil")
		}
	})

	t.Run("DecrementInt8 overflow with negative n", func(t *testing.T) {
		cache.Set("int8Key", int8(127), DefaultExpiration) // math.MaxInt8
		_, err := cache.DecrementInt8("int8Key", -1)
		if err == nil {
			t.Fatalf("Expected overflow error, got nil")
		}
	})

	t.Run("DecrementInt16 overflow with negative n", func(t *testing.T) {
		cache.Set("int16Key", int16(32767), DefaultExpiration) // math.MaxInt16
		_, err := cache.DecrementInt16("int16Key", -1)
		if err == nil {
			t.Fatalf("Expected overflow error, got nil")
		}
	})

	t.Run("DecrementInt32 overflow with negative n", func(t *testing.T) {
		cache.Set("int32Key", int32(2147483647), DefaultExpiration) // math.MaxInt32
		_, err := cache.DecrementInt32("int32Key", -1)
		if err == nil {
			t.Fatalf("Expected overflow error, got nil")
		}
	})

	t.Run("DecrementInt64 overflow with negative n", func(t *testing.T) {
		cache.Set("int64Key", int64(9223372036854775807), DefaultExpiration) // math.MaxInt64
		_, err := cache.DecrementInt64("int64Key", -1)
		if err == nil {
			t.Fatalf("Expected overflow error, got nil")
		}
	})

	t.Run("DecrementUint underflow", func(t *testing.T) {
		cache.Set("uintKey", uint(5), DefaultExpiration)
		_, err := cache.DecrementUint("uintKey", 10)
		if err == nil {
			t.Fatalf("Expected underflow error, got nil")
		}
	})

	t.Run("DecrementUintptr underflow", func(t *testing.T) {
		cache.Set("uintptrKey", uintptr(3), DefaultExpiration)
		_, err := cache.DecrementUintptr("uintptrKey", 5)
		if err == nil {
			t.Fatalf("Expected underflow error, got nil")
		}
	})

	t.Run("DecrementUint8 underflow", func(t *testing.T) {
		cache.Set("uint8Key", uint8(2), DefaultExpiration)
		_, err := cache.DecrementUint8("uint8Key", 5)
		if err == nil {
			t.Fatalf("Expected underflow error, got nil")
		}
	})

	t.Run("DecrementUint16 underflow", func(t *testing.T) {
		cache.Set("uint16Key", uint16(10), DefaultExpiration)
		_, err := cache.DecrementUint16("uint16Key", 20)
		if err == nil {
			t.Fatalf("Expected underflow error, got nil")
		}
	})

	t.Run("DecrementUint32 underflow", func(t *testing.T) {
		cache.Set("uint32Key", uint32(50), DefaultExpiration)
		_, err := cache.DecrementUint32("uint32Key", 100)
		if err == nil {
			t.Fatalf("Expected underflow error, got nil")
		}
	})

	t.Run("DecrementUint64 underflow", func(t *testing.T) {
		cache.Set("uint64Key", uint64(100), DefaultExpiration)
		_, err := cache.DecrementUint64("uint64Key", 200)
		if err == nil {
			t.Fatalf("Expected underflow error, got nil")
		}
	})
}
