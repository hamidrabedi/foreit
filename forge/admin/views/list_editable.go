package views

import (
	"context"
	"fmt"

	"github.com/forgego/forge/admin"
	adminutils "github.com/forgego/forge/admin/utils"
)

// ListEditableView handles inline editing in list views for the new admincore system
type ListEditableView[T any] struct {
	admin *admin.Admin[T]
}

// NewListEditableView creates a new list editable view
func NewListEditableView[T any](admin *admin.Admin[T]) *ListEditableView[T] {
	return &ListEditableView[T]{
		admin: admin,
	}
}

// ListEdit represents a single edit in the list view
type ListEdit[T any] struct {
	ObjectID  int64
	FieldName string
	Value     interface{}
}

// SaveEdits saves edits made in the list view
func (v *ListEditableView[T]) SaveEdits(ctx context.Context, edits []ListEdit[T], user interface{}) error {
	config := v.admin.Config()

	// Check if list_editable is configured
	if config == nil || len(config.ListEditable) == 0 {
		return fmt.Errorf("list_editable is not configured for this admin")
	}

	// Build map of editable field names
	editableFieldMap := make(map[string]bool)
	for _, field := range config.ListEditable {
		fieldName := v.getFieldName(field)
		if fieldName != "" {
			editableFieldMap[fieldName] = true
		}
	}

	// Validate that all edited fields are in ListEditable
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

		// Check permissions
		if !v.admin.HasChangePermission(ctx, user, obj) {
			return fmt.Errorf("permission denied to edit object %d", edit.ObjectID)
		}

		// Set the field value using reflection (since we're using field names)
		if err := adminutils.SetFieldValue(obj, edit.FieldName, edit.Value); err != nil {
			return fmt.Errorf("failed to set field %s on object %d: %w", edit.FieldName, edit.ObjectID, err)
		}

		// Save the object
		if err := v.admin.Manager().Update(ctx, obj); err != nil {
			return fmt.Errorf("failed to update object %d: %w", edit.ObjectID, err)
		}
	}

	return nil
}

// ValidateEdits validates that edits are allowed
func (v *ListEditableView[T]) ValidateEdits(ctx context.Context, edits []ListEdit[T], user interface{}) error {
	config := v.admin.Config()

	if config == nil || len(config.ListEditable) == 0 {
		return fmt.Errorf("list_editable is not configured")
	}

	// Build map of editable fields
	editableFields := make(map[string]bool)
	for _, field := range config.ListEditable {
		fieldName := v.getFieldName(field)
		if fieldName != "" {
			editableFields[fieldName] = true
		}
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

		if !v.admin.HasChangePermission(ctx, user, obj) {
			return fmt.Errorf("permission denied to edit object %d", edit.ObjectID)
		}
	}

	return nil
}

// getFieldName extracts field name from interface{} (string or FieldExpr)
func (v *ListEditableView[T]) getFieldName(field interface{}) string {
	if name, ok := field.(string); ok {
		return name
	}
	// If it's a FieldExpr, we'd need to get the name from it
	// For now, return empty
	return ""
}

// GetEditableFields returns the list of editable field names
func (v *ListEditableView[T]) GetEditableFields() []string {
	config := v.admin.Config()
	if config == nil {
		return []string{}
	}

	result := make([]string, 0, len(config.ListEditable))
	for _, field := range config.ListEditable {
		if name := v.getFieldName(field); name != "" {
			result = append(result, name)
		}
	}
	return result
}
