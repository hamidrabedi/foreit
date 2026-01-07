# Extending the Schema System

This guide explains how to extend the schema system with custom field types, validators, and builders for third-party packages or user-specific needs.

## Architecture Overview

The schema system uses a composition-based architecture with a unified builder pattern:

- **UnifiedFieldBuilder**: Core builder containing all common field options (DBColumn, DBCollation, Required, etc.)
- **Type-Specific Builders**: Wrap `UnifiedFieldBuilder` and provide type-safe methods (Int64FieldBuilder, StringFieldBuilder, etc.)
- **BaseFieldBuilder**: Legacy builder kept for backward compatibility (new code should use type-specific builders)
- **FieldOptions**: Organized configuration with separated concerns (DB, Validation, Presentation)
- **Validators**: Composable validation functions that can be applied to any field
- **Field Registry**: Allows registration of custom field types at runtime

### Recommended Approach

For new code, use the type-specific builders directly:
```go
schema.Int64("id").Primary().AutoIncrement().Build()
schema.String("email").Required().Unique().MaxLength(255).Build()
```

For custom field builders, you can either:
1. Embed `BaseFieldBuilder` (backward compatible)
2. Embed `UnifiedFieldBuilder` (recommended for new code)

## Using Validators

Instead of creating separate field types for specialized formats (like slugs, IP addresses), use validators with base field types:

### Example: Slug Field

```go
import "github.com/forgego/forge/pkg/schema"

func (User) Fields() []schema.Field {
    return []schema.Field{
        schema.String("slug").
            Required().
            Unique().
            MaxLength(100).
            Validators(schema.SlugValidator()).
            Build(),
    }
}
```

### Example: IP Address Field

```go
func (Network) Fields() []schema.Field {
    return []schema.Field{
        schema.String("ip_address").
            Required().
            Validators(schema.IPAddressValidator()).
            Build(),
    }
}
```

### Available Validators

- `SlugValidator()` - Validates URL-friendly slug format
- `IPAddressValidator()` - Validates IPv4/IPv6 addresses
- `UUIDValidator()` - Validates UUID format
- `EmailValidator()` - Validates email format
- `URLValidator()` - Validates URL format
- `MinLengthValidator(min int)` - Validates minimum length
- `MaxLengthValidator(max int)` - Validates maximum length

### Creating Custom Validators

Implement the `Validator` interface:

```go
type MyCustomValidator struct {
    // Your validator state
}

func (v *MyCustomValidator) Validate(value interface{}) error {
    // Your validation logic
    str, ok := value.(string)
    if !ok {
        return fmt.Errorf("expected string, got %T", value)
    }
    
    // Validate the value
    if !isValid(str) {
        return fmt.Errorf("validation failed")
    }
    
    return nil
}

// Usage
func (Model) Fields() []schema.Field {
    return []schema.Field{
        schema.String("field").
            Validators(&MyCustomValidator{}).
            Build(),
    }
}
```

## Creating Custom Field Builders

For more complex field types, you can create custom field builders. You have two options:

### Option 1: Using UnifiedFieldBuilder (Recommended)

```go
package mypackage

import "github.com/forgego/forge/pkg/schema"

// MyCustomFieldBuilder is a custom field builder using the unified architecture
type MyCustomFieldBuilder struct {
    *schema.UnifiedFieldBuilder
    customOption string
}

// NewMyCustomField creates a new custom field builder
func NewMyCustomField(name string) *MyCustomFieldBuilder {
    return &MyCustomFieldBuilder{
        UnifiedFieldBuilder: schema.newUnifiedFieldBuilder(name, schema.TypeString),
        customOption: "",
    }
}
```

### Option 2: Using BaseFieldBuilder (Backward Compatible)

```go
package mypackage

import "github.com/forgego/forge/pkg/schema"

// MyCustomFieldBuilder is a custom field builder (legacy approach)
type MyCustomFieldBuilder struct {
    *schema.BaseFieldBuilder
    customOption string
}

// NewMyCustomField creates a new custom field builder
func NewMyCustomField(name string) *MyCustomFieldBuilder {
    return &MyCustomFieldBuilder{
        BaseFieldBuilder: &schema.BaseFieldBuilder{
            field: schema.Field{
                Name: name,
                Type: schema.TypeString, // or your custom type
            },
        },
    }
}

// CustomOption sets a custom option specific to this field type
func (b *MyCustomFieldBuilder) CustomOption(value string) *MyCustomFieldBuilder {
    b.customOption = value
    return b
}

// Build returns the built field
func (b *MyCustomFieldBuilder) Build() schema.Field {
    // Apply custom logic if needed
    return b.field
}

// Usage
func (Model) Fields() []schema.Field {
    return []schema.Field{
        NewMyCustomField("custom_field").
            Required().
            CustomOption("value").
            DBIndex().
            Build(),
    }
}
```

## Registering Custom Field Types

For field types that should be available globally (e.g., in a third-party package), register them with the field registry:

