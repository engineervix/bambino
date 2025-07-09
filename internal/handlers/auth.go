package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/engineervix/bambino/internal/models"
	"github.com/engineervix/bambino/internal/utils"
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

// UserResponse represents the user response
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Login handles user authentication
func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Basic validation
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username and password are required")
	}

	// Get database from context
	db, ok := c.Get("db").(*gorm.DB)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "database connection error")
	}

	// Find user by username
	var user models.User
	if err := db.Where("username = ?", strings.TrimSpace(req.Username)).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid username or password")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "database error")
	}

	// Verify password
	valid, err := utils.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "password verification error")
	}
	if !valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid username or password")
	}

	// Create session
	if err := utils.CreateUserSession(c, user.ID, user.Username); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "session creation error")
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Message:  "login successful",
		Username: user.Username,
	})
}

// Logout handles user logout
func Logout(c echo.Context) error {
	// Destroy session
	if err := utils.DestroyUserSession(c); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "logout error")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "logout successful",
	})
}

// GetCurrentUser returns the current authenticated user
func GetCurrentUser(c echo.Context) error {
	// Get user from session
	userID, _, err := utils.GetUserSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "not authenticated")
	}

	// Get database from context
	db, ok := c.Get("db").(*gorm.DB)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "database connection error")
	}

	// Verify user still exists in database
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User was deleted, destroy session
			utils.DestroyUserSession(c)
			return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "database error")
	}

	return c.JSON(http.StatusOK, UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
	})
}

// CheckAuth is a simple endpoint to check if user is authenticated
func CheckAuth(c echo.Context) error {
	if !utils.IsAuthenticated(c) {
		return echo.NewHTTPError(http.StatusUnauthorized, "not authenticated")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "authenticated",
	})
}
