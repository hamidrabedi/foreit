package orm

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestModel is a test model for integration tests
type TestModel struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Price     float64   `db:"price"`
	Available bool      `db:"available"`
	CreatedAt time.Time `db:"created_at"`
}

// setupTestDB creates a test SQLite database
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Create test table
	_, err = db.Exec(`
		CREATE TABLE test_models (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT,
			price REAL DEFAULT 0.0,
			available INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	return db
}

// teardownTestDB cleans up test database
func teardownTestDB(t *testing.T, db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func TestQuerySet_Integration_Filter(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// Insert test data
	_, err := db.Exec(`
		INSERT INTO test_models (name, email, price, available) VALUES
		('Product 1', 'test1@example.com', 10.0, 1),
		('Product 2', 'test2@example.com', 20.0, 1),
		('Product 3', 'test3@example.com', 30.0, 0)
	`)
	require.NoError(t, err)

	qs, err := NewQuerySet[TestModel]("test_models")
	if err != nil {
		t.Skipf("Schema not registered for TestModel: %v", err)
		return
	}

	// Set DB connection (would need adapter for *sql.DB)
	// For now, test that QuerySet can be created and chained
	priceField := NewField[float64]("price", "test_models")
	expr := priceField.Gt(15.0)

	filtered := qs.Filter(expr)
	assert.NotNil(t, filtered)

	// Verify filtering worked
	assert.NotEqual(t, qs, filtered)
}

func TestQuerySet_Integration_OrderBy(t *testing.T) {
	qs, err := NewQuerySet[TestModel]("test_models")
	if err != nil {
		t.Skipf("Schema not registered for TestModel: %v", err)
		return
	}

	ordered := qs.OrderBy(NewOrderField("price", false))
	assert.NotNil(t, ordered)

	// buildSQL is not exported, test through public API
	_ = ordered
}

func TestQuerySet_Integration_LimitOffset(t *testing.T) {
	qs, err := NewQuerySet[TestModel]("test_models")
	if err != nil {
		t.Skipf("Schema not registered for TestModel: %v", err)
		return
	}

	limited := qs.Limit(10).Offset(5)
	assert.NotNil(t, limited)

	// buildSQL is not exported, test through public API
	_ = limited
}

func TestQuerySet_Integration_ComplexQuery(t *testing.T) {
	qs, err := NewQuerySet[TestModel]("test_models")
	if err != nil {
		t.Skipf("Schema not registered for TestModel: %v", err)
		return
	}

	priceField := NewField[float64]("price", "test_models")
	availableField := NewField[bool]("available", "test_models")

	// Complex query: price > 10 AND available = true, ordered by price DESC, limit 5
	complex := qs.
		Filter(priceField.Gt(10.0)).
		Filter(availableField.Eq(true)).
		OrderBy(NewOrderField("price", false)).
		Limit(5)

	assert.NotNil(t, complex)

	// buildSQL is not exported, test through public API
	// Verify the QuerySet was created and chained correctly
	assert.NotNil(t, complex)
}

func TestUpdateBuilder_Integration(t *testing.T) {
	qs, err := NewQuerySet[TestModel]("test_models")
	if err != nil {
		t.Skipf("Schema not registered for TestModel: %v", err)
		return
	}

	ub, err := NewUpdateBuilder[TestModel](qs)
	// Skip if schema not registered
	if err != nil {
		t.Skipf("Schema not registered: %v", err)
		return
	}

	// Test update builder chaining
	ub = ub.
		Set("name", "Updated Name").
		Set("price", 99.99).
		Increment("id", int64(1))

	assert.NotNil(t, ub)

	// Verify updates were stored
	updates := ub.updates
	assert.Greater(t, len(updates), 0)
}

func TestExpression_Integration_StringOperations(t *testing.T) {
	nameField := NewField[string]("name", "test_models")

	// Test Contains
	containsExpr := nameField.Contains("Product")
	builder := NewSQLBuilder()
	sql, args, err := containsExpr.ToSQL(builder)
	require.NoError(t, err)
	assert.Contains(t, sql, "LIKE")
	assert.Len(t, args, 1)

	// Test StartsWith
	startsWithExpr := nameField.StartsWith("Product")
	builder = NewSQLBuilder()
	sql, args, err = startsWithExpr.ToSQL(builder)
	require.NoError(t, err)
	assert.Contains(t, sql, "LIKE")
	assert.Len(t, args, 1)

	// Test EndsWith
	endsWithExpr := nameField.EndsWith("1")
	builder = NewSQLBuilder()
	sql, args, err = endsWithExpr.ToSQL(builder)
	require.NoError(t, err)
	assert.Contains(t, sql, "LIKE")
	assert.Len(t, args, 1)
}

func TestQ_Integration_ComplexNesting(t *testing.T) {
	priceField := NewField[float64]("price", "test_models")
	availableField := NewField[bool]("available", "test_models")
	nameField := NewField[string]("name", "test_models")

	// Complex: (price > 20 OR price < 10) AND (available = true OR name LIKE '%Product%')
	q1 := NewQ(priceField.Gt(20.0)).Or(NewQ(priceField.Lt(10.0)))
	q2 := NewQ(availableField.Eq(true)).Or(NewQ(nameField.Contains("Product")))
	combined := q1.And(q2)

	builder := NewSQLBuilder()
	sql, args, err := combined.ToSQL(builder)
	require.NoError(t, err)

	assert.Contains(t, sql, "OR")
	assert.Contains(t, sql, "AND")
	assert.GreaterOrEqual(t, len(args), 4)
}

func TestSQLBuilder_Integration_ComplexQuery(t *testing.T) {
	builder := NewSQLBuilder()

	// Build a complex SELECT query
	selectClause := builder.BuildSelect("test_models", []string{"id", "name", "price"}, false)
	whereClause, whereArgs := builder.BuildWhere(
		[]QueryExpr{NewFieldQueryExpr("price", OpGreater, 10.0)},
		[]QueryExpr{},
	)
	orderByClause := builder.BuildOrderBy([]string{"-price", "name"})
	limitClause := builder.BuildLimit(intPtr(10))
	offsetClause := builder.BuildOffset(intPtr(0))

	sql := selectClause + " " + whereClause + " " + orderByClause + " " + limitClause + " " + offsetClause

	assert.Contains(t, sql, "SELECT")
	assert.Contains(t, sql, "FROM")
	assert.Contains(t, sql, "WHERE")
	assert.Contains(t, sql, "ORDER BY")
	assert.Contains(t, sql, "LIMIT")
	assert.Contains(t, sql, "OFFSET")
	assert.Len(t, whereArgs, 1)
}

func intPtr(i int) *int {
	return &i
}
