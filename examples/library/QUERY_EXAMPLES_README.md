# Query Examples for Library Demo

This file demonstrates comprehensive usage of the `forge` query system with all types of queries.

## Running the Examples

The examples are automatically run when you start the library server:

```bash
cd examples/library
go run .
```

The examples will run if a database connection is available.

## Examples Included

### 1. Simple Queries
- Get all records
- Get by ID
- Simple filters (price > 10, available = true)

### 2. Complex Filtering
- Multiple conditions with AND/OR
- String operations (Contains, StartsWith, EndsWith)
- IN operator
- Range queries (BETWEEN)
- NULL checks
- Exclude queries

### 3. Ordering and Limiting
- Order by single field
- Order by multiple fields
- Pagination (Limit + Offset)
- First and Last records

### 4. Aggregation
- Count queries
- Count with filters
- Exists checks
- Aggregate functions (Avg, Sum, Max, Min)

### 5. Annotations
- Computed fields
- Expression-based annotations

### 6. Values and ValuesList
- Values as maps
- Values as tuples
- Flat single-field lists

### 7. Update Operations
- Simple updates
- Update with expressions
- Increment/Decrement operations

### 8. Delete Operations
- Delete with filters

### 9. Create Operations
- Creating new records

### 10. Complex Customized Queries
- Multi-condition queries
- Queries with annotations
- Chained operations
- Distinct queries

## Key Features Demonstrated

1. **Type Safety**: All field access is type-checked at compile time
2. **Expression System**: Complex queries using Q objects
3. **Chainable API**: All QuerySet methods return new instances
4. **SQL Safety**: Parameter binding prevents SQL injection

## Usage Pattern

```go
// 1. Create manager
manager, _ := query.NewManager[models.Book]("books")
manager.SetDB(database)

// 2. Get field accessor
fa, _ := manager.GetFieldAccessor()

// 3. Build type-safe queries
priceField := fa.Field[float64]("price")
availableField := fa.Field[bool]("available")

// 4. Execute queries
books, err := manager.Filter(
    query.NewQ(priceField.Gt(10.0)).
        And(query.NewQ(availableField.Eq(true))),
).All(ctx)
```

## Notes

- All examples use the `QuerySet` and `Manager` APIs
- The old `QuerySet` and `Manager` APIs still work for backward compatibility
- Examples are designed to be safe (delete operations are commented out)
- Some examples may need database setup to run successfully
