package models

import (
	"errors"
	"time"
)

type Income struct {
	ID          uint64  `db:"id"`
	UserID      uint64  `db:"fk_user_id"`
	Amount      float64 `db:"amount"`
	Description string  `db:"description"`
	Type        IncomeType
	Date        time.Time `db:"date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func NewIncome(id uint64, userId uint64, amount float64, description string, incomeType IncomeType, date time.Time, createdAt time.Time, updatedAt time.Time) *Income {
	return &Income{
		ID:          id,
		UserID:      userId,
		Amount:      amount,
		Description: description,
		Date:        date,
		Type:        incomeType,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// Validate validates the income.
func (r *Income) Validate() error {
	if r.Amount == 0 {
		return errors.New("Amount must be greater than 0")
	}

	if len(r.Description) == 0 {
		return errors.New("Description cannot be empty")
	}

	if r.Type.ID == 0 {
		return errors.New("Type cannot be empty")
	}

	return nil
}

// IncomeType represents the type of income.
type IncomeType struct {
	ID   uint64 `db:"id"`
	Name string `db:"name"`
}

func NewIncomeType(id uint64, name string) *IncomeType {
	return &IncomeType{
		ID:   id,
		Name: name,
	}
}

// Validate validates the income type.
func (rt *IncomeType) Validate() error {
	if len(rt.Name) == 0 {
		return errors.New("Name cannot be empty")
	}

	return nil
}
