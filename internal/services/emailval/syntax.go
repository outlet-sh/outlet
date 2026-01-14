package emailval

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// Simple, user-friendly syntax regex (not full RFC 5322).
	// Matches: local@domain.tld where local and domain have no spaces.
	simpleRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
)

// normalize lowercases and trims whitespace from an email.
func normalize(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// validateSyntax checks email syntax and returns local and domain parts.
func validateSyntax(email string) (local, domain string, err error) {
	if !simpleRegex.MatchString(email) {
		return "", "", errors.New("does not match email pattern")
	}

	local, domain, err = splitEmail(email)
	if err != nil {
		return "", "", err
	}

	// Check for invalid local-part structures.
	if strings.Contains(local, "..") {
		return "", "", errors.New("local-part contains consecutive dots")
	}
	if strings.HasPrefix(local, ".") {
		return "", "", errors.New("local-part starts with dot")
	}
	if strings.HasSuffix(local, ".") {
		return "", "", errors.New("local-part ends with dot")
	}

	return local, domain, nil
}

// splitEmail splits an email into local and domain parts.
func splitEmail(email string) (local, domain string, err error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", "", errors.New("must contain single @")
	}
	local = parts[0]
	domain = parts[1]
	if local == "" {
		return "", "", errors.New("empty local-part")
	}
	if domain == "" {
		return "", "", errors.New("empty domain")
	}
	return local, domain, nil
}
