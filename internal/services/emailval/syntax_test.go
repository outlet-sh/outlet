package emailval

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase conversion",
			input:    "TEST@EXAMPLE.COM",
			expected: "test@example.com",
		},
		{
			name:     "mixed case",
			input:    "Test.User@Example.COM",
			expected: "test.user@example.com",
		},
		{
			name:     "leading whitespace",
			input:    "  user@example.com",
			expected: "user@example.com",
		},
		{
			name:     "trailing whitespace",
			input:    "user@example.com  ",
			expected: "user@example.com",
		},
		{
			name:     "both sides whitespace",
			input:    "  user@example.com  ",
			expected: "user@example.com",
		},
		{
			name:     "tabs and newlines",
			input:    "\t\nuser@example.com\t\n",
			expected: "user@example.com",
		},
		{
			name:     "already normalized",
			input:    "user@example.com",
			expected: "user@example.com",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   \t\n   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalize(tt.input)
			if result != tt.expected {
				t.Errorf("normalize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateSyntax(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		wantLocal   string
		wantDomain  string
		shouldError bool
		errorMsg    string
	}{
		// Valid emails
		{
			name:       "simple valid email",
			email:      "user@example.com",
			wantLocal:  "user",
			wantDomain: "example.com",
		},
		{
			name:       "email with dots in local part",
			email:      "user.name@example.com",
			wantLocal:  "user.name",
			wantDomain: "example.com",
		},
		{
			name:       "email with plus tag",
			email:      "user+tag@example.com",
			wantLocal:  "user+tag",
			wantDomain: "example.com",
		},
		{
			name:       "email with numbers",
			email:      "user123@example.com",
			wantLocal:  "user123",
			wantDomain: "example.com",
		},
		{
			name:       "email with subdomain",
			email:      "user@mail.example.com",
			wantLocal:  "user",
			wantDomain: "mail.example.com",
		},
		{
			name:       "email with hyphen in domain",
			email:      "user@my-example.com",
			wantLocal:  "user",
			wantDomain: "my-example.com",
		},
		{
			name:       "email with underscore in local",
			email:      "user_name@example.com",
			wantLocal:  "user_name",
			wantDomain: "example.com",
		},
		{
			name:       "email with multiple plus signs",
			email:      "user+tag+more@example.com",
			wantLocal:  "user+tag+more",
			wantDomain: "example.com",
		},

		// Invalid emails - missing parts
		{
			name:        "empty string",
			email:       "",
			shouldError: true,
			errorMsg:    "does not match email pattern",
		},
		{
			name:        "no @ symbol",
			email:       "userexample.com",
			shouldError: true,
			errorMsg:    "does not match email pattern",
		},
		{
			name:        "multiple @ symbols",
			email:       "user@example@com",
			shouldError: true,
			errorMsg:    "does not match email pattern", // Regex catches this before splitEmail
		},
		{
			name:        "no domain",
			email:       "user@",
			shouldError: true,
			errorMsg:    "does not match email pattern",
		},
		{
			name:        "no local part",
			email:       "@example.com",
			shouldError: true,
			errorMsg:    "does not match email pattern",
		},
		{
			name:        "no TLD",
			email:       "user@example",
			shouldError: true,
			errorMsg:    "does not match email pattern",
		},

		// Invalid emails - whitespace
		{
			name:        "space in local part",
			email:       "user name@example.com",
			shouldError: true,
			errorMsg:    "does not match email pattern",
		},
		{
			name:        "space in domain",
			email:       "user@example .com",
			shouldError: true,
			errorMsg:    "does not match email pattern",
		},

		// Invalid emails - dot issues
		{
			name:        "consecutive dots in local",
			email:       "user..name@example.com",
			shouldError: true,
			errorMsg:    "consecutive dots",
		},
		{
			name:        "leading dot in local",
			email:       ".user@example.com",
			shouldError: true,
			errorMsg:    "starts with dot",
		},
		{
			name:        "trailing dot in local",
			email:       "user.@example.com",
			shouldError: true,
			errorMsg:    "ends with dot",
		},

		// Edge cases
		{
			name:        "just @",
			email:       "@",
			shouldError: true,
			errorMsg:    "does not match email pattern",
		},
		{
			name:        "only dots",
			email:       "...@example.com",
			shouldError: true,
			errorMsg:    "consecutive dots", // Consecutive dots error is checked first
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			local, domain, err := validateSyntax(tt.email)

			if tt.shouldError {
				if err == nil {
					t.Errorf("validateSyntax(%q) expected error containing %q, got nil", tt.email, tt.errorMsg)
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("validateSyntax(%q) error = %q, want error containing %q", tt.email, err.Error(), tt.errorMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("validateSyntax(%q) unexpected error: %v", tt.email, err)
				return
			}

			if local != tt.wantLocal {
				t.Errorf("validateSyntax(%q) local = %q, want %q", tt.email, local, tt.wantLocal)
			}
			if domain != tt.wantDomain {
				t.Errorf("validateSyntax(%q) domain = %q, want %q", tt.email, domain, tt.wantDomain)
			}
		})
	}
}

func TestSplitEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		wantLocal   string
		wantDomain  string
		shouldError bool
		errorMsg    string
	}{
		{
			name:       "valid email",
			email:      "user@example.com",
			wantLocal:  "user",
			wantDomain: "example.com",
		},
		{
			name:       "complex local part",
			email:      "user.name+tag@subdomain.example.com",
			wantLocal:  "user.name+tag",
			wantDomain: "subdomain.example.com",
		},
		{
			name:        "no @ symbol",
			email:       "userexample.com",
			shouldError: true,
			errorMsg:    "must contain single @",
		},
		{
			name:        "multiple @ symbols",
			email:       "user@domain@example.com",
			shouldError: true,
			errorMsg:    "must contain single @",
		},
		{
			name:        "empty local part",
			email:       "@example.com",
			shouldError: true,
			errorMsg:    "empty local-part",
		},
		{
			name:        "empty domain",
			email:       "user@",
			shouldError: true,
			errorMsg:    "empty domain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			local, domain, err := splitEmail(tt.email)

			if tt.shouldError {
				if err == nil {
					t.Errorf("splitEmail(%q) expected error, got nil", tt.email)
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("splitEmail(%q) error = %q, want error containing %q", tt.email, err.Error(), tt.errorMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("splitEmail(%q) unexpected error: %v", tt.email, err)
				return
			}

			if local != tt.wantLocal {
				t.Errorf("splitEmail(%q) local = %q, want %q", tt.email, local, tt.wantLocal)
			}
			if domain != tt.wantDomain {
				t.Errorf("splitEmail(%q) domain = %q, want %q", tt.email, domain, tt.wantDomain)
			}
		})
	}
}

// contains checks if substr is in s (case-insensitive for error messages).
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
