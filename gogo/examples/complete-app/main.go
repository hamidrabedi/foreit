package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gogo/pkg/auth"
	"github.com/gogo/pkg/endpoints"
	"github.com/gogo/pkg/gogo"
	"github.com/gogo/pkg/routing"
	"github.com/gofiber/fiber/v2"
)

// User model
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" required:"true" min:"3"`
	Email     string    `json:"email" required:"true" email:"true"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository implements the repository interface
type UserRepository struct {
	users  []*User
	nextID int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  make([]*User, 0),
		nextID: 1,
	}
}

func (r *UserRepository) Query() interface{} {
	return r
}

func (r *UserRepository) GetByID(ctx context.Context, id interface{}) (*User, error) {
	idInt := id.(int)
	for _, user := range r.users {
		if user.ID == idInt {
			return user, nil
		}
	}
	return nil, endpoints.ErrNotFound
}

func (r *UserRepository) All(ctx context.Context, query interface{}) ([]*User, error) {
	return r.users, nil
}

func (r *UserRepository) Count(ctx context.Context, query interface{}) (int, error) {
	return len(r.users), nil
}

func (r *UserRepository) Create(ctx context.Context, data *User) (*User, error) {
	user := &User{
		ID:        r.nextID,
		Name:      data.Name,
		Email:     data.Email,
		CreatedAt: time.Now(),
	}
	r.nextID++
	r.users = append(r.users, user)
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, id interface{}, data *User) (*User, error) {
	idInt := id.(int)
	for i, user := range r.users {
		if user.ID == idInt {
			user.Name = data.Name
			user.Email = data.Email
			return user, nil
		}
	}
	return nil, endpoints.ErrNotFound
}

func (r *UserRepository) Delete(ctx context.Context, id interface{}) error {
	idInt := id.(int)
	for i, user := range r.users {
		if user.ID == idInt {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return endpoints.ErrNotFound
}

// UserResource
type UserResource struct {
	*endpoints.BaseResource[*User, interface{}]
}

func NewUserResource() *UserResource {
	repo := NewUserRepository()
	return &UserResource{
		BaseResource: endpoints.NewResource[*User, interface{}](repo),
	}
}

func main() {
	// Create Gogo application
	app, err := gogo.New(&gogo.AppConfig{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		Port:        getEnvInt("PORT", 8080),
		SecretKey:   getEnv("SECRET_KEY", "dev-secret-key"),
		Debug:       getEnvBool("DEBUG", true),
		
		EnableConsole: true,
		EnableWorkers: false,
		EnableCache:   true,
		EnableSessions: true,
		EnableI18n:    false,
		EnableStatic:  false,
		
		ConsolePath: "/admin",
		APIPath:     "/api/v1",
	})
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}
	
	// Setup authentication
	jwtAuth := auth.NewJWT(getEnv("SECRET_KEY", "dev-secret-key"))
	app.SetAuth(jwtAuth, false)
	
	// Register resources
	userResource := NewUserResource()
	app.RegisterResource("users", userResource)
	
	// Add custom routes
	router := app.Router()
	router.Get("/", homeHandler, routing.Name("home"))
	router.Get("/health", healthHandler, routing.Name("health"))
	router.Get("/api", apiInfoHandler, routing.Name("api-info"))
	
	// Start server
	go func() {
		addr := fmt.Sprintf(":%d", getEnvInt("PORT", 8080))
		log.Printf("🚀 Gogo server starting on %s", addr)
		log.Printf("📚 API: http://localhost%s/api/v1", addr)
		log.Printf("🎛️  Admin: http://localhost%s/admin", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	
	log.Println("Server stopped")
}

func homeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Welcome to Gogo Framework",
		"version": "0.1.0",
		"endpoints": fiber.Map{
			"api":   "/api/v1",
			"admin": "/admin",
			"health": "/health",
		},
	})
}

func healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now(),
		"uptime":    "running",
	})
}

func apiInfoHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"version": "v1",
		"resources": []string{
			"users",
		},
		"endpoints": fiber.Map{
			"users": fiber.Map{
				"list":   "GET /api/v1/users",
				"show":   "GET /api/v1/users/:id",
				"create": "POST /api/v1/users",
				"update": "PUT /api/v1/users/:id",
				"delete": "DELETE /api/v1/users/:id",
			},
		},
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

