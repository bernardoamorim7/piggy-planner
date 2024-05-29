package server

import (
	"net/http"
	"os"

	"piggy-planner/cmd/web"
	"piggy-planner/cmd/web/views"

	"piggy-planner/internal/server/handlers"
	"piggy-planner/internal/utils"

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

	// Static assets
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	// Health check
	e.GET("/health", s.healthHandler, utils.IsLocalhost)

	// Index
	e.GET("/", echo.WrapHandler(templ.Handler(web.Base())))

	// Auth
	e.GET("/login", echo.WrapHandler(templ.Handler(web.Login())))
	e.POST("/login", handlers.Login)
	e.GET("/register", echo.WrapHandler(templ.Handler(web.Register())))
	e.POST("/register", handlers.Register)

	// Dashboard
	e.GET("/dashboard", echo.WrapHandler(templ.Handler(views.Dashboard())), handlers.Protected())

	return e
}
