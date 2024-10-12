package services

import (
	"errors"
	"piggy-planner/internal/database"
	"piggy-planner/internal/models"
	"time"
)

type IncomesService interface {
	// Create a income in the database
	Create(income *models.Income) error
	// Get all incomes from the database for a specific user
	GetAll(fkUserId uint64) ([]models.Income, error)
	// Get a income by ID from the database
	GetByID(id uint64) (*models.Income, error)
	// Get a income by description from the database
	GetByDescription(description string) ([]models.Income, error)
	// Get all incomes during a specific period
	GetByPeriod(fkUserId uint64, startDate time.Time, endDate time.Time) ([]models.Income, error)
	// Update a income in the database
	Update(income *models.Income) error
	// Delete a income by ID from the database
	Delete(id uint64) error
}

type incomesService struct {
	DB database.DbService
}

func NewIncomesService(db database.DbService) IncomesService {
	return &incomesService{
		DB: db,
	}
}

func (s *incomesService) Create(income *models.Income) error {
	if income.Amount == 0 {
		return errors.New("Missing income amount")
	}

	if income.Description == "" {
		return errors.New("Missing income description")
	}

	if income.UserID == 0 {
		return errors.New("Missing income user ID")
	}

	date := income.Date.Format("2006-01-02")

	query := "INSERT INTO incomes (fk_user_id, amount, description, date, fk_income_type_id) VALUES (?, ?, ?, ?, ?)"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(income.UserID, income.Amount, income.Description, date, income.Type.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *incomesService) GetAll(fkUserId uint64) ([]models.Income, error) {
	query := `
    SELECT 
        incomes.id, 
        incomes.fk_user_id, 
        incomes.amount, 
        incomes.description, 
        income_types.id AS income_type_id, 
        income_types.name AS income_type_name, 
        DATE(incomes.date) AS date
    FROM 
        incomes 
    INNER JOIN 
        income_types 
    ON 
        incomes.fk_income_type_id = income_types.id 
    WHERE 
        incomes.fk_user_id = ? 
    ORDER BY 
        incomes.date DESC
	`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(fkUserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []models.Income
	for rows.Next() {
		var income models.Income
		var incomeType models.IncomeType
		var date []byte

		err := rows.Scan(&income.ID, &income.UserID, &income.Amount, &income.Description, &incomeType.ID, &incomeType.Name, &date)
		if err != nil {
			return nil, err
		}

		if string(date) != "0000-00-00" {
			income.Date, err = time.Parse("2006-01-02", string(date))
			if err != nil {
				return nil, err
			}
		} else {
			income.Date = time.Time{}
		}

		income.Type = incomeType

		incomes = append(incomes, income)
	}

	return incomes, nil
}

func (s *incomesService) GetByID(id uint64) (*models.Income, error) {
	query := `SELECT 
						incomes.id, 
						incomes.fk_user_id, 
						incomes.amount, 
						incomes.description, 
						income_types.id AS income_type_id, 
						income_types.name AS income_type_name, 
						DATE(incomes.date) AS date
					FROM 
						incomes 
					INNER JOIN 
						income_types 
					ON 
						incomes.fk_income_type_id = income_types.id 
					WHERE 
						incomes.id = ? 
					ORDER BY 
						incomes.date DESC
					LIMIT 1`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	var income models.Income
	var incomeType models.IncomeType
	var date []byte
	err = row.Scan(&income.ID, &income.UserID, &income.Amount, &income.Description, &incomeType.ID, &incomeType.Name, &date)
	if err != nil {
		return nil, err
	}

	if string(date) != "0000-00-00" {
		income.Date, err = time.Parse("2006-01-02", string(date))
		if err != nil {
			return nil, err
		}
	} else {
		income.Date = time.Time{}
	}

	income.Type = incomeType

	return &income, nil
}

func (s *incomesService) GetByDescription(description string) ([]models.Income, error) {
	query := `SELECT 
						incomes.id, 
						incomes.fk_user_id, 
						incomes.amount, 
						incomes.description, 
						income_types.id AS income_type_id, 
						income_types.name AS income_type_name, 
						DATE(incomes.date) AS date
					FROM 
						incomes 
					INNER JOIN 
						income_types 
					ON 
						incomes.fk_income_type_id = income_types.id 
					WHERE 
						incomes.description LIKE ?
					ORDER BY 
						incomes.date DESC
					`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	description = "%" + description + "%"
	row, err := stmt.Query(description)
	if err != nil {
		return nil, err
	}

	var incomes []models.Income
	for row.Next() {
		var income models.Income
		var incomeType models.IncomeType
		var date []byte

		err := row.Scan(&income.ID, &income.UserID, &income.Amount, &income.Description, &incomeType.ID, &incomeType.Name, &date)
		if err != nil {
			return nil, err
		}

		if string(date) != "0000-00-00" {
			income.Date, err = time.Parse("2006-01-02", string(date))
			if err != nil {
				return nil, err
			}
		} else {
			income.Date = time.Time{}
		}

		income.Type = incomeType

		incomes = append(incomes, income)
	}

	if len(incomes) == 0 {
		return nil, errors.New("Income not found")
	}

	return incomes, nil
}

func (s *incomesService) GetByPeriod(fkUserId uint64, startDate time.Time, endDate time.Time) ([]models.Income, error) {
	query := `
	SELECT
		incomes.id,
		incomes.fk_user_id,
		incomes.amount,
		incomes.description,
		income_types.id AS income_type_id,
		income_types.name AS income_type_name,
		DATE(incomes.date) AS date
	FROM
		incomes
	INNER JOIN
		income_types
	ON
		incomes.fk_income_type_id = income_types.id
	WHERE
		incomes.fk_user_id = ?
	AND
		incomes.date >= ?
	AND
		incomes.date <= ?
	ORDER BY
		incomes.date DESC
	`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(fkUserId, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var incomes []models.Income
	for rows.Next() {
		var income models.Income
		var incomeType models.IncomeType
		var date []byte

		err := rows.Scan(&income.ID, &income.UserID, &income.Amount, &income.Description, &incomeType.ID, &incomeType.Name, &date)
		if err != nil {
			return nil, err
		}

		if string(date) != "0000-00-00" {
			income.Date, err = time.Parse("2006-01-02", string(date))
			if err != nil {
				return nil, err
			}
		} else {
			income.Date = time.Time{}
		}

		income.Type = incomeType

		incomes = append(incomes, income)
	}

	return incomes, nil
}

func (s *incomesService) Update(income *models.Income) error {
	if income.Amount == 0 {
		return errors.New("Missing income amount")
	}

	if income.Description == "" {
		return errors.New("Missing income description")
	}

	if income.Type.ID == 0 {
		return errors.New("Missing income type ID")
	}

	if income.UserID == 0 {
		return errors.New("Missing income user ID")
	}

	date := income.Date.Format("2006-01-02")

	query := "UPDATE incomes SET amount = ?, description = ?, fk_income_type_id = ?, date = ? WHERE id = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(income.Amount, income.Description, income.Type.ID, date, income.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *incomesService) Delete(id uint64) error {
	query := "DELETE FROM incomes WHERE id = ?"

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

// IncomeTypesService represents a service for managing income types.
type IncomeTypesService interface {
	// Create a income type in the database
	Create(fkUserId uint64, incomeType *models.IncomeType) error
	// Get all income types from the database
	GetAll(fkUserID uint64) ([]models.IncomeType, error)
	// Get a income type by ID from the database
	GetByID(id uint64) (*models.IncomeType, error)
	// Update a income type in the database
	Update(incomeType *models.IncomeType) error
	// Delete a income type by ID from the database
	Delete(id uint64) error
}

type incomeTypeService struct {
	DB database.DbService
}

func NewIncomeTypesService(db database.DbService) IncomeTypesService {
	return &incomeTypeService{
		DB: db,
	}
}

func (s *incomeTypeService) Create(fkUserId uint64, incomeType *models.IncomeType) error {
	if incomeType.Name == "" {
		return errors.New("Missing income type name")
	}

	query := "INSERT INTO income_types (name, fk_user_id) VALUES (?, ?)"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(incomeType.Name, fkUserId)
	if err != nil {
		return err
	}

	return nil
}

func (s *incomeTypeService) GetAll(userID uint64) ([]models.IncomeType, error) {
	query := "SELECT id, name FROM income_types where fk_user_id = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	var incomeTypes []models.IncomeType
	for rows.Next() {
		var incomeType models.IncomeType
		err := rows.Scan(&incomeType.ID, &incomeType.Name)
		if err != nil {
			return nil, err
		}

		incomeTypes = append(incomeTypes, incomeType)
	}

	return incomeTypes, nil
}

func (s *incomeTypeService) GetByID(id uint64) (*models.IncomeType, error) {
	query := "SELECT id, name FROM income_types WHERE id = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	var incomeType models.IncomeType
	err = row.Scan(&incomeType.ID, &incomeType.Name)
	if err != nil {
		return nil, err
	}

	return &incomeType, nil
}

func (s *incomeTypeService) Update(incomeType *models.IncomeType) error {
	if incomeType.Name == "" {
		return errors.New("Missing income type name")
	}

	query := "UPDATE income_types SET name = ? WHERE id = ?"

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(incomeType.Name, incomeType.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *incomeTypeService) Delete(id uint64) error {
	query := "DELETE FROM income_types WHERE id = ?"

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
