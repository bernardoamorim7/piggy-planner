package services

import (
	"errors"
	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
	"time"
)

type SecurityLogsService interface {
	// GetAll returns all security logs from the database
	GetAll() ([]models.SecurityLog, error)

	// GetByUserName returns all security logs from the database by user name
	GetByUserName(userName string) ([]models.SecurityLog, error)

	// Create a security log in the database
	Create(securityLog *models.SecurityLog) error
}

type securityLogsService struct {
	DB database.DbService
}

func NewSecurityLogsService(db database.DbService) SecurityLogsService {
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
	if log.User.ID == 0 {
		return errors.New("Missing log user ID")
	}

	query := "INSERT INTO security_logs (fk_user_id, action, ip_address, user_agent) VALUES (?, ?, ?, ?)"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(log.User.ID, log.Action, log.IPAdress, log.UserAgent)
	if err != nil {
		return err
	}

	return nil
}

func (s *securityLogsService) GetAll() ([]models.SecurityLog, error) {
	query := `
        SELECT 
            sl.id, 
            u.name AS user, 
            sl.action, 
            sl.ip_address, 
            sl.user_agent, 
            sl.created_at 
        FROM 
            security_logs sl
        JOIN 
            users u ON sl.fk_user_id = u.id
        ORDER BY 
            sl.created_at DESC
    `

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var logs []models.SecurityLog
	for rows.Next() {
		var log models.SecurityLog
		var date string
		var user models.User

		err := rows.Scan(&log.ID, &user.Name, &log.Action, &log.IPAdress, &log.UserAgent, &date)
		if err != nil {
			return nil, err
		}

		log.User = user

		log.CreatedAt, err = time.Parse("2006-01-02 15:04:05", date)
		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (s *securityLogsService) GetByUserName(userName string) ([]models.SecurityLog, error) {
	query := `
        SELECT 
            sl.id, 
            u.name AS user, 
            sl.action, 
            sl.ip_address, 
            sl.user_agent, 
            sl.created_at 
        FROM 
            security_logs sl
        JOIN 
            users u ON sl.fk_user_id = u.id
        WHERE 
            u.name LIKE ?
        ORDER BY 
            sl.created_at DESC
    `

	userName = "%" + userName + "%"

	rows, err := s.DB.Query(query, userName)
	if err != nil {
		return nil, err
	}

	var logs []models.SecurityLog
	for rows.Next() {
		var log models.SecurityLog
		var date string
		var user models.User

		err := rows.Scan(&log.ID, &user.Name, &log.Action, &log.IPAdress, &log.UserAgent, &date)
		if err != nil {
			return nil, err
		}

		log.User = user

		log.CreatedAt, err = time.Parse("2006-01-02 15:04:05", date)
		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	if len(logs) == 0 {
		return nil, errors.New("Security log not found")
	}

	return logs, nil
}
