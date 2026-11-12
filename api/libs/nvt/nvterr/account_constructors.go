package nvterr

// NewAlreadyRegistered returns a 409 with the message "<subject> already
// registered".
//
//	nvterr.NewAlreadyRegistered("SSN")
//
// Prefer the dedicated sentinels (ErrEmailAlreadyRegistered,
// ErrUsernameAlreadyRegistered, ErrPhoneAlreadyRegistered) where they exist.
func NewAlreadyRegistered(subject string) *Error {
	return ErrAlreadyRegistered.WithMessage(subject + " already registered")
}
