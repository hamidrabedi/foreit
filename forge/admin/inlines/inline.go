package inlines

import (
	"context"
	"fmt"

	"github.com/forgego/forge/admin"
)

// InlineManager manages inline editing for related models
type InlineManager[T any, R any] struct {
	admin      *admin.Admin[T]
	inline     admin.Inline[T, R]
	relatedAdmin *admin.Admin[R]
}

// NewInlineManager creates a new inline manager
func NewInlineManager[T any, R any](
	admin *admin.Admin[T],
	inline admin.Inline[T, R],
	relatedAdmin *admin.Admin[R],
) *InlineManager[T, R] {
	return &InlineManager[T, R]{
		admin:        admin,
		inline:       inline,
		relatedAdmin: relatedAdmin,
	}
}

// GetInlineInstances gets related instances for inline editing
func (im *InlineManager[T, R]) GetInlineInstances(ctx context.Context, parentInstance *T) ([]*R, error) {
	// Get parent field value (would use ORM field accessor)
	// For now, this is a placeholder
	_ = parentInstance
	
	// Query related instances
	manager := im.relatedAdmin.Manager()
	qs, err := manager.GetQueryset(ctx)
	if err != nil {
		return nil, err
	}
	
	// Filter by parent (would use proper field expression)
	// qs = qs.Filter(parentField.Eq(parentID))
	
	return qs.All(ctx)
}

// SaveInlineInstances saves inline instances
func (im *InlineManager[T, R]) SaveInlineInstances(ctx context.Context, parentInstance *T, instances []*R) error {
	manager := im.relatedAdmin.Manager()
	
	for _, instance := range instances {
		// Set parent field (would use ORM field accessor)
		// setParentField(instance, parentInstance)
		
		// Save instance
		if err := manager.Create(ctx, instance); err != nil {
			return fmt.Errorf("failed to save inline instance: %w", err)
		}
	}
	
	return nil
}

// DeleteInlineInstance deletes an inline instance
func (im *InlineManager[T, R]) DeleteInlineInstance(ctx context.Context, instance *R) error {
	manager := im.relatedAdmin.Manager()
	return manager.Delete(ctx, instance)
}

// TabularInline renders inline instances in tabular format
func (im *InlineManager[T, R]) TabularInline(ctx context.Context, parentInstance *T) (*TabularInlineData[R], error) {
	instances, err := im.GetInlineInstances(ctx, parentInstance)
	if err != nil {
		return nil, err
	}
	
	return &TabularInlineData[R]{
		Instances: instances,
		Fields:    im.getInlineFields(),
		Extra:     im.inline.GetExtra(),
		MaxNum:    im.inline.GetMaxNum(),
	}, nil
}

// StackedInline renders inline instances in stacked format
func (im *InlineManager[T, R]) StackedInline(ctx context.Context, parentInstance *T) (*StackedInlineData[R], error) {
	instances, err := im.GetInlineInstances(ctx, parentInstance)
	if err != nil {
		return nil, err
	}
	
	return &StackedInlineData[R]{
		Instances: instances,
		Fields:    im.getInlineFields(),
		Extra:     im.inline.GetExtra(),
		MaxNum:    im.inline.GetMaxNum(),
	}, nil
}

// getInlineFields gets field names for inline display
func (im *InlineManager[T, R]) getInlineFields() []string {
	inlineFields := im.inline.GetFields()
	fields := make([]string, 0, len(inlineFields))
	for _, field := range inlineFields {
		if fieldName, ok := field.(string); ok {
			fields = append(fields, fieldName)
		}
	}
	return fields
}

// TabularInlineData contains data for tabular inline
type TabularInlineData[R any] struct {
	Instances []*R
	Fields    []string
	Extra     int
	MaxNum    int
}

// StackedInlineData contains data for stacked inline
type StackedInlineData[R any] struct {
	Instances []*R
	Fields    []string
	Extra     int
	MaxNum    int
}
