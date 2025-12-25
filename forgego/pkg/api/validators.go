package api

import (
	"github.com/forgego/forge/pkg/models/field"
)

type Validator interface {
	Validate(value interface{}) error
}

type ValidatorAdapter struct {
	validator field.Validator
}

func NewValidatorAdapter(v field.Validator) Validator {
	return &ValidatorAdapter{validator: v}
}

func (a *ValidatorAdapter) Validate(value interface{}) error {
	return a.validator.Validate(value)
}

func AdaptValidators(validators []field.Validator) []Validator {
	adapted := make([]Validator, len(validators))
	for i, v := range validators {
		adapted[i] = NewValidatorAdapter(v)
	}
	return adapted
}

