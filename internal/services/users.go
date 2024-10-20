package services

import (
	"errors"
	"time"

	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
)

// UserService defines the interface for interacting with users.
type UserService interface {
	// Create creates a new user.
	Create(user *models.User) error

	// Update updates a user.
	Update(user *models.User) error

	// Delete deletes a user.
	Delete(id uint64) error

	// GetByEmail returns a user by email.
	GetByEmail(email string) (*models.User, error)

	// GetByID returns a user by ID.
	GetByID(id uint64) (*models.User, error)

	// GetAll returns all users.
	GetAll() ([]models.User, error)

	// GetByUserName returns a user by username.
	GetByUserName(userName string) ([]models.User, error)
}

// NewUserService creates a new user service.
func NewUserService(db database.DbService) UserService {
	return &userService{
		DB: db,
	}
}

type userService struct {
	DB database.DbService
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	query := "SELECT id, name, email, password, avatar, is_admin FROM users WHERE email = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	row := *stmt.QueryRow(email)

	user := &models.User{}

	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Avatar, &user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(id uint64) (*models.User, error) {
	query := "SELECT id, name, email, password, avatar, is_admin FROM users WHERE id = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)
	user := &models.User{}
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Avatar, &user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Create(user *models.User) error {
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

	// Check if there are any existing users
	var count uint64
	err := s.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	// If no users exist, set is_admin to true
	if count == 0 {
		user.IsAdmin = true
	} else {
		user.IsAdmin = false
	}

	query := "INSERT INTO users (name, email, password, avatar, is_admin) VALUES (?, ?, ?, ?, ?)"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.Name, user.Email, user.Password, user.Avatar, user.IsAdmin)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) Update(user *models.User) error {
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

	u, err := s.GetByEmail(user.Email)
	if err != nil {
		if u.Email == user.Email {
			return errors.New("E-mail already exists")
		}
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	query := "UPDATE users SET name = ?, email = ?, password = ?, avatar = ?, is_admin = ?, updated_at = ? WHERE id = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	updateTime := time.Now()

	user.Avatar = "https://api.dicebear.com/8.x/thumbs/png?seed=" + user.Name

	_, err = stmt.Exec(user.Name, user.Email, user.Password, user.Avatar, user.IsAdmin, updateTime, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) Delete(id uint64) error {
	query := "DELETE FROM users WHERE id = ?"

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

func (s *userService) GetAll() ([]models.User, error) {
	query := "SELECT id, name, email, password, avatar, is_admin FROM users"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	users := []models.User{}

	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Avatar, &user.IsAdmin)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *userService) GetByUserName(userName string) ([]models.User, error) {
	query := "SELECT id, name, email, password, avatar, is_admin FROM users WHERE name LIKE ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	userName = "%" + userName + "%"

	rows, err := stmt.Query(userName)
	if err != nil {
		return nil, err
	}

	users := []models.User{}

	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Avatar, &user.IsAdmin)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
