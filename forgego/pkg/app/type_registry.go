package app

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/forgego/forge/pkg/models"
)

type ManagerFactory func(*models.DB) interface{}
type SerializerFactory func() interface{}
type ServiceFactory func(interface{}, interface{}) interface{}

var (
	typeRegistry     = make(map[reflect.Type]*TypeRegistryEntry)
	registryMutex    sync.RWMutex
)

type TypeRegistryEntry struct {
	ManagerFactory   ManagerFactory
	SerializerFactory SerializerFactory
	ServiceFactory   ServiceFactory
}

func RegisterType[T any](managerFactory ManagerFactory, serializerFactory SerializerFactory, serviceFactory ServiceFactory) {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	registryMutex.Lock()
	defer registryMutex.Unlock()

	typeRegistry[t] = &TypeRegistryEntry{
		ManagerFactory:   managerFactory,
		SerializerFactory: serializerFactory,
		ServiceFactory:   serviceFactory,
	}
}

func getTypeEntry(baseType reflect.Type) (*TypeRegistryEntry, error) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	entry, ok := typeRegistry[baseType]
	if !ok {
		return nil, fmt.Errorf("type %s not registered in type registry. Use app generate framework-model to generate the model", baseType.Name())
	}
	return entry, nil
}

