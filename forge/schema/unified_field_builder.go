package schema

import "time"

// UnifiedFieldBuilder is a single, unified builder for all field types
// This eliminates duplication by using composition and type-specific method sets
type UnifiedFieldBuilder struct {
	field   Field
	options *FieldOptions
}

// newUnifiedFieldBuilder creates a new unified field builder
func newUnifiedFieldBuilder(name string, fieldType FieldType) *UnifiedFieldBuilder {
	return &UnifiedFieldBuilder{
		field: Field{
			Name: name,
			Type: fieldType,
		},
		options: NewFieldOptions(),
	}
}

// Build returns the built field with all options applied
func (b *UnifiedFieldBuilder) Build() Field {
	b.options.ApplyToField(&b.field)
	return b.field
}

// ============================================================================
// Common methods available on all field types (via BaseFieldBuilder pattern)
// ============================================================================

// DBColumn sets a custom database column name
func (b *UnifiedFieldBuilder) DBColumn(name string) *UnifiedFieldBuilder {
	b.options.DB.Column = name
	return b
}

// DBType sets an explicit database column type override
func (b *UnifiedFieldBuilder) DBType(dbType string) *UnifiedFieldBuilder {
	b.options.DB.Type = dbType
	return b
}

// DBCollation sets the character set collation
func (b *UnifiedFieldBuilder) DBCollation(collation string) *UnifiedFieldBuilder {
	b.options.DB.Collation = collation
	return b
}

// DBComment sets the column comment/description
func (b *UnifiedFieldBuilder) DBComment(comment string) *UnifiedFieldBuilder {
	b.options.DB.Comment = comment
	return b
}

// DBTablespace sets the tablespace for field storage
func (b *UnifiedFieldBuilder) DBTablespace(tablespace string) *UnifiedFieldBuilder {
	b.options.DB.Tablespace = tablespace
	return b
}

// DBDefault sets a database-level default (SQL expression)
func (b *UnifiedFieldBuilder) DBDefault(expr string) *UnifiedFieldBuilder {
	b.options.DB.Default = expr
	return b
}

// DBIndex creates a database index on this field
func (b *UnifiedFieldBuilder) DBIndex() *UnifiedFieldBuilder {
	b.options.DB.Index = true
	return b
}

// GeneratedColumn marks the field as a generated column
func (b *UnifiedFieldBuilder) GeneratedColumn(expression string, stored bool) *UnifiedFieldBuilder {
	b.options.DB.Generated = true
	b.options.DB.GeneratedExpr = expression
	b.options.DB.IsStored = stored
	return b
}

// Required marks the field as required (NOT NULL)
func (b *UnifiedFieldBuilder) Required() *UnifiedFieldBuilder {
	b.options.Validation.Required = true
	return b
}

// Optional marks the field as optional (allows NULL)
func (b *UnifiedFieldBuilder) Optional() *UnifiedFieldBuilder {
	b.options.Validation.Required = false
	return b
}

// Unique marks the field as unique
func (b *UnifiedFieldBuilder) Unique() *UnifiedFieldBuilder {
	b.options.Validation.Unique = true
	return b
}

// Primary marks the field as primary key
func (b *UnifiedFieldBuilder) Primary() *UnifiedFieldBuilder {
	b.field.PrimaryKey = true
	return b
}

// Blank allows the field to be blank
func (b *UnifiedFieldBuilder) Blank() *UnifiedFieldBuilder {
	b.options.Validation.Blank = true
	return b
}

// Editable controls whether the field is editable in forms/admin
func (b *UnifiedFieldBuilder) Editable(editable bool) *UnifiedFieldBuilder {
	b.options.Presentation.Editable = editable
	return b
}

// Serialize controls whether the field is serialized in forms/API
func (b *UnifiedFieldBuilder) Serialize(serialize bool) *UnifiedFieldBuilder {
	b.options.Presentation.Serialize = serialize
	return b
}

// Validators adds one or more validators to the field
func (b *UnifiedFieldBuilder) Validators(validators ...Validator) *UnifiedFieldBuilder {
	b.options.Validation.Validators = append(b.options.Validation.Validators, validators...)
	return b
}

// ValidationTag sets a validation tag
func (b *UnifiedFieldBuilder) ValidationTag(tag string) *UnifiedFieldBuilder {
	b.options.Validation.ValidationTag = tag
	return b
}

// ErrorMessages sets custom error messages per validation type
func (b *UnifiedFieldBuilder) ErrorMessages(messages map[string]string) *UnifiedFieldBuilder {
	b.options.Validation.ErrorMessages = messages
	return b
}

// HelpText sets the help text for the field
func (b *UnifiedFieldBuilder) HelpText(text string) *UnifiedFieldBuilder {
	b.options.Presentation.HelpText = text
	return b
}

// VerboseName sets the human-readable name for the field
func (b *UnifiedFieldBuilder) VerboseName(name string) *UnifiedFieldBuilder {
	b.options.Presentation.VerboseName = name
	return b
}

// UniqueForDate makes the field unique for each date value of the specified date field
func (b *UnifiedFieldBuilder) UniqueForDate(dateField string) *UnifiedFieldBuilder {
	b.options.Validation.UniqueForDate = dateField
	return b
}

