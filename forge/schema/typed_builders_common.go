package schema

// This file contains common method overrides for type-specific builders
// to ensure proper method chaining that returns the correct builder type.

// ============================================================================
// Int64FieldBuilder common method overrides
// ============================================================================

func (b *Int64FieldBuilder) Required() *Int64FieldBuilder {
	b.UnifiedFieldBuilder.Required()
	return b
}

func (b *Int64FieldBuilder) Optional() *Int64FieldBuilder {
	b.UnifiedFieldBuilder.Optional()
	return b
}

func (b *Int64FieldBuilder) Unique() *Int64FieldBuilder {
	b.UnifiedFieldBuilder.Unique()
	return b
}

func (b *Int64FieldBuilder) Primary() *Int64FieldBuilder {
	b.UnifiedFieldBuilder.Primary()
	return b
}

func (b *Int64FieldBuilder) DBIndex() *Int64FieldBuilder {
	b.UnifiedFieldBuilder.DBIndex()
	return b
}

func (b *Int64FieldBuilder) AutoIncrement() *Int64FieldBuilder {
	b.UnifiedFieldBuilder.AutoIncrement()
	return b
}

func (b *Int64FieldBuilder) MaxValue(val float64) *Int64FieldBuilder {
	b.UnifiedFieldBuilder.MaxValue(val)
	return b
}

func (b *Int64FieldBuilder) MinValue(val float64) *Int64FieldBuilder {
	b.UnifiedFieldBuilder.MinValue(val)
	return b
}

func (b *Int64FieldBuilder) VerboseName(name string) *Int64FieldBuilder {
	b.UnifiedFieldBuilder.VerboseName(name)
	return b
}

func (b *Int64FieldBuilder) HelpText(text string) *Int64FieldBuilder {
	b.UnifiedFieldBuilder.HelpText(text)
	return b
}

// ============================================================================
// StringFieldBuilder common method overrides
// ============================================================================

func (b *StringFieldBuilder) Required() *StringFieldBuilder {
	b.UnifiedFieldBuilder.Required()
	return b
}

func (b *StringFieldBuilder) Optional() *StringFieldBuilder {
	b.UnifiedFieldBuilder.Optional()
	return b
}

func (b *StringFieldBuilder) Unique() *StringFieldBuilder {
	b.UnifiedFieldBuilder.Unique()
	return b
}

func (b *StringFieldBuilder) Primary() *StringFieldBuilder {
	b.UnifiedFieldBuilder.Primary()
	return b
}

func (b *StringFieldBuilder) DBIndex() *StringFieldBuilder {
	b.UnifiedFieldBuilder.DBIndex()
	return b
}

func (b *StringFieldBuilder) MaxLength(n int) *StringFieldBuilder {
	b.UnifiedFieldBuilder.MaxLength(n)
	return b
}

func (b *StringFieldBuilder) MinLength(n int) *StringFieldBuilder {
	b.UnifiedFieldBuilder.MinLength(n)
	return b
}

func (b *StringFieldBuilder) Choices(choices ...Choice) *StringFieldBuilder {
	b.UnifiedFieldBuilder.Choices(choices...)
	return b
}

func (b *StringFieldBuilder) ChoicesFromPairs(pairs ...string) *StringFieldBuilder {
	b.UnifiedFieldBuilder.ChoicesFromPairs(pairs...)
	return b
}

func (b *StringFieldBuilder) VerboseName(name string) *StringFieldBuilder {
	b.UnifiedFieldBuilder.VerboseName(name)
	return b
}

func (b *StringFieldBuilder) HelpText(text string) *StringFieldBuilder {
	b.UnifiedFieldBuilder.HelpText(text)
	return b
}

// ============================================================================
// DateTimeFieldBuilder common method overrides
// ============================================================================

func (b *DateTimeFieldBuilder) Required() *DateTimeFieldBuilder {
	b.UnifiedFieldBuilder.Required()
	return b
}

func (b *DateTimeFieldBuilder) Optional() *DateTimeFieldBuilder {
	b.UnifiedFieldBuilder.Optional()
	return b
}

func (b *DateTimeFieldBuilder) AutoNow() *DateTimeFieldBuilder {
	b.UnifiedFieldBuilder.AutoNow()
	return b
}

func (b *DateTimeFieldBuilder) AutoNowAdd() *DateTimeFieldBuilder {
	b.UnifiedFieldBuilder.AutoNowAdd()
	return b
}

func (b *DateTimeFieldBuilder) VerboseName(name string) *DateTimeFieldBuilder {
	b.UnifiedFieldBuilder.VerboseName(name)
	return b
}

func (b *DateTimeFieldBuilder) HelpText(text string) *DateTimeFieldBuilder {
	b.UnifiedFieldBuilder.HelpText(text)
	return b
}

// ============================================================================
// DecimalFieldBuilder common method overrides
// ============================================================================

func (b *DecimalFieldBuilder) Required() *DecimalFieldBuilder {
	b.UnifiedFieldBuilder.Required()
	return b
}

func (b *DecimalFieldBuilder) Optional() *DecimalFieldBuilder {
	b.UnifiedFieldBuilder.Optional()
	return b
}

func (b *DecimalFieldBuilder) MaxDigits(digits int) *DecimalFieldBuilder {
	b.UnifiedFieldBuilder.MaxDigits(digits)
	return b
}

func (b *DecimalFieldBuilder) DecimalPlaces(places int) *DecimalFieldBuilder {
	b.UnifiedFieldBuilder.DecimalPlaces(places)
	return b
}

func (b *DecimalFieldBuilder) MaxValue(val float64) *DecimalFieldBuilder {
	b.UnifiedFieldBuilder.MaxValue(val)
	return b
}

func (b *DecimalFieldBuilder) MinValue(val float64) *DecimalFieldBuilder {
	b.UnifiedFieldBuilder.MinValue(val)
	return b
}

func (b *DecimalFieldBuilder) VerboseName(name string) *DecimalFieldBuilder {
	b.UnifiedFieldBuilder.VerboseName(name)
	return b
}

func (b *DecimalFieldBuilder) HelpText(text string) *DecimalFieldBuilder {
	b.UnifiedFieldBuilder.HelpText(text)
	return b
}

// ============================================================================
// UUIDFieldBuilder common method overrides
// ============================================================================

func (b *UUIDFieldBuilder) Required() *UUIDFieldBuilder {
	b.UnifiedFieldBuilder.Required()
	return b
}

func (b *UUIDFieldBuilder) Optional() *UUIDFieldBuilder {
	b.UnifiedFieldBuilder.Optional()
	return b
}

func (b *UUIDFieldBuilder) Unique() *UUIDFieldBuilder {
	b.UnifiedFieldBuilder.Unique()
	return b
}

func (b *UUIDFieldBuilder) Primary() *UUIDFieldBuilder {
	b.UnifiedFieldBuilder.Primary()
	return b
}

func (b *UUIDFieldBuilder) VerboseName(name string) *UUIDFieldBuilder {
	b.UnifiedFieldBuilder.VerboseName(name)
	return b
}

func (b *UUIDFieldBuilder) HelpText(text string) *UUIDFieldBuilder {
	b.UnifiedFieldBuilder.HelpText(text)
	return b
}
