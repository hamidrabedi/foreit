package orm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/forgego/forge/db"
	"github.com/forgego/forge/errors"
)

// Manager provides type-safe CRUD operations
type Manager[T any] struct {
	tableName string
	schema    *ModelSchema
	db        *db.DB
}

// NewManager creates a new manager
func NewManager[T any](tableName string) (*Manager[T], error) {
	schema, err := GetModelSchema[T]()
	if err != nil {
		return nil, err
	}

	if tableName == "" {
		tableName = schema.TableName
	}

	return &Manager[T]{
		tableName: tableName,
		schema:    schema,
	}, nil
}

// SetDB sets the database connection
func (m *Manager[T]) SetDB(database *db.DB) {
	m.db = database
}

// GetFieldAccessor returns a field accessor for this model
func (m *Manager[T]) GetFieldAccessor() (*FieldAccessor[T], error) {
	return NewFieldAccessor[T]()
}

// Filter returns a QuerySet for filtering
func (m *Manager[T]) Filter(expr Expression) (QuerySet[T], error) {
	qs, err := NewQuerySet[T](m.tableName)
	if err != nil {
		return nil, err
	}

	if m.db != nil {
		qs = qs.SetDB(m.db).(QuerySet[T])
	}

	return qs.Filter(expr), nil
}

// Get retrieves a model by ID
func (m *Manager[T]) Get(ctx context.Context, id int64) (*T, error) {
	if m.db == nil {
		return nil, errors.NewNotImplementedError("Manager.Get() - database connection not set")
	}

	fa, err := m.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	idField := FieldFor[T, int64](fa, "ID")
	if idField.Path() == "" {
		// Try lowercase
		idField = FieldFor[T, int64](fa, "id")
	}

	qs, err := m.Filter(idField.Eq(id))
	if err != nil {
		return nil, err
	}

	return qs.Get(ctx)
}

// All returns all model instances
func (m *Manager[T]) All(ctx context.Context) ([]*T, error) {
	qs, err := NewQuerySet[T](m.tableName)
	if err != nil {
		return nil, err
	}

	if m.db != nil {
		qs = qs.SetDB(m.db).(QuerySet[T])
	}

	return qs.All(ctx)
}

