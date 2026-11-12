// Package nvtval sanitizes and validates request structs and turns validation
// failures into messages that are ready to send to a client.
//
// Describe the rules with struct tags and call SanitizeAndValidate:
//
//	type RegisterRequest struct {
//		Email    string `json:"email"    san:"trim,lower" validate:"required,email"`
//		Password string `json:"password"                  validate:"required,min=8"`
//		Age      int    `json:"age"                       validate:"min=18"`
//	}
//
//	var req RegisterRequest
//	if err := c.Bind().Body(&req); err != nil {
//		return err
//	}
//	if problems := nvtval.SanitizeAndValidate(&req); problems != nil {
//		return c.Status(fiber.StatusUnprocessableEntity).JSON(problems)
//	}
//
// problems is a []FieldError, which marshals to
//
//	[{"field":"Age","message":"Age must be at least 18"}]
//
// Sanitizing uses github.com/go-sanitize/sanitize ("san" tags) and validating
// uses github.com/go-playground/validator ("validate" tags). Use Validate to
// validate without modifying the struct.
//
// Messages for size rules (min, max, len, gt, gte, lt, lte) depend on the
// type of the field: "must be at least 3 characters" for strings, "must
// contain at least 3 items" for slices and maps, "must be at least 3" for
// numbers. Rules without a dedicated message produce
// "<Field> failed on validation rule '<rule>'".
//
// Input that cannot be validated at all (nil, or not a struct) is a
// programming error and panics instead of being reported as valid.
//
// # Layout
//
//	validate.go   SanitizeAndValidate, Validate, FieldError
//	messages.go   message texts per validation rule
package nvtval
