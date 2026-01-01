# forge API Reference

## QuerySet API

### Type-Safe Querying

```go
// Get all active users (direct filter)
users, err := User.Objects.Filter(
    User.Fields.IsActive.Eq(true),
).All(ctx)

// Complex query using And/Or
users, err := User.Objects.Filter(
    orm.And(
        User.Fields.IsActive.Eq(true),
        orm.Or(
            User.Fields.DateJoined.Gt(time.Now().AddDate(0, -1, 0)),
            User.Fields.IsStaff.Eq(true),
        ),
    ),
).OrderBy("-date_joined").Limit(10).All(ctx)
```

### Field Methods (Type-Safe)

```go
// Equality
User.Fields.Username.Eq("john")
User.Fields.Age.Ne(18)

// Comparison
User.Fields.Age.Gt(18)
User.Fields.Age.Gte(18)
User.Fields.Age.Lt(65)
User.Fields.Age.Lte(65)

// Null checks
User.Fields.LastLogin.IsNull()
User.Fields.LastLogin.IsNotNull()

// Membership
User.Fields.Status.In("active", "pending", "approved")
User.Fields.Status.NotIn("deleted", "banned")

// String operations
User.Fields.Username.Contains("john")
User.Fields.Username.StartsWith("admin")
User.Fields.Username.EndsWith(".com")
User.Fields.Username.IContains("JOHN")  // Case-insensitive

// Range
User.Fields.Age.Range(18, 65)

// Date operations
User.Fields.DateJoined.Year(2024)
User.Fields.DateJoined.Month(12)
User.Fields.DateJoined.Day(25)
```

### Runtime Field References

When you don't have type-safe fields, use `F()` or `FieldRef()`:

```go
// Django-style (short)
qs.Filter(orm.F("age").Gt(18))

// Explicit alternative
qs.Filter(orm.FieldRef("age").Gt(18))

// SQL-like Where clause
qs.Filter(orm.Where("age", orm.OpGreater, 18))
qs.Filter(orm.Where("name", orm.OpEquals, "John"))
```

### Boolean Expression Combiners

```go
// And - combine expressions with AND
qs.Filter(orm.And(
    User.Fields.Name.Eq("John"),
    User.Fields.Age.Gt(18),
))

// Or - combine expressions with OR
qs.Filter(orm.Or(
    User.Fields.Age.Gt(18),
    User.Fields.Role.Eq("admin"),
))

// Not - negate an expression
qs.Exclude(orm.Not(User.Fields.Age.Gt(65)))

// Complex combinations
qs.Filter(orm.And(
    User.Fields.Name.Eq("John"),
    orm.Or(
        User.Fields.Age.Gt(18),
        User.Fields.Role.Eq("admin"),
    ),
))
```

### QuerySet Methods

```go
// Filtering
qs := User.Objects.Filter(User.Fields.IsActive.Eq(true))
qs = qs.Exclude(User.Fields.IsDeleted.Eq(true))

// Ordering
qs = qs.OrderBy("username", "-date_joined")  // - means DESC

// Pagination
qs = qs.Limit(10).Offset(20)

// Distinct
qs = qs.Distinct()

// Field selection
qs = qs.Select("id", "username", "email")
qs = qs.Only("username", "email")  // Only these fields
qs = qs.Defer("password", "secret")  // Exclude these fields

// Relations
qs = qs.SelectRelated("profile", "groups")  // JOIN
qs = qs.PrefetchRelated("posts", "comments")  // Separate query

// Aggregation
result, err := qs.Aggregate(
    aggregates.Count("id"),
    aggregates.Avg("age"),
    aggregates.Max("date_joined"),
).Get(ctx)

// Annotations
qs = qs.Annotate(
    annotations.Value("full_name", User.Fields.FirstName + " " + User.Fields.LastName),
)

// Execution
users, err := qs.All(ctx)           // []*User
user, err := qs.Get(ctx)            // *User (expects exactly one)
user, err := qs.First(ctx)          // *User (first result)
user, err := qs.Last(ctx)           // *User (last result)
count, err := qs.Count(ctx)         // int64
exists, err := qs.Exists(ctx)       // bool

// Updates
affected, err := qs.Update(ctx, map[string]interface{}{
    "is_active": false,
})

// Bulk operations
err := qs.BulkUpdate(ctx, []map[string]interface{}{
    {"id": 1, "is_active": true},
    {"id": 2, "is_active": false},
})

err := qs.BulkCreate(ctx, []*User{user1, user2})

// Deletion
affected, err := qs.Delete(ctx)

// Set operations
activeUsers := User.Objects.Filter(User.Fields.IsActive.Eq(true))
staffUsers := User.Objects.Filter(User.Fields.IsStaff.Eq(true))

allUsers := activeUsers.Union(staffUsers)
commonUsers := activeUsers.Intersection(staffUsers)
onlyActive := activeUsers.Difference(staffUsers)
```

### Dynamic Querying

