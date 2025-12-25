package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type StructValidator struct {
	validate *validator.Validate
}

func NewStructValidator() *StructValidator {
	v := validator.New()
	registerCustomValidators(v)

	return &StructValidator{
		validate: v,
	}
}

func (v *StructValidator) Validate(obj interface{}) error {
	if err := v.validate.Struct(obj); err != nil {
		return formatValidationError(err)
	}
	return nil
}

func (v *StructValidator) ValidateVar(field interface{}, tag string) error {
	if err := v.validate.Var(field, tag); err != nil {
		return formatValidationError(err)
	}
	return nil
}

func (v *StructValidator) RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error {
	return v.validate.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}

func (v *StructValidator) RegisterStructValidation(fn validator.StructLevelFunc, types ...interface{}) {
	v.validate.RegisterStructValidation(fn, types...)
}

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

func registerCustomValidators(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
