package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/engineervix/baby-tracker/internal/utils"
)

// AuthConfig defines the configuration for auth middleware
type AuthConfig struct {
	// Skipper defines a function to skip middleware
	Skipper middleware.Skipper

	// Optional custom error handler
	ErrorHandler func(c echo.Context, err error) error
}

// DefaultAuthConfig provides default configuration
var DefaultAuthConfig = AuthConfig{
	Skipper: middleware.DefaultSkipper,
	ErrorHandler: func(c echo.Context, err error) error {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	},
}

// RequireAuth returns auth middleware with default config
func RequireAuth() echo.MiddlewareFunc {
	return RequireAuthWithConfig(DefaultAuthConfig)
}

// RequireAuthWithConfig returns auth middleware with custom config
func RequireAuthWithConfig(config AuthConfig) echo.MiddlewareFunc {
	// Set defaults
	if config.Skipper == nil {
		config.Skipper = DefaultAuthConfig.Skipper
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = DefaultAuthConfig.ErrorHandler
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip if configured to skip
			if config.Skipper(c) {
				return next(c)
			}

			// Check if user is authenticated
			userID, username, err := utils.GetUserSession(c)
			if err != nil {
				return config.ErrorHandler(c, err)
			}

			// Add user info to context for use in handlers
			c.Set("user_id", userID.String()) // Convert UUID to string
			c.Set("username", username)

			return next(c)
		}
	}
}

// RequireAuthJSON returns auth middleware that always returns JSON errors
func RequireAuthJSON() echo.MiddlewareFunc {
	return RequireAuthWithConfig(AuthConfig{
		Skipper: middleware.DefaultSkipper,
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "authentication required",
			})
		},
	})
}

// SkipAuth creates a skipper function that skips auth for specific paths
func SkipAuth(skipPaths ...string) middleware.Skipper {
	skipMap := make(map[string]bool)
	for _, path := range skipPaths {
		skipMap[path] = true
	}

	return func(c echo.Context) bool {
		return skipMap[c.Request().URL.Path]
	}
}

// SkipAuthPrefix creates a skipper function that skips auth for paths with specific prefixes
func SkipAuthPrefix(prefixes ...string) middleware.Skipper {
	return func(c echo.Context) bool {
		path := c.Request().URL.Path
		for _, prefix := range prefixes {
			if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
				return true
			}
		}
		return false
	}
}
