package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========== RateLimitConfig Tests ==========

func TestDefaultAuthRateLimitConfig(t *testing.T) {
	config := DefaultAuthRateLimitConfig()

	assert.Equal(t, 5, config.MaxAttempts, "Default max attempts should be 5")
	assert.Equal(t, 60, config.WindowSeconds, "Default window should be 60 seconds")
	assert.Equal(t, 300, config.BlockDurationSeconds, "Default block duration should be 300 seconds")
}

func TestNewRateLimitMiddleware(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          10,
		WindowSeconds:        30,
		BlockDurationSeconds: 60,
	}

	middleware := NewRateLimitMiddleware(config)

	assert.NotNil(t, middleware, "NewRateLimitMiddleware should return a non-nil middleware")
	assert.Equal(t, 10, middleware.config.MaxAttempts)
	assert.Equal(t, 30, middleware.config.WindowSeconds)
	assert.Equal(t, 60, middleware.config.BlockDurationSeconds)
	assert.NotNil(t, middleware.entries, "Entries map should be initialized")
}

// ========== Rate Limit Enforcement Tests ==========

func TestRateLimitMiddleware_FirstRequest(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          5,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handlerCalled := false
	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.True(t, handlerCalled, "Handler should be called for first request")
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRateLimitMiddleware_WithinRateLimit(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          5,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make 5 requests (the limit)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		req.RemoteAddr = "192.168.1.2:12345"
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Request %d should succeed", i+1)
	}
}

func TestRateLimitMiddleware_ExceedsRateLimit(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          3,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make 3 requests (the limit)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		req.RemoteAddr = "192.168.1.3:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}

	// 4th request should be blocked
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.3:12345"
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "Too many attempts. Please try again later.", response["error"])
	assert.NotNil(t, response["retry_after"])
	assert.NotNil(t, response["blocked_until"])
}

func TestRateLimitMiddleware_BlockedRequestsStayBlocked(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          2,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Exhaust the limit
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		req.RemoteAddr = "192.168.1.4:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
	}

	// Multiple subsequent requests should remain blocked
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		req.RemoteAddr = "192.168.1.4:12345"
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, http.StatusTooManyRequests, rr.Code, "Request should remain blocked")
	}
}

// ========== Window Expiration Tests ==========

func TestRateLimitMiddleware_WindowExpiration(t *testing.T) {
	// Use a very short window for testing
	config := RateLimitConfig{
		MaxAttempts:          2,
		WindowSeconds:        1, // 1 second window
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make 2 requests (reach the limit but don't exceed)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		req.RemoteAddr = "192.168.1.5:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}

	// Wait for window to expire
	time.Sleep(1100 * time.Millisecond)

	// After window expires, should be able to make requests again
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.5:12345"
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Request should succeed after window expires")
}

func TestRateLimitMiddleware_BlockExpiration(t *testing.T) {
	// Use a very short block duration for testing
	config := RateLimitConfig{
		MaxAttempts:          1,
		WindowSeconds:        60,
		BlockDurationSeconds: 1, // 1 second block
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// First request succeeds
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.6:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Second request triggers block
	req = httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.6:12345"
	rr = httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// Wait for block to expire
	time.Sleep(1100 * time.Millisecond)

	// After block expires, should be able to make requests again
	req = httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.6:12345"
	rr = httptest.NewRecorder()

	handler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Request should succeed after block expires")
}

// ========== Different Rate Limits Tests ==========

func TestRateLimitMiddleware_DifferentConfigurations(t *testing.T) {
	testCases := []struct {
		name              string
		maxAttempts       int
		windowSeconds     int
		blockDuration     int
		requestsToSucceed int
	}{
		{"Strict limit", 1, 60, 300, 1},
		{"Standard limit", 5, 60, 300, 5},
		{"Lenient limit", 10, 60, 300, 10},
		{"Very short window", 3, 1, 300, 3},
		{"High volume", 100, 60, 300, 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := RateLimitConfig{
				MaxAttempts:          tc.maxAttempts,
				WindowSeconds:        tc.windowSeconds,
				BlockDurationSeconds: tc.blockDuration,
			}
			middleware := NewRateLimitMiddleware(config)

			handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Should succeed up to maxAttempts
			for i := 0; i < tc.requestsToSucceed; i++ {
				req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
				req.RemoteAddr = "10.0.0.1:12345"
				rr := httptest.NewRecorder()
				handler(rr, req)
				assert.Equal(t, http.StatusOK, rr.Code, "Request %d should succeed", i+1)
			}

			// Next request should be blocked
			req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
			req.RemoteAddr = "10.0.0.1:12345"
			rr := httptest.NewRecorder()
			handler(rr, req)
			assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		})
	}
}

func TestRateLimitMiddleware_DifferentIPsIndependent(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          2,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Exhaust limit for IP1
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
	}

	// IP1 should be blocked
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// IP2 should still work
	req = httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.200:12345"
	rr = httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Different IP should not be affected")
}

// ========== Client IP Extraction Tests ==========

func TestGetClientIP_FromXForwardedFor(t *testing.T) {
	testCases := []struct {
		name       string
		xff        string
		expectedIP string
	}{
		{"Single IP", "203.0.113.50", "203.0.113.50"},
		{"Multiple IPs", "203.0.113.50, 70.41.3.18, 150.172.238.178", "203.0.113.50"},
		{"IPv6", "2001:db8::1", "2001:db8::1"},
		{"Mixed chain", "192.0.2.1, 2001:db8::1", "192.0.2.1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("X-Forwarded-For", tc.xff)

			ip := getClientIP(req)
			assert.Equal(t, tc.expectedIP, ip)
		})
	}
}

