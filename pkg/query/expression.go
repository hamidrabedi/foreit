package query

import (
	"fmt"
	"reflect"
	"strings"
)

// Expression is the base interface for all query expressions
type Expression interface {
	// ToSQL converts the expression to SQL with parameters
	ToSQL(builder *SQLBuilder) (string, []interface{}, error)
	// Resolve validates and resolves the expression against a schema
	Resolve(schema *ModelSchema) error
}

// FieldExpression represents a field reference (like Django's F)
// T is the type of the field value
type FieldExpression[T any] struct {
	fieldPath string // e.g., "price", "author__name"
	table     string
	fieldType reflect.Type
}

// NewField creates a type-safe field expression
func NewField[T any](path, table string) FieldExpression[T] {
	return FieldExpression[T]{
		fieldPath: path,
		table:     table,
		fieldType: reflect.TypeOf((*T)(nil)).Elem(),
	}
}

// Path returns the field path
func (f FieldExpression[T]) Path() string {
	return f.fieldPath
}

// Table returns the table name
func (f FieldExpression[T]) Table() string {
	return f.table
}

// ToSQL converts field expression to SQL
func (f FieldExpression[T]) ToSQL(builder *SQLBuilder) (string, []interface{}, error) {
	// Escape identifier - handle table prefix if needed
	if f.table != "" && !strings.Contains(f.fieldPath, ".") {
		escaped := EscapeIdentifier(f.table) + "." + EscapeIdentifier(f.fieldPath)
		return escaped, nil, nil
	}
	escaped := EscapeIdentifier(f.fieldPath)
	return escaped, nil, nil
}

// Resolve validates the field exists in schema
func (f FieldExpression[T]) Resolve(schema *ModelSchema) error {
	// Validate field path exists
	parts := splitFieldPath(f.fieldPath)
	if len(parts) == 0 {
		return fmt.Errorf("empty field path")
	}
	
	// Check first part exists
	field := schema.GetField(parts[0])
	if field == nil {
		return fmt.Errorf("field %s not found in model", parts[0])
	}
	
	// Validate type matches
	if field.Type != f.fieldType {
		return fmt.Errorf("field %s has type %v, expected %v", parts[0], field.Type, f.fieldType)
	}
	
	// TODO: Validate relation paths for nested fields
	
	return nil
}

// Arithmetic operations

// Add creates an addition expression
func (f FieldExpression[T]) Add(other FieldExpression[T]) CombinedExpression[T] {
	return CombinedExpression[T]{
		Left:  f,
		Op:    OpAdd,
		Right: other,
	}
}

// Sub creates a subtraction expression
func (f FieldExpression[T]) Sub(other FieldExpression[T]) CombinedExpression[T] {
	return CombinedExpression[T]{
		Left:  f,
		Op:    OpSub,
		Right: other,
	}
}

// Mul creates a multiplication expression
func (f FieldExpression[T]) Mul(other FieldExpression[T]) CombinedExpression[T] {
	return CombinedExpression[T]{
		Left:  f,
		Op:    OpMul,
		Right: other,
	}
}

// Div creates a division expression
func (f FieldExpression[T]) Div(other FieldExpression[T]) CombinedExpression[T] {
	return CombinedExpression[T]{
		Left:  f,
		Op:    OpDiv,
		Right: other,
	}
}

// Comparison operations

// Eq creates an equality comparison
func (f FieldExpression[T]) Eq(val T) ComparisonExpression[T] {
	return ComparisonExpression[T]{
		Field: f,
		Op:    OpEquals,
		Value: val,
	}
}

// Ne creates a not-equal comparison
func (f FieldExpression[T]) Ne(val T) ComparisonExpression[T] {
	return ComparisonExpression[T]{
		Field: f,
		Op:    OpNotEquals,
		Value: val,
	}
}

// Gt creates a greater-than comparison
func (f FieldExpression[T]) Gt(val T) ComparisonExpression[T] {
	return ComparisonExpression[T]{
		Field: f,
		Op:    OpGreater,
		Value: val,
	}
}

