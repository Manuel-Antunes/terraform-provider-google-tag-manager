package api

import (
	"sync"
	"time"
)

// TestCoordinator helps manage API rate limits by controlling the timing between API calls
type TestCoordinator struct {
	mutex           sync.Mutex
	lastRequestTime time.Time
	minDelay        time.Duration
}

// NewTestCoordinator creates a new coordinator with the specified minimum delay between operations
func NewTestCoordinator(minDelay time.Duration) *TestCoordinator {
	return &TestCoordinator{
		lastRequestTime: time.Now().Add(-minDelay), // Allow immediate first call
		minDelay:        minDelay,
	}
}

// WaitBeforeRequest waits if necessary to ensure minimum delay between API requests
func (c *TestCoordinator) WaitBeforeRequest() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	elapsed := time.Since(c.lastRequestTime)
	if elapsed < c.minDelay {
		sleepTime := c.minDelay - elapsed
		time.Sleep(sleepTime)
	}
	c.lastRequestTime = time.Now()
}

var (
	// Global test coordinator with 2 second delay between API calls
	// This helps prevent rate limit errors when running tests
	GlobalTestCoordinator = NewTestCoordinator(2 * time.Second)
)

// Helper function to generate a unique test name based on current time
func testName(prefix string) string {
	return prefix + "-" + time.Now().Format("20060102-150405")
}
