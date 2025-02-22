package server

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"piggy-planner/web"

	"piggy-planner/internal/handlers"
	locales "piggy-planner/internal/i18n"
	"piggy-planner/internal/middlewares"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/invopop/ctxi18n"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())

	// Secret Key and Session Store
	secretKey := os.Getenv("PIGGY_SECRET")
	if secretKey == "" {
		secretKey = generateSecretKey()
		saveSecretKeyToEnv(secretKey)
	}
	sessionStore := sessions.NewCookieStore([]byte(secretKey))
	isProduction := os.Getenv("PIGGY_ENV") == "production"
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	}
	e.Use(session.Middleware(sessionStore))

	// FileServer
	fileServer := http.FileServer(http.FS(web.Files))

	// Custom Middleware
	e.Use(middlewares.GetSessionVars())
	e.Use(middlewares.RequestLogger())

	// i18n
	e.Use(middlewares.I18NMiddleware())
	if err := ctxi18n.Load(locales.LocaleFS); err != nil {
		e.Logger.Errorf("error loading locales: %v", err)
	}

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

	e.GET("/modals/users/update/:id", handlers.UpdateUserModalHandler, middlewares.Protected(), middlewares.AdminOnly())
	e.GET("/modals/users/delete/:id", handlers.DeleteUserModalHandler, middlewares.Protected(), middlewares.AdminOnly())

	// Admin Views
	e.GET("/users", handlers.UsersHandler, middlewares.Protected(), middlewares.AdminOnly())
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

	// Users
	api.GET("/users", handlers.GetAllUsers, middlewares.Protected(), middlewares.AdminOnly())
	api.POST("/users/search", handlers.GetUserByName, middlewares.Protected(), middlewares.AdminOnly())
	api.GET("/users/user/:id", handlers.GetUserByID, middlewares.Protected(), middlewares.AdminOnly())
	api.PUT("/users", handlers.UpdateUser, middlewares.Protected(), middlewares.AdminOnly())
	api.DELETE("/users", handlers.DeleteUser, middlewares.Protected(), middlewares.AdminOnly())

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
	api.GET("/stats/incomes-month-chart", handlers.IncomesPerMonthChartHandler, middlewares.Protected())
	api.GET("/stats/expenses-month-chart", handlers.ExpensesPerMonthChartHandler, middlewares.Protected())

	// Admin
	// Security
	api.GET("/security", handlers.GetAllSecurityLogs, middlewares.Protected(), middlewares.AdminOnly())
	api.POST("/security/search", handlers.GetSecurityLogsByUserName, middlewares.Protected(), middlewares.AdminOnly())

	// Requests
	api.GET("/requests", handlers.RequestLogsHandler, middlewares.Protected(), middlewares.AdminOnly())
	api.GET("/requests/history", handlers.RequestHistoryHandler, middlewares.Protected(), middlewares.AdminOnly())

	return e
}

func generateSecretKey() string {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Fatalf("Error generating secret key: %v", err)
	}
	return base64.StdEncoding.EncodeToString(key)
}

func saveSecretKeyToEnv(secretKey string) {
	file, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening .env file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("PIGGY_SECRET=\"%s\"\n", secretKey)); err != nil {
		log.Fatalf("Error writing to .env file: %v", err)
	}
}
