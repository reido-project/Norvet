package nvtval

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/go-sanitize/sanitize"
)

var (
	validate  = validator.New()
	sanitizer = mustNewSanitizer()
)

// FieldError describes one failed validation rule. It is meant to be sent to
// the client as JSON.
type FieldError struct {
	// Field is the name of the struct field that failed.
	Field string `json:"field"`

	// Message says in plain language what is wrong with it.
	Message string `json:"message"`
}

func mustNewSanitizer() *sanitize.Sanitizer {
	s, err := sanitize.New()
	if err != nil {
		panic(fmt.Sprintf("nvtval: cannot create sanitizer: %v", err))
	}
	return s
}

// SanitizeAndValidate sanitizes the struct data points to in place (see the
// go-sanitize "san" struct tag), then validates it (see the validator
// "validate" struct tag). It returns nil when the data is valid, otherwise one
// FieldError per failed rule.
//
// data must be a non-nil pointer to a struct; anything else is a programming
// error and panics.
func SanitizeAndValidate(data any) []FieldError {
	if v := reflect.ValueOf(data); v.Kind() != reflect.Pointer || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("nvtval: SanitizeAndValidate needs a non-nil pointer to a struct, got %T", data))
	}

	if err := sanitizer.Sanitize(data); err != nil {
		panic(fmt.Sprintf("nvtval: sanitize %T: %v", data, err))
	}

	return Validate(data)
}

// Validate validates data (a struct or a pointer to one) against its
// "validate" struct tags without changing it. It returns nil when the data is
// valid, otherwise one FieldError per failed rule.
//
// Input that cannot be validated at all, such as nil or a value that is not a
// struct, is a programming error and panics rather than being reported as
// valid.
func Validate(data any) []FieldError {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	failures, ok := errors.AsType[validator.ValidationErrors](err)
	if !ok {
		panic(fmt.Sprintf("nvtval: cannot validate %T: %v", data, err))
	}

	out := make([]FieldError, len(failures))
	for i, failure := range failures {
		out[i] = FieldError{
			Field:   failure.Field(),
			Message: message(failure.Tag(), failure.Field(), failure.Param(), failure.Kind()),
		}
	}
	return out
}
