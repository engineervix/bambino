package utils

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	SessionName = "baby-tracker-session"
	UserIDKey   = "user_id"
	UsernameKey = "username"
)

// SessionConfig holds session configuration
type SessionConfig struct {
	Secret   string
	MaxAge   int // in seconds
	HttpOnly bool
	Secure   bool // set to true in production with HTTPS
}

// CreateSessionStore creates a new session store
func CreateSessionStore(config SessionConfig) sessions.Store {
	store := sessions.NewCookieStore([]byte(config.Secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   config.MaxAge,
		HttpOnly: config.HttpOnly,
		Secure:   config.Secure,
		SameSite: http.SameSiteStrictMode,
	}
	return store
}

// CreateUserSession creates a new session for the user
func CreateUserSession(c echo.Context, userID uuid.UUID, username string) error {
	sess, err := session.Get(SessionName, c)
	if err != nil {
		return err
	}

	sess.Values[UserIDKey] = userID.String()
	sess.Values[UsernameKey] = username

	return sess.Save(c.Request(), c.Response())
}

// GetUserSession retrieves user information from session
func GetUserSession(c echo.Context) (uuid.UUID, string, error) {
	sess, err := session.Get(SessionName, c)
	if err != nil {
		return uuid.Nil, "", err
	}

	userIDStr, ok := sess.Values[UserIDKey].(string)
	if !ok {
		return uuid.Nil, "", ErrSessionNotFound
	}

	username, ok := sess.Values[UsernameKey].(string)
	if !ok {
		return uuid.Nil, "", ErrSessionNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, "", err
	}

	return userID, username, nil
}

// DestroyUserSession removes the user session
func DestroyUserSession(c echo.Context) error {
	sess, err := session.Get(SessionName, c)
	if err != nil {
		return err
	}

	// Set MaxAge to -1 to delete the session
	sess.Options.MaxAge = -1
	delete(sess.Values, UserIDKey)
	delete(sess.Values, UsernameKey)

	return sess.Save(c.Request(), c.Response())
}

// IsAuthenticated checks if the user is authenticated
func IsAuthenticated(c echo.Context) bool {
	_, _, err := GetUserSession(c)
	return err == nil
}

// GetSessionExpiry calculates session expiry time
func GetSessionExpiry(maxAge int) time.Time {
	return time.Now().Add(time.Duration(maxAge) * time.Second)
}

// Custom errors
var (
	ErrSessionNotFound = echo.NewHTTPError(401, "session not found")
	ErrSessionInvalid  = echo.NewHTTPError(401, "session invalid")
)
