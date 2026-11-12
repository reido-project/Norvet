# nvterr

Typed, hierarchical REST API errors for Go. Each error is a stable
machine-readable **code**, an **HTTP status**, and a **client-safe message**.
Internal details (database errors, stack traces) ride along as a hidden
**cause** and never reach the response.

- Standard library only: works the same with Fiber v3, Fiber v2, Gin and `net/http`.
- Over 260 ready-made errors, browsable by typing `nvterr.Err` in your editor. See [ERRORS.md](ERRORS.md).
- Immutable: the `Err*` sentinels are safe to share and return directly.

```go
import "api/libs/nvt/nvterr"
```

## Quick start

```go
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	u, err := s.repo.Find(ctx, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nvterr.ErrUserNotFound
	case err != nil:
		return nil, nvterr.ErrDatabaseFailed.WithCause(err) // client sees only "internal server error"
	}
	return u, nil
}
```

Return the error as-is up to your error handler (below). There is nothing else
to wire up.

## Responding to the client

`nvterr.From(err)` turns **any** error into an `*nvterr.Error`:

| Input | Result |
|---|---|
| an `*nvterr.Error` anywhere in the chain | returned as-is |
| `context.Canceled` | `client_closed_request` (499) |
| `context.DeadlineExceeded` | `gateway_timeout` (504) |
| anything else | `internal` (500), original error kept as hidden cause |

`nvterr.IsKnown(err)` reports whether `From` can classify `err` (an `*nvterr.Error` in the chain, or a context
cancellation/deadline) rather than falling back to the generic 500.

`nvterr.FromStatus(code)` maps a bare HTTP status to the matching sentinel. Use
it for errors the framework generates itself (404 for an unmatched route, 405,
413...), which would otherwise be reported as 500.

The JSON form of an error is `{"code", "message", "violations"?, "meta"?}`. How
you wrap it (`{"error": ...}`, RFC 9457, ...) is up to you.

### Fiber v3

```go
app := fiber.New(fiber.Config{ErrorHandler: errorHandler})

func errorHandler(c fiber.Ctx, err error) error {
	var (
		e  *nvterr.Error
		fe *fiber.Error
	)
	switch {
	case errors.As(err, &e) && e != nil: // already one of ours
	case errors.As(err, &fe) && fe != nil: // router 404/405, body limit, ...
		e = nvterr.FromStatus(fe.Code)
	default:
		e = nvterr.From(err)
	}

	if nvterr.IsServerError(e) {
		slog.Error("request failed", "err", e, "method", c.Method(), "path", c.Path()) // includes the hidden cause
	}
	if s := e.RetryAfterSeconds(); s > 0 {
		c.Set(fiber.HeaderRetryAfter, strconv.Itoa(s))
	}
	return c.Status(e.Status()).JSON(fiber.Map{"error": e})
}
```

### Fiber v2

Same handler; only the signature differs:

```go
func errorHandler(c *fiber.Ctx, err error) error { /* identical body */ }
```

### Gin

```go
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}
		e := nvterr.From(c.Errors.Last().Err)
		if s := e.RetryAfterSeconds(); s > 0 {
			c.Header("Retry-After", strconv.Itoa(s))
		}
		c.AbortWithStatusJSON(e.Status(), gin.H{"error": e})
	}
}

// In a handler:  _ = c.Error(nvterr.NewNotFound("invoice"))
// Unmatched routes never reach the middleware:
//   r.NoRoute(func(c *gin.Context) { c.AbortWithStatusJSON(404, gin.H{"error": nvterr.ErrRouteNotFound}) })
```

## The error tree

A specific error matches **itself and every ancestor**, never the reverse:

```go
err := nvterr.ErrEmailAlreadyRegistered
errors.Is(err, nvterr.ErrEmailAlreadyRegistered)    // true
errors.Is(err, nvterr.ErrAlreadyRegistered)         // true
errors.Is(err, nvterr.ErrConflict)                  // true
errors.Is(err, nvterr.ErrUsernameAlreadyRegistered) // false
```

Children keep their parent's HTTP status. Matching is by code, so modified
copies (`NewNotFound("x")`, `WithMessage(...)`) still match their sentinel.
Catch a whole family with one check: `errors.Is(err, nvterr.ErrNotFound)`.

## Constructors

Use a sentinel when one exists; use a constructor to put a subject in the message.

