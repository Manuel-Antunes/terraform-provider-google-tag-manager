package provider

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
	// Global test coordinator with 3 second delay between API calls in provider tests
	// These tests tend to be more resource intensive, so we use a slightly longer delay
	GlobalTestCoordinator = NewTestCoordinator(3 * time.Second)
)
