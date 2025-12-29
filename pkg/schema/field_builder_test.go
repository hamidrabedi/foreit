package schema

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestUnifiedFieldBuilder tests the unified field builder functionality
func TestUnifiedFieldBuilder(t *testing.T) {
	builder := newUnifiedFieldBuilder("test_field", TypeString)
	
	// Test common methods
	builder.Required().Unique().DBIndex()
	builder.MaxLength(100).MinLength(5)
	builder.VerboseName("Test Field").HelpText("This is a test field")
	
	field := builder.Build()
	
	if field.Name != "test_field" {
		t.Errorf("Expected field name 'test_field', got '%s'", field.Name)
	}
	if field.Type != TypeString {
		t.Errorf("Expected field type TypeString, got %v", field.Type)
	}
	if !field.Required {
		t.Error("Expected field to be required")
	}
	if !field.Unique {
		t.Error("Expected field to be unique")
	}
	if !field.DBIndex {
		t.Error("Expected field to have DB index")
	}
	if field.MaxLength == nil || *field.MaxLength != 100 {
		t.Error("Expected MaxLength to be 100")
	}
	if field.MinLength == nil || *field.MinLength != 5 {
		t.Error("Expected MinLength to be 5")
	}
	if field.VerboseName != "Test Field" {
		t.Errorf("Expected VerboseName 'Test Field', got '%s'", field.VerboseName)
	}
	if field.HelpText != "This is a test field" {
		t.Errorf("Expected HelpText 'This is a test field', got '%s'", field.HelpText)
	}
}

// TestInt64FieldBuilder tests Int64 field builder
func TestInt64FieldBuilder(t *testing.T) {
	field := Int64("id").
		Primary().
		AutoIncrement().
		Default(int64(0)).
		MaxValue(1000).
		MinValue(0).
		Required().
		Build()
	
	if field.Name != "id" {
		t.Errorf("Expected field name 'id', got '%s'", field.Name)
	}
	if field.Type != TypeInt64 {
		t.Errorf("Expected field type TypeInt64, got %v", field.Type)
	}
	if !field.PrimaryKey {
		t.Error("Expected field to be primary key")
	}
	if !field.AutoIncrement {
		t.Error("Expected field to be auto-increment")
	}
	if field.Default.(int64) != 0 {
		t.Error("Expected default value to be 0")
	}
	if field.MaxValue == nil || *field.MaxValue != 1000 {
		t.Error("Expected MaxValue to be 1000")
	}
	if field.MinValue == nil || *field.MinValue != 0 {
		t.Error("Expected MinValue to be 0")
	}
}

// TestStringFieldBuilder tests String field builder
func TestStringFieldBuilder(t *testing.T) {
	field := String("email").
		Required().
		Unique().
		MaxLength(255).
		MinLength(5).
		Default("test@example.com").
		Choices(Choice{Value: "active", Label: "Active"}).
		WriteOnly().
		Build()
	
	if field.Name != "email" {
		t.Errorf("Expected field name 'email', got '%s'", field.Name)
	}
	if field.Type != TypeString {
		t.Errorf("Expected field type TypeString, got %v", field.Type)
	}
	if !field.Required {
		t.Error("Expected field to be required")
	}
	if !field.Unique {
		t.Error("Expected field to be unique")
	}
	if field.MaxLength == nil || *field.MaxLength != 255 {
		t.Error("Expected MaxLength to be 255")
	}
	if field.MinLength == nil || *field.MinLength != 5 {
		t.Error("Expected MinLength to be 5")
	}
	if field.Default.(string) != "test@example.com" {
		t.Error("Expected default value to be 'test@example.com'")
	}
	if len(field.Choices) != 1 {
		t.Error("Expected 1 choice")
	}
	if field.Serialize {
		t.Error("Expected field to be write-only (not serialized), but Serialize is true")
	}
}

// TestDateTimeFieldBuilder tests DateTime field builder
func TestDateTimeFieldBuilder(t *testing.T) {
	now := time.Now()
	field := DateTime("created_at").
		AutoNowAdd().
		Default(now).
		Required().
		Build()
	
	if field.Name != "created_at" {
		t.Errorf("Expected field name 'created_at', got '%s'", field.Name)
	}
	if field.Type != TypeDateTime {
		t.Errorf("Expected field type TypeDateTime, got %v", field.Type)
	}
	if !field.AutoNowAdd {
		t.Error("Expected field to have AutoNowAdd")
	}
	if !field.Required {
		t.Error("Expected field to be required")
	}
}

