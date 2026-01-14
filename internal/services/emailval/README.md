# emailval - Email Validation Package

A layered email validation library for Go with syntax, DNS, disposable domain, role-based, and optional SMTP verification.

## Quick Start

```go
import (
    "context"
    "time"

    "outlet/internal/services/emailval"
)

// Basic validation (syntax + MX only - RECOMMENDED for most use cases)
ctx := context.Background()
result, err := emailval.Validate(ctx, "user@example.com", nil)

if result.Valid() {
    // Email passed all checks
}
```

## Validation Levels

The package performs validation in layers, stopping at the first failure:

| Level          | What it checks                                         | Network?   | Risk     |
| -------------- | ------------------------------------------------------ | ---------- | -------- |
| **Syntax**     | RFC-like pattern, no consecutive dots, valid structure | No         | None     |
| **Domain**     | MX records exist for the domain                        | Yes (DNS)  | None     |
| **Disposable** | Domain is not a known throwaway service                | No         | None     |
| **Role**       | Address is not role-based (info@, support@, etc.)      | No         | None     |
| **SMTP**       | Mailbox actually exists on the server                  | Yes (SMTP) | **HIGH** |

## Recommended Configuration

### For Lead Forms / Signups (RECOMMENDED)

```go
result, err := emailval.Validate(ctx, email, &emailval.Options{
    Timeout:         5 * time.Second,
    CheckDisposable: true,  // Block throwaway emails
    AllowRole:       false, // Block info@, support@, etc.
    CheckSMTP:       false, // DO NOT enable unless you understand the risks
})

if !result.Valid() {
    // Reject or flag the email
    log.Printf("Invalid email: %v", result.Messages)
}
```

### For High-Value Leads Only

```go
result, err := emailval.Validate(ctx, email, &emailval.Options{
    Timeout:         10 * time.Second,
    CheckDisposable: true,
    AllowRole:       false,
    CheckSMTP:       true, // Enable with caution
})
```

### For Newsletter Signups (Lenient)

```go
result, err := emailval.Validate(ctx, email, &emailval.Options{
    Timeout:         3 * time.Second,
    CheckDisposable: false, // Allow disposable for low-value signups
    AllowRole:       true,  // Allow role-based
    CheckSMTP:       false,
})
```

## Understanding Results

```go
result, _ := emailval.Validate(ctx, email, opts)

// Check overall validity
if result.Valid() {
    // All enabled checks passed
}

// Inspect specific checks
result.SyntaxOK    // Passed syntax validation
result.DomainOK    // Domain has MX records
result.HasMX       // Synonym for DomainOK
result.Disposable  // Domain is disposable (if checked)
result.RoleBased   // Address is role-based (if checked)

// SMTP results (only populated if CheckSMTP: true)
result.SMTPChecked      // Whether SMTP check was attempted
result.SMTPDeliverable  // *bool: nil=unknown, true/false=verified

// Failure information
result.FailedAt   // *Level: which check failed first (nil if valid)
result.Messages   // []string: human-readable failure reasons

// Normalized email
result.Normalized // Lowercased, trimmed version
```

## SMTP Verification: Warnings

> **⚠️ You use Amazon SES for sending.** SMTP verification probes come from your server IP, not SES. If your server IP gets blacklisted, it won't affect SES sending—but it will break SMTP verification entirely. More importantly, enabling SMTP checks adds latency and complexity with minimal benefit when you already have SES bounce handling.

**SMTP verification can get your IP blacklisted.** Only use it if you:

1. Understand the risks
2. Have a dedicated IP for verification
3. Are checking low volumes (< 100/day)
4. Have implemented additional caching

### Why SMTP is Risky

- Mail servers track probe attempts
- High volume = automatic blacklisting
- Some servers (Gmail, Outlook) are especially aggressive
- Greylisting causes false negatives
- Many servers accept all addresses at SMTP level (catch-all)

### Built-in Protections

The package includes some safeguards:

- **Rate limiting**: 2 second minimum between probes to the same domain
- **Greylist detection**: Treats temporary rejections as "likely valid"
- **Max retries**: Only tries 2 MX servers before giving up
- **Proper EHLO**: Uses modern SMTP etiquette

### These Are NOT Enough For Production

For production SMTP verification, you should:

1. **Cache results** - Don't probe the same email twice
2. **Use a dedicated IP** - Separate from your sending infrastructure
3. **Configure proper identity** - Set real HELO hostname and FROM domain
4. **Consider a SaaS** - ZeroBounce, NeverBounce, Hunter.io, etc.

## Custom Disposable Domain Provider

The default list covers ~80 common disposable domains. For comprehensive coverage:

```go
// Load from file, database, or external API
domains := []string{"tempmail.com", "throwaway.email", ...}
provider := emailval.NewStaticDisposableProvider(domains)

result, err := emailval.Validate(ctx, email, &emailval.Options{
    CheckDisposable:    true,
    DisposableProvider: provider,
})
```

### Recommended External Lists

- https://github.com/disposable-email-domains/disposable-email-domains
- https://github.com/ivolo/disposable-email-domains

## Role-Based Detection

The package detects common role prefixes:

