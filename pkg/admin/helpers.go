package admin

import (
	"fmt"
	"net/http"
)

// checkPermission checks if the current user has permission for an action
func checkPermission(r *http.Request, action string, modelName string) error {
	user, ok := GetUser(r)
	if !ok {
		return fmt.Errorf("user not authenticated")
	}

	// Build permission string (e.g., "admin.add_user", "admin.change_user")
	perm := fmt.Sprintf("admin.%s_%s", action, modelName)
	
	if !HasPerm(user, perm) {
		return fmt.Errorf("permission denied: %s", perm)
	}

	return nil
}

// requirePermission is a middleware that checks permission before allowing access
func requirePermission(action string, modelName string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := checkPermission(r, action, modelName); err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			next(w, r)
		}
	}
}

// canAdd checks if user can add instances of a model
func canAdd(r *http.Request, modelName string) bool {
	user, ok := GetUser(r)
	if !ok {
		return false
	}
	perm := fmt.Sprintf("admin.add_%s", modelName)
	return HasPerm(user, perm)
}

// canChange checks if user can change instances of a model
func canChange(r *http.Request, modelName string) bool {
	user, ok := GetUser(r)
	if !ok {
		return false
	}
	perm := fmt.Sprintf("admin.change_%s", modelName)
	return HasPerm(user, perm)
}

// canDelete checks if user can delete instances of a model
func canDelete(r *http.Request, modelName string) bool {
	user, ok := GetUser(r)
	if !ok {
		return false
	}
	perm := fmt.Sprintf("admin.delete_%s", modelName)
	return HasPerm(user, perm)
}

// canView checks if user can view instances of a model
func canView(r *http.Request, modelName string) bool {
	user, ok := GetUser(r)
	if !ok {
		return false
	}
	perm := fmt.Sprintf("admin.view_%s", modelName)
	return HasPerm(user, perm)
}


