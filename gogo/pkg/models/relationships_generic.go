package models

import (
	"context"
	"fmt"
)

// Relationship defines a type-safe relationship between models
type Relationship[Parent Model, Child Model] interface {
	// Load loads the related model(s) for a parent
	Load(ctx context.Context, parent *Parent) ([]*Child, error)
	
	// Get returns related models (single for OneToOne/ManyToOne, slice for OneToMany/ManyToMany)
	Get(ctx context.Context, parent *Parent) (interface{}, error)
	
	// Set sets the relationship
	Set(ctx context.Context, parent *Parent, children []*Child) error
}

// OneToMany represents a one-to-many relationship (type-safe!)
type OneToMany[Parent Model, Child Model] struct {
	parentField *FieldRef[Parent]
	childField  *FieldRef[Child]
	childManager Manager[Child]
	foreignKey   string
}

// NewOneToMany creates a type-safe one-to-many relationship
func NewOneToMany[Parent Model, Child Model](
	parentField *FieldRef[Parent],
	childField *FieldRef[Child],
	childManager Manager[Child],
	foreignKey string,
) *OneToMany[Parent, Child] {
	return &OneToMany[Parent, Child]{
		parentField: parentField,
		childField:  childField,
		childManager: childManager,
		foreignKey:   foreignKey,
	}
}

// Load loads all children for a parent (type-safe!)
func (r *OneToMany[Parent, Child]) Load(ctx context.Context, parent *Parent) ([]*Child, error) {
	parentID := (*parent).GetID()
	
	// Use type-safe queryset
	children, err := r.childManager.Filter(ctx).
		Filter(Q{r.foreignKey: parentID}).
		All(ctx)
	
	// Returns []*Child - fully type-safe!
	return children, err
}

// Get returns related children
func (r *OneToMany[Parent, Child]) Get(ctx context.Context, parent *Parent) (interface{}, error) {
	return r.Load(ctx, parent)
}

// Set sets the relationship
func (r *OneToMany[Parent, Child]) Set(ctx context.Context, parent *Parent, children []*Child) error {
	parentID := (*parent).GetID()
	
	for _, child := range children {
		// Set foreign key on child
		// This would need reflection or a SetField method
		// For now, assume child has a SetField method
		if setter, ok := any(child).(interface{ SetField(string, interface{}) }); ok {
			setter.SetField(r.foreignKey, parentID)
		}
		
		if err := r.childManager.Save(ctx, child); err != nil {
			return err
		}
	}
	
	return nil
}

// ManyToOne represents a many-to-one relationship (type-safe!)
type ManyToOne[Parent Model, Child Model] struct {
	parentField *FieldRef[Parent]
	childField  *FieldRef[Child]
	childManager Manager[Child]
	foreignKey   string
}

// NewManyToOne creates a type-safe many-to-one relationship
func NewManyToOne[Parent Model, Child Model](
	parentField *FieldRef[Parent],
	childField *FieldRef[Child],
	childManager Manager[Child],
	foreignKey string,
) *ManyToOne[Parent, Child] {
	return &ManyToOne[Parent, Child]{
		parentField: parentField,
		childField:  childField,
		childManager: childManager,
		foreignKey:   foreignKey,
	}
}

// Load loads the related child for a parent (type-safe!)
func (r *ManyToOne[Parent, Child]) Load(ctx context.Context, parent *Parent) (*Child, error) {
	// Get foreign key value from parent
	parentID := (*parent).GetID()
	
	// Get child ID from parent's foreign key field
	// This would need reflection or a GetField method
	// For now, assume parent has a GetField method
	var childID interface{}
	if getter, ok := any(parent).(interface{ GetField(string) interface{} }); ok {
		childID = getter.GetField(r.foreignKey)
	}
	
	if childID == nil {
		return nil, nil
	}
	
	// Use type-safe manager
	child, err := r.childManager.Get(ctx, childID)
	
	// Returns *Child - fully type-safe!
	return child, err
}

// Get returns the related child
func (r *ManyToOne[Parent, Child]) Get(ctx context.Context, parent *Parent) (interface{}, error) {
	return r.Load(ctx, parent)
}

