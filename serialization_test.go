package cache

import (
	"bytes"
	"encoding/gob"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache_Save(t *testing.T) {
	t.Run("Save cache items to io.Writer", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key1", "value1", NoExpiration)
		c.Set("key2", "value2", NoExpiration)

		var buf bytes.Buffer
		err := c.Save(&buf)

		assert.NoError(t, err)

		// Decode the saved data to verify
		dec := gob.NewDecoder(&buf)
		items := map[string]Item{}
		err = dec.Decode(&items)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(items))
		assert.Equal(t, "value1", items["key1"].Object)
		assert.Equal(t, "value2", items["key2"].Object)
	})
}

func TestCache_SaveFile(t *testing.T) {
	t.Run("Save cache items to file", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key1", "value1", NoExpiration)
		c.Set("key2", "value2", NoExpiration)

		fileName := "test_cache_save.gob"
		defer os.Remove(fileName) // Clean up

		err := c.SaveFile(fileName)
		assert.NoError(t, err)

		// Load the file and verify its contents
		file, err := os.Open(fileName)
		assert.NoError(t, err)
		defer file.Close()

		dec := gob.NewDecoder(file)
		items := map[string]Item{}
		err = dec.Decode(&items)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(items))
		assert.Equal(t, "value1", items["key1"].Object)
		assert.Equal(t, "value2", items["key2"].Object)
	})

	t.Run("SaveFile with invalid path", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key1", "value1", NoExpiration)

		// Try to save to an invalid directory
		err := c.SaveFile("/invalid/path/to/file.gob")
		assert.Error(t, err)
	})
}

func TestCache_Load(t *testing.T) {
	t.Run("Load cache items from io.Reader", func(t *testing.T) {
		c := New(DefaultExpiration, 0)

		// Prepare data to load
		items := map[string]Item{
			"key1": {Object: "value1", Expiration: 0},
			"key2": {Object: "value2", Expiration: 0},
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(&items)
		assert.NoError(t, err)

		err = c.Load(&buf)
		assert.NoError(t, err)

		// Verify loaded data
		val, found := c.Get("key1")
		assert.True(t, found)
		assert.Equal(t, "value1", val)

		val, found = c.Get("key2")
		assert.True(t, found)
		assert.Equal(t, "value2", val)
	})

	t.Run("Load with invalid data", func(t *testing.T) {
		c := New(DefaultExpiration, 0)

		// Try to load invalid data
		var buf bytes.Buffer
		buf.WriteString("invalid gob data")

		err := c.Load(&buf)
		assert.Error(t, err)
	})

	t.Run("Load does not overwrite existing non-expired items", func(t *testing.T) {
		c := New(DefaultExpiration, 0)
		c.Set("key1", "existing_value", NoExpiration)

		// Prepare data to load with the same key
		items := map[string]Item{
			"key1": {Object: "new_value", Expiration: 0},
			"key2": {Object: "value2", Expiration: 0},
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(&items)
		assert.NoError(t, err)

		err = c.Load(&buf)
		assert.NoError(t, err)

		// Verify that key1 still has the original value
		val, found := c.Get("key1")
		assert.True(t, found)
		assert.Equal(t, "existing_value", val)

		// Verify that key2 was loaded
		val, found = c.Get("key2")
		assert.True(t, found)
		assert.Equal(t, "value2", val)
	})
}

func TestCache_LoadFile(t *testing.T) {
	t.Run("Load cache items from file", func(t *testing.T) {
		c := New(DefaultExpiration, 0)

		// Prepare data to load
		items := map[string]Item{
			"key1": {Object: "value1", Expiration: 0},
			"key2": {Object: "value2", Expiration: 0},
		}
		fileName := "test_cache_load.gob"
		defer os.Remove(fileName) // Clean up

		file, err := os.Create(fileName)
		assert.NoError(t, err)

		enc := gob.NewEncoder(file)
		err = enc.Encode(&items)
		assert.NoError(t, err)
		file.Close()

		err = c.LoadFile(fileName)
		assert.NoError(t, err)

		// Verify loaded data
		val, found := c.Get("key1")
		assert.True(t, found)
		assert.Equal(t, "value1", val)

		val, found = c.Get("key2")
		assert.True(t, found)
		assert.Equal(t, "value2", val)
	})

	t.Run("LoadFile with non-existent file", func(t *testing.T) {
		c := New(DefaultExpiration, 0)

		err := c.LoadFile("nonexistent_file.gob")
		assert.Error(t, err)
	})
}

// TestCache_SaveGobPanic covers the recover path: gob.Register panics on
// unregisterable value types (e.g. funcs) and Save must surface it as an
// error instead of crashing.
func TestCache_SaveGobPanic(t *testing.T) {
	c := New(NoExpiration, 0)
	c.Set("fn", func() {}, NoExpiration)

	var buf bytes.Buffer
	err := c.Save(&buf)
	assert.Error(t, err)
}

// FuzzCacheLoad feeds arbitrary data into Load: it must return an error or
// succeed, but never panic. Run with: go test -fuzz=FuzzCacheLoad -fuzztime=30s
func FuzzCacheLoad(f *testing.F) {
	c := New(NoExpiration, 0)
	c.Set("key", "value", NoExpiration)
	var buf bytes.Buffer
	if err := c.Save(&buf); err != nil {
		f.Fatal(err)
	}
	f.Add(buf.Bytes())
	f.Add([]byte{})
	f.Add([]byte("not a gob stream"))
	f.Fuzz(func(t *testing.T, data []byte) {
		c := New(NoExpiration, 0)
		_ = c.Load(bytes.NewReader(data))
	})
}

// FuzzShardedCacheLoad is the sharded counterpart of FuzzCacheLoad.
func FuzzShardedCacheLoad(f *testing.F) {
	sc := newShardedCache(2, NoExpiration)
	sc.Set("key", "value", NoExpiration)
	var buf bytes.Buffer
	if err := sc.Save(&buf); err != nil {
		f.Fatal(err)
	}
	f.Add(buf.Bytes())
	f.Add([]byte{})
	f.Add([]byte("not a gob stream"))
	f.Fuzz(func(t *testing.T, data []byte) {
		sc := newShardedCache(2, NoExpiration)
		_ = sc.Load(bytes.NewReader(data))
	})
}