```go
// String-based queries
users, err := User.Objects.FilterDynamic(
    query.Q("is_active", true),
    query.Q("age__gt", 18),  // age > 18
    query.Q("username__contains", "john"),
).All(ctx)

// Complex dynamic queries
users, err := User.Objects.FilterDynamic(
    query.Q("is_active", true).And(
        query.Q("date_joined__gte", lastMonth),
    ),
).All(ctx)
```

## Manager API

```go
// Get by ID
user, err := User.Objects.Get(ctx, 1)

// Get all
users, err := User.Objects.All(ctx)

// Create
user := &User{
    Username: "john",
    Email: "john@example.com",
}
err := User.Objects.Create(ctx, user)

// Update
user.Username = "jane"
err := User.Objects.Update(ctx, user)

// Save (create or update)
err := user.Save(ctx)  // Instance method

// Delete
err := User.Objects.Delete(ctx, user)
err := user.Delete(ctx)  // Instance method

// Filter returns QuerySet
qs := User.Objects.Filter(User.Fields.IsActive.Equals(true))
```

## Expression API

```go
// Build complex queries using And/Or/Not functions
expr := orm.And(
    User.Fields.IsActive.Eq(true),
    User.Fields.IsStaff.Eq(true),
)

// Or combine
expr = orm.Or(
    User.Fields.IsActive.Eq(true),
    User.Fields.DateJoined.Gt(lastMonth),
)

// Negation
expr = orm.Not(User.Fields.Age.Gt(65))

// Complex combinations
expr = orm.And(
    User.Fields.Name.Eq("John"),
    orm.Or(
        User.Fields.Age.Gt(18),
        User.Fields.Role.Eq("admin"),
    ),
)

// Use in queries
qs.Filter(expr)
```

## Aggregates

```go
// Count
count := aggregates.Count("id")

// Sum
total := aggregates.Sum("price")

// Average
avg := aggregates.Avg("age")

// Max
max := aggregates.Max("date_joined")

// Min
min := aggregates.Min("date_joined")

// Usage
result, err := User.Objects.Aggregate(
    aggregates.Count("id"),
    aggregates.Avg("age"),
).Get(ctx)
```

## Values and ValuesList

```go
// Values - returns []map[string]interface{}
results, err := User.Objects.Filter(...).Values("username", "email")

// ValuesList - returns [][]interface{}
results, err := User.Objects.Filter(...).ValuesList("username", "email")

// Flat ValuesList - returns []interface{} (single field)
usernames, err := User.Objects.Filter(...).ValuesList("username", flat=true)
```

## Relations

```go
// Access relations
user, _ := User.Objects.Get(ctx, 1)
posts := user.Posts  // []*Post
profile := user.Profile  // *UserProfile

// Prefetch relations (efficient)
users, _ := User.Objects.PrefetchRelated("posts", "profile").All(ctx)

// Select related (JOIN)
users, _ := User.Objects.SelectRelated("profile").All(ctx)

// Filter by relation
posts, _ := Post.Objects.Filter(
    Post.Fields.Author.Eq(user.ID),
).All(ctx)
```

## Transactions

```go
// With transaction
err := db.WithTx(ctx, func(tx *db.Tx) error {
    user := &User{Username: "john"}
    if err := User.Objects.Create(ctx, user); err != nil {
        return err
    }
    
    post := &Post{Author: user, Title: "Hello"}
    return Post.Objects.Create(ctx, post)
})

// Manual transaction
tx, _ := db.BeginTx(ctx, nil)
defer tx.Rollback()

// Use tx in queries
// ... operations ...

tx.Commit()
```

## Validation

```go
// Validate model
err := user.Validate()

// Custom validation in hooks
func (User) Hooks() schema.ModelHooks {
    return schema.ModelHooks{
        Clean: func(ctx context.Context, instance interface{}) error {
            user := instance.(*User)
            if user.Username == "" {
                return errors.New("username is required")
            }
            return nil
        },
    }
}
```

## REST API Features (Planned)

> **Note:** The REST API system structure is in place, but full implementation is pending QuerySet integration.

### Pagination

ViewSets will support automatic pagination through query parameters:

```
GET /api/v1/posts?page=1&page_size=20
```

**Paginated Response Format:**
```json
{
  "count": 100,
  "next": "http://example.com/api/v1/posts?page=2&page_size=20",
  "previous": null,
  "results": [...]
}
```

### Filtering

Filter sets will allow filtering by field:

```
GET /api/v1/posts?published=true&title__icontains=django
```

**Filter Functions:**
- `Exact` - Exact match
- `IExact` - Case-insensitive exact match
- `Contains` - Contains substring
- `IContains` - Case-insensitive contains
- `StartsWith`, `EndsWith`
- `In` - Value in list
- `Range` - Range filter
- `IsNull`, `IsNotNull`

### Ordering

Use the `ordering` parameter:

```
GET /api/v1/posts?ordering=-created_at,title
```

### Search

Full-text search via `search` parameter:

```
GET /api/v1/posts?search=django tutorial
```

**Implementation Status:**
- ✅ Pagination structure complete
- ✅ Filter structure complete
- ✅ Ordering structure complete
- ✅ Search structure complete
- 🚧 QuerySet integration pending
- 🚧 Database query execution pending

