package models

import "errors"

type Objective struct {
	ID          uint64 `db:"id"`
	UserID      uint64 `db:"fk_user_id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Deadline    string `db:"deadline"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

func NewObjective(id uint64, userId uint64, title string, description string, deadline string, createdAt string, updatedAt string) *Objective {
	return &Objective{
		ID:          id,
		UserID:      userId,
		Title:       title,
		Description: description,
		Deadline:    deadline,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// Validate validates the objective.
func (o *Objective) Validate() error {
	if len(o.Title) == 0 {
		return errors.New("Title cannot be empty")
	}

	if len(o.Description) == 0 {
		return errors.New("Description cannot be empty")
	}

	if len(o.Deadline) == 0 {
		return errors.New("Deadline cannot be empty")
	}

	return nil
}
