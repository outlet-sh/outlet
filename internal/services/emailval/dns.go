package emailval

import (
	"context"
	"net"
)

// checkMX looks up MX records for a domain.
// Returns true if at least one MX record exists.
func checkMX(ctx context.Context, domain string) (bool, error) {
	// Create a resolver that respects context cancellation.
	resolver := &net.Resolver{}

	mxRecords, err := resolver.LookupMX(ctx, domain)
	if err != nil {
		// Check if it's a "no such host" error vs network error.
		if dnsErr, ok := err.(*net.DNSError); ok && dnsErr.IsNotFound {
			return false, nil
		}
		return false, err
	}

	return len(mxRecords) > 0, nil
}

// getMXHosts returns MX hosts for a domain, sorted by preference.
func getMXHosts(ctx context.Context, domain string) ([]string, error) {
	resolver := &net.Resolver{}

	mxRecords, err := resolver.LookupMX(ctx, domain)
	if err != nil {
		return nil, err
	}

	hosts := make([]string, 0, len(mxRecords))
	for _, mx := range mxRecords {
		hosts = append(hosts, mx.Host)
	}

	return hosts, nil
}
