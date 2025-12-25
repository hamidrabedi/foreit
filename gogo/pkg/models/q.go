package models

import (
	"fmt"
)

// Q represents a query condition (like Django's Q objects)
type Q struct {
	conditions []FilterCondition
	connector  string // "AND" or "OR"
	negated    bool
}

// NewQ creates a new Q object
func NewQ() *Q {
	return &Q{
		conditions: make([]FilterCondition, 0),
		connector:  "AND",
		negated:    false,
	}
}

// Q creates a new Q object (convenience function)
func QFilter(field string, operator string, value interface{}) *Q {
	return NewQ().And(field, operator, value)
}

// And adds an AND condition
func (q *Q) And(field string, operator string, value interface{}) *Q {
	q.conditions = append(q.conditions, FilterCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	q.connector = "AND"
	return q
}

// Or adds an OR condition
func (q *Q) Or(field string, operator string, value interface{}) *Q {
	q.conditions = append(q.conditions, FilterCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	q.connector = "OR"
	return q
}

// Not negates the condition
func (q *Q) Not() *Q {
	q.negated = !q.negated
	return q
}

// Combine combines two Q objects
func (q *Q) Combine(other *Q, connector string) *Q {
	combined := NewQ()
	combined.conditions = append(combined.conditions, q.conditions...)
	combined.conditions = append(combined.conditions, other.conditions...)
	combined.connector = connector
	return combined
}

// AndQ combines with AND
func (q *Q) AndQ(other *Q) *Q {
	return q.Combine(other, "AND")
}

// OrQ combines with OR
func (q *Q) OrQ(other *Q) *Q {
	return q.Combine(other, "OR")
}

// toFilterConditions converts Q to filter conditions
func (q *Q) toFilterConditions() []FilterCondition {
	return q.conditions
}

// F represents a field reference (like Django's F expressions)
type F struct {
	field string
}

// NewF creates a new F expression
func NewF(field string) *F {
	return &F{field: field}
}

// F creates a new F expression (convenience function)
func FField(field string) *F {
	return NewF(field)
}

// String returns the field name
func (f *F) String() string {
	return f.field
}

// Field returns the field name
func (f *F) Field() string {
	return f.field
}

// Operations for F expressions
func (f *F) Eq(value interface{}) *Q {
	return QFilter(f.field, "=", value)
}

func (f *F) Ne(value interface{}) *Q {
	return QFilter(f.field, "!=", value)
}

func (f *F) Gt(value interface{}) *Q {
	return QFilter(f.field, ">", value)
}

func (f *F) Gte(value interface{}) *Q {
	return QFilter(f.field, ">=", value)
}

func (f *F) Lt(value interface{}) *Q {
	return QFilter(f.field, "<", value)
}

func (f *F) Lte(value interface{}) *Q {
	return QFilter(f.field, "<=", value)
}

func (f *F) In(values []interface{}) *Q {
	return QFilter(f.field, "IN", values)
}

func (f *F) Contains(value string) *Q {
	return QFilter(f.field, "LIKE", fmt.Sprintf("%%%s%%", value))
}

func (f *F) IContains(value string) *Q {
	return QFilter(f.field, "ILIKE", fmt.Sprintf("%%%s%%", value))
}

func (f *F) StartsWith(value string) *Q {
	return QFilter(f.field, "LIKE", fmt.Sprintf("%s%%", value))
}

func (f *F) EndsWith(value string) *Q {
	return QFilter(f.field, "LIKE", fmt.Sprintf("%%%s", value))
}

func (f *F) IsNull() *Q {
	return QFilter(f.field, "IS NULL", nil)
}

func (f *F) IsNotNull() *Q {
	return QFilter(f.field, "IS NOT NULL", nil)
}

