package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/forgego/forge/pkg/admin"
	"github.com/forgego/forge/pkg/api"
	"github.com/forgego/forge/pkg/i18n"
	"github.com/forgego/forge/pkg/middleware"
	"github.com/forgego/forge/pkg/models"
	"github.com/forgego/forge/pkg/rest"
	"github.com/forgego/forge/pkg/settings"
	"github.com/forgego/forge/pkg/workers"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	fiber                  *fiber.App
	settings               *AppSettings
	modelsDB               *models.DB
	router              *rest.Router
	serviceIntegration *admin.ServiceIntegration
	workersStarted     bool
}

// New creates a new App instance and automatically initializes enabled framework features
// If settings is nil, it will attempt to load from config files or environment variables
// Returns error if any required initialization fails
func New(appSettings *AppSettings) (*App, error) {
	if appSettings == nil {
		// Try to load from config files or environment
		loaded, err := settings.Load[AppSettings]()
		if err != nil {
			// If loading fails, use defaults
			appSettings = DefaultAppSettings()
		} else {
			appSettings = loaded
		}
	}
	
	app := &App{
		fiber: fiber.New(fiber.Config{
			ErrorHandler: rest.HandleError,
		}),
		settings: appSettings,
	}
	
	// Automatically initialize framework features based on settings
	if err := app.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize app: %w", err)
	}
	
	return app, nil
}

// initialize sets up all enabled framework features
func (a *App) initialize() error {
	// 1. Set up middleware (always enabled for framework)
	a.setupMiddleware()
	
	// 2. Connect to database if URL is provided
	if a.settings.Database.URL != "" {
		if err := a.connectDB(); err != nil {
			return fmt.Errorf("database connection failed: %w", err)
		}
	}
	
	// 3. Set up API router if path is configured
	if a.settings.API.Path != "" {
		a.setupAPIRouter()
	}
	
	// 4. Set up i18n if enabled
	if a.settings.I18n.Enable {
		if err := a.setupI18n(); err != nil {
			return fmt.Errorf("i18n setup failed: %w", err)
		}
	}
	
	// 5. Set up static files if enabled
	if a.settings.Static.Enable {
		a.setupStatic()
	}
	
	// 6. Set up workers if enabled
	if a.settings.Workers.Enable {
		if err := a.setupWorkers(); err != nil {
			return fmt.Errorf("workers setup failed: %w", err)
		}
	}
	
	// 7. Set up admin if enabled (requires database)
	if a.settings.Admin.Enable {
		if err := a.setupAdmin(); err != nil {
			return fmt.Errorf("admin setup failed: %w", err)
		}
	}
	
	return nil
}

func (a *App) setupMiddleware() {
	a.fiber.Use(middleware.Logging())
	a.fiber.Use(middleware.Recovery())
	a.fiber.Use(middleware.CORS())
	a.fiber.Use(middleware.RequestID())
	a.fiber.Use(middleware.SecurityHeaders())
	
	if a.settings.Security.Debug {
		a.fiber.Use(middleware.RateLimit(1000, time.Minute))
	} else {
		a.fiber.Use(middleware.RateLimit(100, time.Minute))
	}
}

func (a *App) connectDB() error {
	driver := a.settings.Database.Driver
	if driver == "" {
		driver = "postgres"
	}
	
	modelsDB, err := models.NewDB(driver, a.settings.Database.URL)
	if err != nil {
		return err
	}
	a.modelsDB = modelsDB
	return nil
}

func (a *App) setupAPIRouter() {
	if a.router == nil {
		a.router = rest.NewRouter(a.fiber, a.settings.API.Path)
	}
}

func (a *App) setupI18n() error {
	if a.settings.I18n.LocalesPath == "" {
		return fmt.Errorf("i18n.locales_path is required when i18n is enabled")
	}
	
	loader := i18n.NewJSONLoader(a.settings.I18n.LocalesPath)
		bundles, err := loader.Load()
		if err != nil {
		return fmt.Errorf("failed to load i18n bundles: %w", err)
	}
	
	manager := i18n.NewManager(a.settings.I18n.DefaultLocale)
			for locale, bundle := range bundles {
				manager.RegisterBundle(locale, bundle)
			}
			i18n.SetDefaultManager(manager)
	return nil
	}
	
func (a *App) setupStatic() {
	if a.settings.Static.Root != "" {
		a.fiber.Static(a.settings.Static.Path, a.settings.Static.Root, fiber.Static{
			Compress: true,
		})
	}
}

