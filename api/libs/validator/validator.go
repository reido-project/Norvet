package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validation = validator.New()

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Validate(data any) []ErrorResponse {
	err := validation.Struct(data)
	if err == nil {
		return nil
	}

	var out []ErrorResponse

	if validationErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
		for _, fieldErr := range validationErrors {
			var element ErrorResponse
			element.Field = fieldErr.Field()

			param := fieldErr.Param()
			if param != "" {
				element.Message = translateWithParam(fieldErr.Tag(), fieldErr.Field(), param)
			} else {
				element.Message = translateTagOnly(fieldErr.Tag(), fieldErr.Field())
			}

			out = append(out, element)
		}
	}

	return out
}
