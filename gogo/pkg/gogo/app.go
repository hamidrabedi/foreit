package gogo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gogo/pkg/auth"
	"github.com/gogo/pkg/cache"
	"github.com/gogo/pkg/console"
	"github.com/gogo/pkg/endpoints"
	"github.com/gogo/pkg/i18n"
	"github.com/gogo/pkg/orm"
	"github.com/gogo/pkg/pipeline"
	"github.com/gogo/pkg/routing"
	"github.com/gogo/pkg/sessions"
	"github.com/gogo/pkg/settings"
	"github.com/gogo/pkg/static"
	"github.com/gogo/pkg/workers"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	fiber     *fiber.App
	config    *AppConfig
	client    *orm.Client
	router    *routing.Router
	apiRouter *endpoints.Router
}

type AppConfig struct {
	Settings    interface{}
	DatabaseURL string
	Port        int
	SecretKey   string
	Debug       bool
	
	EnableConsole  bool
	EnableWorkers  bool
	EnableCache    bool
	EnableSessions bool
	EnableI18n     bool
	EnableStatic   bool
	
	ConsolePath   string
	APIPath       string
	StaticPath    string
	StaticRoot    string
	LocalesPath   string
	DefaultLocale string
	
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	WorkersConcurrency int
}

func DefaultConfig() *AppConfig {
	return &AppConfig{
		Port: 8080,
		EnableConsole: true,
		EnableWorkers: true,
		EnableCache: true,
		EnableSessions: true,
		EnableI18n: false,
		EnableStatic: true,
		ConsolePath: "/console",
		APIPath: "/api/v1",
		StaticPath: "/static",
		StaticRoot: "./public",
		LocalesPath: "./locales",
		DefaultLocale: "en",
		RedisAddr: "localhost:6379",
		RedisPassword: "",
		RedisDB: 0,
		WorkersConcurrency: 10,
	}
}

func New(config *AppConfig) (*App, error) {
	if config == nil {
		config = DefaultConfig()
	}
	
	app := &App{
		fiber: fiber.New(fiber.Config{
			ErrorHandler: endpoints.HandleError,
		}),
		config: config,
	}
	
	if config.DatabaseURL != "" {
		client, err := orm.NewClient("postgres", config.DatabaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		app.client = client
	}
	
	app.router = routing.NewRouter(app.fiber)
	app.apiRouter = endpoints.NewRouter(app.fiber, config.APIPath)
	
	app.setupMiddleware()
	app.setupModules()
	
	return app, nil
}

func (a *App) setupMiddleware() {
	a.fiber.Use(pipeline.Logging())
	a.fiber.Use(pipeline.Recovery())
	a.fiber.Use(pipeline.CORS())
	a.fiber.Use(pipeline.RequestID())
	a.fiber.Use(pipeline.SecurityHeaders())
	
	if a.config.Debug {
		a.fiber.Use(pipeline.RateLimit(1000, time.Minute))
	} else {
		a.fiber.Use(pipeline.RateLimit(100, time.Minute))
	}
}

func (a *App) setupModules() {
	if a.config.EnableSessions {
		sessionStore := sessions.NewMemoryStore(sessions.DefaultConfig())
		a.fiber.Use(sessions.Middleware(sessionStore))
	}
	
	if a.config.EnableI18n && a.config.LocalesPath != "" {
		if err := i18n.Load(a.config.LocalesPath, a.config.DefaultLocale); err != nil {
			log.Printf("Warning: Failed to load translations: %v", err)
		} else {
			a.fiber.Use(i18n.Middleware())
		}
	}
	
	if a.config.EnableStatic && a.config.StaticRoot != "" {
		a.fiber.Use(static.New(static.Config{
			Root:     a.config.StaticRoot,
			Prefix:   a.config.StaticPath,
			Compress: true,
		}))
	}
	
	if a.config.EnableConsole {
		console.InstallRoutes(a.fiber, a.config.ConsolePath)
	}
	
	if a.config.EnableWorkers {
		queue, err := workers.NewAsynqQueue(
			a.config.RedisAddr,
			a.config.RedisPassword,
			a.config.RedisDB,
		)
		if err == nil {
			workers.SetDefaultQueue(queue)
			concurrency := a.config.WorkersConcurrency
			if concurrency <= 0 {
				concurrency = 10
			}
			workers.Start(context.Background(), concurrency)
			workers.StartScheduler(context.Background())
		} else {
			log.Printf("Warning: Failed to initialize workers: %v", err)
		}
	}
	
	if a.config.EnableCache {
		if a.config.RedisAddr != "" {
			redisStore, err := cache.NewRedisStore(
				a.config.RedisAddr,
				a.config.RedisPassword,
				a.config.RedisDB,
			)
			if err == nil {
				cache.SetDefaultStore(redisStore)
			} else {
				log.Printf("Warning: Failed to initialize Redis cache, using memory store: %v", err)
				cache.SetDefaultStore(cache.NewTaggedMemoryStore())
			}
		} else {
			cache.SetDefaultStore(cache.NewTaggedMemoryStore())
		}
	}
}

func (a *App) RegisterResource[T any, Q any](name string, resource endpoints.Resource[T, Q]) {
	a.apiRouter.RegisterResource(name, resource)
}

func (a *App) RegisterConsole[T any](console console.Console[T]) {
	console.Register(console)
}

func (a *App) Use(handler fiber.Handler) {
	a.fiber.Use(handler)
}

func (a *App) Get(path string, handler fiber.Handler, options ...routing.RouteOption) {
	a.router.Get(path, handler, options...)
}

func (a *App) Post(path string, handler fiber.Handler, options ...routing.RouteOption) {
	a.router.Post(path, handler, options...)
}

func (a *App) Put(path string, handler fiber.Handler, options ...routing.RouteOption) {
	a.router.Put(path, handler, options...)
}

func (a *App) Delete(path string, handler fiber.Handler, options ...routing.RouteOption) {
	a.router.Delete(path, handler, options...)
}

func (a *App) Group(prefix string, handlers ...fiber.Handler) *fiber.Router {
	return a.router.Group(prefix, handlers...)
}

func (a *App) SetAuth(authenticator auth.Authenticator, required bool) {
	if required {
		a.fiber.Use(auth.Middleware(authenticator))
	} else {
		a.fiber.Use(auth.Optional(authenticator))
	}
}

func (a *App) Client() *orm.Client {
	return a.client
}

func (a *App) Router() *routing.Router {
	return a.router
}

func (a *App) APIRouter() *endpoints.Router {
	return a.apiRouter
}

func (a *App) Fiber() *fiber.App {
	return a.fiber
}

func (a *App) Listen(addr string) error {
	if addr == "" {
		addr = fmt.Sprintf(":%d", a.config.Port)
	}

	log.Printf("Gogo server starting on %s", addr)
	return a.fiber.Listen(addr)
}

func (a *App) Shutdown(ctx context.Context) error {
	var errs []error

	if a.config.EnableWorkers {
		if err := workers.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("workers stop: %w", err))
		}
		workers.StopScheduler()
	}

	if a.client != nil {
		if err := a.client.Close(); err != nil {
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

