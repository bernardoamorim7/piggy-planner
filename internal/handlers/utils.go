package handlers

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// This custom Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
func render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	userID := ctx.Get("userID")
	name := ctx.Get("name")
	avatar := ctx.Get("avatar")

	// Create a new context that includes the values
	newCtx := context.WithValue(ctx.Request().Context(), "userID", userID)
	newCtx = context.WithValue(newCtx, "name", name)
	newCtx = context.WithValue(newCtx, "avatar", avatar)

	if err := t.Render(newCtx, buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
