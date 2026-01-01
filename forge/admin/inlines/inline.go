package inlines

import (
	"context"
	"fmt"
	"reflect"

	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/orm"
)

// InlineManager manages inline editing for related models
type InlineManager[T any, R any] struct {
	admin        *admin.Admin[T]
	inline       admin.Inline[T, R]
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
	// Get parent ID
	parentID := getID(parentInstance)
	if parentID == nil {
		return nil, fmt.Errorf("parent instance has no ID")
	}

	// Resolve parent field name
	fieldName := im.resolveParentFieldName()
	if fieldName == "" {
		return nil, fmt.Errorf("invalid parent field configuration")
	}

	// Query related instances
	manager := im.relatedAdmin.Manager()

	// Create filter expression: fieldName = parentID
	// We handle both ID matching and object matching logic implicitly by ensuring
	// we filter by the correct column. If fieldName refers to relation, ORM should handle ID.
	expr := genericExpression{
		field: fieldName,
		op:    orm.OpEquals,
		value: parentID,
	}

	qs, err := manager.Manager().Filter(expr)
	if err != nil {
		return nil, err
	}

	return qs.All(ctx)
}

// SaveInlineInstances saves inline instances
func (im *InlineManager[T, R]) SaveInlineInstances(ctx context.Context, parentInstance *T, instances []*R) error {
	manager := im.relatedAdmin.Manager()

	// Resolve parent field name
	fieldName := im.resolveParentFieldName()
	if fieldName == "" {
		return fmt.Errorf("invalid parent field configuration")
	}

	for _, instance := range instances {
		// Set parent field on the instance
		if err := im.setParentField(instance, fieldName, parentInstance); err != nil {
			return fmt.Errorf("failed to set parent field: %w", err)
		}

		// Check if instance already has an ID (update vs create)
		isNew := true
		if id := getID(instance); id != nil {
			// If generic ID check returns non-zero/non-nil, it's an update
			// (simplified check, assume int64 > 0)
			if intVal, ok := id.(int64); ok && intVal > 0 {
				isNew = false
			}
		}

		// Save instance
		if isNew {
			if err := manager.Create(ctx, instance); err != nil {
				return fmt.Errorf("failed to create inline instance: %w", err)
			}
		} else {
			if err := manager.Update(ctx, instance); err != nil {
				return fmt.Errorf("failed to update inline instance: %w", err)
			}
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

// Helper methods

func (im *InlineManager[T, R]) resolveParentFieldName() string {
	parentField := im.inline.GetParentField()
	if name, ok := parentField.(string); ok {
		return name
	}
	// Try reflection for Named interface (like FieldExpr)
	val := reflect.ValueOf(parentField)
	method := val.MethodByName("Name")
	if method.IsValid() && method.Type().NumOut() == 1 && method.Type().Out(0).Kind() == reflect.String {
		results := method.Call(nil)
		return results[0].String()
	}
	return ""
}

func (im *InlineManager[T, R]) setParentField(instance *R, fieldName string, parentV *T) error {
	val := reflect.ValueOf(instance).Elem()
	field := val.FieldByName(fieldName)

	if !field.IsValid() {
		return fmt.Errorf("field %s not found on inline model", fieldName)
	}

	if !field.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}

	// Determine what to set: the parent object or the parent ID?
	// Check field type
	if field.Type() == reflect.TypeOf(parentV) {
		// Field expects the pointer to parent struct
		field.Set(reflect.ValueOf(parentV))
		return nil
	}

	// Check if field expects ID
	parentID := getID(parentV)
	if parentID == nil {
		return fmt.Errorf("parent has no ID")
	}

	parentIDVal := reflect.ValueOf(parentID)
	if field.Type() == parentIDVal.Type() {
		field.Set(parentIDVal)
		return nil
	}

	// Try converting ID (e.g. int to int64)
	if field.Kind() == reflect.Int64 && parentIDVal.Kind() == reflect.Int {
		field.SetInt(parentIDVal.Int())
		return nil
	}

	return fmt.Errorf("type mismatch for parent field %s: expected %v, got %v or %T",
		fieldName, field.Type(), reflect.TypeOf(parentV), parentID)
}

func getID(instance interface{}) interface{} {
	// Try ModelWithID interface
	type ModelWithID interface {
		GetID() int64
	}
	if m, ok := instance.(ModelWithID); ok {
		return m.GetID()
	}

	// Reflection fallback
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	field := val.FieldByName("ID")
	if !field.IsValid() {
		field = val.FieldByName("id")
	}

	if field.IsValid() {
		return field.Interface()
	}
	return nil
}

// genericExpression allows filtering without strict schema type checking
// This is useful for admin inlines where we know the relationship exists but
// don't have easy access to the exact field type for NewField[T].
type genericExpression struct {
	field string
	op    orm.Operator
	value interface{}
}

func (e genericExpression) ToSQL(builder *orm.SQLBuilder) (string, []interface{}, error) {
	// Add argument to builder
	placeholder := builder.AddArg(e.value)
	// Return SQL fragment and the args used in this fragment
	// Note: builder.AddArg already added it to builder.args, but ToSQL interface requires returning them too?
	// Based on expression.go, yes.
	return fmt.Sprintf("%s %s %s", orm.EscapeIdentifier(e.field), e.op, placeholder), []interface{}{e.value}, nil
}

func (e genericExpression) Resolve(schema *orm.ModelSchema) error {
	// Skip strict validation or implement basic existence check
	if schema.GetField(e.field) == nil {
		// Try to see if it's a relation field (basic check)
		// For now, allow it to pass to let ORM handle it at SQL level or error later
		return nil
	}
	return nil
}
