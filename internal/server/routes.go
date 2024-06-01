package server

import (
	"net/http"
	"os"

	"piggy-planner/cmd/web"

	"piggy-planner/internal/handlers"
	"piggy-planner/internal/middlewares"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SECRET")))))
	fileServer := http.FileServer(http.FS(web.Files))
	e.Use(middlewares.GetSessionVars())

	// Static assets
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	// Health check
	e.GET("/health", s.healthHandler, middlewares.IsLocalhost())

	// Auth
	e.GET("/login", echo.WrapHandler(templ.Handler(web.Login())), middlewares.RedirectIfLoggedIn())
	e.POST("/login", handlers.Login)
	e.GET("/register", echo.WrapHandler(templ.Handler(web.Register())), middlewares.RedirectIfLoggedIn())
	e.POST("/register", handlers.Register)
	e.POST("/logout", handlers.Logout)

	// Index
	e.GET("/", handlers.DashboardHandler, middlewares.Protected())

	// Views
	e.GET("/profile", handlers.ProfileHandler, middlewares.Protected())

	return e
}
