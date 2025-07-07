package cmd

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"

	"github.com/engineervix/baby-tracker/internal/config"
	"github.com/engineervix/baby-tracker/internal/database"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Baby Tracker server",
	Long:  `Starts the web server for the Baby Tracker application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get port flag value
		port, _ := cmd.Flags().GetString("port")
		runServer(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Add serve-specific flags
	serveCmd.Flags().StringP("port", "p", "", "Port to run the server on (overrides config)")
}

func runServer(overridePort string) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(db, cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configure CORS
	corsConfig := middleware.CORSConfig{
		AllowOrigins:     strings.Split(cfg.AllowedOrigins, ","),
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true, // Important for session cookies
	}
	e.Use(middleware.CORSWithConfig(corsConfig))

	// Store config and db in context for handlers
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			c.Set("config", cfg)
			return next(c)
		}
	})

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Baby Tracker API",
			"version": "1.0.0",
		})
	})

	// API routes group
	// api := e.Group("/api")

	// TODO: Add API routes here
	// api.POST("/auth/login", handlers.Login)
	// api.POST("/auth/logout", handlers.Logout)
	// api.GET("/auth/me", handlers.GetCurrentUser)

	// Serve static files in production
	if cfg.Env == "production" {
		e.Static("/", "web/dist")
		// Catch-all route for Vue Router
		e.GET("/*", func(c echo.Context) error {
			return c.File("web/dist/index.html")
		})
	}

	// Start server
	port := cfg.Port
	if overridePort != "" {
		port = overridePort
	}

	log.Printf("Starting server on port %s in %s mode", port, cfg.Env)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
