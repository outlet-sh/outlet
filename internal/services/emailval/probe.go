package emailval

import "context"

// ProbeResult is the outcome of an SMTP RCPT probe for one address.
type ProbeResult struct {
	// Deliverable is nil when the server gave no definitive answer.
	Deliverable *bool
	// CatchAll is true when the domain accepted a random mailbox.
	CatchAll bool
	// SubStatus explains a rejection or an unknown result (Sub* constants), empty when accepted.
	SubStatus string
}

// probeMailbox probes email on domain's MX hosts (preference order).
// It owns proxy selection, per-provider rate limits, catch-all detection and response parsing.
func probeMailbox(ctx context.Context, domain, email string, mxHosts []string, opts *Options) ProbeResult {
	ok, err := checkSMTPMailbox(ctx, domain, email, opts)
	if err != nil {
		return ProbeResult{SubStatus: SubExceptionOccurred}
	}
	if !ok {
		return ProbeResult{Deliverable: &ok, SubStatus: SubMailboxNotFound}
	}
	return ProbeResult{Deliverable: &ok}
}
