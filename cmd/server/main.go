package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/engineervix/baby-tracker/internal/config"
	"github.com/engineervix/baby-tracker/internal/database"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	dbConfig := database.NewConfig()
	if err := database.Connect(dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run database migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Create Echo instance
	e := echo.New()

	// Hide Echo banner
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS middleware with configuration
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{cfg.AllowedOrigins},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Request ID middleware
	e.Use(middleware.RequestID())

	// Secure middleware
	if cfg.Env == "production" {
		e.Use(middleware.Secure())
	}

	// Routes
	setupRoutes(e)

	// Start server with graceful shutdown
	go func() {
		port := cfg.Port
		if port == "" {
			port = "8080"
		}

		log.Printf("Starting server on port %s in %s mode", port, cfg.Env)
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown complete")
}

func setupRoutes(e *echo.Echo) {
	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API routes group
	api := e.Group("/api")

	// Public routes (no auth required)
	api.POST("/auth/login", handleLogin)

	// Protected routes (auth required) - to be added
	// auth := api.Group("", authMiddleware)
	// auth.GET("/auth/me", handleGetCurrentUser)
	// auth.POST("/auth/logout", handleLogout)
	// ... other protected routes

	// Serve static files in production
	if os.Getenv("ENV") == "production" {
		e.Static("/", "web/dist")
		e.File("/*", "web/dist/index.html")
	}
}

// Placeholder handlers - to be implemented
func handleLogin(c echo.Context) error {
	// TODO: Implement login logic
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login endpoint - to be implemented",
	})
}
