package nvterr

import (
	"net/http"
)

// Domain "request": malformed requests and HTTP protocol-level errors.

var (
	ErrBadRequest = define("bad_request", http.StatusBadRequest, "bad request")

	// Body.
	ErrMalformedBody = define("malformed_body", http.StatusBadRequest, "request body is malformed")
	ErrInvalidJSON   = defineChild(ErrMalformedBody, "invalid_json", "request body is not valid JSON")
	ErrInvalidXML    = defineChild(ErrMalformedBody, "invalid_xml", "request body is not valid XML")
	ErrInvalidForm   = defineChild(ErrMalformedBody, "invalid_form", "request form data is malformed")
	ErrEmptyBody     = define("empty_body", http.StatusBadRequest, "request body is required")
	ErrUnknownField  = define("unknown_field", http.StatusBadRequest, "request contains an unknown field")

	// Parameters.
	ErrMissingParameter = define("missing_parameter", http.StatusBadRequest, "required parameter is missing")
	ErrInvalidParameter = define("invalid_parameter", http.StatusBadRequest, "parameter is invalid")
	ErrInvalidID        = defineChild(ErrInvalidParameter, "invalid_id", "identifier is invalid")
	ErrInvalidUUID      = defineChild(ErrInvalidID, "invalid_uuid", "identifier is not a valid UUID")

	// Headers.
	ErrMissingHeader         = define("missing_header", http.StatusBadRequest, "required header is missing")
	ErrMissingContentType    = defineChild(ErrMissingHeader, "missing_content_type", "Content-Type header is required")
	ErrIdempotencyKeyMissing = defineChild(ErrMissingHeader, "idempotency_key_missing", "Idempotency-Key header is required")
	ErrInvalidHeader         = define("invalid_header", http.StatusBadRequest, "header is invalid")

	// Idempotency key sent again with a different payload.
	ErrIdempotencyKeyReused = define("idempotency_key_reused", http.StatusUnprocessableEntity, "idempotency key was already used with a different request")

	ErrAPIVersionUnsupported = define("api_version_unsupported", http.StatusBadRequest, "API version is not supported")

	// Protocol.
	ErrMethodNotAllowed    = define("method_not_allowed", http.StatusMethodNotAllowed, "method not allowed")
	ErrNotAcceptable       = define("not_acceptable", http.StatusNotAcceptable, "no acceptable representation available")
	ErrRequestTimeout      = define("request_timeout", http.StatusRequestTimeout, "request timed out")
	ErrClientClosedRequest = define("client_closed_request", StatusClientClosedRequest, "client closed the request")
	ErrPayloadTooLarge     = define("payload_too_large", http.StatusRequestEntityTooLarge, "payload too large")

	ErrUnsupportedMediaType = define("unsupported_media_type", http.StatusUnsupportedMediaType, "unsupported media type")
)
