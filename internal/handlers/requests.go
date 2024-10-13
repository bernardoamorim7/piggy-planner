package handlers

import (
	"net/http"

	"piggy-planner/internal/middlewares"
	requestComponents "piggy-planner/web/components/requests"

	"github.com/labstack/echo/v4"
)

func RequestLogsHandler(c echo.Context) error {
	requestLogs := middlewares.GetRequestLogs()

	return render(c, http.StatusOK, requestComponents.RequestLogRows(requestLogs))
}
