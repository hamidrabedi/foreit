package orm

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/forgego/forge/db"
	forgeerrors "github.com/forgego/forge/errors"
	"github.com/forgego/forge/internal/testutils"
	"github.com/forgego/forge/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type aggregateTestItem struct {
	schema.BaseSchema
	ID     int64   `db:"id"`
	Amount float64 `db:"amount"`
	Kind   string  `db:"kind"`
}

type aggregateBlobItem struct {
	schema.BaseSchema
	ID   int64  `db:"id"`
	Text string `db:"text"`
	Data []byte `db:"data"`
}

func (aggregateBlobItem) Meta() schema.Meta { return schema.Meta{TableName: "aggregate_blob_items"} }

func (aggregateBlobItem) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64Field("id", schema.Primary(), schema.AutoIncrement()),
		schema.StringField("text"),
		schema.BytesField("data"),
	}
}

type aggregateValuesWrapper[T any] struct{ QuerySet[T] }

func (aggregateValuesWrapper[T]) AggregateValues(context.Context, ...Aggregate) (map[string]any, error) {
	return map[string]any{"sentinel": true}, nil
}

func (aggregateTestItem) Meta() schema.Meta { return schema.Meta{TableName: "aggregate_test_items"} }

func (aggregateTestItem) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64Field("id", schema.Primary(), schema.AutoIncrement()),
		schema.Float64Field("amount"),
		schema.StringField("kind"),
	}
}

func assertAggregateValues(t *testing.T, database *db.DB, table string) {
	t.Helper()
	ctx := context.Background()
	querySet, err := NewQuerySet[aggregateTestItem](table)
	require.NoError(t, err)
	base := querySet.SetDB(database).(*BaseQuerySet[aggregateTestItem])

	values, err := base.AggregateValues(ctx, Count("id"), Sum("amount"), Avg("amount"), Min("amount"), Max("amount"))
	require.NoError(t, err)
	assert.Equal(t, int64(3), values["count"])
	assert.Equal(t, 600.0, values["sum"])
	assert.Equal(t, 200.0, values["avg"])
	assert.Equal(t, 100.0, values["min"])
	assert.Equal(t, 300.0, values["max"])

	filtered := base.Filter(F("kind").Eq("a")).(*BaseQuerySet[aggregateTestItem])
	values, err = filtered.AggregateValues(ctx, Count("id"), Sum("amount"), Avg("amount"), Min("amount"), Max("amount"))
	require.NoError(t, err)
	assert.Equal(t, int64(2), values["count"])
	assert.Equal(t, 300.0, values["sum"])
	assert.Equal(t, 150.0, values["avg"])
	assert.Equal(t, 100.0, values["min"])
	assert.Equal(t, 200.0, values["max"])
}

func TestAggregateValuesSQLite(t *testing.T) {
	database := setupRelationTestDB(t)
	defer database.Close()
	_, err := database.Exec(`CREATE TABLE aggregate_test_items (id INTEGER PRIMARY KEY, amount REAL, kind TEXT)`)
	require.NoError(t, err)
	_, err = database.Exec(`INSERT INTO aggregate_test_items (id, amount, kind) VALUES (1, 100, 'a'), (2, 200, 'a'), (3, 300, 'b')`)
	require.NoError(t, err)

	assertAggregateValues(t, database, "aggregate_test_items")

	querySet, err := NewQuerySet[aggregateTestItem]("aggregate_test_items")
	require.NoError(t, err)
	values, err := querySet.SetDB(database).(*BaseQuerySet[aggregateTestItem]).Filter(F("id").Eq(-1)).(*BaseQuerySet[aggregateTestItem]).AggregateValues(context.Background(), Count("id"), Sum("amount"), Avg("amount"), Min("amount"), Max("amount"))
	require.NoError(t, err)
	assert.Equal(t, int64(0), values["count"])
	assert.Nil(t, values["sum"])
	assert.Nil(t, values["avg"])
	assert.Nil(t, values["min"])
	assert.Nil(t, values["max"])
}

func TestAggregateValuesPostgreSQL(t *testing.T) {
	sqlDB := testutils.SetupTestDB(t)
	t.Cleanup(func() { _ = sqlDB.Close() })
	database := &db.DB{DB: sqlDB, Driver: "postgres"}
	table := fmt.Sprintf("aggregate_test_items_%d", time.Now().UnixNano())
	escapedTable := EscapeIdentifier(table)
	_, err := database.Exec(fmt.Sprintf(`CREATE TABLE %s (id BIGINT PRIMARY KEY, amount DOUBLE PRECISION, kind TEXT)`, escapedTable))
	require.NoError(t, err)
	_, err = database.Exec(fmt.Sprintf(`INSERT INTO %s (id, amount, kind) VALUES (1, 100, 'a'), (2, 200, 'a'), (3, 300, 'b')`, escapedTable))
	require.NoError(t, err)
	assertAggregateValues(t, database, table)
}

