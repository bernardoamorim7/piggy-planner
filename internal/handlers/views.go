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

func SettingsHandler(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return render(c, http.StatusOK, views.Settings())
	} else {
		return render(c, http.StatusOK, web.Base())
	}
}

func NotFoundHandler(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return render(c, http.StatusNotFound, web.NotFound())
	} else {
		return render(c, http.StatusNotFound, web.Base())
	}
}

func IncomesHandler(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return render(c, http.StatusOK, views.Incomes())
	} else {
		return render(c, http.StatusOK, web.Base())
	}
}

func ExpensesHandler(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return render(c, http.StatusOK, views.Expenses())
	} else {
		return render(c, http.StatusOK, web.Base())
	}
}

func ObjectivesHandler(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "" {
		return render(c, http.StatusOK, views.Objectives())
	} else {
		return render(c, http.StatusOK, web.Base())
	}
}
