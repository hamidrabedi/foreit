package models

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Model defines the interface that all models must implement
// Similar to Django's models.Model
type Model interface {
	// Save saves the model instance (create or update)
	Save(ctx context.Context) error
	
	// Delete deletes the model instance
	Delete(ctx context.Context) error
	
	// String returns a string representation of the model
	String() string
	
	// IsNew returns true if the model hasn't been saved yet
	IsNew() bool
	
	// GetID returns the primary key value
	GetID() interface{}
	
	// SetID sets the primary key value
	SetID(id interface{})
	
	// GetCreatedAt returns the creation timestamp
	GetCreatedAt() *time.Time
	
	// GetUpdatedAt returns the last update timestamp
	GetUpdatedAt() *time.Time
}

// ModelWithHooks extends Model with lifecycle hooks
type ModelWithHooks interface {
	Model
	
	// BeforeSave is called before saving (create or update)
	BeforeSave(ctx context.Context) error
	
	// AfterSave is called after saving (create or update)
	AfterSave(ctx context.Context) error
	
	// BeforeDelete is called before deleting
	BeforeDelete(ctx context.Context) error
	
	// AfterDelete is called after deleting
	AfterDelete(ctx context.Context) error
	
	// BeforeCreate is called before creating
	BeforeCreate(ctx context.Context) error
	
	// AfterCreate is called after creating
	AfterCreate(ctx context.Context) error
	
	// BeforeUpdate is called before updating
	BeforeUpdate(ctx context.Context) error
	
	// AfterUpdate is called after updating
	AfterUpdate(ctx context.Context) error
}

// BaseModel provides a base implementation that can be embedded
// Similar to Django's models.Model base class
type BaseModel struct {
	ID        interface{} `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	isNew     bool
	manager   Manager
}

// NewBaseModel creates a new base model
func NewBaseModel() *BaseModel {
	now := time.Now()
	return &BaseModel{
		isNew:     true,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

// GetID returns the primary key
func (m *BaseModel) GetID() interface{} {
	return m.ID
}

// SetID sets the primary key
func (m *BaseModel) SetID(id interface{}) {
	m.ID = id
	m.isNew = false
}

// IsNew returns true if the model hasn't been saved
func (m *BaseModel) IsNew() bool {
	return m.isNew
}

// GetCreatedAt returns the creation timestamp
func (m *BaseModel) GetCreatedAt() *time.Time {
	return m.CreatedAt
}

// GetUpdatedAt returns the update timestamp
func (m *BaseModel) GetUpdatedAt() *time.Time {
	return m.UpdatedAt
}

// String returns a string representation
func (m *BaseModel) String() string {
	if m.ID != nil {
		return "Model(id=" + toString(m.ID) + ")"
	}
	return "Model(unsaved)"
}

// Save saves the model (delegates to manager)
func (m *BaseModel) Save(ctx context.Context) error {
	if m.manager == nil {
		return ErrNoManager
	}
	
	// Validate before saving
	if validatable, ok := m.toModel().(Validatable); ok {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}
	
	if hookModel, ok := m.toModel().(ModelWithHooks); ok {
		if err := hookModel.BeforeSave(ctx); err != nil {
			return err
		}
		
		if m.isNew {
			if err := hookModel.BeforeCreate(ctx); err != nil {
				return err
			}
		} else {
			if err := hookModel.BeforeUpdate(ctx); err != nil {
				return err
			}
		}
	}
	
	// Send signals
	if m.isNew {
		PreCreate.Send(ctx, m.toModel(), SignalPreCreate)
	} else {
		PreUpdate.Send(ctx, m.toModel(), SignalPreUpdate)
	}
	PreSave.Send(ctx, m.toModel(), SignalPreSave)
	
	now := time.Now()
	if m.isNew {
		if m.CreatedAt == nil {
			m.CreatedAt = &now
		}
	}
	m.UpdatedAt = &now
	
	if err := m.manager.Save(ctx, m.toModel()); err != nil {
		return err
	}
	
	m.isNew = false
	
	// Send signals
	if !m.isNew {
		PostCreate.Send(ctx, m.toModel(), SignalPostCreate)
	} else {
		PostUpdate.Send(ctx, m.toModel(), SignalPostUpdate)
	}
	PostSave.Send(ctx, m.toModel(), SignalPostSave)
	
	if hookModel, ok := m.toModel().(ModelWithHooks); ok {
		if err := hookModel.AfterSave(ctx); err != nil {
			return err
		}
		
		if !m.isNew {
			if err := hookModel.AfterCreate(ctx); err != nil {
				return err
			}
		} else {
			if err := hookModel.AfterUpdate(ctx); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// Delete deletes the model (delegates to manager)
func (m *BaseModel) Delete(ctx context.Context) error {
	if m.manager == nil {
		return ErrNoManager
	}
	
	// Send signal
	PreDelete.Send(ctx, m.toModel(), SignalPreDelete)
	
	if hookModel, ok := m.toModel().(ModelWithHooks); ok {
		if err := hookModel.BeforeDelete(ctx); err != nil {
			return err
		}
	}
	
	if err := m.manager.Delete(ctx, m.toModel()); err != nil {
		return err
	}
	
	// Send signal
	PostDelete.Send(ctx, m.toModel(), SignalPostDelete)
	
	if hookModel, ok := m.toModel().(ModelWithHooks); ok {
		if err := hookModel.AfterDelete(ctx); err != nil {
			return err
		}
	}
	
	return nil
}

// SetManager sets the manager for this model
func (m *BaseModel) SetManager(manager Manager) {
	m.manager = manager
}

// toModel converts BaseModel to Model interface
func (m *BaseModel) toModel() Model {
	return m
}

// toString converts a value to string
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Manager handles database operations for models
type Manager interface {
	// Save saves a model (create or update)
	Save(ctx context.Context, model Model) error
	
	// Delete deletes a model
	Delete(ctx context.Context, model Model) error
	
	// Get retrieves a model by ID
	Get(ctx context.Context, id interface{}) (Model, error)
	
	// All returns all models
	All(ctx context.Context) ([]Model, error)
	
	// Filter returns a queryset for filtering
	Filter(ctx context.Context) QuerySet
}

// QuerySet represents a chainable query set
type QuerySet interface {
	// Filter adds a filter condition
	Filter(condition interface{}) QuerySet
	
	// Exclude adds an exclude condition
	Exclude(condition interface{}) QuerySet
	
	// OrderBy adds ordering
	OrderBy(fields ...string) QuerySet
	
	// Limit limits the number of results
	Limit(n int) QuerySet
	
	// Offset sets the offset
	Offset(n int) QuerySet
	
	// Get retrieves a single result
	Get(ctx context.Context) (Model, error)
	
	// All retrieves all results
	All(ctx context.Context) ([]Model, error)
	
	// Count returns the count
	Count(ctx context.Context) (int, error)
	
	// Exists checks if any results exist
	Exists(ctx context.Context) (bool, error)
	
	// First returns the first result
	First(ctx context.Context) (Model, error)
	
	// Last returns the last result
	Last(ctx context.Context) (Model, error)
}

// Validatable interface for models that can validate themselves
type Validatable interface {
	Validate() error
}

var (
	ErrNoManager = errors.New("model has no manager set")
)

