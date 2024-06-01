package middlewares

import (
	"log/slog"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func GetSessionVars() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("piggysession", c)
			if err != nil {
				return err
			}

			if sess.Values["userID"] != nil {
				if value, ok := sess.Values["userID"].(int64); ok {
					c.Set("userID", value)
				} else {
					slog.WarnContext(c.Request().Context(), "userID is not of type int64")
				}
			}

			if sess.Values["name"] != nil {
				if value, ok := sess.Values["name"].(string); ok {
					c.Set("name", value)
				} else {
					slog.WarnContext(c.Request().Context(), "name is not of type string")
				}
			}

			if sess.Values["avatar"] != nil {
				if value, ok := sess.Values["avatar"].(string); ok {
					c.Set("avatar", value)
				} else {
					slog.WarnContext(c.Request().Context(), "avatar is not of type string")
				}
			}

			return next(c)
		}
	}
}
