package util

import (
	"sync"
	"time"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	store map[string][]time.Time
	mu    sync.Mutex
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		store: make(map[string][]time.Time),
	}
}

// CheckLimit checks if the request is within rate limit
// Returns true if allowed, false if exceeded
func (rl *RateLimiter) CheckLimit(key string, limit int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	requests := rl.store[key]

	// Remove old requests outside the time window
	valid := []time.Time{}
	for _, t := range requests {
		if now.Sub(t) < window {
			valid = append(valid, t)
		}
	}

	// Check if limit exceeded
	if len(valid) >= limit {
		return false // Rate limit exceeded
	}

	// Add current request
	valid = append(valid, now)
	rl.store[key] = valid

	return true // Allowed
}

// Cleanup removes old entries from the store
// Should be called periodically to prevent memory leak
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	maxWindow := 24 * time.Hour // Keep last 24 hours

	for key, requests := range rl.store {
		valid := []time.Time{}
		for _, t := range requests {
			if now.Sub(t) < maxWindow {
				valid = append(valid, t)
			}
		}

		if len(valid) == 0 {
			delete(rl.store, key)
		} else {
			rl.store[key] = valid
		}
	}
}
