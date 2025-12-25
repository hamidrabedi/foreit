package models

import (
	"context"
	"errors"
)

// Relationship represents a relationship between models (SQLAlchemy-inspired)
type Relationship interface {
	// Load loads the related model(s)
	Load(ctx context.Context, model Model) error
	
	// Get returns the related model(s)
	Get(ctx context.Context, model Model) (interface{}, error)
}

// RelationshipType represents the type of relationship
type RelationshipType string

const (
	RelationshipOneToOne   RelationshipType = "one_to_one"
	RelationshipOneToMany  RelationshipType = "one_to_many"
	RelationshipManyToOne  RelationshipType = "many_to_one"
	RelationshipManyToMany RelationshipType = "many_to_many"
)

// ForeignKey represents a foreign key relationship
type ForeignKey struct {
	Field      string
	TargetType string
	OnDelete   string // CASCADE, SET_NULL, RESTRICT
	OnUpdate   string
}

// OneToOne represents a one-to-one relationship
type OneToOne struct {
	Field      string
	TargetType string
	ForeignKey *ForeignKey
	Lazy       bool
}

// Load loads the related model
func (r *OneToOne) Load(ctx context.Context, model Model) error {
	// Implementation depends on the manager
	return errors.New("not implemented")
}

// Get returns the related model
func (r *OneToOne) Get(ctx context.Context, model Model) (interface{}, error) {
	return nil, errors.New("not implemented")
}

// OneToMany represents a one-to-many relationship
type OneToMany struct {
	Field      string
	TargetType string
	ForeignKey *ForeignKey
	Lazy       bool
}

// Load loads the related models
func (r *OneToMany) Load(ctx context.Context, model Model) error {
	return errors.New("not implemented")
}

// Get returns the related models
func (r *OneToMany) Get(ctx context.Context, model Model) (interface{}, error) {
	return nil, errors.New("not implemented")
}

// ManyToOne represents a many-to-one relationship
type ManyToOne struct {
	Field      string
	TargetType string
	ForeignKey *ForeignKey
	Lazy       bool
}

// Load loads the related model
func (r *ManyToOne) Load(ctx context.Context, model Model) error {
	return errors.New("not implemented")
}

// Get returns the related model
func (r *ManyToOne) Get(ctx context.Context, model Model) (interface{}, error) {
	return nil, errors.New("not implemented")
}

// EagerLoad loads relationships eagerly (like select_related/prefetch_related)
func EagerLoad(ctx context.Context, model Model, relationships ...string) error {
	// Implementation would use the manager to load relationships
	return errors.New("not implemented")
}

// LazyLoad loads relationships lazily (on access)
func LazyLoad(ctx context.Context, model Model, relationship string) (interface{}, error) {
	// Implementation would use the manager to load the relationship
	return nil, errors.New("not implemented")
}

