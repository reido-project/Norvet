package nvterr

// NewNotFound returns a 404 with the message "<subject> not found".
//
//	nvterr.NewNotFound("waldo")
func NewNotFound(subject string) *Error { return ErrNotFound.WithMessage(subject + " not found") }

// NewAlreadyExists returns a 409 with the message "<subject> already exists".
func NewAlreadyExists(subject string) *Error {
	return ErrAlreadyExists.WithMessage(subject + " already exists")
}

// NewConflict returns a 409 with a custom reason as the message.
func NewConflict(reason string) *Error { return ErrConflict.WithMessage(reason) }

// NewInUse returns a 409 with the message "<subject> is in use".
func NewInUse(subject string) *Error { return ErrInUse.WithMessage(subject + " is in use") }

// NewHasDependents returns a 409 with the message "<subject> has dependent
// <dependents>".
//
//	nvterr.NewHasDependents("project", "tasks") // "project has dependent tasks"
func NewHasDependents(subject, dependents string) *Error {
	return ErrHasDependents.WithMessage(subject + " has dependent " + dependents)
}

// NewInvalidState returns a 409 with the message "<subject> is <current>,
// which does not allow this operation" and the state as the "current_state"
// meta entry.
func NewInvalidState(subject, current string) *Error {
	return ErrInvalidState.
		WithMessage(subject+" is "+current+", which does not allow this operation").
		WithMeta("current_state", current)
}

// NewInvalidTransition returns a 409 with the message "<subject> cannot go
// from <from> to <to>" and both states as the "from" and "to" meta entries.
func NewInvalidTransition(subject, from, to string) *Error {
	return ErrInvalidTransition.
		WithMessage(subject+" cannot go from "+from+" to "+to).
		WithMeta("from", from).
		WithMeta("to", to)
}

// NewGone returns a 410 with the message "<subject> is no longer available".
func NewGone(subject string) *Error {
	return ErrGone.WithMessage(subject + " is no longer available")
}

// NewExpired returns a 410 with the message "<subject> has expired".
func NewExpired(subject string) *Error { return ErrExpired.WithMessage(subject + " has expired") }

// NewLocked returns a 423 with the message "<subject> is locked".
func NewLocked(subject string) *Error { return ErrLocked.WithMessage(subject + " is locked") }
