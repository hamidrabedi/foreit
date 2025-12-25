package api

import (
	"fmt"

	"github.com/forgego/forge/pkg/models"
	"github.com/forgego/forge/pkg/models/field"
)

type FieldAdapter struct {
	modelFieldDef models.FieldDefinition
	readOnly      bool
	writeOnly     bool
	source        string
	fieldName     string
	toRepr        func(interface{}) (interface{}, error)
	toInternal    func(interface{}) (interface{}, error)
	validators    []field.Validator
}

func NewFieldAdapterFromModelField(
	modelFieldDef models.FieldDefinition,
	fieldName string,
) *FieldAdapter {
	adapter := &FieldAdapter{
		modelFieldDef: modelFieldDef,
		fieldName:     fieldName,
		source:        modelFieldDef.GetName(),
		readOnly:      false,
		writeOnly:     false,
		validators:    []field.Validator{},
	}
	
	if withValidators, ok := modelFieldDef.(models.FieldDefinitionWithValidators); ok {
		validators := withValidators.GetValidators()
		for _, v := range validators {
			if validator, ok := v.(field.Validator); ok {
				adapter.validators = append(adapter.validators, validator)
			}
		}
	}
	
	adapter.toRepr = adapter.defaultToRepresentation
	adapter.toInternal = adapter.defaultToInternalValue
	
	return adapter
}

func (a *FieldAdapter) GetName() string {
	return a.fieldName
}

func (a *FieldAdapter) GetSource() string {
	if a.source != "" {
		return a.source
	}
	return a.fieldName
}

func (a *FieldAdapter) GetReadOnly() bool {
	return a.readOnly
}

func (a *FieldAdapter) GetWriteOnly() bool {
	return a.writeOnly
}

func (a *FieldAdapter) GetRequired() bool {
	for _, v := range a.validators {
		if _, ok := v.(field.RequiredValidator); ok {
			return true
		}
	}
	return false
}

func (a *FieldAdapter) GetDefault() interface{} {
	return nil
}

func (a *FieldAdapter) GetAllowNull() bool {
	return true
}

func (a *FieldAdapter) ToRepresentation(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}
	return a.toRepr(value)
}

func (a *FieldAdapter) ToInternalValue(data interface{}) (interface{}, error) {
	if data == nil {
		if !a.GetAllowNull() {
			return nil, fmt.Errorf("this field may not be null")
		}
		return nil, nil
	}
	
	value, err := a.toInternal(data)
	if err != nil {
		return nil, err
	}
	
	for _, validator := range a.validators {
		if err := validator.Validate(value); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
	}
	
	return value, nil
}

func (a *FieldAdapter) Validate(value interface{}) error {
	if a.GetRequired() && value == nil {
		return fmt.Errorf("this field is required")
	}
	if !a.GetAllowNull() && value == nil {
		return fmt.Errorf("this field may not be null")
	}
	
	for _, validator := range a.validators {
		if err := validator.Validate(value); err != nil {
			return err
		}
	}
	
	return nil
}

func (a *FieldAdapter) ReadOnly(readOnly bool) *FieldAdapter {
	a.readOnly = readOnly
	return a
}

func (a *FieldAdapter) WriteOnly(writeOnly bool) *FieldAdapter {
	a.writeOnly = writeOnly
	return a
}

func (a *FieldAdapter) Source(source string) *FieldAdapter {
	a.source = source
	return a
}

func (a *FieldAdapter) WithValidators(validators ...field.Validator) *FieldAdapter {
	a.validators = append(a.validators, validators...)
	return a
}

func (a *FieldAdapter) WithToRepresentation(fn func(interface{}) (interface{}, error)) *FieldAdapter {
	a.toRepr = fn
	return a
}

func (a *FieldAdapter) WithToInternalValue(fn func(interface{}) (interface{}, error)) *FieldAdapter {
	a.toInternal = fn
	return a
}

func (a *FieldAdapter) GetModelFieldDefinition() models.FieldDefinition {
	return a.modelFieldDef
}

func (a *FieldAdapter) defaultToRepresentation(value interface{}) (interface{}, error) {
	return value, nil
}

func (a *FieldAdapter) defaultToInternalValue(data interface{}) (interface{}, error) {
	return data, nil
}

func AdaptModelFields(modelFields []models.FieldDefinition) map[string]Field {
	adapters := make(map[string]Field)
	
	for _, modelField := range modelFields {
		fieldName := modelField.GetName()
		adapter := NewFieldAdapterFromModelField(modelField, fieldName)
		adapters[fieldName] = adapter
	}
	
	return adapters
}

