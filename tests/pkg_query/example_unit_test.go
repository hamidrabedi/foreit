package query

import (
	"testing"

	"github.com/forgego/forge/pkg/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example unit test for FieldExpression
func TestFieldExpression_ToSQL(t *testing.T) {
	tests := []struct {
		name     string
		field    query.FieldExpression[string]
		expected string
	}{
		{
			name:     "simple field",
			field:    query.NewField[string]("name", "users"),
			expected: `"users"."name"`,
		},
		{
			name:     "field without table",
			field:    query.NewField[string]("name", ""),
			expected: `"name"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := query.NewSQLBuilder()
			sql, _, err := tt.field.ToSQL(builder)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
		})
	}
}

// Example unit test for ComparisonExpression
func TestComparisonExpression_ToSQL(t *testing.T) {
	tests := []struct {
		name     string
		field    query.FieldExpression[float64]
		op       query.Operator
		value    interface{}
		expected string
	}{
		{
			name:     "equals",
			field:    query.NewField[float64]("price", "books"),
			op:       query.OpEquals,
			value:    10.0,
			expected: `"books"."price" = $1`,
		},
		{
			name:     "greater than",
			field:    query.NewField[float64]("price", "books"),
			op:       query.OpGreater,
			value:    10.0,
			expected: `"books"."price" > $1`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := query.ComparisonExpression[float64]{
				Field: tt.field,
				Op:    tt.op,
				Value: tt.value,
			}

			builder := query.NewSQLBuilder()
			sql, args, err := expr.ToSQL(builder)
			require.NoError(t, err)
			assert.Contains(t, sql, tt.expected)
			assert.Len(t, args, 1)
			assert.Equal(t, tt.value, args[0])
		})
	}
}

// Example unit test for Q object
func TestQ_And(t *testing.T) {
	priceField := query.NewField[float64]("price", "books")
	availableField := query.NewField[bool]("available", "books")

	q1 := query.NewQ(priceField.Gt(10.0))
	q2 := query.NewQ(availableField.Eq(true))
	combined := q1.And(q2)

	builder := query.NewSQLBuilder()
	sql, args, err := combined.ToSQL(builder)
	require.NoError(t, err)
	
	// Should contain both conditions
	assert.Contains(t, sql, "price")
	assert.Contains(t, sql, "available")
	assert.Contains(t, sql, "AND")
	assert.Len(t, args, 2)
}

// Example unit test for SQLBuilder
func TestSQLBuilder_AddArg(t *testing.T) {
	builder := query.NewSQLBuilder()
	
	placeholder1 := builder.AddArg("value1")
	placeholder2 := builder.AddArg("value2")
	
	assert.Equal(t, "$1", placeholder1)
	assert.Equal(t, "$2", placeholder2)
	assert.Equal(t, []interface{}{"value1", "value2"}, builder.Args())
}

// Example unit test for EscapeIdentifier
func TestEscapeIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple identifier",
			input:    "users",
			expected: `"users"`,
		},
		{
			name:     "identifier with quotes",
			input:    `user"name`,
			expected: `"user""name"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := query.EscapeIdentifier(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
