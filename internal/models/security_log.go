package models

import (
	"time"
)

type action string

const (
	REGISTER        action = "register"
	LOGIN           action = "login"
	LOGOUT          action = "logout"
	PASSWORD_CHANGE action = "password_change"
	PASSWORD_RESET  action = "password_reset"
)

type SecurityLog struct {
	ID        uint64 `db:"id"`
	User      User
	Action    action    `db:"action"`
	IPAdress  string    `db:"ip_address"`
	UserAgent string    `db:"user_agent"`
	CreatedAt time.Time `db:"created_at"`
	UpdateAt  time.Time `db:"updated_at"`
}

func NewSecurityLog(id uint64, user User, action action, ipAdress string, userAgent string, createdAt time.Time, updatedAt time.Time) *SecurityLog {
	return &SecurityLog{
		ID:        id,
		User:      user,
		Action:    action,
		IPAdress:  ipAdress,
		UserAgent: userAgent,
		CreatedAt: createdAt,
		UpdateAt:  updatedAt,
	}
}