// Gte creates a greater-than-or-equal comparison
func (f FieldExpression[T]) Gte(val T) ComparisonExpression[T] {
	return ComparisonExpression[T]{
		Field: f,
		Op:    OpGreaterOrEqual,
		Value: val,
	}
}

// Lt creates a less-than comparison
func (f FieldExpression[T]) Lt(val T) ComparisonExpression[T] {
	return ComparisonExpression[T]{
		Field: f,
		Op:    OpLess,
		Value: val,
	}
}

// Lte creates a less-than-or-equal comparison
func (f FieldExpression[T]) Lte(val T) ComparisonExpression[T] {
	return ComparisonExpression[T]{
		Field: f,
		Op:    OpLessOrEqual,
		Value: val,
	}
}

// In creates an IN clause comparison
func (f FieldExpression[T]) In(vals ...T) ComparisonExpression[T] {
	return ComparisonExpression[T]{
		Field: f,
		Op:    OpIn,
		Value: vals,
	}
}

// IsNull creates an IS NULL comparison
func (f FieldExpression[T]) IsNull() ComparisonExpression[T] {
	return ComparisonExpression[T]{
		Field: f,
		Op:    OpIsNull,
		Value: nil,
	}
}

// IsNotNull creates an IS NOT NULL comparison
func (f FieldExpression[T]) IsNotNull() ComparisonExpression[T] {
	return ComparisonExpression[T]{
		Field: f,
		Op:    OpIsNotNull,
		Value: nil,
	}
}

// Range creates a BETWEEN comparison
func (f FieldExpression[T]) Range(min, max T) ComparisonExpression[T] {
	return ComparisonExpression[T]{
		Field: f,
		Op:    OpRange,
		Value: []T{min, max},
	}
}

// String-specific operations (only for FieldExpression[string])

// Contains creates a LIKE '%value%' comparison (string only)
func (f FieldExpression[string]) Contains(val string) ComparisonExpression[string] {
	return ComparisonExpression[string]{
		Field: f,
		Op:    OpContains,
		Value: val,
	}
}

// StartsWith creates a LIKE 'value%' comparison (string only)
func (f FieldExpression[string]) StartsWith(val string) ComparisonExpression[string] {
	return ComparisonExpression[string]{
		Field: f,
		Op:    OpStartsWith,
		Value: val,
	}
}

// EndsWith creates a LIKE '%value' comparison (string only)
func (f FieldExpression[string]) EndsWith(val string) ComparisonExpression[string] {
	return ComparisonExpression[string]{
		Field: f,
		Op:    OpEndsWith,
		Value: val,
	}
}

// IContains creates an ILIKE '%value%' comparison (case-insensitive, string only)
func (f FieldExpression[string]) IContains(val string) ComparisonExpression[string] {
	return ComparisonExpression[string]{
		Field: f,
		Op:    OpIContains,
		Value: val,
	}
}

// IExact creates an ILIKE 'value' comparison (case-insensitive exact, string only)
func (f FieldExpression[string]) IExact(val string) ComparisonExpression[string] {
	return ComparisonExpression[string]{
		Field: f,
		Op:    OpIExact,
		Value: val,
	}
}

// CombinedExpression represents arithmetic operations between fields
type CombinedExpression[T any] struct {
	Left  FieldExpression[T]
	Op    ArithmeticOperator
	Right FieldExpression[T]
}

// ArithmeticOperator represents arithmetic operations
type ArithmeticOperator string

const (
	OpAdd ArithmeticOperator = "+"
	OpSub ArithmeticOperator = "-"
	OpMul ArithmeticOperator = "*"
	OpDiv ArithmeticOperator = "/"
	OpMod ArithmeticOperator = "%"
)

// ToSQL converts combined expression to SQL
func (c CombinedExpression[T]) ToSQL(builder *SQLBuilder) (string, []interface{}, error) {
	leftSQL, leftArgs, err := c.Left.ToSQL(builder)
	if err != nil {
		return "", nil, err
	}
	
	rightSQL, rightArgs, err := c.Right.ToSQL(builder)
	if err != nil {
		return "", nil, err
	}
	
	sql := fmt.Sprintf("(%s %s %s)", leftSQL, c.Op, rightSQL)
	args := append(leftArgs, rightArgs...)
	
	return sql, args, nil
}

