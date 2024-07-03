package cache

import (
	"bytes"
	"encoding/gob"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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
}