// Set sets the relationship
func (r *ManyToOne[Parent, Child]) Set(ctx context.Context, parent *Parent, child *Child) error {
	childID := (*child).GetID()
	
	// Set foreign key on parent
	if setter, ok := any(parent).(interface{ SetField(string, interface{}) }); ok {
		setter.SetField(r.foreignKey, childID)
	}
	
	return nil
}

// OneToOne represents a one-to-one relationship (type-safe!)
type OneToOne[Parent Model, Child Model] struct {
	parentField *FieldRef[Parent]
	childField  *FieldRef[Child]
	childManager Manager[Child]
	foreignKey   string
	reverse      bool // true if foreign key is on child, false if on parent
}

// NewOneToOne creates a type-safe one-to-one relationship
func NewOneToOne[Parent Model, Child Model](
	parentField *FieldRef[Parent],
	childField *FieldRef[Child],
	childManager Manager[Child],
	foreignKey string,
	reverse bool,
) *OneToOne[Parent, Child] {
	return &OneToOne[Parent, Child]{
		parentField: parentField,
		childField:  childField,
		childManager: childManager,
		foreignKey:   foreignKey,
		reverse:      reverse,
	}
}

// Load loads the related child for a parent (type-safe!)
func (r *OneToOne[Parent, Child]) Load(ctx context.Context, parent *Parent) (*Child, error) {
	parentID := (*parent).GetID()
	
	if r.reverse {
		// Foreign key is on child
		child, err := r.childManager.Filter(ctx).
			Filter(Q{r.foreignKey: parentID}).
			First(ctx)
		return child, err
	} else {
		// Foreign key is on parent
		if getter, ok := any(parent).(interface{ GetField(string) interface{} }); ok {
			childID := getter.GetField(r.foreignKey)
			if childID == nil {
				return nil, nil
			}
			return r.childManager.Get(ctx, childID)
		}
	}
	
	return nil, fmt.Errorf("unable to load one-to-one relationship")
}

// Get returns the related child
func (r *OneToOne[Parent, Child]) Get(ctx context.Context, parent *Parent) (interface{}, error) {
	return r.Load(ctx, parent)
}

// Set sets the relationship
func (r *OneToOne[Parent, Child]) Set(ctx context.Context, parent *Parent, child *Child) error {
	parentID := (*parent).GetID()
	
	if r.reverse {
		// Set foreign key on child
		if setter, ok := any(child).(interface{ SetField(string, interface{}) }); ok {
			setter.SetField(r.foreignKey, parentID)
		}
		return r.childManager.Save(ctx, child)
	} else {
		// Set foreign key on parent
		childID := (*child).GetID()
		if setter, ok := any(parent).(interface{ SetField(string, interface{}) }); ok {
			setter.SetField(r.foreignKey, childID)
		}
		return nil
	}
}

// EagerLoad loads relationships eagerly (type-safe!)
func EagerLoad[Parent Model, Child Model](
	ctx context.Context,
	parents []*Parent,
	relationship *OneToMany[Parent, Child],
) (map[interface{}][]*Child, error) {
	result := make(map[interface{}][]*Child)
	
	// Collect all parent IDs
	parentIDs := make([]interface{}, 0, len(parents))
	for _, parent := range parents {
		parentIDs = append(parentIDs, (*parent).GetID())
	}
	
	// Load all children in one query
	children, err := relationship.childManager.Filter(ctx).
		Filter(Q{relationship.foreignKey + "__in": parentIDs}).
		All(ctx)
	
	if err != nil {
		return nil, err
	}
	
	// Group children by parent ID
	for _, child := range children {
		// Get parent ID from child's foreign key
		if getter, ok := any(child).(interface{ GetField(string) interface{} }); ok {
			parentID := getter.GetField(relationship.foreignKey)
			result[parentID] = append(result[parentID], child)
		}
	}
	
	return result, nil
}

// PrefetchRelated prefetches related objects (type-safe!)
func PrefetchRelated[Parent Model, Child Model](
	ctx context.Context,
	parents []*Parent,
	relationship Relationship[Parent, Child],
) error {
	for _, parent := range parents {
		_, err := relationship.Load(ctx, parent)
		if err != nil {
			return err
		}
	}
	return nil
}

