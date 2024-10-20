package middlewares

import (
	"fmt"

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
				if value, ok := sess.Values["userID"].(uint64); ok {
					c.Set("userID", value)
				} else {
					c.Logger().Debug(c.Request().Context(), "userID is not of type uint64 ("+fmt.Sprintf("%T", sess.Values["userID"])+")")
				}
			}

			if sess.Values["name"] != nil {
				if value, ok := sess.Values["name"].(string); ok {
					c.Set("name", value)
				} else {
					c.Logger().Debug(c.Request().Context(), "name is not of type string ("+fmt.Sprintf("%T", sess.Values["name"])+")")
				}
			}

			if sess.Values["email"] != nil {
				if value, ok := sess.Values["email"].(string); ok {
					c.Set("email", value)
				} else {
					c.Logger().Debug(c.Request().Context(), "email is not of type string ("+fmt.Sprintf("%T", sess.Values["email"])+")")
				}
			}

			if sess.Values["avatar"] != nil {
				if value, ok := sess.Values["avatar"].(string); ok {
					c.Set("avatar", value)
				} else {
					c.Logger().Debug(c.Request().Context(), "avatar is not of type string ("+fmt.Sprintf("%T", sess.Values["avatar"])+")")
				}
			}

			if sess.Values["is_admin"] != nil {
				if value, ok := sess.Values["is_admin"].(bool); ok {
					c.Set("is_admin", value)
				} else {
					c.Logger().Debug(c.Request().Context(), "is_admin is not of type bool ("+fmt.Sprintf("%T", sess.Values["is_admin"])+")")
				}
			}

			return next(c)
		}
	}
}
