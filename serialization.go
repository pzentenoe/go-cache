package cache

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

// Save Write the cache's items (using Gob) to an io.Writer.
func (c *Cache) Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, v := range c.items {
		gob.Register(v.Object)
	}
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
	if err := c.Save(fp); err != nil {
		return err
	}
	return nil
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
