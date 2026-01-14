package email

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSlidingWindowLimiter_Allow_UnderLimit(t *testing.T) {
	limiter := NewSlidingWindowLimiter(time.Second, 5)

	// Should allow 5 requests
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}
}

func TestSlidingWindowLimiter_Allow_AtLimit(t *testing.T) {
	limiter := NewSlidingWindowLimiter(time.Second, 3)

	// Use up the limit
	for i := 0; i < 3; i++ {
		limiter.Allow()
	}

	// Next request should be blocked
	if limiter.Allow() {
		t.Error("Request should be blocked when at limit")
	}
}

func TestSlidingWindowLimiter_Allow_WindowSlides(t *testing.T) {
	limiter := NewSlidingWindowLimiter(50*time.Millisecond, 2)

	// Use up the limit
	limiter.Allow()
	limiter.Allow()

	// Should be blocked
	if limiter.Allow() {
		t.Error("Request should be blocked")
	}

	// Wait for window to slide
	time.Sleep(60 * time.Millisecond)

	// Should be allowed again
	if !limiter.Allow() {
		t.Error("Request should be allowed after window slides")
	}
}

func TestSlidingWindowLimiter_Wait_ImmediateWhenUnderLimit(t *testing.T) {
	limiter := NewSlidingWindowLimiter(time.Second, 5)
	ctx := context.Background()

	start := time.Now()
	err := limiter.Wait(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Wait should not error: %v", err)
	}
	if elapsed > 10*time.Millisecond {
		t.Errorf("Wait should be immediate, took %v", elapsed)
	}
}

func TestSlidingWindowLimiter_Wait_BlocksWhenAtLimit(t *testing.T) {
	limiter := NewSlidingWindowLimiter(100*time.Millisecond, 2)
	ctx := context.Background()

	// Use up the limit
	limiter.Wait(ctx)
	limiter.Wait(ctx)

	// Next wait should block
	start := time.Now()
	err := limiter.Wait(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Wait should not error: %v", err)
	}
	if elapsed < 50*time.Millisecond {
		t.Errorf("Wait should block, only took %v", elapsed)
	}
}

func TestSlidingWindowLimiter_Wait_RespectsContextCancellation(t *testing.T) {
	limiter := NewSlidingWindowLimiter(time.Second, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Use up the limit
	limiter.Wait(context.Background())

	// Next wait should be cancelled
	err := limiter.Wait(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("Wait should return context error, got: %v", err)
	}
}

func TestSlidingWindowLimiter_SetRate(t *testing.T) {
	limiter := NewSlidingWindowLimiter(time.Second, 2)

	// Use up initial limit
	limiter.Allow()
	limiter.Allow()
	if limiter.Allow() {
		t.Error("Should be blocked at limit 2")
	}

	// Increase the rate
	limiter.SetRate(5)

	// Should now be allowed
	if !limiter.Allow() {
		t.Error("Should be allowed after rate increase")
	}
}

func TestSlidingWindowLimiter_SetEnabled(t *testing.T) {
	limiter := NewSlidingWindowLimiter(time.Second, 1)

	// Use up limit
	limiter.Allow()
	if limiter.Allow() {
		t.Error("Should be blocked")
	}

	// Disable limiter
	limiter.SetEnabled(false)

	// Should now allow all requests
	if !limiter.Allow() {
		t.Error("Should be allowed when disabled")
	}
	if !limiter.Allow() {
		t.Error("Should be allowed when disabled")
	}
}

func TestSlidingWindowLimiter_Stats(t *testing.T) {
	limiter := NewSlidingWindowLimiter(time.Second, 5)

	current, max := limiter.Stats()
	if current != 0 || max != 5 {
		t.Errorf("Initial stats should be 0/5, got %d/%d", current, max)
	}

	limiter.Allow()
	limiter.Allow()

	current, max = limiter.Stats()
	if current != 2 || max != 5 {
		t.Errorf("Stats should be 2/5, got %d/%d", current, max)
	}
}

func TestSlidingWindowLimiter_Concurrent(t *testing.T) {
	limiter := NewSlidingWindowLimiter(100*time.Millisecond, 10)
	ctx := context.Background()

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	// Launch 20 goroutines simultaneously
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow() {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Should have allowed at most 10
	if successCount > 10 {
		t.Errorf("Should allow at most 10, allowed %d", successCount)
	}

	// Wait for window to reset
	time.Sleep(150 * time.Millisecond)

	// Now use Wait - should eventually allow all 10
	var wg2 sync.WaitGroup
	waitSuccessCount := 0
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			if err := limiter.Wait(ctx); err == nil {
				mu.Lock()
				waitSuccessCount++
				mu.Unlock()
			}
		}()
	}

	wg2.Wait()

	if waitSuccessCount != 10 {
		t.Errorf("Wait should succeed for all 10, got %d", waitSuccessCount)
	}
}

func TestCampaignRateLimiter_GlobalLimit(t *testing.T) {
	limiter := NewCampaignRateLimiter(5)
	ctx := context.Background()

	// Should allow 5 requests across any campaign
	for i := 0; i < 5; i++ {
		if err := limiter.Wait(ctx, "campaign-1"); err != nil {
			t.Errorf("Request %d should be allowed: %v", i+1, err)
		}
	}

	// Verify stats
	current, max := limiter.GlobalStats()
	if current != 5 || max != 5 {
		t.Errorf("Stats should be 5/5, got %d/%d", current, max)
	}
}

func TestCampaignRateLimiter_SetGlobalRate(t *testing.T) {
	limiter := NewCampaignRateLimiter(2)

	// Use up limit
	limiter.WaitGlobal(context.Background())
	limiter.WaitGlobal(context.Background())

	current, max := limiter.GlobalStats()
	if current != 2 {
		t.Errorf("Current should be 2, got %d", current)
	}

	// Increase rate
	limiter.SetGlobalRate(10)

	_, max = limiter.GlobalStats()
	if max != 10 {
		t.Errorf("Max should be 10, got %d", max)
	}
}

func TestCampaignRateLimiter_WaitGlobal_RespectsContext(t *testing.T) {
	limiter := NewCampaignRateLimiter(1)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Use up limit
	limiter.WaitGlobal(context.Background())

	// Should timeout
	err := limiter.WaitGlobal(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("Should timeout, got: %v", err)
	}
}
