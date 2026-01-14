package emailval

import (
	"strings"
)

// DisposableProvider reports whether a domain is disposable/temporary.
type DisposableProvider interface {
	IsDisposable(domain string) bool
}

// inMemoryDisposableProvider is a map-backed provider.
type inMemoryDisposableProvider struct {
	set map[string]struct{}
}

func (p inMemoryDisposableProvider) IsDisposable(domain string) bool {
	_, ok := p.set[strings.ToLower(domain)]
	return ok
}

// NewStaticDisposableProvider builds a provider from a list of domains.
func NewStaticDisposableProvider(domains []string) DisposableProvider {
	set := make(map[string]struct{}, len(domains))
	for _, d := range domains {
		set[strings.ToLower(strings.TrimSpace(d))] = struct{}{}
	}
	return inMemoryDisposableProvider{set: set}
}

// DefaultDisposableProvider returns a provider with a minimal default list.
// In production, you should inject your own with a larger, periodically refreshed list.
func DefaultDisposableProvider() DisposableProvider {
	return NewStaticDisposableProvider(defaultDisposableDomains)
}

// defaultDisposableDomains is a minimal list of known disposable email domains.
// Expand this list or use an external service for comprehensive coverage.
var defaultDisposableDomains = []string{
	// Popular disposable services
	"mailinator.com",
	"guerrillamail.com",
	"guerrillamail.org",
	"guerrillamail.net",
	"guerrillamail.biz",
	"guerrillamail.de",
	"10minutemail.com",
	"10minutemail.net",
	"tempmail.com",
	"tempmail.net",
	"temp-mail.org",
	"throwaway.email",
	"throwawaymail.com",
	"fakeinbox.com",
	"fakemailgenerator.com",
	"getnada.com",
	"getairmail.com",
	"dispostable.com",
	"mailnesia.com",
	"maildrop.cc",
	"yopmail.com",
	"yopmail.fr",
	"sharklasers.com",
	"trashmail.com",
	"trashmail.net",
	"trashmail.org",
	"spamgourmet.com",
	"mytrashmail.com",
	"mailcatch.com",
	"mailexpire.com",
	"tempinbox.com",
	"tempr.email",
	"discard.email",
	"discardmail.com",
	"spambog.com",
	"spambog.de",
	"spambog.ru",
	"emailondeck.com",
	"mohmal.com",
	"burnermail.io",
	"mintemail.com",
	"tempail.com",
	"emailfake.com",
	"crazymailing.com",
}