- **Contact**: info, contact, hello, enquiries
- **Support**: support, help, helpdesk, customerservice
- **Admin**: admin, webmaster, postmaster, root
- **Sales**: sales, marketing, press, partnerships
- **Team**: team, staff, office, hr, jobs
- **Technical**: noreply, abuse, security, privacy
- **Financial**: billing, accounts, finance, payments

```go
// Check manually without full validation
if emailval.IsRoleBased("info") {
    // This is a role-based local part
}
```

## Error Handling

The `Validate` function always returns a `*Result`, even on errors:

```go
result, err := emailval.Validate(ctx, email, opts)

// err is only set for unexpected errors (context cancelled, etc.)
// Validation failures are in result.FailedAt and result.Messages

if err != nil {
    // Context cancelled, timeout, etc.
    log.Printf("Validation error: %v", err)
}

if !result.Valid() {
    // Email failed validation
    log.Printf("Failed at level %d: %v", *result.FailedAt, result.Messages)
}
```

## Testing

For testing, you can inject a mock SMTP dialer:

```go
type mockDialer struct{}

func (m mockDialer) DialSMTP(ctx context.Context, addr string) (*textproto.Conn, error) {
    // Return mock connection or error
}

result, err := emailval.Validate(ctx, email, &emailval.Options{
    CheckSMTP:  true,
    SMTPDialer: mockDialer{},
})
```

## Performance Considerations

| Check      | Latency  | Concurrency Safe   |
| ---------- | -------- | ------------------ |
| Syntax     | < 1ms    | Yes                |
| DNS/MX     | 10-100ms | Yes                |
| Disposable | < 1ms    | Yes                |
| Role       | < 1ms    | Yes                |
| SMTP       | 1-10s    | Yes (rate limited) |

### Recommended Patterns

```go
// 1. Fail fast - check syntax before network calls
if !simpleEmailRegex.MatchString(email) {
    return errors.New("invalid email format")
}

// 2. Use short timeouts for user-facing validation
ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
defer cancel()

// 3. Cache results for repeat validations
// (implement your own cache layer)
```

## File Structure

```
emailval/
├── emailval.go     # Main Validate function, types, options
├── syntax.go       # Syntax validation, normalization
├── dns.go          # MX record lookup
├── disposable.go   # Disposable domain detection
├── role.go         # Role-based email detection
├── smtp.go         # SMTP mailbox verification
└── README.md       # This file
```

## Summary: What to Enable

| Use Case          | Disposable | Role | SMTP    |
| ----------------- | ---------- | ---- | ------- |
| Lead generation   | Yes        | No   | No      |
| Newsletter signup | No         | Yes  | No      |
| Account creation  | Yes        | No   | No      |
| Enterprise sales  | Yes        | No   | Maybe\* |
| Contact form      | No         | Yes  | No      |

\*Only with caching and dedicated infrastructure

## Amazon SES Integration

Since you use Amazon SES for sending, here's the recommended approach:

### Let SES Handle Bounces

SES provides robust bounce handling. Instead of SMTP verification upfront, leverage SES feedback:

1. **Configure SNS notifications** for bounces and complaints
2. **Track bounce rates** per email domain
3. **Auto-suppress** addresses that bounce

### Recommended Validation Strategy with SES

```go
// Pre-send validation (before queuing for SES)
result, err := emailval.Validate(ctx, email, &emailval.Options{
    Timeout:         3 * time.Second,
    CheckDisposable: true,  // Block throwaway before wasting SES quota
    AllowRole:       false, // Role-based rarely converts
    CheckSMTP:       false, // Let SES handle deliverability
})

if !result.Valid() {
    // Don't waste SES sends on obviously bad emails
    return fmt.Errorf("invalid email: %v", result.Messages)
}

// Queue for SES...
```

### Why This Works

| Layer       | What It Catches     | Cost                 |
| ----------- | ------------------- | -------------------- |
| Syntax + MX | Typos, fake domains | Free, instant        |
| Disposable  | Throwaway signups   | Free, instant        |
| Role-based  | Low-value leads     | Free, instant        |
| SES bounce  | Invalid mailboxes   | Counts against quota |

You're filtering out ~90% of bad emails before they ever hit SES, and letting SES's infrastructure handle the edge cases it's designed for.

### Handling SES Bounces

When SES reports a bounce, update your suppression list:

```go
// In your SNS bounce handler
func handleSESBounce(bounceNotification BounceNotification) {
    for _, recipient := range bounceNotification.Bounce.BouncedRecipients {
        // Add to your suppression list
        db.SuppressEmail(recipient.EmailAddress, "ses_bounce")
    }
}
```

Then check suppression before sending:

```go
if db.IsEmailSuppressed(email) {
    return errors.New("email is suppressed due to previous bounce")
}
```

## Common Mistakes

1. **Enabling SMTP for all signups** - This will get you blacklisted
2. **Not caching SMTP results** - Same email = same probe = wasted risk
3. **Blocking all role-based** - Sometimes info@ is the right contact
4. **Trusting SMTP too much** - Catch-all servers accept everything
5. **Using localhost as HELO** - Many servers reject this
