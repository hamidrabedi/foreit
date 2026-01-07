package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/forgego/forge/schema"
)

func TestSchemaFields_AllFieldTypes(t *testing.T) {
	tests := []struct {
		name      string
		builder   func(string) schema.Field
		fieldType schema.FieldType
	}{
		{"Int64", schema.Int64, schema.TypeInt64},
		{"Int32", schema.Int32, schema.TypeInt32},
		{"String", schema.String, schema.TypeString},
		{"Text", schema.Text, schema.TypeText},
		{"Bool", schema.Bool, schema.TypeBool},
		{"Time", schema.Time, schema.TypeTime},
		{"Date", schema.Date, schema.TypeDate},
		{"DateTime", schema.DateTime, schema.TypeDateTime},
		{"Email", schema.Email, schema.TypeEmail},
		{"URL", schema.URL, schema.TypeURL},
		{"Float32", schema.Float32, schema.TypeFloat32},
		{"Float64", schema.Float64, schema.TypeFloat64},
		{"Decimal", schema.Decimal, schema.TypeDecimal},
		{"JSON", schema.JSON, schema.TypeJSON},
		{"Bytes", schema.Bytes, schema.TypeBytes},
		{"UUID", schema.UUID, schema.TypeUUID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := tt.builder("test_field")
			assert.Equal(t, "test_field", field.Name)
			assert.Equal(t, tt.fieldType, field.Type)
		})
	}
}

