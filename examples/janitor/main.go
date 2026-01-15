package main

import (
	"fmt"
	"time"

	"github.com/pzentenoe/go-cache"
)

func main() {
	fmt.Println("=== Janitor Control Example ===\n")

	// Demonstrate basic janitor functionality
	demonstrateBasicJanitor()

	// Demonstrate pause/resume
	demonstratePauseResume()

	// Demonstrate interval change
	demonstrateIntervalChange()

	// Demonstrate combined control
	demonstrateCombinedControl()
}

func demonstrateBasicJanitor() {
	fmt.Println("=== Basic Janitor (Automatic Cleanup) ===")

	// Create cache with 100ms expiration and 50ms cleanup interval
	c := cache.New(100*time.Millisecond, 50*time.Millisecond)

	// Add items
	for i := 0; i < 10; i++ {
		c.Set(fmt.Sprintf("key%d", i), i, cache.DefaultExpiration)
	}

	fmt.Printf("Added 10 items, cache size: %d\n", c.ItemCount())

	// Wait for expiration
	fmt.Println("Waiting for items to expire...")
	time.Sleep(150 * time.Millisecond)

	// Janitor should have cleaned them up
	fmt.Printf("After expiration, cache size: %d\n", c.ItemCount())
	fmt.Println("✓ Janitor automatically cleaned expired items\n")
}

func demonstratePauseResume() {
	fmt.Println("=== Pause/Resume Janitor ===")

	c := cache.New(100*time.Millisecond, 25*time.Millisecond)

	// Add items
	c.Set("temp1", "value1", 100*time.Millisecond)
	c.Set("temp2", "value2", 100*time.Millisecond)
	fmt.Printf("Added 2 items, cache size: %d\n", c.ItemCount())

	// Pause janitor
	fmt.Println("Pausing janitor...")
	c.PauseJanitor()

	// Wait for items to expire
	fmt.Println("Waiting for items to expire (janitor paused)...")
	time.Sleep(150 * time.Millisecond)

	// Items are expired but not deleted
	fmt.Printf("After expiration (paused), cache still has %d items (expired but not cleaned)\n", c.ItemCount())

	// Resume janitor
	fmt.Println("Resuming janitor...")
	c.ResumeJanitor()

	// Wait for janitor to run
	time.Sleep(50 * time.Millisecond)

	// Now items should be cleaned
	fmt.Printf("After resume, cache size: %d\n", c.ItemCount())
	fmt.Println("✓ Janitor resumed and cleaned expired items\n")
}

func demonstrateIntervalChange() {
	fmt.Println("=== Dynamic Interval Change ===")

	// Start with slow cleanup (500ms interval)
	c := cache.New(50*time.Millisecond, 500*time.Millisecond)

	fmt.Println("Created cache with 500ms cleanup interval")

	// Add item with 50ms expiration
	c.Set("fast_expire", "value", 50*time.Millisecond)
	fmt.Println("Added item with 50ms expiration")

	// Change to fast cleanup (10ms interval)
	fmt.Println("Changing janitor interval to 10ms...")
	c.SetJanitorInterval(10 * time.Millisecond)

	// Wait for expiration and cleanup
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("Cache size after fast cleanup: %d\n", c.ItemCount())
	fmt.Println("✓ Interval change allowed faster cleanup\n")
}

func demonstrateCombinedControl() {
	fmt.Println("=== Combined Janitor Control ===")

	c := cache.New(100*time.Millisecond, 50*time.Millisecond)

	// Phase 1: Normal operation
	fmt.Println("Phase 1: Normal operation")
	c.Set("data1", "value1", 100*time.Millisecond)
	fmt.Printf("  Cache size: %d\n", c.ItemCount())

	// Phase 2: Pause for maintenance
	fmt.Println("\nPhase 2: Pausing for maintenance simulation")
	c.PauseJanitor()
	c.Set("data2", "value2", 100*time.Millisecond)
	c.Set("data3", "value3", 100*time.Millisecond)
	fmt.Printf("  Added more items during pause, size: %d\n", c.ItemCount())

	// Wait for some items to expire
	time.Sleep(120 * time.Millisecond)
	fmt.Printf("  After expiration (paused), size: %d (expired but not cleaned)\n", c.ItemCount())

	// Phase 3: Speed up cleanup and resume
	fmt.Println("\nPhase 3: Speeding up cleanup and resuming")
	c.SetJanitorInterval(10 * time.Millisecond)
	c.ResumeJanitor()

	time.Sleep(50 * time.Millisecond)
	fmt.Printf("  After resume with fast cleanup, size: %d\n", c.ItemCount())
	fmt.Println("✓ Successfully controlled janitor through multiple phases\n")

	// Demonstrate manual cleanup
	fmt.Println("=== Manual Cleanup (DeleteExpired) ===")
	c2 := cache.New(50*time.Millisecond, 0) // No automatic janitor

	// Add items
	for i := 0; i < 5; i++ {
		c2.Set(fmt.Sprintf("item%d", i), i, 50*time.Millisecond)
	}
	fmt.Printf("Added 5 items (no automatic janitor)\n")

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("After expiration, size: %d (no cleanup)\n", c2.ItemCount())

	// Manual cleanup
	fmt.Println("Calling DeleteExpired() manually...")
	c2.DeleteExpired()
	fmt.Printf("After manual cleanup, size: %d\n", c2.ItemCount())
	fmt.Println("✓ Manual cleanup works without janitor")

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\nKey Features Demonstrated:")
	fmt.Println("- PauseJanitor() - Pause automatic cleanup")
	fmt.Println("- ResumeJanitor() - Resume automatic cleanup")
	fmt.Println("- SetJanitorInterval() - Change cleanup frequency")
	fmt.Println("- DeleteExpired() - Manual cleanup")
}
