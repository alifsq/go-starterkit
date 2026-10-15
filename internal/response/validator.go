package response

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationError(err error) []ErrorDetail {
	if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
		out := make([]ErrorDetail, len(ve))
		for i, fe := range ve {
			out[i] = ErrorDetail{
				Field:   fe.Field(),
				Message: getErrorMessage(fe),
			}
		}
		return out
	}
	return nil
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		if fe.Kind() == reflect.String {
			return "Minimum " + fe.Param() + " characters required"
		}
		return "Minimum value is " + fe.Param()
	case "max":
		if fe.Kind() == reflect.String {
			return "Maximum " + fe.Param() + " characters allowed"
		}
		return "Maximum value is " + fe.Param()
	case "oneof":
		return "Must be one of the following: " + strings.ReplaceAll(fe.Param(), " ", ", ")
	case "eqfield":
		return "Must match the " + toSnakeCase(fe.Param()) + " field"
	}
	return "Invalid field value"
}

// toSnakeCase 
func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}