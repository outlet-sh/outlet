package emailval

import (
	"testing"
)

func TestIsRoleBased(t *testing.T) {
	tests := []struct {
		name     string
		local    string
		expected bool
	}{
		// General contact
		{name: "info", local: "info", expected: true},
		{name: "contact", local: "contact", expected: true},
		{name: "hello", local: "hello", expected: true},
		{name: "hi", local: "hi", expected: true},
		{name: "enquiries", local: "enquiries", expected: true},
		{name: "inquiries", local: "inquiries", expected: true},

		// Support
		{name: "support", local: "support", expected: true},
		{name: "help", local: "help", expected: true},
		{name: "helpdesk", local: "helpdesk", expected: true},
		{name: "customerservice", local: "customerservice", expected: true},
		{name: "customer-service", local: "customer-service", expected: true},
		{name: "customercare", local: "customercare", expected: true},
		{name: "customer-care", local: "customer-care", expected: true},

		// Admin
		{name: "admin", local: "admin", expected: true},
		{name: "administrator", local: "administrator", expected: true},
		{name: "webmaster", local: "webmaster", expected: true},
		{name: "postmaster", local: "postmaster", expected: true},
		{name: "hostmaster", local: "hostmaster", expected: true},
		{name: "root", local: "root", expected: true},
		{name: "sysadmin", local: "sysadmin", expected: true},

		// Sales & Marketing
		{name: "sales", local: "sales", expected: true},
		{name: "marketing", local: "marketing", expected: true},
		{name: "press", local: "press", expected: true},
		{name: "media", local: "media", expected: true},
		{name: "pr", local: "pr", expected: true},
		{name: "partnerships", local: "partnerships", expected: true},
		{name: "partners", local: "partners", expected: true},
		{name: "advertising", local: "advertising", expected: true},
		{name: "ads", local: "ads", expected: true},

		// Team addresses
		{name: "team", local: "team", expected: true},
		{name: "staff", local: "staff", expected: true},
		{name: "office", local: "office", expected: true},
		{name: "hr", local: "hr", expected: true},
		{name: "jobs", local: "jobs", expected: true},
		{name: "careers", local: "careers", expected: true},
		{name: "recruiting", local: "recruiting", expected: true},
		{name: "recruitment", local: "recruitment", expected: true},

		// Technical
		{name: "noreply", local: "noreply", expected: true},
		{name: "no-reply", local: "no-reply", expected: true},
		{name: "donotreply", local: "donotreply", expected: true},
		{name: "do-not-reply", local: "do-not-reply", expected: true},
		{name: "mailer-daemon", local: "mailer-daemon", expected: true},
		{name: "abuse", local: "abuse", expected: true},
		{name: "security", local: "security", expected: true},
		{name: "privacy", local: "privacy", expected: true},
		{name: "legal", local: "legal", expected: true},
		{name: "compliance", local: "compliance", expected: true},

		// Financial
		{name: "billing", local: "billing", expected: true},
		{name: "invoices", local: "invoices", expected: true},
		{name: "accounts", local: "accounts", expected: true},
		{name: "accounting", local: "accounting", expected: true},
		{name: "finance", local: "finance", expected: true},
		{name: "payments", local: "payments", expected: true},

		// Subscriptions
		{name: "newsletter", local: "newsletter", expected: true},
		{name: "subscribe", local: "subscribe", expected: true},
		{name: "unsubscribe", local: "unsubscribe", expected: true},
		{name: "notifications", local: "notifications", expected: true},
		{name: "alerts", local: "alerts", expected: true},
		{name: "updates", local: "updates", expected: true},

		// Feedback
		{name: "feedback", local: "feedback", expected: true},
		{name: "suggestions", local: "suggestions", expected: true},
		{name: "complaints", local: "complaints", expected: true},

		// Personal names (not role-based)
		{name: "john", local: "john", expected: false},
		{name: "jane.doe", local: "jane.doe", expected: false},
		{name: "jsmith", local: "jsmith", expected: false},
		{name: "firstname.lastname", local: "firstname.lastname", expected: false},
		{name: "user123", local: "user123", expected: false},
		{name: "random", local: "random", expected: false},

		// Empty and edge cases
		{name: "empty string", local: "", expected: false},
		{name: "single character", local: "a", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRoleBased(tt.local)
			if result != tt.expected {
				t.Errorf("IsRoleBased(%q) = %v, want %v", tt.local, result, tt.expected)
			}
		})
	}
}

