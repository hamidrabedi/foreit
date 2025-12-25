# Type Safety in Gogo

## Philosophy

**Everything is type-safe. No strings. No reflection in the API. Compile-time safety.**

## Core Principles

1. **Generics First** - Use Go 1.18+ generics for all model operations
2. **Field References** - Type-safe field references, not string names
3. **No Type Assertions** - Results are typed, not interfaces
4. **Compile-Time Checks** - Errors caught at compile time, not runtime

## Type-Safe Components

### 1. Generic Manager

```go
// Type-safe manager - returns *User, not Model
userManager := NewBaseManager[*User](repo, meta)

// All operations are type-safe
user, err := userManager.Get(ctx, 1)        // Returns *User
users, err := userManager.All(ctx)          // Returns []*User
err := userManager.Create(ctx, user)        // Takes *User
err := userManager.Update(ctx, user)        // Takes *User
err := userManager.Delete(ctx, user)        // Takes *User
```

### 2. Type-Safe QuerySet

```go
// Type-safe queryset - returns []*User
users, err := userManager.Filter(ctx).
    Filter(UserEmail.Contains("@example.com")).  // Type-safe field ref
    Filter(UserName.Eq("John")).                 // Type-safe field ref
    OrderBy("-created_at").
    Limit(10).
    All(ctx)

// users is []*User - direct access!
for _, user := range users {
    fmt.Println(user.Email)  // No type assertion needed
    fmt.Println(user.Name)   // Compile-time safe
}
```

### 3. Type-Safe Field References

```go
// Define field references once (compile-time checked)
var (
    UserEmail    = NewFieldRef[string]("email")
    UserName     = NewFieldRef[string]("name")
    UserID       = NewFieldRef[int]("id")
    UserIsActive = NewFieldRef[bool]("is_active")
)

// Use in queries - type-safe!
users, _ := userManager.Filter(ctx).
    Filter(UserEmail.Contains("@example.com")).  // String methods only on string fields
    Filter(UserID.Gt(100)).                      // Numeric methods only on numeric fields
    Filter(UserIsActive.Eq(true)).               // Bool methods only on bool fields
    All(ctx)
```

### 4. Type-Safe Relationships

```go
// One-to-Many: User has many Posts
userPosts := NewOneToMany[*User, *Post](
    UserID,
    PostUserID,
    postManager,
    "user_id",
)

// Load posts - returns []*Post, type-safe!
posts, err := userPosts.Load(ctx, user)

// posts is []*Post - no type assertions!
for _, post := range posts {
    fmt.Println(post.Title)  // Direct access
}

// Many-to-One: Post belongs to User
postAuthor := NewManyToOne[*Post, *User](
    PostUserID,
    UserID,
    userManager,
    "user_id",
)

// Load author - returns *User, type-safe!
author, err := postAuthor.Load(ctx, post)

// author is *User - no type assertions!
fmt.Println(author.Email)
```

### 5. Type-Safe Admin

```go
// Define field references for admin
var (
    UserEmail    = admin.NewFieldRef[*User]("email")
    UserName     = admin.NewFieldRef[*User]("name")
    UserID       = admin.NewFieldRef[*User]("id")
    UserIsActive = admin.NewFieldRef[*User]("is_active")
)

// Type-safe admin definition - no strings!
var UserAdmin = admin.NewBaseModelAdmin[*User]("User").
    SetListDisplay(UserID, UserEmail, UserName, UserIsActive).  // FieldRef, not strings!
    SetListDisplayLinks(UserEmail).
    SetListEditable(UserIsActive).
    SetSearchFields(UserEmail, UserName).                       // Type-safe!
    SetListFilter(
        admin.FilterSpec[*User]{
            Field: UserIsActive,  // FieldRef, not string!
            Type:  admin.FilterTypeBoolean,
        },
    ).
    SetFields(UserEmail, UserName, UserIsActive).
    SetReadonlyFields(UserID)

// Custom admin with type-safe overrides
type CustomUserAdmin struct {
    *admin.BaseModelAdmin[*User]
}

func (a *CustomUserAdmin) SaveModel(ctx context.Context, obj *User, form interface{}, change bool) error {
    // obj is *User - type-safe!
    _ = obj.Email
    _ = obj.Name
    
    return a.BaseModelAdmin.SaveModel(ctx, obj, form, change)
}

func (a *CustomUserAdmin) GetQueryset(ctx context.Context) models.QuerySet[*User] {
    // Returns type-safe QuerySet[*User]
    return userManager.Filter(ctx)
}
```

### 6. Type-Safe Model Definitions

```go
// Type-safe model builder
var UserModel = NewModelBuilder[*User]("User").
    Int("id").Required().Indexed().
    String("email").Required().Unique().
    String("name").Required().
    Time("created_at").Default(time.Now()).
    Bool("is_active").Default(true).
    Build()

// All field definitions are typed
// String("email") returns *Field[string]
// Int("id") returns *Field[int]
// etc.
```

## Benefits

### ✅ Compile-Time Safety
- Field name typos caught at compile time
- Type mismatches caught at compile time
- Method calls validated at compile time

### ✅ IDE Support
- Full autocomplete for field references
- Type checking in real-time
- Refactoring support

### ✅ No Runtime Errors
- No "field not found" errors
- No type assertion panics
- No string-based lookups

### ✅ Performance
- No reflection in hot paths
- Direct struct field access
- Compiler optimizations

## Migration from String-Based

### Before (String-Based):
```go
// ❌ Runtime errors possible
users, _ := manager.Filter(ctx).
    Filter(map[string]interface{}{"emial": "@example.com"}).  // Typo!
    All(ctx)

user := users[0].(*User)  // Type assertion needed
```

### After (Type-Safe):
```go
// ✅ Compile-time checked
users, _ := userManager.Filter(ctx).
    Filter(UserEmail.Contains("@example.com")).  // Typo = compile error!
    All(ctx)

// users is []*User - no type assertion!
user := users[0]
```

## Best Practices

1. **Define field references once** - Create them as package-level variables
2. **Use generics everywhere** - Manager[T], QuerySet[T], ModelAdmin[T]
3. **Avoid interface{}** - Use specific types or generics
4. **Leverage type inference** - Let Go infer types where possible
5. **Test at compile time** - If it compiles, types are correct

