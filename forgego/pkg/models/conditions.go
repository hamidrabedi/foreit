package models

import (
	"fmt"

	"github.com/uptrace/bun"
)

// Condition represents a SQL WHERE condition.
type Condition interface {
	ToSQL() (string, []interface{})
}

// StringCondition represents a string comparison condition.
type StringCondition struct {
	field string
	op    string
	value interface{}
}

// ToSQL converts the condition to SQL.
func (c *StringCondition) ToSQL() (string, []interface{}) {
	return "? ? ?", []interface{}{bun.Ident(c.field), bun.Safe(c.op), c.value}
}

// NewStringCondition creates a new string condition.
func NewStringCondition(field, op string, value interface{}) Condition {
	return &StringCondition{field: field, op: op, value: value}
}

// IntCondition represents an integer comparison condition.
type IntCondition struct {
	field string
	op    string
	value interface{}
}

// ToSQL converts the condition to SQL.
func (c *IntCondition) ToSQL() (string, []interface{}) {
	return "? ? ?", []interface{}{bun.Ident(c.field), bun.Safe(c.op), c.value}
}

// NewIntCondition creates a new integer condition.
func NewIntCondition(field, op string, value interface{}) Condition {
	return &IntCondition{field: field, op: op, value: value}
}

// BoolCondition represents a boolean comparison condition.
type BoolCondition struct {
	field string
	op    string
	value interface{}
}

// ToSQL converts the condition to SQL.
func (c *BoolCondition) ToSQL() (string, []interface{}) {
	return "? ? ?", []interface{}{bun.Ident(c.field), bun.Safe(c.op), c.value}
}

// NewBoolCondition creates a new boolean condition.
func NewBoolCondition(field, op string, value interface{}) Condition {
	return &BoolCondition{field: field, op: op, value: value}
}

// InCondition represents an IN clause condition.
type InCondition struct {
	field  string
	values interface{}
}

// ToSQL converts the condition to SQL.
func (c *InCondition) ToSQL() (string, []interface{}) {
	return "? IN (?)", []interface{}{bun.Ident(c.field), bun.In(c.values)}
}

// NewInCondition creates a new IN condition.
func NewInCondition(field string, values interface{}) Condition {
	return &InCondition{field: field, values: values}
}

// IsNullCondition represents an IS NULL or IS NOT NULL condition.
type IsNullCondition struct {
	field string
	null  bool
}

// ToSQL converts the condition to SQL.
func (c *IsNullCondition) ToSQL() (string, []interface{}) {
	if c.null {
		return "? IS NULL", []interface{}{bun.Ident(c.field)}
	}
	return "? IS NOT NULL", []interface{}{bun.Ident(c.field)}
}

// NewIsNullCondition creates a new IS NULL condition.
func NewIsNullCondition(field string, null bool) Condition {
	return &IsNullCondition{field: field, null: null}
}

// OrderDirection represents the direction of ordering.
type OrderDirection int

const (
	OrderAsc OrderDirection = iota
	OrderDesc
)

// OrderBy represents an ordering clause.
type OrderBy struct {
	Field     string
	Direction OrderDirection
}

// ToSQL converts the OrderBy to SQL.
func (o OrderBy) ToSQL() string {
	if o.Direction == OrderDesc {
		return fmt.Sprintf("%s DESC", o.Field)
	}
	return fmt.Sprintf("%s ASC", o.Field)
}

