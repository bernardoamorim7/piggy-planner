package handlers

import (
	"net/http"
	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
	"piggy-planner/internal/services"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func CreateIncome(c echo.Context) error {
	userIdStr := c.Param("userID")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return err
	}

	amountStr := c.FormValue("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return err
	}

	description := c.FormValue("description")

	incomeTypeStr := c.FormValue("incomeType")
	incomeType, err := strconv.ParseUint(incomeTypeStr, 10, 64)
	if err != nil {
		return err
	}
	incomeTypeModel := models.IncomeType{
		ID: incomeType,
	}

	dateStr := c.FormValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	db := database.New()

	incomeService := services.NewIncomeService(db)

	income := &models.Income{
		UserID:      userId,
		Amount:      amount,
		Description: description,
		Type:        incomeTypeModel,
		Date:        date,
	}

	err = incomeService.Create(income)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func GetAllIncomes(c echo.Context) error {
	userIdStr := c.Param("userID")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return err
	}

	db := database.New()

	incomeService := services.NewIncomeService(db)

	incomes, err := incomeService.GetAll(userId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, incomes)
}

func GetIncome(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	db := database.New()

	incomeService := services.NewIncomeService(db)

	income, err := incomeService.GetByID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, income)
}

func UpdateIncome(c echo.Context) error {
	idStr := c.Param("id")
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

	incomeTypeStr := c.FormValue("incomeType")
	incomeType, err := strconv.ParseUint(incomeTypeStr, 10, 64)
	if err != nil {
		return err
	}
	incomeTypeModel := models.IncomeType{
		ID: incomeType,
	}

	dateStr := c.FormValue("date")
	date, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return err
	}

	db := database.New()

	incomeService := services.NewIncomeService(db)

	income := &models.Income{
		ID:          id,
		Amount:      amount,
		Description: description,
		Type:        incomeTypeModel,
		Date:        date,
	}

	err = incomeService.Update(income)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func DeleteIncome(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	db := database.New()

	incomeService := services.NewIncomeService(db)

	err = incomeService.Delete(id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func CreateIncomeType(c echo.Context) error {
	userIdStr := c.Param("userID")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return err
	}

	name := c.FormValue("name")

	db := database.New()

	incomeTypeService := services.NewIncomeTypeService(db)

	incomeType := &models.IncomeType{
		Name: name,
	}

	err = incomeTypeService.Create(userId, incomeType)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func GetAllIncomeTypes(c echo.Context) error {
	userId := c.Param("userID")
	userIdInt, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return err
	}

	db := database.New()

	incomeTypeService := services.NewIncomeTypeService(db)

	incomeTypes, err := incomeTypeService.GetAll(userIdInt)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, incomeTypes)
}

func GetIncomeType(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	db := database.New()

	incomeTypeService := services.NewIncomeTypeService(db)

	incomeType, err := incomeTypeService.GetByID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, incomeType)
}

func UpdateIncomeType(c echo.Context) error {
	incomeTypeIdStr := c.Param("id")
	incomeTypeId, err := strconv.ParseUint(incomeTypeIdStr, 10, 64)
	if err != nil {
		return err
	}

	name := c.FormValue("name")

	db := database.New()

	incomeTypeService := services.NewIncomeTypeService(db)

	incomeType := &models.IncomeType{
		ID:   incomeTypeId,
		Name: name,
	}

	err = incomeTypeService.Update(incomeType)

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func DeleteIncomeType(c echo.Context) error {
	incomeTypeIdStr := c.Param("id")
	incomeTypeId, err := strconv.ParseUint(incomeTypeIdStr, 10, 64)
	if err != nil {
		return err
	}

	db := database.New()

	incomeTypeService := services.NewIncomeTypeService(db)

	err = incomeTypeService.Delete(incomeTypeId)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
