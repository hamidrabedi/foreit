package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/forgego/forge/pkg/models"
)

// ViewSet is the core interface for DRF-style viewsets
type ViewSet[T any] interface {
	List(c *fiber.Ctx) error
	Retrieve(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Destroy(c *fiber.Ctx) error
}


// FilterBackend applies filters to a QuerySet based on request
type FilterBackend[T any] interface {
	Apply(c *fiber.Ctx, qs models.QuerySet[T]) models.QuerySet[T]
}

// PaginationBackend handles pagination
type PaginationBackend[T any] interface {
	Paginate(ctx context.Context, qs models.QuerySet[T], c *fiber.Ctx) ([]*T, int64, error)
}

// Permission checks if a request is allowed
type Permission interface {
	HasPermission(c *fiber.Ctx) bool
}


// RegisterViewSet registers a ViewSet with Fiber router
// Creates routes:
//   GET    /resource/      -> List
//   POST   /resource/      -> Create
//   GET    /resource/:id  -> Retrieve
//   PUT    /resource/:id  -> Update
//   PATCH  /resource/:id  -> Update
//   DELETE /resource/:id  -> Destroy
func RegisterViewSet[T any](app *fiber.App, path string, vs ViewSet[T]) {
	// Ensure path doesn't end with /
	if path != "" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	
	// List and Create
	app.Get(path+"/", vs.List)
	app.Post(path+"/", vs.Create)
	
	// Retrieve, Update, Destroy
	app.Get(path+"/:id", vs.Retrieve)
	app.Put(path+"/:id", vs.Update)
	app.Patch(path+"/:id", vs.Update)
	app.Delete(path+"/:id", vs.Destroy)
}

