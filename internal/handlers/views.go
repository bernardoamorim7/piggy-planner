package handlers

import (
	"net/http"
	"piggy-planner/cmd/web"
	"piggy-planner/cmd/web/views"

	"github.com/labstack/echo/v4"
)

func DashboardHandler(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return render(c, http.StatusOK, views.Dashboard())
	} else {
		return render(c, http.StatusOK, web.Base())
	}
}

func ProfileHandler(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return render(c, http.StatusOK, views.Profile())
	} else {
		return render(c, http.StatusOK, web.Base())
	}
}