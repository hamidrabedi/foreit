# forge API Reference

## QuerySet API

### Type-Safe Querying

```go
// Get all active users
users, err := User.Objects.Filter(
    User.Fields.IsActive.Equals(true),
).All(ctx)

// Complex query
users, err := User.Objects.Filter(
    User.Fields.IsActive.Equals(true).And(
        User.Fields.DateJoined.Greater(time.Now().AddDate(0, -1, 0)).Or(
            User.Fields.IsStaff.Equals(true),
        ),
    ),
).OrderBy("-date_joined").Limit(10).All(ctx)
```

### FieldExpr Methods

```go
// Equality
User.Fields.Username.Equals("john")
User.Fields.Age.NotEquals(18)

// Comparison
User.Fields.Age.Greater(18)
User.Fields.Age.GreaterOrEqual(18)
User.Fields.Age.Less(65)
User.Fields.Age.LessOrEqual(65)

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

### QuerySet Methods

```go
// Filtering
qs := User.Objects.Filter(User.Fields.IsActive.Equals(true))
qs = qs.Exclude(User.Fields.IsDeleted.Equals(true))

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
activeUsers := User.Objects.Filter(User.Fields.IsActive.Equals(true))
staffUsers := User.Objects.Filter(User.Fields.IsStaff.Equals(true))

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

## QueryExpr API

```go
// Build complex queries
q := User.Fields.IsActive.Equals(true).
    And(User.Fields.IsStaff.Equals(true)).
    Or(User.Fields.DateJoined.Greater(lastMonth))

// Negation
q = q.Not()

// Combine queries
q1 := User.Fields.IsActive.Equals(true)
q2 := User.Fields.IsStaff.Equals(true)
q := q1.And(q2)
q = q1.Or(q2)
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
    Post.Fields.Author.Equals(user.ID),
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

