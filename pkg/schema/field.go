package schema

import "time"

// FieldType represents the type of a field
type FieldType int

const (
	TypeInt64 FieldType = iota
	TypeInt32
	TypeString
	TypeBool
	TypeTime
	TypeDate
	TypeDateTime
	TypeFloat32
	TypeFloat64
	TypeDecimal
	TypeText
	TypeEmail
	TypeURL
	TypeUUID
	TypeJSON
	TypeBytes
	TypeForeignKey
	TypeManyToMany
	TypeOneToOne
)

// Field represents a model field with ALL Django field options
type Field struct {
	CustomField    CustomField
	Default        interface{}
	DBDefault      string   // Database-level default (SQL expression)
	MaxLength      *int
	MinValue       *float64
	MaxValue       *float64
	MinLength      *int
	MaxDigits      *int      // For Decimal fields
	DecimalPlaces  *int      // For Decimal fields
	HelpText       string
	VerboseName    string
	Name           string
	DBColumn       string
	DBType         string    // Explicit database column type override
	DBCollation    string    // Character set collation (PostgreSQL, MySQL)
	DBComment      string    // Column comment/description
	DBTablespace   string    // Tablespace for field storage (PostgreSQL)
	ValidationTag  string
	Choices        []Choice
	Validators     []Validator
	ErrorMessages  map[string]string // Custom error messages per validation type
	Type           FieldType
	AutoIncrement  bool
	PrimaryKey     bool
	Blank          bool
	Required       bool
	Unique         bool
	Editable       bool
	Serialize      bool      // Control field serialization in forms/API
	AutoNow        bool
	AutoNowAdd     bool
	DBIndex        bool
	UniqueForDate  string    // Unique constraint scoped to date field
	UniqueForMonth string    // Unique constraint scoped to month
	UniqueForYear  string    // Unique constraint scoped to year
	Generated      bool      // Generated column (STORED/VIRTUAL)
	GeneratedExpr  string    // Expression for generated column
	IsStored       bool      // For generated columns: STORED vs VIRTUAL
}

// Choice represents a choice option for ChoiceField
type Choice struct {
	Value string
	Label string
}

// Validator is the interface for field validators
type Validator interface {
	Validate(value interface{}) error
}

// CustomField is the interface for custom field types
type CustomField interface {
	GetName() string
	GetType() FieldType
	Validate(value interface{}) error
}

// FieldBuilder provides a fluent interface for building fields
type FieldBuilder struct {
	field Field
}

// Int64 creates a new Int64 field builder
func Int64(name string) *Int64FieldBuilder {
	return &Int64FieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeInt64,
			},
		},
	}
}

// Int64FieldBuilder is a builder for Int64 fields
type Int64FieldBuilder struct {
	*BaseFieldBuilder
}

// AutoIncrement marks the field as auto-incrementing
func (b *Int64FieldBuilder) AutoIncrement() *Int64FieldBuilder {
	b.field.AutoIncrement = true
	return b
}

// Default sets the default value for the field
func (b *Int64FieldBuilder) Default(value int64) *Int64FieldBuilder {
	b.field.Default = value
	return b
}

// MaxValue sets the maximum value constraint
func (b *Int64FieldBuilder) MaxValue(val float64) *Int64FieldBuilder {
	b.field.MaxValue = &val
	return b
}

// MinValue sets the minimum value constraint
func (b *Int64FieldBuilder) MinValue(val float64) *Int64FieldBuilder {
	b.field.MinValue = &val
	return b
}

// Build returns the built field
func (b *Int64FieldBuilder) Build() Field {
	return b.field
}

// String creates a new String field builder
func String(name string) *StringFieldBuilder {
	return &StringFieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeString,
			},
		},
	}
}

// StringFieldBuilder is a builder for String fields
type StringFieldBuilder struct {
	*BaseFieldBuilder
}

// WriteOnly marks the field as write-only (not serialized in API responses)
func (b *StringFieldBuilder) WriteOnly() *StringFieldBuilder {
	b.field.Editable = true
	return b
}

