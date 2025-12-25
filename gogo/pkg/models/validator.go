package models

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator provides Pydantic-inspired validation
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	v := validator.New()
	
	// Register custom validators
	v.RegisterValidation("min_length", validateMinLength)
	v.RegisterValidation("max_length", validateMaxLength)
	v.RegisterValidation("email", validateEmail)
	v.RegisterValidation("url", validateURL)
	
	return &Validator{validate: v}
}

// Validate validates a model struct
func (v *Validator) Validate(model interface{}) error {
	if err := v.validate.Struct(model); err != nil {
		return formatValidationError(err)
	}
	return nil
}

// ValidateField validates a single field
func (v *Validator) ValidateField(model interface{}, fieldName string) error {
	return v.validate.Var(getFieldValue(model, fieldName), getFieldTag(model, fieldName, "validate"))
}

// GetValidationErrors returns detailed validation errors
func (v *Validator) GetValidationErrors(model interface{}) map[string][]string {
	errs := make(map[string][]string)
	
	if err := v.validate.Struct(model); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validationErrors {
				field := err.Field()
				if errs[field] == nil {
					errs[field] = make([]string, 0)
				}
				errs[field] = append(errs[field], formatFieldError(err))
			}
		}
	}
	
	return errs
}

// FieldValidator is a function that validates a field value
type FieldValidator func(value interface{}) error

// RegisterFieldValidator registers a custom field validator
func (v *Validator) RegisterFieldValidator(name string, validator FieldValidator) {
	v.validate.RegisterValidation(name, func(fl validator.FieldLevel) bool {
		return validator(fl.Field().Interface()) == nil
	})
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
}

func formatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		messages := make([]string, 0, len(validationErrors))
		for _, err := range validationErrors {
			messages = append(messages, formatFieldError(err))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
	}
	return err
}

func formatFieldError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", err.Field(), err.Param())
	case "min_length":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max_length":
		return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", err.Field())
	default:
		return fmt.Sprintf("%s failed validation: %s", err.Field(), err.Tag())
	}
}

func validateMinLength(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	minLen, _ := parseInt(fl.Param())
	return len(value) >= minLen
}

func validateMaxLength(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	maxLen, _ := parseInt(fl.Param())
	return len(value) <= maxLen
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func validateURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func getFieldValue(model interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}
	return field.Interface()
}

func getFieldTag(model interface{}, fieldName, tagName string) string {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	field, ok := typ.FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(tagName)
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// DefaultValidator is the default validator instance
var DefaultValidator = NewValidator()

