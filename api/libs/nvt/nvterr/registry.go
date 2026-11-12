package nvterr

// registered holds every error created by define/defineChild, in package
// initialization order. It is written only during package initialization and
// read-only afterwards, so it needs no lock. It exists for tests and for
// generating ERRORS.md.
var registered []*Error

// registry returns a copy of the registered errors.
func registry() []*Error { return append([]*Error(nil), registered...) }

// define registers a top-level error.
func define(code string, status int, message string) *Error {
	e := New(code, status, message)
	registered = append(registered, e)
	return e
}

// defineChild registers an error that specializes parent: it keeps the
// parent's HTTP status and also matches the parent in errors.Is.
func defineChild(parent *Error, code, message string) *Error {
	e := NewChild(parent, code, message)
	registered = append(registered, e)
	return e
}
