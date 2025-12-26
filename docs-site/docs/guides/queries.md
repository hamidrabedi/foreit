---
sidebar_position: 2
---

# Queries Guide

The forge QuerySet API provides a powerful, type-safe way to query your database. It's similar to Django's QuerySet but with Go's type safety.

## Basic Queries

### Get All Objects

```go
ctx := context.Background()

// Get all users
users, err := User.Objects.All(ctx)
```

### Get a Single Object

```go
// Get by ID
user, err := User.Objects.Get(ctx, 1)

// Get first object
user, err := User.Objects.First(ctx)

// Get last object
user, err := User.Objects.Last(ctx)
```

### Filtering

```go
// Filter by field
users, err := User.Objects.
    Filter(User.Fields.IsActive.Equals(true)).
    All(ctx)

// Multiple filters (AND)
users, err := User.Objects.
    Filter(User.Fields.IsActive.Equals(true)).
    Filter(User.Fields.IsStaff.Equals(true)).
    All(ctx)

// Exclude
users, err := User.Objects.
    Exclude(User.Fields.IsDeleted.Equals(true)).
    All(ctx)
```

## Field Expressions

Field expressions provide type-safe field access:

### Equality

```go
User.Fields.Username.Equals("john")
User.Fields.Age.NotEquals(18)
```

### Comparison

```go
User.Fields.Age.Greater(18)
User.Fields.Age.GreaterOrEqual(18)
User.Fields.Age.Less(65)
User.Fields.Age.LessOrEqual(65)
```

### Null Checks

```go
User.Fields.LastLogin.IsNull()
User.Fields.LastLogin.IsNotNull()
```

### Membership

```go
User.Fields.Status.In("active", "pending", "approved")
User.Fields.Status.NotIn("deleted", "banned")
```

### String Operations

```go
User.Fields.Username.Contains("john")
User.Fields.Username.StartsWith("admin")
User.Fields.Username.EndsWith(".com")
User.Fields.Username.IContains("JOHN")  // Case-insensitive
```

### Range

```go
User.Fields.Age.Range(18, 65)
User.Fields.CreatedAt.Range(startDate, endDate)
```

## Complex Queries

### Combining Conditions

```go
// AND
users, err := User.Objects.
    Filter(
        User.Fields.IsActive.Equals(true).
            And(User.Fields.DateJoined.Greater(lastMonth)),
    ).
    All(ctx)

// OR
users, err := User.Objects.
    Filter(
        User.Fields.IsActive.Equals(true).
            Or(User.Fields.IsStaff.Equals(true)),
    ).
    All(ctx)

// NOT
users, err := User.Objects.
    Filter(
        User.Fields.IsActive.Equals(true).Not(),
    ).
    All(ctx)
```

### Ordering

```go
// Single field
users, err := User.Objects.OrderBy("username").All(ctx)

// Multiple fields
users, err := User.Objects.OrderBy("username", "-date_joined").All(ctx)

// Descending (use - prefix)
users, err := User.Objects.OrderBy("-date_joined").All(ctx)
```

### Limiting and Offsetting

```go
// Limit
users, err := User.Objects.Limit(10).All(ctx)

// Offset
users, err := User.Objects.Offset(20).All(ctx)

// Pagination
users, err := User.Objects.
    Limit(10).
    Offset(20).
    All(ctx)
```

### Distinct

```go
users, err := User.Objects.Distinct().All(ctx)
```

## Aggregations

### Count

```go
count, err := User.Objects.
    Filter(User.Fields.IsActive.Equals(true)).
    Count(ctx)
```

### Exists

```go
exists, err := User.Objects.
    Filter(User.Fields.Username.Equals("john")).
    Exists(ctx)
```

### Aggregates

```go
import "github.com/forgego/forge/pkg/query/aggregates"

result, err := User.Objects.
    Aggregate(
        aggregates.Count("id"),
        aggregates.Avg("age"),
        aggregates.Max("date_joined"),
        aggregates.Min("date_joined"),
        aggregates.Sum("score"),
    ).
    Get(ctx)
```

