package middlewares

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func GetSessionVars(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("piggysession", c)
		if err != nil {
			return err
		}

		c.Set("userID", sess.Values["userID"].(int64))
		c.Set("name", sess.Values["name"].(string))
		c.Set("avatar", sess.Values["avatar"].(string))

		return next(c)
	}

}
