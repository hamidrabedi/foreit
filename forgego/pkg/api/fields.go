package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/forgego/forge/pkg/models/field"
)

type Field interface {
	ToRepresentation(value interface{}) (interface{}, error)
	ToInternalValue(data interface{}) (interface{}, error)
	Validate(value interface{}) error
	GetReadOnly() bool
	GetWriteOnly() bool
	GetRequired() bool
	GetDefault() interface{}
	GetSource() string
	GetAllowNull() bool
}

type BaseField struct {
	ReadOnly   bool
	WriteOnly  bool
	Required   bool
	Default    interface{}
	Source     string
	AllowNull  bool
	validators []Validator
	modelValidators []field.Validator
}

func NewBaseField() *BaseField {
	return &BaseField{
		ReadOnly:   false,
		WriteOnly:  false,
		Required:   false,
		Default:    nil,
		Source:     "",
		AllowNull:  true,
		validators: []Validator{},
	}
}

func (f *BaseField) GetReadOnly() bool {
	return f.ReadOnly
}

func (f *BaseField) GetWriteOnly() bool {
	return f.WriteOnly
}

func (f *BaseField) GetRequired() bool {
	return f.Required
}

func (f *BaseField) GetDefault() interface{} {
	return f.Default
}

func (f *BaseField) GetSource() string {
	return f.Source
}

func (f *BaseField) GetAllowNull() bool {
	return f.AllowNull
}

func (f *BaseField) Validate(value interface{}) error {
	if f.Required && value == nil {
		return fmt.Errorf("this field is required")
	}
	if !f.AllowNull && value == nil {
		return fmt.Errorf("this field may not be null")
	}
	for _, validator := range f.validators {
		if err := validator.Validate(value); err != nil {
			return err
		}
	}
	for _, validator := range f.modelValidators {
		if err := validator.Validate(value); err != nil {
			return err
		}
	}
	return nil
}

func (f *BaseField) WithValidators(validators ...Validator) *BaseField {
	f.validators = append(f.validators, validators...)
	return f
}

func (f *BaseField) WithModelValidators(validators ...field.Validator) *BaseField {
	f.modelValidators = append(f.modelValidators, validators...)
	return f
}

type StringField struct {
	*BaseField
	MaxLength *int
	MinLength *int
	AllowBlank bool
}

func NewStringField() *StringField {
	return &StringField{
		BaseField: NewBaseField(),
		AllowBlank: true,
	}
}

func (f *StringField) ToRepresentation(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}
	if str, ok := value.(string); ok {
		return str, nil
	}
	return fmt.Sprintf("%v", value), nil
}

func (f *StringField) ToInternalValue(data interface{}) (interface{}, error) {
	if data == nil {
		if f.Default != nil {
			return f.Default, nil
		}
		if !f.AllowNull {
			return nil, fmt.Errorf("this field may not be null")
		}
		return nil, nil
	}
	if str, ok := data.(string); ok {
		if !f.AllowBlank && str == "" {
			return nil, fmt.Errorf("this field may not be blank")
		}
		if f.MinLength != nil && len(str) < *f.MinLength {
			return nil, fmt.Errorf("ensure this field has at least %d characters", *f.MinLength)
		}
		if f.MaxLength != nil && len(str) > *f.MaxLength {
			return nil, fmt.Errorf("ensure this field has no more than %d characters", *f.MaxLength)
		}
		return str, nil
	}
	return nil, fmt.Errorf("expected string, got %T", data)
}

func (f *StringField) MaxLengthOption(length int) *StringField {
	f.MaxLength = &length
	return f
}

func (f *StringField) MinLengthOption(length int) *StringField {
	f.MinLength = &length
	return f
}

func (f *StringField) AllowBlankOption(allow bool) *StringField {
	f.AllowBlank = allow
	return f
}

type IntegerField struct {
	*BaseField
	MaxValue *int64
	MinValue *int64
}

func NewIntegerField() *IntegerField {
	return &IntegerField{
		BaseField: NewBaseField(),
	}
}

func (f *IntegerField) ToRepresentation(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	}
	return nil, fmt.Errorf("cannot convert %T to integer", value)
}

func (f *IntegerField) ToInternalValue(data interface{}) (interface{}, error) {
	if data == nil {
		if f.Default != nil {
			return f.Default, nil
		}
		if !f.AllowNull {
			return nil, fmt.Errorf("this field may not be null")
		}
		return nil, nil
	}
	var val int64
	switch v := data.(type) {
	case int:
		val = int64(v)
	case int8:
		val = int64(v)
	case int16:
		val = int64(v)
	case int32:
		val = int64(v)
	case int64:
		val = v
	case uint:
		val = int64(v)
	case uint8:
		val = int64(v)
	case uint16:
		val = int64(v)
	case uint32:
		val = int64(v)
	case uint64:
		val = int64(v)
	case float64:
		val = int64(v)
	case float32:
		val = int64(v)
	default:
		return nil, fmt.Errorf("expected integer, got %T", data)
	}
	if f.MinValue != nil && val < *f.MinValue {
		return nil, fmt.Errorf("ensure this value is greater than or equal to %d", *f.MinValue)
	}
	if f.MaxValue != nil && val > *f.MaxValue {
		return nil, fmt.Errorf("ensure this value is less than or equal to %d", *f.MaxValue)
	}
	return val, nil
}

func (f *IntegerField) MaxValueOption(value int64) *IntegerField {
	f.MaxValue = &value
	return f
}

func (f *IntegerField) MinValueOption(value int64) *IntegerField {
	f.MinValue = &value
	return f
}

