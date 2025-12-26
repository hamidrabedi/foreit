package query

import (
	"context"
	"fmt"
	"reflect"

	"github.com/forgego/forge/pkg/db"
	"github.com/forgego/forge/pkg/errors"
)

// Manager provides type-safe CRUD operations for a model
type Manager[T any] struct {
	tableName string
	db        *db.DB
}

// NewManager creates a new generic manager for a model type
func NewManager[T any](tableName string) *Manager[T] {
	return &Manager[T]{
		tableName: tableName,
	}
}

// SetDB sets the database connection for this manager
func (m *Manager[T]) SetDB(database *db.DB) {
	m.db = database
}

// Filter returns a QuerySet for filtering
func (m *Manager[T]) Filter(expr ...QueryExpr) QuerySet[T] {
	qs := NewBaseQuerySet[T](m.tableName)
	if m.db != nil {
		qs.SetDB(m.db)
	}
	if len(expr) > 0 {
		return qs.Filter(expr[0])
	}
	return qs
}

// Get retrieves a model by ID
func (m *Manager[T]) Get(ctx context.Context, id int64) (*T, error) {
	if m.db == nil {
		return nil, errors.NewNotImplementedError("Manager.Get() - database connection not set")
	}
	return m.Filter(NewFieldQueryExpr("id", OpEquals, id)).Get(ctx)
}

// All returns all model instances
func (m *Manager[T]) All(ctx context.Context) ([]*T, error) {
	return m.Filter().All(ctx)
}

// Create creates a new model instance
func (m *Manager[T]) Create(ctx context.Context, instance *T) error {
	if m.db == nil {
		return errors.NewNotImplementedError("Manager.Create() - database connection not set")
	}

	// Run BeforeCreate hooks
	if err := ExecuteHooks(ctx, instance, "BeforeCreate"); err != nil {
		return err
	}

	// Validate instance
	if err := ValidateInstance(instance); err != nil {
		return err
	}

	// Build and execute INSERT
	sql, args, _, err := BuildInsertSQL(instance, m.tableName)
	if err != nil {
		return fmt.Errorf("failed to build insert SQL: %w", err)
	}

	// Execute INSERT and get generated ID
	id, err := ExecuteInsert(ctx, m.db, sql, args)
	if err != nil {
		return err
	}

	// Set the ID on the instance using type assertion if it implements ModelWithID
	if modelWithID, ok := any(instance).(ModelWithID); ok {
		modelWithID.SetID(id)
	} else {
		// Fallback to reflection only if needed
		instanceValue := reflect.ValueOf(instance).Elem()
		idField := instanceValue.FieldByName("ID")
		if !idField.IsValid() {
			idField = instanceValue.FieldByName("id")
		}
		if idField.IsValid() && idField.CanSet() {
			idField.SetInt(id)
		}
	}

	// Run AfterCreate hooks
	if err := ExecuteHooks(ctx, instance, "AfterCreate"); err != nil {
		return err
	}

	return nil
}

// Update updates an existing model instance
func (m *Manager[T]) Update(ctx context.Context, instance *T) error {
	if m.db == nil {
		return errors.NewNotImplementedError("Manager.Update() - database connection not set")
	}

	// Get ID value using type assertion if available
	var id int64
	if modelWithID, ok := any(instance).(ModelWithID); ok {
		id = modelWithID.GetID()
	} else {
		// Fallback to reflection
		idValue, err := GetIDValue(instance, "id")
		if err != nil {
			return fmt.Errorf("failed to get ID: %w", err)
		}
		var ok bool
		id, ok = idValue.(int64)
		if !ok {
			return fmt.Errorf("ID must be int64, got %T", idValue)
		}
	}

	if id == 0 {
		return errors.NewInvalidInputError("id", "ID must be non-zero for update")
	}

	// Run BeforeUpdate hooks
	if err := ExecuteHooks(ctx, instance, "BeforeUpdate"); err != nil {
		return err
	}

	// Validate instance
	if err := ValidateInstance(instance); err != nil {
		return err
	}

	// Build and execute UPDATE
	sql, args, err := BuildUpdateSQL(instance, m.tableName, "id")
	if err != nil {
		return fmt.Errorf("failed to build update SQL: %w", err)
	}

	// Execute UPDATE
	rowsAffected, err := ExecuteUpdate(ctx, m.db, sql, args)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("model", id)
	}

	// Run AfterUpdate hooks
	if err := ExecuteHooks(ctx, instance, "AfterUpdate"); err != nil {
		return err
	}

	return nil
}

// Save saves a model (create or update)
func (m *Manager[T]) Save(ctx context.Context, instance *T) error {
	// Use type assertion if available, otherwise fallback to reflection
	if modelWithID, ok := any(instance).(ModelWithID); ok {
		if modelWithID.GetID() != 0 {
			return m.Update(ctx, instance)
		}
		return m.Create(ctx, instance)
	}
	
	// Fallback to reflection
	idValue, err := GetIDValue(instance, "id")
	if err == nil {
		if id, ok := idValue.(int64); ok && id != 0 {
			return m.Update(ctx, instance)
		}
	}
	return m.Create(ctx, instance)
}

// Delete deletes a model instance
func (m *Manager[T]) Delete(ctx context.Context, instance *T) error {
	if m.db == nil {
		return errors.NewNotImplementedError("Manager.Delete() - database connection not set")
	}

	// Get ID value using type assertion if available
	var id int64
	if modelWithID, ok := any(instance).(ModelWithID); ok {
		id = modelWithID.GetID()
	} else {
		// Fallback to reflection
		idValue, err := GetIDValue(instance, "id")
		if err != nil {
			return fmt.Errorf("failed to get ID: %w", err)
		}
		var ok bool
		id, ok = idValue.(int64)
		if !ok {
			return fmt.Errorf("ID must be int64, got %T", idValue)
		}
	}

	if id == 0 {
		return errors.NewInvalidInputError("id", "ID must be non-zero for delete")
	}

	// Run BeforeDelete hooks
	if err := ExecuteHooks(ctx, instance, "BeforeDelete"); err != nil {
		return err
	}

	// Build and execute DELETE
	sql, args := BuildDeleteSQL(m.tableName, "id", id)

	// Execute DELETE
	rowsAffected, err := ExecuteDelete(ctx, m.db, sql, args)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("model", id)
	}

	// Run AfterDelete hooks
	if err := ExecuteHooks(ctx, instance, "AfterDelete"); err != nil {
		return err
	}

	return nil
}

