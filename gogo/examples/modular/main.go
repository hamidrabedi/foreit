package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gogo/pkg/auth"
	"github.com/gogo/pkg/console"
	"github.com/gogo/pkg/endpoints"
	"github.com/gogo/pkg/orm"
	"github.com/gogo/pkg/pipeline"
	"github.com/gogo/pkg/routing"
	"github.com/gogo/pkg/settings"
	"github.com/gogo/pkg/sessions"
	"github.com/gogo/pkg/workers"
	"github.com/gofiber/fiber/v2"
)

// AppSettings defines application configuration
type AppSettings struct {
	DatabaseURL string `env:"DATABASE_URL" required:"true"`
	Port        int    `env:"PORT" default:"8080"`
	Debug       bool   `env:"DEBUG" default:"false"`
	SecretKey   string `env:"SECRET_KEY" required:"true"`
}

func main() {
	// Load settings
	cfg, err := settings.Load[AppSettings]()
	if err != nil {
		log.Fatalf("Failed to load settings: %v", err)
	}

	// Setup database (placeholder - would use actual Ent client)
	client, err := orm.NewClient("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Close()

	// Setup Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: endpoints.HandleError,
	})

	// Pipeline - Middleware
	app.Use(pipeline.Logging())
	app.Use(pipeline.Recovery())
	app.Use(pipeline.CORS())
	app.Use(pipeline.RequestID())
	app.Use(pipeline.SecurityHeaders())
	app.Use(pipeline.RateLimit(100, time.Minute))

	// Sessions
	sessionStore := sessions.NewMemoryStore(sessions.DefaultConfig())
	app.Use(sessions.Middleware(sessionStore))

	// Auth
	jwtAuth := auth.NewJWT(cfg.SecretKey)
	app.Use(auth.Optional(jwtAuth))

	// Routing
	router := routing.NewRouter(app)

	// API Endpoints
	apiRouter := endpoints.NewRouter(app, "/api/v1")

	// Example: Register a resource (would use actual Ent models)
	// userResource := &UserResource{
	//     Resource: endpoints.NewResource[models.User](client),
	// }
	// apiRouter.RegisterResource("users", userResource)

	// Console (Admin)
	console.Register[interface{}](nil) // Placeholder
	console.InstallRoutes(app, "/console")

	// Workers - Setup Asynq queue
	queue, err := workers.NewAsynqQueue("localhost:6379", "", 0)
	if err != nil {
		log.Printf("Warning: Failed to initialize workers queue: %v", err)
	} else {
		workers.SetDefaultQueue(queue)
		workers.Register("send_email", &EmailJobHandler{})
		workers.Start(context.Background(), 5)
		defer workers.Stop()

		// Scheduler
		workers.Schedule("0 6 * * *", &DailyReportJob{})
		workers.StartScheduler(context.Background())
		defer workers.StopScheduler()
	}

	// Routes
	router.Get("/", homeHandler, routing.Name("home"))
	router.Get("/health", healthHandler, routing.Name("health"))

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Handlers
func homeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Welcome to Gogo Framework",
		"version": "0.1.0",
	})
}

func healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"time":   time.Now(),
	})
}

// Example job handler
type EmailJobHandler struct{}

func (h *EmailJobHandler) Handle(ctx context.Context, job workers.Job) error {
	// Process email job
	return nil
}

// Example scheduled job
type DailyReportJob struct {
	workers.BaseJob
}

func (j *DailyReportJob) Execute(ctx context.Context) error {
	// Generate daily report
	return nil
}

