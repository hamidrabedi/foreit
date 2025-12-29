package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgego/forge/pkg/schema"
)

// TestSchemaBuilders_AllFieldTypes tests that all field type builders work correctly
func TestSchemaBuilders_AllFieldTypes(t *testing.T) {
	tests := []struct {
		name      string
		builder   func(string) interface{ Build() schema.Field }
		fieldType schema.FieldType
	}{
		{"Int64", func(n string) interface{ Build() schema.Field } { return schema.Int64(n) }, schema.TypeInt64},
		{"Int32", func(n string) interface{ Build() schema.Field } { return schema.Int32(n) }, schema.TypeInt32},
		{"String", func(n string) interface{ Build() schema.Field } { return schema.String(n) }, schema.TypeString},
		{"Text", func(n string) interface{ Build() schema.Field } { return schema.Text(n) }, schema.TypeText},
		{"Bool", func(n string) interface{ Build() schema.Field } { return schema.Bool(n) }, schema.TypeBool},
		{"Time", func(n string) interface{ Build() schema.Field } { return schema.Time(n) }, schema.TypeTime},
		{"Date", func(n string) interface{ Build() schema.Field } { return schema.Date(n) }, schema.TypeDate},
		{"DateTime", func(n string) interface{ Build() schema.Field } { return schema.DateTime(n) }, schema.TypeDateTime},
		{"Email", func(n string) interface{ Build() schema.Field } { return schema.Email(n) }, schema.TypeEmail},
		{"URL", func(n string) interface{ Build() schema.Field } { return schema.URL(n) }, schema.TypeURL},
		{"Float32", func(n string) interface{ Build() schema.Field } { return schema.Float32(n) }, schema.TypeFloat32},
		{"Float64", func(n string) interface{ Build() schema.Field } { return schema.Float64(n) }, schema.TypeFloat64},
		{"Decimal", func(n string) interface{ Build() schema.Field } { return schema.Decimal(n) }, schema.TypeDecimal},
		{"JSON", func(n string) interface{ Build() schema.Field } { return schema.JSON(n) }, schema.TypeJSON},
		{"Bytes", func(n string) interface{ Build() schema.Field } { return schema.Bytes(n) }, schema.TypeBytes},
		{"UUID", func(n string) interface{ Build() schema.Field } { return schema.UUID(n) }, schema.TypeUUID},
		{"Int", func(n string) interface{ Build() schema.Field } { return schema.Int(n) }, schema.TypeInt64}, // Alias
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := tt.builder("test_field").Build()
			assert.Equal(t, "test_field", field.Name)
			assert.Equal(t, tt.fieldType, field.Type)
		})
	}
}

// TestSchemaBuilders_MethodChaining tests that method chaining works correctly
func TestSchemaBuilders_MethodChaining(t *testing.T) {
	field := schema.Int64("id").
		Primary().
		AutoIncrement().
		Required().
		Unique().
		DBIndex().
		VerboseName("ID").
		HelpText("Primary key identifier").
		DBColumn("user_id").
		MaxValue(1000).
		MinValue(1).
		Build()

	assert.True(t, field.PrimaryKey)
	assert.True(t, field.AutoIncrement)
	assert.True(t, field.Required)
	assert.True(t, field.Unique)
	assert.True(t, field.DBIndex)
	assert.Equal(t, "ID", field.VerboseName)
	assert.Equal(t, "Primary key identifier", field.HelpText)
	assert.Equal(t, "user_id", field.DBColumn)
	assert.NotNil(t, field.MaxValue)
	assert.NotNil(t, field.MinValue)
}

// TestSchemaBuilders_StringFieldOptions tests string field specific options
func TestSchemaBuilders_StringFieldOptions(t *testing.T) {
	field := schema.String("email").
		Required().
		Unique().
		MaxLength(255).
		MinLength(5).
		Choices(
			schema.Choice{Value: "active", Label: "Active"},
			schema.Choice{Value: "inactive", Label: "Inactive"},
		).
		WriteOnly().
		Build()

	assert.True(t, field.Required)
	assert.True(t, field.Unique)
	assert.NotNil(t, field.MaxLength)
	assert.Equal(t, 255, *field.MaxLength)
	assert.NotNil(t, field.MinLength)
	assert.Equal(t, 5, *field.MinLength)
	assert.Len(t, field.Choices, 2)
	assert.False(t, field.Serialize) // WriteOnly sets Serialize to false
}

// TestSchemaBuilders_ChoicesFromPairs tests the ChoicesFromPairs convenience method
func TestSchemaBuilders_ChoicesFromPairs(t *testing.T) {
	field := schema.String("status").
		ChoicesFromPairs("active", "Active", "inactive", "Inactive").
		Build()

	assert.Len(t, field.Choices, 2)
	assert.Equal(t, "active", field.Choices[0].Value)
	assert.Equal(t, "Active", field.Choices[0].Label)
	assert.Equal(t, "inactive", field.Choices[1].Value)
	assert.Equal(t, "Inactive", field.Choices[1].Label)
}

