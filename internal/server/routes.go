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

	// 404
	e.RouteNotFound("/*", handlers.NotFoundHandler)

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
	e.GET("/settings", handlers.SettingsHandler, middlewares.Protected())
	e.GET("/incomes", handlers.IncomesHandler, middlewares.Protected())
	e.GET("/expenses", handlers.ExpensesHandler, middlewares.Protected())
	e.GET("/objectives", handlers.ObjectivesHandler, middlewares.Protected())

	// API
	api := e.Group("/api")

	// Incomes
	api.POST("/incomes/:userID", handlers.CreateIncome, middlewares.Protected())
	api.GET("/incomes/:userID", handlers.GetAllIncomes, middlewares.Protected())
	api.GET("/incomes/income/:id", handlers.GetIncome, middlewares.Protected())
	api.PUT("/incomes/:id", handlers.UpdateIncome, middlewares.Protected())
	api.DELETE("/incomes/:id", handlers.DeleteIncome, middlewares.Protected())

	// Incomes types
	api.POST("/incomes/types/:userID", handlers.CreateIncomeType, middlewares.Protected())
	api.GET("/incomes/types/:userID", handlers.GetAllIncomeTypes, middlewares.Protected())
	api.GET("/incomes/types/type/:id", handlers.GetIncomeType, middlewares.Protected())
	api.PUT("/incomes/types/:id", handlers.UpdateIncomeType, middlewares.Protected())
	api.DELETE("/incomes/types/:id", handlers.DeleteIncomeType, middlewares.Protected())

	return e
}