// TestDecimalFieldBuilder tests Decimal field builder
func TestDecimalFieldBuilder(t *testing.T) {
	field := Decimal("price").
		MaxDigits(12).
		DecimalPlaces(2).
		Default(0.0).
		MaxValue(999999999.99).
		MinValue(0.0).
		Required().
		Build()
	
	if field.Name != "price" {
		t.Errorf("Expected field name 'price', got '%s'", field.Name)
	}
	if field.Type != TypeDecimal {
		t.Errorf("Expected field type TypeDecimal, got %v", field.Type)
	}
	if field.MaxDigits == nil || *field.MaxDigits != 12 {
		t.Error("Expected MaxDigits to be 12")
	}
	if field.DecimalPlaces == nil || *field.DecimalPlaces != 2 {
		t.Error("Expected DecimalPlaces to be 2")
	}
	if field.Default.(float64) != 0.0 {
		t.Error("Expected default value to be 0.0")
	}
}

// TestUUIDFieldBuilder tests UUID field builder
func TestUUIDFieldBuilder(t *testing.T) {
	testUUID := uuid.New()
	field := UUID("uuid").
		Required().
		Unique().
		Primary().
		DefaultUUID(testUUID).
		Build()
	
	if field.Name != "uuid" {
		t.Errorf("Expected field name 'uuid', got '%s'", field.Name)
	}
	if field.Type != TypeUUID {
		t.Errorf("Expected field type TypeUUID, got %v", field.Type)
	}
	if !field.Required {
		t.Error("Expected field to be required")
	}
	if !field.Unique {
		t.Error("Expected field to be unique")
	}
	if !field.PrimaryKey {
		t.Error("Expected field to be primary key")
	}
	if field.Default.(string) != testUUID.String() {
		t.Error("Expected default value to match UUID string")
	}
}

// TestFieldOptions tests FieldOptions composition
func TestFieldOptions(t *testing.T) {
	options := NewFieldOptions()
	
	options.DB.Column = "custom_column"
	options.DB.Type = "VARCHAR(255)"
	options.DB.Index = true
	
	options.Validation.Required = true
	options.Validation.Unique = true
	options.Validation.MaxLength = intPtr(100)
	
	options.Presentation.VerboseName = "Test Field"
	options.Presentation.HelpText = "Help text"
	
	field := Field{Name: "test", Type: TypeString}
	options.ApplyToField(&field)
	
	if field.DBColumn != "custom_column" {
		t.Errorf("Expected DBColumn 'custom_column', got '%s'", field.DBColumn)
	}
	if field.DBType != "VARCHAR(255)" {
		t.Errorf("Expected DBType 'VARCHAR(255)', got '%s'", field.DBType)
	}
	if !field.DBIndex {
		t.Error("Expected DBIndex to be true")
	}
	if !field.Required {
		t.Error("Expected Required to be true")
	}
	if !field.Unique {
		t.Error("Expected Unique to be true")
	}
	if field.MaxLength == nil || *field.MaxLength != 100 {
		t.Error("Expected MaxLength to be 100")
	}
	if field.VerboseName != "Test Field" {
		t.Errorf("Expected VerboseName 'Test Field', got '%s'", field.VerboseName)
	}
}

// TestBackwardCompatibility tests that old API still works
func TestBackwardCompatibility(t *testing.T) {
	// Test that all field type constructors work
	_ = Int64("int64_field")
	_ = Int32("int32_field")
	_ = String("string_field")
	_ = Text("text_field")
	_ = Bool("bool_field")
	_ = Time("time_field")
	_ = Date("date_field")
	_ = DateTime("datetime_field")
	_ = Email("email_field")
	_ = URL("url_field")
	_ = Float32("float32_field")
	_ = Float64("float64_field")
	_ = Decimal("decimal_field")
	_ = JSON("json_field")
	_ = Bytes("bytes_field")
	_ = UUID("uuid_field")
	_ = Int("int_field") // Alias for Int64
	
	// Test that common methods work on all builders
	intField := Int64("test").Required().Unique().DBIndex().Build()
	if !intField.Required || !intField.Unique || !intField.DBIndex {
		t.Error("Common methods should work on all field builders")
	}
	
	stringField := String("test").Required().Unique().DBIndex().Build()
	if !stringField.Required || !stringField.Unique || !stringField.DBIndex {
		t.Error("Common methods should work on all field builders")
	}
}

// TestChaining tests method chaining
func TestChaining(t *testing.T) {
	field := Int64("id").
		Primary().
		AutoIncrement().
		Required().
		Unique().
		DBIndex().
		VerboseName("ID").
		HelpText("Primary key").
		Build()
	
	if field.Name != "id" {
		t.Error("Method chaining should work correctly")
	}
	if !field.PrimaryKey || !field.AutoIncrement || !field.Required {
		t.Error("All chained methods should be applied")
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}
