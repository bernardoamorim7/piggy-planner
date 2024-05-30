package models

type action string

const (
	REGISTER        action = "register"
	LOGIN           action = "login"
	LOGOUT          action = "logout"
	PASSWORD_CHANGE action = "password_change"
	PASSWORD_RESET  action = "password_reset"
)

type SecurityLog struct {
	ID        int64
	UserID    int64
	Action    action
	IPAdress  string
	UserAgent string
}

func NewSecurityLog(id int64, userId int64, action action, ipAdress string, userAgent string) *SecurityLog {
	return &SecurityLog{
		ID:        id,
		UserID:    userId,
		Action:    action,
		IPAdress:  ipAdress,
		UserAgent: userAgent,
	}
}
