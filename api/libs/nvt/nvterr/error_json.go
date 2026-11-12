package nvterr

import "encoding/json"

// Violation describes a single invalid field in a request.
type Violation struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	// Code is an optional machine-readable reason, e.g. "too_short".
	Code string `json:"code,omitempty"`
}

// MarshalJSON renders the client-facing representation:
//
//	{"code": "...", "message": "...", "violations": [...], "meta": {...}}
//
// violations and meta are omitted when empty. Status and cause are
// intentionally left out: the former belongs on the HTTP response, the latter
// must not leak.
func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code       string         `json:"code"`
		Message    string         `json:"message"`
		Violations []Violation    `json:"violations,omitempty"`
		Meta       map[string]any `json:"meta,omitempty"`
	}{e.code, e.message, e.violations, e.meta})
}
