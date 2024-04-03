package response

import (
	"github.com/go-playground/validator/v10"
)

type ErrorMsg struct {
	Code    int    `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// GetErrorMsg returns the error message based on the validation tag.
func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "max":
		return "length of " + fe.Field() + " must not be over."
	case "email":
		return "email field is not in the correct format."
	case "eq=owner|eq=admin|eq=user":
		return fe.Field() + " must be one of the allowed."
	}
	return "Unknown error"
}
