package main

import (
	"fmt"
	"time"

	"github.com/pzentenoe/go-cache"
)

func main() {
	fmt.Println("=== Basic Cache Example ===")

	// Create a cache with a default expiration time of 5 minutes,
	// and which purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)

	// Set the value of the key "foo" to "bar", with the default expiration time
	c.Set("foo", "bar", cache.DefaultExpiration)
	fmt.Println("Set 'foo' = 'bar'")

	// Get the value associated with the key "foo"
	if val, found := c.Get("foo"); found {
		fmt.Printf("Get 'foo' = %v\n\n", val)
	}

	// Set with custom expiration
	fmt.Println("Setting 'temp' with 2 second expiration...")
	c.Set("temp", "expires soon", 2*time.Second)

	// Get immediately
	if val, found := c.Get("temp"); found {
		fmt.Printf("Immediately after set: 'temp' = %v\n", val)
	}

	// Wait for expiration
	fmt.Println("Waiting 3 seconds for expiration...")
	time.Sleep(3 * time.Second)

	// Try to get after expiration
	if _, found := c.Get("temp"); !found {
		fmt.Println("After expiration: 'temp' not found (expired)")
	}

	// Set with no expiration
	c.Set("permanent", "never expires", cache.NoExpiration)
	fmt.Println("Set 'permanent' with NoExpiration")

	// Use SetDefault (uses cache's default expiration)
	c.SetDefault("user:123", map[string]string{
		"name":  "John Doe",
		"email": "john@example.com",
	})
	fmt.Println("Set user data with SetDefault")

	// Get with expiration info
	if val, expiration, found := c.GetWithExpiration("user:123"); found {
		fmt.Printf("User data: %v\n", val)
		if !expiration.IsZero() {
			fmt.Printf("Expires at: %v\n", expiration.Format(time.RFC3339))
		}
	}

	// Item count
	fmt.Printf("\nTotal items in cache: %d\n", c.ItemCount())

	// Add (only if not exists)
	err := c.Add("foo", "new value", cache.DefaultExpiration)
	if err != nil {
		fmt.Printf("Add 'foo' failed: %v (already exists)\n", err)
	}

	err = c.Add("new_key", "new value", cache.DefaultExpiration)
	if err == nil {
		fmt.Println("Add 'new_key' succeeded")
	}

	// Replace (only if exists)
	err = c.Replace("foo", "replaced value", cache.DefaultExpiration)
	if err == nil {
		fmt.Println("Replace 'foo' succeeded")
		if val, found := c.Get("foo"); found {
			fmt.Printf("New value: %v\n", val)
		}
	}

	// Delete
	c.Delete("new_key")
	fmt.Println("\nDeleted 'new_key'")

	// Flush all
	fmt.Printf("Items before flush: %d\n", c.ItemCount())
	c.Flush()
	fmt.Printf("Items after flush: %d\n", c.ItemCount())

	fmt.Println("\n=== Example Complete ===")
}
