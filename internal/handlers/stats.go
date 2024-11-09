package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"piggy-planner/internal/database"
	"piggy-planner/internal/services"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
)

func BalanceHandler(c echo.Context) error {
	userID := c.Get("userID").(uint64)

	var balance float64

	db, err := database.New()
	if err != nil {
		return err
	}

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

// 		db, err := database.New()
// if err != nil {
// 	return err
// }

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

	db, err := database.New()
	if err != nil {
		return err
	}

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

	db, err := database.New()
	if err != nil {
		return err
	}

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

	db, err := database.New()
	if err != nil {
		return err
	}

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

	db, err := database.New()
	if err != nil {
		return err
	}

	incomesService := services.NewIncomesService(db)

	incomes, err := incomesService.GetAll(userID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Map to hold unique income types and their associated amounts
	incomeMap := make(map[string][]float64)

	for _, income := range incomes {
		incomeMap[income.Type.Name] = append(incomeMap[income.Type.Name], income.Amount)
	}

	// Prepare labels and values for the chart
	var labels []string

	for label := range incomeMap {
		labels = append(labels, label)
	}

	// Sort the labels by alphabetical order so it's always the same order
	// when the chart is rendered
	sort.Strings(labels)
	sortedValues := make([]float64, len(labels))
	for i, label := range labels {
		// Sum the amounts for each expense type
		var total float64
		for _, amount := range incomeMap[label] {
			total += amount
		}
		sortedValues[i] = total
	}

	// Convert labels and values to JSON format for the chart
	labelsJSON, _ := json.Marshal(labels)
	valuesJSON, _ := json.Marshal(sortedValues)

	chart := fmt.Sprintf(`
    <canvas id="incomesChart" width="400" height="400"></canvas>
    <script>
        var ctx = document.getElementById('incomesChart').getContext('2d');
        var myChart = new Chart(ctx, {
            type: 'pie',
            data: {
                labels: %s,
                datasets: [{
                    label: 'Incomes',
                    data: %s,
                    backgroundColor: [
                        'rgba(255, 99, 132, 0.2)',
                        'rgba(54, 162, 235, 0.2)',
                        'rgba(255, 206, 86, 0.2)',
                        'rgba(50, 220, 220, 0.2)',
                        'rgba(153, 102, 255, 0.2)',
                        'rgba(255, 159, 64, 0.2)',
						'rgba(140, 160, 50, 0.2)',
						'rgba(188, 143, 143, 0.2)',
						'rgba(160, 160, 160, 0.2)'
                    ],
                    borderColor: [
                        'rgba(255, 99, 132, 1)',
                        'rgba(54, 162, 235, 1)',
                        'rgba(255, 206, 86, 1)',
                        'rgba(50, 220, 220, 1)',
                        'rgba(153, 102, 255, 1)',
                        'rgba(255, 159, 64, 1)',
      					'rgba(140, 180, 50, 1)',
						'rgba(188, 143, 143, 1)',
						'rgba(160, 160, 160, 1)'
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
`, string(labelsJSON), string(valuesJSON))

	return c.HTML(http.StatusOK, chart)
}

func ExpensesChartHandler(c echo.Context) error {
	userID := c.Get("userID").(uint64)

	db, err := database.New()
	if err != nil {
		return err
	}

	expensesService := services.NewExpensesService(db)

	expenses, err := expensesService.GetAll(userID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Map to hold unique income types and their associated amounts
	expenseMap := make(map[string][]float64)

	for _, expense := range expenses {
		expenseMap[expense.Type.Name] = append(expenseMap[expense.Type.Name], expense.Amount)
	}

	// Prepare labels and values for the chart
	var labels []string

	for label := range expenseMap {
		labels = append(labels, label)
	}

	// Sort the labels by alphabetical order so it's always the same order
	// when the chart is rendered
	sort.Strings(labels)
	sortedValues := make([]float64, len(labels))
	for i, label := range labels {
		// Sum the amounts for each expense type
		var total float64
		for _, amount := range expenseMap[label] {
			total += amount
		}
		sortedValues[i] = total
	}

	// Convert labels and values to JSON format for the chart
	labelsJSON, _ := json.Marshal(labels)
	valuesJSON, _ := json.Marshal(sortedValues)

	chart := fmt.Sprintf(`
    <canvas id="expensesChart" width="400" height="400"></canvas>
    <script>
        var ctx = document.getElementById('expensesChart').getContext('2d');
        var myChart = new Chart(ctx, {
            type: 'pie',
            data: {
                labels: %s,
                datasets: [{
                    label: 'Expenses',
                    data: %s,
                    backgroundColor: [
                        'rgba(255, 99, 132, 0.2)',
                        'rgba(54, 162, 235, 0.2)',
                        'rgba(255, 206, 86, 0.2)',
                        'rgba(50, 220, 220, 0.2)',
                        'rgba(153, 102, 255, 0.2)',
                        'rgba(255, 159, 64, 0.2)',
						'rgba(140, 160, 50, 0.2)',
						'rgba(188, 143, 143, 0.2)',
						'rgba(160, 160, 160, 0.2)'
                    ],
                    borderColor: [
                        'rgba(255, 99, 132, 1)',
                        'rgba(54, 162, 235, 1)',
                        'rgba(255, 206, 86, 1)',
                        'rgba(50, 220, 220, 1)',
                        'rgba(153, 102, 255, 1)',
                        'rgba(255, 159, 64, 1)',
      					'rgba(140, 180, 50, 1)',
						'rgba(188, 143, 143, 1)',
						'rgba(160, 160, 160, 1)'
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
`, string(labelsJSON), string(valuesJSON))

	return c.HTML(http.StatusOK, chart)
}
