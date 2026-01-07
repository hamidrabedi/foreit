package schema

import (
	"fmt"
	"sync"
)

// RelationFactory creates a new relation instance for the given name/target.
type RelationFactory func(name, to string) Relation

var (
	relationRegistry = make(map[string]RelationFactory)
	relationMutex    sync.RWMutex
)

// RegisterRelationType registers a custom relation type.
func RegisterRelationType(typeName string, relationType RelationType, factory RelationFactory) error {
	relationMutex.Lock()
	defer relationMutex.Unlock()

	if typeName == "" {
		return fmt.Errorf("relation type name cannot be empty")
	}

	if _, exists := relationRegistry[typeName]; exists {
		return fmt.Errorf("relation type %q is already registered", typeName)
	}

	if factory == nil {
		factory = func(name, to string) Relation {
			return newRelation(name, to, relationType)
		}
	}

	relationRegistry[typeName] = factory
	return nil
}

// UnregisterRelationType removes a registered relation type.
func UnregisterRelationType(typeName string) {
	relationMutex.Lock()
	defer relationMutex.Unlock()
	delete(relationRegistry, typeName)
}

// NewRelation creates a new relation instance from a registered relation type.
func NewRelation(typeName, name, to string) (Relation, error) {
	relationMutex.RLock()
	factory, exists := relationRegistry[typeName]
	relationMutex.RUnlock()
	if !exists {
		return Relation{}, fmt.Errorf("relation type %q is not registered", typeName)
	}

	return factory(name, to), nil
}

// ListRegisteredRelationTypes returns all registered custom relation types.
func ListRegisteredRelationTypes() []string {
	relationMutex.RLock()
	defer relationMutex.RUnlock()

	types := make([]string, 0, len(relationRegistry))
	for typeName := range relationRegistry {
		types = append(types, typeName)
	}
	return types
}
