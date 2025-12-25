package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// Logging returns a logging middleware
func Logging() fiber.Handler {
	return logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	})
}

// Recovery returns a recovery middleware that handles panics
func Recovery() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}

// CORS returns a CORS middleware
func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	})
}

// RequestID returns a middleware that adds a request ID
func RequestID() fiber.Handler {
	return requestid.New()
}

// Timing adds request timing to context
func Timing() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := NewContext(c)
		startTime := time.Now()
		ctx.Set("start_time", startTime)
		
		err := c.Next()
		
		if startTimeVal, ok := ctx.Get("start_time"); ok {
			if startTime, ok := startTimeVal.(time.Time); ok {
				duration := time.Since(startTime)
				ctx.Set("duration", duration)
				log.Printf("Request took %v", duration)
			}
		}
		
		return err
	}
}

// UserContext extracts user from request and adds to context
func UserContext(userExtractor func(*fiber.Ctx) interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if userExtractor != nil {
			user := userExtractor(c)
			if user != nil {
				c.Locals("user", user)
			}
		}
		return c.Next()
	}
}

// Compress adds compression middleware
func Compress() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})
}

// SecurityHeaders adds security headers
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		return c.Next()
	}
}

// RateLimit provides rate limiting (basic implementation)
func RateLimit(maxRequests int, window time.Duration) fiber.Handler {
	// Simple in-memory rate limiter
	// In production, use a proper rate limiting library
	requests := make(map[string][]time.Time)
	
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		now := time.Now()
		
		// Clean old requests
		if times, ok := requests[ip]; ok {
			validTimes := make([]time.Time, 0)
			for _, t := range times {
				if now.Sub(t) < window {
					validTimes = append(validTimes, t)
				}
			}
			requests[ip] = validTimes
			
			if len(validTimes) >= maxRequests {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Rate limit exceeded",
				})
			}
		}
		
		// Record this request
		if requests[ip] == nil {
			requests[ip] = make([]time.Time, 0)
		}
		requests[ip] = append(requests[ip], now)
		
		return c.Next()
	}
}

// Chain chains multiple middleware handlers
func Chain(handlers ...fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		for _, handler := range handlers {
			if err := handler(c); err != nil {
				return err
			}
		}
		return c.Next()
	}
}

