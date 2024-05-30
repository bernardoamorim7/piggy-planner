package services

import (
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
	// TODO
	return nil
}
