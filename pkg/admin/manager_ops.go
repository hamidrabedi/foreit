package admin

import (
	"context"
	"fmt"
	"reflect"
)

// ManagerOps provides helper functions for manager operations using reflection
type ManagerOps struct{}

// NewManagerOps creates a new ManagerOps instance
func NewManagerOps() *ManagerOps {
	return &ManagerOps{}
}

// GetInstance retrieves an instance by ID using the manager
func (m *ManagerOps) GetInstance(ctx context.Context, manager interface{}, id int64) (interface{}, error) {
	managerValue := reflect.ValueOf(manager)
	getMethod := managerValue.MethodByName("Get")
	if !getMethod.IsValid() {
		return nil, fmt.Errorf("manager does not have Get method")
	}

	results := getMethod.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(id),
	})

	if len(results) < 2 {
		return nil, fmt.Errorf("invalid Get method signature")
	}

	if !results[1].IsNil() {
		if err, ok := results[1].Interface().(error); ok && err != nil {
			return nil, fmt.Errorf("failed to get instance: %w", err)
		}
	}

	return results[0].Interface(), nil
}

// CreateInstance creates a new instance using the manager
func (m *ManagerOps) CreateInstance(ctx context.Context, manager interface{}, instance interface{}) error {
	managerValue := reflect.ValueOf(manager)
	createMethod := managerValue.MethodByName("Create")
	if !createMethod.IsValid() {
		return fmt.Errorf("manager does not have Create method")
	}

	results := createMethod.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(instance),
	})

	if len(results) > 0 && !results[0].IsNil() {
		if err, ok := results[0].Interface().(error); ok && err != nil {
			return fmt.Errorf("failed to create instance: %w", err)
		}
	}

	return nil
}

// UpdateInstance updates an instance using the manager
func (m *ManagerOps) UpdateInstance(ctx context.Context, manager interface{}, instance interface{}) error {
	managerValue := reflect.ValueOf(manager)
	updateMethod := managerValue.MethodByName("Update")
	if !updateMethod.IsValid() {
		return fmt.Errorf("manager does not have Update method")
	}

	results := updateMethod.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(instance),
	})

	if len(results) > 0 && !results[0].IsNil() {
		if err, ok := results[0].Interface().(error); ok && err != nil {
			return fmt.Errorf("failed to update instance: %w", err)
		}
	}

	return nil
}

// DeleteInstance deletes an instance using the manager
func (m *ManagerOps) DeleteInstance(ctx context.Context, manager interface{}, instance interface{}) error {
	managerValue := reflect.ValueOf(manager)
	deleteMethod := managerValue.MethodByName("Delete")
	if !deleteMethod.IsValid() {
		return fmt.Errorf("manager does not have Delete method")
	}

	results := deleteMethod.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(instance),
	})

	if len(results) > 0 && !results[0].IsNil() {
		if err, ok := results[0].Interface().(error); ok && err != nil {
			return fmt.Errorf("failed to delete instance: %w", err)
		}
	}

	return nil
}

// GetAllInstances retrieves all instances using the manager
func (m *ManagerOps) GetAllInstances(ctx context.Context, manager interface{}) ([]interface{}, error) {
	managerValue := reflect.ValueOf(manager)
	allMethod := managerValue.MethodByName("All")
	if !allMethod.IsValid() {
		return nil, fmt.Errorf("manager does not have All method")
	}

	results := allMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
	if len(results) < 2 {
		return nil, fmt.Errorf("invalid All method signature")
	}

	if !results[1].IsNil() {
		if err, ok := results[1].Interface().(error); ok && err != nil {
			return nil, fmt.Errorf("failed to get all instances: %w", err)
		}
	}

	// Convert slice to []interface{}
	objectsValue := reflect.ValueOf(results[0].Interface())
	if objectsValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("All method did not return a slice")
	}

	objectsInterface := make([]interface{}, objectsValue.Len())
	for i := 0; i < objectsValue.Len(); i++ {
		objectsInterface[i] = objectsValue.Index(i).Interface()
	}

	return objectsInterface, nil
}

// GetCount gets the count of instances using the manager
func (m *ManagerOps) GetCount(ctx context.Context, manager interface{}) (int64, error) {
	// Try to get count via All() method
	instances, err := m.GetAllInstances(ctx, manager)
	if err != nil {
		return 0, err
	}
	return int64(len(instances)), nil
}

// CreateNewInstance creates a new instance of the model type
func CreateNewInstance(model interface{}) interface{} {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	instanceValue := reflect.New(modelType)
	return instanceValue.Interface()
}

// getIDFromInstance extracts ID from an instance
func getIDFromInstance(instance interface{}) int64 {
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}
	idField := instanceValue.FieldByName("ID")
	if !idField.IsValid() {
		idField = instanceValue.FieldByName("id")
	}
	if idField.IsValid() && idField.CanInterface() {
		if idVal, ok := idField.Interface().(int64); ok {
			return idVal
		}
	}
	return 0
}

