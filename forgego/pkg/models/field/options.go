package field

// Option represents a field configuration option.
type Option interface{}

// PrimaryKey marks the field as a primary key
func PrimaryKey() Option {
	return primaryKeyOption{}
}

type primaryKeyOption struct{}

// Required marks the field as required
func Required() Option {
	return requiredOption{}
}

type requiredOption struct{}

// Optional marks the field as optional
func Optional() Option {
	return optionalOption{}
}

type optionalOption struct{}

// Unique marks the field as unique
func Unique() Option {
	return uniqueOption{}
}

type uniqueOption struct{}

// Index marks the field to be indexed
func Index() Option {
	return indexOption{}
}

type indexOption struct{}

// Default sets a default value
func Default(value interface{}) Option {
	return defaultOption{Value: value}
}

type defaultOption struct {
	Value interface{}
}

// MaxLength sets maximum length (for strings)
func MaxLength(n int) Option {
	return maxLengthOption{Max: n}
}

type maxLengthOption struct {
	Max int
}

// MinLength sets minimum length (for strings)
func MinLength(n int) Option {
	return minLengthOption{Min: n}
}

type minLengthOption struct {
	Min int
}

// MaxValue sets maximum value (for numbers)
func MaxValue(n interface{}) Option {
	return maxValueOption{Max: n}
}

type maxValueOption struct {
	Max interface{}
}

// MinValue sets minimum value (for numbers)
func MinValue(n interface{}) Option {
	return minValueOption{Min: n}
}

type minValueOption struct {
	Min interface{}
}

// Email adds email validator
func Email() Option {
	return emailOption{}
}

type emailOption struct{}

// AutoNow sets field to current time on update
func AutoNow() Option {
	return autoNowOption{}
}

type autoNowOption struct{}

// AutoNowAdd sets field to current time on create
func AutoNowAdd() Option {
	return autoNowAddOption{}
}

type autoNowAddOption struct{}

// Column sets custom database column name
func Column(name string) Option {
	return columnOption{Name: name}
}

type columnOption struct {
	Name string
}

// RelatedModel sets the related model for foreign keys
func RelatedModel(model string) Option {
	return relatedModelOption{Model: model}
}

type relatedModelOption struct {
	Model string
}

// HasOption checks if a field has a specific option using a custom check function.
func HasOption[T any, V any](f Field[T, V], check func(interface{}) bool) bool {
	for _, opt := range f.Options {
		if check(opt) {
			return true
		}
	}
	return false
}

// IsPrimaryKey checks if field is primary key
func IsPrimaryKey[T any, V any](f Field[T, V]) bool {
	return HasOption(f, func(opt interface{}) bool {
		_, ok := opt.(primaryKeyOption)
		return ok
	})
}

// IsRequired checks if field is required
func IsRequired[T any, V any](f Field[T, V]) bool {
	return HasOption(f, func(opt interface{}) bool {
		_, ok := opt.(requiredOption)
		return ok
	})
}

// IsUnique checks if field is unique
func IsUnique[T any, V any](f Field[T, V]) bool {
	return HasOption(f, func(opt interface{}) bool {
		_, ok := opt.(uniqueOption)
		return ok
	})
}

// IsIndexed checks if field is indexed
func IsIndexed[T any, V any](f Field[T, V]) bool {
	return HasOption(f, func(opt interface{}) bool {
		_, ok := opt.(indexOption)
		return ok
	})
}

// GetDefaultValue returns the default value for a field, or nil if not set.
func GetDefaultValue[T any, V any](f Field[T, V]) interface{} {
	for _, opt := range f.Options {
		if d, ok := opt.(defaultOption); ok {
			return d.Value
		}
	}
	return nil
}

// GetMaxLength returns the maximum length for a field, or nil if not set.
func GetMaxLength[T any, V any](f Field[T, V]) *int {
	for _, opt := range f.Options {
		if m, ok := opt.(maxLengthOption); ok {
			return &m.Max
		}
	}
	return nil
}

// GetMinLength returns the minimum length for a field, or nil if not set.
func GetMinLength[T any, V any](f Field[T, V]) *int {
	for _, opt := range f.Options {
		if m, ok := opt.(minLengthOption); ok {
			return &m.Min
		}
	}
	return nil
}

// GetMaxValue returns the maximum value for a field, or nil if not set.
func GetMaxValue[T any, V any](f Field[T, V]) interface{} {
	for _, opt := range f.Options {
		if m, ok := opt.(maxValueOption); ok {
			return m.Max
		}
	}
	return nil
}

// GetMinValue returns the minimum value for a field, or nil if not set.
func GetMinValue[T any, V any](f Field[T, V]) interface{} {
	for _, opt := range f.Options {
		if m, ok := opt.(minValueOption); ok {
			return m.Min
		}
	}
	return nil
}

// IsAutoNow checks if field auto-updates on save
func IsAutoNow[T any, V any](f Field[T, V]) bool {
	return HasOption(f, func(opt interface{}) bool {
		_, ok := opt.(autoNowOption)
		return ok
	})
}

// IsAutoNowAdd checks if field auto-sets on create
func IsAutoNowAdd[T any, V any](f Field[T, V]) bool {
	return HasOption(f, func(opt interface{}) bool {
		_, ok := opt.(autoNowAddOption)
		return ok
	})
}

// GetRelatedModel returns the related model name for a foreign key field, or empty string if not set.
func GetRelatedModel[T any, V any](f Field[T, V]) string {
	for _, opt := range f.Options {
		if r, ok := opt.(relatedModelOption); ok {
			return r.Model
		}
	}
	return ""
}
