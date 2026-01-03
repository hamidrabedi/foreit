package main

import (
	"context"
	"log"
	"time"

	"github.com/etemademan/backend/internal/handler"
	"github.com/etemademan/backend/internal/repository"
	"github.com/etemademan/backend/internal/usecase"
	"github.com/etemademan/backend/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Database connection
	dbPool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer dbPool.Close()

	// Initialize Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Repositories
	userRepo := repository.NewPostgresUserRepository(dbPool)
	websiteRepo := repository.NewPostgresWebsiteRepository(dbPool)
	reviewRepo := repository.NewPostgresReviewRepository(dbPool)

	// Usecases
	userUsecase := usecase.NewUserUsecase(userRepo, cfg.JWTSecret, time.Second*10)
	websiteUsecase := usecase.NewWebsiteUsecase(websiteRepo, reviewRepo, time.Second*10)
	reviewUsecase := usecase.NewReviewUsecase(reviewRepo, time.Second*10)

	// Handlers
	handler.NewUserHandler(app, userUsecase)
	handler.NewWebsiteHandler(app, websiteUsecase)
	handler.NewReviewHandler(app, reviewUsecase)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Etemademan API is running",
		})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
