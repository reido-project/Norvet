package nvterr

import (
	"context"
	"errors"
	"net/http"
)

// From converts any error into an *Error suitable for responding with.
//
//   - an *Error anywhere in the chain is returned as-is
//   - context.Canceled maps to ErrClientClosedRequest
//   - context.DeadlineExceeded maps to ErrGatewayTimeout
//   - anything else becomes ErrInternal with err as its (hidden) cause, so
//     internal details never reach clients
//
// From(nil) returns nil.
func From(err error) *Error {
	if err == nil {
		return nil
	}
	if e := find(err); e != nil {
		return e
	}
	return fallback(err).WithCause(err)
}

// StatusOf returns the HTTP status for err: 200 for nil, 500 for errors that
// are not nvterr errors.
func StatusOf(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return classify(err).status
}

// CodeOf returns the machine-readable code for err: "" for nil, "internal" for
// errors that are not nvterr errors.
func CodeOf(err error) string {
	if err == nil {
		return ""
	}
	return classify(err).code
}

// IsServerError reports whether err maps to a 5xx status. Useful to decide
// what to log at error level and what to alert on.
func IsServerError(err error) bool { return StatusOf(err) >= 500 }

// IsClientError reports whether err maps to a 4xx status.
func IsClientError(err error) bool {
	s := StatusOf(err)
	return s >= 400 && s < 500
}

// IsKnown reports whether From can classify err without falling back to a
// generic internal error: err is, or wraps, an *Error, context.Canceled or
// context.DeadlineExceeded. IsKnown(nil) is false.
//
// Use it to tell "an error this package understands" from "something
// unexpected" without caring which sentinel it is.
func IsKnown(err error) bool {
	if err == nil {
		return false
	}
	return find(err) != nil || fallback(err) != ErrInternal
}

// find returns the first non-nil *Error in err's chain (depth-first, like
// errors.As), or nil. It walks the chain by hand because errors.As allocates
// on every call, and this sits on the path of every error response.
func find(err error) *Error {
	for err != nil {
		if e, ok := err.(*Error); ok && e != nil {
			return e
		}
		switch w := err.(type) {
		case interface{ Unwrap() error }:
			err = w.Unwrap()
		case interface{ Unwrap() []error }:
			for _, inner := range w.Unwrap() {
				if e := find(inner); e != nil {
					return e
				}
			}
			return nil
		default:
			return nil
		}
	}
	return nil
}

// fallback picks the sentinel for an error that carries no *Error.
func fallback(err error) *Error {
	switch {
	case errors.Is(err, context.Canceled):
		return ErrClientClosedRequest
	case errors.Is(err, context.DeadlineExceeded):
		return ErrGatewayTimeout
	}
	return ErrInternal
}

// classify resolves err to an *Error without allocating: unlike From it does
// not attach the cause, so it is meant for reading Status or Code only.
func classify(err error) *Error {
	if e := find(err); e != nil {
		return e
	}
	return fallback(err)
}
