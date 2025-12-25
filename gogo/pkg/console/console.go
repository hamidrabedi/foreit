package console

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

// Console defines the interface for admin console
type Console[T any] interface {
	// ListDisplay returns fields to display in list view
	ListDisplay() []string
	
	// SearchFields returns fields to search in
	SearchFields() []string
	
	// ListFilters returns filters for list view
	ListFilters() []Filter
	
	// Actions returns custom actions
	Actions() []Action
}

// ModelConsole provides a base console implementation
type ModelConsole[T any] struct {
	options *Options
}

// Options configures a console
type Options struct {
	ListDisplay  []string
	SearchFields []string
	Filters      []Filter
	Actions      []Action
	ReadOnly     bool
}

// NewModelConsole creates a new model console
func NewModelConsole[T any](options *Options) *ModelConsole[T] {
	if options == nil {
		options = &Options{}
	}
	return &ModelConsole[T]{
		options: options,
	}
}

// ListDisplay returns fields to display
func (c *ModelConsole[T]) ListDisplay() []string {
	if c.options != nil && len(c.options.ListDisplay) > 0 {
		return c.options.ListDisplay
	}
	return []string{} // Will be auto-generated from Ent schema
}

// SearchFields returns searchable fields
func (c *ModelConsole[T]) SearchFields() []string {
	if c.options != nil && len(c.options.SearchFields) > 0 {
		return c.options.SearchFields
	}
	return []string{}
}

// ListFilters returns filters
func (c *ModelConsole[T]) ListFilters() []Filter {
	if c.options != nil && len(c.options.Filters) > 0 {
		return c.options.Filters
	}
	return []Filter{}
}

// Actions returns custom actions
func (c *ModelConsole[T]) Actions() []Action {
	if c.options != nil && len(c.options.Actions) > 0 {
		return c.options.Actions
	}
	return []Filter{}
}

// Filter represents a filter in the console
type Filter struct {
	Field    string
	Type     FilterType
	Choices  []Choice
	Label    string
}

// FilterType represents the type of filter
type FilterType string

const (
	FilterTypeText     FilterType = "text"
	FilterTypeNumber   FilterType = "number"
	FilterTypeDate     FilterType = "date"
	FilterTypeDateTime FilterType = "datetime"
	FilterTypeBoolean  FilterType = "boolean"
	FilterTypeChoice   FilterType = "choice"
)

// Choice represents a filter choice
type Choice struct {
	Value string
	Label string
}

// Action represents a custom action
type Action struct {
	Name        string
	Label       string
	Handler     func(context.Context, []interface{}) error
	Description string
	Icon        string
}

// Registry stores registered consoles
type Registry struct {
	consoles map[string]ConsoleInfo
}

type ConsoleInfo struct {
	Console Console[interface{}]
	Name    string
}

var globalRegistry = &Registry{
	consoles: make(map[string]ConsoleInfo),
}

// Register registers a console for a model
func Register[T any](console Console[T]) {
	var zero T
	name := getTypeName(zero)
	
	// Wrap typed console as interface{}
	wrapper := &consoleWrapper[T]{console: console}
	globalRegistry.consoles[name] = ConsoleInfo{
		Console: wrapper,
		Name:    name,
	}
}

// consoleWrapper wraps a typed console as interface{}
type consoleWrapper[T any] struct {
	console Console[T]
}

func (w *consoleWrapper[T]) ListDisplay() []string {
	return w.console.ListDisplay()
}

func (w *consoleWrapper[T]) SearchFields() []string {
	return w.console.SearchFields()
}

func (w *consoleWrapper[T]) ListFilters() []Filter {
	return w.console.ListFilters()
}

func (w *consoleWrapper[T]) Actions() []Action {
	return w.console.Actions()
}

// getTypeName gets the type name (simplified)
func getTypeName(v interface{}) string {
	// In production, use reflection or generics properly
	return "unknown"
}

// InstallRoutes installs console routes on a Fiber app
func InstallRoutes(app *fiber.App, basePath string) {
	if basePath == "" {
		basePath = "/console"
	}
	
	// List all models
	app.Get(basePath, listModels)
	
	// Model-specific routes
	app.Get(basePath+"/:model", listModel)
	app.Get(basePath+"/:model/:id", showModel)
	app.Get(basePath+"/:model/new", newModel)
	app.Post(basePath+"/:model", createModel)
	app.Get(basePath+"/:model/:id/edit", editModel)
	app.Put(basePath+"/:model/:id", updateModel)
	app.Delete(basePath+"/:model/:id", deleteModel)
}

// listModels handles GET /console
func listModels(c *fiber.Ctx) error {
	models := make([]fiber.Map, 0)
	for name, info := range globalRegistry.consoles {
		models = append(models, fiber.Map{
			"name":         name,
			"list_display": info.Console.ListDisplay(),
			"url":          "/console/" + name,
		})
	}
	return c.JSON(fiber.Map{
		"models": models,
		"count":  len(models),
	})
}