### Example: Registering a Custom Field Type

```go
package mypackage

import (
    "github.com/forgego/forge/pkg/schema"
)

func init() {
    // Register the custom field type
    schema.RegisterFieldType("my_custom_type", func(name string) interface{} {
        return NewMyCustomField(name)
    })
}

// MyCustomFieldBuilder must implement CustomFieldBuilder interface
func (b *MyCustomFieldBuilder) GetFieldType() schema.FieldType {
    return schema.TypeString // or your custom type
}
```

### Using Registered Field Types

```go
import "github.com/forgego/forge/pkg/schema"

func (Model) Fields() []schema.Field {
    builder, err := schema.NewFieldBuilder("my_custom_type", "field_name")
    if err != nil {
        // Handle error
    }
    
    field := builder.
        Required().
        DBIndex().
        Build()
    
    return []schema.Field{field}
}
```

## Best Practices

### 1. Prefer Validators Over Custom Field Types

For simple format validation (slugs, IPs, emails), use validators with base field types rather than creating separate field builders.

**Good:**
```go
schema.String("slug").Validators(schema.SlugValidator())
```

**Avoid:**
```go
// Don't create a separate SlugFieldBuilder unless you need complex behavior
```

### 2. Embed UnifiedFieldBuilder or BaseFieldBuilder

For new code, embed `*UnifiedFieldBuilder` to use the modern architecture:
```go
type MyFieldBuilder struct {
    *schema.UnifiedFieldBuilder
    // Your custom fields
}
```

For backward compatibility, you can still embed `*BaseFieldBuilder`:
```go
type MyFieldBuilder struct {
    *schema.BaseFieldBuilder
    // Your custom fields
}
```

### 3. Keep Type-Specific Methods Minimal

Only add methods that are truly specific to your field type. Common options (Required, DBIndex, etc.) are already available through `UnifiedFieldBuilder` or `BaseFieldBuilder`.

### 4. Document Your Extensions

If creating a third-party package, document:
- What the field type does
- How to register it
- Example usage
- Any special requirements

## Migration from Specialized Field Types

If you were using specialized field types that have been removed:

### Before (removed):
```go
schema.Slug("slug").Required().Unique()
schema.IPAddress("ip").Required()
```

### After (using validators):
```go
schema.String("slug").
    Required().
    Unique().
    Validators(schema.SlugValidator())

schema.String("ip").
    Required().
    Validators(schema.IPAddressValidator())
```

## Advanced: Custom Field Types

For truly custom field types (beyond what validators can provide), you can:

1. Define a new `FieldType` constant (if modifying the core)
2. Create a custom builder that embeds `BaseFieldBuilder`
3. Register it with the field registry
4. Implement custom database mapping logic in your migration layer

However, this should be rare - most use cases can be handled with validators and the existing field types.

## Examples

### Example 1: Phone Number Field with Validator

```go
type PhoneValidator struct{}

func (v *PhoneValidator) Validate(value interface{}) error {
    str, ok := value.(string)
    if !ok {
        return fmt.Errorf("expected string")
    }
    // Validate phone format
    return nil
}

func (User) Fields() []schema.Field {
    return []schema.Field{
        schema.String("phone").
            Required().
            Validators(&PhoneValidator{}).
            Build(),
    }
}
```

### Example 2: Custom Field Builder with Registration

```go
package myfields

import "github.com/forgego/forge/pkg/schema"

type ColorFieldBuilder struct {
    *schema.BaseFieldBuilder
}

func Color(name string) *ColorFieldBuilder {
    return &ColorFieldBuilder{
        BaseFieldBuilder: &schema.BaseFieldBuilder{
            field: schema.Field{
                Name: name,
                Type: schema.TypeString,
            },
        },
    }
}

func (b *ColorFieldBuilder) Build() schema.Field {
    return b.field
}

func (b *ColorFieldBuilder) GetFieldType() schema.FieldType {
    return schema.TypeString
}

func init() {
    schema.RegisterFieldType("color", func(name string) interface{} {
        return Color(name)
    })
}
```

## Summary

- Use **validators** for format validation (slugs, IPs, emails)
- Use **custom field builders** for complex field types with special behavior
- Use **field registry** for globally available custom field types
- For new code: embed `UnifiedFieldBuilder` to inherit common functionality
- For backward compatibility: embed `BaseFieldBuilder` (still supported)
- Keep custom builders minimal - only add type-specific methods

## Architecture Notes

The schema package has been refactored to use a unified builder architecture:
- **UnifiedFieldBuilder**: Single source of truth for all common field methods
- **Type-specific builders**: Provide type-safe wrappers (Int64FieldBuilder, StringFieldBuilder, etc.)
- **FieldOptions**: Separated configuration (DB, Validation, Presentation)
- **Backward compatible**: All existing code continues to work without changes

See `pkg/schema/ARCHITECTURE.md` for detailed architecture documentation.
