package schema

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// Type-specific builders that wrap UnifiedFieldBuilder for type safety
// These eliminate duplication by delegating to
// the unified builder while providing type-safe APIs.
// ============================================================================

// Int64FieldBuilder provides type-safe methods for Int64 fields
type Int64FieldBuilder struct {
	CommonMethods[*Int64FieldBuilder]
}

// Int64 creates a new Int64 field builder
func Int64(name string) *Int64FieldBuilder {
	b := &Int64FieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeInt64)
	initCommonMethods(&b.CommonMethods, ufb, func() *Int64FieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *Int64FieldBuilder) Default(value int64) *Int64FieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *Int64FieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// Int32FieldBuilder provides type-safe methods for Int32 fields
type Int32FieldBuilder struct {
	CommonMethods[*Int32FieldBuilder]
}

// Int32 creates a new Int32 field builder
func Int32(name string) *Int32FieldBuilder {
	b := &Int32FieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeInt32)
	initCommonMethods(&b.CommonMethods, ufb, func() *Int32FieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *Int32FieldBuilder) Default(value int32) *Int32FieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *Int32FieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// StringFieldBuilder provides type-safe methods for String fields
type StringFieldBuilder struct {
	CommonMethods[*StringFieldBuilder]
}

// String creates a new String field builder
func String(name string) *StringFieldBuilder {
	b := &StringFieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeString)
	initCommonMethods(&b.CommonMethods, ufb, func() *StringFieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *StringFieldBuilder) Default(value string) *StringFieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *StringFieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// Text creates a new Text field builder (same as String but TypeText)
func Text(name string) *StringFieldBuilder {
	b := &StringFieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeText)
	initCommonMethods(&b.CommonMethods, ufb, func() *StringFieldBuilder { return b })
	return b
}

// EmailFieldBuilder provides type-safe methods for Email fields
type EmailFieldBuilder struct {
	CommonMethods[*EmailFieldBuilder]
}

// Email creates a new Email field builder
func Email(name string) *EmailFieldBuilder {
	b := &EmailFieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeEmail)
	initCommonMethods(&b.CommonMethods, ufb, func() *EmailFieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *EmailFieldBuilder) Default(value string) *EmailFieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *EmailFieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// URLFieldBuilder provides type-safe methods for URL fields
type URLFieldBuilder struct {
	CommonMethods[*URLFieldBuilder]
}

// URL creates a new URL field builder
func URL(name string) *URLFieldBuilder {
	b := &URLFieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeURL)
	initCommonMethods(&b.CommonMethods, ufb, func() *URLFieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *URLFieldBuilder) Default(value string) *URLFieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *URLFieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// BoolFieldBuilder provides type-safe methods for Bool fields
type BoolFieldBuilder struct {
	CommonMethods[*BoolFieldBuilder]
}

// Bool creates a new Bool field builder
func Bool(name string) *BoolFieldBuilder {
	b := &BoolFieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeBool)
	initCommonMethods(&b.CommonMethods, ufb, func() *BoolFieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *BoolFieldBuilder) Default(value bool) *BoolFieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *BoolFieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// TimeFieldBuilder provides type-safe methods for Time fields
type TimeFieldBuilder struct {
	CommonMethods[*TimeFieldBuilder]
}

// Time creates a new Time field builder
func Time(name string) *TimeFieldBuilder {
	b := &TimeFieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeTime)
	initCommonMethods(&b.CommonMethods, ufb, func() *TimeFieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *TimeFieldBuilder) Default(value time.Time) *TimeFieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *TimeFieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// DateFieldBuilder provides type-safe methods for Date fields
type DateFieldBuilder struct {
	CommonMethods[*DateFieldBuilder]
}

// Date creates a new Date field builder
func Date(name string) *DateFieldBuilder {
	b := &DateFieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeDate)
	initCommonMethods(&b.CommonMethods, ufb, func() *DateFieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *DateFieldBuilder) Default(value time.Time) *DateFieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *DateFieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// DateTimeFieldBuilder provides type-safe methods for DateTime fields
