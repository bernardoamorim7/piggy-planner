package middlewares

import (
	"net/http"
	"piggy-planner/internal/i18n"
	"strings"

	"github.com/invopop/ctxi18n"
	"github.com/labstack/echo/v4"
)

// I18NMiddleware is a middleware that sets the locale for the request.
func I18NMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			acceptLang := c.Request().Header.Get("Accept-Language")
			locale := "en"
			if acceptLang != "" {
				locale = strings.Split(acceptLang, ",")[0]
				// Extract the language code (e.g., "en-US" -> "en")
				locale = strings.Split(locale, "-")[0]
			}

			// Check if the locale is supported, otherwise default to "en"
			if !i18n.SupportedLocales[locale] {
				locale = "en"
			}

			// Create a new context with the locale
			localeCtx, err := ctxi18n.WithLocale(c.Request().Context(), locale)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			// Update the request with the new context
			req := c.Request().WithContext(localeCtx)
			c.SetRequest(req)

			return next(c)
		}
	}
}
