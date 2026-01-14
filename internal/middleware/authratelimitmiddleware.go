package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// AuthRateLimitMiddleware provides rate limiting for authentication endpoints
// to prevent brute force attacks. Uses in-memory storage for tracking attempts.
type AuthRateLimitMiddleware struct {
	maxAttempts          int
	windowSeconds        int
	blockDurationSeconds int
	entries              map[string]*authRateLimitEntry
	mu                   sync.RWMutex
}

// authRateLimitEntry tracks attempts for an IP
type authRateLimitEntry struct {
	attempts  int
	firstSeen time.Time
	blockedAt time.Time
}

// NewAuthRateLimitMiddleware creates a new rate limit middleware with sensible defaults
// for auth endpoints: 5 attempts per minute, 5 minute block after exceeding.
func NewAuthRateLimitMiddleware() *AuthRateLimitMiddleware {
	rl := &AuthRateLimitMiddleware{
		maxAttempts:          5,   // 5 attempts
		windowSeconds:        60,  // per minute
		blockDurationSeconds: 300, // block for 5 minutes after exceeding
		entries:              make(map[string]*authRateLimitEntry),
	}

	// Start cleanup goroutine to remove old entries
	go rl.cleanup()

	return rl
}

// cleanup periodically removes stale entries to prevent memory leaks
func (m *AuthRateLimitMiddleware) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for ip, entry := range m.entries {
			// Remove entries that are older than the block duration
			if now.Sub(entry.firstSeen) > time.Duration(m.blockDurationSeconds)*time.Second {
				delete(m.entries, ip)
			}
		}
		m.mu.Unlock()
	}
}

// getClientIP extracts the client IP from the request, checking proxy headers
func (m *AuthRateLimitMiddleware) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for reverse proxies)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take the first IP in the chain
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// Handle wraps an HTTP handler with rate limiting based on client IP.
// Returns 429 Too Many Requests if the rate limit is exceeded.
func (m *AuthRateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := m.getClientIP(r)

		m.mu.Lock()
		entry, exists := m.entries[ip]
		now := time.Now()

		if !exists {
			// First request from this IP
			m.entries[ip] = &authRateLimitEntry{
				attempts:  1,
				firstSeen: now,
			}
			m.mu.Unlock()
			next(w, r)
			return
		}

		// Check if currently blocked
		if !entry.blockedAt.IsZero() {
			blockExpiry := entry.blockedAt.Add(time.Duration(m.blockDurationSeconds) * time.Second)
			if now.Before(blockExpiry) {
				m.mu.Unlock()
				remainingSeconds := int(blockExpiry.Sub(now).Seconds())
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", fmt.Sprintf("%d", remainingSeconds))
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":         "Too many attempts. Please try again later.",
					"retry_after":   remainingSeconds,
					"blocked_until": blockExpiry.Format(time.RFC3339),
				})
				return
			}
			// Block expired, reset the entry
			entry.attempts = 1
			entry.firstSeen = now
			entry.blockedAt = time.Time{}
			m.mu.Unlock()
			next(w, r)
			return
		}

		// Check if window has expired
		windowExpiry := entry.firstSeen.Add(time.Duration(m.windowSeconds) * time.Second)
		if now.After(windowExpiry) {
			// Window expired, reset
			entry.attempts = 1
			entry.firstSeen = now
			m.mu.Unlock()
			next(w, r)
			return
		}

		// Within window, check attempts
		entry.attempts++
		if entry.attempts > m.maxAttempts {
			// Rate limit exceeded, block the IP
			entry.blockedAt = now
			m.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", fmt.Sprintf("%d", m.blockDurationSeconds))
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":         "Too many attempts. Please try again later.",
				"retry_after":   m.blockDurationSeconds,
				"blocked_until": now.Add(time.Duration(m.blockDurationSeconds) * time.Second).Format(time.RFC3339),
			})
			return
		}

		m.mu.Unlock()
		next(w, r)
	}
}

// RecordFailure allows logic layers to record a failed authentication attempt.
// This enables tracking failed logins even when the IP hasn't exceeded request limits.
func (m *AuthRateLimitMiddleware) RecordFailure(r *http.Request) {
	ip := m.getClientIP(r)

	m.mu.Lock()
	defer m.mu.Unlock()

	entry, exists := m.entries[ip]
	if !exists {
		m.entries[ip] = &authRateLimitEntry{
			attempts:  1,
			firstSeen: time.Now(),
		}
		return
	}

	// Don't increment if already blocked
	if !entry.blockedAt.IsZero() {
		return
	}

	entry.attempts++
	if entry.attempts > m.maxAttempts {
		entry.blockedAt = time.Now()
	}
}

// ResetForIP clears rate limit tracking for a specific IP (e.g., after successful login)
func (m *AuthRateLimitMiddleware) ResetForIP(r *http.Request) {
	ip := m.getClientIP(r)

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.entries, ip)
}