func TestAggregateValuesRejectedBeforeSQL(t *testing.T) {
	database := setupRelationTestDB(t)
	require.NoError(t, database.Close())
	querySet, err := NewQuerySet[aggregateTestItem]("aggregate_test_items")
	require.NoError(t, err)
	base := querySet.SetDB(database).(*BaseQuerySet[aggregateTestItem])

	for _, aggs := range [][]Aggregate{
		{{Func: "SUM", Field: "unknown"}},
		{{Func: "STDDEV", Field: "amount"}},
		nil,
		{{Name: "same", Func: "SUM", Field: "amount"}, {Name: "same", Func: "AVG", Field: "amount"}},
	} {
		_, err := base.AggregateValues(context.Background(), aggs...)
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "database is closed")
	}
	_, err = base.AggregateValues(context.Background(), Aggregate{Func: "SUM", Field: "unknown"})
	assert.True(t, forgeerrors.IsInvalidInput(err))
	_, err = base.AggregateValues(context.Background(), Aggregate{Func: "STDDEV", Field: "amount"})
	assert.True(t, forgeerrors.IsNotImplemented(err))
}

func TestAggregateValuesPackageHelper(t *testing.T) {
	database := setupRelationTestDB(t)
	defer database.Close()
	_, err := database.Exec(`CREATE TABLE aggregate_test_items (id INTEGER PRIMARY KEY, amount REAL, kind TEXT)`)
	require.NoError(t, err)
	_, err = database.Exec(`INSERT INTO aggregate_test_items (id, amount, kind) VALUES (1, 42, 'a')`)
	require.NoError(t, err)
	querySet, err := NewQuerySet[aggregateTestItem]("aggregate_test_items")
	require.NoError(t, err)
	var asInterface QuerySet[aggregateTestItem] = querySet.SetDB(database)
	values, err := AggregateValues(context.Background(), asInterface, Sum("amount"))
	require.NoError(t, err)
	assert.Equal(t, 42.0, values["sum"])
}

func TestAggregateValuesPackageHelperUsesExportedCapability(t *testing.T) {
	querySet, err := NewQuerySet[aggregateTestItem]("aggregate_test_items")
	require.NoError(t, err)

	values, err := AggregateValues(context.Background(), aggregateValuesWrapper[aggregateTestItem]{querySet}, Count("id"))
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"sentinel": true}, values)

	_, err = AggregateValues(context.Background(), struct{ QuerySet[aggregateTestItem] }{querySet}, Count("id"))
	assert.True(t, forgeerrors.IsNotImplemented(err))
}

func TestAggregateValuesToManyFilterAggregatesBaseRows(t *testing.T) {
	database := setupToManyDedupeDB(t)
	defer database.Close()
	_, err := database.Exec(`UPDATE customers SET credit = CASE id WHEN 1 THEN 10 WHEN 2 THEN 20 WHEN 3 THEN 30 END`)
	require.NoError(t, err)

	querySet, err := NewQuerySet[TestCustomer]("customers")
	require.NoError(t, err)
	values, err := querySet.SetDB(database).Filter(F("orders__total").Gt(100.0)).(*BaseQuerySet[TestCustomer]).AggregateValues(context.Background(), Count("id"), Sum("credit"))
	require.NoError(t, err)
	assert.Equal(t, int64(2), values["count"])
	assert.Equal(t, 30.0, values["sum"])
}

