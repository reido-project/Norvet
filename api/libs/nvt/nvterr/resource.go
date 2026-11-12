package nvterr

import (
	"net/http"
)

// Domain "resource": lookup, conflicts, state, concurrency.

var (
	// 404
	ErrNotFound      = define("not_found", http.StatusNotFound, "resource not found")
	ErrRouteNotFound = defineChild(ErrNotFound, "route_not_found", "endpoint not found")

	// 410
	ErrGone             = define("gone", http.StatusGone, "resource is no longer available")
	ErrAPIVersionSunset = defineChild(ErrGone, "api_version_sunset", "this API version is no longer available")
	ErrExpired          = define("expired", http.StatusGone, "resource has expired")
	ErrLinkExpired      = defineChild(ErrExpired, "link_expired", "link has expired")

	// 409
	ErrConflict = define("conflict", http.StatusConflict, "request conflicts with the current state of the resource")

	ErrAlreadyExists = defineChild(ErrConflict, "already_exists", "resource already exists")
	ErrSlugTaken     = defineChild(ErrAlreadyExists, "slug_taken", "slug is already taken")
	ErrNameTaken     = defineChild(ErrAlreadyExists, "name_taken", "name is already taken")

	ErrInvalidState      = defineChild(ErrConflict, "invalid_state", "resource is in an invalid state for this operation")
	ErrInvalidTransition = defineChild(ErrInvalidState, "invalid_transition", "invalid state transition")
	ErrNotActive         = defineChild(ErrInvalidState, "not_active", "resource is not active")
	ErrAlreadyActive     = defineChild(ErrInvalidState, "already_active", "resource is already active")
	ErrAlreadyEnabled    = defineChild(ErrInvalidState, "already_enabled", "already enabled")
	ErrAlreadyDisabled   = defineChild(ErrInvalidState, "already_disabled", "already disabled")
	ErrAlreadyVerified   = defineChild(ErrInvalidState, "already_verified", "already verified")
	ErrAlreadyAccepted   = defineChild(ErrInvalidState, "already_accepted", "already accepted")
	ErrAlreadyCompleted  = defineChild(ErrInvalidState, "already_completed", "already completed")
	ErrAlreadyCanceled   = defineChild(ErrInvalidState, "already_canceled", "already canceled")
	ErrAlreadyArchived   = defineChild(ErrInvalidState, "already_archived", "already archived")
	ErrAlreadyDeleted    = defineChild(ErrInvalidState, "already_deleted", "already deleted")
	ErrAlreadyProcessed  = defineChild(ErrInvalidState, "already_processed", "already processed")

	ErrInUse                  = defineChild(ErrConflict, "in_use", "resource is in use and cannot be modified")
	ErrHasDependents          = defineChild(ErrConflict, "has_dependents", "resource has dependent resources")
	ErrConcurrentModification = defineChild(ErrConflict, "concurrent_modification", "resource was modified by another request")
	ErrVersionConflict        = defineChild(ErrConcurrentModification, "version_conflict", "resource version does not match")
	ErrRequestInProgress      = defineChild(ErrConflict, "request_in_progress", "an identical request is already in progress")

	// Conditional requests (If-Match / If-Unmodified-Since).
	ErrPreconditionFailed   = define("precondition_failed", http.StatusPreconditionFailed, "precondition failed")
	ErrETagMismatch         = defineChild(ErrPreconditionFailed, "etag_mismatch", "ETag does not match the current version")
	ErrPreconditionRequired = define("precondition_required", http.StatusPreconditionRequired, "precondition header is required")

	ErrLocked = define("locked", http.StatusLocked, "resource is locked")
)
