package handlers

import (
	"net/http"
	"strings"

	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
	"piggy-planner/internal/services"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	email := strings.TrimSpace(c.FormValue("email"))
	password := strings.TrimSpace(c.FormValue("password"))

	if !models.ValidateEmail(email) {
		return c.String(http.StatusBadRequest, "Invalid email")
	}

	if !models.ValidatePassword(password) {
		return c.String(http.StatusBadRequest, "Password must be at least 8 characters long")
	}

	db := database.New()
	defer db.Close()

	userService := services.NewUserService(db)

	user, err := userService.GetByEmail(email)
	if err != nil {
		return c.String(http.StatusNotFound, "Account not found")
	}

	if !user.ComparePassword(password) {
		return c.String(http.StatusUnauthorized, "Wrong password")
	}

	sess, _ := session.Get("piggysession", c)
	sess.Values["authenticated"] = true
	sess.Values["userID"] = user.ID
	sess.Values["name"] = user.Name
	sess.Values["avatar"] = user.Avatar
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func Register(c echo.Context) error {
	name := strings.TrimSpace(c.FormValue("name"))
	email := strings.TrimSpace(c.FormValue("email"))
	password := strings.TrimSpace(c.FormValue("password"))
	passwordConfirm := strings.TrimSpace(c.FormValue("password_confirm"))

	if password != passwordConfirm {
		return c.String(http.StatusBadRequest, "Passwords do not match")
	}

	if !models.ValidateEmail(email) {
		return c.String(http.StatusBadRequest, "Invalid email")
	}

	if !models.ValidatePassword(password) {
		return c.String(http.StatusBadRequest, "Password must be at least 8 characters long")
	}

	db := database.New()
	defer db.Close()

	userService := services.NewUserService(db)

	_, err := userService.GetByEmail(email)
	if err == nil {
		return c.String(http.StatusConflict, "Account already exists")
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
		Avatar:   "https://api.dicebear.com/8.x/thumbs/png?seed=" + name,
	}

	err = userService.Create(&user)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusOK)
}