func TestAggregateValuesDeferredFilterErrorOverridesAggregateChainError(t *testing.T) {
	database := setupRelationTestDB(t)
	require.NoError(t, database.Close())
	querySet, err := NewQuerySet[aggregateTestItem]("aggregate_test_items")
	require.NoError(t, err)
	base := querySet.SetDB(database).(*BaseQuerySet[aggregateTestItem])

	_, err = base.Aggregate(Count("id")).Filter(F("unknown_field").Eq(1)).(*BaseQuerySet[aggregateTestItem]).AggregateValues(context.Background(), Count("id"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid filter expression")
	assert.NotContains(t, err.Error(), "database is closed")

	_, err = base.Aggregate(Count("id")).All(context.Background())
	assert.True(t, forgeerrors.IsNotImplemented(err))
}

func TestAggregateValuesMinMaxPreserveBinaryValues(t *testing.T) {
	database := setupRelationTestDB(t)
	defer database.Close()
	_, err := database.Exec(`CREATE TABLE aggregate_blob_items (id INTEGER PRIMARY KEY, text TEXT, data BLOB)`)
	require.NoError(t, err)
	_, err = database.Exec(`INSERT INTO aggregate_blob_items (id, text, data) VALUES (1, 'alpha', x'01FF'), (2, 'zulu', x'02AA')`)
	require.NoError(t, err)

	querySet, err := NewQuerySet[aggregateBlobItem]("aggregate_blob_items")
	require.NoError(t, err)
	values, err := querySet.SetDB(database).(*BaseQuerySet[aggregateBlobItem]).AggregateValues(context.Background(), Aggregate{Name: "max_data", Field: "data", Func: string(AggMax)}, Aggregate{Name: "max_text", Field: "text", Func: string(AggMax)})
	require.NoError(t, err)
	assert.Equal(t, []byte{0x02, 0xAA}, values["max_data"])
	assert.Equal(t, "zulu", values["max_text"])
}

func TestAggregateValuesRelationJoinPreservedInOuterQuery(t *testing.T) {
	database := setupToManyDedupeDB(t)
	defer database.Close()

	ctx := context.Background()
	querySet, err := NewQuerySet[TestCustomer]("customers")
	require.NoError(t, err)

	t.Run("count to-many without filter", func(t *testing.T) {
		base := querySet.SetDB(database).(*BaseQuerySet[TestCustomer])
		values, err := base.AggregateValues(ctx, Count("orders__id"))
		require.NoError(t, err)
		assert.Equal(t, int64(5), values["count"])
	})

	t.Run("count to-many filtered by to-many", func(t *testing.T) {
		filtered := querySet.SetDB(database).Filter(F("orders__total").Gt(100.0)).(*BaseQuerySet[TestCustomer])
		values, err := filtered.AggregateValues(ctx, Count("orders__id"))
		require.NoError(t, err)
		assert.Equal(t, int64(4), values["count"])
	})

	t.Run("count base model filtered by to-many", func(t *testing.T) {
		filtered := querySet.SetDB(database).Filter(F("orders__total").Gt(100.0)).(*BaseQuerySet[TestCustomer])
		values, err := filtered.AggregateValues(ctx, Count("id"))
		require.NoError(t, err)
		assert.Equal(t, int64(2), values["count"])
	})

	t.Run("assert built SQL for case 2", func(t *testing.T) {
		filtered := querySet.SetDB(database).Filter(F("orders__total").Gt(100.0)).(*BaseQuerySet[TestCustomer])
		resolved, err := filtered.resolveAggregates([]Aggregate{Count("orders__id")})
		require.NoError(t, err)
		sql, _, err := filtered.buildAggregateSQL(resolved)
		require.NoError(t, err)

		subqueryMarker := `WHERE "customers"."id" IN (`
		require.Contains(t, sql, subqueryMarker)
		splitIdx := strings.Index(sql, subqueryMarker)
		outerSQL := sql[:splitIdx]
		subquerySQL := sql[splitIdx:]

		const expectedOrderJoin = `LEFT JOIN "orders" AS "orders" ON "orders"."customer_id" = "customers"."id"`
		assert.Contains(t, outerSQL, expectedOrderJoin, "outer query must contain orders join")
		assert.Contains(t, subquerySQL, expectedOrderJoin, "subquery must contain orders join")
		assert.Equal(t, 2, strings.Count(sql, expectedOrderJoin), "orders join must appear once in outer and once in subquery")
	})

	t.Run("count to-one filtered by to-many", func(t *testing.T) {
		filtered := querySet.SetDB(database).Filter(F("orders__total").Gt(100.0)).(*BaseQuerySet[TestCustomer])
		values, err := filtered.AggregateValues(ctx, Count("company__id"))
		require.NoError(t, err)
		assert.Equal(t, int64(2), values["count"])

		resolved, err := filtered.resolveAggregates([]Aggregate{Count("company__id")})
		require.NoError(t, err)
		sql, _, err := filtered.buildAggregateSQL(resolved)
		require.NoError(t, err)

		subqueryMarker := `WHERE "customers"."id" IN (`
		require.Contains(t, sql, subqueryMarker)
		splitIdx := strings.Index(sql, subqueryMarker)
		outerSQL := sql[:splitIdx]
		subquerySQL := sql[splitIdx:]

		const expectedCompanyJoin = `LEFT JOIN "companies" AS "company" ON "company"."id" = "customers"."company_id"`
		const expectedOrderJoin = `LEFT JOIN "orders" AS "orders" ON "orders"."customer_id" = "customers"."id"`
		assert.Contains(t, outerSQL, expectedCompanyJoin, "outer query must contain company join")
		assert.NotContains(t, outerSQL, expectedOrderJoin, "outer query must not contain orders join")
		assert.Contains(t, subquerySQL, expectedOrderJoin, "subquery must contain orders join")
		assert.NotContains(t, subquerySQL, expectedCompanyJoin, "subquery must not contain company join")
	})
}
