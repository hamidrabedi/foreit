package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator with framework-specific methods
type Validator struct {
	*validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	v := validator.New()
	return &Validator{Validate: v}
}

// ValidateStruct validates a struct
func (v *Validator) ValidateStruct(s interface{}) error {
	if err := v.Struct(s); err != nil {
		return formatValidationErrors(err)
	}
	return nil
}

// ValidateField validates a single field
func (v *Validator) ValidateField(field interface{}, tag string) error {
	if err := v.Var(field, tag); err != nil {
		return formatValidationErrors(err)
	}
	return nil
}

// formatValidationErrors formats validator errors into a readable format
func formatValidationErrors(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errMsg string
		for i, fe := range validationErrors {
			if i > 0 {
				errMsg += "; "
			}
			errMsg += fmt.Sprintf("%s: %s", fe.Field(), getErrorMessage(fe))
		}
		return fmt.Errorf("validation failed: %s", errMsg)
	}
	return err
}

// getErrorMessage returns a user-friendly error message for a validation error
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s characters", fe.Param())
	case "numeric":
		return "must be numeric"
	case "alpha":
		return "must contain only letters"
	case "alphanum":
		return "must contain only letters and numbers"
	default:
		return fmt.Sprintf("failed validation for tag '%s'", fe.Tag())
	}
}

// RegisterCustomValidator registers a custom validation function
func (v *Validator) RegisterCustomValidator(tag string, fn validator.Func) error {
	return v.RegisterValidation(tag, fn)
}
