package relations

import (
	"context"
	"fmt"
	"reflect"

	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/orm"
)

// RelationManager manages related models in a separate tab/section
type RelationManager[T any, R any] struct {
	name         string
	label        string
	admin        *admin.Admin[T]
	relatedAdmin *admin.Admin[R]
	relationName string // Field name of the relation
	relationType string // "ForeignKey", "OneToOne", "ManyToMany"
}

// NewRelationManager creates a new relation manager
func NewRelationManager[T any, R any](
	name string,
	admin *admin.Admin[T],
	relatedAdmin *admin.Admin[R],
	relationName string,
) *RelationManager[T, R] {
	return &RelationManager[T, R]{
		name:         name,
		label:        name,
		admin:        admin,
		relatedAdmin: relatedAdmin,
		relationName: relationName,
	}
}

// Name returns the relation manager name
func (rm *RelationManager[T, R]) Name() string {
	return rm.name
}

// Label returns the relation manager label
func (rm *RelationManager[T, R]) Label() string {
	return rm.label
}

// GetQueryset gets the queryset for related objects
func (rm *RelationManager[T, R]) GetQueryset(ctx context.Context, parent *T) (orm.QuerySet[R], error) {
	// Get the related manager
	relatedManager := rm.relatedAdmin.Manager()
	if relatedManager == nil {
		return nil, fmt.Errorf("related admin has no manager")
	}

	// Get parent ID
	parentID := rm.getParentID(parent)
	if parentID == nil {
		return nil, fmt.Errorf("could not get parent ID")
	}

	// Create queryset filtered by relation
	// This is simplified - full implementation would:
	// 1. Get the relation field from parent
	// 2. Filter related objects by that field
	// 3. Return filtered queryset

	qs, err := orm.NewQuerySet[R]("")
	if err != nil {
		return nil, err
	}

	// Filter by relation field
	// In a full implementation, this would use:
	// fieldExpr := orm.NewField[int64](rm.relationName, "")
	// qs = qs.Filter(fieldExpr.Eq(parentID))

	return qs, nil
}

// getParentID gets the ID from the parent object
func (rm *RelationManager[T, R]) getParentID(parent *T) interface{} {
	val := reflect.ValueOf(parent)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	f := val.FieldByName("ID")
	if f.IsValid() {
		return f.Interface()
	}
	return nil
}

// CanCreate checks if user can create related objects
func (rm *RelationManager[T, R]) CanCreate(ctx context.Context, user interface{}, parent *T) bool {
	return rm.relatedAdmin.HasAddPermission(ctx, user)
}

// CanEdit checks if user can edit related objects
func (rm *RelationManager[T, R]) CanEdit(ctx context.Context, user interface{}, parent *T, related *R) bool {
	return rm.relatedAdmin.HasChangePermission(ctx, user, related)
}

// CanDelete checks if user can delete related objects
func (rm *RelationManager[T, R]) CanDelete(ctx context.Context, user interface{}, parent *T, related *R) bool {
	return rm.relatedAdmin.HasDeletePermission(ctx, user, related)
}

// CanView checks if user can view related objects
func (rm *RelationManager[T, R]) CanView(ctx context.Context, user interface{}, parent *T, related *R) bool {
	return rm.relatedAdmin.HasViewPermission(ctx, user, related)
}

// RelationManagerConfig contains configuration for a relation manager
type RelationManagerConfig[T any, R any] struct {
	Name         string
	Label        string
	RelationName string
	DisplayFields []string
	SearchFields []string
	Filters      []interface{}
	Ordering     []string
	PerPage      int
}

// NewRelationManagerConfig creates a new relation manager config
func NewRelationManagerConfig[T any, R any](name, relationName string) *RelationManagerConfig[T, R] {
	return &RelationManagerConfig[T, R]{
		Name:         name,
		Label:        name,
		RelationName: relationName,
		PerPage:      25,
	}
}
