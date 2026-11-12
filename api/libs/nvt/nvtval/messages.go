package nvtval

import (
	"fmt"
	"reflect"
)

// message returns the human-readable message for a failed validation rule.
// tag and param come from the validator, field is the name of the field and
// kind its type: size rules (min, max, ...) read differently for strings,
// collections and numbers.
func message(tag, field, param string, kind reflect.Kind) string {
	if param == "" {
		return tagMessage(tag, field)
	}
	return paramMessage(tag, field, param, kind)
}

// sizeText holds the message templates of a size rule. Each takes the field
// name and the rule's parameter.
type sizeText struct{ str, items, num string }

var sizeRules = map[string]sizeText{
	"min": {"%s must be at least %s characters", "%s must contain at least %s items", "%s must be at least %s"},
	"max": {"%s must not exceed %s characters", "%s must not contain more than %s items", "%s must not exceed %s"},
	"len": {"%s must be exactly %s characters long", "%s must contain exactly %s items", "%s must be exactly %s"},
	"gt":  {"%s must be longer than %s characters", "%s must contain more than %s items", "%s must be greater than %s"},
	"lt":  {"%s must be shorter than %s characters", "%s must contain fewer than %s items", "%s must be less than %s"},
}

func init() {
	sizeRules["gte"] = sizeRules["min"]
	sizeRules["lte"] = sizeRules["max"]
}

func paramMessage(tag, field, param string, kind reflect.Kind) string {
	if text, ok := sizeRules[tag]; ok {
		switch kind {
		case reflect.String:
			return fmt.Sprintf(text.str, field, param)
		case reflect.Slice, reflect.Array, reflect.Map:
			return fmt.Sprintf(text.items, field, param)
		}
		return fmt.Sprintf(text.num, field, param)
	}

	switch tag {
	case "oneof":
		return fmt.Sprintf("%s must be one of the following values: [%s]", field, param)
	case "eq":
		return fmt.Sprintf("%s must be equal to %s", field, param)
	case "ne":
		return fmt.Sprintf("%s must not be equal to %s", field, param)
	case "eqfield":
		return fmt.Sprintf("%s didn't match", field)
	case "nefield":
		return fmt.Sprintf("%s cannot be the same as %s", field, param)
	case "gtfield":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "ltfield":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "datetime":
		return fmt.Sprintf("%s must match the date/time format: %s", field, param)
	case "contains":
		return fmt.Sprintf("%s must contain '%s'", field, param)
	case "startswith":
		return fmt.Sprintf("%s must start with '%s'", field, param)
	case "endswith":
		return fmt.Sprintf("%s must end with '%s'", field, param)
	case "excludes":
		if param == " " {
			return fmt.Sprintf("%s cannot contain space in between characters", field)
		}
		return fmt.Sprintf("%s cannot contain character '%s'", field, param)
	}
	return fmt.Sprintf("%s failed on validation rule '%s' with value '%s'", field, tag, param)
}

func tagMessage(tag, field string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required and must not be empty", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "alpha", "alphaunicode":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum", "alphanumunicode":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "numeric":
		return fmt.Sprintf("%s must be a number", field)
	case "number":
		return fmt.Sprintf("%s must contain only digits", field)
	case "uuid", "uuid3", "uuid4", "uuid5":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "boolean":
		return fmt.Sprintf("%s must be true or false", field)
	case "lowercase":
		return fmt.Sprintf("%s must be in lowercase", field)
	case "uppercase":
		return fmt.Sprintf("%s must be in uppercase", field)
	case "ascii":
		return fmt.Sprintf("%s must contain only ASCII characters", field)
	case "hexadecimal":
		return fmt.Sprintf("%s must be a hexadecimal value", field)
	case "base64":
		return fmt.Sprintf("%s must be base64 encoded", field)
	case "json":
		return fmt.Sprintf("%s must be valid JSON", field)
	case "ip":
		return fmt.Sprintf("%s must be a valid IP address", field)
	case "ipv4":
		return fmt.Sprintf("%s must be a valid IPv4 address", field)
	case "ipv6":
		return fmt.Sprintf("%s must be a valid IPv6 address", field)
	case "e164":
		return fmt.Sprintf("%s must be a valid phone number in international format (E.164)", field)
	}
	return fmt.Sprintf("%s failed on validation rule '%s'", field, tag)
}
