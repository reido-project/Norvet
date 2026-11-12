package nvterr

import (
	"strconv"
	"time"
)

// NewRateLimited returns a 429 with a Retry-After hint. The wait is exposed
// both as RetryAfter and as the "retry_after_seconds" meta entry (rounded up
// to whole seconds).
func NewRateLimited(retryAfter time.Duration) *Error {
	return withRetryAfterMeta(ErrRateLimited, retryAfter)
}

// NewTooManyAttempts returns a 429 for repeated attempts (login, OTP) with a
// Retry-After hint. See NewRateLimited.
func NewTooManyAttempts(retryAfter time.Duration) *Error {
	return withRetryAfterMeta(ErrTooManyAttempts, retryAfter)
}

// NewQuotaExceeded returns a 403 with the message "<subject> quota exceeded".
func NewQuotaExceeded(subject string) *Error {
	return ErrQuotaExceeded.WithMessage(subject + " quota exceeded")
}

// NewLimitExceeded returns a 403 with the message "<subject> limit of <limit>
// reached" and the limit as the "limit" meta entry.
//
//	nvterr.NewLimitExceeded("project", 5) // "project limit of 5 reached"
func NewLimitExceeded(subject string, limit int) *Error {
	return ErrLimitExceeded.
		WithMessage(subject+" limit of "+strconv.Itoa(limit)+" reached").
		WithMeta("limit", limit)
}

// NewFeatureDisabled returns a 403 with the message "<subject> is disabled".
func NewFeatureDisabled(subject string) *Error {
	return ErrFeatureDisabled.WithMessage(subject + " is disabled")
}

func withRetryAfterMeta(base *Error, retryAfter time.Duration) *Error {
	return base.WithRetryAfter(retryAfter).WithMeta("retry_after_seconds", ceilSeconds(retryAfter))
}
