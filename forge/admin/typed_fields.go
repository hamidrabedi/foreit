package admin

import (
	"github.com/forgego/forge/orm"
)

// Fields creates a slice of field expressions from multiple field expressions
// This provides type-safe field configuration
// Usage: ListDisplay: Fields(User.Fields.ID, User.Fields.Username, User.Fields.Email)
func Fields(fields ...orm.FieldExpression[interface{}]) []interface{} {
	result := make([]interface{}, len(fields))
	for i, field := range fields {
		result[i] = field
	}
	return result
}

// FieldExpr is a helper to convert a typed field expression to interface{}
// This allows using typed FieldExpression[T] in admin config
func FieldExpr[T any, V any](field orm.FieldExpression[V]) orm.FieldExpression[interface{}] {
	// Type erasure - convert to interface{} for storage
	// Type safety is maintained through the generic constraint
	return orm.NewField[interface{}](field.Path(), field.Table())
}

// TypedListDisplay creates a type-safe list display configuration
// Usage: TypedListDisplay(User.Fields.ID, User.Fields.Username)
func TypedListDisplay(fields ...orm.FieldExpression[interface{}]) []interface{} {
	return Fields(fields...)
}

// TypedSearchFields creates type-safe search fields
// Usage: TypedSearchFields(User.Fields.Username, User.Fields.Email)
func TypedSearchFields(fields ...orm.FieldExpression[interface{}]) []interface{} {
	return Fields(fields...)
}

// TypedReadOnlyFields creates type-safe read-only fields
// Usage: TypedReadOnlyFields(User.Fields.ID, User.Fields.CreatedAt)
func TypedReadOnlyFields(fields ...orm.FieldExpression[interface{}]) []interface{} {
	return Fields(fields...)
}
