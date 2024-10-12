package handlers

import (
	"net/http"
	"strconv"
	"time"

	expensesComponents "piggy-planner/cmd/web/components/expenses"
	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
	"piggy-planner/internal/services"

	"github.com/labstack/echo/v4"
)

func CreateExpense(c echo.Context) error {
	userId := c.Get("userID").(uint64)

	amountStr := c.FormValue("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return err
	}

	description := c.FormValue("description")

	expenseTypeStr := c.FormValue("expenseType")
	expenseType, err := strconv.ParseUint(expenseTypeStr, 10, 64)
	if err != nil {
		return err
	}
	expenseTypeModel := models.ExpenseType{
		ID: expenseType,
	}

	dateStr := c.FormValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	db, err := database.New()
	if err != nil {
		return err
	}

	expensesService := services.NewExpensesService(db)

	expense := &models.Expense{
		UserID:      userId,
		Amount:      amount,
		Description: description,
		Type:        expenseTypeModel,
		Date:        date,
	}

	err = expensesService.Create(expense)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func GetAllExpenses(c echo.Context) error {
	userId := c.Get("userID").(uint64)

	db, err := database.New()
	if err != nil {
		return err
	}

	expensesService := services.NewExpensesService(db)

	expenses, err := expensesService.GetAll(userId)
	if err != nil {
		return err
	}

	if len(expenses) == 0 {
		_ = render(c, http.StatusNotFound, expensesComponents.NotFoundExpenses())
		return nil
	}

	for i := range expenses {
		_ = render(c, http.StatusOK, expensesComponents.ExpenseRow(expenses[i]))
	}

	return nil
}

func GetExpense(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	db, err := database.New()
	if err != nil {
		return err
	}

	expensesService := services.NewExpensesService(db)

	expense, err := expensesService.GetByID(id)
	if err != nil {
		return err
	}

	_ = render(c, http.StatusOK, expensesComponents.ExpenseRow(*expense))

	return nil
}

func GetExpensesByDescription(c echo.Context) error {
	description := c.FormValue("search")

	db, er := database.New()
	if er != nil {
		return er
	}

	expensesService := services.NewExpensesService(db)

	var (
		expenses []models.Expense
		err      error
	)

	if description == "" {
		userID := c.Get("userID").(uint64)
		expenses, err = expensesService.GetAll(userID)
		if err != nil {
			if err.Error() == "Expenses not found" {
				_ = render(c, http.StatusNotFound, expensesComponents.NotFoundExpenses())
				return nil
			} else {
				return err
			}
		}
	} else {
		expenses, err = expensesService.GetByDescription(description)
		if err != nil {
			if err.Error() == "Expenses not found" {
				_ = render(c, http.StatusNotFound, expensesComponents.NotFoundExpenses())
				return nil
			} else {
				return err
			}
		}
	}

	if len(expenses) == 0 {
		_ = render(c, http.StatusNotFound, expensesComponents.NotFoundExpenses())
		return nil
	}

	for i := range expenses {
		_ = render(c, http.StatusOK, expensesComponents.ExpenseRow(expenses[i]))
	}

	return nil
}

func UpdateExpense(c echo.Context) error {
	userID := c.Get("userID").(uint64)

	idStr := c.FormValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	amountStr := c.FormValue("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return err
	}

	description := c.FormValue("description")

	expenseTypeStr := c.FormValue("expenseType")
	expenseType, err := strconv.ParseUint(expenseTypeStr, 10, 64)
	if err != nil {
		return err
	}
	expenseTypeModel := models.ExpenseType{
		ID: expenseType,
	}

	dateStr := c.FormValue("date")

	var date time.Time
	if dateStr != "0000-00-00" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}
	} else {
		date = time.Time{}
	}

	db, err := database.New()
	if err != nil {
		return err
	}

	expensesService := services.NewExpensesService(db)

	expense := &models.Expense{
		ID:          id,
		UserID:      userID,
		Amount:      amount,
		Description: description,
		Type:        expenseTypeModel,
		Date:        date,
	}

	err = expensesService.Update(expense)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func DeleteExpense(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	db, err := database.New()
	if err != nil {
		return err
	}

	expensesService := services.NewExpensesService(db)

	err = expensesService.Delete(id)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func CreateExpenseType(c echo.Context) error {
	userId := c.Get("userID").(uint64)

	name := c.FormValue("name")

	db, err := database.New()
	if err != nil {
		return err
	}

	expensesTypesService := services.NewExpenseTypesService(db)

	expenseType := &models.ExpenseType{
		Name: name,
	}

	err = expensesTypesService.Create(userId, expenseType)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func GetAllExpenseTypes(c echo.Context) error {
	userId := c.Get("userID").(uint64)

	db, err := database.New()
	if err != nil {
		return err
	}

	expensesTypesService := services.NewExpenseTypesService(db)

	expenseTypes, err := expensesTypesService.GetAll(userId)
	if err != nil {
		return err
	}

	return render(c, http.StatusOK, expensesComponents.ExpenseTypesOptions(expenseTypes))
}

func GetExpenseType(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	db, err := database.New()
	if err != nil {
		return err
	}

	expensesTypesService := services.NewExpenseTypesService(db)

	expenseType, err := expensesTypesService.GetByID(id)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, expenseType)
}

func UpdateExpenseType(c echo.Context) error {
	expenseTypeIdStr := c.Param("id")
	expenseTypeId, err := strconv.ParseUint(expenseTypeIdStr, 10, 64)
	if err != nil {
		return err
	}

	name := c.FormValue("name")

	db, err := database.New()
	if err != nil {
		return err
	}

	expenseTypesService := services.NewExpenseTypesService(db)

	expenseType := &models.ExpenseType{
		ID:   expenseTypeId,
		Name: name,
	}

	err = expenseTypesService.Update(expenseType)

	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func DeleteExpenseType(c echo.Context) error {
	formParams, err := c.FormParams()
	if err != nil {
		return err
	}

	expenseTypeIdStr := formParams.Get("expenseID")

	expenseTypeId, err := strconv.ParseUint(expenseTypeIdStr, 10, 64)
	if err != nil {
		return err
	}

	db, err := database.New()
	if err != nil {
		return err
	}

	expensesTypesService := services.NewExpenseTypesService(db)

	err = expensesTypesService.Delete(expenseTypeId)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func CreateExpenseModalHandler(c echo.Context) error {
	return render(c, http.StatusOK, expensesComponents.CreateExpenseModal())
}

func UpdateExpenseModalHandler(c echo.Context) error {
	expenseIDStr := c.Param("id")
	expenseID, err := strconv.ParseUint(expenseIDStr, 10, 64)
	if err != nil {
		return err
	}

	db, err := database.New()
	if err != nil {
		return err
	}

	expensesService := services.NewExpensesService(db)

	expense, err := expensesService.GetByID(expenseID)
	if err != nil {
		return err
	}

	return render(c, http.StatusOK, expensesComponents.UpdateExpenseModal(*expense))
}

func DeleteExpenseModalHandler(c echo.Context) error {
	expenseIDStr := c.Param("id")
	expenseID, err := strconv.ParseUint(expenseIDStr, 10, 64)
	if err != nil {
		return err
	}

	return render(c, http.StatusOK, expensesComponents.DeleteExpenseModal(expenseID))
}

func CreateExpenseTypeModalHandler(c echo.Context) error {
	return render(c, http.StatusOK, expensesComponents.CreateExpenseTypeModal())
}
