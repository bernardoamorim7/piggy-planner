package handlers

import (
	"net/http"

	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
	"piggy-planner/internal/services"
	securityLogsComponents "piggy-planner/web/components/security"

	"github.com/labstack/echo/v4"
)

func GetAllSecurityLogs(c echo.Context) error {
	db, err := database.New()
	if err != nil {
		return err
	}

	securityLogsService := services.NewSecurityLogsService(db)

	logs, err := securityLogsService.GetAll()
	if err != nil {
		return err
	}

	if len(logs) == 0 {
		_ = render(c, http.StatusNotFound, securityLogsComponents.NotFoundSecurityLogs())
		return nil
	}

	for i := range logs {
		_ = render(c, http.StatusOK, securityLogsComponents.Row(logs[i]))
	}

	return nil
}

func GetSecurityLogsByUserName(c echo.Context) error {
	userName := c.FormValue("search")

	db, er := database.New()
	if er != nil {
		return er
	}

	securityLogsService := services.NewSecurityLogsService(db)

	var (
		logs []models.SecurityLog
		err  error
	)

	if userName == "" {
		logs, err = securityLogsService.GetAll()
		if err != nil {
			return err
		}
	} else {
		logs, err = securityLogsService.GetByUserName(userName)
		if err != nil {
			if err.Error() == "Security log not found" {
				_ = render(c, http.StatusNotFound, securityLogsComponents.NotFoundSecurityLogs())
				return nil
			} else {
				return err
			}
		}
	}

	if len(logs) == 0 {
		_ = render(c, http.StatusNotFound, securityLogsComponents.NotFoundSecurityLogs())
		return nil
	}

	for i := range logs {
		_ = render(c, http.StatusOK, securityLogsComponents.Row(logs[i]))
	}

	return nil
}