func TestIsRoleBased_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name     string
		local    string
		expected bool
	}{
		{name: "INFO uppercase", local: "INFO", expected: true},
		{name: "Admin mixed case", local: "Admin", expected: true},
		{name: "SUPPORT uppercase", local: "SUPPORT", expected: true},
		{name: "NoReply mixed case", local: "NoReply", expected: true},
		{name: "SALES uppercase", local: "SALES", expected: true},
		{name: "BiLLiNG random case", local: "BiLLiNG", expected: true},
		{name: "JOHN uppercase personal", local: "JOHN", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRoleBased(tt.local)
			if result != tt.expected {
				t.Errorf("IsRoleBased(%q) = %v, want %v", tt.local, result, tt.expected)
			}
		})
	}
}

func TestIsRoleBased_WithPlusTag(t *testing.T) {
	tests := []struct {
		name     string
		local    string
		expected bool
	}{
		{
			name:     "info with plus tag",
			local:    "info+newsletter",
			expected: true,
		},
		{
			name:     "admin with plus tag",
			local:    "admin+company",
			expected: true,
		},
		{
			name:     "support with complex tag",
			local:    "support+ticket-123",
			expected: true,
		},
		{
			name:     "noreply with tag",
			local:    "noreply+bounce",
			expected: true,
		},
		{
			name:     "personal name with plus tag",
			local:    "john+work",
			expected: false,
		},
		{
			name:     "random user with tag",
			local:    "randomuser+test",
			expected: false,
		},
		{
			name:     "just plus sign",
			local:    "+tag",
			expected: false,
		},
		{
			name:     "multiple plus signs - first part checked",
			local:    "info+tag+more",
			expected: true,
		},
		{
			name:     "uppercase with tag",
			local:    "SALES+campaign",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRoleBased(tt.local)
			if result != tt.expected {
				t.Errorf("IsRoleBased(%q) = %v, want %v", tt.local, result, tt.expected)
			}
		})
	}
}

func TestIsRoleBased_SimilarButNotRole(t *testing.T) {
	// These look similar to role-based but are not exact matches
	tests := []struct {
		name     string
		local    string
		expected bool
	}{
		{name: "infos plural", local: "infos", expected: false},
		{name: "myinfo prefix", local: "myinfo", expected: false},
		{name: "infomail", local: "infomail", expected: false},
		{name: "adminuser", local: "adminuser", expected: false},
		{name: "myadmin", local: "myadmin", expected: false},
		{name: "support2", local: "support2", expected: false},
		{name: "presales", local: "presales", expected: false},
		{name: "marketing1", local: "marketing1", expected: false},
		{name: "teama", local: "teama", expected: false},
		{name: "billing-dept", local: "billing-dept", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRoleBased(tt.local)
			if result != tt.expected {
				t.Errorf("IsRoleBased(%q) = %v, want %v (should not match as it's not an exact role prefix)", tt.local, result, tt.expected)
			}
		})
	}
}

func TestRolePrefixesCoverage(t *testing.T) {
	// Ensure all defined role prefixes are detected
	for _, prefix := range rolePrefixes {
		t.Run("prefix_"+prefix, func(t *testing.T) {
			if !IsRoleBased(prefix) {
				t.Errorf("IsRoleBased(%q) = false, but %q is in rolePrefixes", prefix, prefix)
			}
		})
	}
}

// BenchmarkIsRoleBased benchmarks the role-based check performance.
func BenchmarkIsRoleBased(b *testing.B) {
	b.Run("role_based_address", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			IsRoleBased("support")
		}
	})

	b.Run("personal_address", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			IsRoleBased("john.doe")
		}
	})

	b.Run("with_plus_tag", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			IsRoleBased("info+newsletter")
		}
	})
}
