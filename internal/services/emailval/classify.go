package emailval

// Validation statuses stored on contacts and written by `outlet validate`.
const (
	StatusValid      = "valid"
	StatusRisky      = "risky"
	StatusInvalid    = "invalid"
	StatusDisposable = "disposable"
	StatusUnknown    = "unknown"
)

// Classify maps a Result to a single send decision.
// catchAll marks a domain that accepted a random address, so its RCPT answers prove nothing.
func Classify(r *Result, catchAll bool) string {
	if r.FailedAt != nil && (*r.FailedAt == LevelSyntax || *r.FailedAt == LevelDomain) {
		return StatusInvalid
	}
	if r.SMTPDeliverable != nil && !*r.SMTPDeliverable {
		return StatusInvalid
	}
	if r.Disposable {
		return StatusDisposable
	}
	if r.RoleBased || catchAll {
		return StatusRisky
	}
	if r.SMTPChecked && r.SMTPDeliverable == nil {
		return StatusUnknown
	}
	return StatusValid
}
