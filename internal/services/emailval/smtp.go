package emailval

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"sync"
	"time"
)

// SMTPDialer defines how to connect to SMTP servers.
type SMTPDialer interface {
	DialSMTP(ctx context.Context, address string) (*textproto.Conn, error)
}

// DefaultSMTPDialer dials plain TCP with a timeout.
type DefaultSMTPDialer struct {
	Timeout time.Duration
}

// DialSMTP connects to an SMTP server.
func (d DefaultSMTPDialer) DialSMTP(ctx context.Context, address string) (*textproto.Conn, error) {
	timeout := d.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	dialer := &net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, err
	}
	return textproto.NewConn(conn), nil
}

// smtpRateLimiter provides per-domain rate limiting to avoid blacklisting.
var smtpRateLimiter = &rateLimiter{
	domains:     make(map[string]time.Time),
	minInterval: 2 * time.Second, // Minimum 2 seconds between probes to same domain
}

type rateLimiter struct {
	mu          sync.Mutex
	domains     map[string]time.Time
	minInterval time.Duration
}

func (r *rateLimiter) wait(ctx context.Context, domain string) error {
	r.mu.Lock()
	lastProbe, exists := r.domains[domain]
	r.mu.Unlock()

	if exists {
		waitTime := r.minInterval - time.Since(lastProbe)
		if waitTime > 0 {
			select {
			case <-time.After(waitTime):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	r.mu.Lock()
	r.domains[domain] = time.Now()
	r.mu.Unlock()

	return nil
}

// SMTPConfig provides configuration for SMTP verification.
// Use this to customize behavior and avoid blacklisting.
type SMTPConfig struct {
	// HeloHostname is the hostname to use in HELO/EHLO.
	// Should be a valid, resolvable hostname you control.
	// If empty, defaults to "localhost".
	HeloHostname string

	// FromDomain is the domain to use in MAIL FROM.
	// Should be a domain you control with valid SPF.
	// If empty, uses HeloHostname.
	FromDomain string

	// SkipRateLimiting disables the built-in rate limiter.
	// Only set this if you implement your own rate limiting.
	SkipRateLimiting bool

	// MaxRetries is the number of MX servers to try before giving up.
	// Defaults to 2.
	MaxRetries int
}

// DefaultSMTPConfig returns a safe default SMTP configuration.
func DefaultSMTPConfig() *SMTPConfig {
	return &SMTPConfig{
		HeloHostname: "localhost",
		FromDomain:   "localhost",
		MaxRetries:   2,
	}
}

// checkSMTPMailbox verifies if a mailbox exists via SMTP.
// It connects to the domain's MX servers and checks RCPT TO response.
// This never sends actual email, only probes deliverability.
//
// IMPORTANT: SMTP verification should be used sparingly to avoid blacklisting.
// Consider:
// - Rate limiting (built-in, 2s between probes to same domain)
// - Only checking high-value submissions
// - Caching results to avoid repeated probes
// - Using a dedicated IP for verification
func checkSMTPMailbox(ctx context.Context, domain, fullEmail string, opts *Options) (bool, error) {
	mxHosts, err := getMXHosts(ctx, domain)
	if err != nil || len(mxHosts) == 0 {
		return false, fmt.Errorf("no MX records for SMTP check: %w", err)
	}

	// Apply rate limiting to prevent blacklisting.
	if err := smtpRateLimiter.wait(ctx, domain); err != nil {
		return false, fmt.Errorf("rate limit wait cancelled: %w", err)
	}

	dialer := opts.SMTPDialer
	if dialer == nil {
		dialer = DefaultSMTPDialer{Timeout: opts.Timeout}
	}

	smtpCfg := DefaultSMTPConfig()
	maxRetries := smtpCfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 2
	}

	// Try MX hosts in order until one responds cleanly.
	var lastErr error
	tried := 0
	for _, host := range mxHosts {
		if tried >= maxRetries {
			break
		}
		tried++

		// Remove trailing dot from MX host.
		host = strings.TrimSuffix(host, ".")
		addr := net.JoinHostPort(host, "25")

		deliverable, err := smtpProbe(ctx, dialer, addr, fullEmail, smtpCfg)
		if err == nil {
			return deliverable, nil
		}

		// Check if this is a greylisting response - don't count as failure.
		if isGreylistError(err) {
			// Greylisting means the server is probably valid but wants us to retry later.
			// Treat as "unknown but likely valid".
			return true, nil
		}

		lastErr = err
	}

	if lastErr != nil {
		return false, fmt.Errorf("all MX probes failed: %w", lastErr)
	}
	return false, errors.New("all MX probes failed")
}

// isGreylistError checks if an error looks like a greylisting response.
func isGreylistError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	// Common greylisting indicators.
	return strings.Contains(msg, "451") ||
		strings.Contains(msg, "greylist") ||
		strings.Contains(msg, "try again") ||
		strings.Contains(msg, "temporarily")
}

// smtpProbe performs a single SMTP probe to check if an address is deliverable.
func smtpProbe(ctx context.Context, dialer SMTPDialer, addr, target string, cfg *SMTPConfig) (bool, error) {
	conn, err := dialer.DialSMTP(ctx, addr)
	if err != nil {
		return false, fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close()

	// Read server greeting.
	if err := expectCode(conn, "220"); err != nil {
		return false, fmt.Errorf("greeting: %w", err)
	}

	// Send EHLO first (preferred), fall back to HELO.
	heloHost := cfg.HeloHostname
	if heloHost == "" {
		heloHost = "localhost"
	}

	if err := conn.PrintfLine("EHLO %s", heloHost); err != nil {
		return false, fmt.Errorf("EHLO send: %w", err)
	}

	// EHLO might have multi-line response, read until we get final line.
	line, err := readFinalResponse(conn)
	if err != nil {
		return false, fmt.Errorf("EHLO response: %w", err)
	}

	// If EHLO failed, try HELO.
	if !strings.HasPrefix(line, "250") {
		if err := conn.PrintfLine("HELO %s", heloHost); err != nil {
			return false, fmt.Errorf("HELO send: %w", err)
		}
		if err := expectCode(conn, "250"); err != nil {
			return false, fmt.Errorf("HELO response: %w", err)
		}
	}

	// Send MAIL FROM with a proper address.
	fromDomain := cfg.FromDomain
	if fromDomain == "" {
		fromDomain = heloHost
	}
	if err := conn.PrintfLine("MAIL FROM:<verify@%s>", fromDomain); err != nil {
		return false, fmt.Errorf("MAIL FROM send: %w", err)
	}
	if err := expectCode(conn, "250"); err != nil {
		return false, fmt.Errorf("MAIL FROM response: %w", err)
	}

	// Send RCPT TO - this is the key check.
	if err := conn.PrintfLine("RCPT TO:<%s>", target); err != nil {
		return false, fmt.Errorf("RCPT TO send: %w", err)
	}

	line, err = conn.ReadLine()
	if err != nil {
		return false, fmt.Errorf("RCPT TO response: %w", err)
	}

	// Send QUIT regardless of result.
	_ = conn.PrintfLine("QUIT")

	// 250, 251 = OK (mailbox exists or will be forwarded).
	// 252 = Cannot verify, but will attempt delivery.
	if strings.HasPrefix(line, "250") || strings.HasPrefix(line, "251") || strings.HasPrefix(line, "252") {
		return true, nil
	}

	// 550, 551, 552, 553 = Mailbox does not exist or rejected.
	if strings.HasPrefix(line, "55") {
		return false, nil
	}

	// 450, 451, 452 = Temporary failures (greylisting, rate limiting).
	if strings.HasPrefix(line, "45") {
		return false, fmt.Errorf("temporary rejection (greylist): %s", line)
	}

	// For other codes, treat as unknown.
	return false, fmt.Errorf("unexpected RCPT response: %s", line)
}

// readFinalResponse reads SMTP responses until it gets the final line.
// SMTP multi-line responses have "250-" for continuation and "250 " for final.
func readFinalResponse(conn *textproto.Conn) (string, error) {
	for {
		line, err := conn.ReadLine()
		if err != nil {
			return "", err
		}
		// Check if this is the final line (code followed by space, not hyphen).
		if len(line) >= 4 && line[3] == ' ' {
			return line, nil
		}
		// Continue reading if it's a continuation line (code followed by hyphen).
	}
}

// expectCode reads a line and checks if it starts with the expected code.
func expectCode(conn *textproto.Conn, code string) error {
	line, err := conn.ReadLine()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(line, code) {
		return fmt.Errorf("expected %s, got: %s", code, line)
	}
	return nil
}
