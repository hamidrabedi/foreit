# Enterprise Demo - Comprehensive ORM Examples

This package demonstrates all features of the Forge ORM with a complex enterprise schema.

## Structure

- **`models.go`**: Complex schema definitions with relationships, indexes, and constraints
- **`repositories.go`**: Repository pattern with comprehensive data access methods
- **`services.go`**: Business logic layer with complex queries
- **`query_examples.go`**: Comprehensive examples of all ORM features

## Features Demonstrated

### 1. Type-Safe Queries
- Field expressions with compile-time type checking
- Manager CRUD operations
- QuerySet API with full type safety

### 2. Complex Filtering
- Q objects for AND/OR combinations
- String operations (Contains, StartsWith, EndsWith)
- Range queries (BETWEEN)
- IN operator
- NULL checks

### 3. Ordering and Pagination
- Single and multiple field ordering
- Ascending/descending
- Limit and Offset
- First/Last operations

### 4. Aggregates
- Count, Sum, Avg, Min, Max
- Grouped aggregates
- Aggregate with filters

### 5. Field Selection
- Select specific fields
- Only certain fields
- Defer large fields

### 6. Values Queries
- Values (returns maps)
- ValuesList (returns tuples)
- Flat values

### 7. Update Operations
- Simple updates
- UpdateBuilder with expressions
- Increment/Decrement
- Bulk updates

### 8. Relations
- SelectRelated (JOIN queries)
- PrefetchRelated (separate queries)
- Related field filtering

### 9. Set Operations
- Union
- Intersection
- Difference

### 10. Advanced Features
- Distinct queries
- Exclude queries
- Annotations (computed fields)
- Complex business logic queries

## Usage

### Running Query Examples

```go
import (
    "enterprise-demo/app/enterprise"
    "github.com/forgego/forge/db"
)

// Initialize database
db, _ := db.NewDBFromConfig(cfg)

// Run all query examples
enterprise.RunAllQueryExamples(db)
```

### Using Repositories

```go
// Create repository
orgRepo, _ := enterprise.NewOrganizationRepository(db)

// Get active organizations
orgs, _ := orgRepo.GetActiveOrganizations(ctx)

// Search organizations
results, _ := orgRepo.SearchOrganizations(ctx, "search term")

// Get with pagination
orgs, total, _ := orgRepo.GetOrganizationsWithStats(ctx, 10, 0)
```

### Using Services

```go
// Create service
orgService := enterprise.NewOrganizationService(orgRepo)

// Get organization statistics
stats, _ := orgService.GetOrganizationStats(ctx, orgID)
```

## Repository Pattern

Each model has a repository with:
- Basic CRUD operations
- Complex filtering methods
- Search functionality
- Pagination support
- Business-specific queries

## Service Layer

Services provide:
- Business logic
- Multi-model operations
- Statistics and aggregations
- Complex workflows

## Query Examples

The `query_examples.go` file demonstrates:
- All ORM features
- Best practices
- Complex query patterns
- Type-safe query construction
- Dynamic query building

## Models

The enterprise schema includes:
- **Organization**: Companies with subscription tiers
- **SubscriptionTier**: Subscription plans
- **Department**: Organizational departments
- **Employee**: Employees with skills and projects
- **Project**: Projects with tasks
- **Client**: External clients
- **Task**: Project tasks
- **Skill**: Employee skills (many-to-many)

All models have:
- Complex relationships (ForeignKey, ManyToMany)
- Indexes for performance
- Constraints for data integrity
- Lifecycle hooks support

## Type Safety

All queries are type-safe:
- Field expressions are typed
- Compile-time checking
- No runtime type errors
- IDE autocomplete support

## Performance

Examples demonstrate:
- Efficient query construction
- Field selection optimization
- Pagination for large datasets
- Index usage
- Relation optimization
