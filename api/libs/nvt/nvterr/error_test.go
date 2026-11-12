package nvterr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"testing"
	"time"
)

func TestIs_Hierarchy(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{"self", ErrEmailAlreadyRegistered, ErrEmailAlreadyRegistered, true},
		{"parent", ErrEmailAlreadyRegistered, ErrAlreadyRegistered, true},
		{"grandparent", ErrEmailAlreadyRegistered, ErrAlreadyExists, true},
		{"root", ErrEmailAlreadyRegistered, ErrConflict, true},
		{"sibling", ErrEmailAlreadyRegistered, ErrUsernameAlreadyRegistered, false},
		{"never upwards", ErrConflict, ErrEmailAlreadyRegistered, false},
		{"cross-tree parent", ErrSessionNotFound, ErrNotFound, true},
		{"modified copy matches sentinel", NewNotFound("waldo"), ErrNotFound, true},
		{"wrapped by fmt.Errorf", fmt.Errorf("ctx: %w", ErrTokenExpired), ErrTokenExpired, true},
		{"cause is searched too", ErrInternal.WithCause(ErrNotFound), ErrNotFound, true},
		{"non-nvterr target", ErrNotFound, errors.New("not_found"), false},
		{"typed nil target", ErrNotFound, (*Error)(nil), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.err, tt.target); got != tt.want {
				t.Fatalf("errors.Is(%v, %v) = %v, want %v", tt.err, tt.target, got, tt.want)
			}
		})
	}
}

func TestNewChild(t *testing.T) {
	invoiceLocked := NewChild(ErrLocked, "invoice_locked", "invoice is locked")
	if invoiceLocked.Status() != http.StatusLocked {
		t.Errorf("status = %d, want inherited %d", invoiceLocked.Status(), http.StatusLocked)
	}
	if !errors.Is(invoiceLocked, ErrLocked) {
		t.Error("child must match its parent")
	}
	if errors.Is(ErrLocked, invoiceLocked) {
		t.Error("parent must not match its child")
	}
	if slices.ContainsFunc(registry(), func(e *Error) bool { return e.Code() == "invoice_locked" }) {
		t.Error("NewChild must not register the error")
	}
}

func TestImmutability(t *testing.T) {
	before := *ErrNotFound
	_ = ErrNotFound.WithMessage("x").WithCause(errors.New("c")).WithMeta("k", 1).
		WithViolations(Violation{Field: "f"}).WithRetryAfter(time.Minute)
	after := *ErrNotFound
	if before.message != after.message || before.cause != after.cause || before.meta != nil ||
		after.meta != nil || after.violations != nil || after.retryAfter != 0 {
		t.Fatalf("sentinel was mutated: %+v", after)
	}

	a := ErrNotFound.WithMeta("a", 1)
	b := a.WithMeta("b", 2)
	if _, leaked := a.Meta()["b"]; leaked {
		t.Fatal("WithMeta on a copy must not change the original")
	}
	if len(b.Meta()) != 2 {
		t.Fatalf("meta = %v, want both keys", b.Meta())
	}
}

func TestAccessorsReturnCopies(t *testing.T) {
	e := NewValidation(Violation{Field: "a", Message: "bad"}).WithMeta("k", "v")
	e.Violations()[0].Field = "mutated"
	e.Meta()["k"] = "mutated"
	if e.Violations()[0].Field != "a" || e.Meta()["k"] != "v" {
		t.Fatal("accessors must return copies")
	}
	if ErrNotFound.Meta() != nil || ErrNotFound.Violations() != nil {
		t.Fatal("empty meta and violations must be nil")
	}
}

func TestWithViolationsAppends(t *testing.T) {
	e := NewInvalidField("a", "x").WithViolations(Violation{Field: "b", Message: "y"})
	if got := len(e.Violations()); got != 2 {
		t.Fatalf("violations = %d, want 2", got)
	}
}

func TestErrorString(t *testing.T) {
	if got := ErrNotFound.Error(); got != "resource not found" {
		t.Errorf("got %q", got)
	}
	if got := ErrInternal.WithCause(errors.New("boom")).Error(); got != "internal server error: boom" {
		t.Errorf("got %q", got)
	}
	var nilErr *Error
	if got := nilErr.Error(); got != "<nil>" {
		t.Errorf("nil receiver: got %q", got)
	}
}

