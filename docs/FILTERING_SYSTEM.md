# Filtering System Documentation

## Overview

The filtering system provides comprehensive, type-safe filtering capabilities for the admin interface, REST APIs, and direct ORM usage. It supports deep relation filtering, boolean tree composition, custom filters, and advanced security features.

## Key Features

- **Filter AST** - Single source of truth for serialization and persistence
- **Deep Relations** - Filter across relationships (e.g., `author__company__country`)
- **Boolean Tree** - Complex AND/OR/NOT composition
- **Runtime Validation** - Uses existing schema system (no codegen)
- **Security Hardening** - Field/lookup whitelisting, cost-based throttles
- **Persisted Filters** - Save, share, and version filter configurations
- **Query Optimization** - Intelligent JOIN vs EXISTS vs subquery selection
- **Custom Filters** - Extensible custom operations
- **Shared Filters** - Reusable filters across API and Admin

## Basic Usage

### Creating a FilterSet

```go
import "github.com/forgego/forge/forge/filter"

type User struct {
    ID       int64
    Username string
    Email    string
    IsActive bool
}

// Create a FilterSet
fs, err := filter.NewFilterSet[User]()
if err != nil {
    log.Fatal(err)
}

// Set base queryset
fs = fs.WithQueryset(User.Objects)
```

### Simple Filtering

```go
// Filter by username
fs.Where("username").Contains("john")

// Filter by exact match
fs.Where("email").Equals("test@example.com")

// Filter by boolean
fs.Where("is_active").Equals(true)

// Filter with IN clause
fs.Where("status").In("active", "pending", "approved")
```

### Boolean Tree Composition

```go
// AND group
fs.AndGroup(func(q *filter.QueryBuilder[User]) {
    q.Where("is_active").Equals(true)
    q.Where("email").EndsWith("@example.com")
})

// OR group
fs.OrGroup(func(q *filter.QueryBuilder[User]) {
    q.Where("username").Contains("admin")
    q.OrFilter("email").Contains("admin")
})
```

### Deep Relation Filtering

```go
// Filter users by author's company's country
fs.Where("author__company__country").Equals("USA")

// Filter by related model field
fs.Where("author__is_active").Equals(true)
```

### Applying Filters

```go
ctx := context.Background()
ast := fs.GetAST()

// Apply AST to queryset
filteredQS, err := fs.ApplyAST(ctx, ast)
if err != nil {
    log.Fatal(err)
}

// Execute query
users, err := filteredQS.All(ctx)
```

## Declarative API

### Defining a FilterSet

```go
type UserFilterSet struct {
    *filter.FilterSet[User]
    Username *filters.CharFilter[User]
    Email    *filters.CharFilter[User]
    IsActive *filters.BooleanFilter[User]
    Author   *filters.RelatedFilter[User, Author]
}

func NewUserFilterSet() *UserFilterSet {
    fs, _ := filter.NewFilterSet[User]()
    return &UserFilterSet{
        FilterSet: fs,
        Username:  filters.NewCharFilter[User]("username").IContains(),
        Email:     filters.NewCharFilter[User]("email").Contains(),
        IsActive:  filters.NewBooleanFilter[User]("is_active"),
    }
}
```

## API Integration

### Using with ViewSets

```go
import "github.com/forgego/forge/forge/api"

// Create FilterSet
fs, _ := filter.NewFilterSet[User]()
fs = fs.WithQueryset(User.Objects)

// Create integration
integration := api.NewFilterSetIntegration(fs)

// In ViewSet List method
func (vs *UserViewSet) List(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Apply filters
    qs, err := integration.ApplyToViewSet(ctx, r, vs.Queryset)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Get filter metadata
    metadata := integration.GetFilterMetadata(r)
    
    // Execute query and return with metadata
    // ...
}
```

### Query Parameters

The API automatically parses query parameters:

```
GET /api/users/?username__contains=john&is_active=true&email__endswith=@example.com
```

Supported formats:
- `field=value` - Exact match
- `field__lookup=value` - Specific lookup (contains, gt, lt, etc.)
- `field__in=val1,val2,val3` - IN clause
- `field__range=min,max` - Range query

## Admin Integration

### Using with Admin ListView