func TestGetClientIP_FromXRealIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "10.0.0.50")

	ip := getClientIP(req)
	assert.Equal(t, "10.0.0.50", ip)
}

func TestGetClientIP_FallbackToRemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.0.1:8080"

	ip := getClientIP(req)
	assert.Equal(t, "192.168.0.1:8080", ip)
}

func TestGetClientIP_XForwardedForTakesPrecedence(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.50")
	req.Header.Set("X-Real-IP", "10.0.0.1")
	req.RemoteAddr = "127.0.0.1:8080"

	ip := getClientIP(req)
	assert.Equal(t, "203.0.113.50", ip, "X-Forwarded-For should take precedence")
}

func TestGetClientIP_XRealIPTakesPrecedenceOverRemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "10.0.0.1")
	req.RemoteAddr = "127.0.0.1:8080"

	ip := getClientIP(req)
	assert.Equal(t, "10.0.0.1", ip, "X-Real-IP should take precedence over RemoteAddr")
}

// ========== Concurrent Request Handling Tests ==========

func TestRateLimitMiddleware_ConcurrentRequests(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          50,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	var successCount int
	var blockedCount int
	var mu sync.Mutex

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
			req.RemoteAddr = "192.168.1.50:12345"
			rr := httptest.NewRecorder()

			handler(rr, req)

			mu.Lock()
			if rr.Code == http.StatusOK {
				successCount++
			} else if rr.Code == http.StatusTooManyRequests {
				blockedCount++
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	// Should have exactly 50 successes and 50 blocked
	assert.Equal(t, 50, successCount, "Should have exactly maxAttempts successes")
	assert.Equal(t, 50, blockedCount, "Remaining requests should be blocked")
}

func TestRateLimitMiddleware_ConcurrentDifferentIPs(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          5,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var wg sync.WaitGroup
	numIPs := 10
	requestsPerIP := 5

	var successCount int
	var mu sync.Mutex

	for ip := 0; ip < numIPs; ip++ {
		for req := 0; req < requestsPerIP; req++ {
			wg.Add(1)
			ipAddr := ip // Capture for goroutine
			go func() {
				defer wg.Done()

				req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
				req.RemoteAddr = "192.168.1." + string(rune('0'+ipAddr)) + ":12345"
				rr := httptest.NewRecorder()

				handler(rr, req)

				mu.Lock()
				if rr.Code == http.StatusOK {
					successCount++
				}
				mu.Unlock()
			}()
		}
	}

	wg.Wait()

	// All requests should succeed since each IP stays within limit
	assert.Equal(t, numIPs*requestsPerIP, successCount, "All requests should succeed within limits")
}

// ========== Response Headers Tests ==========

func TestRateLimitMiddleware_ResponseHeaders_WhenBlocked(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          1,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// First request
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.7:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)

	// Second request triggers block
	req = httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.7:12345"
	rr = httptest.NewRecorder()

	handler(rr, req)

	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.NotEmpty(t, rr.Header().Get("Retry-After"), "Retry-After header should be set")
}

func TestRateLimitMiddleware_ResponseBody_WhenBlocked(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          1,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Exhaust limit
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.8:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)

	// Trigger block
	req = httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.8:12345"
	rr = httptest.NewRecorder()

	handler(rr, req)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "Too many attempts. Please try again later.", response["error"])
	assert.Equal(t, float64(300), response["retry_after"])

	// blocked_until should be a valid RFC3339 timestamp
	blockedUntil, ok := response["blocked_until"].(string)
	require.True(t, ok, "blocked_until should be a string")
	_, err = time.Parse(time.RFC3339, blockedUntil)
	assert.NoError(t, err, "blocked_until should be valid RFC3339 format")
}

// ========== RecordFailure Tests ==========

func TestRateLimitMiddleware_RecordFailure_NewIP(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          5,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.9:12345"

	// Record a failure for a new IP
	middleware.RecordFailure(req)

	// Entry should exist
	middleware.mu.RLock()
	entry, exists := middleware.entries["192.168.1.9:12345"]
	middleware.mu.RUnlock()

	assert.True(t, exists, "Entry should exist after RecordFailure")
	assert.Equal(t, 1, entry.attempts)
}