// Default sets the default value for the field
func (b *StringFieldBuilder) Default(value string) *StringFieldBuilder {
	b.field.Default = value
	return b
}

// MaxLength sets the maximum length constraint
func (b *StringFieldBuilder) MaxLength(n int) *StringFieldBuilder {
	b.field.MaxLength = &n
	return b
}

// MinLength sets the minimum length constraint
func (b *StringFieldBuilder) MinLength(n int) *StringFieldBuilder {
	b.field.MinLength = &n
	return b
}

// Choices sets predefined choices for the field
func (b *StringFieldBuilder) Choices(choices ...Choice) *StringFieldBuilder {
	b.field.Choices = choices
	return b
}

// ChoicesFromPairs is a convenience method that creates choices from value-label pairs
// Usage: ChoicesFromPairs("active", "Active", "inactive", "Inactive")
func (b *StringFieldBuilder) ChoicesFromPairs(pairs ...string) *StringFieldBuilder {
	b.field.Choices = Choices(pairs...)
	return b
}

// Build returns the built field
func (b *StringFieldBuilder) Build() Field {
	return b.field
}

// Bool creates a new Bool field builder
func Bool(name string) *BoolFieldBuilder {
	return &BoolFieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeBool,
			},
		},
	}
}

// BoolFieldBuilder is a builder for Bool fields
type BoolFieldBuilder struct {
	*BaseFieldBuilder
}

// Default sets the default value for the field
func (b *BoolFieldBuilder) Default(value bool) *BoolFieldBuilder {
	b.field.Default = value
	return b
}

// Build returns the built field
func (b *BoolFieldBuilder) Build() Field {
	return b.field
}

// Time creates a new Time field builder
func Time(name string) *TimeFieldBuilder {
	return &TimeFieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeTime,
			},
		},
	}
}

// TimeFieldBuilder is a builder for Time fields
type TimeFieldBuilder struct {
	*BaseFieldBuilder
}

// AutoNow sets the field to automatically update on save
func (b *TimeFieldBuilder) AutoNow() *TimeFieldBuilder {
	b.field.AutoNow = true
	return b
}

// AutoNowAdd sets the field to automatically set on creation only
func (b *TimeFieldBuilder) AutoNowAdd() *TimeFieldBuilder {
	b.field.AutoNowAdd = true
	return b
}

// Default sets the default value for the field
func (b *TimeFieldBuilder) Default(value time.Time) *TimeFieldBuilder {
	b.field.Default = value
	return b
}

// Build returns the built field
func (b *TimeFieldBuilder) Build() Field {
	return b.field
}

// Text creates a new Text field builder (same as String but TypeText)
func Text(name string) *StringFieldBuilder {
	return &StringFieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeText,
			},
		},
	}
}

// Int creates a new Int field builder (alias for Int64)
func Int(name string) *Int64FieldBuilder {
	return Int64(name)
}

// Email creates a new Email field builder
func Email(name string) *EmailFieldBuilder {
	return &EmailFieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeEmail,
			},
		},
	}
}

// EmailFieldBuilder is a builder for Email fields
type EmailFieldBuilder struct {
	*BaseFieldBuilder
}

// Default sets the default value for the field
func (b *EmailFieldBuilder) Default(value string) *EmailFieldBuilder {
	b.field.Default = value
	return b
}

// MaxLength sets the maximum length constraint
func (b *EmailFieldBuilder) MaxLength(n int) *EmailFieldBuilder {
	b.field.MaxLength = &n
	return b
}

// Build returns the built field
func (b *EmailFieldBuilder) Build() Field {
	return b.field
}

// Float64 creates a new Float64 field builder
func Float64(name string) *Float64FieldBuilder {
	return &Float64FieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeFloat64,
			},
		},
	}
}

// Float64FieldBuilder is a builder for Float64 fields
type Float64FieldBuilder struct {
	*BaseFieldBuilder
}

// Default sets the default value for the field
func (b *Float64FieldBuilder) Default(value float64) *Float64FieldBuilder {
	b.field.Default = value
	return b
}

