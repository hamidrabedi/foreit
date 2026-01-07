package orm

import (
	"reflect"
	"sync"

	"github.com/forgego/forge/db"
)

type managerRegistry struct {
	mu    sync.RWMutex
	items map[reflect.Type]reflect.Value
}

var globalManagerRegistry = &managerRegistry{
	items: make(map[reflect.Type]reflect.Value),
}

var (
	defaultDBMu sync.RWMutex
	defaultDB   *db.DB
)

// RegisterManagerFor registers a manager for type T.
func RegisterManagerFor[T any](manager *Manager[T]) {
	if manager == nil {
		return
	}

	var zero T
	typ := reflect.TypeOf(zero)
	if typ == nil {
		return
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	globalManagerRegistry.mu.Lock()
	globalManagerRegistry.items[typ] = reflect.ValueOf(manager)
	globalManagerRegistry.mu.Unlock()
}

// SetDefaultDB sets the default database connection for managers without an explicit DB.
func SetDefaultDB(database *db.DB) {
	defaultDBMu.Lock()
	defaultDB = database
	defaultDBMu.Unlock()
}

// GetDefaultDB returns the default database connection.
func GetDefaultDB() *db.DB {
	defaultDBMu.RLock()
	defer defaultDBMu.RUnlock()
	return defaultDB
}

// GetManagerForType returns the registered manager for a model type.
func GetManagerForType(modelType reflect.Type) reflect.Value {
	if modelType == nil {
		return reflect.Value{}
	}
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	globalManagerRegistry.mu.RLock()
	manager, ok := globalManagerRegistry.items[modelType]
	globalManagerRegistry.mu.RUnlock()
	if !ok {
		return reflect.Value{}
	}
	return manager
}

// SetDBForAllManagers sets the database connection for all registered managers.
func SetDBForAllManagers(database *db.DB) {
	if database == nil {
		return
	}

	SetDefaultDB(database)

	globalManagerRegistry.mu.RLock()
	defer globalManagerRegistry.mu.RUnlock()

	for _, manager := range globalManagerRegistry.items {
		if !manager.IsValid() {
			continue
		}
		if setter, ok := manager.Interface().(interface{ SetDB(*db.DB) }); ok {
			setter.SetDB(database)
		}
	}
}



