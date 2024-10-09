package domain

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()

	if err := validate.RegisterValidation("enum_potential", func(fl validator.FieldLevel) bool {
		potential := fl.Field().String()
		return potential == "minimum" || potential == "recommended" || potential == "maximum"
	}); err != nil {
		fmt.Printf("Error registering validation enum_potential: %v\n", err)
	}

	if err := validate.RegisterValidation("enum_os", func(fl validator.FieldLevel) bool {
		os := fl.Field().String()
		return os == "windows" || os == "mac" || os == "linux"
	}); err != nil {
		fmt.Printf("Error registering validation enum_os: %v\n", err)
	}
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
		case "enum_potential":
			errorMessages = append(errorMessages, fmt.Sprintf("%s must be one of 'minimum', 'recommended', or 'maximum'", fieldName))
		case "enum_os":
			errorMessages = append(errorMessages, fmt.Sprintf("%s must be one of 'windows', 'mac', or 'linux'", fieldName))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("%s is not valid", fieldName))
		}
	}

	return fmt.Errorf("%s", strings.Join(errorMessages, ", "))
}
