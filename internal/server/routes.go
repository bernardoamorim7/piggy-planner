package server

import (
	"net/http"
	"os"

	"piggy-planner/web"

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
	//e.GET("/settings", handlers.SettingsHandler, middlewares.Protected())
	e.GET("/incomes", handlers.IncomesHandler, middlewares.Protected())
	e.GET("/expenses", handlers.ExpensesHandler, middlewares.Protected())
	//e.GET("/objectives", handlers.ObjectivesHandler, middlewares.Protected())

	// Modals
	e.GET("/modals/incomes/create", handlers.CreateIncomeModalHandler, middlewares.Protected())
	e.GET("/modals/incomes/update/:id", handlers.UpdateIncomeModalHandler, middlewares.Protected())
	e.GET("/modals/incomes/delete/:id", handlers.DeleteIncomeModalHandler, middlewares.Protected())

	e.GET("/modals/incomes/types/create", handlers.CreateIncomeTypeModalHandler, middlewares.Protected())

	e.GET("/modals/expenses/create", handlers.CreateExpenseModalHandler, middlewares.Protected())
	e.GET("/modals/expenses/update/:id", handlers.UpdateExpenseModalHandler, middlewares.Protected())
	e.GET("/modals/expenses/delete/:id", handlers.DeleteExpenseModalHandler, middlewares.Protected())

	e.GET("/modals/expenses/types/create", handlers.CreateExpenseTypeModalHandler, middlewares.Protected())

	// Admin Views
	e.GET("/security", handlers.SecurityHandler, middlewares.Protected(), middlewares.AdminOnly())
	e.GET("/database", handlers.DatabaseHandler, middlewares.Protected(), middlewares.AdminOnly())
	e.GET("/requests", handlers.RequestsHandler, middlewares.Protected(), middlewares.AdminOnly())

	// API
	api := e.Group("/api")

	// Incomes
	api.POST("/incomes", handlers.CreateIncome, middlewares.Protected())
	api.GET("/incomes", handlers.GetAllIncomes, middlewares.Protected())
	api.GET("/incomes/income/:id", handlers.GetIncome, middlewares.Protected())
	api.POST("/incomes/search", handlers.GetIncomesByDescription, middlewares.Protected())
	api.PUT("/incomes", handlers.UpdateIncome, middlewares.Protected())
	api.DELETE("/incomes", handlers.DeleteIncome, middlewares.Protected())

	// Incomes types
	api.POST("/incomes/types", handlers.CreateIncomeType, middlewares.Protected())
	api.GET("/incomes/types", handlers.GetAllIncomeTypes, middlewares.Protected())
	api.GET("/incomes/types/type/:id", handlers.GetIncomeType, middlewares.Protected())
	api.PUT("/incomes/types", handlers.UpdateIncomeType, middlewares.Protected())
	api.DELETE("/incomes/types", handlers.DeleteIncomeType, middlewares.Protected())

	// Expenses
	api.POST("/expenses", handlers.CreateExpense, middlewares.Protected())
	api.GET("/expenses", handlers.GetAllExpenses, middlewares.Protected())
	api.GET("/expenses/expense/:id", handlers.GetExpense, middlewares.Protected())
	api.POST("/expenses/search", handlers.GetExpensesByDescription, middlewares.Protected())
	api.PUT("/expenses", handlers.UpdateExpense, middlewares.Protected())
	api.DELETE("/expenses", handlers.DeleteExpense, middlewares.Protected())

	// Expenses types
	api.POST("/expenses/types", handlers.CreateExpenseType, middlewares.Protected())
	api.GET("/expenses/types", handlers.GetAllExpenseTypes, middlewares.Protected())
	api.GET("/expenses/types/type/:id", handlers.GetExpenseType, middlewares.Protected())
	api.PUT("/expenses/types", handlers.UpdateExpenseType, middlewares.Protected())
	api.DELETE("/expenses/types", handlers.DeleteExpenseType, middlewares.Protected())

	// Objectives
	//api.POST("/objectives", handlers.CreateObjective, middlewares.Protected())
	//api.GET("/objectives", handlers.GetAllObjectives, middlewares.Protected())
	//api.GET("/objectives/objective/:id", handlers.GetObjective, middlewares.Protected())
	//api.PUT("/objectives", handlers.UpdateObjective, middlewares.Protected())
	//api.DELETE("/objectives", handlers.DeleteObjective, middlewares.Protected())

	// Objectives types
	//api.POST("/objectives/types", handlers.CreateObjectiveType, middlewares.Protected())
	//api.GET("/objectives/types", handlers.GetAllObjectiveTypes, middlewares.Protected())
	//api.GET("/objectives/types/type/:id", handlers.GetObjectiveType, middlewares.Protected())
	//api.PUT("/objectives/types", handlers.UpdateObjectiveType, middlewares.Protected())
	//api.DELETE("/objectives/types", handlers.DeleteObjectiveType, middlewares.Protected())

	// Dasbhboard Stats
	// Stat cards
	api.GET("/stats/balance", handlers.BalanceHandler, middlewares.Protected())
	//api.GET("/stats/debt", handlers.DebtHandler, middlewares.Protected())
	api.GET("/stats/total-expenses", handlers.TotalExpensesHandler, middlewares.Protected())
	api.GET("/stats/current-month-incomes", handlers.CurrentMonthIncomesHandler, middlewares.Protected())
	api.GET("/stats/current-month-expenses", handlers.CurrentMonthExpensesHandler, middlewares.Protected())

	// Charts
	api.GET("/stats/incomes-chart", handlers.IncomesChartHandler, middlewares.Protected())
	api.GET("/stats/expenses-chart", handlers.ExpensesChartHandler, middlewares.Protected())

	// Admin
	// Security
	api.GET("/security", handlers.GetAllSecurityLogs, middlewares.Protected(), middlewares.AdminOnly())
	api.POST("/security/search", handlers.GetSecurityLogsByUserName, middlewares.Protected(), middlewares.AdminOnly())

	return e
}
