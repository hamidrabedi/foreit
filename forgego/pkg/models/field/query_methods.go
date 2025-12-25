package field

import (
	"fmt"
	"time"

	"github.com/forgego/forge/pkg/models"
)

// Eq returns an equality condition for string fields.
func (f Field[T, string]) Eq(value string) models.Condition {
	return models.NewStringCondition(f.Column, "=", value)
}

// Ne returns a not-equal condition for string fields.
func (f Field[T, string]) Ne(value string) models.Condition {
	return models.NewStringCondition(f.Column, "!=", value)
}

// Contains returns a LIKE condition for substring matching.
func (f Field[T, string]) Contains(value string) models.Condition {
	pattern := fmt.Sprintf("%%%v%%", value)
	return models.NewStringCondition(f.Column, "LIKE", pattern)
}

// IContains returns a case-insensitive LIKE condition for substring matching.
func (f Field[T, string]) IContains(value string) models.Condition {
	pattern := fmt.Sprintf("%%%v%%", value)
	return models.NewStringCondition(f.Column, "ILIKE", pattern)
}

// StartsWith returns a case-insensitive starts-with condition.
func (f Field[T, string]) StartsWith(value string) models.Condition {
	pattern := fmt.Sprintf("%v%%", value)
	return models.NewStringCondition(f.Column, "ILIKE", pattern)
}

// EndsWith returns a case-insensitive ends-with condition.
func (f Field[T, string]) EndsWith(value string) models.Condition {
	pattern := fmt.Sprintf("%%%v", value)
	return models.NewStringCondition(f.Column, "ILIKE", pattern)
}

// In returns an IN condition for matching any value in the slice.
func (f Field[T, string]) In(values []string) models.Condition {
	return models.NewInCondition(f.Column, values)
}

// EqInt returns an equality condition for int fields.
func (f Field[T, int]) EqInt(value int) models.Condition {
	return models.NewIntCondition(f.Column, "=", value)
}

// NeInt returns a not-equal condition for int fields.
func (f Field[T, int]) NeInt(value int) models.Condition {
	return models.NewIntCondition(f.Column, "!=", value)
}

// GtInt returns a greater-than condition for int fields.
func (f Field[T, int]) GtInt(value int) models.Condition {
	return models.NewIntCondition(f.Column, ">", value)
}

// GteInt returns a greater-than-or-equal condition for int fields.
func (f Field[T, int]) GteInt(value int) models.Condition {
	return models.NewIntCondition(f.Column, ">=", value)
}

// LtInt returns a less-than condition for int fields.
func (f Field[T, int]) LtInt(value int) models.Condition {
	return models.NewIntCondition(f.Column, "<", value)
}

// LteInt returns a less-than-or-equal condition for int fields.
func (f Field[T, int]) LteInt(value int) models.Condition {
	return models.NewIntCondition(f.Column, "<=", value)
}

// InInt returns an IN condition for int fields.
func (f Field[T, int]) InInt(values []int) models.Condition {
	return models.NewInCondition(f.Column, values)
}

// EqInt64 returns an equality condition for int64 fields.
func (f Field[T, int64]) EqInt64(value int64) models.Condition {
	return models.NewIntCondition(f.Column, "=", value)
}

// NeInt64 returns a not-equal condition for int64 fields.
func (f Field[T, int64]) NeInt64(value int64) models.Condition {
	return models.NewIntCondition(f.Column, "!=", value)
}

// GtInt64 returns a greater-than condition for int64 fields.
func (f Field[T, int64]) GtInt64(value int64) models.Condition {
	return models.NewIntCondition(f.Column, ">", value)
}

// GteInt64 returns a greater-than-or-equal condition for int64 fields.
func (f Field[T, int64]) GteInt64(value int64) models.Condition {
	return models.NewIntCondition(f.Column, ">=", value)
}

// LtInt64 returns a less-than condition for int64 fields.
func (f Field[T, int64]) LtInt64(value int64) models.Condition {
	return models.NewIntCondition(f.Column, "<", value)
}

// LteInt64 returns a less-than-or-equal condition for int64 fields.
func (f Field[T, int64]) LteInt64(value int64) models.Condition {
	return models.NewIntCondition(f.Column, "<=", value)
}

// InInt64 returns an IN condition for int64 fields.
func (f Field[T, int64]) InInt64(values []int64) models.Condition {
	return models.NewInCondition(f.Column, values)
}

// IsTrue returns a condition matching true values for bool fields.
func (f Field[T, bool]) IsTrue() models.Condition {
	return models.NewIntCondition(f.Column, "=", 1)
}

// IsFalse returns a condition matching false values for bool fields.
func (f Field[T, bool]) IsFalse() models.Condition {
	return models.NewIntCondition(f.Column, "=", 0)
}

// TimeFieldQuery provides query methods for time.Time fields.
type TimeFieldQuery[T any] struct {
	column string
}

// ForTimeField creates a query helper for a time.Time field.
func ForTimeField[T any](f Field[T, time.Time]) TimeFieldQuery[T] {
	return TimeFieldQuery[T]{column: f.Column}
}

// Eq returns an equality condition for time fields.
func (q TimeFieldQuery[T]) Eq(value time.Time) models.Condition {
	return models.NewStringCondition(q.column, "=", value)
}

// Gt returns a greater-than condition for time fields.
func (q TimeFieldQuery[T]) Gt(value time.Time) models.Condition {
	return models.NewStringCondition(q.column, ">", value)
}

// Gte returns a greater-than-or-equal condition for time fields.
func (q TimeFieldQuery[T]) Gte(value time.Time) models.Condition {
	return models.NewStringCondition(q.column, ">=", value)
}

// Lt returns a less-than condition for time fields.
func (q TimeFieldQuery[T]) Lt(value time.Time) models.Condition {
	return models.NewStringCondition(q.column, "<", value)
}

// Lte returns a less-than-or-equal condition for time fields.
func (q TimeFieldQuery[T]) Lte(value time.Time) models.Condition {
	return models.NewStringCondition(q.column, "<=", value)
}

// EqFloat64 returns an equality condition for float64 fields.
func (f Field[T, float64]) EqFloat64(value float64) models.Condition {
	return models.NewIntCondition(f.Column, "=", value)
}

// NeFloat64 returns a not-equal condition for float64 fields.
func (f Field[T, float64]) NeFloat64(value float64) models.Condition {
	return models.NewIntCondition(f.Column, "!=", value)
}

// GtFloat64 returns a greater-than condition for float64 fields.
func (f Field[T, float64]) GtFloat64(value float64) models.Condition {
	return models.NewIntCondition(f.Column, ">", value)
}

// GteFloat64 returns a greater-than-or-equal condition for float64 fields.
func (f Field[T, float64]) GteFloat64(value float64) models.Condition {
	return models.NewIntCondition(f.Column, ">=", value)
}

// LtFloat64 returns a less-than condition for float64 fields.
func (f Field[T, float64]) LtFloat64(value float64) models.Condition {
	return models.NewIntCondition(f.Column, "<", value)
}

// LteFloat64 returns a less-than-or-equal condition for float64 fields.
func (f Field[T, float64]) LteFloat64(value float64) models.Condition {
	return models.NewIntCondition(f.Column, "<=", value)
}
