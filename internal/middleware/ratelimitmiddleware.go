package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// RateLimitConfig defines rate limiting parameters
type RateLimitConfig struct {
	// MaxAttempts is the maximum number of attempts allowed within the window
	MaxAttempts int
	// WindowSeconds is the time window in seconds
	WindowSeconds int
	// BlockDurationSeconds is how long to block after exceeding the limit
	BlockDurationSeconds int
}

// DefaultAuthRateLimitConfig returns sensible defaults for auth endpoints
func DefaultAuthRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		MaxAttempts:          5,   // 5 attempts
		WindowSeconds:        60,  // per minute
		BlockDurationSeconds: 300, // block for 5 minutes after exceeding
	}
}

// rateLimitEntry tracks attempts for an IP
type rateLimitEntry struct {
	attempts  int
	firstSeen time.Time
	blockedAt time.Time
}

// RateLimitMiddleware provides rate limiting based on IP address
type RateLimitMiddleware struct {
	config  RateLimitConfig
	entries map[string]*rateLimitEntry
	mu      sync.RWMutex
}

// NewRateLimitMiddleware creates a new rate limit middleware with the given config
func NewRateLimitMiddleware(config RateLimitConfig) *RateLimitMiddleware {
	rl := &RateLimitMiddleware{
		config:  config,
		entries: make(map[string]*rateLimitEntry),
	}

	// Start cleanup goroutine to remove old entries
	go rl.cleanup()

	return rl
}

// cleanup periodically removes stale entries
func (m *RateLimitMiddleware) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for ip, entry := range m.entries {
			// Remove entries that are older than the block duration
			if now.Sub(entry.firstSeen) > time.Duration(m.config.BlockDurationSeconds)*time.Second {
				delete(m.entries, ip)
			}
		}
		m.mu.Unlock()
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for reverse proxies)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take the first IP in the chain
		if idx := len(xff); idx > 0 {
			for i := 0; i < len(xff); i++ {
				if xff[i] == ',' {
					return xff[:i]
				}
			}
			return xff
		}
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// Handle wraps an HTTP handler with rate limiting
func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		m.mu.Lock()
		entry, exists := m.entries[ip]
		now := time.Now()

		if !exists {
			// First request from this IP
			m.entries[ip] = &rateLimitEntry{
				attempts:  1,
				firstSeen: now,
			}
			m.mu.Unlock()
			next(w, r)
			return
		}

		// Check if currently blocked
		if !entry.blockedAt.IsZero() {
			blockExpiry := entry.blockedAt.Add(time.Duration(m.config.BlockDurationSeconds) * time.Second)
			if now.Before(blockExpiry) {
				m.mu.Unlock()
				remainingSeconds := int(blockExpiry.Sub(now).Seconds())
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", string(rune(remainingSeconds)))
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
		windowExpiry := entry.firstSeen.Add(time.Duration(m.config.WindowSeconds) * time.Second)
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
		if entry.attempts > m.config.MaxAttempts {
			// Rate limit exceeded, block the IP
			entry.blockedAt = now
			m.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", string(rune(m.config.BlockDurationSeconds)))
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":         "Too many attempts. Please try again later.",
				"retry_after":   m.config.BlockDurationSeconds,
				"blocked_until": now.Add(time.Duration(m.config.BlockDurationSeconds) * time.Second).Format(time.RFC3339),
			})
			return
		}

		m.mu.Unlock()
		next(w, r)
	}
}

// RecordFailure allows logic layers to record a failed authentication attempt
// This allows the rate limiter to track failed logins even for correct IPs with wrong passwords
func (m *RateLimitMiddleware) RecordFailure(r *http.Request) {
	ip := getClientIP(r)

	m.mu.Lock()
	defer m.mu.Unlock()

	entry, exists := m.entries[ip]
	if !exists {
		m.entries[ip] = &rateLimitEntry{
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
	if entry.attempts > m.config.MaxAttempts {
		entry.blockedAt = time.Now()
	}
}

// ResetForIP clears rate limit tracking for a specific IP (e.g., after successful login)
func (m *RateLimitMiddleware) ResetForIP(r *http.Request) {
	ip := getClientIP(r)

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.entries, ip)
}
