# Filters Package

Django-filters-like functionality with full type safety.

## Features

- ✅ **Type-safe filters** - FieldRef-based, no strings
- ✅ **Multiple lookups** - exact, contains, icontains, in, range, etc.
- ✅ **Composable** - Works with type-safe QuerySet
- ✅ **No reflection** - Compile-time checked

## Usage

### Define FilterSet

```go
import (
    "github.com/gogo/pkg/models"
    "github.com/gogo/pkg/filters"
)

// Define field references
var (
    UserEmail = models.NewFieldRef[string]("email")
    UserName  = models.NewFieldRef[string]("name")
    UserAge   = models.NewFieldRef[int]("age")
)

// Create type-safe filter set
filterSet := filters.NewFilterSet[*User](userManager).
    AddFilter("email", UserEmail, 
        filters.LookupExact, 
        filters.LookupIContains,
        filters.LookupStartswith,
    ).
    AddFilter("name", UserName,
        filters.LookupExact,
        filters.LookupIContains,
    ).
    AddFilter("age", UserAge,
        filters.LookupExact,
        filters.LookupGt,
        filters.LookupGte,
        filters.LookupLt,
        filters.LookupLte,
        filters.LookupRange,
    )
```

### Use FilterSet

```go
// Filter from request parameters
params := map[string]interface{}{
    "email__icontains": "@example.com",
    "age__gte": 18,
    "age__lte": 65,
}

qs := filterSet.Filter(ctx, params)

// Returns type-safe QuerySet[*User]
users, err := qs.All(ctx)

// users is []*User - type-safe!
for _, user := range users {
    fmt.Println(user.Email)
}
```

### Available Lookups

- `exact` - Exact match
- `iexact` - Case-insensitive exact match
- `contains` - Contains substring
- `icontains` - Case-insensitive contains
- `in` - Value in list
- `gt`, `gte`, `lt`, `lte` - Numeric comparisons
- `startswith`, `istartswith` - Starts with
- `endswith`, `iendswith` - Ends with
- `range` - Value range (for dates/numbers)
- `date`, `year`, `month`, `day` - Date lookups
- `isnull` - IS NULL / IS NOT NULL
- `regex`, `iregex` - Regular expression match

## Integration with Admin

FilterSets can be used in admin for list filtering:

```go
var UserAdmin = admin.NewBaseModelAdmin[*User]("User").
    SetListFilter(
        admin.FilterSpec[*User]{
            Field: UserEmail,
            Type:  admin.FilterTypeList,
            Lookups: []admin.Lookup{
                admin.LookupExact,
                admin.LookupIContains,
            },
        },
    )
```

