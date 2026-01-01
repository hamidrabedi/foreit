package schema

// CommonMethods provides chainable methods that work for all field builders
// Using generics, we implement these methods once and they work for all builder types.
// 
// Usage pattern:
//   type ConcreteBuilder struct {
//       CommonMethods[*ConcreteBuilder]
//   }
//   
//   func NewConcreteBuilder() *ConcreteBuilder {
//       b := &ConcreteBuilder{}
//       ufb := newUnifiedFieldBuilder(...)
//       initCommonMethods(&b.CommonMethods, ufb, func() *ConcreteBuilder { return b })
//       return b
//   }
//
// The chainFn is set to return the concrete builder instance.
type CommonMethods[T any] struct {
	*UnifiedFieldBuilder
	chainFn func() T // Function that returns the concrete builder for chaining
}

// initCommonMethods initializes CommonMethods with the chain function
// This should be called by each builder's constructor
func initCommonMethods[T any](cm *CommonMethods[T], ufb *UnifiedFieldBuilder, chainFn func() T) {
	cm.UnifiedFieldBuilder = ufb
	cm.chainFn = chainFn
}

// Note: The actual implementation requires each concrete builder to:
// 1. Embed CommonMethods[*ConcreteBuilder]  
// 2. Implement chain() *ConcreteBuilder { return b }
// 3. CommonMethods methods will be called on the concrete builder and return T
//    by calling the concrete builder's chain() method

// ============================================================================
// Common chainable methods available on all field types
// ============================================================================

// Required marks the field as required (NOT NULL)
func (c *CommonMethods[T]) Required() T {
	c.UnifiedFieldBuilder.Required()
	return c.chainFn()
}

// Optional marks the field as optional (allows NULL)
func (c *CommonMethods[T]) Optional() T {
	c.UnifiedFieldBuilder.Optional()
	return c.chainFn()
}

// Unique marks the field as unique
func (c *CommonMethods[T]) Unique() T {
	c.UnifiedFieldBuilder.Unique()
	return c.chainFn()
}

// Primary marks the field as primary key
func (c *CommonMethods[T]) Primary() T {
	c.UnifiedFieldBuilder.Primary()
	return c.chainFn()
}

// DBIndex creates a database index on this field
func (c *CommonMethods[T]) DBIndex() T {
	c.UnifiedFieldBuilder.DBIndex()
	return c.chainFn()
}

// VerboseName sets the human-readable name for the field
func (c *CommonMethods[T]) VerboseName(name string) T {
	c.UnifiedFieldBuilder.VerboseName(name)
	return c.chainFn()
}

// HelpText sets the help text for the field
func (c *CommonMethods[T]) HelpText(text string) T {
	c.UnifiedFieldBuilder.HelpText(text)
	return c.chainFn()
}

// DBColumn sets a custom database column name
func (c *CommonMethods[T]) DBColumn(name string) T {
	c.UnifiedFieldBuilder.DBColumn(name)
	return c.chainFn()
}

// DBType sets an explicit database column type override
func (c *CommonMethods[T]) DBType(dbType string) T {
	c.UnifiedFieldBuilder.DBType(dbType)
	return c.chainFn()
}

// DBCollation sets the character set collation
func (c *CommonMethods[T]) DBCollation(collation string) T {
	c.UnifiedFieldBuilder.DBCollation(collation)
	return c.chainFn()
}

// DBComment sets the column comment/description
func (c *CommonMethods[T]) DBComment(comment string) T {
	c.UnifiedFieldBuilder.DBComment(comment)
	return c.chainFn()
}

// DBTablespace sets the tablespace for field storage
func (c *CommonMethods[T]) DBTablespace(tablespace string) T {
	c.UnifiedFieldBuilder.DBTablespace(tablespace)
	return c.chainFn()
}

// DBDefault sets a database-level default (SQL expression)
func (c *CommonMethods[T]) DBDefault(expr string) T {
	c.UnifiedFieldBuilder.DBDefault(expr)
	return c.chainFn()
}

// GeneratedColumn marks the field as a generated column
func (c *CommonMethods[T]) GeneratedColumn(expression string, stored bool) T {
	c.UnifiedFieldBuilder.GeneratedColumn(expression, stored)
	return c.chainFn()
}

// Blank allows the field to be blank
func (c *CommonMethods[T]) Blank() T {
	c.UnifiedFieldBuilder.Blank()
	return c.chainFn()
}

// Editable controls whether the field is editable in forms/admin
func (c *CommonMethods[T]) Editable(editable bool) T {
	c.UnifiedFieldBuilder.Editable(editable)
	return c.chainFn()
}