| Constructor | Status | Message |
|---|---|---|
| `NewNotFound(subject)` | 404 | `<subject> not found` |
| `NewAlreadyExists(subject)` | 409 | `<subject> already exists` |
| `NewAlreadyRegistered(subject)` | 409 | `<subject> already registered` |
| `NewConflict(reason)` | 409 | `<reason>` |
| `NewInUse(subject)` | 409 | `<subject> is in use` |
| `NewHasDependents(subject, dependents)` | 409 | `<subject> has dependent <dependents>` |
| `NewInvalidState(subject, current)` | 409 | `<subject> is <current>, which does not allow this operation` |
| `NewInvalidTransition(subject, from, to)` | 409 | `<subject> cannot go from <from> to <to>` |
| `NewGone(subject)` / `NewExpired(subject)` | 410 | `<subject> is no longer available` / `has expired` |
| `NewLocked(subject)` | 423 | `<subject> is locked` |
| `NewForbidden(action)` | 403 | `you are not allowed to <action>` |
| `NewMissingScope(scopes...)` | 403 | `token is missing required scope: ...` |
| `NewQuotaExceeded(subject)` | 403 | `<subject> quota exceeded` |
| `NewLimitExceeded(subject, limit)` | 403 | `<subject> limit of <limit> reached` |
| `NewFeatureDisabled(subject)` | 403 | `<subject> is disabled` |
| `NewMissingParameter(name)` | 400 | `missing required parameter '<name>'` |
| `NewInvalidParameter(name, reason)` | 400 | `invalid parameter '<name>': <reason>` |
| `NewMissingHeader(name)` | 400 | `missing required header '<name>'` |
| `NewMethodNotAllowed(allowed...)` | 405 | `method not allowed`, `allowed` in meta |
| `NewPayloadTooLarge(maxBytes)` | 413 | `payload exceeds the maximum size of <n> bytes` |
| `NewUnsupportedMediaType(got, accepted...)` | 415 | `unsupported media type '<got>', expected one of: ...` |
| `NewValidation(violations...)` | 422 | `validation failed` + violations |
| `NewInvalidField(field, reason)` | 422 | `validation failed` + one violation |
| `NewRequiredField(field)` | 422 | `field is required` + one violation |
| `NewNotSupported(subject)` | 422 | `<subject> is not supported` |
| `NewBusinessRule(reason)` | 422 | `<reason>` |
| `NewRateLimited(retryAfter)` | 429 | `too many requests` + Retry-After |
| `NewTooManyAttempts(retryAfter)` | 429 | `too many attempts, try again later` + Retry-After |
| `NewInternal(cause)` | 500 | `internal server error` |
| `NewDependencyFailed(service, cause)` | 502 | `<service> failed` |
| `NewDependencyUnavailable(service, cause)` | 503 | `<service> is unavailable` |
| `NewDependencyTimeout(service, cause)` | 504 | `<service> timed out` |
| `NewUnavailable(subject)` | 503 | `<subject> is temporarily unavailable` |

### Validation errors

```go
return nvterr.NewValidation(
	nvterr.Violation{Field: "email", Message: "is not a valid address", Code: "invalid_email"},
	nvterr.Violation{Field: "age", Message: "must be at least 18"},
)
```

```json
{
  "code": "validation_failed",
  "message": "validation failed",
  "violations": [
    {"field": "email", "message": "is not a valid address", "code": "invalid_email"},
    {"field": "age", "message": "must be at least 18"}
  ]
}
```

### Rate limits

```go
err := nvterr.NewRateLimited(1500 * time.Millisecond)
err.RetryAfter()        // 1.5s
err.RetryAfterSeconds() // 2: rounded up for the Retry-After header, so clients never retry early
```

## Modifying an error

Every `With*` method returns a copy; the original is never touched.

| Method | Effect |
|---|---|
| `WithMessage(s)` | replace the client-safe message |
| `WithCause(err)` | attach the internal cause (logs and `errors.Is/As` only) |
| `WithMeta(key, value)` | add extra context for clients |
| `WithViolations(v...)` | append per-field violations |
| `WithRetryAfter(d)` | set the Retry-After hint |

## Custom errors

```go
// Standalone: matches only its own code.
var ErrTeapot = nvterr.New("teapot", http.StatusTeapot, "I'm a teapot")

// Child of an existing error: inherits the status and matches the parent.
var ErrInvoiceLocked = nvterr.NewChild(nvterr.ErrLocked, "invoice_locked", "invoice is locked")
errors.Is(ErrInvoiceLocked, nvterr.ErrLocked) // true
```

Custom errors are not registered and do not appear in `ERRORS.md`.

## Security notes

- **Never put internals in a message.** Attach them with `WithCause`; the cause
  appears in `Error()` (for logs) but not in `Message()` or the JSON.
- All 500-status errors share the generic message `internal server error`. A
  test enforces it.
- On login endpoints return plain `ErrInvalidCredentials`. The more specific
  children (`ErrCurrentPasswordIncorrect`, ...) tell an attacker which part was
  wrong.

## Maintaining the package

Add a `define("code", status, "message")` or `defineChild(parent, "code", "message")`
line to the right domain file, then:

```sh
go generate ./...   # refreshes ERRORS.md
go test ./...
```

The tests fail if `ERRORS.md` is stale, a code is duplicated or not snake_case,
a child's status differs from its parent's, or a 500 leaks a specific message.
The file layout is described in [doc.go](doc.go).
