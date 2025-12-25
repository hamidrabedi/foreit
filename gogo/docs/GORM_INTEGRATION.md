# GORM Integration - Type-Safe Queries

## Philosophy

**No string-based lookups. Direct struct field references. Direct to GORM.**

Unlike Django (which uses strings because Python), we leverage Go's type system to:
1. Reference struct fields directly (not strings)
2. Build conditions that map straight to GORM queries
3. Zero string parsing or conversion

## Architecture

```
User Code
    ↓
FieldRef (struct field metadata)
    ↓
Condition (GORM query function)
    ↓
GORMQuerySet (GORM *gorm.DB)
    ↓
Database
```

## Usage

### 1. Define Model with GORM Tags

```go
type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Email     string    `gorm:"uniqueIndex;not null" json:"email"`
    Name      string    `gorm:"not null" json:"name"`
    IsActive  bool      `gorm:"default:true" json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 2. Register Model

```go
func init() {
    RegisterModel[User](map[string]string{
        "ID":        "id",
        "Email":     "email",
        "Name":      "name",
        "IsActive":  "is_active",
        "CreatedAt": "created_at",
        "UpdatedAt": "updated_at",
    })
}
```

### 3. Create Type-Safe Field References

```go
var (
    UserEmail    = NewStringFieldRef[User]("Email")
    UserName     = NewStringFieldRef[User]("Name")
    UserID       = NewFieldRef[uint, User]("ID")
    UserIsActive = NewFieldRef[bool, User]("IsActive")
)
```

### 4. Build Queries (No Strings!)

```go
manager := NewGORMManager(db, User{})
queryset := manager.Filter(ctx)

// Type-safe conditions - direct to GORM
activeUsers := queryset.
    Filter(NewCondition[User](UserIsActive.ApplyEq(true))).
    Filter(NewCondition[User](UserEmail.ApplyIContains("example")))

users, err := activeUsers.All(ctx)
```

### 5. Combine Conditions

```go
emailCondition := NewCondition[User](UserEmail.ApplyIContains("example"))
activeCondition := NewCondition[User](UserIsActive.ApplyEq(true))

// AND
combined := emailCondition.And(activeCondition)

// OR
either := emailCondition.Or(activeCondition)

results, err := queryset.Filter(combined).All(ctx)
```

## Key Benefits

1. **✅ Type-Safe**: Field references are compile-time checked
2. **✅ No Strings**: Direct struct field access
3. **✅ Direct to GORM**: No string parsing, goes straight to `*gorm.DB`
4. **✅ IDE Support**: Full autocomplete and type checking
5. **✅ Performance**: Zero reflection overhead in hot paths

## Available Operations

### Numeric/String/Bool Fields
- `ApplyEq(value)` - Equality
- `ApplyNe(value)` - Not equal
- `ApplyGt(value)` - Greater than
- `ApplyGte(value)` - Greater than or equal
- `ApplyLt(value)` - Less than
- `ApplyLte(value)` - Less than or equal
- `ApplyIn(values)` - IN clause
- `ApplyIsNull()` - IS NULL
- `ApplyIsNotNull()` - IS NOT NULL

### String Fields Only
- `ApplyContains(value)` - LIKE '%value%'
- `ApplyIContains(value)` - ILIKE '%value%'
- `ApplyStartsWith(value)` - LIKE 'value%'
- `ApplyEndsWith(value)` - LIKE '%value'

## Migrations

Use `golang-migrate` for versioned migrations:

```bash
# Create migration
gogo makemigrations

# Apply migrations
gogo migrate

# Rollback
gogo migrate down 1
```

Migrations are SQL files, versioned, and work with any ORM.

