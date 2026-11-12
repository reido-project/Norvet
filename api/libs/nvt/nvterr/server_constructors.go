package nvterr

// NewInternal returns a 500 with a generic client message. The cause is kept
// for logs and never returned to clients.
func NewInternal(cause error) *Error { return ErrInternal.WithCause(cause) }

// NewDependencyFailed returns a 502 with the message "<service> failed". The
// cause is kept for logs and never returned to clients.
func NewDependencyFailed(service string, cause error) *Error {
	return ErrDependencyFailed.WithMessage(service + " failed").WithCause(cause)
}

// NewDependencyUnavailable returns a 503 with the message "<service> is
// unavailable". The cause is kept for logs and never returned to clients.
func NewDependencyUnavailable(service string, cause error) *Error {
	return ErrDependencyUnavailable.WithMessage(service + " is unavailable").WithCause(cause)
}

// NewDependencyTimeout returns a 504 with the message "<service> timed out".
// The cause is kept for logs and never returned to clients.
func NewDependencyTimeout(service string, cause error) *Error {
	return ErrGatewayTimeout.WithMessage(service + " timed out").WithCause(cause)
}

// NewUnavailable returns a 503 with the message "<subject> is temporarily
// unavailable".
func NewUnavailable(subject string) *Error {
	return ErrServiceUnavailable.WithMessage(subject + " is temporarily unavailable")
}