type DateTimeFieldBuilder struct {
	CommonMethods[*DateTimeFieldBuilder]
}

// DateTime creates a new DateTime field builder
func DateTime(name string) *DateTimeFieldBuilder {
	b := &DateTimeFieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeDateTime)
	initCommonMethods(&b.CommonMethods, ufb, func() *DateTimeFieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *DateTimeFieldBuilder) Default(value time.Time) *DateTimeFieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *DateTimeFieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// Float64FieldBuilder provides type-safe methods for Float64 fields
type Float64FieldBuilder struct {
	CommonMethods[*Float64FieldBuilder]
}

// Float64 creates a new Float64 field builder
func Float64(name string) *Float64FieldBuilder {
	b := &Float64FieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeFloat64)
	initCommonMethods(&b.CommonMethods, ufb, func() *Float64FieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *Float64FieldBuilder) Default(value float64) *Float64FieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *Float64FieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// Float32FieldBuilder provides type-safe methods for Float32 fields
type Float32FieldBuilder struct {
	CommonMethods[*Float32FieldBuilder]
}

// Float32 creates a new Float32 field builder
func Float32(name string) *Float32FieldBuilder {
	b := &Float32FieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeFloat32)
	initCommonMethods(&b.CommonMethods, ufb, func() *Float32FieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *Float32FieldBuilder) Default(value float32) *Float32FieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *Float32FieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// DecimalFieldBuilder provides type-safe methods for Decimal fields
type DecimalFieldBuilder struct {
	CommonMethods[*DecimalFieldBuilder]
}

// Decimal creates a new Decimal field builder
func Decimal(name string) *DecimalFieldBuilder {
	b := &DecimalFieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeDecimal)
	initCommonMethods(&b.CommonMethods, ufb, func() *DecimalFieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *DecimalFieldBuilder) Default(value float64) *DecimalFieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *DecimalFieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// JSONFieldBuilder provides type-safe methods for JSON fields
type JSONFieldBuilder struct {
	CommonMethods[*JSONFieldBuilder]
}

// JSON creates a new JSON field builder
func JSON(name string) *JSONFieldBuilder {
	b := &JSONFieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeJSON)
	initCommonMethods(&b.CommonMethods, ufb, func() *JSONFieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *JSONFieldBuilder) Default(value interface{}) *JSONFieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *JSONFieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// BytesFieldBuilder provides type-safe methods for Bytes fields
type BytesFieldBuilder struct {
	CommonMethods[*BytesFieldBuilder]
}

// Bytes creates a new Bytes field builder
func Bytes(name string) *BytesFieldBuilder {
	b := &BytesFieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeBytes)
	initCommonMethods(&b.CommonMethods, ufb, func() *BytesFieldBuilder { return b })
	return b
}

// Default sets the default value for the field
func (b *BytesFieldBuilder) Default(value []byte) *BytesFieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field
func (b *BytesFieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// UUIDFieldBuilder provides type-safe methods for UUID fields
type UUIDFieldBuilder struct {
	CommonMethods[*UUIDFieldBuilder]
}

// UUID creates a new UUID field builder
func UUID(name string) *UUIDFieldBuilder {
	b := &UUIDFieldBuilder{}
	ufb := newUnifiedFieldBuilder(name, TypeUUID)
	initCommonMethods(&b.CommonMethods, ufb, func() *UUIDFieldBuilder { return b })
	return b
}

// DefaultUUID sets a default UUID value
func (b *UUIDFieldBuilder) DefaultUUID(value uuid.UUID) *UUIDFieldBuilder {
	b.setDefault(value.String())
	return b
}

// DefaultNewUUID sets default to generate new UUID
func (b *UUIDFieldBuilder) DefaultNewUUID() *UUIDFieldBuilder {
	b.setDefault(uuid.New().String())
	return b
}

// Build returns the built field
func (b *UUIDFieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}
