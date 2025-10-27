package client

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter *rate.Limiter // golang.org/x/time/rate
	mu      sync.RWMutex

	// Track daily usage (resets at account's midnight)
	dailyLimit     int
	dailyRemaining int
	dailyResetTime time.Time
}

// RateLimiter manages rate limiting for API requests
func NewRateLimiter(maxBurst int) *RateLimiter {
	// Convert requests per 10 seconds to per-second rate
	requestsPerSecond := float64(maxBurst) / 10.0

	return &RateLimiter{
		limiter:        rate.NewLimiter(rate.Limit(requestsPerSecond), maxBurst),
		dailyLimit:     250000,
		dailyRemaining: 250000,
		dailyResetTime: time.Now().AddDate(0, 0, 1),
	}
}

// Wait blocks until a token is available for use
func (rl *RateLimiter) Wait(ctx context.Context) error {
	return rl.limiter.Wait(ctx)
}

// UpdateFromResponse updates the rate limiter state from response headers
func (rl *RateLimiter) UpdateFromResponse(resp *Response) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.dailyRemaining = resp.RateLimit.DailyRemaining
	rl.dailyLimit = resp.RateLimit.DailyLimit
}

// CheckDailyLimit returns true if daily quota is available
func (rl *RateLimiter) CheckDailyLimit() bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return rl.dailyRemaining > 0
}

// GetDailyRemaining return the current daily remaining quota
func (rl *RateLimiter) GetDailyRemaining() int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return rl.dailyRemaining
}
