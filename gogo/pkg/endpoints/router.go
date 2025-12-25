package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
)

type Router struct {
	app       *fiber.App
	resources map[string]ResourceHandler
	basePath  string
}

type ResourceHandler struct {
	Resource Resource[interface{}, interface{}]
	Name     string
}

func NewRouter(app *fiber.App, basePath string) *Router {
	if basePath == "" {
		basePath = "/api"
	}
	return &Router{
		app:       app,
		resources: make(map[string]ResourceHandler),
		basePath:  basePath,
	}
}

func (r *Router) RegisterResource[T any, Q any](name string, resource Resource[T, Q]) {
	wrapper := &resourceWrapper[T, Q]{resource: resource}
	
	r.resources[name] = ResourceHandler{
		Resource: wrapper,
		Name:     name,
	}
	
	r.installRoutes(name, wrapper)
}

func (r *Router) installRoutes[T any, Q any](name string, resource Resource[T, Q]) {
	prefix := r.basePath + "/" + name
	
	r.app.Get(prefix, func(c *fiber.Ctx) error {
		ctx := &Context{Ctx: c, Request: c.Request()}
		results, err := resource.Index(ctx)
		if err != nil {
			return handleError(c, err)
		}
		return c.JSON(fiber.Map{"data": results})
	})
	
	r.app.Get(prefix+"/:id", func(c *fiber.Ctx) error {
		ctx := &Context{Ctx: c, Request: c.Request()}
		result, err := resource.Show(ctx)
		if err != nil {
			return handleError(c, err)
		}
		return c.JSON(fiber.Map{"data": result})
	})
	
	r.app.Post(prefix, func(c *fiber.Ctx) error {
		ctx := &Context{Ctx: c, Request: c.Request()}
		result, err := resource.Create(ctx)
		if err != nil {
			return handleError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
	})
	
	r.app.Put(prefix+"/:id", func(c *fiber.Ctx) error {
		ctx := &Context{Ctx: c, Request: c.Request()}
		result, err := resource.Update(ctx)
		if err != nil {
			return handleError(c, err)
		}
		return c.JSON(fiber.Map{"data": result})
	})
	
	r.app.Delete(prefix+"/:id", func(c *fiber.Ctx) error {
		ctx := &Context{Ctx: c, Request: c.Request()}
		if err := resource.Destroy(ctx); err != nil {
			return handleError(c, err)
		}
		return c.SendStatus(fiber.StatusNoContent)
	})
}

type resourceWrapper[T any, Q any] struct {
	resource Resource[T, Q]
}

func (w *resourceWrapper[T, Q]) Index(ctx *Context) ([]interface{}, error) {
	results, err := w.resource.Index(ctx)
	if err != nil {
		return nil, err
	}
	return toInterfaceSlice(results), nil
}

func (w *resourceWrapper[T, Q]) Show(ctx *Context) (interface{}, error) {
	return w.resource.Show(ctx)
}

func (w *resourceWrapper[T, Q]) Create(ctx *Context) (interface{}, error) {
	return w.resource.Create(ctx)
}

func (w *resourceWrapper[T, Q]) Update(ctx *Context) (interface{}, error) {
	return w.resource.Update(ctx)
}

func (w *resourceWrapper[T, Q]) Destroy(ctx *Context) error {
	return w.resource.Destroy(ctx)
}

func toInterfaceSlice(slice interface{}) []interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil
	}
	
	result := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}
	return result
}

func handleError(c *fiber.Ctx, err error) error {
	return HandleError(c, err)
}

func (r *Router) GetResource(name string) (ResourceHandler, bool) {
	resource, ok := r.resources[name]
	return resource, ok
}

func (r *Router) ListResources() []string {
	names := make([]string, 0, len(r.resources))
	for name := range r.resources {
		names = append(names, name)
	}
	return names
}

func (r *Router) Group(prefix string, handlers ...fiber.Handler) *fiber.Router {
	group := r.app.Group(r.basePath + prefix)
	for _, handler := range handlers {
		group.Use(handler)
	}
	return &group
}

func resourceNameToPath(name string) string {
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('-')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