func TestSchemaFields_MethodChaining(t *testing.T) {
	field := schema.Int64("id").
		WithPrimary().
		WithAutoIncrement().
		WithRequired().
		WithUnique().
		WithDBIndex().
		WithVerboseName("ID").
		WithHelpText("Primary key identifier").
		WithDBColumn("user_id").
		WithMaxValue(1000).
		WithMinValue(1)

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

func TestSchemaFields_StringOptions(t *testing.T) {
	field := schema.String("email").
		WithRequired().
		WithUnique().
		WithMaxLength(255).
		WithMinLength(5).
		WithChoices(
			schema.Choice{Value: "active", Label: "Active"},
			schema.Choice{Value: "inactive", Label: "Inactive"},
		).
		WithWriteOnly()

	assert.True(t, field.Required)
	assert.True(t, field.Unique)
	assert.NotNil(t, field.MaxLength)
	assert.Equal(t, 255, *field.MaxLength)
	assert.NotNil(t, field.MinLength)
	assert.Equal(t, 5, *field.MinLength)
	assert.Len(t, field.Choices, 2)
	assert.False(t, field.Serialize)
}

func TestSchemaFields_ChoicesFromPairs(t *testing.T) {
	field := schema.String("status").
		WithChoicesFromPairs("active", "Active", "inactive", "Inactive")

	assert.Len(t, field.Choices, 2)
	assert.Equal(t, "active", field.Choices[0].Value)
	assert.Equal(t, "Active", field.Choices[0].Label)
	assert.Equal(t, "inactive", field.Choices[1].Value)
	assert.Equal(t, "Inactive", field.Choices[1].Label)
}

func TestSchemaFields_TemporalOptions(t *testing.T) {
	t.Run("DateTime with AutoNow", func(t *testing.T) {
		field := schema.DateTime("updated_at").
			WithAutoNow()

		assert.True(t, field.AutoNow)
		assert.False(t, field.AutoNowAdd)
	})

	t.Run("DateTime with AutoNowAdd", func(t *testing.T) {
		field := schema.DateTime("created_at").
			WithAutoNowAdd()

		assert.False(t, field.AutoNow)
		assert.True(t, field.AutoNowAdd)
	})

	t.Run("Date with both auto options", func(t *testing.T) {
		field := schema.Date("modified_at").
			WithAutoNow().
			WithAutoNowAdd()

		assert.True(t, field.AutoNow)
		assert.True(t, field.AutoNowAdd)
	})
}

func TestSchemaFields_DecimalOptions(t *testing.T) {
	field := schema.Decimal("price").
		WithMaxDigits(12).
		WithDecimalPlaces(2).
		WithMaxValue(999999999.99).
		WithMinValue(0.0).
		WithRequired()

	assert.NotNil(t, field.MaxDigits)
	assert.Equal(t, 12, *field.MaxDigits)
	assert.NotNil(t, field.DecimalPlaces)
	assert.Equal(t, 2, *field.DecimalPlaces)
	assert.NotNil(t, field.MaxValue)
	assert.NotNil(t, field.MinValue)
	assert.True(t, field.Required)
}

func TestSchemaFields_DefaultValues(t *testing.T) {
	t.Run("Int64 default", func(t *testing.T) {
		field := schema.Int64("count").WithDefault(int64(10))
		assert.Equal(t, int64(10), field.Default)
	})

	t.Run("String default", func(t *testing.T) {
		field := schema.String("name").WithDefault("John")
		assert.Equal(t, "John", field.Default)
	})

	t.Run("Bool default", func(t *testing.T) {
		field := schema.Bool("active").WithDefault(true)
		assert.Equal(t, true, field.Default)
	})

	t.Run("Float64 default", func(t *testing.T) {
		field := schema.Float64("rate").WithDefault(0.5)
		assert.Equal(t, 0.5, field.Default)
	})
}

func TestSchemaFields_DatabaseOptions(t *testing.T) {
	field := schema.String("title").
		WithDBColumn("post_title").
		WithDBType("VARCHAR(500)").
		WithDBCollation("utf8mb4_unicode_ci").
		WithDBComment("Post title field").
		WithDBTablespace("tablespace1").
		WithDBDefault("''").
		WithDBIndex()

	assert.Equal(t, "post_title", field.DBColumn)
	assert.Equal(t, "VARCHAR(500)", field.DBType)
	assert.Equal(t, "utf8mb4_unicode_ci", field.DBCollation)
	assert.Equal(t, "Post title field", field.DBComment)
	assert.Equal(t, "tablespace1", field.DBTablespace)
	assert.Equal(t, "''", field.DBDefault)
	assert.True(t, field.DBIndex)
}

func TestSchemaFields_GeneratedColumns(t *testing.T) {
	field := schema.String("full_name").
		WithGeneratedColumn("first_name || ' ' || last_name", true)

	assert.True(t, field.Generated)
	assert.Equal(t, "first_name || ' ' || last_name", field.GeneratedExpr)
	assert.True(t, field.IsStored)
}

func TestSchemaFields_UniqueConstraints(t *testing.T) {
	field := schema.String("slug").
		WithUnique().
		WithUniqueForDate("published_at").
		WithUniqueForMonth("published_at").
		WithUniqueForYear("published_at")

	assert.True(t, field.Unique)
	assert.Equal(t, "published_at", field.UniqueForDate)
	assert.Equal(t, "published_at", field.UniqueForMonth)
	assert.Equal(t, "published_at", field.UniqueForYear)
}

func TestSchemaFields_Validators(t *testing.T) {
	validator := &testValidator{name: "test"}
	field := schema.String("email").
		WithValidators(validator)

	assert.Len(t, field.Validators, 1)
	assert.Equal(t, validator, field.Validators[0])
}

func TestSchemaFields_ComplexField(t *testing.T) {
	field := schema.Decimal("total_price").
		WithRequired().
		WithUnique().
		WithMaxDigits(12).
		WithDecimalPlaces(2).
		WithMaxValue(999999999.99).
		WithMinValue(0.0).
		WithDefault(0.0).
		WithDBIndex().
		WithVerboseName("Total Price").
		WithHelpText("Total price including tax and shipping").
		WithDBColumn("total_price_usd").
		WithDBComment("Stored in USD")

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

func TestSchemaFields_FieldOptionsComposition(t *testing.T) {
	options := schema.NewFieldOptions()

	options.DB.Column = "custom_col"
	options.DB.Type = "TEXT"
	options.DB.Index = true

	options.Validation.Required = true
	options.Validation.Unique = true
	options.Validation.MaxLength = intPtr(100)

	options.Presentation.VerboseName = "Custom Field"
	options.Presentation.HelpText = "Help text"

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

// Helper function
func intPtr(i int) *int {
	return &i
}

type testValidator struct {
	name string
}

func (v *testValidator) Validate(value interface{}) error {
	return nil
}