// Create creates a new model instance
func (m *Manager[T]) Create(ctx context.Context, instance *T) error {
	if m.db == nil {
		return errors.NewNotImplementedError("Manager.Create() - database connection not set")
	}

	// Run BeforeCreate hooks
	if hookable, ok := any(instance).(interface{ BeforeCreate(context.Context) error }); ok {
		if err := hookable.BeforeCreate(ctx); err != nil {
			return fmt.Errorf("BeforeCreate hook failed: %w", err)
		}
	}

	// Run BeforeSave hook
	if hookable, ok := any(instance).(interface{ BeforeSave(context.Context) error }); ok {
		if err := hookable.BeforeSave(ctx); err != nil {
			return fmt.Errorf("BeforeSave hook failed: %w", err)
		}
	}

	// Validate instance
	if validatable, ok := any(instance).(interface{ Clean() error }); ok {
		if err := validatable.Clean(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
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

	// Set the ID on the instance type-safely
	if modelWithID, ok := any(instance).(ModelWithID); ok {
		modelWithID.SetID(id)
	} else {
		// Fallback to reflection
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
	if hookable, ok := any(instance).(interface{ AfterCreate(context.Context) error }); ok {
		if err := hookable.AfterCreate(ctx); err != nil {
			return fmt.Errorf("AfterCreate hook failed: %w", err)
		}
	}

	// Run AfterSave hook
	if hookable, ok := any(instance).(interface{ AfterSave(context.Context) error }); ok {
		if err := hookable.AfterSave(ctx); err != nil {
			return fmt.Errorf("AfterSave hook failed: %w", err)
		}
	}

	return nil
}

// BulkCreate creates multiple model instances efficiently using a single INSERT statement
func (m *Manager[T]) BulkCreate(ctx context.Context, instances []*T) error {
	if m.db == nil {
		return errors.NewNotImplementedError("Manager.BulkCreate() - database connection not set")
	}

	if len(instances) == 0 {
		return nil
	}

	// Run BeforeCreate hooks for all instances
	for _, instance := range instances {
		if hookable, ok := any(instance).(interface{ BeforeCreate(context.Context) error }); ok {
			if err := hookable.BeforeCreate(ctx); err != nil {
				return fmt.Errorf("BeforeCreate hook failed: %w", err)
			}
		}

		// Run BeforeSave hook
		if hookable, ok := any(instance).(interface{ BeforeSave(context.Context) error }); ok {
			if err := hookable.BeforeSave(ctx); err != nil {
				return fmt.Errorf("BeforeSave hook failed: %w", err)
			}
		}

		// Validate instance
		if validatable, ok := any(instance).(interface{ Clean() error }); ok {
			if err := validatable.Clean(); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}
		}
	}

	// Convert []*T to []interface{} for BuildBulkInsertSQL
	instancesInterface := make([]interface{}, len(instances))
	for i, instance := range instances {
		instancesInterface[i] = instance
	}

	// Build bulk INSERT SQL
	sql, args, _, err := BuildBulkInsertSQL(instancesInterface, m.tableName)
	if err != nil {
		return fmt.Errorf("failed to build bulk insert SQL: %w", err)
	}

	// Execute bulk INSERT and get generated IDs
	ids, err := ExecuteBulkInsert(ctx, m.db, sql, args)
	if err != nil {
		return err
	}

	// Set IDs on instances
	for i, instance := range instances {
		if i < len(ids) {
			if modelWithID, ok := any(instance).(ModelWithID); ok {
				modelWithID.SetID(ids[i])
			} else {
				// Fallback to reflection
				instanceValue := reflect.ValueOf(instance).Elem()
				idField := instanceValue.FieldByName("ID")
				if !idField.IsValid() {
					idField = instanceValue.FieldByName("id")
				}
				if idField.IsValid() && idField.CanSet() {
					idField.SetInt(ids[i])
				}
			}
		}
	}

	// Run AfterCreate hooks for all instances
	for _, instance := range instances {
		if hookable, ok := any(instance).(interface{ AfterCreate(context.Context) error }); ok {
			if err := hookable.AfterCreate(ctx); err != nil {
				return fmt.Errorf("AfterCreate hook failed: %w", err)
			}
		}

		// Run AfterSave hook
		if hookable, ok := any(instance).(interface{ AfterSave(context.Context) error }); ok {
			if err := hookable.AfterSave(ctx); err != nil {
				return fmt.Errorf("AfterSave hook failed: %w", err)
			}
		}
	}

	return nil
}

// Update updates an existing model instance
func (m *Manager[T]) Update(ctx context.Context, instance *T) error {
	if m.db == nil {
		return errors.NewNotImplementedError("Manager.Update() - database connection not set")
	}

	// Get ID value type-safely
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
	if hookable, ok := any(instance).(interface{ BeforeUpdate(context.Context) error }); ok {
		if err := hookable.BeforeUpdate(ctx); err != nil {
			return fmt.Errorf("BeforeUpdate hook failed: %w", err)
		}
	}

	// Run BeforeSave hook
	if hookable, ok := any(instance).(interface{ BeforeSave(context.Context) error }); ok {
		if err := hookable.BeforeSave(ctx); err != nil {
			return fmt.Errorf("BeforeSave hook failed: %w", err)
		}
	}

	// Validate instance
	if validatable, ok := any(instance).(interface{ Clean() error }); ok {
		if err := validatable.Clean(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
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
	if hookable, ok := any(instance).(interface{ AfterUpdate(context.Context) error }); ok {
		if err := hookable.AfterUpdate(ctx); err != nil {
			return fmt.Errorf("AfterUpdate hook failed: %w", err)
		}
	}

	// Run AfterSave hook
	if hookable, ok := any(instance).(interface{ AfterSave(context.Context) error }); ok {
		if err := hookable.AfterSave(ctx); err != nil {
			return fmt.Errorf("AfterSave hook failed: %w", err)
		}
	}

	return nil
}

// Save saves a model (create or update)
func (m *Manager[T]) Save(ctx context.Context, instance *T) error {
	// Use type assertion if available
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

	// Get ID value type-safely
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
	if hookable, ok := any(instance).(interface{ BeforeDelete(context.Context) error }); ok {
		if err := hookable.BeforeDelete(ctx); err != nil {
			return fmt.Errorf("BeforeDelete hook failed: %w", err)
		}
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
	if hookable, ok := any(instance).(interface{ AfterDelete(context.Context) error }); ok {
		if err := hookable.AfterDelete(ctx); err != nil {
			return fmt.Errorf("AfterDelete hook failed: %w", err)
		}
	}

	return nil
}

// UpdateFields updates specific fields type-safely
func (m *Manager[T]) UpdateFields(ctx context.Context, id int64, updates UpdateMap) error {
	// Validate all fields exist
	for fieldName := range updates {
		fieldInfo := m.schema.GetField(fieldName)
		if fieldInfo == nil {
			return fmt.Errorf("field %s not found", fieldName)
		}
	}

	// Get QuerySet and use UpdateBuilder
	fa, err := m.GetFieldAccessor()
	if err != nil {
		return err
	}

	idField := FieldFor[T, int64](fa, "ID")
	if idField.Path() == "" {
		idField = FieldFor[T, int64](fa, "id")
	}

	qs, err := m.Filter(idField.Eq(id))
	if err != nil {
		return err
	}

	ub, err := qs.UpdateBuilder()
	if err != nil {
		return err
	}

	// Apply updates
	for fieldName, value := range updates {
		// Type checking happens in Set
		ub.updates[fieldName] = value
	}

	_, err = ub.Execute(ctx)
	return err
}
