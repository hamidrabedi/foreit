package admin

import (
	"context"
	"fmt"
)

// ListEditableView handles inline editing in list views
type ListEditableView[T any] struct {
	admin *Admin[T]
}

// NewListEditableView creates a new list editable view
func NewListEditableView[T any](admin *Admin[T]) *ListEditableView[T] {
	return &ListEditableView[T]{
		admin: admin,
	}
}

// SaveEdits saves edits made in the list view
func (v *ListEditableView[T]) SaveEdits(ctx context.Context, edits []ListEdit[T]) error {
	config := v.admin.Config()

	// Check if list_editable is configured
	if len(config.ListEditable) == 0 {
		return fmt.Errorf("list_editable is not configured for this admin")
	}

	// Validate that all edited fields are in ListEditable
	editableFieldMap := make(map[string]bool)
	for _, field := range config.ListEditable {
		editableFieldMap[field.Name()] = true
	}

	for _, edit := range edits {
		if !editableFieldMap[edit.FieldName] {
			return fmt.Errorf("field %s is not editable in list view", edit.FieldName)
		}
	}

	// Apply edits to each object
	for _, edit := range edits {
		// Get the object
		obj, err := v.admin.Manager().Get(ctx, edit.ObjectID)
		if err != nil {
			return fmt.Errorf("failed to get object %d: %w", edit.ObjectID, err)
		}

		// Find the field expression
		var fieldExpr FieldExpr[T, interface{}]
		for _, f := range config.ListEditable {
			if f.Name() == edit.FieldName {
				fieldExpr = f
				break
			}
		}

		if fieldExpr.Name() == "" {
			return fmt.Errorf("field %s not found", edit.FieldName)
		}

		// Set the field value
		fieldExpr.Set(obj, edit.Value)

		// Save the object
		if err := v.admin.Manager().Update(ctx, obj); err != nil {
			return fmt.Errorf("failed to update object %d: %w", edit.ObjectID, err)
		}
	}

	return nil
}

// ListEdit represents a single edit in the list view
type ListEdit[T any] struct {
	ObjectID  int64
	FieldName string
	Value     interface{}
}

// ValidateEdits validates that edits are allowed
func (v *ListEditableView[T]) ValidateEdits(ctx context.Context, edits []ListEdit[T]) error {
	config := v.admin.Config()

	if len(config.ListEditable) == 0 {
		return fmt.Errorf("list_editable is not configured")
	}

	// Build map of editable fields
	editableFields := make(map[string]bool)
	for _, field := range config.ListEditable {
		editableFields[field.Name()] = true
	}

	// Validate each edit
	for _, edit := range edits {
		if !editableFields[edit.FieldName] {
			return fmt.Errorf("field %s is not in list_editable", edit.FieldName)
		}

		// Check permissions
		obj, err := v.admin.Manager().Get(ctx, edit.ObjectID)
		if err != nil {
			return fmt.Errorf("failed to get object %d: %w", edit.ObjectID, err)
		}

		// Get user from context (would need to be set by middleware)
		// For now, just check basic permission
		if !v.admin.HasChangePermission(ctx, nil, obj) {
			return fmt.Errorf("permission denied to edit object %d", edit.ObjectID)
		}
	}

	return nil
}
