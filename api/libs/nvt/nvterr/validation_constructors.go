package nvterr

// NewValidation returns a 422 carrying one or more per-field violations.
func NewValidation(violations ...Violation) *Error {
	return ErrValidationFailed.WithViolations(violations...)
}

// NewInvalidField returns a 422 validation error for a single field.
//
//	nvterr.NewInvalidField("age", "must be at least 18")
func NewInvalidField(field, reason string) *Error {
	return NewValidation(Violation{Field: field, Message: reason})
}

// NewRequiredField returns a 422 validation error for a missing field. The
// violation carries the code "required".
func NewRequiredField(field string) *Error {
	return ErrRequired.WithViolations(Violation{Field: field, Message: "is required", Code: "required"})
}

// NewNotSupported returns a 422 with the message "<subject> is not supported".
func NewNotSupported(subject string) *Error {
	return ErrNotSupported.WithMessage(subject + " is not supported")
}

// NewBusinessRule returns a 422 with a custom explanation of the violated
// rule as the message.
func NewBusinessRule(reason string) *Error { return ErrBusinessRuleViolation.WithMessage(reason) }
