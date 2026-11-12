package nvterr

import (
	"net/http"
)

// Domain "validation": field validation, credential quality, and
// list-endpoint query problems.

var (
	ErrValidationFailed      = define("validation_failed", http.StatusUnprocessableEntity, "validation failed")
	ErrUnprocessable         = define("unprocessable", http.StatusUnprocessableEntity, "request cannot be processed")
	ErrBusinessRuleViolation = define("business_rule_violation", http.StatusUnprocessableEntity, "business rule violation")
	ErrNotSupported          = define("not_supported", http.StatusUnprocessableEntity, "operation not supported")
)

// Field-level validation. Everything here is a child of ErrValidationFailed.
var (
	ErrRequired = defineChild(ErrValidationFailed, "required", "field is required")

	// Format.
	ErrInvalidFormat      = defineChild(ErrValidationFailed, "invalid_format", "field has an invalid format")
	ErrInvalidEmail       = defineChild(ErrInvalidFormat, "invalid_email", "email address is invalid")
	ErrInvalidPhoneNumber = defineChild(ErrInvalidFormat, "invalid_phone_number", "phone number is invalid")
	ErrInvalidURL         = defineChild(ErrInvalidFormat, "invalid_url", "URL is invalid")
	ErrInvalidIPAddress   = defineChild(ErrInvalidFormat, "invalid_ip_address", "IP address is invalid")
	ErrInvalidUsername    = defineChild(ErrInvalidFormat, "invalid_username", "username is invalid")
	ErrInvalidPostalCode  = defineChild(ErrInvalidFormat, "invalid_postal_code", "postal code is invalid")
	ErrInvalidCountryCode = defineChild(ErrInvalidFormat, "invalid_country_code", "country code is invalid")
	ErrInvalidCurrency    = defineChild(ErrInvalidFormat, "invalid_currency", "currency code is invalid")
	ErrInvalidLocale      = defineChild(ErrInvalidFormat, "invalid_locale", "locale is invalid")
	ErrInvalidTimezone    = defineChild(ErrInvalidFormat, "invalid_timezone", "timezone is invalid")
	ErrInvalidBase64      = defineChild(ErrInvalidFormat, "invalid_base64", "value is not valid base64")

	// Dates.
	ErrInvalidDate      = defineChild(ErrValidationFailed, "invalid_date", "date is invalid")
	ErrDateInPast       = defineChild(ErrInvalidDate, "date_in_past", "date must not be in the past")
	ErrDateInFuture     = defineChild(ErrInvalidDate, "date_in_future", "date must not be in the future")
	ErrInvalidDateRange = defineChild(ErrInvalidDate, "invalid_date_range", "end date must not be before start date")

	// Size and range.
	ErrInvalidLength = defineChild(ErrValidationFailed, "invalid_length", "value has an invalid length")
	ErrValueTooShort = defineChild(ErrInvalidLength, "value_too_short", "value is too short")
	ErrValueTooLong  = defineChild(ErrInvalidLength, "value_too_long", "value is too long")

	ErrValueOutOfRange = defineChild(ErrValidationFailed, "value_out_of_range", "value is out of range")
	ErrValueTooSmall   = defineChild(ErrValueOutOfRange, "value_too_small", "value is below the minimum")
	ErrValueTooLarge   = defineChild(ErrValueOutOfRange, "value_too_large", "value is above the maximum")

	ErrTooFewItems  = defineChild(ErrValidationFailed, "too_few_items", "too few items")
	ErrTooManyItems = defineChild(ErrValidationFailed, "too_many_items", "too many items")

	// Allowed values.
	ErrValueNotAllowed = defineChild(ErrValidationFailed, "value_not_allowed", "value is not allowed")
	ErrInvalidEnum     = defineChild(ErrValueNotAllowed, "invalid_enum", "value is not one of the allowed options")
	ErrDuplicateValue  = defineChild(ErrValidationFailed, "duplicate_value", "duplicate value")
)

// Credential quality.
var (
	ErrPasswordTooWeak          = defineChild(ErrValidationFailed, "password_too_weak", "password does not meet the strength requirements")
	ErrPasswordTooShort         = defineChild(ErrPasswordTooWeak, "password_too_short", "password is too short")
	ErrPasswordTooLong          = defineChild(ErrPasswordTooWeak, "password_too_long", "password is too long")
	ErrPasswordCompromised      = defineChild(ErrPasswordTooWeak, "password_compromised", "password appears in a known data breach")
	ErrPasswordContainsUserInfo = defineChild(ErrPasswordTooWeak, "password_contains_user_info", "password must not contain your name or email")
	ErrPasswordReused           = defineChild(ErrValidationFailed, "password_reused", "password was used recently and cannot be reused")
	ErrPasswordSameAsCurrent    = defineChild(ErrPasswordReused, "password_same_as_current", "new password must differ from the current password")
)

// List endpoints (400: the query itself is wrong).
var (
	ErrInvalidPagination = defineChild(ErrInvalidParameter, "invalid_pagination", "invalid pagination parameters")
	ErrInvalidPage       = defineChild(ErrInvalidPagination, "invalid_page", "invalid page number")
	ErrInvalidPageSize   = defineChild(ErrInvalidPagination, "invalid_page_size", "invalid page size")
	ErrInvalidCursor     = defineChild(ErrInvalidPagination, "invalid_cursor", "invalid or expired pagination cursor")

	ErrInvalidSort      = defineChild(ErrInvalidParameter, "invalid_sort", "invalid sort parameter")
	ErrInvalidSortField = defineChild(ErrInvalidSort, "invalid_sort_field", "cannot sort by this field")
	ErrInvalidSortOrder = defineChild(ErrInvalidSort, "invalid_sort_order", "sort order must be asc or desc")

	ErrInvalidFilter      = defineChild(ErrInvalidParameter, "invalid_filter", "invalid filter parameter")
	ErrInvalidSearchQuery = defineChild(ErrInvalidParameter, "invalid_search_query", "invalid search query")
)
