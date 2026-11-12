package nvterr

import (
	"slices"
	"strconv"
	"strings"
)

// NewMissingParameter returns a 400 with the message "missing required
// parameter '<name>'" and the name as the "parameter" meta entry.
func NewMissingParameter(name string) *Error {
	return ErrMissingParameter.
		WithMessage("missing required parameter '"+name+"'").
		WithMeta("parameter", name)
}

// NewInvalidParameter returns a 400 with the message "invalid parameter
// '<name>': <reason>" and the name as the "parameter" meta entry.
func NewInvalidParameter(name, reason string) *Error {
	return ErrInvalidParameter.
		WithMessage("invalid parameter '"+name+"': "+reason).
		WithMeta("parameter", name)
}

// NewMissingHeader returns a 400 with the message "missing required header
// '<name>'" and the name as the "header" meta entry.
func NewMissingHeader(name string) *Error {
	return ErrMissingHeader.
		WithMessage("missing required header '"+name+"'").
		WithMeta("header", name)
}

// NewMethodNotAllowed returns a 405 listing the allowed methods as the
// "allowed" meta entry. Also set the Allow header on the response. With no
// methods it returns ErrMethodNotAllowed unchanged.
func NewMethodNotAllowed(allowed ...string) *Error {
	if len(allowed) == 0 {
		return ErrMethodNotAllowed
	}
	return ErrMethodNotAllowed.WithMeta("allowed", slices.Clone(allowed))
}

// NewPayloadTooLarge returns a 413 stating the limit in bytes (also exposed
// as the "max_bytes" meta entry).
func NewPayloadTooLarge(maxBytes int64) *Error {
	return ErrPayloadTooLarge.
		WithMessage("payload exceeds the maximum size of "+strconv.FormatInt(maxBytes, 10)+" bytes").
		WithMeta("max_bytes", maxBytes)
}

// NewUnsupportedMediaType returns a 415 naming the offending type and, when
// given, the accepted ones (also exposed as the "accepted" meta entry).
//
//	nvterr.NewUnsupportedMediaType("text/plain", "application/json")
//	// "unsupported media type 'text/plain', expected one of: application/json"
func NewUnsupportedMediaType(got string, accepted ...string) *Error {
	message := "unsupported media type '" + got + "'"
	if len(accepted) == 0 {
		return ErrUnsupportedMediaType.WithMessage(message)
	}
	return ErrUnsupportedMediaType.
		WithMessage(message+", expected one of: "+strings.Join(accepted, ", ")).
		WithMeta("accepted", slices.Clone(accepted))
}
