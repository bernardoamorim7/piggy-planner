package services

import (
	"errors"

	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
)

var db database.Service = database.New()

type User = models.User

// UserService defines the interface for interacting with users.
type UserService interface {
	// Create creates a new user.
	Create(user *User) error

	// GetByEmail returns a user by email.
	GetByEmail(email string) (*User, error)

	// GetByID returns a user by ID.
	GetByID(id int64) (*User, error)

	// Update updates a user.
	Update(user *User) error

	// Delete deletes a user.
	Delete(id int64) error
}

// NewUserService creates a new user service.
func NewUserService() UserService {
	return &userService{}
}

type userService struct{}

func (s *userService) Create(user *User) error {
	if user.Email == "" {
		return errors.New("Email is required")
	}
	if user.Name == "" {
		return errors.New("Name is required")
	}
	if user.Password == "" {
		return errors.New("Password is required")
	}

	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	query := "INSERT INTO user (name, email, password) VALUES (?, ?, ?)"

	stmt, err := db.Prepare(query)
	if err != nil {
		defer stmt.Close()
		return err
	}

	_, err = stmt.Exec(user.Name, user.Email, user.Password)
	if err != nil {
		defer stmt.Close()
		return err
	}

	defer stmt.Close()
	return nil
}

func (s *userService) GetByEmail(email string) (*User, error) {
	query := "SELECT id, name, email, password FROM user WHERE email = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		defer stmt.Close()
		return nil, err
	}

	row := stmt.QueryRow(email)
	user := &User{}
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		defer stmt.Close()
		return nil, err
	}

	defer stmt.Close()
	return user, nil
}

func (s *userService) GetByID(id int64) (*User, error) {
	query := "SELECT id, name, email, password FROM user WHERE id = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		defer stmt.Close()
		return nil, err
	}

	row := stmt.QueryRow(id)
	user := &User{}
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		defer stmt.Close()
		return nil, err
	}

	defer stmt.Close()
	return user, nil
}

func (s *userService) Update(user *User) error {
	if user.ID == 0 {
		return errors.New("ID is required")
	}
	if user.Email == "" {
		return errors.New("Email is required")
	}
	if user.Name == "" {
		return errors.New("Name is required")
	}
	if user.Password == "" {
		return errors.New("Password is required")
	}

	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	query := "UPDATE user SET name = ?, email = ?, password = ? WHERE id = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		defer stmt.Close()
		return err
	}

	_, err = stmt.Exec(user.Name, user.Email, user.Password, user.ID)
	if err != nil {
		defer stmt.Close()
		return err
	}

	defer stmt.Close()
	return nil
}

func (s *userService) Delete(id int64) error {
	query := "DELETE FROM user WHERE id = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		defer stmt.Close()
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		defer stmt.Close()
		return err
	}

	defer stmt.Close()
	return nil
}
