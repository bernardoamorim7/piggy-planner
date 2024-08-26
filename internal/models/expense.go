package models

type Expense struct {
	ID          uint64 `db:"id"`
	UserID      uint64 `db:"fk_user_id"`
	Amount      float64 `db:"amount"`
	Description string `db:"description"`
	Type        string `db:"fk_expense_type"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

func NewExpense(id uint64, userId uint64, amount float64, description string, expenseType string, createdAt string, updatedAt string) *Expense {
	return &Expense{
		ID:          id,
		UserID:      userId,
		Amount:      amount,
		Description: description,
		Type:        expenseType,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// ExpenseType represents the type of expense.
type ExpenseType struct {
	ID   uint64 `db:"id"`
	Name string `db:"name"`
}

func NewExpenseType(id uint64, name string) *ExpenseType {
	return &ExpenseType{
		ID:   id,
		Name: name,
	}
}
