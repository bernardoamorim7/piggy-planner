package handlers

import (
	"errors"
	"net/http"
	"strings"

	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
	"piggy-planner/internal/services"

	"github.com/gorilla/sessions"
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

	db, err := database.New()
	if err != nil {
		return err
	}

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
	sess.Values["email"] = user.Email
	sess.Values["avatar"] = user.Avatar
	sess.Values["is_admin"] = user.IsAdmin

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	// Create security log
	logService := services.NewSecurityLogsService(db)
	log := models.SecurityLog{
		User:      *user,
		Action:    "login",
		IPAdress:  c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}
	_ = logService.Create(&log)

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

	db, err := database.New()
	if err != nil {
		return err
	}

	userService := services.NewUserService(db)

	_, err = userService.GetByEmail(email)
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

	newUser, err := userService.GetByEmail(user.Email)
	if err != nil {
		return err
	}

	// Create security log
	logService := services.NewSecurityLogsService(db)
	log := models.SecurityLog{
		User:      *newUser,
		Action:    "register",
		IPAdress:  c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}
	_ = logService.Create(&log)

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusOK)
}

func Logout(c echo.Context) error {
	sess, _ := session.Get("piggysession", c)

	userID := sess.Values["userID"].(uint64)

	db, err := database.New()
	if err != nil {
		return err
	}

	user := &models.User{
		ID: userID,
	}

	// Create security log
	logService := services.NewSecurityLogsService(db)
	log := models.SecurityLog{
		User:      *user,
		Action:    "logout",
		IPAdress:  c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}
	_ = logService.Create(&log)

	// Clear session
	sess.Options = &sessions.Options{
		MaxAge:   -1,
		HttpOnly: true,
	}
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return errors.New("Failed to delete session")
	}

	// Clear cookie
	c.SetCookie(&http.Cookie{
		Name:   "piggysession",
		Value:  "",
		MaxAge: -1,
	})

	// Clear possible context variables
	c.Set("authenticated", nil)
	c.Set("userID", nil)
	c.Set("name", nil)
	c.Set("email", nil)
	c.Set("avatar", nil)
	c.Set("is_admin", nil)

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}
