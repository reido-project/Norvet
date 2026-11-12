package nvterr

import (
	"maps"
	"slices"
	"time"
)

// The With* methods never mutate the receiver; each returns a modified copy.
// They are what make it safe to share the package-level sentinels.

// WithMessage returns a copy with a different client-safe message.
func (e *Error) WithMessage(message string) *Error {
	c := *e
	c.message = message
	return &c
}

// WithCause returns a copy that wraps cause for logging and errors.Is/As.
// The cause is never exposed to clients.
func (e *Error) WithCause(cause error) *Error {
	c := *e
	c.cause = cause
	return &c
}

// WithMeta returns a copy with an extra key/value for clients. An existing
// key is overwritten.
func (e *Error) WithMeta(key string, value any) *Error {
	c := *e
	c.meta = make(map[string]any, len(e.meta)+1)
	maps.Copy(c.meta, e.meta)
	c.meta[key] = value
	return &c
}

// WithViolations returns a copy with per-field violations appended to any the
// error already carries.
func (e *Error) WithViolations(v ...Violation) *Error {
	c := *e
	c.violations = slices.Concat(e.violations, v)
	return &c
}

// WithRetryAfter returns a copy carrying a Retry-After hint. It sets only the
// hint (see RetryAfter and RetryAfterSeconds); NewRateLimited and
// NewTooManyAttempts also expose it in Meta.
func (e *Error) WithRetryAfter(d time.Duration) *Error {
	c := *e
	c.retryAfter = d
	return &c
}
