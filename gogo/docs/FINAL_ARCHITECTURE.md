# Final Architecture Decision

## ORM & Migration Solution

### ✅ Decision: GORM + golang-migrate

**Why this combination:**
1. **GORM**: Mature, handles SQL/scanning/relationships - we don't rebuild ORM
2. **golang-migrate**: Industry standard for versioned migrations (like Django)
3. **Type-Safe Wrapper**: We build a thin layer on top for Django-like API

## Architecture

```
User Code (Type-Safe)
    ↓
FieldRef (struct field metadata, cached)
    ↓
Condition (GORM query function)
    ↓
GORMQuerySet (wraps *gorm.DB)
    ↓
GORM (handles SQL/scanning)
    ↓
Database

Migrations:
    ↓
golang-migrate (versioned SQL files)
    ↓
Database
```

## Key Principles

### 1. No String-Based Lookups
- ❌ **NOT**: `Filter("email__icontains", "example")`
- ✅ **YES**: `Filter(NewCondition[User](UserEmail.ApplyIContains("example")))`

### 2. Direct Struct Field References
- Field references map directly to struct fields
- Column names extracted from GORM tags or struct field names
- Cached metadata for performance

### 3. Direct to GORM
- Conditions are functions that modify `*gorm.DB`
- No string parsing or conversion
- Goes straight to GORM's query builder

### 4. Versioned Migrations
- SQL files in `migrations/` directory
- Versioned like Django: `0001_initial.up.sql`, `0001_initial.down.sql`
- Managed by `golang-migrate`

## Usage Example

```go
// 1. Define model with GORM tags
type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Email     string    `gorm:"uniqueIndex;not null" json:"email"`
    Name      string    `gorm:"not null" json:"name"`
    IsActive  bool      `gorm:"default:true" json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// 2. Register model
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

// 3. Create type-safe field references
var (
    UserEmail    = NewStringFieldRef[User]("Email")
    UserName     = NewStringFieldRef[User]("Name")
    UserID       = NewFieldRef[uint, User]("ID")
    UserIsActive = NewFieldRef[bool, User]("IsActive")
)

// 4. Use with GORM
db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
manager := NewGORMManager(db, User{})
queryset := manager.Filter(ctx)

// Type-safe queries - no strings!
activeUsers := queryset.
    Filter(NewCondition[User](UserIsActive.ApplyEq(true))).
    Filter(NewCondition[User](UserEmail.ApplyIContains("example")))

users, err := activeUsers.All(ctx)
```

## Migration Workflow

```bash
# 1. Create migration files manually or with tool
# migrations/0001_initial.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

# migrations/0001_initial.down.sql
DROP TABLE users;

# 2. Apply migrations
gogo migrate up

# 3. Check version
gogo migrate version

# 4. Rollback if needed
gogo migrate down
```

## Benefits

1. **✅ Type-Safe**: Compile-time checked field references
2. **✅ No Strings**: Direct struct field access
3. **✅ Direct to GORM**: No string parsing overhead
4. **✅ Versioned Migrations**: Like Django, industry standard
5. **✅ Mature Backend**: GORM handles all the hard parts
6. **✅ Production Ready**: Both tools are battle-tested

## What We Build vs. What We Use

### We Build:
- Type-safe field references (`FieldRef`)
- Condition builder (`Condition`)
- QuerySet wrapper (`GORMQuerySet`)
- Manager wrapper (`GORMManager`)
- Migration CLI integration

### We Use:
- **GORM**: SQL building, row scanning, relationships
- **golang-migrate**: Versioned migrations

This gives us the best of both worlds: Django-like DX with Go's type safety and proven backend tools.

