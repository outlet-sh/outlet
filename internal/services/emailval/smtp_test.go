package emailval

import (
	"context"
	"errors"
	"net/textproto"
	"testing"
	"time"
)

func TestIsGreylistError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "451 error code",
			err:      errors.New("451 Try again later"),
			expected: true,
		},
		{
			name:     "greylist in message",
			err:      errors.New("Greylisted, please retry"),
			expected: true,
		},
		{
			name:     "try again message",
			err:      errors.New("Please try again in a few minutes"),
			expected: true,
		},
		{
			name:     "temporarily unavailable",
			err:      errors.New("Service temporarily unavailable"),
			expected: true,
		},
		{
			name:     "regular error",
			err:      errors.New("connection refused"),
			expected: false,
		},
		{
			name:     "550 user not found",
			err:      errors.New("550 User not found"),
			expected: false,
		},
		{
			name:     "timeout error",
			err:      errors.New("connection timeout"),
			expected: false,
		},
		{
			name:     "mixed case greylist",
			err:      errors.New("GREYLISTED - retry later"),
			expected: true,
		},
		{
			name:     "451 with greylist",
			err:      errors.New("451 Greylisting in effect"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isGreylistError(tt.err)
			if result != tt.expected {
				t.Errorf("isGreylistError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestDefaultSMTPConfig(t *testing.T) {
	cfg := DefaultSMTPConfig()

	if cfg == nil {
		t.Fatal("DefaultSMTPConfig() returned nil")
	}

	if cfg.HeloHostname != "localhost" {
		t.Errorf("HeloHostname = %q, want %q", cfg.HeloHostname, "localhost")
	}

	if cfg.FromDomain != "localhost" {
		t.Errorf("FromDomain = %q, want %q", cfg.FromDomain, "localhost")
	}

	if cfg.MaxRetries != 2 {
		t.Errorf("MaxRetries = %d, want %d", cfg.MaxRetries, 2)
	}

	if cfg.SkipRateLimiting != false {
		t.Errorf("SkipRateLimiting = %v, want %v", cfg.SkipRateLimiting, false)
	}
}

func TestDefaultSMTPDialer(t *testing.T) {
	dialer := DefaultSMTPDialer{Timeout: 1 * time.Second}

	// Test with a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := dialer.DialSMTP(ctx, "localhost:25")
	if err == nil {
		t.Error("Expected error with cancelled context, got nil")
	}
}

func TestDefaultSMTPDialer_DefaultTimeout(t *testing.T) {
	dialer := DefaultSMTPDialer{} // No timeout set

	// The default timeout should be 10 seconds
	// We can't easily test the actual value without reflection,
	// but we can verify it doesn't panic
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Should fail quickly due to context timeout
	_, err := dialer.DialSMTP(ctx, "localhost:25")
	if err == nil {
		t.Error("Expected error with very short context timeout")
	}
}

func TestRateLimiter_Wait(t *testing.T) {
	rl := &rateLimiter{
		domains:     make(map[string]time.Time),
		minInterval: 10 * time.Millisecond,
	}

	ctx := context.Background()

	// First call should not wait
	start := time.Now()
	err := rl.wait(ctx, "test.com")
	if err != nil {
		t.Fatalf("First wait returned error: %v", err)
	}
	elapsed := time.Since(start)
	if elapsed > 5*time.Millisecond {
		t.Errorf("First wait took %v, expected near-instant", elapsed)
	}

	// Second call to same domain should wait
	start = time.Now()
	err = rl.wait(ctx, "test.com")
	if err != nil {
		t.Fatalf("Second wait returned error: %v", err)
	}
	elapsed = time.Since(start)
	if elapsed < 5*time.Millisecond {
		t.Errorf("Second wait took %v, expected at least ~10ms", elapsed)
	}

	// Call to different domain should not wait
	start = time.Now()
	err = rl.wait(ctx, "other.com")
	if err != nil {
		t.Fatalf("Wait for other domain returned error: %v", err)
	}
	elapsed = time.Since(start)
	if elapsed > 5*time.Millisecond {
		t.Errorf("Wait for other domain took %v, expected near-instant", elapsed)
	}
}

func TestRateLimiter_WaitCancelled(t *testing.T) {
	rl := &rateLimiter{
		domains:     make(map[string]time.Time),
		minInterval: 1 * time.Second, // Long interval
	}

	// First call to set the timestamp
	ctx := context.Background()
	err := rl.wait(ctx, "test.com")
	if err != nil {
		t.Fatalf("First wait returned error: %v", err)
	}

	// Second call with cancelled context should return immediately with error
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	start := time.Now()
	err = rl.wait(ctx, "test.com")
	elapsed := time.Since(start)

	if err == nil {
		t.Error("Expected error with cancelled context")
	}

	if elapsed > 100*time.Millisecond {
		t.Errorf("Wait with cancelled context took %v, expected immediate return", elapsed)
	}
}

// MockSMTPDialer is a test double for SMTPDialer
type MockSMTPDialer struct {
	DialFunc func(ctx context.Context, address string) (*textproto.Conn, error)
}

func (m MockSMTPDialer) DialSMTP(ctx context.Context, address string) (*textproto.Conn, error) {
	if m.DialFunc != nil {
		return m.DialFunc(ctx, address)
	}
	return nil, errors.New("dial not implemented")
}

func TestSMTPDialer_Interface(t *testing.T) {
	// Verify DefaultSMTPDialer implements SMTPDialer
	var _ SMTPDialer = DefaultSMTPDialer{}
	var _ SMTPDialer = &MockSMTPDialer{}
}

func TestSMTPConfig_Fields(t *testing.T) {
	cfg := &SMTPConfig{
		HeloHostname:     "mail.example.com",
		FromDomain:       "example.com",
		SkipRateLimiting: true,
		MaxRetries:       5,
	}

	if cfg.HeloHostname != "mail.example.com" {
		t.Errorf("HeloHostname = %q, want %q", cfg.HeloHostname, "mail.example.com")
	}

	if cfg.FromDomain != "example.com" {
		t.Errorf("FromDomain = %q, want %q", cfg.FromDomain, "example.com")
	}

	if cfg.SkipRateLimiting != true {
		t.Errorf("SkipRateLimiting = %v, want %v", cfg.SkipRateLimiting, true)
	}

	if cfg.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d, want %d", cfg.MaxRetries, 5)
	}
}

// Note: Testing checkSMTPMailbox and smtpProbe would require network access
// or more sophisticated mocking of the textproto.Conn. These are integration
// tests that should be run in a controlled environment.

func TestExpectCode_EdgeCases(t *testing.T) {
	// This tests the expectCode function behavior conceptually
	// Actual testing would require a mock textproto.Conn

	// Document expected behavior:
	// - expectCode("220") should accept "220 Welcome"
	// - expectCode("250") should accept "250 OK"
	// - expectCode("250") should reject "251 Forward"
}

func TestReadFinalResponse_EdgeCases(t *testing.T) {
	// This tests the readFinalResponse function behavior conceptually
	// Actual testing would require a mock textproto.Conn

	// Document expected behavior:
	// - "250-First line" followed by "250 Last line" should return "250 Last line"
	// - "250 Single line" should return immediately
}

// BenchmarkRateLimiter benchmarks the rate limiter
func BenchmarkRateLimiter(b *testing.B) {
	rl := &rateLimiter{
		domains:     make(map[string]time.Time),
		minInterval: 1 * time.Nanosecond, // Very short for benchmarking
	}
	ctx := context.Background()

	b.Run("same_domain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rl.wait(ctx, "bench.com")
		}
	})

	b.Run("different_domains", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rl.wait(ctx, "domain"+string(rune(i%26+65))+".com")
		}
	})
}
