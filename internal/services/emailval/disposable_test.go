package emailval

import (
	"testing"
)

func TestInMemoryDisposableProvider_IsDisposable(t *testing.T) {
	provider := NewStaticDisposableProvider([]string{
		"mailinator.com",
		"10minutemail.com",
		"TEMPMAIL.COM", // Test case insensitivity
	})

	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		{
			name:     "known disposable domain",
			domain:   "mailinator.com",
			expected: true,
		},
		{
			name:     "known disposable uppercase input",
			domain:   "MAILINATOR.COM",
			expected: true,
		},
		{
			name:     "known disposable mixed case input",
			domain:   "Mailinator.Com",
			expected: true,
		},
		{
			name:     "another known disposable",
			domain:   "10minutemail.com",
			expected: true,
		},
		{
			name:     "uppercase in list matched lowercase",
			domain:   "tempmail.com",
			expected: true,
		},
		{
			name:     "legitimate domain",
			domain:   "gmail.com",
			expected: false,
		},
		{
			name:     "legitimate domain outlook",
			domain:   "outlook.com",
			expected: false,
		},
		{
			name:     "custom domain",
			domain:   "mycompany.com",
			expected: false,
		},
		{
			name:     "empty domain",
			domain:   "",
			expected: false,
		},
		{
			name:     "subdomain of disposable",
			domain:   "subdomain.mailinator.com",
			expected: false, // Provider uses exact match
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.IsDisposable(tt.domain)
			if result != tt.expected {
				t.Errorf("IsDisposable(%q) = %v, want %v", tt.domain, result, tt.expected)
			}
		})
	}
}

func TestNewStaticDisposableProvider(t *testing.T) {
	tests := []struct {
		name         string
		domains      []string
		testDomain   string
		shouldMatch  bool
	}{
		{
			name:        "empty list",
			domains:     []string{},
			testDomain:  "mailinator.com",
			shouldMatch: false,
		},
		{
			name:        "single domain",
			domains:     []string{"test.com"},
			testDomain:  "test.com",
			shouldMatch: true,
		},
		{
			name:        "domains with whitespace",
			domains:     []string{"  test.com  ", "\tother.com\n"},
			testDomain:  "test.com",
			shouldMatch: true,
		},
		{
			name:        "domains with whitespace - other domain",
			domains:     []string{"  test.com  ", "\tother.com\n"},
			testDomain:  "other.com",
			shouldMatch: true,
		},
		{
			name:        "nil list",
			domains:     nil,
			testDomain:  "test.com",
			shouldMatch: false,
		},
		{
			name:        "duplicate domains",
			domains:     []string{"test.com", "test.com", "TEST.COM"},
			testDomain:  "test.com",
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewStaticDisposableProvider(tt.domains)
			result := provider.IsDisposable(tt.testDomain)
			if result != tt.shouldMatch {
				t.Errorf("NewStaticDisposableProvider(%v).IsDisposable(%q) = %v, want %v",
					tt.domains, tt.testDomain, result, tt.shouldMatch)
			}
		})
	}
}

func TestDefaultDisposableProvider(t *testing.T) {
	provider := DefaultDisposableProvider()

	// Test that the default provider has some expected domains
	knownDisposableDomains := []string{
		"mailinator.com",
		"guerrillamail.com",
		"10minutemail.com",
		"tempmail.com",
		"yopmail.com",
		"maildrop.cc",
		"trashmail.com",
	}

	for _, domain := range knownDisposableDomains {
		t.Run("known_"+domain, func(t *testing.T) {
			if !provider.IsDisposable(domain) {
				t.Errorf("DefaultDisposableProvider().IsDisposable(%q) = false, want true", domain)
			}
		})
	}

	// Test that legitimate domains are not marked as disposable
	legitimateDomains := []string{
		"gmail.com",
		"yahoo.com",
		"outlook.com",
		"hotmail.com",
		"icloud.com",
		"protonmail.com",
		"example.com",
		"company.co.uk",
	}

	for _, domain := range legitimateDomains {
		t.Run("legitimate_"+domain, func(t *testing.T) {
			if provider.IsDisposable(domain) {
				t.Errorf("DefaultDisposableProvider().IsDisposable(%q) = true, want false", domain)
			}
		})
	}
}

func TestDefaultDisposableDomainsCompleteness(t *testing.T) {
	// Verify the default list contains expected entries
	expectedInList := map[string]bool{
		"mailinator.com":     true,
		"guerrillamail.com":  true,
		"10minutemail.com":   true,
		"tempmail.com":       true,
		"yopmail.com":        true,
		"sharklasers.com":    true,
		"burnermail.io":      true,
	}

	provider := DefaultDisposableProvider()

	for domain, expected := range expectedInList {
		result := provider.IsDisposable(domain)
		if result != expected {
			t.Errorf("defaultDisposableDomains should contain %q: got %v, want %v", domain, result, expected)
		}
	}
}

// BenchmarkIsDisposable benchmarks the disposable check performance.
func BenchmarkIsDisposable(b *testing.B) {
	provider := DefaultDisposableProvider()

	b.Run("disposable_domain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			provider.IsDisposable("mailinator.com")
		}
	})

	b.Run("legitimate_domain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			provider.IsDisposable("gmail.com")
		}
	})
}
