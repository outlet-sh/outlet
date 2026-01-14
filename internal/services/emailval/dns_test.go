package emailval

import (
	"context"
	"testing"
	"time"
)

func TestCheckMX_KnownDomains(t *testing.T) {
	// Skip in short mode or CI to avoid network dependencies
	if testing.Short() {
		t.Skip("Skipping DNS tests in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		domain   string
		wantMX   bool
		wantErr  bool
	}{
		{
			name:    "gmail.com has MX records",
			domain:  "gmail.com",
			wantMX:  true,
			wantErr: false,
		},
		{
			name:    "yahoo.com has MX records",
			domain:  "yahoo.com",
			wantMX:  true,
			wantErr: false,
		},
		{
			name:    "non-existent domain",
			domain:  "this-domain-definitely-does-not-exist-12345.com",
			wantMX:  false,
			wantErr: false, // No error, just no MX records
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasMX, err := checkMX(ctx, tt.domain)

			if tt.wantErr && err == nil {
				t.Errorf("checkMX(%q) expected error, got nil", tt.domain)
				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("checkMX(%q) unexpected error: %v", tt.domain, err)
				return
			}

			if hasMX != tt.wantMX {
				t.Errorf("checkMX(%q) = %v, want %v", tt.domain, hasMX, tt.wantMX)
			}
		})
	}
}

func TestCheckMX_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := checkMX(ctx, "gmail.com")
	if err == nil {
		t.Log("checkMX with cancelled context may or may not return error depending on timing")
	}
}

func TestCheckMX_ContextTimeout(t *testing.T) {
	// Very short timeout to force timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Give context time to expire
	time.Sleep(1 * time.Millisecond)

	_, err := checkMX(ctx, "gmail.com")
	// With expired context, the lookup should fail or return quickly
	if err != nil {
		t.Logf("Expected behavior: error with expired context: %v", err)
	}
}

func TestGetMXHosts_KnownDomains(t *testing.T) {
	// Skip in short mode or CI to avoid network dependencies
	if testing.Short() {
		t.Skip("Skipping DNS tests in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name      string
		domain    string
		wantHosts bool // true if we expect at least one host
		wantErr   bool
	}{
		{
			name:      "gmail.com returns MX hosts",
			domain:    "gmail.com",
			wantHosts: true,
			wantErr:   false,
		},
		{
			name:      "outlook.com returns MX hosts",
			domain:    "outlook.com",
			wantHosts: true,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hosts, err := getMXHosts(ctx, tt.domain)

			if tt.wantErr && err == nil {
				t.Errorf("getMXHosts(%q) expected error, got nil", tt.domain)
				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("getMXHosts(%q) unexpected error: %v", tt.domain, err)
				return
			}

			if tt.wantHosts && len(hosts) == 0 {
				t.Errorf("getMXHosts(%q) returned no hosts, expected at least one", tt.domain)
			}

			if !tt.wantHosts && len(hosts) > 0 {
				t.Errorf("getMXHosts(%q) returned hosts, expected none", tt.domain)
			}

			// Verify hosts are non-empty strings
			for i, host := range hosts {
				if host == "" {
					t.Errorf("getMXHosts(%q) returned empty host at index %d", tt.domain, i)
				}
			}
		})
	}
}

func TestGetMXHosts_NonExistentDomain(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping DNS tests in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hosts, err := getMXHosts(ctx, "this-domain-definitely-does-not-exist-12345.com")

	// Should return error or empty hosts
	if err == nil && len(hosts) > 0 {
		t.Errorf("getMXHosts for non-existent domain returned hosts: %v", hosts)
	}
}

func TestGetMXHosts_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := getMXHosts(ctx, "gmail.com")
	if err == nil {
		t.Log("getMXHosts with cancelled context may or may not return error depending on timing")
	}
}

// TestMXHostsOrder verifies that MX hosts are returned (note: Go's LookupMX returns them sorted by preference)
func TestGetMXHosts_Order(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping DNS tests in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hosts, err := getMXHosts(ctx, "gmail.com")
	if err != nil {
		t.Skipf("Skipping order test due to DNS error: %v", err)
	}

	if len(hosts) < 2 {
		t.Skip("Need at least 2 MX hosts to test ordering")
	}

	// Just verify we got multiple hosts - Go's LookupMX returns them sorted by preference
	t.Logf("MX hosts for gmail.com (in order): %v", hosts)
}

// Integration test that uses both checkMX and getMXHosts
func TestMXFunctions_Consistency(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping DNS tests in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	domain := "gmail.com"

	hasMX, err := checkMX(ctx, domain)
	if err != nil {
		t.Fatalf("checkMX(%q) error: %v", domain, err)
	}

	hosts, err := getMXHosts(ctx, domain)
	if err != nil {
		t.Fatalf("getMXHosts(%q) error: %v", domain, err)
	}

	// If checkMX says there are MX records, getMXHosts should return hosts
	if hasMX && len(hosts) == 0 {
		t.Errorf("checkMX returned true but getMXHosts returned no hosts")
	}

	// If getMXHosts returns hosts, checkMX should return true
	if len(hosts) > 0 && !hasMX {
		t.Errorf("getMXHosts returned hosts but checkMX returned false")
	}
}

// BenchmarkCheckMX benchmarks DNS MX lookups
func BenchmarkCheckMX(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping DNS benchmark in short mode")
	}

	ctx := context.Background()

	b.Run("gmail.com", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			checkMX(ctx, "gmail.com")
		}
	})
}

// BenchmarkGetMXHosts benchmarks DNS MX host retrieval
func BenchmarkGetMXHosts(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping DNS benchmark in short mode")
	}

	ctx := context.Background()

	b.Run("gmail.com", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			getMXHosts(ctx, "gmail.com")
		}
	})
}
