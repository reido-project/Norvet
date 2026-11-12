package nvterr

import (
	"errors"
	"maps"
	"time"
)

// Error is a REST API error: a stable machine-readable code, an HTTP status,
// and a message that is safe to show to API clients.
//
// The wrapped cause (WithCause) is for logs only; it is never part of the
// JSON representation or Message().
//
// Errors are immutable. Every With* method returns a modified copy, so the
// package-level Err* sentinels can be shared freely between goroutines.
//
// Errors form a tree. A specific error matches itself and every ancestor, so
// errors.Is(ErrEmailAlreadyRegistered, ErrAlreadyRegistered) is true, and so is
// errors.Is(..., ErrAlreadyExists) and errors.Is(..., ErrConflict). The
// reverse is never true. Matching is by Code, so a modified copy of a
// sentinel (WithMessage, WithMeta, ...) still matches the sentinel it came
// from.
type Error struct {
	code       string
	status     int
	message    string
	violations []Violation
	meta       map[string]any
	retryAfter time.Duration
	cause      error
	parent     *Error
}

// New creates a custom, standalone error for application-specific cases. It
// has no parent, so it matches only errors with the same code. Use NewChild to
// make it part of an existing branch of the tree.
//
// New does not register the error: it will not appear in ERRORS.md.
func New(code string, status int, message string) *Error {
	return &Error{code: code, status: status, message: message}
}

// NewChild creates a custom error that specializes parent: it inherits the
// parent's HTTP status and matches the parent (and the parent's ancestors) in
// errors.Is. parent must not be nil.
//
//	var ErrInvoiceLocked = nvterr.NewChild(nvterr.ErrLocked, "invoice_locked", "invoice is locked")
//	errors.Is(ErrInvoiceLocked, nvterr.ErrLocked) // true
func NewChild(parent *Error, code, message string) *Error {
	e := New(code, parent.status, message)
	e.parent = parent
	return e
}

// Code returns the stable machine-readable code, e.g. "not_found".
func (e *Error) Code() string { return e.code }

// Status returns the HTTP status code.
func (e *Error) Status() int { return e.status }

// Message returns the client-safe message.
func (e *Error) Message() string { return e.message }

// Parent returns the broader error this one specializes, or nil.
func (e *Error) Parent() *Error { return e.parent }

// Violations returns a copy of the per-field validation problems, if any.
func (e *Error) Violations() []Violation {
	if len(e.violations) == 0 {
		return nil
	}
	return append([]Violation(nil), e.violations...)
}

// Meta returns a copy of the extra structured context for clients, or nil.
func (e *Error) Meta() map[string]any {
	if len(e.meta) == 0 {
		return nil
	}
	return maps.Clone(e.meta)
}

// RetryAfter returns how long the client should wait before retrying, or 0.
// Prefer RetryAfterSeconds when setting the Retry-After header.
func (e *Error) RetryAfter() time.Duration { return e.retryAfter }

// RetryAfterSeconds returns RetryAfter rounded up to whole seconds, or 0 when
// there is no hint. It is the value to put in the Retry-After header: rounding
// up guarantees a client that obeys it never retries too early.
func (e *Error) RetryAfterSeconds() int { return int(ceilSeconds(e.retryAfter)) }

// Error implements error. It includes the cause, so it is suitable for logs
// but not for clients; use Message for that.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

// Unwrap exposes the underlying cause to errors.Is / errors.As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// Is reports whether target is an *Error with the same code as e or any of
// e's ancestors. It is called by errors.Is.
func (e *Error) Is(target error) bool {
	var t *Error
	ok := errors.As(target, &t)
	if !ok || t == nil {
		return false
	}
	for c := e; c != nil; c = c.parent {
		if c.code == t.code {
			return true
		}
	}
	return false
}

// ceilSeconds converts d to whole seconds, rounding up. Non-positive
// durations yield 0.
func ceilSeconds(d time.Duration) int64 {
	if d <= 0 {
		return 0
	}
	s := int64(d / time.Second)
	if d%time.Second != 0 {
		s++
	}
	return s
}
