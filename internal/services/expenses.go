package services

import (
	"errors"
	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
	"time"
)

type ExpensesService interface {
	// Create a expense in the database
	Create(expense *models.Expense) error
	// Get all expenses from the database for a specific user
	GetAll(fkUserId uint64) ([]models.Expense, error)
	// Get a expense by ID from the database
	GetByID(id uint64) (*models.Expense, error)
	// Get a expense by description from the database
	GetByDescription(description string) ([]models.Expense, error)
	// Update a expense in the database
	Update(expense *models.Expense) error
	// Delete a expense by ID from the database
	Delete(id uint64) error
}

type expensesService struct {
	DB database.Service
}

func NewExpensesService(db database.Service) ExpensesService {
	return &expensesService{
		DB: db,
	}
}

func (s *expensesService) Create(expense *models.Expense) error {
	if expense.Amount == 0 {
		return errors.New("Missing expense amount")
	}

	if expense.Description == "" {
		return errors.New("Missing expense description")
	}

	if expense.UserID == 0 {
		return errors.New("Missing expense user ID")
	}

	query := "INSERT INTO expenses (fk_user_id, amount, description, date, fk_expense_type_id) VALUES (?, ?, ?, ?, ?)"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(expense.UserID, expense.Amount, expense.Description, expense.Date, expense.Type.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *expensesService) GetAll(fkUserId uint64) ([]models.Expense, error) {
	query := `
    SELECT 
        expenses.id, 
        expenses.fk_user_id, 
        expenses.amount, 
        expenses.description, 
        expense_types.id AS expense_type_id, 
        expense_types.name AS expense_type_name, 
        expenses.date 
    FROM 
        expenses 
    INNER JOIN 
        expense_types 
    ON 
        expenses.fk_expense_type_id = expense_types.id 
    WHERE 
        expenses.fk_user_id = ? 
    ORDER BY 
        expenses.date DESC
	`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(fkUserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		var expenseType models.ExpenseType
		var date []byte

		err := rows.Scan(&expense.ID, &expense.UserID, &expense.Amount, &expense.Description, &expenseType.ID, &expenseType.Name, &date)
		if err != nil {
			return nil, err
		}

		if string(date) != "0000-00-00" {
			expense.Date, err = time.Parse("2006-01-02", string(date))
			if err != nil {
				return nil, err
			}
		} else {
			expense.Date = time.Time{}
		}

		expense.Type = expenseType

		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func (s *expensesService) GetByID(id uint64) (*models.Expense, error) {
	query := `SELECT 
						expenses.id, 
						expenses.fk_user_id, 
						expenses.amount, 
						expenses.description, 
						expense_types.id AS expense_type_id, 
						expense_types.name AS expense_type_name, 
						expenses.date 
					FROM 
						expenses 
					INNER JOIN 
						expense_types 
					ON 
						expenses.fk_expense_type_id = expense_types.id 
					WHERE 
						expenses.id = ? 
					ORDER BY 
						expenses.date DESC
					LIMIT 1`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	var expense models.Expense
	var expenseType models.ExpenseType
	var date []byte
	err = row.Scan(&expense.ID, &expense.UserID, &expense.Amount, &expense.Description, &expenseType.ID, &expenseType.Name, &date)
	if err != nil {
		return nil, err
	}

	if string(date) != "0000-00-00" {
		expense.Date, err = time.Parse("2006-01-02", string(date))
		if err != nil {
			return nil, err
		}
	} else {
		expense.Date = time.Time{}
	}

	expense.Type = expenseType

	return &expense, nil
}

func (s *expensesService) GetByDescription(description string) ([]models.Expense, error) {
	query := `SELECT 
						expenses.id, 
						expenses.fk_user_id, 
						expenses.amount, 
						expenses.description, 
						expense_types.id AS expense_type_id, 
						expense_types.name AS expense_type_name, 
						expenses.date 
					FROM 
						expenses 
					INNER JOIN 
						expense_types 
					ON 
						expenses.fk_expense_type_id = expense_types.id 
					WHERE 
						expenses.description LIKE ?
					ORDER BY 
						expenses.date DESC
					`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	description = "%" + description + "%"
	row, err := stmt.Query(description)
	if err != nil {
		return nil, err
	}

	var expenses []models.Expense
	for row.Next() {
		var expense models.Expense
		var expenseType models.ExpenseType
		var date []byte

		err := row.Scan(&expense.ID, &expense.UserID, &expense.Amount, &expense.Description, &expenseType.ID, &expenseType.Name, &date)
		if err != nil {
			return nil, err
		}

		if string(date) != "0000-00-00" {
			expense.Date, err = time.Parse("2006-01-02", string(date))
			if err != nil {
				return nil, err
			}
		} else {
			expense.Date = time.Time{}
		}

		expense.Type = expenseType

		expenses = append(expenses, expense)
	}

	if len(expenses) == 0 {
		return nil, errors.New("Expenses not found")
	}

	return expenses, nil
}

func (s *expensesService) Update(expense *models.Expense) error {
	if expense.Amount == 0 {
		return errors.New("Missing expense amount")
	}

	if expense.Description == "" {
		return errors.New("Missing expense description")
	}

	if expense.Type.ID == 0 {
		return errors.New("Missing expense type ID")
	}

	if expense.UserID == 0 {
		return errors.New("Missing expense user ID")
	}

	query := "UPDATE expenses SET amount = ?, description = ?, fk_expense_type_id = ?, date = ? WHERE id = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(expense.Amount, expense.Description, expense.Type.ID, expense.Date, expense.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *expensesService) Delete(id uint64) error {
	query := "DELETE FROM expenses WHERE id = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

// Expense represents a service for managing expense types.
type ExpenseTypesService interface {
	// Create a Expenses type in the database
	Create(fkUserId uint64, expenseType *models.ExpenseType) error
	// Get all expense types from the database
	GetAll(fkUserID uint64) ([]models.ExpenseType, error)
	// Get a expense type by ID from the database
	GetByID(id uint64) (*models.ExpenseType, error)
	// Update a expense type in the database
	Update(expenseType *models.ExpenseType) error
	// Delete a expense type by ID from the database
	Delete(id uint64) error
}

type expenseTypesService struct {
	DB database.Service
}

func NewExpenseTypesService(db database.Service) ExpenseTypesService {
	return &expenseTypesService{
		DB: db,
	}
}

func (s *expenseTypesService) Create(fkUserId uint64, expenseType *models.ExpenseType) error {
	if expenseType.Name == "" {
		return errors.New("Missing expense type name")
	}

	query := "INSERT INTO expense_types (name, fk_user_id) VALUES (?, ?)"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(expenseType.Name, fkUserId)
	if err != nil {
		return err
	}

	return nil
}

func (s *expenseTypesService) GetAll(userID uint64) ([]models.ExpenseType, error) {
	query := "SELECT id, name FROM expense_types where fk_user_id = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	var expenseTypes []models.ExpenseType
	for rows.Next() {
		var expenseType models.ExpenseType
		err := rows.Scan(&expenseType.ID, &expenseType.Name)
		if err != nil {
			return nil, err
		}

		expenseTypes = append(expenseTypes, expenseType)
	}

	return expenseTypes, nil
}

func (s *expenseTypesService) GetByID(id uint64) (*models.ExpenseType, error) {
	query := "SELECT id, name FROM expense_types WHERE id = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	var expenseType models.ExpenseType
	err = row.Scan(&expenseType.ID, &expenseType.Name)
	if err != nil {
		return nil, err
	}

	return &expenseType, nil
}

func (s *expenseTypesService) Update(expenseType *models.ExpenseType) error {
	if expenseType.Name == "" {
		return errors.New("Missing expense type name")
	}

	query := "UPDATE expense_types SET name = ? WHERE id = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(expenseType.Name, expenseType.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *expenseTypesService) Delete(id uint64) error {
	query := "DELETE FROM expense_types WHERE id = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
