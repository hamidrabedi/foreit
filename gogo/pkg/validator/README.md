# Validator Module

Validation using [go-playground/validator/v10](https://github.com/go-playground/validator), the most popular Go validation library.

## Usage

### Basic Validation

```go
import "github.com/gogo/pkg/validator"

type User struct {
    Name  string `validate:"required,min=3,max=100"`
    Email string `validate:"required,email"`
    Age   int    `validate:"gte=18,lte=120"`
}

user := &User{
    Name:  "John",
    Email: "john@example.com",
    Age:   25,
}

v := validator.New()
if err := v.Validate(user); err != nil {
    // Handle validation error
}
```

### Available Tags

- `required` - Field is required
- `email` - Must be valid email
- `min=value` - Minimum value/length
- `max=value` - Maximum value/length
- `len=value` - Exact length
- `gte=value` - Greater than or equal
- `lte=value` - Less than or equal
- `gt=value` - Greater than
- `lt=value` - Less than
- `oneof=val1 val2` - Must be one of the values
- `url` - Must be valid URL
- `uuid` - Must be valid UUID
- `alpha` - Only letters
- `alphanum` - Letters and numbers
- `numeric` - Numeric characters

### Custom Validators

```go
v := validator.New()
v.RegisterValidation("customtag", func(fl validator.FieldLevel) bool {
    // Custom validation logic
    return true
})
```

### Struct-Level Validation

```go
v.RegisterStructValidation(func(sl validator.StructLevel) {
    // Validate entire struct
}, MyStruct{})
```
