package models

// FieldAccessor provides type-safe field access for relationships
type FieldAccessor interface {
	// GetField gets a field value by name
	GetField(name string) interface{}
	
	// SetField sets a field value by name
	SetField(name string, value interface{})
}

// FieldAccessorImpl provides a default implementation using reflection
// This is used internally for relationships
type FieldAccessorImpl struct {
	model Model
}

// NewFieldAccessor creates a field accessor for a model
func NewFieldAccessor(model Model) *FieldAccessorImpl {
	return &FieldAccessorImpl{model: model}
}

// GetField gets a field value (uses reflection internally, but API is type-safe)
func (a *FieldAccessorImpl) GetField(name string) interface{} {
	// This would use reflection to get the field
	// For now, return nil - actual implementation would use reflect package
	return nil
}

// SetField sets a field value (uses reflection internally, but API is type-safe)
func (a *FieldAccessorImpl) SetField(name string, value interface{}) {
	// This would use reflection to set the field
	// Actual implementation would use reflect package
}

// WithFieldAccessor wraps a model to provide field access
func WithFieldAccessor[T Model](model *T) *T {
	// Models can implement FieldAccessor directly for better performance
	// This is just a helper
	return model
}

