package middlewares

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// Protected is a middleware that allows only authenticated user to access certains pages/views
func Protected() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("piggysession", c)
			if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
				return c.Redirect(http.StatusSeeOther, "/login")
			}
			return next(c)
		}
	}
}

// RedirectIfLoggedIn is a middleware that redirects a user to the dashboard if they are already logged in.
func RedirectIfLoggedIn() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("piggysession", c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "server error")
			}

			if auth, ok := sess.Values["authenticated"].(bool); ok && auth {
				c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
				c.Response().Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
				c.Response().Header().Set("Expires", "0")                                         // Proxies.
				return c.Redirect(http.StatusSeeOther, "/dashboard")
			}

			return next(c)
		}
	}
}