// MaxValue sets the maximum value constraint
func (b *Float64FieldBuilder) MaxValue(val float64) *Float64FieldBuilder {
	b.field.MaxValue = &val
	return b
}

// MinValue sets the minimum value constraint
func (b *Float64FieldBuilder) MinValue(val float64) *Float64FieldBuilder {
	b.field.MinValue = &val
	return b
}

// Build returns the built field
func (b *Float64FieldBuilder) Build() Field {
	return b.field
}

// Decimal creates a new Decimal field builder
func Decimal(name string) *DecimalFieldBuilder {
	return &DecimalFieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeDecimal,
			},
		},
	}
}

// DecimalFieldBuilder is a builder for Decimal fields
type DecimalFieldBuilder struct {
	*BaseFieldBuilder
	maxDigits     int
	decimalPlaces int
}

// MaxDigits sets the maximum number of digits
func (b *DecimalFieldBuilder) MaxDigits(digits int) *DecimalFieldBuilder {
	b.maxDigits = digits
	return b
}

// DecimalPlaces sets the number of decimal places
func (b *DecimalFieldBuilder) DecimalPlaces(places int) *DecimalFieldBuilder {
	b.decimalPlaces = places
	return b
}

// Default sets the default value for the field
func (b *DecimalFieldBuilder) Default(value float64) *DecimalFieldBuilder {
	b.field.Default = value
	return b
}

// MaxValue sets the maximum value constraint
func (b *DecimalFieldBuilder) MaxValue(val float64) *DecimalFieldBuilder {
	b.field.MaxValue = &val
	return b
}

// MinValue sets the minimum value constraint
func (b *DecimalFieldBuilder) MinValue(val float64) *DecimalFieldBuilder {
	b.field.MinValue = &val
	return b
}

// Build returns the built field
func (b *DecimalFieldBuilder) Build() Field {
	if b.maxDigits > 0 {
		b.field.MaxDigits = &b.maxDigits
	}
	if b.decimalPlaces > 0 {
		b.field.DecimalPlaces = &b.decimalPlaces
	}
	return b.field
}

// JSON creates a new JSON field builder
func JSON(name string) *JSONFieldBuilder {
	return &JSONFieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeJSON,
			},
		},
	}
}

// JSONFieldBuilder is a builder for JSON fields
type JSONFieldBuilder struct {
	*BaseFieldBuilder
}

// Default sets the default value for the field
func (b *JSONFieldBuilder) Default(value interface{}) *JSONFieldBuilder {
	b.field.Default = value
	return b
}

// Build returns the built field
func (b *JSONFieldBuilder) Build() Field {
	return b.field
}

// Date creates a new Date field builder
func Date(name string) *DateFieldBuilder {
	return &DateFieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeDate,
			},
		},
	}
}

// DateFieldBuilder is a builder for Date fields
type DateFieldBuilder struct {
	*BaseFieldBuilder
}

// AutoNow sets the field to automatically update on save
func (b *DateFieldBuilder) AutoNow() *DateFieldBuilder {
	b.field.AutoNow = true
	return b
}

// AutoNowAdd sets the field to automatically set on creation only
func (b *DateFieldBuilder) AutoNowAdd() *DateFieldBuilder {
	b.field.AutoNowAdd = true
	return b
}

// Default sets the default value for the field
func (b *DateFieldBuilder) Default(value time.Time) *DateFieldBuilder {
	b.field.Default = value
	return b
}

// Build returns the built field
func (b *DateFieldBuilder) Build() Field {
	return b.field
}

// DateTime creates a new DateTime field builder
func DateTime(name string) *DateTimeFieldBuilder {
	return &DateTimeFieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeDateTime,
			},
		},
	}
}

// DateTimeFieldBuilder is a builder for DateTime fields
type DateTimeFieldBuilder struct {
	*BaseFieldBuilder
}

// AutoNow sets the field to automatically update on save
func (b *DateTimeFieldBuilder) AutoNow() *DateTimeFieldBuilder {
	b.field.AutoNow = true
	return b
}