// UniqueForMonth makes the field unique for each month value of the specified date field
func (b *UnifiedFieldBuilder) UniqueForMonth(dateField string) *UnifiedFieldBuilder {
	b.options.Validation.UniqueForMonth = dateField
	return b
}

// UniqueForYear makes the field unique for each year value of the specified date field
func (b *UnifiedFieldBuilder) UniqueForYear(dateField string) *UnifiedFieldBuilder {
	b.options.Validation.UniqueForYear = dateField
	return b
}

// ============================================================================
// Type-specific methods (only available on compatible field types)
// ============================================================================

// Default sets the default value (type-safe variants below)
func (b *UnifiedFieldBuilder) setDefault(value interface{}) *UnifiedFieldBuilder {
	b.field.Default = value
	return b
}

// MaxLength sets the maximum length constraint (for strings, bytes, arrays)
func (b *UnifiedFieldBuilder) MaxLength(n int) *UnifiedFieldBuilder {
	b.options.Validation.MaxLength = &n
	return b
}

// MinLength sets the minimum length constraint (for strings, bytes, arrays)
func (b *UnifiedFieldBuilder) MinLength(n int) *UnifiedFieldBuilder {
	b.options.Validation.MinLength = &n
	return b
}

// MaxValue sets the maximum value constraint (for numeric types)
func (b *UnifiedFieldBuilder) MaxValue(val float64) *UnifiedFieldBuilder {
	b.options.Validation.MaxValue = &val
	return b
}

// MinValue sets the minimum value constraint (for numeric types)
func (b *UnifiedFieldBuilder) MinValue(val float64) *UnifiedFieldBuilder {
	b.options.Validation.MinValue = &val
	return b
}

// MaxDigits sets the maximum number of digits (for Decimal fields)
func (b *UnifiedFieldBuilder) MaxDigits(digits int) *UnifiedFieldBuilder {
	b.options.Validation.MaxDigits = &digits
	return b
}

// DecimalPlaces sets the number of decimal places (for Decimal fields)
func (b *UnifiedFieldBuilder) DecimalPlaces(places int) *UnifiedFieldBuilder {
	b.options.Validation.DecimalPlaces = &places
	return b
}

// Choices sets predefined choices for the field
func (b *UnifiedFieldBuilder) Choices(choices ...Choice) *UnifiedFieldBuilder {
	b.options.Validation.Choices = choices
	return b
}

// ChoicesFromPairs creates choices from value-label pairs
func (b *UnifiedFieldBuilder) ChoicesFromPairs(pairs ...string) *UnifiedFieldBuilder {
	b.options.Validation.Choices = Choices(pairs...)
	return b
}

// AutoIncrement marks the field as auto-incrementing (for integer fields)
func (b *UnifiedFieldBuilder) AutoIncrement() *UnifiedFieldBuilder {
	b.field.AutoIncrement = true
	return b
}

// AutoNow sets the field to automatically update on save (for temporal fields)
func (b *UnifiedFieldBuilder) AutoNow() *UnifiedFieldBuilder {
	b.field.AutoNow = true
	return b
}

// AutoNowAdd sets the field to automatically set on creation only (for temporal fields)
func (b *UnifiedFieldBuilder) AutoNowAdd() *UnifiedFieldBuilder {
	b.field.AutoNowAdd = true
	return b
}

// WriteOnly marks the field as write-only (not serialized in API responses)
func (b *UnifiedFieldBuilder) WriteOnly() *UnifiedFieldBuilder {
	b.options.Presentation.Serialize = false
	return b
}

// ============================================================================
// Type-safe Default methods for each field type
// ============================================================================

// DefaultInt64 sets the default value for Int64 fields
func (b *UnifiedFieldBuilder) DefaultInt64(value int64) *UnifiedFieldBuilder {
	return b.setDefault(value)
}

// DefaultInt32 sets the default value for Int32 fields
func (b *UnifiedFieldBuilder) DefaultInt32(value int32) *UnifiedFieldBuilder {
	return b.setDefault(value)
}

// DefaultString sets the default value for String fields
func (b *UnifiedFieldBuilder) DefaultString(value string) *UnifiedFieldBuilder {
	return b.setDefault(value)
}

// DefaultBool sets the default value for Bool fields
func (b *UnifiedFieldBuilder) DefaultBool(value bool) *UnifiedFieldBuilder {
	return b.setDefault(value)
}

// DefaultTime sets the default value for Time fields
func (b *UnifiedFieldBuilder) DefaultTime(value time.Time) *UnifiedFieldBuilder {
	return b.setDefault(value)
}

// DefaultFloat64 sets the default value for Float64 fields
func (b *UnifiedFieldBuilder) DefaultFloat64(value float64) *UnifiedFieldBuilder {
	return b.setDefault(value)
}

// DefaultFloat32 sets the default value for Float32 fields
func (b *UnifiedFieldBuilder) DefaultFloat32(value float32) *UnifiedFieldBuilder {
	return b.setDefault(value)
}

// DefaultJSON sets the default value for JSON fields
func (b *UnifiedFieldBuilder) DefaultJSON(value interface{}) *UnifiedFieldBuilder {
	return b.setDefault(value)
}

// DefaultBytes sets the default value for Bytes fields
func (b *UnifiedFieldBuilder) DefaultBytes(value []byte) *UnifiedFieldBuilder {
	return b.setDefault(value)
}