func TestRateLimitMiddleware_RecordFailure_ExistingIP(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          5,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.10:12345"

	// Record multiple failures
	middleware.RecordFailure(req)
	middleware.RecordFailure(req)
	middleware.RecordFailure(req)

	middleware.mu.RLock()
	entry, _ := middleware.entries["192.168.1.10:12345"]
	middleware.mu.RUnlock()

	assert.Equal(t, 3, entry.attempts)
}

func TestRateLimitMiddleware_RecordFailure_TriggersBlock(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          2,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.11:12345"

	// Record failures to exceed limit
	middleware.RecordFailure(req)
	middleware.RecordFailure(req)
	middleware.RecordFailure(req) // This should trigger block

	middleware.mu.RLock()
	entry, _ := middleware.entries["192.168.1.11:12345"]
	middleware.mu.RUnlock()

	assert.False(t, entry.blockedAt.IsZero(), "IP should be blocked after exceeding limit")
}

func TestRateLimitMiddleware_RecordFailure_DoesNotIncrementWhenBlocked(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          2,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.12:12345"

	// Exceed limit and get blocked
	for i := 0; i < 5; i++ {
		middleware.RecordFailure(req)
	}

	middleware.mu.RLock()
	entry, _ := middleware.entries["192.168.1.12:12345"]
	attemptsAfterBlock := entry.attempts
	middleware.mu.RUnlock()

	// Record more failures - should not increment
	middleware.RecordFailure(req)
	middleware.RecordFailure(req)

	middleware.mu.RLock()
	entry, _ = middleware.entries["192.168.1.12:12345"]
	middleware.mu.RUnlock()

	assert.Equal(t, attemptsAfterBlock, entry.attempts, "Attempts should not increment when blocked")
}

// ========== ResetForIP Tests ==========

func TestRateLimitMiddleware_ResetForIP(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          5,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make some requests
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		req.RemoteAddr = "192.168.1.13:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
	}

	// Verify entry exists
	middleware.mu.RLock()
	_, exists := middleware.entries["192.168.1.13:12345"]
	middleware.mu.RUnlock()
	assert.True(t, exists)

	// Reset for this IP
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.13:12345"
	middleware.ResetForIP(req)

	// Verify entry is gone
	middleware.mu.RLock()
	_, exists = middleware.entries["192.168.1.13:12345"]
	middleware.mu.RUnlock()
	assert.False(t, exists, "Entry should be removed after ResetForIP")
}

func TestRateLimitMiddleware_ResetForIP_AllowsNewRequests(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          2,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Exhaust limit and get blocked
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		req.RemoteAddr = "192.168.1.14:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
	}

	// Verify blocked
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.14:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// Reset
	middleware.ResetForIP(req)

	// Should be able to make requests again
	req = httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.14:12345"
	rr = httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Request should succeed after reset")
}

func TestRateLimitMiddleware_ResetForIP_NonExistentIP(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          5,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.15:12345"

	// Should not panic when resetting non-existent IP
	assert.NotPanics(t, func() {
		middleware.ResetForIP(req)
	})
}

// ========== Edge Cases ==========

func TestRateLimitMiddleware_EmptyRemoteAddr(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          5,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handlerCalled := false
	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = ""
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.True(t, handlerCalled, "Handler should still be called with empty RemoteAddr")
}

func TestRateLimitMiddleware_ZeroConfig(t *testing.T) {
	// Test with zero values - should effectively block everything after first request
	config := RateLimitConfig{
		MaxAttempts:          0,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// First request goes through (new IP)
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.16:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Second request should be blocked (0 max attempts means blocked immediately after first)
	req = httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.16:12345"
	rr = httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
}

func TestRateLimitMiddleware_DifferentHTTPMethods(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          5,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			// Use different IP for each method test
			ipAddr := "10.0.0." + string(rune('1'+len(method))) + ":12345"

			handlerCalled := false
			handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(method, "/api/endpoint", nil)
			req.RemoteAddr = ipAddr
			rr := httptest.NewRecorder()

			handler(rr, req)

			assert.True(t, handlerCalled, "Handler should be called for %s method", method)
			assert.Equal(t, http.StatusOK, rr.Code)
		})
	}
}

// ========== Stress Test ==========

func TestRateLimitMiddleware_HighVolumeRequests(t *testing.T) {
	config := RateLimitConfig{
		MaxAttempts:          1000,
		WindowSeconds:        60,
		BlockDurationSeconds: 300,
	}
	middleware := NewRateLimitMiddleware(config)

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make 1000 requests from single IP
	for i := 0; i < 1000; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		req.RemoteAddr = "192.168.1.17:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "Request %d should succeed", i+1)
	}

	// 1001st should be blocked
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.17:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
}
