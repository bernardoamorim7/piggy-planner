package handlers

import (
	"net/http"

	"piggy-planner/internal/models"
	"piggy-planner/internal/services"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if !models.ValidateEmail(email) {
		return c.String(http.StatusBadRequest, "Invalid email")
	}

	if !models.ValidatePassword(password) {
		return c.String(http.StatusBadRequest, "Password must be at least 8 characters long")
	}

	userService := services.NewUserService()
	user, err := userService.GetByEmail(email)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Account not found")
	}

	if !user.ComparePassword(password) {
		return c.String(http.StatusUnauthorized, "Wrong password")
	}

	sess, _ := session.Get("piggysession", c)
	sess.Values["authenticated"] = true
	sess.Values["userID"] = user.ID
	sess.Values["name"] = user.Name
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/dashboard")
}

func Protected() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("piggysession", c)
			if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
				return c.String(http.StatusUnauthorized, "You are not logged in")
			}
			return next(c)
		}
	}
}

func Register(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	passwordConfirm := c.FormValue("password_confirm")

	if password != passwordConfirm {
		return c.String(http.StatusBadRequest, "Passwords do not match")
	}

	if !models.ValidateEmail(email) {
		return c.String(http.StatusBadRequest, "Invalid email")
	}

	if !models.ValidatePassword(password) {
		return c.String(http.StatusBadRequest, "Password must be at least 8 characters long")
	}

	userService := services.NewUserService()
	_, err := userService.GetByEmail(email)
	if err == nil {
		return c.String(http.StatusConflict, "Account already exists")
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	err = userService.Create(&user)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}
