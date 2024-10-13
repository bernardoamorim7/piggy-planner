package models

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system.
type User struct {
	ID       uint64 `db:"id"`       // Unique identifier
	Name     string `db:"name"`     // Name of the user
	Email    string `db:"email"`    // Email address of the user
	Password string `db:"password"` // Hashed password of the user
	Avatar   string `db:"avatar"`   // URL of Dicebear's API generated avatar
	IsAdmin  bool   `db:"is_admin"` // Whether the user is an admin or not (default: false)
}

// NewUser creates a new user.
func NewUser(name, email, password, avatar string, isAdmin bool) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
		Avatar:   avatar,
		IsAdmin:  isAdmin,
	}
}

// Validate validates the user.
func (u *User) Validate() error {
	if !ValidateEmail(u.Email) {
		return errors.New("Invalid email")
	}

	if !ValidatePassword(u.Password) {
		return errors.New("Password must be at least 8 characters")
	}

	if !ValidateName(u.Name) {
		return errors.New("Invalid name")
	}

	return nil
}

// HashPassword hashes the user's password.
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

// ComparePassword compares the user's password with the provided password.
func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ValidateEmail validates an email address.
func ValidateEmail(email string) bool {
	email = strings.TrimSpace(email)

	m, err := regexp.Match(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, []byte(email))
	if err != nil {
		return false
	}

	return m
}

// ValidatePassword validates a password.
func ValidatePassword(password string) bool {
	return len(password) >= 8
}

// ValidateName validates a name.
func ValidateName(name string) bool {
	if len(name) < 1 || len(name) > 255 {
		return false
	}

	m, err := regexp.Match(`^[a-zA-Z0-9._%+-]{1,255}$`, []byte(name))
	if err != nil {
		return false
	}

	return m
}
