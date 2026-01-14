package emailval

import (
	"strings"
)

// rolePrefixes are common role-based email prefixes.
// These addresses typically represent a function rather than a person.
var rolePrefixes = []string{
	// General contact
	"info",
	"contact",
	"hello",
	"hi",
	"enquiries",
	"inquiries",

	// Support
	"support",
	"help",
	"helpdesk",
	"customerservice",
	"customer-service",
	"customercare",
	"customer-care",

	// Admin
	"admin",
	"administrator",
	"webmaster",
	"postmaster",
	"hostmaster",
	"root",
	"sysadmin",

	// Sales & Marketing
	"sales",
	"marketing",
	"press",
	"media",
	"pr",
	"partnerships",
	"partners",
	"advertising",
	"ads",

	// Team addresses
	"team",
	"staff",
	"office",
	"hr",
	"jobs",
	"careers",
	"recruiting",
	"recruitment",

	// Technical
	"noreply",
	"no-reply",
	"donotreply",
	"do-not-reply",
	"mailer-daemon",
	"abuse",
	"security",
	"privacy",
	"legal",
	"compliance",

	// Financial
	"billing",
	"invoices",
	"accounts",
	"accounting",
	"finance",
	"payments",

	// Subscriptions
	"newsletter",
	"subscribe",
	"unsubscribe",
	"notifications",
	"alerts",
	"updates",

	// Feedback
	"feedback",
	"suggestions",
	"complaints",
}

// IsRoleBased checks if the local part of an email is a role-based address.
func IsRoleBased(local string) bool {
	// Strip +tag if present (e.g., info+tag@example.com -> info).
	if idx := strings.Index(local, "+"); idx > 0 {
		local = local[:idx]
	}
	local = strings.ToLower(local)

	for _, prefix := range rolePrefixes {
		if local == prefix {
			return true
		}
	}
	return false
}
