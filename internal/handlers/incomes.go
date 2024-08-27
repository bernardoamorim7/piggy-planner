package handlers

import (
	"net/http"
	"strconv"
	"time"

	incomeComponents "piggy-planner/cmd/web/components/incomes"
	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
	"piggy-planner/internal/services"

	"github.com/labstack/echo/v4"
)

func CreateIncome(c echo.Context) error {
	userId := c.Get("userID").(uint64)

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

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func GetAllIncomes(c echo.Context) error {
	userId := c.Get("userID").(uint64)

	db := database.New()

	incomeService := services.NewIncomeService(db)

	incomes, err := incomeService.GetAll(userId)
	if err != nil {
		return err
	}

	if len(incomes) == 0 {
		return c.HTML(http.StatusOK, "No incomes found")
	}

	for i := range incomes {
		_ = render(c, http.StatusOK, incomeComponents.IncomeRow(incomes[i]))
	}

	return nil
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

	_ = render(c, http.StatusOK, incomeComponents.IncomeRow(*income))

	return nil
}

func GetIncomesByDescription(c echo.Context) error {
	description := c.FormValue("search")

	db := database.New()

	incomeService := services.NewIncomeService(db)

	var (
		incomes []models.Income
		err     error
	)

	if description == "" {
		userID := c.Get("userID").(uint64)
		incomes, err = incomeService.GetAll(userID)
		if err != nil {
			return err
		}
	} else {
		incomes, err = incomeService.GetByDescription(description)
		if err != nil {
			if err.Error() == "Income not found" {
				_ = render(c, http.StatusNotFound, incomeComponents.NotFoundIncomes())
				return nil
			} else {
				return err
			}
		}
	}

	if len(incomes) == 0 {
		_ = render(c, http.StatusNotFound, incomeComponents.NotFoundIncomes())
		return nil
	}

	for i := range incomes {
		_ = render(c, http.StatusOK, incomeComponents.IncomeRow(incomes[i]))
	}

	return nil
}

func UpdateIncome(c echo.Context) error {
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

	incomeTypeStr := c.FormValue("incomeType")
	incomeType, err := strconv.ParseUint(incomeTypeStr, 10, 64)
	if err != nil {
		return err
	}
	incomeTypeModel := models.IncomeType{
		ID: incomeType,
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

	db := database.New()

	incomeService := services.NewIncomeService(db)

	income := &models.Income{
		ID:          id,
		UserID:      userID,
		Amount:      amount,
		Description: description,
		Type:        incomeTypeModel,
		Date:        date,
	}

	err = incomeService.Update(income)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func DeleteIncome(c echo.Context) error {
	idStr := c.QueryParam("id")
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

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func CreateIncomeType(c echo.Context) error {
	userId := c.Get("userID").(uint64)

	name := c.FormValue("name")

	db := database.New()

	incomeTypeService := services.NewIncomeTypeService(db)

	incomeType := &models.IncomeType{
		Name: name,
	}

	err := incomeTypeService.Create(userId, incomeType)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func GetAllIncomeTypes(c echo.Context) error {
	userId := c.Get("userID").(uint64)

	db := database.New()

	incomeTypeService := services.NewIncomeTypeService(db)

	incomeTypes, err := incomeTypeService.GetAll(userId)
	if err != nil {
		return err
	}

	return render(c, http.StatusOK, incomeComponents.IncomeTypesOptions(incomeTypes))
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

	c.Response().Header().Set("Content-Type", "application/json")
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

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func DeleteIncomeType(c echo.Context) error {
	formParams, err := c.FormParams()
	if err != nil {
		return err
	}

	incomeTypeIdStr := formParams.Get("incomeID")

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

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func CreateIncomeModalHandler(c echo.Context) error {
	return render(c, http.StatusOK, incomeComponents.CreateIncomeModal())
}

func UpdateIncomeModalHandler(c echo.Context) error {
	incomeIDStr := c.Param("id")
	incomeID, err := strconv.ParseUint(incomeIDStr, 10, 64)
	if err != nil {
		return err
	}

	db := database.New()

	incomeService := services.NewIncomeService(db)

	income, err := incomeService.GetByID(incomeID)
	if err != nil {
		return err
	}

	return render(c, http.StatusOK, incomeComponents.UpdateIncomeModal(*income))
}

func DeleteIncomeModalHandler(c echo.Context) error {
	incomeIDStr := c.Param("id")
	incomeID, err := strconv.ParseUint(incomeIDStr, 10, 64)
	if err != nil {
		return err
	}

	return render(c, http.StatusOK, incomeComponents.DeleteIncomeModal(incomeID))
}

func CreateIncomeTypeModalHandler(c echo.Context) error {
	return render(c, http.StatusOK, incomeComponents.CreateIncomeTypeModal())
}