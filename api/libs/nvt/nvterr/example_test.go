package nvterr_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"api/libs/nvt/nvterr"
)

func Example_hierarchy() {
	err := nvterr.ErrEmailAlreadyRegistered

	fmt.Println(errors.Is(err, nvterr.ErrEmailAlreadyRegistered))
	fmt.Println(errors.Is(err, nvterr.ErrAlreadyRegistered))
	fmt.Println(errors.Is(err, nvterr.ErrConflict))
	fmt.Println(errors.Is(err, nvterr.ErrUsernameAlreadyRegistered))
	// Output:
	// true
	// true
	// true
	// false
}

func ExampleNewNotFound() {
	err := nvterr.NewNotFound("waldo")

	fmt.Println(err.Status(), err.Code(), err.Message())
	fmt.Println(errors.Is(err, nvterr.ErrNotFound))
	// Output:
	// 404 not_found waldo not found
	// true
}

func ExampleNewValidation() {
	err := nvterr.NewValidation(
		nvterr.Violation{Field: "email", Message: "is not a valid address", Code: "invalid_email"},
		nvterr.Violation{Field: "age", Message: "must be at least 18"},
	)

	body, _ := json.MarshalIndent(err, "", "  ")
	fmt.Println(err.Status())
	fmt.Println(string(body))
	// Output:
	// 422
	// {
	//   "code": "validation_failed",
	//   "message": "validation failed",
	//   "violations": [
	//     {
	//       "field": "email",
	//       "message": "is not a valid address",
	//       "code": "invalid_email"
	//     },
	//     {
	//       "field": "age",
	//       "message": "must be at least 18"
	//     }
	//   ]
	// }
}

func ExampleFrom() {
	dbErr := errors.New(`pq: password authentication failed for user "app"`)

	e := nvterr.From(dbErr)
	fmt.Println(e.Status(), e.Message()) // what the client sees
	fmt.Println(e.Error())               // what the log sees

	fmt.Println(nvterr.From(context.DeadlineExceeded).Code())
	// Output:
	// 500 internal server error
	// internal server error: pq: password authentication failed for user "app"
	// gateway_timeout
}

func ExampleError_WithCause() {
	var dbErr = errors.New("connection reset")

	err := nvterr.ErrDatabaseFailed.WithCause(dbErr)

	fmt.Println(err.Message())
	fmt.Println(errors.Is(err, dbErr))
	// Output:
	// internal server error
	// true
}

func ExampleNewRateLimited() {
	err := nvterr.NewRateLimited(1500 * time.Millisecond)

	// Value for the Retry-After header: rounded up, never early.
	fmt.Println(err.RetryAfterSeconds())
	// Output: 2
}

func ExampleNewChild() {
	var ErrInvoiceLocked = nvterr.NewChild(nvterr.ErrLocked, "invoice_locked", "invoice is locked")

	fmt.Println(ErrInvoiceLocked.Status())
	fmt.Println(errors.Is(ErrInvoiceLocked, nvterr.ErrLocked))
	// Output:
	// 423
	// true
}

func ExampleFromStatus() {
	// For errors a framework generated itself, e.g. Fiber's 404 for an
	// unmatched route, which carries only a status code.
	e := nvterr.FromStatus(404)

	fmt.Println(e.Code(), errors.Is(e, nvterr.ErrNotFound))
	// Output: not_found true
}
