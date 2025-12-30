package admin

import (
	"context"
	"fmt"
)

// DetailView represents a type-safe detail view
type DetailView[T any] struct {
	admin    *Admin[T]
	instance *T
}

// DetailData contains data for rendering the detail view
type DetailData[T any] struct {
	Instance *T
	Fields   []FieldDisplay[T]
}

// FieldDisplay represents a field to display
type FieldDisplay[T any] struct {
	Name  string
	Label string
	Value interface{}
}

// NewDetailView creates a new detail view
func NewDetailView[T any](admin *Admin[T], instance *T) *DetailView[T] {
	return &DetailView[T]{
		admin:    admin,
		instance: instance,
	}
}

// Render renders the detail view and returns the data
func (v *DetailView[T]) Render(ctx context.Context) (*DetailData[T], error) {
	config := v.admin.Config()
	fields := make([]FieldDisplay[T], 0, len(config.ListDisplay))

	// Use list display fields if available, otherwise use all fields
	displayFields := config.ListDisplay
	if len(displayFields) == 0 {
		// TODO: Get all fields from schema
		return &DetailData[T]{
			Instance: v.instance,
			Fields:   fields,
		}, nil
	}

	// Build field displays
	for _, field := range displayFields {
		value := field.Get(v.instance)
		fields = append(fields, FieldDisplay[T]{
			Name:  field.Name(),
			Label: field.Name(), // TODO: Get label from field definition
			Value: value,
		})
	}

	return &DetailData[T]{
		Instance: v.instance,
		Fields:   fields,
	}, nil
}

// GetDetailViewByID gets an instance by ID and creates a detail view
func GetDetailViewByID[T any](admin *Admin[T], ctx context.Context, id int64) (*DetailView[T], error) {
	instance, err := admin.Manager().Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	return NewDetailView(admin, instance), nil
}
