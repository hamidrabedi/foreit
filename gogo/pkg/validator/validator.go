package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator
type Validator struct {
	validate *validator.Validate
}

// New creates a new validator instance
func New() *Validator {
	v := validator.New()
	
	// Register custom validators
	registerCustomValidators(v)
	
	return &Validator{
		validate: v,
	}
}

// Validate validates a struct
func (v *Validator) Validate(obj interface{}) error {
	if err := v.validate.Struct(obj); err != nil {
		return formatValidationError(err)
	}
	return nil
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	if err := v.validate.Var(field, tag); err != nil {
		return formatValidationError(err)
	}
	return nil
}

// RegisterValidation registers a custom validation function
func (v *Validator) RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error {
	return v.validate.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}

// RegisterStructValidation registers a struct-level validation
func (v *Validator) RegisterStructValidation(fn validator.StructLevelFunc, types ...interface{}) {
	v.validate.RegisterStructValidation(fn, types...)
}

// formatValidationError formats validation errors into a readable format
func formatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []string
		for _, e := range validationErrors {
			errors = append(errors, formatFieldError(e))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
	}
	return err
}

// formatFieldError formats a single field error
func formatFieldError(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()
	param := e.Param()
	
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	default:
		return fmt.Sprintf("%s failed validation for tag '%s'", field, tag)
	}
}

// registerCustomValidators registers custom validation functions
func registerCustomValidators(v *validator.Validate) {
	// Register custom tag name function to use JSON tags
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
