// Package forge provides admin interface functionality
package forge

import (
	pkgAdmin "github.com/forgego/forge/pkg/admin"
)

// RegisterAdminRoutes registers admin routes on the router
// sessionManager and database are optional - if provided, login/logout routes will be registered
func RegisterAdminRoutes(router *Router, path string, sessionManager interface{}, database interface{}) {
	pkgAdmin.RegisterAdminRoutes(router, path, sessionManager, database)
}

// RegisterModel registers a model for admin auto-generation
func RegisterModel(model interface{}) error {
	return pkgAdmin.RegisterModel(model)
}

// RegisterModelWithOptions registers a model with custom admin options
func RegisterModelWithOptions(model interface{}, options ...AdminOption) error {
	return pkgAdmin.RegisterModelWithOptions(model, options...)
}

// AdminOption is a function that configures an admin model
type AdminOption = pkgAdmin.AdminOption

// WithListDisplay sets the list display fields
func WithListDisplay(fields ...interface{}) AdminOption {
	return pkgAdmin.WithListDisplay(fields...)
}

// WithListFilter sets the list filter fields
func WithListFilter(fields ...interface{}) AdminOption {
	return pkgAdmin.WithListFilter(fields...)
}

// WithSearchFields sets the search fields
func WithSearchFields(fields ...interface{}) AdminOption {
	return pkgAdmin.WithSearchFields(fields...)
}

// WithReadOnlyFields sets the read-only fields
func WithReadOnlyFields(fields ...interface{}) AdminOption {
	return pkgAdmin.WithReadOnlyFields(fields...)
}
