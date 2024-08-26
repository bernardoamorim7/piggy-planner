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
	ID        uint64 `db:"id"`
	UserID    uint64 `db:"fk_user_id"`
	Action    action `db:"action"`
	IPAdress  string `db:"ip_address"`
	UserAgent string `db:"user_agent"`
}

func NewSecurityLog(id uint64, userId uint64, action action, ipAdress string, userAgent string) *SecurityLog {
	return &SecurityLog{
		ID:        id,
		UserID:    userId,
		Action:    action,
		IPAdress:  ipAdress,
		UserAgent: userAgent,
	}
}