// AutoNowAdd sets the field to automatically set on creation only
func (b *DateTimeFieldBuilder) AutoNowAdd() *DateTimeFieldBuilder {
	b.field.AutoNowAdd = true
	return b
}

// Default sets the default value for the field
func (b *DateTimeFieldBuilder) Default(value time.Time) *DateTimeFieldBuilder {
	b.field.Default = value
	return b
}

// Build returns the built field
func (b *DateTimeFieldBuilder) Build() Field {
	return b.field
}

// URL creates a new URL field builder
func URL(name string) *URLFieldBuilder {
	return &URLFieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeURL,
			},
		},
	}
}

// URLFieldBuilder is a builder for URL fields
type URLFieldBuilder struct {
	*BaseFieldBuilder
}

// Default sets the default value for the field
func (b *URLFieldBuilder) Default(value string) *URLFieldBuilder {
	b.field.Default = value
	return b
}

// MaxLength sets the maximum length constraint
func (b *URLFieldBuilder) MaxLength(n int) *URLFieldBuilder {
	b.field.MaxLength = &n
	return b
}

// Build returns the built field
func (b *URLFieldBuilder) Build() Field {
	return b.field
}

// Int32 creates a new Int32 field builder
func Int32(name string) *Int32FieldBuilder {
	return &Int32FieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeInt32,
			},
		},
	}
}

// Int32FieldBuilder is a builder for Int32 fields
type Int32FieldBuilder struct {
	*BaseFieldBuilder
}

// AutoIncrement marks the field as auto-incrementing
func (b *Int32FieldBuilder) AutoIncrement() *Int32FieldBuilder {
	b.field.AutoIncrement = true
	return b
}

// Default sets the default value for the field
func (b *Int32FieldBuilder) Default(value int32) *Int32FieldBuilder {
	b.field.Default = value
	return b
}

// MaxValue sets the maximum value constraint
func (b *Int32FieldBuilder) MaxValue(val float64) *Int32FieldBuilder {
	b.field.MaxValue = &val
	return b
}

// MinValue sets the minimum value constraint
func (b *Int32FieldBuilder) MinValue(val float64) *Int32FieldBuilder {
	b.field.MinValue = &val
	return b
}

// Build returns the built field
func (b *Int32FieldBuilder) Build() Field {
	return b.field
}

// Float32 creates a new Float32 field builder
func Float32(name string) *Float32FieldBuilder {
	return &Float32FieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeFloat32,
			},
		},
	}
}

// Float32FieldBuilder is a builder for Float32 fields
type Float32FieldBuilder struct {
	*BaseFieldBuilder
}

// Default sets the default value for the field
func (b *Float32FieldBuilder) Default(value float32) *Float32FieldBuilder {
	b.field.Default = value
	return b
}

// MaxValue sets the maximum value constraint
func (b *Float32FieldBuilder) MaxValue(val float64) *Float32FieldBuilder {
	b.field.MaxValue = &val
	return b
}

// MinValue sets the minimum value constraint
func (b *Float32FieldBuilder) MinValue(val float64) *Float32FieldBuilder {
	b.field.MinValue = &val
	return b
}

// Build returns the built field
func (b *Float32FieldBuilder) Build() Field {
	return b.field
}

// Bytes creates a new Bytes field builder
func Bytes(name string) *BytesFieldBuilder {
	return &BytesFieldBuilder{
		BaseFieldBuilder: &BaseFieldBuilder{
			field: Field{
				Name: name,
				Type: TypeBytes,
			},
		},
	}
}

// BytesFieldBuilder is a builder for Bytes fields
type BytesFieldBuilder struct {
	*BaseFieldBuilder
}

// Default sets the default value for the field
func (b *BytesFieldBuilder) Default(value []byte) *BytesFieldBuilder {
	b.field.Default = value
	return b
}

// MaxLength sets the maximum length constraint
func (b *BytesFieldBuilder) MaxLength(n int) *BytesFieldBuilder {
	b.field.MaxLength = &n
	return b
}

// Build returns the built field
func (b *BytesFieldBuilder) Build() Field {
	return b.field
}

