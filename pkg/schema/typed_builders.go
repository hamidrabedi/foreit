package schema

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// Type-specific builders that wrap UnifiedFieldBuilder for type safety
// and backward compatibility. These eliminate duplication by delegating to
// the unified builder while providing type-safe APIs.
// ============================================================================

// Int64FieldBuilder provides type-safe methods for Int64 fields
type Int64FieldBuilder struct {
	*UnifiedFieldBuilder
}

// Int64 creates a new Int64 field builder
func Int64(name string) *Int64FieldBuilder {
	return &Int64FieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeInt64),
	}
}

// Default sets the default value for the field
func (b *Int64FieldBuilder) Default(value int64) *Int64FieldBuilder {
	b.setDefault(value)
	return b
}

// Int32FieldBuilder provides type-safe methods for Int32 fields
type Int32FieldBuilder struct {
	*UnifiedFieldBuilder
}

// Int32 creates a new Int32 field builder
func Int32(name string) *Int32FieldBuilder {
	return &Int32FieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeInt32),
	}
}

// Default sets the default value for the field
func (b *Int32FieldBuilder) Default(value int32) *Int32FieldBuilder {
	b.setDefault(value)
	return b
}

// StringFieldBuilder provides type-safe methods for String fields
type StringFieldBuilder struct {
	*UnifiedFieldBuilder
}

// String creates a new String field builder
func String(name string) *StringFieldBuilder {
	return &StringFieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeString),
	}
}

// Default sets the default value for the field
func (b *StringFieldBuilder) Default(value string) *StringFieldBuilder {
	b.setDefault(value)
	return b
}

// WriteOnly marks the field as write-only (not serialized in API responses)
func (b *StringFieldBuilder) WriteOnly() *StringFieldBuilder {
	b.UnifiedFieldBuilder.WriteOnly()
	return b
}

// Text creates a new Text field builder (same as String but TypeText)
func Text(name string) *StringFieldBuilder {
	builder := newUnifiedFieldBuilder(name, TypeText)
	return &StringFieldBuilder{
		UnifiedFieldBuilder: builder,
	}
}

// EmailFieldBuilder provides type-safe methods for Email fields
type EmailFieldBuilder struct {
	*UnifiedFieldBuilder
}

// Email creates a new Email field builder
func Email(name string) *EmailFieldBuilder {
	return &EmailFieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeEmail),
	}
}

// Default sets the default value for the field
func (b *EmailFieldBuilder) Default(value string) *EmailFieldBuilder {
	b.setDefault(value)
	return b
}

// URLFieldBuilder provides type-safe methods for URL fields
type URLFieldBuilder struct {
	*UnifiedFieldBuilder
}

// URL creates a new URL field builder
func URL(name string) *URLFieldBuilder {
	return &URLFieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeURL),
	}
}

// Default sets the default value for the field
func (b *URLFieldBuilder) Default(value string) *URLFieldBuilder {
	b.setDefault(value)
	return b
}

// BoolFieldBuilder provides type-safe methods for Bool fields
type BoolFieldBuilder struct {
	*UnifiedFieldBuilder
}

// Bool creates a new Bool field builder
func Bool(name string) *BoolFieldBuilder {
	return &BoolFieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeBool),
	}
}

// Default sets the default value for the field
func (b *BoolFieldBuilder) Default(value bool) *BoolFieldBuilder {
	b.setDefault(value)
	return b
}

// TimeFieldBuilder provides type-safe methods for Time fields
type TimeFieldBuilder struct {
	*UnifiedFieldBuilder
}

// Time creates a new Time field builder
func Time(name string) *TimeFieldBuilder {
	return &TimeFieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeTime),
	}
}

// Default sets the default value for the field
func (b *TimeFieldBuilder) Default(value time.Time) *TimeFieldBuilder {
	b.setDefault(value)
	return b
}