// TestSchemaBuilders_TemporalFields tests temporal field specific options
func TestSchemaBuilders_TemporalFields(t *testing.T) {
	t.Run("DateTime with AutoNow", func(t *testing.T) {
		field := schema.DateTime("updated_at").
			AutoNow().
			Build()

		assert.True(t, field.AutoNow)
		assert.False(t, field.AutoNowAdd)
	})

	t.Run("DateTime with AutoNowAdd", func(t *testing.T) {
		field := schema.DateTime("created_at").
			AutoNowAdd().
			Build()

		assert.False(t, field.AutoNow)
		assert.True(t, field.AutoNowAdd)
	})

	t.Run("Date with both auto options", func(t *testing.T) {
		field := schema.Date("modified_at").
			AutoNow().
			AutoNowAdd().
			Build()

		assert.True(t, field.AutoNow)
		assert.True(t, field.AutoNowAdd)
	})
}

// TestSchemaBuilders_DecimalFieldOptions tests decimal field specific options
func TestSchemaBuilders_DecimalFieldOptions(t *testing.T) {
	field := schema.Decimal("price").
		MaxDigits(12).
		DecimalPlaces(2).
		MaxValue(999999999.99).
		MinValue(0.0).
		Required().
		Build()

	assert.NotNil(t, field.MaxDigits)
	assert.Equal(t, 12, *field.MaxDigits)
	assert.NotNil(t, field.DecimalPlaces)
	assert.Equal(t, 2, *field.DecimalPlaces)
	assert.NotNil(t, field.MaxValue)
	assert.NotNil(t, field.MinValue)
	assert.True(t, field.Required)
}

// TestSchemaBuilders_UUIDFieldOptions tests UUID field specific options
func TestSchemaBuilders_UUIDFieldOptions(t *testing.T) {
	t.Run("UUID with DefaultUUID", func(t *testing.T) {
		field := schema.UUID("uuid").
			Required().
			Unique().
			Primary().
			Build()

		assert.True(t, field.Required)
		assert.True(t, field.Unique)
		assert.True(t, field.PrimaryKey)
		assert.Equal(t, schema.TypeUUID, field.Type)
	})
}

// TestSchemaBuilders_DefaultValues tests default value setting for different types
func TestSchemaBuilders_DefaultValues(t *testing.T) {
	t.Run("Int64 default", func(t *testing.T) {
		field := schema.Int64("count").Default(int64(10)).Build()
		assert.Equal(t, int64(10), field.Default)
	})

	t.Run("String default", func(t *testing.T) {
		field := schema.String("name").Default("John").Build()
		assert.Equal(t, "John", field.Default)
	})

	t.Run("Bool default", func(t *testing.T) {
		field := schema.Bool("active").Default(true).Build()
		assert.Equal(t, true, field.Default)
	})

	t.Run("Float64 default", func(t *testing.T) {
		field := schema.Float64("rate").Default(0.5).Build()
		assert.Equal(t, 0.5, field.Default)
	})
}

// TestSchemaBuilders_DatabaseOptions tests database-level options
func TestSchemaBuilders_DatabaseOptions(t *testing.T) {
	field := schema.String("title").
		DBColumn("post_title").
		DBType("VARCHAR(500)").
		DBCollation("utf8mb4_unicode_ci").
		DBComment("Post title field").
		DBTablespace("tablespace1").
		DBDefault("''").
		DBIndex().
		Build()

	assert.Equal(t, "post_title", field.DBColumn)
	assert.Equal(t, "VARCHAR(500)", field.DBType)
	assert.Equal(t, "utf8mb4_unicode_ci", field.DBCollation)
	assert.Equal(t, "Post title field", field.DBComment)
	assert.Equal(t, "tablespace1", field.DBTablespace)
	assert.Equal(t, "''", field.DBDefault)
	assert.True(t, field.DBIndex)
}

// TestSchemaBuilders_GeneratedColumns tests generated column options
func TestSchemaBuilders_GeneratedColumns(t *testing.T) {
	field := schema.String("full_name").
		GeneratedColumn("first_name || ' ' || last_name", true).
		Build()

	assert.True(t, field.Generated)
	assert.Equal(t, "first_name || ' ' || last_name", field.GeneratedExpr)
	assert.True(t, field.IsStored)
}

