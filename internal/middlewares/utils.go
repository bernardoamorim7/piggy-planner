package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// IsLocalhost is a middleware that only allows requests from localhost.
func IsLocalhost() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			if ip != "127.0.0.1" {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			return next(c)
		}
	}
}
