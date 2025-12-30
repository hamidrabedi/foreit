package http

import (
	"github.com/forgego/forge/admin"
)

// RegisterAdminForHTTP registers an admin instance for HTTP handlers
// This should be called after registering with admin.Register()
func RegisterAdminForHTTP[T any](admin *admin.Admin[T]) {
	RegisterAdmin(admin)
}
