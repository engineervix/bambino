package utils

import (
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
)

// CaptureError captures an error using Sentry with Echo context
func CaptureError(c echo.Context, err error, message string) {
	if hub := sentryecho.GetHubFromContext(c); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("handler", c.Path())
			scope.SetTag("method", c.Request().Method)
			scope.SetExtra("message", message)
			if c.Get("user") != nil {
				scope.SetUser(sentry.User{
					ID: c.Get("user").(string),
				})
			}
			hub.CaptureException(err)
		})
	}
}

// CaptureMessage captures a message using Sentry with Echo context
func CaptureMessage(c echo.Context, message string, level sentry.Level) {
	if hub := sentryecho.GetHubFromContext(c); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("handler", c.Path())
			scope.SetTag("method", c.Request().Method)
			if c.Get("user") != nil {
				scope.SetUser(sentry.User{
					ID: c.Get("user").(string),
				})
			}
			hub.CaptureMessage(message)
		})
	}
}
