package models

import (
	"context"
	"fmt"
)

// ModelValidator validates a model instance
type ModelValidator interface {
	Validate(model interface{}) error
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("field '%s': %s", e.Field, e.Message)
	}
	return e.Message
}

// ValidateModel validates a model instance according to its definition.
// Validation order:
//  1. Model-level Clean() method (if defined)
//  2. Model-level validators
// Field-level validation is handled by the field's validators when accessing field values.
func ValidateModel[T any](m *ModelDefinition[T], instance *T) error {
	if m.hooks != nil && m.hooks.Clean != nil {
		if err := m.hooks.Clean(instance); err != nil {
			return err
		}
	}

	for _, validator := range m.validators {
		if err := validator.Validate(instance); err != nil {
			return err
		}
	}

	return nil
}

// RunHooks executes model lifecycle hooks in the correct order.
// Hook types: "BeforeCreate", "AfterCreate", "BeforeUpdate", "AfterUpdate", "BeforeDelete", "AfterDelete"
func RunHooks[T any](m *ModelDefinition[T], ctx context.Context, instance *T, hookType string) error {
	if m.hooks == nil {
		return nil
	}

	switch hookType {
	case "BeforeCreate":
		if m.hooks.BeforeSave != nil {
			if err := m.hooks.BeforeSave(ctx, instance); err != nil {
				return err
			}
		}
		if m.hooks.BeforeCreate != nil {
			return m.hooks.BeforeCreate(ctx, instance)
		}
	case "AfterCreate":
		if m.hooks.AfterCreate != nil {
			if err := m.hooks.AfterCreate(ctx, instance); err != nil {
				return err
			}
		}
		if m.hooks.AfterSave != nil {
			return m.hooks.AfterSave(ctx, instance)
		}
	case "BeforeUpdate":
		if m.hooks.BeforeSave != nil {
			if err := m.hooks.BeforeSave(ctx, instance); err != nil {
				return err
			}
		}
		if m.hooks.BeforeUpdate != nil {
			return m.hooks.BeforeUpdate(ctx, instance)
		}
	case "AfterUpdate":
		if m.hooks.AfterUpdate != nil {
			if err := m.hooks.AfterUpdate(ctx, instance); err != nil {
				return err
			}
		}
		if m.hooks.AfterSave != nil {
			return m.hooks.AfterSave(ctx, instance)
		}
	case "BeforeDelete":
		if m.hooks.BeforeDelete != nil {
			return m.hooks.BeforeDelete(ctx, instance)
		}
	case "AfterDelete":
		if m.hooks.AfterDelete != nil {
			return m.hooks.AfterDelete(ctx, instance)
		}
	}

	return nil
}