// Resolve validates the combined expression
func (c CombinedExpression[T]) Resolve(schema *ModelSchema) error {
	if err := c.Left.Resolve(schema); err != nil {
		return err
	}
	return c.Right.Resolve(schema)
}

// ComparisonExpression represents comparison operations
type ComparisonExpression[T any] struct {
	Field FieldExpression[T]
	Op    Operator
	Value interface{}
}

// ToSQL converts comparison expression to SQL
func (c ComparisonExpression[T]) ToSQL(builder *SQLBuilder) (string, []interface{}, error) {
	fieldSQL, _, err := c.Field.ToSQL(builder)
	if err != nil {
		return "", nil, err
	}
	
	// Build SQL based on operator
	var sql string
	var args []interface{}
	
	// Handle operators that share the same string value using if-else
	if c.Op == OpContains {
		// LIKE '%value%'
		if strVal, ok := c.Value.(string); ok {
			pattern := "%" + strVal + "%"
			placeholder := builder.AddArg(pattern)
			sql = fmt.Sprintf("%s LIKE %s", fieldSQL, placeholder)
			args = []interface{}{pattern}
		} else {
			return "", nil, fmt.Errorf("Contains operator requires string value")
		}
	} else if c.Op == OpStartsWith {
		// LIKE 'value%'
		if strVal, ok := c.Value.(string); ok {
			pattern := strVal + "%"
			placeholder := builder.AddArg(pattern)
			sql = fmt.Sprintf("%s LIKE %s", fieldSQL, placeholder)
			args = []interface{}{pattern}
		} else {
			return "", nil, fmt.Errorf("StartsWith operator requires string value")
		}
	} else if c.Op == OpEndsWith {
		// LIKE '%value'
		if strVal, ok := c.Value.(string); ok {
			pattern := "%" + strVal
			placeholder := builder.AddArg(pattern)
			sql = fmt.Sprintf("%s LIKE %s", fieldSQL, placeholder)
			args = []interface{}{pattern}
		} else {
			return "", nil, fmt.Errorf("EndsWith operator requires string value")
		}
	} else if c.Op == OpIContains {
		// ILIKE '%value%'
		if strVal, ok := c.Value.(string); ok {
			pattern := "%" + strVal + "%"
			placeholder := builder.AddArg(pattern)
			sql = fmt.Sprintf("%s ILIKE %s", fieldSQL, placeholder)
			args = []interface{}{pattern}
		} else {
			return "", nil, fmt.Errorf("IContains operator requires string value")
		}
	} else if c.Op == OpIExact {
		// ILIKE 'value'
		if strVal, ok := c.Value.(string); ok {
			placeholder := builder.AddArg(strVal)
			sql = fmt.Sprintf("%s ILIKE %s", fieldSQL, placeholder)
			args = []interface{}{strVal}
		} else {
			return "", nil, fmt.Errorf("IExact operator requires string value")
		}
	} else {
		// Use switch for operators with unique values
		switch c.Op {
		case OpIsNull:
			sql = fmt.Sprintf("%s IS NULL", fieldSQL)
		case OpIsNotNull:
			sql = fmt.Sprintf("%s IS NOT NULL", fieldSQL)
		case OpIn:
			values, ok := c.Value.([]T)
			if !ok {
				return "", nil, fmt.Errorf("IN operator requires slice value")
			}
			placeholders := make([]string, len(values))
			for i, val := range values {
				placeholders[i] = builder.AddArg(val)
				args = append(args, val)
			}
			sql = fmt.Sprintf("%s IN (%s)", fieldSQL, strings.Join(placeholders, ", "))
		case OpNotIn:
			values, ok := c.Value.([]T)
			if !ok {
				return "", nil, fmt.Errorf("NOT IN operator requires slice value")
			}
			placeholders := make([]string, len(values))
			for i, val := range values {
				placeholders[i] = builder.AddArg(val)
				args = append(args, val)
			}
			sql = fmt.Sprintf("%s NOT IN (%s)", fieldSQL, strings.Join(placeholders, ", "))
		case OpRange:
			values, ok := c.Value.([]T)
			if !ok || len(values) != 2 {
				return "", nil, fmt.Errorf("Range operator requires slice of 2 values")
			}
			placeholder1 := builder.AddArg(values[0])
			placeholder2 := builder.AddArg(values[1])
			sql = fmt.Sprintf("%s BETWEEN %s AND %s", fieldSQL, placeholder1, placeholder2)
			args = []interface{}{values[0], values[1]}
		default:
			// Standard operators (=, !=, >, >=, <, <=)
			placeholder := builder.AddArg(c.Value)
			sql = fmt.Sprintf("%s %s %s", fieldSQL, c.Op, placeholder)
			args = []interface{}{c.Value}
		}
	}
	
	return sql, args, nil
}

