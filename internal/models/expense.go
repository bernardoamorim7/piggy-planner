package models

import (
	"errors"
	"time"
)

type Expense struct {
	ID          uint64  `db:"id"`
	UserID      uint64  `db:"fk_user_id"`
	Amount      float64 `db:"amount"`
	Description string  `db:"description"`
	Type        ExpenseType
	Date        time.Time `db:"date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func NewExpense(id uint64, userId uint64, amount float64, description string, expenseType ExpenseType, date time.Time, createdAt time.Time, updatedAt time.Time) *Expense {
	return &Expense{
		ID:          id,
		UserID:      userId,
		Amount:      amount,
		Description: description,
		Type:        expenseType,
		Date:        date,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// Validate validates the expense.
func (e *Expense) Validate() error {
	if e.Amount == 0 {
		return errors.New("Amount must be greater than 0")
	}

	if len(e.Description) == 0 {
		return errors.New("Description cannot be empty")
	}

	if e.Type.ID == 0 {
		return errors.New("Type cannot be empty")
	}

	return nil
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

// Validate validates the expense type.
func (et *ExpenseType) Validate() error {
	if len(et.Name) == 0 {
		return errors.New("Name cannot be empty")
	}

	return nil
}