func (a *App) setupWorkers() error {
	queue, err := workers.NewAsynqQueue(
		a.settings.Workers.RedisAddr,
		a.settings.Workers.RedisPassword,
		a.settings.Workers.RedisDB,
		)
	if err != nil {
		return fmt.Errorf("failed to create worker queue: %w", err)
	}
	
			workers.SetDefaultQueue(queue)
	concurrency := a.settings.Workers.Concurrency
			if concurrency <= 0 {
				concurrency = 10
			}
	
			workers.Start(context.Background(), concurrency)
			workers.StartScheduler(context.Background())
	a.workersStarted = true
	return nil
}

func (a *App) setupAdmin() error {
	if a.modelsDB == nil {
		return fmt.Errorf("database connection required for admin - set database.url in settings")
	}
	
	adminRegistry := admin.NewRegistry()
	rendererRegistry := admin.NewRendererRegistry("1.0", true)
	templates := admin.NewTemplateEngine(a.settings.Admin.TemplatePath, rendererRegistry)
	if err := templates.LoadTemplates(); err != nil {
		return fmt.Errorf("failed to load admin templates: %w", err)
		}
	
	a.serviceIntegration = admin.NewServiceIntegration(
		adminRegistry, rendererRegistry, templates, a.modelsDB,
	)
	return nil
}


// RegisterResource registers a REST resource
func RegisterResource[T any, Q any](a *App, name string, resource rest.Resource[T, Q]) {
	if a.router == nil {
		a.setupAPIRouter()
	}
	rest.RegisterResource[T, Q](a.router, name, resource)
}

// RegisterModel registers a model with the API
// Database must be configured in settings
func RegisterModel[T any](a *App, path string, modelDef *models.ModelDefinition[T]) error {
	if a.modelsDB == nil {
		return fmt.Errorf("database not connected - set database.url in settings")
	}
	api.RegisterModel[T](a.fiber, path, a.modelsDB, modelDef)
	return nil
}

// RegisterModelAdmin registers a model with the admin interface
// Admin must be enabled in settings
func RegisterModelAdmin[T any](a *App, modelDef *models.ModelDefinition[T], options ...admin.Option) error {
	if a.serviceIntegration == nil {
		return fmt.Errorf("admin not enabled - set admin.enable=true in settings")
	}
	return admin.RegisterModelFromApp[T](a.serviceIntegration, modelDef, options...)
}

func (a *App) Use(handler fiber.Handler) {
	a.fiber.Use(handler)
}

func (a *App) Get(path string, handler fiber.Handler, middleware ...fiber.Handler) {
	handlers := append([]fiber.Handler{handler}, middleware...)
	a.fiber.Get(path, handlers...)
}

func (a *App) Post(path string, handler fiber.Handler, middleware ...fiber.Handler) {
	handlers := append([]fiber.Handler{handler}, middleware...)
	a.fiber.Post(path, handlers...)
}

func (a *App) Put(path string, handler fiber.Handler, middleware ...fiber.Handler) {
	handlers := append([]fiber.Handler{handler}, middleware...)
	a.fiber.Put(path, handlers...)
}

func (a *App) Delete(path string, handler fiber.Handler, middleware ...fiber.Handler) {
	handlers := append([]fiber.Handler{handler}, middleware...)
	a.fiber.Delete(path, handlers...)
}

func (a *App) Group(prefix string, handlers ...fiber.Handler) fiber.Router {
	return a.fiber.Group(prefix, handlers...)
}


func (a *App) DB() *models.DB {
	return a.modelsDB
}

func (a *App) APIRouter() *rest.Router {
	return a.router
}

func (a *App) Fiber() *fiber.App {
	return a.fiber
}

func (a *App) ServiceIntegration() *admin.ServiceIntegration {
	return a.serviceIntegration
}

func (a *App) Listen(addr string) error {
	if addr == "" {
		addr = fmt.Sprintf("%s:%d", a.settings.Server.Host, a.settings.Server.Port)
	}

	log.Printf("Gogo server starting on %s", addr)
	return a.fiber.Listen(addr)
}

// Shutdown gracefully shuts down the app and all initialized features
func (a *App) Shutdown(ctx context.Context) error {
	var errs []error

	// Shutdown workers if they were started
	if a.workersStarted {
		if err := workers.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("workers stop: %w", err))
		}
		workers.StopScheduler()
	}

	if a.modelsDB != nil {
		if err := a.modelsDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("database close: %w", err))
		}
	}

	if err := a.fiber.Shutdown(); err != nil {
		errs = append(errs, fmt.Errorf("fiber shutdown: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	return nil
}

// Settings returns the app settings
func (a *App) Settings() *AppSettings {
	return a.settings
}

