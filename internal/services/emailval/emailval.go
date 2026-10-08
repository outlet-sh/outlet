// Package emailval provides layered email validation utilities.
// It performs syntax, domain/MX, disposable, role-based, and optional SMTP checks.
package emailval

import (
	"context"
	"fmt"
	"time"
)

// Level represents a validation layer.
type Level int

const (
	LevelSyntax Level = iota + 1
	LevelDomain
	LevelDisposable
	LevelRole
	LevelSMTP
)

// Options configures email validation behavior.
type Options struct {
	// Timeout is the maximum duration for all network work (DNS, SMTP).
	// Defaults to 3 seconds if not set.
	Timeout time.Duration

	// CheckSMTP enables SMTP mailbox verification.
	CheckSMTP bool

	// CheckDisposable enables disposable/temporary domain detection.
	CheckDisposable bool

	// AllowRole allows role-based emails like info@, support@, admin@.
	AllowRole bool

	// DisposableProvider is an optional custom disposable domain provider.
	DisposableProvider DisposableProvider

	// SMTPDialer is an optional custom SMTP dialer (for testing or proxying).
	SMTPDialer SMTPDialer

	// SMTP configures the HELO hostname and MAIL FROM domain for SMTP checks.
	// Defaults to DefaultSMTPConfig() (localhost), which most MX servers reject.
	SMTP *SMTPConfig

	// Proxies are SOCKS5 URLs (socks5://user:pass@host:1080) of dedicated probe hosts.
	// Probes rotate across healthy proxies; empty means dial directly from this host.
	Proxies []string

	// CheckDomainAge looks up domain registration age over RDAP.
	CheckDomainAge bool
}

// Result is a structured report of validation.
type Result struct {
	Input      string
	Normalized string // lowercased, trimmed

	SyntaxOK   bool
	DomainOK   bool
	HasMX      bool
	Disposable bool
	RoleBased  bool

	SMTPChecked     bool
	SMTPDeliverable *bool // nil = unknown, true/false if SMTP check done

	// FailedAt is the first failing level, if any.
	FailedAt *Level

	// Messages contains human readable reasons for failures or warnings.
	Messages []string

	// ZeroBounce-compatible verdict (see status.go).
	Status        string `json:"status"`
	SubStatus     string `json:"sub_status"`
	FreeEmail     bool   `json:"free_email"`
	DidYouMean    string `json:"did_you_mean"`
	Account       string `json:"account"`
	Domain        string `json:"domain"`
	DomainAgeDays int    `json:"domain_age_days"` // 0 = unknown
	SMTPProvider  string `json:"smtp_provider"`   // google, microsoft, yahoo, proofpoint, mimecast, ... or ""
	MXRecord      string `json:"mx_record"`       // preferred MX host
	CatchAll      bool   `json:"catch_all"`
}

func (r *Result) addMessage(msg string, args ...any) {
	r.Messages = append(r.Messages, fmt.Sprintf(msg, args...))
}

// Valid returns true if no validation level failed.
func (r *Result) Valid() bool {
	return r.FailedAt == nil
}

// Validate performs layered validation on an email.
// It never sends a message, only checks for deliverability.
// It never returns nil Result; on error you still get a partially filled Result.
func Validate(ctx context.Context, email string, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}

	res := &Result{
		Input:      email,
		Normalized: normalize(email),
	}
	defer setVerdict(res)

	if res.Normalized == "" {
		lvl := LevelSyntax
		res.FailedAt = &lvl
		res.addMessage("empty email")
		return res, nil
	}

	// Attach timeout to context for network work.
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 1) Syntax check
	local, domain, err := validateSyntax(res.Normalized)
	if err != nil {
		lvl := LevelSyntax
		res.FailedAt = &lvl
		res.addMessage("syntax error: %v", err)
		return res, nil
	}
	res.SyntaxOK = true

	// 2) Domain MX check
	hasMX, err := checkMX(ctx, domain)
	if err != nil || !hasMX {
		lvl := LevelDomain
		res.FailedAt = &lvl
		res.addMessage("no MX records for domain %q", domain)
		return res, nil
	}
	res.DomainOK = true
	res.HasMX = true

	// 3) Disposable detection
	if opts.CheckDisposable {
		provider := opts.DisposableProvider
		if provider == nil {
			provider = DefaultDisposableProvider()
		}
		if provider.IsDisposable(domain) {
			res.Disposable = true
			lvl := LevelDisposable
			res.FailedAt = &lvl
			res.addMessage("disposable or temporary domain")
		}
	}

	// 4) Role-based check
	if !opts.AllowRole && IsRoleBased(local) {
		res.RoleBased = true
		lvl := LevelRole
		if res.FailedAt == nil {
			res.FailedAt = &lvl
		}
		res.addMessage("role-based address: %s", local)
	}

	// 5) Optional SMTP check
	if opts.CheckSMTP {
		deliverable, err := checkSMTPMailbox(ctx, domain, res.Normalized, opts)
		res.SMTPChecked = true
		if err != nil {
			// Treat as unknown, do not hard-fail on transient SMTP errors.
			res.SMTPDeliverable = nil
			res.addMessage("SMTP check error: %v", err)
		} else {
			res.SMTPDeliverable = &deliverable
			if !deliverable {
				lvl := LevelSMTP
				if res.FailedAt == nil {
					res.FailedAt = &lvl
				}
				res.addMessage("SMTP server reports mailbox undeliverable")
			}
		}
	}

	return res, nil
}

// setVerdict fills Status and SubStatus from the layer results.
func setVerdict(r *Result) {
	switch {
	case r.FailedAt != nil && *r.FailedAt == LevelSyntax:
		r.Status, r.SubStatus = StatusInvalid, SubFailedSyntaxCheck
	case r.FailedAt != nil && *r.FailedAt == LevelDomain:
		r.Status, r.SubStatus = StatusInvalid, SubNoDNSEntries
	case r.SMTPDeliverable != nil && !*r.SMTPDeliverable:
		r.Status, r.SubStatus = StatusInvalid, SubMailboxNotFound
	case r.Disposable:
		r.Status, r.SubStatus = StatusDoNotMail, SubDisposable
	case r.RoleBased:
		r.Status, r.SubStatus = StatusDoNotMail, SubRoleBased
	case r.SMTPChecked && r.SMTPDeliverable == nil:
		r.Status, r.SubStatus = StatusUnknown, SubExceptionOccurred
	default:
		r.Status = StatusValid
	}
}
