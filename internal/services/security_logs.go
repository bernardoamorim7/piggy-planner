package services

import (
	"errors"
	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
)

type SecurityLogsService interface {
	// Create a security log in the database
	Create(securityLog *models.SecurityLog) error
}

type securityLogsService struct {
	DB database.Service
}

func NewSecurityLogsService(db database.Service) SecurityLogsService {
	return &securityLogsService{
		DB: db,
	}
}

func (s *securityLogsService) Create(log *models.SecurityLog) error {
	if log.Action == "" {
		return errors.New("Missing log action")
	}
	if log.IPAdress == "" {
		return errors.New("Missing log IP Address")
	}
	if log.UserAgent == "" {
		return errors.New("Missing log user agent")
	}
	if log.UserID == 0 {
		return errors.New("Missing log user ID")
	}

	query := "INSERT INTO security_logs (fk_user_id, action, ip_address, user_agent) VALUES (?, ?, ?, ?)"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(log.UserID, log.Action, log.IPAdress, log.UserAgent)
	if err != nil {
		return err
	}

	return nil
}
