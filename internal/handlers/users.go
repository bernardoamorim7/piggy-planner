package handlers

import (
	"net/http"
	"strconv"

	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
	"piggy-planner/internal/services"
	usersComponents "piggy-planner/web/components/users"

	"github.com/labstack/echo/v4"
)

func GetAllUsers(c echo.Context) error {
	db, err := database.New()
	if err != nil {
		return err
	}

	userService := services.NewUserService(db)

	users, err := userService.GetAll()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		_ = render(c, http.StatusNotFound, usersComponents.UsersNotFound())
		return nil
	}

	for i := range users {
		_ = render(c, http.StatusOK, usersComponents.UserRow(users[i]))
	}

	return nil
}

func GetUserByName(c echo.Context) error {
	userName := c.FormValue("search")

	db, er := database.New()
	if er != nil {
		return er
	}

	usersService := services.NewUserService(db)

	var (
		users []models.User
		err   error
	)

	if userName == "" {
		users, err = usersService.GetAll()
		if err != nil {
			return err
		}
	} else {
		users, err = usersService.GetByUserName(userName)
		if err != nil {
			if err.Error() == "User not found" {
				_ = render(c, http.StatusNotFound, usersComponents.UsersNotFound())
				return nil
			} else {
				return err
			}
		}
	}

	if len(users) == 0 {
		_ = render(c, http.StatusNotFound, usersComponents.UsersNotFound())
		return nil
	}

	for i := range users {
		_ = render(c, http.StatusOK, usersComponents.UserRow(users[i]))
	}

	return nil
}

func GetUserByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	db, er := database.New()
	if er != nil {
		return er
	}

	usersService := services.NewUserService(db)

	user, err := usersService.GetByID(id)
	if err != nil {
		return err
	}

	_ = render(c, http.StatusOK, usersComponents.UserRow(*user))

	return nil
}

func UpdateUser(c echo.Context) error {
	idStr := c.FormValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	db, err := database.New()
	if err != nil {
		return err
	}

	usersService := services.NewUserService(db)

	user, err := usersService.GetByID(id)
	if err != nil {
		return err
	}

	user.Name = c.FormValue("name")

	// Check if the email is already in use
	oldEmail := c.FormValue("oldEmail")
	email := c.FormValue("email")
	if oldEmail != email {
		if _, err := usersService.GetByEmail(email); err == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Email already in use")
		}
	}
	user.Email = email

	if c.FormValue("password") != "" {
		user.Password = c.FormValue("password")
	}

	user.IsAdmin = c.FormValue("isAdmin") == "true"

	er := usersService.Update(user)
	if er != nil {
		return er
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func DeleteUser(c echo.Context) error {
	idStr := c.FormValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	db, er := database.New()
	if er != nil {
		return er
	}

	usersService := services.NewUserService(db)

	err = usersService.Delete(id)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func UpdateUserModalHandler(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	db, er := database.New()
	if er != nil {
		return er
	}

	usersService := services.NewUserService(db)

	user, err := usersService.GetByID(id)
	if err != nil {
		return err
	}

	_ = render(c, http.StatusOK, usersComponents.UpdateUserModal(*user))

	return nil
}

func DeleteUserModalHandler(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	db, er := database.New()
	if er != nil {
		return er
	}

	usersService := services.NewUserService(db)

	user, err := usersService.GetByID(id)
	if err != nil {
		return err
	}

	_ = render(c, http.StatusOK, usersComponents.DeleteUserModal(*user))

	return nil
}
