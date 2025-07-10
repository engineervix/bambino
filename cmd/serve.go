package cmd

import (
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/engineervix/bambino/internal/assets"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"

	"github.com/engineervix/bambino/internal/config"
	"github.com/engineervix/bambino/internal/database"
	"github.com/engineervix/bambino/internal/handlers"
	authMiddleware "github.com/engineervix/bambino/internal/middleware"
	"github.com/engineervix/bambino/internal/utils"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Bambino server",
	Long:  `Starts the web server for the Bambino application.`,
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

	// Initialize Sentry
	if cfg.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.SentryDSN,
			Environment:      cfg.Env,
			TracesSampleRate: cfg.SentryTracesSampleRate,
			SendDefaultPII:   true, // Adds request headers and IP for users
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				// Filter out sensitive information in development
				if cfg.Env == "development" {
					if hint.Context != nil {
						if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
							// You can modify the event here based on the request
							_ = req // Use the request if needed
						}
					}
				}
				return event
			},
		})
		if err != nil {
			log.Printf("Sentry initialization failed: %v", err)
		} else {
			log.Printf("Sentry initialized successfully")
			// Flush buffered events on exit
			defer sentry.Flush(2 * time.Second)
		}
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

	// Configure IP extraction for Traefik proxy
	// Traefik uses X-Forwarded-For header and runs in Docker network
	e.IPExtractor = echo.ExtractIPFromXFFHeader(
		echo.TrustLoopback(true),   // Trust loopback addresses
		echo.TrustLinkLocal(true),  // Trust link-local addresses
		echo.TrustPrivateNet(true), // Trust private network addresses (Docker networks)
	)

	// Basic middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Add Sentry middleware if configured
	if cfg.SentryDSN != "" {
		e.Use(sentryecho.New(sentryecho.Options{
			Repanic:         true,
			WaitForDelivery: false,
			Timeout:         3 * time.Second,
		}))
	}

	// Configure CORS
	corsConfig := middleware.CORSConfig{
		AllowOrigins:     strings.Split(cfg.AllowedOrigins, ","),
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true, // Important for session cookies
	}
	e.Use(middleware.CORSWithConfig(corsConfig))

	// Session middleware
	sessionStore := utils.CreateSessionStore(utils.SessionConfig{
		Secret:   cfg.SessionSecret,
		MaxAge:   cfg.SessionMaxAge,
		HttpOnly: true,
		Secure:   cfg.Env == "production", // Only secure in production
	})
	e.Use(session.Middleware(sessionStore))

	// Store config and db in context for handlers
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			c.Set("config", cfg)
			return next(c)
		}
	})

	// Public routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Bambino API",
			"version": "1.0.0",
		})
	})

	// Auth routes (public)
	auth := e.Group("/api/auth")
	auth.POST("/login", handlers.Login)
	auth.POST("/logout", handlers.Logout)
	auth.GET("/check", handlers.CheckAuth)

	// Auth routes (protected)
	authProtected := e.Group("/api/auth")
	authProtected.Use(authMiddleware.RequireAuthJSON())
	authProtected.GET("/me", handlers.GetCurrentUser)

	// Protected API routes
	api := e.Group("/api")
	api.Use(authMiddleware.RequireAuthJSON())

	// Baby routes
	api.GET("/babies", handlers.GetBabies)

	// Activity routes
	api.GET("/activities", handlers.GetActivities)
	api.POST("/activities", handlers.CreateActivity)
	api.GET("/activities/:id", handlers.GetActivity)
	api.PUT("/activities/:id", handlers.UpdateActivity)
	api.DELETE("/activities/:id", handlers.DeleteActivity)

	// Timer endpoints
	api.POST("/activities/timer/start", handlers.StartActivityTimer)
	api.PUT("/activities/timer/:id/stop", handlers.StopActivityTimer)

	// Statistics routes
	api.GET("/stats/daily", handlers.GetDailyStats)
	api.GET("/stats/recent", handlers.GetRecentStats)
	api.GET("/stats/weekly", handlers.GetWeeklyStats)

	// Serve static files in production
	if cfg.Env == "production" {
		web, err := fs.Sub(assets.Assets, "dist")
		if err != nil {
			log.Fatal("failed to create sub filesystem", err)
		}
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       ".",
			Index:      "index.html",
			HTML5:      true,
			Filesystem: http.FS(web),
		}))
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
