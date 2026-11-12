# nvtval

Sanitizes and validates request structs, and turns failures into client-ready messages.

```go
type RegisterRequest struct {
    Email    string `json:"email"    san:"trim,lower" validate:"required,email"`
    Password string `json:"password"                  validate:"required,min=8"`
    Age      int    `json:"age"                       validate:"min=18"`
}

var req RegisterRequest
if err := c.Bind().Body(&req); err != nil {
    return err
}
if problems := nvtval.SanitizeAndValidate(&req); problems != nil {
    return c.Status(fiber.StatusUnprocessableEntity).JSON(problems)
}
```

```json
[{"field":"Age","message":"Age must be at least 18"}]
```

| Function | Purpose |
|---|---|
| `SanitizeAndValidate(&v)` | sanitize `v` in place (`san` tags), then validate it (`validate` tags) |
| `Validate(v)` | validate only, without modifying `v` |

Both return `nil` when the data is valid, otherwise one `FieldError{Field, Message}` per failed rule. `Field` is the struct field name.

## Messages

Size rules adapt to the field type:

| Rule | string | slice / map | number |
|---|---|---|---|
| `min=3` | must be at least 3 characters | must contain at least 3 items | must be at least 3 |
| `max=3` | must not exceed 3 characters | must not contain more than 3 items | must not exceed 3 |
| `len=3` | must be exactly 3 characters long | must contain exactly 3 items | must be exactly 3 |
| `gt` `gte` `lt` `lte` | longer than, at least, shorter than, at most | more than, at least, fewer than, at most | greater than, at least, less than, at most |

Other rules with a message: `required`, `email`, `url`, `alpha`, `alphanum`, `numeric`, `number`, `uuid*`, `boolean`, `lowercase`, `uppercase`, `ascii`, `hexadecimal`, `base64`, `json`, `ip`, `ipv4`, `ipv6`, `e164`, `oneof`, `eq`, `ne`, `eqfield`, `nefield`, `gtfield`, `ltfield`, `datetime`, `contains`, `startswith`, `endswith`, `excludes`. Any other rule gives `<Field> failed on validation rule '<rule>'`.

## Programming errors panic

`SanitizeAndValidate` needs a non-nil pointer to a struct. `Validate` needs a struct or a pointer to one. Anything else (nil, a string, a nil pointer) panics with a message naming the type, instead of being reported as valid. Fiber's recover middleware turns that into a 500.

Requires Go 1.26 (`errors.AsType`).

## Renamed from the previous version

| Before | Now |
|---|---|
| `Response` | `FieldError` (same JSON) |
