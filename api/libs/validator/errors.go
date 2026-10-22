package validator

import "fmt"

func translateWithParam(tag string, field string, param string) string {
	switch tag {
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, param)
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of the following values: [%s]", field, param)
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
	default:
		return fmt.Sprintf("Failed validation rule '%s' with value '%s'", tag, param)
	}
}

func translateTagOnly(tag string, field string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required and must not be empty", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "numeric":
		return fmt.Sprintf("%s must be a number", field)
	case "uuid", "uuid3", "uuid4", "uuid5":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "boolean":
		return fmt.Sprintf("%s must be true or false", field)
	case "lowercase":
		return fmt.Sprintf("%s must be in lowercase", field)
	case "uppercase":
		return fmt.Sprintf("%s must be in uppercase", field)
	default:
		return fmt.Sprintf("%s failed on validation rule '%s'", field, tag)
	}
}
