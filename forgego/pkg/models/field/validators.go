package field

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Validator validates a field value.
type Validator interface {
	Validate(value interface{}) error
}

// EmailValidator validates email format using regex.
type EmailValidator struct{}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Validate checks if the value is a valid email address.
func (v EmailValidator) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return errors.New("email must be a string")
	}
	if str == "" {
		return nil
	}
	if !emailRegex.MatchString(str) {
		return fmt.Errorf("invalid email format: %s", str)
	}
	return nil
}

// MinValueValidator validates that a numeric value is not less than the minimum.
type MinValueValidator struct {
	Min interface{}
}

// Validate checks if the value is greater than or equal to the minimum.
func (v MinValueValidator) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	switch val := value.(type) {
	case int:
		min, ok := v.Min.(int)
		if !ok {
			return errors.New("min value type mismatch")
		}
		if val < min {
			return fmt.Errorf("value %d is less than minimum %d", val, min)
		}
	case int64:
		min, ok := v.Min.(int64)
		if !ok {
			return errors.New("min value type mismatch")
		}
		if val < min {
			return fmt.Errorf("value %d is less than minimum %d", val, min)
		}
	case float64:
		min, ok := v.Min.(float64)
		if !ok {
			return errors.New("min value type mismatch")
		}
		if val < min {
			return fmt.Errorf("value %f is less than minimum %f", val, min)
		}
	}
	return nil
}

// MaxValueValidator validates that a numeric value is not greater than the maximum.
type MaxValueValidator struct {
	Max interface{}
}

// Validate checks if the value is less than or equal to the maximum.
func (v MaxValueValidator) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	switch val := value.(type) {
	case int:
		max, ok := v.Max.(int)
		if !ok {
			return errors.New("max value type mismatch")
		}
		if val > max {
			return fmt.Errorf("value %d is greater than maximum %d", val, max)
		}
	case int64:
		max, ok := v.Max.(int64)
		if !ok {
			return errors.New("max value type mismatch")
		}
		if val > max {
			return fmt.Errorf("value %d is greater than maximum %d", val, max)
		}
	case float64:
		max, ok := v.Max.(float64)
		if !ok {
			return errors.New("max value type mismatch")
		}
		if val > max {
			return fmt.Errorf("value %f is greater than maximum %f", val, max)
		}
	}
	return nil
}

// MinLengthValidator validates that a string has a minimum length.
type MinLengthValidator struct {
	Min int
}

// Validate checks if the string length is at least the minimum.
func (v MinLengthValidator) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return errors.New("min length validator only works with strings")
	}
	if len(str) < v.Min {
		return fmt.Errorf("string length %d is less than minimum %d", len(str), v.Min)
	}
	return nil
}

// MaxLengthValidator validates that a string has a maximum length.
type MaxLengthValidator struct {
	Max int
}

// Validate checks if the string length is at most the maximum.
func (v MaxLengthValidator) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return errors.New("max length validator only works with strings")
	}
	if len(str) > v.Max {
		return fmt.Errorf("string length %d is greater than maximum %d", len(str), v.Max)
	}
	return nil
}

// RequiredValidator validates that a value is not empty.
type RequiredValidator struct{}

// Validate checks if the value is present and not empty.
func (v RequiredValidator) Validate(value interface{}) error {
	if value == nil {
		return errors.New("field is required")
	}
	switch val := value.(type) {
	case string:
		if strings.TrimSpace(val) == "" {
			return errors.New("field is required")
		}
	case int, int64, float64, bool:
		return nil
	}
	return nil
}

// ValidateField validates a field value using all its validators.
func ValidateField[T any, V any](f Field[T, V], value interface{}) error {
	for _, validator := range f.Validators {
		if err := validator.Validate(value); err != nil {
			return err
		}
	}
	return nil
}

