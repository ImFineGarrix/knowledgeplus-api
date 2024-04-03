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
		return "lenght of " + fe.Field() + " must not be over."
	case "email":
		return "email field is not correct format."
	case "eq":
		if fe.Field() == "role" {
			return fe.Field() + " must be correct role."
		} else {
			return fe.Field() + " must be equal to " + fe.Param()
		}
	}
	return "Unknown error"
}
