package nvterr

import (
	"slices"
	"strings"
)

// NewForbidden returns a 403 with the message "you are not allowed to <action>".
//
//	nvterr.NewForbidden("delete this invoice")
func NewForbidden(action string) *Error {
	return ErrForbidden.WithMessage("you are not allowed to " + action)
}

// NewMissingScope returns a 403 naming the scope(s) the token lacks. The
// scopes are also exposed as the "required_scopes" meta entry. With no
// scopes it returns ErrInsufficientScope unchanged.
func NewMissingScope(scopes ...string) *Error {
	if len(scopes) == 0 {
		return ErrInsufficientScope
	}
	return ErrInsufficientScope.
		WithMessage("token is missing required scope: "+strings.Join(scopes, ", ")).
		WithMeta("required_scopes", slices.Clone(scopes))
}
