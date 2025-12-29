package schema

// FieldTrait represents a capability that a field type can have
// This uses the Strategy pattern to handle type-specific behaviors
type FieldTrait interface {
	// Apply applies the trait to the field builder
	Apply(builder *UnifiedFieldBuilder)
}

// NumericTrait provides numeric field capabilities (MinValue, MaxValue, AutoIncrement)
type NumericTrait struct{}

func (t *NumericTrait) Apply(builder *UnifiedFieldBuilder) {
	// Numeric traits are applied via builder methods
}

// StringTrait provides string field capabilities (MinLength, MaxLength, Choices)
type StringTrait struct{}

func (t *StringTrait) Apply(builder *UnifiedFieldBuilder) {
	// String traits are applied via builder methods
}

// TemporalTrait provides temporal field capabilities (AutoNow, AutoNowAdd)
type TemporalTrait struct{}

func (t *TemporalTrait) Apply(builder *UnifiedFieldBuilder) {
	// Temporal traits are applied via builder methods
}

// DecimalTrait provides decimal field capabilities (MaxDigits, DecimalPlaces)
type DecimalTrait struct{}

func (t *DecimalTrait) Apply(builder *UnifiedFieldBuilder) {
	// Decimal traits are applied via builder methods
}

// FieldTypeInfo contains metadata about a field type
type FieldTypeInfo struct {
	Type         FieldType
	Traits       []FieldTrait
	DefaultValue interface{}
}

// GetFieldTypeInfo returns the FieldTypeInfo for a given FieldType
func GetFieldTypeInfo(fieldType FieldType) FieldTypeInfo {
	infos := map[FieldType]FieldTypeInfo{
		TypeInt64: {
			Type:   TypeInt64,
			Traits: []FieldTrait{&NumericTrait{}},
		},
		TypeInt32: {
			Type:   TypeInt32,
			Traits: []FieldTrait{&NumericTrait{}},
		},
		TypeString: {
			Type:   TypeString,
			Traits: []FieldTrait{&StringTrait{}},
		},
		TypeText: {
			Type:   TypeText,
			Traits: []FieldTrait{&StringTrait{}},
		},
		TypeEmail: {
			Type:   TypeEmail,
			Traits: []FieldTrait{&StringTrait{}},
		},
		TypeURL: {
			Type:   TypeURL,
			Traits: []FieldTrait{&StringTrait{}},
		},
		TypeBool: {
			Type: TypeBool,
		},
		TypeTime: {
			Type:   TypeTime,
			Traits: []FieldTrait{&TemporalTrait{}},
		},
		TypeDate: {
			Type:   TypeDate,
			Traits: []FieldTrait{&TemporalTrait{}},
		},
		TypeDateTime: {
			Type:   TypeDateTime,
			Traits: []FieldTrait{&TemporalTrait{}},
		},
		TypeFloat32: {
			Type:   TypeFloat32,
			Traits: []FieldTrait{&NumericTrait{}},
		},
		TypeFloat64: {
			Type:   TypeFloat64,
			Traits: []FieldTrait{&NumericTrait{}},
		},
		TypeDecimal: {
			Type:   TypeDecimal,
			Traits: []FieldTrait{&NumericTrait{}, &DecimalTrait{}},
		},
		TypeUUID: {
			Type: TypeUUID,
		},
		TypeJSON: {
			Type: TypeJSON,
		},
		TypeBytes: {
			Type:   TypeBytes,
			Traits: []FieldTrait{&StringTrait{}}, // Bytes can have length constraints
		},
	}

	if info, ok := infos[fieldType]; ok {
		return info
	}
	return FieldTypeInfo{Type: fieldType}
}
