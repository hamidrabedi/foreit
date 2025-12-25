package auth

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

// Policy defines authorization policies for a model
type Policy[T any] interface {
	// CanView checks if user can view an object
	CanView(user interface{}, obj T) bool
	
	// CanEdit checks if user can edit an object
	CanEdit(user interface{}, obj T) bool
	
	// CanDelete checks if user can delete an object
	CanDelete(user interface{}, obj T) bool
	
	// CanCreate checks if user can create objects
	CanCreate(user interface{}) bool
}

// PolicyRegistry stores policies for different models
type PolicyRegistry struct {
	policies map[string]Policy[interface{}]
}

var globalRegistry = &PolicyRegistry{
	policies: make(map[string]Policy[interface{}]),
}

// Register registers a policy for a model type
func Register[T any](policy Policy[T]) {
	var zero T
	typeName := getTypeName(zero)
	
	// Wrap typed policy as interface{}
	wrapper := &policyWrapper[T]{policy: policy}
	globalRegistry.policies[typeName] = wrapper
}

// Get retrieves a policy for a type
func Get[T any]() (Policy[T], bool) {
	var zero T
	typeName := getTypeName(zero)
	
	policy, ok := globalRegistry.policies[typeName]
	if !ok {
		return nil, false
	}
	
	// Unwrap to typed policy
	if wrapper, ok := policy.(*policyWrapper[T]); ok {
		return wrapper.policy, true
	}
	
	return nil, false
}

// policyWrapper wraps a typed policy as interface{}
type policyWrapper[T any] struct {
	policy Policy[T]
}

func (w *policyWrapper[T]) CanView(user interface{}, obj interface{}) bool {
	if typedObj, ok := obj.(T); ok {
		return w.policy.CanView(user, typedObj)
	}
	return false
}

func (w *policyWrapper[T]) CanEdit(user interface{}, obj interface{}) bool {
	if typedObj, ok := obj.(T); ok {
		return w.policy.CanEdit(user, typedObj)
	}
	return false
}

func (w *policyWrapper[T]) CanDelete(user interface{}, obj interface{}) bool {
	if typedObj, ok := obj.(T); ok {
		return w.policy.CanDelete(user, typedObj)
	}
	return false
}

func (w *policyWrapper[T]) CanCreate(user interface{}) bool {
	return w.policy.CanCreate(user)
}

// Can checks if a user can perform an action on an object
func Can[T any](user interface{}, action string, obj T) (bool, error) {
	policy, ok := Get[T]()
	if !ok {
		return false, errors.New("no policy registered for type")
	}
	
	switch action {
	case "view":
		return policy.CanView(user, obj), nil
	case "edit", "update":
		return policy.CanEdit(user, obj), nil
	case "delete":
		return policy.CanDelete(user, obj), nil
	case "create":
		return policy.CanCreate(user), nil
	default:
		return false, errors.New("unknown action")
	}
}

// Require checks if user can perform action, returns error if not
func Require[T any](c *fiber.Ctx, action string, obj T) error {
	user := c.Locals("user")
	if user == nil {
		return errors.New("user not authenticated")
	}
	
	allowed, err := Can[T](user, action, obj)
	if err != nil {
		return err
	}
	
	if !allowed {
		return errors.New("forbidden")
	}
	
	return nil
}

// getTypeName gets the type name of a value
func getTypeName(v interface{}) string {
	// Simple implementation - in production might use reflection or generics
	return "unknown"
}

