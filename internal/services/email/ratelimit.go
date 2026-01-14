package email

import (
	"context"
	"sync"
	"time"
)

// SlidingWindowLimiter implements a sliding window rate limiter
// More accurate than token bucket for SES rate limits
type SlidingWindowLimiter struct {
	mu           sync.Mutex
	windowSize   time.Duration
	maxRequests  int
	timestamps   []time.Time
	enabled      bool
}

// NewSlidingWindowLimiter creates a new sliding window rate limiter
// windowSize: duration of the sliding window (e.g., 1 second for SES)
// maxRequests: max requests allowed in the window
func NewSlidingWindowLimiter(windowSize time.Duration, maxRequests int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		windowSize:  windowSize,
		maxRequests: maxRequests,
		timestamps:  make([]time.Time, 0, maxRequests),
		enabled:     true,
	}
}

// Wait blocks until a request is allowed
// Returns error if context is cancelled
func (l *SlidingWindowLimiter) Wait(ctx context.Context) error {
	if !l.enabled {
		return nil
	}

	for {
		waitTime := l.timeUntilAllowed()
		if waitTime <= 0 {
			l.record()
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			// Retry after waiting
		}
	}
}

// Allow checks if a request is allowed without blocking
func (l *SlidingWindowLimiter) Allow() bool {
	if !l.enabled {
		return true
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.cleanup()
	if len(l.timestamps) < l.maxRequests {
		l.timestamps = append(l.timestamps, time.Now())
		return true
	}
	return false
}

// timeUntilAllowed returns how long to wait before next request is allowed
func (l *SlidingWindowLimiter) timeUntilAllowed() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cleanup()
	if len(l.timestamps) < l.maxRequests {
		return 0
	}

	// Wait until oldest timestamp expires
	oldest := l.timestamps[0]
	return time.Until(oldest.Add(l.windowSize))
}

// record adds current timestamp to the window
func (l *SlidingWindowLimiter) record() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timestamps = append(l.timestamps, time.Now())
}

// cleanup removes expired timestamps
func (l *SlidingWindowLimiter) cleanup() {
	cutoff := time.Now().Add(-l.windowSize)
	newIdx := 0
	for i, ts := range l.timestamps {
		if ts.After(cutoff) {
			newIdx = i
			break
		}
		newIdx = i + 1
	}
	if newIdx > 0 {
		l.timestamps = l.timestamps[newIdx:]
	}
}

// SetRate updates the rate limit
func (l *SlidingWindowLimiter) SetRate(maxRequests int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.maxRequests = maxRequests
}

// SetEnabled enables or disables the limiter
func (l *SlidingWindowLimiter) SetEnabled(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = enabled
}

// Stats returns current limiter statistics
func (l *SlidingWindowLimiter) Stats() (currentCount, maxCount int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cleanup()
	return len(l.timestamps), l.maxRequests
}

// CampaignRateLimiter manages per-campaign rate limits
type CampaignRateLimiter struct {
	mu        sync.RWMutex
	limiters  map[string]*SlidingWindowLimiter
	global    *SlidingWindowLimiter
	perSecond int
}

// NewCampaignRateLimiter creates a new campaign rate limiter
// globalRate: overall rate limit (e.g., 14/sec for SES sandbox)
// perCampaignRate: rate limit per campaign (0 = no per-campaign limit)
func NewCampaignRateLimiter(globalRate int) *CampaignRateLimiter {
	return &CampaignRateLimiter{
		limiters:  make(map[string]*SlidingWindowLimiter),
		global:    NewSlidingWindowLimiter(time.Second, globalRate),
		perSecond: globalRate,
	}
}

// Wait waits for rate limit for a specific campaign
func (c *CampaignRateLimiter) Wait(ctx context.Context, campaignID string) error {
	// Global rate limit
	if err := c.global.Wait(ctx); err != nil {
		return err
	}
	return nil
}

// WaitGlobal waits for the global rate limit only
func (c *CampaignRateLimiter) WaitGlobal(ctx context.Context) error {
	return c.global.Wait(ctx)
}

// SetGlobalRate updates the global rate limit
func (c *CampaignRateLimiter) SetGlobalRate(rate int) {
	c.global.SetRate(rate)
	c.perSecond = rate
}

// GlobalStats returns global limiter statistics
func (c *CampaignRateLimiter) GlobalStats() (current, max int) {
	return c.global.Stats()
}
