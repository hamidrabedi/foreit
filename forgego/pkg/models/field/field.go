package field

// Field represents a typed field definition with getter/setter functions.
// T is the model type, V is the field value type.
type Field[T any, V any] struct {
	Name       string
	Column     string
	Getter     func(*T) *V
	Setter     func(*T, V)
	Options    []Option
	Validators []Validator
}

// NewField creates a new Field with the given name, column, getter, setter, and options.
func NewField[T any, V any](
	name string,
	column string,
	getter func(*T) *V,
	setter func(*T, V),
	opts ...Option,
) Field[T, V] {
	f := Field[T, V]{
		Name:       name,
		Column:     column,
		Getter:     getter,
		Setter:     setter,
		Options:    []Option{},
		Validators: []Validator{},
	}

	// Apply options
	for _, opt := range opts {
		f.applyOption(opt)
	}

	return f
}

// applyOption applies an option to the field.
func (f *Field[T, V]) applyOption(opt Option) {
	switch o := opt.(type) {
	case columnOption:
		f.Column = o.Name
	case emailOption:
		f.Validators = append(f.Validators, EmailValidator{})
	case maxLengthOption:
		f.Validators = append(f.Validators, MaxLengthValidator{Max: o.Max})
	case minLengthOption:
		f.Validators = append(f.Validators, MinLengthValidator{Min: o.Min})
	case maxValueOption:
		f.Validators = append(f.Validators, MaxValueValidator{Max: o.Max})
	case minValueOption:
		f.Validators = append(f.Validators, MinValueValidator{Min: o.Min})
	case requiredOption:
		f.Validators = append(f.Validators, RequiredValidator{})
	default:
		f.Options = append(f.Options, opt)
	}
}

// Get returns the field value from the model instance using the getter function.
func (f Field[T, V]) Get(instance *T) *V {
	if f.Getter == nil {
		return nil
	}
	return f.Getter(instance)
}

// Set sets the field value on the model instance using the setter function.
func (f Field[T, V]) Set(instance *T, value V) {
	if f.Setter != nil {
		f.Setter(instance, value)
	}
}

// WithOptions returns a new Field with additional options applied.
func (f Field[T, V]) WithOptions(opts ...Option) Field[T, V] {
	for _, opt := range opts {
		f.applyOption(opt)
	}
	return f
}

// WithValidators returns a new Field with additional validators.
func (f Field[T, V]) WithValidators(validators ...Validator) Field[T, V] {
	f.Validators = append(f.Validators, validators...)
	return f
}

// GetName returns the field name. Implements FieldDefinition interface.
func (f Field[T, V]) GetName() string {
	return f.Name
}

// GetColumn returns the database column name. Implements FieldDefinition interface.
func (f Field[T, V]) GetColumn() string {
	return f.Column
}

// GetValidators returns validators as []interface{} to implement FieldDefinitionWithValidators.
func (f Field[T, V]) GetValidators() []interface{} {
	validators := make([]interface{}, len(f.Validators))
	for i, v := range f.Validators {
		validators[i] = v
	}
	return validators
}

