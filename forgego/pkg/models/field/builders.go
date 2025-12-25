package field

import "time"

// StringField creates a string field with getter/setter functions.
func StringField[T any](
	getter func(*T) *string,
	setter func(*T, string),
	opts ...Option,
) Field[T, string] {
	return NewField[T, string]("", "", getter, setter, opts...)
}

// IntField creates an integer field with getter/setter functions.
func IntField[T any](
	getter func(*T) *int,
	setter func(*T, int),
	opts ...Option,
) Field[T, int] {
	return NewField[T, int]("", "", getter, setter, opts...)
}

// Int64Field creates an int64 field with getter/setter functions.
func Int64Field[T any](
	getter func(*T) *int64,
	setter func(*T, int64),
	opts ...Option,
) Field[T, int64] {
	return NewField[T, int64]("", "", getter, setter, opts...)
}

// BoolField creates a boolean field with getter/setter functions.
func BoolField[T any](
	getter func(*T) *bool,
	setter func(*T, bool),
	opts ...Option,
) Field[T, bool] {
	return NewField[T, bool]("", "", getter, setter, opts...)
}

// TimeField creates a time.Time field with getter/setter functions.
func TimeField[T any](
	getter func(*T) *time.Time,
	setter func(*T, time.Time),
	opts ...Option,
) Field[T, time.Time] {
	return NewField[T, time.Time]("", "", getter, setter, opts...)
}

// TimePtrField creates a *time.Time field with getter/setter functions.
func TimePtrField[T any](
	getter func(*T) **time.Time,
	setter func(*T, *time.Time),
	opts ...Option,
) Field[T, *time.Time] {
	return NewField[T, *time.Time]("", "", func(t *T) *(*time.Time) {
		ptr := getter(t)
		return ptr
	}, setter, opts...)
}

// Float64Field creates a float64 field with getter/setter functions.
func Float64Field[T any](
	getter func(*T) *float64,
	setter func(*T, float64),
	opts ...Option,
) Field[T, float64] {
	return NewField[T, float64]("", "", getter, setter, opts...)
}

// JSONField creates a JSON field with generic type V and getter/setter functions.
func JSONField[T any, V any](
	getter func(*T) *V,
	setter func(*T, V),
	opts ...Option,
) Field[T, V] {
	return NewField[T, V]("", "", getter, setter, opts...)
}

// UUIDField creates a UUID field (string type) with getter/setter functions.
func UUIDField[T any](
	getter func(*T) *string,
	setter func(*T, string),
	opts ...Option,
) Field[T, string] {
	return NewField[T, string]("", "", getter, setter, opts...)
}

// ForeignKeyField creates a foreign key field (int type) with getter/setter functions.
func ForeignKeyField[T any](
	getter func(*T) *int,
	setter func(*T, int),
	relatedModel string,
	opts ...Option,
) Field[T, int] {
	opts = append(opts, RelatedModel(relatedModel))
	return NewField[T, int]("", "", getter, setter, opts...)
}

