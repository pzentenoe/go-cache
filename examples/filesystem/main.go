package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pzentenoe/go-cache"
)

const cacheFile = "app_cache.gob"

func main() {
	fmt.Println("=== Filesystem Cache Example ===")

	// 1. Warm start: create cache from file if it exists
	c := warmStart()

	// 2. Populate with new data
	populateCache(c)

	// 3. Periodic save in background
	stopSaver := startPeriodicSave(c, 10*time.Second)

	// 4. Graceful shutdown on signal
	setupGracefulShutdown(c, stopSaver)

	// 5. Simulate application work
	simulateWork(c)

	// 6. Final save and cleanup
	shutdown(c, stopSaver)
}

// warmStart loads cache from disk if available, otherwise starts fresh.
func warmStart() *cache.Cache {
	fmt.Println("\n--- Warm Start ---")

	c := cache.New(5*time.Minute, 1*time.Minute)

	if _, err := os.Stat(cacheFile); err == nil {
		fmt.Printf("Found existing cache file: %s\n", cacheFile)
		if err := c.LoadFile(cacheFile); err != nil {
			log.Printf("Warning: could not load cache from %s: %v (starting fresh)", cacheFile, err)
		} else {
			fmt.Printf("Loaded %d items from disk\n", c.ItemCount())
			return c
		}
	} else {
		fmt.Println("No cache file found, starting with empty cache")
	}

	return c
}

// populateCache adds sample data to the cache.
func populateCache(c *cache.Cache) {
	fmt.Println("\n--- Populating Cache ---")

	c.Set("app:version", "1.3.0", cache.NoExpiration)
	c.Set("user:1", map[string]string{"name": "Alice", "role": "admin"}, 10*time.Minute)
	c.Set("user:2", map[string]string{"name": "Bob", "role": "user"}, 10*time.Minute)
	c.Set("counter:visits", int64(0), cache.NoExpiration)

	fmt.Printf("Cache now has %d items\n", c.ItemCount())
}

// startPeriodicSave saves cache to disk at regular intervals.
// Returns a stop channel to signal the goroutine to exit.
func startPeriodicSave(c *cache.Cache, interval time.Duration) chan struct{} {
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := atomicSave(c); err != nil {
					log.Printf("Periodic save failed: %v", err)
				} else {
					fmt.Printf("[periodic] Saved %d items to %s\n", c.ItemCount(), cacheFile)
				}
			case <-stop:
				return
			}
		}
	}()
	fmt.Printf("Periodic save started (every %v)\n", interval)
	return stop
}

// atomicSave writes to a temp file first, then renames for crash safety.
func atomicSave(c *cache.Cache) error {
	tmpFile := cacheFile + ".tmp"
	if err := c.SaveFile(tmpFile); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("save to temp file: %w", err)
	}
	if err := os.Rename(tmpFile, cacheFile); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

// setupGracefulShutdown handles OS signals for clean shutdown.
func setupGracefulShutdown(c *cache.Cache, stopSaver chan struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal: %v\n", sig)
		shutdown(c, stopSaver)
		os.Exit(0)
	}()
}

// simulateWork simulates application activity.
func simulateWork(c *cache.Cache) {
	fmt.Println("\n--- Simulating Application Work ---")
	fmt.Println("(Press Ctrl+C for graceful shutdown, or wait for auto-completion)")

	for i := 0; i < 5; i++ {
		c.Increment("counter:visits", 1)

		if val, found := c.Get("counter:visits"); found {
			fmt.Printf("  Visit #%v processed\n", val)
		}

		time.Sleep(1 * time.Second)
	}
}

// shutdown performs clean shutdown: save to disk and stop janitor.
func shutdown(c *cache.Cache, stopSaver chan struct{}) {
	fmt.Println("\n--- Shutdown ---")

	// Stop periodic saver
	select {
	case stopSaver <- struct{}{}:
	default:
	}

	// Final save
	c.DeleteExpired()
	if err := atomicSave(c); err != nil {
		log.Printf("Final save failed: %v", err)
	} else {
		fmt.Printf("Final save: %d items written to %s\n", c.ItemCount(), cacheFile)
	}

	// Stop janitor
	c.Close()
	fmt.Println("Cache closed")

	// Cleanup demo file
	os.Remove(cacheFile)
	fmt.Println("Demo file cleaned up")

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\nPatterns demonstrated:")
	fmt.Println("  - Warm start (load from disk if available)")
	fmt.Println("  - Atomic save (write to temp, then rename)")
	fmt.Println("  - Periodic background save")
	fmt.Println("  - Graceful shutdown with final save")
	fmt.Println("  - Explicit Close() for janitor cleanup")
}