// TestSchemaBuilders_UniqueConstraints tests unique constraint options
func TestSchemaBuilders_UniqueConstraints(t *testing.T) {
	field := schema.String("slug").
		Unique().
		UniqueForDate("published_at").
		UniqueForMonth("published_at").
		UniqueForYear("published_at").
		Build()

	assert.True(t, field.Unique)
	assert.Equal(t, "published_at", field.UniqueForDate)
	assert.Equal(t, "published_at", field.UniqueForMonth)
	assert.Equal(t, "published_at", field.UniqueForYear)
}

// TestSchemaBuilders_Validators tests validator addition
func TestSchemaBuilders_Validators(t *testing.T) {
	validator := &testValidator{name: "test"}
	field := schema.String("email").
		Validators(validator).
		Build()

	assert.Len(t, field.Validators, 1)
	assert.Equal(t, validator, field.Validators[0])
}

// TestSchemaBuilders_ComplexField tests a complex field with many options
func TestSchemaBuilders_ComplexField(t *testing.T) {
	field := schema.Decimal("total_price").
		Required().
		Unique().
		MaxDigits(12).
		DecimalPlaces(2).
		MaxValue(999999999.99).
		MinValue(0.0).
		Default(0.0).
		DBIndex().
		VerboseName("Total Price").
		HelpText("Total price including tax and shipping").
		DBColumn("total_price_usd").
		DBComment("Stored in USD").
		Build()

	assert.True(t, field.Required)
	assert.True(t, field.Unique)
	assert.NotNil(t, field.MaxDigits)
	assert.Equal(t, 12, *field.MaxDigits)
	assert.NotNil(t, field.DecimalPlaces)
	assert.Equal(t, 2, *field.DecimalPlaces)
	assert.NotNil(t, field.MaxValue)
	assert.NotNil(t, field.MinValue)
	assert.Equal(t, 0.0, field.Default)
	assert.True(t, field.DBIndex)
	assert.Equal(t, "Total Price", field.VerboseName)
	assert.Equal(t, "Total price including tax and shipping", field.HelpText)
	assert.Equal(t, "total_price_usd", field.DBColumn)
	assert.Equal(t, "Stored in USD", field.DBComment)
}

// TestSchemaBuilders_FieldOptionsComposition tests FieldOptions composition
func TestSchemaBuilders_FieldOptionsComposition(t *testing.T) {
	options := schema.NewFieldOptions()

	// Set DB options
	options.DB.Column = "custom_col"
	options.DB.Type = "TEXT"
	options.DB.Index = true

	// Set validation options
	options.Validation.Required = true
	options.Validation.Unique = true
	options.Validation.MaxLength = intPtr(100)

	// Set presentation options
	options.Presentation.VerboseName = "Custom Field"
	options.Presentation.HelpText = "Help text"

	// Apply to field
	field := schema.Field{Name: "test", Type: schema.TypeString}
	options.ApplyToField(&field)

	assert.Equal(t, "custom_col", field.DBColumn)
	assert.Equal(t, "TEXT", field.DBType)
	assert.True(t, field.DBIndex)
	assert.True(t, field.Required)
	assert.True(t, field.Unique)
	assert.NotNil(t, field.MaxLength)
	assert.Equal(t, 100, *field.MaxLength)
	assert.Equal(t, "Custom Field", field.VerboseName)
	assert.Equal(t, "Help text", field.HelpText)
}

// TestSchemaBuilders_BackwardCompatibility tests that old API patterns still work
func TestSchemaBuilders_BackwardCompatibility(t *testing.T) {
	// Test that all field constructors work
	fields := []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("name").Required().MaxLength(100).Build(),
		schema.Bool("active").Default(true).Build(),
		schema.DateTime("created_at").AutoNowAdd().Build(),
		schema.Decimal("price").MaxDigits(10).DecimalPlaces(2).Build(),
		schema.UUID("uuid").Required().Unique().Build(),
	}

	assert.Len(t, fields, 6)
	for _, field := range fields {
		assert.NotEmpty(t, field.Name)
		assert.NotZero(t, field.Type)
	}
}

// TestSchemaBuilders_IntAlias tests that Int() alias works correctly
func TestSchemaBuilders_IntAlias(t *testing.T) {
	field1 := schema.Int("id").Build()
	field2 := schema.Int64("id").Build()

	assert.Equal(t, schema.TypeInt64, field1.Type)
	assert.Equal(t, schema.TypeInt64, field2.Type)
	assert.Equal(t, field1.Type, field2.Type)
}

// TestSchemaBuilders_TextAlias tests that Text() returns StringFieldBuilder with TypeText
func TestSchemaBuilders_TextAlias(t *testing.T) {
	field := schema.Text("content").Build()

	assert.Equal(t, schema.TypeText, field.Type)
	assert.Equal(t, "content", field.Name)
}

// Helper function
func intPtr(i int) *int {
	return &i
}

// testValidator is a simple validator for testing
type testValidator struct {
	name string
}

func (v *testValidator) Validate(value interface{}) error {
	return nil
}