type FloatField struct {
	*BaseField
	MaxValue *float64
	MinValue *float64
}

func NewFloatField() *FloatField {
	return &FloatField{
		BaseField: NewBaseField(),
	}
}

func (f *FloatField) ToRepresentation(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}
	switch v := value.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	}
	return nil, fmt.Errorf("cannot convert %T to float", value)
}

func (f *FloatField) ToInternalValue(data interface{}) (interface{}, error) {
	if data == nil {
		if f.Default != nil {
			return f.Default, nil
		}
		if !f.AllowNull {
			return nil, fmt.Errorf("this field may not be null")
		}
		return nil, nil
	}
	var val float64
	switch v := data.(type) {
	case float32:
		val = float64(v)
	case float64:
		val = v
	case int:
		val = float64(v)
	case int64:
		val = float64(v)
	case uint:
		val = float64(v)
	case uint64:
		val = float64(v)
	default:
		return nil, fmt.Errorf("expected float, got %T", data)
	}
	if f.MinValue != nil && val < *f.MinValue {
		return nil, fmt.Errorf("ensure this value is greater than or equal to %f", *f.MinValue)
	}
	if f.MaxValue != nil && val > *f.MaxValue {
		return nil, fmt.Errorf("ensure this value is less than or equal to %f", *f.MaxValue)
	}
	return val, nil
}

func (f *FloatField) MaxValueOption(value float64) *FloatField {
	f.MaxValue = &value
	return f
}

func (f *FloatField) MinValueOption(value float64) *FloatField {
	f.MinValue = &value
	return f
}

type BooleanField struct {
	*BaseField
}

func NewBooleanField() *BooleanField {
	return &BooleanField{
		BaseField: NewBaseField(),
	}
}

func (f *BooleanField) ToRepresentation(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}
	if b, ok := value.(bool); ok {
		return b, nil
	}
	return nil, fmt.Errorf("cannot convert %T to boolean", value)
}

func (f *BooleanField) ToInternalValue(data interface{}) (interface{}, error) {
	if data == nil {
		if f.Default != nil {
			return f.Default, nil
		}
		if !f.AllowNull {
			return nil, fmt.Errorf("this field may not be null")
		}
		return nil, nil
	}
	if b, ok := data.(bool); ok {
		return b, nil
	}
	return nil, fmt.Errorf("expected boolean, got %T", data)
}

type DateTimeField struct {
	*BaseField
	Format string
}

func NewDateTimeField() *DateTimeField {
	return &DateTimeField{
		BaseField: NewBaseField(),
		Format:    time.RFC3339,
	}
}

func (f *DateTimeField) ToRepresentation(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}
	if t, ok := value.(time.Time); ok {
		return t.Format(f.Format), nil
	}
	if t, ok := value.(*time.Time); ok && t != nil {
		return t.Format(f.Format), nil
	}
	return nil, fmt.Errorf("cannot convert %T to datetime", value)
}

func (f *DateTimeField) ToInternalValue(data interface{}) (interface{}, error) {
	if data == nil {
		if f.Default != nil {
			return f.Default, nil
		}
		if !f.AllowNull {
			return nil, fmt.Errorf("this field may not be null")
		}
		return nil, nil
	}
	if str, ok := data.(string); ok {
		t, err := time.Parse(f.Format, str)
		if err != nil {
			return nil, fmt.Errorf("invalid datetime format: %w", err)
		}
		return t, nil
	}
	return nil, fmt.Errorf("expected string, got %T", data)
}

func (f *DateTimeField) FormatOption(format string) *DateTimeField {
	f.Format = format
	return f
}

type SerializerMethodField struct {
	*BaseField
	Method string
}

func NewSerializerMethodField(method string) *SerializerMethodField {
	return &SerializerMethodField{
		BaseField: NewBaseField(),
		Method:    method,
	}
}

func (f *SerializerMethodField) ToRepresentation(value interface{}) (interface{}, error) {
	return value, nil
}

func (f *SerializerMethodField) ToInternalValue(data interface{}) (interface{}, error) {
	return nil, fmt.Errorf("SerializerMethodField is read-only")
}

type NestedField struct {
	*BaseField
	serializerFunc func(interface{}) (interface{}, error)
	fromCreateFunc func([]byte) (interface{}, error)
}

func NewNestedField[T any](serializer Serializer[T]) *NestedField {
	return &NestedField{
		BaseField: NewBaseField(),
		serializerFunc: func(obj interface{}) (interface{}, error) {
			if obj == nil {
				return nil, nil
			}
			if typedObj, ok := obj.(*T); ok {
				return serializer.ToRepresentation(typedObj)
			}
			return nil, fmt.Errorf("type mismatch")
		},
		fromCreateFunc: func(body []byte) (interface{}, error) {
			result, err := serializer.FromCreate(body)
			if err != nil {
				return nil, err
			}
			return result, nil
		},
	}
}

func (f *NestedField) ToRepresentation(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}
	return f.serializerFunc(value)
}

func (f *NestedField) ToInternalValue(data interface{}) (interface{}, error) {
	if data == nil {
		if f.Default != nil {
			return f.Default, nil
		}
		if !f.AllowNull {
			return nil, fmt.Errorf("this field may not be null")
		}
		return nil, nil
	}
	if dataMap, ok := data.(map[string]interface{}); ok {
		jsonData, err := json.Marshal(dataMap)
		if err != nil {
			return nil, err
		}
		return f.fromCreateFunc(jsonData)
	}
	return nil, fmt.Errorf("expected object, got %T", data)
}

type ValidatorFunc func(value interface{}) error

func (f ValidatorFunc) Validate(value interface{}) error {
	return f(value)
}

