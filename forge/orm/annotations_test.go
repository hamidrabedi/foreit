package orm

import (
	"fmt"
	"testing"

	"github.com/forgego/forge/db/dialect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnnotation_QueryExprAndExpressionRendering(t *testing.T) {
	fieldExpr := NewFieldExpr[float64]("CreatedAt", "")
	annFromField := NewAnnotation("high_price", fieldExpr.Greater(10.0))

	field := NewField[float64]("CreatedAt", "")
	annFromExpr := NewExpressionAnnotation("high_price", field.Gt(10.0))

	t.Run("PostgreSQL builder produces identical SELECT fragment", func(t *testing.T) {
		qs1, err := NewQuerySet[testModel]("test_table")
		require.NoError(t, err)
		qs1 = qs1.Annotate(annFromField)
		base1 := qs1.(*BaseQuerySet[testModel])
		b1 := NewSQLBuilderWithDialect(dialect.NewPostgreSQLDialect())
		select1 := base1.buildSelectClause(b1, false)

		qs2, err := NewQuerySet[testModel]("test_table")
		require.NoError(t, err)
		qs2 = qs2.Annotate(annFromExpr)
		base2 := qs2.(*BaseQuerySet[testModel])
		b2 := NewSQLBuilderWithDialect(dialect.NewPostgreSQLDialect())
		select2 := base2.buildSelectClause(b2, false)

		assert.Contains(t, select1, `CreatedAt > $1 AS "high_price"`)
		assert.Contains(t, select2, `"CreatedAt" > $1 AS "high_price"`)
		assert.Equal(t, []interface{}{10.0}, b1.Args())
	})

	t.Run("SQLite builder produces identical SELECT fragment", func(t *testing.T) {
		qs1, err := NewQuerySet[testModel]("test_table")
		require.NoError(t, err)
		qs1 = qs1.Annotate(annFromField)
		base1 := qs1.(*BaseQuerySet[testModel])
		b1 := NewSQLBuilderWithDialect(dialect.NewSQLiteDialect())
		select1 := base1.buildSelectClause(b1, false)

		qs2, err := NewQuerySet[testModel]("test_table")
		require.NoError(t, err)
		qs2 = qs2.Annotate(annFromExpr)
		base2 := qs2.(*BaseQuerySet[testModel])
		b2 := NewSQLBuilderWithDialect(dialect.NewSQLiteDialect())
		select2 := base2.buildSelectClause(b2, false)

		assert.Contains(t, select1, `CreatedAt > ? AS "high_price"`)
		assert.Contains(t, select2, `"CreatedAt" > ? AS "high_price"`)
		assert.Equal(t, []interface{}{10.0}, b1.Args())
	})
}

func TestAnnotation_EqualityComparisonEquivalence(t *testing.T) {
	fieldExpr := NewFieldExpr[bool]("available", "")
	annOld := NewAnnotation("is_available", fieldExpr.Equals(true))

	field := NewField[bool]("available", "")
	annNew := NewExpressionAnnotation("is_available", field.Eq(true))

	t.Run("PostgreSQL equality fragment", func(t *testing.T) {
		bOld := NewSQLBuilderWithDialect(dialect.NewPostgreSQLDialect())
		sqlOld, _, err := annOld.Expression.ToSQL(bOld)
		require.NoError(t, err)
		fragOld := fmt.Sprintf("%s AS %s", sqlOld, EscapeIdentifier(annOld.Name))

		bNew := NewSQLBuilderWithDialect(dialect.NewPostgreSQLDialect())
		sqlNew, _, err := annNew.Expression.ToSQL(bNew)
		require.NoError(t, err)
		fragNew := fmt.Sprintf("%s AS %s", sqlNew, EscapeIdentifier(annNew.Name))

		assert.Equal(t, `available = $1 AS "is_available"`, fragOld)
		assert.Equal(t, `"available" = $1 AS "is_available"`, fragNew)
	})

	t.Run("SQLite equality fragment", func(t *testing.T) {
		bOld := NewSQLBuilderWithDialect(dialect.NewSQLiteDialect())
		sqlOld, _, err := annOld.Expression.ToSQL(bOld)
		require.NoError(t, err)
		fragOld := fmt.Sprintf("%s AS %s", sqlOld, EscapeIdentifier(annOld.Name))

		bNew := NewSQLBuilderWithDialect(dialect.NewSQLiteDialect())
		sqlNew, _, err := annNew.Expression.ToSQL(bNew)
		require.NoError(t, err)
		fragNew := fmt.Sprintf("%s AS %s", sqlNew, EscapeIdentifier(annNew.Name))

		assert.Equal(t, `available = ? AS "is_available"`, fragOld)
		assert.Equal(t, `"available" = ? AS "is_available"`, fragNew)
	})
}

func TestAnnotation_UnsetExpressionFieldFallback(t *testing.T) {
	ann := AnnotationExpr{
		Name: "custom_flag",
		Expr: NewFieldQueryExpr("price", OpGreater, 50.0),
	}

	qs, err := NewQuerySet[testModel]("test_table")
	require.NoError(t, err)
	qs = qs.Annotate(ann)
	base := qs.(*BaseQuerySet[testModel])

	b := NewSQLBuilderWithDialect(dialect.NewPostgreSQLDialect())
	selectClause := base.buildSelectClause(b, false)
	assert.Contains(t, selectClause, `price > $1 AS "custom_flag"`)
	assert.Equal(t, []interface{}{50.0}, b.Args())
}

func TestSQLBuilder_BuildWhere_AdapterEquivalence(t *testing.T) {
	conditions := []QueryExpr{
		NewFieldQueryExpr("price", OpGreater, 10.0),
		NewFieldQueryExpr("available", OpEquals, true),
	}
	excludes := []QueryExpr{
		NewFieldQueryExpr("deleted", OpEquals, true),
	}

	t.Run("PostgreSQL BuildWhere", func(t *testing.T) {
		b := NewSQLBuilderWithDialect(dialect.NewPostgreSQLDialect())
		where, args := b.BuildWhere(conditions, excludes)
		assert.Equal(t, `WHERE price > $1 AND available = $2 AND NOT (deleted = $3)`, where)
		assert.Equal(t, []interface{}{10.0, true, true}, args)
	})

	t.Run("SQLite BuildWhere", func(t *testing.T) {
		b := NewSQLBuilderWithDialect(dialect.NewSQLiteDialect())
		where, args := b.BuildWhere(conditions, excludes)
		assert.Equal(t, `WHERE price > ? AND available = ? AND NOT (deleted = ?)`, where)
		assert.Equal(t, []interface{}{10.0, true, true}, args)
	})
}