## Field Selection

### Select Specific Fields

```go
// Only these fields
users, err := User.Objects.
    Only("username", "email").
    All(ctx)

// Exclude these fields
users, err := User.Objects.
    Defer("password", "secret").
    All(ctx)
```

## Relations

### Accessing Relations

```go
// Get user with posts
user, err := User.Objects.Get(ctx, 1)
posts := user.Posts  // []*Post

// Prefetch related objects (efficient)
users, err := User.Objects.
    PrefetchRelated("posts", "profile").
    All(ctx)

// Select related (JOIN)
users, err := User.Objects.
    SelectRelated("profile").
    All(ctx)
```

### Filtering by Relations

```go
// Get posts by author
posts, err := Post.Objects.
    Filter(Post.Fields.Author.Equals(userID)).
    All(ctx)

// Get users with posts
users, err := User.Objects.
    Filter(User.Fields.Posts.Count().Greater(0)).
    All(ctx)
```

## Updates

### Update Single Object

```go
user.Username = "newusername"
err := User.Objects.Update(ctx, user)
```

### Bulk Update

```go
affected, err := User.Objects.
    Filter(User.Fields.IsActive.Equals(false)).
    Update(ctx, map[string]interface{}{
        "is_active": true,
    })
```

## Deletion

### Delete Single Object

```go
err := User.Objects.Delete(ctx, user)
```

### Bulk Delete

```go
affected, err := User.Objects.
    Filter(User.Fields.IsDeleted.Equals(true)).
    Delete(ctx)
```

## Transactions

### Using Transactions

```go
import "github.com/forgego/forge/pkg/db"

err := db.WithTx(ctx, func(tx *db.Tx) error {
    user := &User{Username: "john"}
    if err := User.Objects.Create(ctx, user); err != nil {
        return err
    }
    
    post := &Post{Author: user, Title: "Hello"}
    return Post.Objects.Create(ctx, post)
})
```

## Query Optimization

### Use SelectRelated for Foreign Keys

```go
// Good: Uses JOIN
posts, err := Post.Objects.
    SelectRelated("author").
    All(ctx)

// Bad: N+1 queries
posts, err := Post.Objects.All(ctx)
for _, post := range posts {
    author := post.Author  // Separate query for each post
}
```

### Use PrefetchRelated for Many Relations

```go
// Good: Two queries total
users, err := User.Objects.
    PrefetchRelated("posts").
    All(ctx)

// Bad: N+1 queries
users, err := User.Objects.All(ctx)
for _, user := range users {
    posts := user.Posts  // Separate query for each user
}
```

### Use Only/Defer for Large Models

```go
// Good: Only fetch needed fields
users, err := User.Objects.
    Only("username", "email").
    All(ctx)

// Bad: Fetch all fields including large text fields
users, err := User.Objects.All(ctx)
```

## Common Patterns

### Pagination

```go
func GetUsers(page, pageSize int) ([]*User, int64, error) {
    ctx := context.Background()
    
    // Get total count
    total, err := User.Objects.Count(ctx)
    if err != nil {
        return nil, 0, err
    }
    
    // Get page
    users, err := User.Objects.
        Limit(pageSize).
        Offset((page - 1) * pageSize).
        OrderBy("-date_joined").
        All(ctx)
    
    return users, total, err
}
```

### Search

```go
func SearchUsers(query string) ([]*User, error) {
    ctx := context.Background()
    
    return User.Objects.
        Filter(
            User.Fields.Username.Contains(query).
                Or(User.Fields.Email.Contains(query)),
        ).
        All(ctx)
}
```

## Next Steps

- [API Reference](reference/queryset) - Complete QuerySet API
- [Manager Reference](reference/manager) - Manager methods
- [Field Reference](reference/fields) - Field expression methods

