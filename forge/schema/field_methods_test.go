package schema

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestFieldChaining(t *testing.T) {
	field := String("email").
		WithRequired().
		WithUnique().
		WithMaxLength(255).
		WithMinLength(5).
		WithDefault("test@example.com").
		WithChoices(Choice{Value: "active", Label: "Active"}).
		WithWriteOnly()

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
		t.Error("Expected field to be write-only (not serialized)")
	}
}

func TestInt64Field(t *testing.T) {
	field := Int64("id").
		WithPrimary().
		WithAutoIncrement().
		WithDefault(int64(0)).
		WithMaxValue(1000).
		WithMinValue(0).
		WithRequired()

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

func TestDateTimeField(t *testing.T) {
	now := time.Now()
	field := DateTime("created_at").
		WithAutoNowAdd().
		WithDefault(now).
		WithRequired()

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

func TestDecimalField(t *testing.T) {
	field := Decimal("price").
		WithMaxDigits(12).
		WithDecimalPlaces(2).
		WithDefault(0.0).
		WithMaxValue(999999999.99).
		WithMinValue(0.0).
		WithRequired()

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

func TestUUIDField(t *testing.T) {
	testUUID := uuid.New()
	field := UUID("uuid").
		WithRequired().
		WithUnique().
		WithPrimary().
		WithDefaultUUID(testUUID)

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

func TestFieldDefaults(t *testing.T) {
	field := String("name")
	if !field.Editable {
		t.Error("Expected Editable to default to true")
	}
	if !field.Serialize {
		t.Error("Expected Serialize to default to true")
	}
}

func TestFieldTypeConstructors(t *testing.T) {
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

	intField := Int64("test").WithRequired().WithUnique().WithDBIndex()
	if !intField.Required || !intField.Unique || !intField.DBIndex {
		t.Error("Common methods should work on all field types")
	}
}