func TestMarshalJSON_HidesCauseAndStatus(t *testing.T) {
	e := NewInvalidField("email", "is wrong").
		WithMeta("hint", "x").
		WithCause(errors.New("secret db detail"))
	raw, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"code":"validation_failed","message":"validation failed","violations":[{"field":"email","message":"is wrong"}],"meta":{"hint":"x"}}`
	if string(raw) != want {
		t.Fatalf("got  %s\nwant %s", raw, want)
	}

	raw, _ = json.Marshal(ErrNotFound)
	if string(raw) != `{"code":"not_found","message":"resource not found"}` {
		t.Fatalf("empty fields must be omitted, got %s", raw)
	}
}

func TestFrom(t *testing.T) {
	cause := errors.New("pq: connection refused")
	var typedNil *Error

	tests := []struct {
		name       string
		err        error
		wantCode   string
		wantStatus int
		hidden     bool // err must be kept as an unexposed cause
	}{
		{"direct", ErrNotFound, "not_found", 404, false},
		{"wrapped", fmt.Errorf("load: %w", ErrNotFound), "not_found", 404, false},
		{"canceled", context.Canceled, "client_closed_request", 499, true},
		{"deadline", fmt.Errorf("q: %w", context.DeadlineExceeded), "gateway_timeout", 504, true},
		{"unknown", cause, "internal", 500, true},
		{"typed nil in interface", typedNil, "internal", 500, false},
		{"joined", errors.Join(cause, fmt.Errorf("w: %w", ErrConflict)), "conflict", 409, false},
		{"joined with context.Canceled", errors.Join(cause, context.Canceled), "client_closed_request", 499, true},
		{"joined, nothing known", errors.Join(cause, errors.New("y")), "internal", 500, true},
		{"typed nil then real", fmt.Errorf("w: %w", errors.Join(typedNil, ErrLocked)), "locked", 423, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := From(tt.err)
			if got.Code() != tt.wantCode || got.Status() != tt.wantStatus {
				t.Fatalf("got %s/%d, want %s/%d", got.Code(), got.Status(), tt.wantCode, tt.wantStatus)
			}
			if tt.hidden {
				if !errors.Is(got, tt.err) {
					t.Fatalf("cause %v not reachable from %v", tt.err, got)
				}
				if got.Message() == cause.Error() {
					t.Fatalf("cause text leaked into the client message: %q", got.Message())
				}
			}
		})
	}

	if From(nil) != nil {
		t.Error("From(nil) must be nil")
	}
	if got := From(cause); got.Message() != "internal server error" {
		t.Errorf("cause text leaked into message: %q", got.Message())
	}
	if got := From(ErrNotFound); got != ErrNotFound {
		t.Error("an *Error must be returned as-is, not copied")
	}
}

func TestStatusCodeAndClassifiers(t *testing.T) {
	tests := []struct {
		err    error
		status int
		code   string
		server bool
		client bool
	}{
		{nil, 200, "", false, false},
		{ErrNotFound, 404, "not_found", false, true},
		{ErrClientClosedRequest, 499, "client_closed_request", false, true},
		{NewDependencyFailed("mailer", errors.New("x")), 502, "dependency_failed", true, false},
		{errors.New("x"), 500, "internal", true, false},
		{context.Canceled, 499, "client_closed_request", false, true},
	}
	for _, tt := range tests {
		if got := StatusOf(tt.err); got != tt.status {
			t.Errorf("StatusOf(%v) = %d, want %d", tt.err, got, tt.status)
		}
		if got := CodeOf(tt.err); got != tt.code {
			t.Errorf("CodeOf(%v) = %q, want %q", tt.err, got, tt.code)
		}
		if got := IsServerError(tt.err); got != tt.server {
			t.Errorf("IsServerError(%v) = %v", tt.err, got)
		}
		if got := IsClientError(tt.err); got != tt.client {
			t.Errorf("IsClientError(%v) = %v", tt.err, got)
		}
	}
}

func TestClassifiersDoNotAllocate(t *testing.T) {
	plain := errors.New("x")
	wrapped := fmt.Errorf("w: %w", ErrNotFound)
	for name, fn := range map[string]func(){
		"StatusOf plain":   func() { _ = StatusOf(plain) },
		"CodeOf wrapped":   func() { _ = CodeOf(wrapped) },
		"IsServerError":    func() { _ = IsServerError(plain) },
		"StatusOf nvterr":  func() { _ = StatusOf(ErrNotFound) },
		"StatusOf ctx":     func() { _ = StatusOf(context.Canceled) },
		"errors.Is nvterr": func() { _ = errors.Is(ErrEmailAlreadyRegistered, ErrConflict) },
	} {
		if n := testing.AllocsPerRun(100, fn); n != 0 {
			t.Errorf("%s: %v allocs/op, want 0", name, n)
		}
	}
}

func TestRetryAfter(t *testing.T) {
	tests := []struct {
		in   time.Duration
		want int
	}{
		{-time.Second, 0},
		{0, 0},
		{time.Nanosecond, 1},
		{400 * time.Millisecond, 1},
		{time.Second, 1},
		{1500 * time.Millisecond, 2},
		{90 * time.Second, 90},
		{time.Duration(1<<63 - 1), 9223372037}, // no overflow at the extreme
	}
	for _, tt := range tests {
		e := ErrRateLimited.WithRetryAfter(tt.in)
		if got := e.RetryAfterSeconds(); got != tt.want {
			t.Errorf("RetryAfterSeconds(%v) = %d, want %d", tt.in, got, tt.want)
		}
		if e.RetryAfter() != tt.in {
			t.Errorf("RetryAfter(%v) altered", tt.in)
		}
	}

	e := NewRateLimited(1500 * time.Millisecond)
	if e.Meta()["retry_after_seconds"] != int64(2) || e.RetryAfterSeconds() != 2 {
		t.Errorf("NewRateLimited meta/seconds mismatch: %v", e.Meta())
	}
	e = NewTooManyAttempts(30 * time.Second)
	if e.Code() != "too_many_attempts" || e.Meta()["retry_after_seconds"] != int64(30) {
		t.Errorf("NewTooManyAttempts: %v %v", e.Code(), e.Meta())
	}
}

func TestConstructors(t *testing.T) {
	tests := []struct {
		name    string
		err     *Error
		code    string
		status  int
		message string
	}{
		{"NewNotFound", NewNotFound("waldo"), "not_found", 404, "waldo not found"},
		{"NewAlreadyExists", NewAlreadyExists("slug"), "already_exists", 409, "slug already exists"},
		{"NewAlreadyRegistered", NewAlreadyRegistered("SSN"), "already_registered", 409, "SSN already registered"},
		{"NewConflict", NewConflict("clash"), "conflict", 409, "clash"},
		{"NewInUse", NewInUse("role"), "in_use", 409, "role is in use"},
		{"NewHasDependents", NewHasDependents("project", "tasks"), "has_dependents", 409, "project has dependent tasks"},
		{"NewInvalidState", NewInvalidState("order", "shipped"), "invalid_state", 409, "order is shipped, which does not allow this operation"},
		{"NewInvalidTransition", NewInvalidTransition("order", "paid", "draft"), "invalid_transition", 409, "order cannot go from paid to draft"},
		{"NewGone", NewGone("post"), "gone", 410, "post is no longer available"},
		{"NewExpired", NewExpired("link"), "expired", 410, "link has expired"},
		{"NewLocked", NewLocked("doc"), "locked", 423, "doc is locked"},
		{"NewForbidden", NewForbidden("delete this invoice"), "forbidden", 403, "you are not allowed to delete this invoice"},
		{"NewQuotaExceeded", NewQuotaExceeded("storage"), "quota_exceeded", 403, "storage quota exceeded"},
		{"NewLimitExceeded", NewLimitExceeded("project", 5), "limit_exceeded", 403, "project limit of 5 reached"},
		{"NewFeatureDisabled", NewFeatureDisabled("exports"), "feature_disabled", 403, "exports is disabled"},
		{"NewMissingParameter", NewMissingParameter("id"), "missing_parameter", 400, "missing required parameter 'id'"},
		{"NewInvalidParameter", NewInvalidParameter("page", "must be positive"), "invalid_parameter", 400, "invalid parameter 'page': must be positive"},
		{"NewMissingHeader", NewMissingHeader("X-Tenant"), "missing_header", 400, "missing required header 'X-Tenant'"},
		{"NewPayloadTooLarge", NewPayloadTooLarge(1024), "payload_too_large", 413, "payload exceeds the maximum size of 1024 bytes"},
		{"NewUnsupportedMediaType", NewUnsupportedMediaType("text/plain"), "unsupported_media_type", 415, "unsupported media type 'text/plain'"},
		{"NewUnsupportedMediaType accepted", NewUnsupportedMediaType("text/plain", "application/json", "text/csv"), "unsupported_media_type", 415, "unsupported media type 'text/plain', expected one of: application/json, text/csv"},
		{"NewNotSupported", NewNotSupported("bulk delete"), "not_supported", 422, "bulk delete is not supported"},
		{"NewBusinessRule", NewBusinessRule("cart is empty"), "business_rule_violation", 422, "cart is empty"},
		{"NewMissingScope", NewMissingScope("a", "b"), "insufficient_scope", 403, "token is missing required scope: a, b"},
		{"NewMissingScope none", NewMissingScope(), "insufficient_scope", 403, ErrInsufficientScope.Message()},
		{"NewMethodNotAllowed none", NewMethodNotAllowed(), "method_not_allowed", 405, ErrMethodNotAllowed.Message()},
		{"NewDependencyFailed", NewDependencyFailed("mailer", errors.New("x")), "dependency_failed", 502, "mailer failed"},
		{"NewDependencyUnavailable", NewDependencyUnavailable("redis", errors.New("x")), "dependency_unavailable", 503, "redis is unavailable"},
		{"NewDependencyTimeout", NewDependencyTimeout("stripe", errors.New("x")), "gateway_timeout", 504, "stripe timed out"},
		{"NewUnavailable", NewUnavailable("search"), "service_unavailable", 503, "search is temporarily unavailable"},
		{"NewInternal", NewInternal(errors.New("x")), "internal", 500, "internal server error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code() != tt.code || tt.err.Status() != tt.status || tt.err.Message() != tt.message {
				t.Fatalf("got %s/%d/%q\nwant %s/%d/%q", tt.err.Code(), tt.err.Status(), tt.err.Message(), tt.code, tt.status, tt.message)
			}
		})
	}
}

func TestConstructorMeta(t *testing.T) {
	if m := NewInvalidTransition("o", "a", "b").Meta(); m["from"] != "a" || m["to"] != "b" {
		t.Errorf("transition meta: %v", m)
	}
	if m := NewInvalidState("o", "x").Meta(); m["current_state"] != "x" {
		t.Errorf("state meta: %v", m)
	}
	if m := NewLimitExceeded("p", 3).Meta(); m["limit"] != 3 {
		t.Errorf("limit meta: %v", m)
	}
	if m := NewPayloadTooLarge(9).Meta(); m["max_bytes"] != int64(9) {
		t.Errorf("max_bytes meta: %v", m)
	}
	v := NewRequiredField("email").Violations()
	if len(v) != 1 || v[0] != (Violation{Field: "email", Message: "is required", Code: "required"}) {
		t.Errorf("required violation: %+v", v)
	}
	if !errors.Is(NewRequiredField("x"), ErrValidationFailed) {
		t.Error("ErrRequired must match ErrValidationFailed")
	}
}

func TestConstructorsDoNotAliasCallerSlices(t *testing.T) {
	scopes := []string{"read", "write"}
	e := NewMissingScope(scopes...)
	scopes[0] = "mutated"
	if got := e.Meta()["required_scopes"].([]string); got[0] != "read" {
		t.Errorf("required_scopes aliased the caller's slice: %v", got)
	}

	methods := []string{"GET", "POST"}
	e = NewMethodNotAllowed(methods...)
	methods[0] = "mutated"
	if got := e.Meta()["allowed"].([]string); got[0] != "GET" {
		t.Errorf("allowed aliased the caller's slice: %v", got)
	}

	accepted := []string{"application/json"}
	e = NewUnsupportedMediaType("x/y", accepted...)
	accepted[0] = "mutated"
	if got := e.Meta()["accepted"].([]string); got[0] != "application/json" {
		t.Errorf("accepted aliased the caller's slice: %v", got)
	}
}

func TestFromStatus(t *testing.T) {
	tests := []struct {
		status   int
		wantCode string
	}{
		{400, "bad_request"},
		{404, "not_found"},
		{405, "method_not_allowed"},
		{413, "payload_too_large"},
		{422, "unprocessable"},
		{429, "rate_limited"},
		{499, "client_closed_request"},
		{500, "internal"},
		{502, "dependency_failed"},
		{503, "service_unavailable"},
		{504, "gateway_timeout"},
		{418, "client_error"},
		{507, "server_error"},
		{200, "internal"},
		{0, "internal"},
		{999, "internal"},
	}
	for _, tt := range tests {
		got := FromStatus(tt.status)
		if got.Code() != tt.wantCode {
			t.Errorf("FromStatus(%d).Code() = %q, want %q", tt.status, got.Code(), tt.wantCode)
		}
	}
	if e := FromStatus(418); e.Status() != 418 || e.Message() != "I'm a teapot" {
		t.Errorf("unmapped 4xx must keep status and text, got %d %q", e.Status(), e.Message())
	}
	for status, sentinel := range sentinelByStatus {
		if sentinel.Status() != status {
			t.Errorf("sentinelByStatus[%d] = %s has status %d", status, sentinel.Code(), sentinel.Status())
		}
		if sentinel.Parent() != nil {
			t.Errorf("sentinelByStatus[%d] = %s should be a root error", status, sentinel.Code())
		}
	}
}

var codePattern = regexp.MustCompile(`^[a-z][a-z0-9]*(_[a-z0-9]+)*$`)

// TestRegistryInvariants guards the hand-written definitions.
func TestRegistryInvariants(t *testing.T) {
	all := registry()
	if len(all) == 0 {
		t.Fatal("registry is empty")
	}
	seen := make(map[string]bool, len(all))
	for _, e := range all {
		if seen[e.Code()] {
			t.Errorf("duplicate code %q", e.Code())
		}
		seen[e.Code()] = true

		if !codePattern.MatchString(e.Code()) {
			t.Errorf("code %q is not snake_case", e.Code())
		}
		if e.Message() == "" {
			t.Errorf("%s: empty message", e.Code())
		}
		if e.Status() < 400 || e.Status() > 599 {
			t.Errorf("%s: status %d is not an error status", e.Code(), e.Status())
		}
		if p := e.Parent(); p != nil {
			if p.Status() != e.Status() {
				t.Errorf("%s: status %d differs from parent %s (%d)", e.Code(), e.Status(), p.Code(), p.Status())
			}
			if !seen[p.Code()] && !slices.Contains(all, p) {
				t.Errorf("%s: parent %s is not registered", e.Code(), p.Code())
			}
		}
		// 500 messages must never describe the failure to clients.
		if e.Status() == http.StatusInternalServerError && e.Message() != "internal server error" {
			t.Errorf("%s: a 500 must use the generic message, got %q", e.Code(), e.Message())
		}
	}
}

func TestIsKnown(t *testing.T) {
	var typedNil *Error
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"sentinel", ErrNotFound, true},
		{"wrapped", fmt.Errorf("w: %w", ErrNotFound), true},
		{"canceled", context.Canceled, true},
		{"deadline", fmt.Errorf("w: %w", context.DeadlineExceeded), true},
		{"unknown", errors.New("x"), false},
		{"typed nil", typedNil, false},
		{"joined known", errors.Join(errors.New("x"), ErrLocked), true},
	}
	for _, tt := range tests {
		if got := IsKnown(tt.err); got != tt.want {
			t.Errorf("%s: IsKnown = %v, want %v", tt.name, got, tt.want)
		}
	}
	if n := testing.AllocsPerRun(100, func() { _ = IsKnown(errors.New("x")) }); n > 1 {
		t.Errorf("IsKnown allocates %v times beyond the error itself", n)
	}
}
