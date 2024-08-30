package handlers

import (
	"fmt"
	"net/http"
	"piggy-planner/internal/database"
	"piggy-planner/internal/services"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func BalanceHandler(c echo.Context) error {
	userID := c.Get("userID").(uint64)

	var balance float64

	db := database.New()

	incomesService := services.NewIncomesService(db)

	incomes, err := incomesService.GetAll(userID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	for _, income := range incomes {
		balance += income.Amount
	}

	expensesService := services.NewExpensesService(db)

	expenses, err := expensesService.GetAll(userID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	for _, expense := range expenses {
		balance -= expense.Amount
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("%.2f", balance))
}

// func DebtHandler(c echo.Context) error {
// 	userID := c.Get("userID").(uint64)

// 	var debt float64

// 	db := database.New()

// 	debtsService := services.NewDebtsService(db)

// 	debts, err := debtsService.GetAll(userID)
// 	if err != nil {
// 		return c.NoContent(http.StatusInternalServerError)
// 	}

// 	for _, debt := range debts {
// 		debt += debt.Amount
// 	}

// 	return c.HTML(http.StatusOK, fmt.Sprintf("%.2f", debt))
// }

func TotalExpensesHandler(c echo.Context) error {
	userID := c.Get("userID").(uint64)

	db := database.New()

	expensesService := services.NewExpensesService(db)

	expenses, err := expensesService.GetAll(userID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	var expensesAmount float64
	for _, expense := range expenses {
		expensesAmount += expense.Amount
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("%.2f", expensesAmount))
}

func CurrentMonthIncomesHandler(c echo.Context) error {
	userID := c.Get("userID").(uint64)

	now := time.Now()

	// Calculate the first day of the current month
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Calculate the last day of the current month
	endDate := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 0, -1)

	db := database.New()

	incomesService := services.NewIncomesService(db)

	incomes, err := incomesService.GetByPeriod(userID, startDate, endDate)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	var incomesAmount float64
	for _, income := range incomes {
		incomesAmount += income.Amount
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("%.2f", incomesAmount))
}

func CurrentMonthExpensesHandler(c echo.Context) error {
	userID := c.Get("userID").(uint64)

	now := time.Now()

	// Calculate the first day of the current month
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Calculate the last day of the current month
	endDate := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 0, -1)

	db := database.New()

	expensesService := services.NewExpensesService(db)

	expenses, err := expensesService.GetByPeriod(userID, startDate, endDate)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	var expensesAmount float64
	for _, expense := range expenses {
		expensesAmount += expense.Amount
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("%.2f", expensesAmount))
}

func IncomesChartHandler(c echo.Context) error {
	userID := c.Get("userID").(uint64)

	db := database.New()

	incomesService := services.NewIncomesService(db)

	incomes, err := incomesService.GetAll(userID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	var labels []string
	var values []float64

	for _, income := range incomes {
		labels = append(labels, income.Type.Name)
		values = append(values, income.Amount)
	}

	labelsArr := strings.ReplaceAll(fmt.Sprintf("%q", labels), "\" \"", "\", \"")

	var valuesStr []string
	for _, v := range values {
		valuesStr = append(valuesStr, fmt.Sprintf("%.2f", v))
	}
	valuesArr := strings.ReplaceAll(fmt.Sprintf("%q", valuesStr), "\" \"", "\",\"")

	chart := fmt.Sprintf(`
		    <canvas id="incomesChart" width="400" height="400"></canvas>
            <script>
                var ctx = document.getElementById('incomesChart').getContext('2d');
                var myChart = new Chart(ctx, {
                    type: 'pie',
                    data: {
                        labels: %v,
                        datasets: [{
                            label: 'Incomes',
                            data: %v,
                            backgroundColor: [
                                'rgba(255, 99, 132, 0.2)',
                                'rgba(54, 162, 235, 0.2)',
                                'rgba(255, 206, 86, 0.2)',
                                'rgba(75, 192, 192, 0.2)',
                                'rgba(153, 102, 255, 0.2)',
                                'rgba(255, 159, 64, 0.2)'
                            ],
                            borderColor: [
                                'rgba(255, 99, 132, 1)',
                                'rgba(54, 162, 235, 1)',
                                'rgba(255, 206, 86, 1)',
                                'rgba(75, 192, 192, 1)',
                                'rgba(153, 102, 255, 1)',
                                'rgba(255, 159, 64, 1)'
                            ],
                            borderWidth: 1
                        }]
                    },
                    options: {
                        responsive: true,
                        plugins: {
                            legend: {
                                position: 'top',
                            },
                        }
                    }
                });
            </script>
        </body>
        </html>
    `, labelsArr, valuesArr)

	return c.HTML(http.StatusOK, chart)
}

func ExpensesChartHandler(c echo.Context) error {
	userID := c.Get("userID").(uint64)

	db := database.New()

	expensesService := services.NewExpensesService(db)

	expenses, err := expensesService.GetAll(userID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	var labels []string
	var values []float64

	for _, expense := range expenses {
		labels = append(labels, expense.Type.Name)
		values = append(values, expense.Amount)
	}

	labelsArr := strings.ReplaceAll(fmt.Sprintf("%q", labels), "\" \"", "\", \"")

	var valuesStr []string
	for _, v := range values {
		valuesStr = append(valuesStr, fmt.Sprintf("%.2f", v))
	}
	valuesArr := strings.ReplaceAll(fmt.Sprintf("%q", valuesStr), "\" \"", "\",\"")

	chart := fmt.Sprintf(`
		    <canvas id="expensesChart" width="400" height="400"></canvas>
            <script>
                var ctx = document.getElementById('expensesChart').getContext('2d');
                var myChart = new Chart(ctx, {
                    type: 'pie',
                    data: {
                        labels: %v,
                        datasets: [{
                            label: 'Expenses',
                            data: %v,
                            backgroundColor: [
                                'rgba(255, 99, 132, 0.2)',
                                'rgba(54, 162, 235, 0.2)',
                                'rgba(255, 206, 86, 0.2)',
                                'rgba(75, 192, 192, 0.2)',
                                'rgba(153, 102, 255, 0.2)',
                                'rgba(255, 159, 64, 0.2)'
                            ],
                            borderColor: [
                                'rgba(255, 99, 132, 1)',
                                'rgba(54, 162, 235, 1)',
                                'rgba(255, 206, 86, 1)',
                                'rgba(75, 192, 192, 1)',
                                'rgba(153, 102, 255, 1)',
                                'rgba(255, 159, 64, 1)'
                            ],
                            borderWidth: 1
                        }]
                    },
                    options: {
                        responsive: true,
                        plugins: {
                            legend: {
                                position: 'top',
                            },
                        }
                    }
                });
            </script>
        </body>
        </html>
    `, labelsArr, valuesArr)

	return c.HTML(http.StatusOK, chart)
}
