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

// ========== Constructor Tests ==========

func TestNewAuthRateLimitMiddleware(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	assert.NotNil(t, middleware, "NewAuthRateLimitMiddleware should return a non-nil middleware")
	assert.Equal(t, 5, middleware.maxAttempts, "Default max attempts should be 5")
	assert.Equal(t, 60, middleware.windowSeconds, "Default window should be 60 seconds")
	assert.Equal(t, 300, middleware.blockDurationSeconds, "Default block duration should be 300 seconds")
	assert.NotNil(t, middleware.entries, "Entries map should be initialized")
}

// ========== Rate Limit Enforcement Tests ==========

func TestAuthRateLimitMiddleware_FirstRequest(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	handlerCalled := false
	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.True(t, handlerCalled, "Handler should be called for first request")
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthRateLimitMiddleware_WithinRateLimit(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make 5 requests (the default limit)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.RemoteAddr = "192.168.1.2:12345"
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Request %d should succeed", i+1)
	}
}

func TestAuthRateLimitMiddleware_ExceedsRateLimit(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make 5 requests (the limit)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.RemoteAddr = "192.168.1.3:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}

	// 6th request should be blocked
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
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

func TestAuthRateLimitMiddleware_BlockedRequestsStayBlocked(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Exhaust the limit (5 requests) + 1 to trigger block
	for i := 0; i < 6; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.RemoteAddr = "192.168.1.4:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
	}

	// Multiple subsequent requests should remain blocked
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.RemoteAddr = "192.168.1.4:12345"
		rr := httptest.NewRecorder()

		handler(rr, req)

		assert.Equal(t, http.StatusTooManyRequests, rr.Code, "Request should remain blocked")
	}
}

// ========== Window Expiration Tests ==========

func TestAuthRateLimitMiddleware_WindowExpiration(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	// Manually set short window for testing
	middleware.windowSeconds = 1

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make 5 requests (reach the limit but don't exceed)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.RemoteAddr = "192.168.1.5:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}

	// Wait for window to expire
	time.Sleep(1100 * time.Millisecond)

	// After window expires, should be able to make requests again
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.5:12345"
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Request should succeed after window expires")
}

func TestAuthRateLimitMiddleware_BlockExpiration(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	// Manually set short block duration for testing
	middleware.blockDurationSeconds = 1
	middleware.maxAttempts = 1

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// First request succeeds
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.6:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Second request triggers block
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.6:12345"
	rr = httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// Wait for block to expire
	time.Sleep(1100 * time.Millisecond)

	// After block expires, should be able to make requests again
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.6:12345"
	rr = httptest.NewRecorder()

	handler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Request should succeed after block expires")
}

// ========== Different Endpoints Tests ==========

func TestAuthRateLimitMiddleware_DifferentIPsIndependent(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 2

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Exhaust limit for IP1
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
	}

	// IP1 should be blocked
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// IP2 should still work
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.200:12345"
	rr = httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Different IP should not be affected")
}

func TestAuthRateLimitMiddleware_SharedAcrossEndpoints(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 3

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make requests to different endpoints but same IP
	endpoints := []string{"/api/auth/login", "/api/auth/register", "/api/auth/forgot-password"}

	for i, endpoint := range endpoints {
		req := httptest.NewRequest(http.MethodPost, endpoint, nil)
		req.RemoteAddr = "192.168.1.101:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "Request %d to %s should succeed", i+1, endpoint)
	}

	// 4th request (to any endpoint) should be blocked
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.101:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code, "Rate limit should apply across endpoints")
}

// ========== Client IP Extraction Tests ==========

func TestAuthRateLimitMiddleware_GetClientIP_FromXForwardedFor(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

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

			ip := middleware.getClientIP(req)
			assert.Equal(t, tc.expectedIP, ip)
		})
	}
}

func TestAuthRateLimitMiddleware_GetClientIP_FromXRealIP(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "10.0.0.50")

	ip := middleware.getClientIP(req)
	assert.Equal(t, "10.0.0.50", ip)
}

func TestAuthRateLimitMiddleware_GetClientIP_FallbackToRemoteAddr(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.0.1:8080"

	ip := middleware.getClientIP(req)
	assert.Equal(t, "192.168.0.1:8080", ip)
}