// Serialize controls whether the field is serialized in forms/API
func (c *CommonMethods[T]) Serialize(serialize bool) T {
	c.UnifiedFieldBuilder.Serialize(serialize)
	return c.chainFn()
}

// Validators adds one or more validators to the field
func (c *CommonMethods[T]) Validators(validators ...Validator) T {
	c.UnifiedFieldBuilder.Validators(validators...)
	return c.chainFn()
}

// ValidationTag sets a validation tag
func (c *CommonMethods[T]) ValidationTag(tag string) T {
	c.UnifiedFieldBuilder.ValidationTag(tag)
	return c.chainFn()
}

// ErrorMessages sets custom error messages per validation type
func (c *CommonMethods[T]) ErrorMessages(messages map[string]string) T {
	c.UnifiedFieldBuilder.ErrorMessages(messages)
	return c.chainFn()
}

// UniqueForDate makes the field unique for each date value of the specified date field
func (c *CommonMethods[T]) UniqueForDate(dateField string) T {
	c.UnifiedFieldBuilder.UniqueForDate(dateField)
	return c.chainFn()
}

// UniqueForMonth makes the field unique for each month value of the specified date field
func (c *CommonMethods[T]) UniqueForMonth(dateField string) T {
	c.UnifiedFieldBuilder.UniqueForMonth(dateField)
	return c.chainFn()
}

// UniqueForYear makes the field unique for each year value of the specified date field
func (c *CommonMethods[T]) UniqueForYear(dateField string) T {
	c.UnifiedFieldBuilder.UniqueForYear(dateField)
	return c.chainFn()
}

// ============================================================================
// Type-specific methods (available on compatible field types)
// ============================================================================

// MaxLength sets the maximum length constraint (for strings, bytes, arrays)
func (c *CommonMethods[T]) MaxLength(n int) T {
	c.UnifiedFieldBuilder.MaxLength(n)
	return c.chainFn()
}

// MinLength sets the minimum length constraint (for strings, bytes, arrays)
func (c *CommonMethods[T]) MinLength(n int) T {
	c.UnifiedFieldBuilder.MinLength(n)
	return c.chainFn()
}

// MaxValue sets the maximum value constraint (for numeric types)
func (c *CommonMethods[T]) MaxValue(val float64) T {
	c.UnifiedFieldBuilder.MaxValue(val)
	return c.chainFn()
}

// MinValue sets the minimum value constraint (for numeric types)
func (c *CommonMethods[T]) MinValue(val float64) T {
	c.UnifiedFieldBuilder.MinValue(val)
	return c.chainFn()
}

// MaxDigits sets the maximum number of digits (for Decimal fields)
func (c *CommonMethods[T]) MaxDigits(digits int) T {
	c.UnifiedFieldBuilder.MaxDigits(digits)
	return c.chainFn()
}

// DecimalPlaces sets the number of decimal places (for Decimal fields)
func (c *CommonMethods[T]) DecimalPlaces(places int) T {
	c.UnifiedFieldBuilder.DecimalPlaces(places)
	return c.chainFn()
}

// Choices sets predefined choices for the field
func (c *CommonMethods[T]) Choices(choices ...Choice) T {
	c.UnifiedFieldBuilder.Choices(choices...)
	return c.chainFn()
}

// ChoicesFromPairs creates choices from value-label pairs
func (c *CommonMethods[T]) ChoicesFromPairs(pairs ...string) T {
	c.UnifiedFieldBuilder.ChoicesFromPairs(pairs...)
	return c.chainFn()
}

// AutoIncrement marks the field as auto-incrementing (for integer fields)
func (c *CommonMethods[T]) AutoIncrement() T {
	c.UnifiedFieldBuilder.AutoIncrement()
	return c.chainFn()
}

// AutoNow sets the field to automatically update on save (for temporal fields)
func (c *CommonMethods[T]) AutoNow() T {
	c.UnifiedFieldBuilder.AutoNow()
	return c.chainFn()
}

// AutoNowAdd sets the field to automatically set on creation only (for temporal fields)
func (c *CommonMethods[T]) AutoNowAdd() T {
	c.UnifiedFieldBuilder.AutoNowAdd()
	return c.chainFn()
}

// WriteOnly marks the field as write-only (not serialized in API responses)
func (c *CommonMethods[T]) WriteOnly() T {
	c.UnifiedFieldBuilder.WriteOnly()
	return c.chainFn()
}
