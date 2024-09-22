package domain

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
}

func FormatValidationError(err error) error {
	var errorMessages []string

	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		switch err.Tag() {
		case "required":
			errorMessages = append(errorMessages, fmt.Sprintf("%s is a required field", fieldName))
		case "email":
			errorMessages = append(errorMessages, fmt.Sprintf("%s must be a valid email address", fieldName))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("%s is not valid", fieldName))
		}
	}

	return fmt.Errorf("%s", strings.Join(errorMessages, ", "))
}