// Resolve validates the comparison expression
func (c ComparisonExpression[T]) Resolve(schema *ModelSchema) error {
	return c.Field.Resolve(schema)
}

// ValueExpression represents a literal value
type ValueExpression[T any] struct {
	value T
}

// NewValue creates a value expression
func NewValue[T any](val T) ValueExpression[T] {
	return ValueExpression[T]{value: val}
}

// ToSQL converts value expression to SQL
func (v ValueExpression[T]) ToSQL(builder *SQLBuilder) (string, []interface{}, error) {
	placeholder := builder.AddArg(v.value)
	return placeholder, []interface{}{v.value}, nil
}

// Resolve validates the value expression (always valid)
func (v ValueExpression[T]) Resolve(schema *ModelSchema) error {
	return nil
}

// Q is for complex query composition (like Django's Q)
type Q struct {
	expressions []Expression
	connector   Connector // AND, OR
	negated     bool
}

// NewQ creates a new Q object from an expression
func NewQ(expr Expression) *Q {
	return &Q{
		expressions: []Expression{expr},
		connector:   ConnectorAnd,
	}
}

// And combines with AND
func (q *Q) And(other *Q) *Q {
	return &Q{
		expressions: []Expression{q, other},
		connector:   ConnectorAnd,
	}
}

// Or combines with OR
func (q *Q) Or(other *Q) *Q {
	return &Q{
		expressions: []Expression{q, other},
		connector:   ConnectorOr,
	}
}

// Not negates the expression
func (q *Q) Not() *Q {
	return &Q{
		expressions: q.expressions,
		connector:   q.connector,
		negated:     !q.negated,
	}
}

// ToSQL converts Q to SQL
func (q *Q) ToSQL(builder *SQLBuilder) (string, []interface{}, error) {
	if len(q.expressions) == 0 {
		return "1=1", nil, nil
	}
	
	var parts []string
	var allArgs []interface{}
	
	for _, expr := range q.expressions {
		sql, args, err := expr.ToSQL(builder)
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, fmt.Sprintf("(%s)", sql))
		allArgs = append(allArgs, args...)
	}
	
	combinedSQL := ""
	if len(parts) > 0 {
		combinedSQL = parts[0]
		for i := 1; i < len(parts); i++ {
			combinedSQL = fmt.Sprintf("%s %s %s", combinedSQL, q.connector, parts[i])
		}
	}
	
	if q.negated {
		combinedSQL = fmt.Sprintf("NOT (%s)", combinedSQL)
	}
	
	return combinedSQL, allArgs, nil
}

// Resolve validates all expressions in Q
func (q *Q) Resolve(schema *ModelSchema) error {
	for _, expr := range q.expressions {
		if err := expr.Resolve(schema); err != nil {
			return err
		}
	}
	return nil
}

// Connector represents how expressions are combined
type Connector string

const (
	ConnectorAnd Connector = "AND"
	ConnectorOr  Connector = "OR"
)

// Helper function to split field path (e.g., "author__name" -> ["author", "name"])
func splitFieldPath(path string) []string {
	// Simple implementation - can be enhanced for more complex paths
	parts := []string{path}
	// TODO: Handle double underscore for actual underscores
	return parts
}
