package cache

import (
	"encoding/gob"
	"errors"
	"io"
	"os"
	"reflect"
)

// registerGobTypes registers each distinct value type with Gob once.
// Registering the same type twice is a no-op, so deduping by type avoids
// repeated calls through Gob's global lock when the cache is large.
func registerGobTypes(items map[string]Item, registered map[reflect.Type]struct{}) {
	for _, v := range items {
		t := reflect.TypeOf(v.Object)
		if _, ok := registered[t]; ok {
			continue
		}
		registered[t] = struct{}{}
		gob.Register(v.Object)
	}
}

// Save Write the cache's items (using Gob) to an io.Writer.
func (c *Cache) Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = errors.New(errGobRegistration)
		}
	}()
	c.mu.RLock()
	defer c.mu.RUnlock()
	registerGobTypes(c.items, make(map[reflect.Type]struct{}))
	err = enc.Encode(&c.items)
	return
}

// SaveFile Save the cache's items to the given filename, creating the file if it
// doesn't exist, and overwriting it if it does.
func (c *Cache) SaveFile(fname string) error {
	fp, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer fp.Close()
	return c.Save(fp)
}

// Load Add (Gob-serialized) cache items from an io.Reader, excluding any items with
// keys that already exist (and haven't expired) in the current cache.
func (c *Cache) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	items := map[string]Item{}
	if err := dec.Decode(&items); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range items {
		if ov, found := c.items[k]; !found || ov.Expired() {
			c.items[k] = v
		}
	}
	return nil
}

// LoadFile Load and add cache items from the given filename, excluding any items with
// keys that already exist in the current cache.
func (c *Cache) LoadFile(fname string) error {
	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fp.Close()
	return c.Load(fp)
}
