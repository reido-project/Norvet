// Package nvterr provides production-grade REST API errors for the Norvet
// framework.
//
// It is framework-agnostic (standard library only), so the same errors work
// with Fiber v3, Fiber v2, Gin or net/http. See README.md for ready-made
// error handlers.
//
// Everything lives in this one flat package, so typing "nvterr." in an editor
// lists every error and helper:
//
//	nvterr.Err...  sentinel errors, for errors.Is checks and for returning as-is
//	nvterr.New...  constructors that put a subject in the message
//
// Errors form a tree from generic to specific. A specific error matches all
// of its ancestors, never the other way round:
//
//	err := nvterr.ErrEmailAlreadyRegistered
//	errors.Is(err, nvterr.ErrEmailAlreadyRegistered)    // true
//	errors.Is(err, nvterr.ErrAlreadyRegistered)         // true
//	errors.Is(err, nvterr.ErrConflict)                  // true
//	errors.Is(err, nvterr.ErrUsernameAlreadyRegistered) // false
//
// Use the generic constructors for subjects without a dedicated sentinel:
//
//	nvterr.NewNotFound("waldo")           // 404 "waldo not found"
//	nvterr.NewAlreadyRegistered("SSN")    // 409 "SSN already registered"
//	nvterr.ErrInternal.WithCause(dbErr)   // 500, cause never leaves the server
//
// In a response handler, nvterr.From(err) turns any error into an *Error
// (unknown errors become a generic 500), so internals never leak. Errors
// that a framework produced itself and that carry only a status code can be
// mapped with nvterr.FromStatus.
//
// # Custom errors
//
// New creates a standalone error; NewChild creates one that sits under an
// existing error, inherits its status and matches it in errors.Is:
//
//	var ErrInvoiceLocked = nvterr.NewChild(nvterr.ErrLocked, "invoice_locked", "invoice is locked")
//
// # Layout
//
// Files are grouped by domain; there are no sub-packages (the error tree
// crosses domains, which would make them import each other) and no generated
// code except ERRORS.md:
//
//	request*.go       malformed requests, headers, HTTP protocol errors
//	validation*.go    field validation, passwords, pagination/sort/filter
//	auth*.go          authentication (401), authorization (403), tokens, MFA
//	account*.go       sign-up identity: username/email/phone, users
//	organization.go   organizations, members, roles, invites, tenants
//	resource*.go      not found, conflicts, state, concurrency, gone/expired
//	file.go           uploads and stored files
//	limits*.go        rate limits, quotas, plan limits, feature gating
//	billing.go        payment required (402)
//	server*.go        5xx: internal, upstream, unavailable, timeouts
//
//	error.go            the Error type, accessors, Is, New, NewChild
//	error_modifiers.go  the With* methods
//	error_json.go       Violation and the JSON representation
//	convert.go          From, StatusOf, CodeOf, IsKnown, IsServerError, IsClientError
//	status.go           FromStatus and StatusClientClosedRequest
//	registry.go         define/defineChild bookkeeping
//
// Adding an error: add a define(...) or defineChild(parent, ...) line to the
// right domain file, then run "go generate ./..." in this directory to
// refresh ERRORS.md. "go test" fails while ERRORS.md is out of date and
// checks that codes are unique and that children share their parent's status.
package nvterr