func TestAuthRateLimitMiddleware_GetClientIP_XForwardedForTakesPrecedence(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.50")
	req.Header.Set("X-Real-IP", "10.0.0.1")
	req.RemoteAddr = "127.0.0.1:8080"

	ip := middleware.getClientIP(req)
	assert.Equal(t, "203.0.113.50", ip, "X-Forwarded-For should take precedence")
}

func TestAuthRateLimitMiddleware_GetClientIP_XRealIPTakesPrecedenceOverRemoteAddr(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "10.0.0.1")
	req.RemoteAddr = "127.0.0.1:8080"

	ip := middleware.getClientIP(req)
	assert.Equal(t, "10.0.0.1", ip, "X-Real-IP should take precedence over RemoteAddr")
}

// ========== Concurrent Request Handling Tests ==========

func TestAuthRateLimitMiddleware_ConcurrentRequests(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 50

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

			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
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

func TestAuthRateLimitMiddleware_ConcurrentDifferentIPs(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var wg sync.WaitGroup
	numIPs := 10
	requestsPerIP := 5 // Within limit

	var successCount int
	var mu sync.Mutex

	for ip := 0; ip < numIPs; ip++ {
		for req := 0; req < requestsPerIP; req++ {
			wg.Add(1)
			ipAddr := ip // Capture for goroutine
			go func() {
				defer wg.Done()

				req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
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

func TestAuthRateLimitMiddleware_ResponseHeaders_WhenBlocked(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 1

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// First request
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.7:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)

	// Second request triggers block
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.7:12345"
	rr = httptest.NewRecorder()

	handler(rr, req)

	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.NotEmpty(t, rr.Header().Get("Retry-After"), "Retry-After header should be set")
}

func TestAuthRateLimitMiddleware_RetryAfterHeader_Format(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 1
	middleware.blockDurationSeconds = 300

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Exhaust limit
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.20:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)

	// Trigger block
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.20:12345"
	rr = httptest.NewRecorder()

	handler(rr, req)

	retryAfter := rr.Header().Get("Retry-After")
	assert.NotEmpty(t, retryAfter, "Retry-After header should be set")
	// The header should be the block duration as a string
	assert.Equal(t, "300", retryAfter)
}

func TestAuthRateLimitMiddleware_ResponseBody_WhenBlocked(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 1
	middleware.blockDurationSeconds = 300

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Exhaust limit
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.8:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)

	// Trigger block
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
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

func TestAuthRateLimitMiddleware_ResponseBody_RemainingSeconds(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 1
	middleware.blockDurationSeconds = 60

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Exhaust limit and trigger block
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.21:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)

	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.21:12345"
	rr = httptest.NewRecorder()
	handler(rr, req)

	// Wait a second and make another request
	time.Sleep(1 * time.Second)

	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.21:12345"
	rr = httptest.NewRecorder()
	handler(rr, req)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	// retry_after should be less than block duration (since we waited 1 second)
	retryAfter := response["retry_after"].(float64)
	assert.Less(t, retryAfter, float64(60), "retry_after should decrease over time")
	assert.Greater(t, retryAfter, float64(55), "retry_after should not decrease too much")
}

// ========== RecordFailure Tests ==========

func TestAuthRateLimitMiddleware_RecordFailure_NewIP(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
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

func TestAuthRateLimitMiddleware_RecordFailure_ExistingIP(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
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

func TestAuthRateLimitMiddleware_RecordFailure_TriggersBlock(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.11:12345"

	// Record failures to exceed limit (5 default + 1 to trigger)
	for i := 0; i < 6; i++ {
		middleware.RecordFailure(req)
	}

	middleware.mu.RLock()
	entry, _ := middleware.entries["192.168.1.11:12345"]
	middleware.mu.RUnlock()

	assert.False(t, entry.blockedAt.IsZero(), "IP should be blocked after exceeding limit")
}

func TestAuthRateLimitMiddleware_RecordFailure_DoesNotIncrementWhenBlocked(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.12:12345"

	// Exceed limit and get blocked
	for i := 0; i < 10; i++ {
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

func TestAuthRateLimitMiddleware_RecordFailure_CombinedWithRequests(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 3

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make 2 requests
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.RemoteAddr = "192.168.1.22:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}

	// Record a failure
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.22:12345"
	middleware.RecordFailure(req)

	// Next request should trigger block
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.22:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
}

// ========== ResetForIP Tests ==========

func TestAuthRateLimitMiddleware_ResetForIP(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make some requests
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
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
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.13:12345"
	middleware.ResetForIP(req)

	// Verify entry is gone
	middleware.mu.RLock()
	_, exists = middleware.entries["192.168.1.13:12345"]
	middleware.mu.RUnlock()
	assert.False(t, exists, "Entry should be removed after ResetForIP")
}

func TestAuthRateLimitMiddleware_ResetForIP_AllowsNewRequests(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 2

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Exhaust limit and get blocked
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.RemoteAddr = "192.168.1.14:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
	}

	// Verify blocked
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.14:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// Reset
	middleware.ResetForIP(req)

	// Should be able to make requests again
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.14:12345"
	rr = httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Request should succeed after reset")
}

func TestAuthRateLimitMiddleware_ResetForIP_NonExistentIP(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.15:12345"

	// Should not panic when resetting non-existent IP
	assert.NotPanics(t, func() {
		middleware.ResetForIP(req)
	})
}

func TestAuthRateLimitMiddleware_ResetForIP_UseCaseSuccessfulLogin(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 3

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Simulate 2 failed login attempts
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.23:12345"
	for i := 0; i < 2; i++ {
		rr := httptest.NewRecorder()
		handler(rr, req)
		middleware.RecordFailure(req) // Failed login
	}

	// Simulate successful login - reset rate limit
	middleware.ResetForIP(req)

	// User should now have full quota again
	for i := 0; i < 3; i++ {
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "Request %d should succeed after successful login", i+1)
	}
}

// ========== Edge Cases ==========

func TestAuthRateLimitMiddleware_EmptyRemoteAddr(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

	handlerCalled := false
	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = ""
	rr := httptest.NewRecorder()

	handler(rr, req)

	assert.True(t, handlerCalled, "Handler should still be called with empty RemoteAddr")
}

func TestAuthRateLimitMiddleware_DifferentHTTPMethods(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()

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

			req := httptest.NewRequest(method, "/api/auth/endpoint", nil)
			req.RemoteAddr = ipAddr
			rr := httptest.NewRecorder()

			handler(rr, req)

			assert.True(t, handlerCalled, "Handler should be called for %s method", method)
			assert.Equal(t, http.StatusOK, rr.Code)
		})
	}
}

// ========== Integration Scenario Tests ==========

func TestAuthRateLimitMiddleware_BruteForceScenario(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 5
	middleware.blockDurationSeconds = 60

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	attackerIP := "192.168.1.100:12345"

	// Attacker makes 5 rapid login attempts
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.RemoteAddr = attackerIP
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}

	// 6th attempt should be blocked
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = attackerIP
	rr := httptest.NewRecorder()
	handler(rr, req)

	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, "Too many attempts. Please try again later.", response["error"])
}

func TestAuthRateLimitMiddleware_LegitimateUserAfterBlock(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 3
	middleware.blockDurationSeconds = 1 // Short for testing

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	userIP := "192.168.1.101:12345"

	// User forgets password, makes multiple attempts
	for i := 0; i < 4; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.RemoteAddr = userIP
		rr := httptest.NewRecorder()
		handler(rr, req)
	}

	// User is now blocked
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = userIP
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// User waits for block to expire
	time.Sleep(1100 * time.Millisecond)

	// User can try again with correct password
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = userIP
	rr = httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// ========== Stress Test ==========

func TestAuthRateLimitMiddleware_HighVolumeRequests(t *testing.T) {
	middleware := NewAuthRateLimitMiddleware()
	middleware.maxAttempts = 1000

	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make 1000 requests from single IP
	for i := 0; i < 1000; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.RemoteAddr = "192.168.1.17:12345"
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "Request %d should succeed", i+1)
	}

	// 1001st should be blocked
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "192.168.1.17:12345"
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
}