// DateFieldBuilder provides type-safe methods for Date fields
type DateFieldBuilder struct {
	*UnifiedFieldBuilder
}

// Date creates a new Date field builder
func Date(name string) *DateFieldBuilder {
	return &DateFieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeDate),
	}
}

// Default sets the default value for the field
func (b *DateFieldBuilder) Default(value time.Time) *DateFieldBuilder {
	b.setDefault(value)
	return b
}

// DateTimeFieldBuilder provides type-safe methods for DateTime fields
type DateTimeFieldBuilder struct {
	*UnifiedFieldBuilder
}

// DateTime creates a new DateTime field builder
func DateTime(name string) *DateTimeFieldBuilder {
	return &DateTimeFieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeDateTime),
	}
}

// Default sets the default value for the field
func (b *DateTimeFieldBuilder) Default(value time.Time) *DateTimeFieldBuilder {
	b.setDefault(value)
	return b
}

// Float64FieldBuilder provides type-safe methods for Float64 fields
type Float64FieldBuilder struct {
	*UnifiedFieldBuilder
}

// Float64 creates a new Float64 field builder
func Float64(name string) *Float64FieldBuilder {
	return &Float64FieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeFloat64),
	}
}

// Default sets the default value for the field
func (b *Float64FieldBuilder) Default(value float64) *Float64FieldBuilder {
	b.setDefault(value)
	return b
}

// Float32FieldBuilder provides type-safe methods for Float32 fields
type Float32FieldBuilder struct {
	*UnifiedFieldBuilder
}

// Float32 creates a new Float32 field builder
func Float32(name string) *Float32FieldBuilder {
	return &Float32FieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeFloat32),
	}
}

// Default sets the default value for the field
func (b *Float32FieldBuilder) Default(value float32) *Float32FieldBuilder {
	b.setDefault(value)
	return b
}

// DecimalFieldBuilder provides type-safe methods for Decimal fields
type DecimalFieldBuilder struct {
	*UnifiedFieldBuilder
}

// Decimal creates a new Decimal field builder
func Decimal(name string) *DecimalFieldBuilder {
	return &DecimalFieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeDecimal),
	}
}

// Default sets the default value for the field
func (b *DecimalFieldBuilder) Default(value float64) *DecimalFieldBuilder {
	b.setDefault(value)
	return b
}

// Build returns the built field (overrides to ensure MaxDigits/DecimalPlaces are applied)
func (b *DecimalFieldBuilder) Build() Field {
	return b.UnifiedFieldBuilder.Build()
}

// JSONFieldBuilder provides type-safe methods for JSON fields
type JSONFieldBuilder struct {
	*UnifiedFieldBuilder
}

// JSON creates a new JSON field builder
func JSON(name string) *JSONFieldBuilder {
	return &JSONFieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeJSON),
	}
}

// Default sets the default value for the field
func (b *JSONFieldBuilder) Default(value interface{}) *JSONFieldBuilder {
	b.setDefault(value)
	return b
}

// BytesFieldBuilder provides type-safe methods for Bytes fields
type BytesFieldBuilder struct {
	*UnifiedFieldBuilder
}

// Bytes creates a new Bytes field builder
func Bytes(name string) *BytesFieldBuilder {
	return &BytesFieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeBytes),
	}
}

// Default sets the default value for the field
func (b *BytesFieldBuilder) Default(value []byte) *BytesFieldBuilder {
	b.setDefault(value)
	return b
}

// UUIDFieldBuilder provides type-safe methods for UUID fields
type UUIDFieldBuilder struct {
	*UnifiedFieldBuilder
}

// UUID creates a new UUID field builder
func UUID(name string) *UUIDFieldBuilder {
	return &UUIDFieldBuilder{
		UnifiedFieldBuilder: newUnifiedFieldBuilder(name, TypeUUID),
	}
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

// Int is an alias for Int64 (backward compatibility)
func Int(name string) *Int64FieldBuilder {
	return Int64(name)
}