```go
import "github.com/forgego/forge/forge/admin"

// Create FilterSet
fs, _ := filter.NewFilterSet[User]()
fs = fs.WithQueryset(User.Objects)

// Create integration
integration := admin.NewFilterSetIntegration(fs)

// In ListView Render method
func (lv *ListView[User]) Render(ctx context.Context) (*ListData[User], error) {
    // Apply filters from request
    r := getRequestFromContext(ctx)
    if err := integration.ApplyToListView(ctx, r, lv); err != nil {
        return nil, err
    }
    
    // Get filter sidebar data
    sidebarData, _ := integration.GetFilterSidebarData(ctx, lv.queryset)
    
    // Render with filters
    // ...
}
```

## Custom Filters

### Registering a Custom Filter

```go
// Register custom filter handler
filter.RegisterCustom("distance", &filter.CustomFilterHandler{
    ID:   "distance_v1",
    Name: "Distance Filter",
    Handler: func(value interface{}) (orm.Expression, error) {
        // Custom logic to create ORM expression
        // Must return parameterized expression
        return customDistanceExpression(value), nil
    },
    AllowedRoles: []string{"admin"},
    Cost: 5,
})

// Use in FilterSet
customFilter, _ := filter.NewCustomFilter[User]("location", "distance")
fs.AddFilter("location_distance", customFilter)
```

## Persisted Filters

### Saving a Filter

```go
storage := filter.NewInMemoryFilterStorage()

// Save current filter
saved, err := filter.SaveFilter(fs, "Active Users", "Users who are active", storage)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Saved filter ID: %s\n", saved.ID)
```

### Loading a Filter

```go
// Load saved filter
loaded, err := filter.LoadFilter[User](saved.ID, storage)
if err != nil {
    log.Fatal(err)
}

// Apply to queryset
ctx := context.Background()
ast := loaded.GetAST()
filteredQS, err := loaded.ApplyAST(ctx, ast)
```

### Preview Mode

```go
// Preview filter with sample count
count, err := filter.PreviewFilter(fs, ctx, 1000)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Filter would return approximately %d results\n", count)
```

## Security

### Security Configuration

```go
security := filter.NewSecurityConfig()
security.MaxJoinDepth = 3
security.MaxConditions = 50
security.MaxORBranches = 20
security.CostThreshold = 100
security.TimeoutDuration = 30 * time.Second

// Set allowed fields per role
security.AllowedFields["user"] = []string{"username", "email", "is_active"}

// Set allowed lookups per field
security.AllowedLookups["username"] = []string{"exact", "contains", "startswith"}

fs = fs.WithSecurity(security)
```

### Audit Logging

```go
logger := filter.NewDefaultAuditLogger()

// Log filter execution
logger.Log(&filter.AuditLog{
    FilterID:      "filter_123",
    UserID:        "user_456",
    Action:        "execute",
    Cost:          25,
    ExecutionTime: 150 * time.Millisecond,
    Denied:        false,
    Parameters:    filter.MaskParameters(params),
})
```

## Query Optimization

### Cost Estimation

```go
optimizer := filter.NewQueryOptimizer()

ast := fs.GetAST()
plan, err := optimizer.Optimize(ast)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Strategy: %s\n", plan.Strategy)
fmt.Printf("Estimated Cost: %d\n", plan.EstimatedCost)
fmt.Printf("SQL Preview: %s\n", plan.SQLPreview)
```

## Dialect Support

The filtering system supports multiple SQL dialects:

- **PostgreSQL** - Uses ILIKE, JSON operators, similarity
- **MySQL** - Uses LOWER LIKE, JSON functions
- **SQLite** - Uses LOWER LIKE, JSON functions

Dialect adapters are automatically selected based on your database configuration.

## Best Practices

1. **Use Runtime Validation** - Let the schema system validate paths at runtime
2. **Set Security Limits** - Configure MaxJoinDepth, MaxConditions, etc.
3. **Use Filter Presets** - Create helper methods for common filter combinations
4. **Monitor Performance** - Use metrics and alerts to track filter usage
5. **Cache Expensive Filters** - Use persisted filters with caching for heavy queries
6. **Validate Input** - Always validate filter values before applying

## Examples

See `examples/` directory for complete examples of:
- Basic filtering
- Deep relations
- Boolean tree composition
- Custom filters
- Admin integration
- API integration
