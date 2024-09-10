package applications

import (
	"database/sql"
	"errors"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(connection *db.DbConnection) *Store {
	return &Store{db: connection.DB}
}

func (s *Store) GetUserApplications(userId int, pagination service.PaginationParams) ([]types.Application, error) {
	rows, err := s.db.Query(
		"SELECT * FROM job_applications WHERE user_id = $1 OFFSET $2 LIMIT $3",
		userId, pagination.GetOffset(), pagination.Count)
	if err != nil {
		return []types.Application{}, err
	}

	applications := make([]types.Application, 0)
	for rows.Next() {
		a, err := scanApplicationRow(rows)
		if err != nil {
			return []types.Application{}, err
		}
		applications = append(applications, a)
	}

	return applications, nil
}

func (s *Store) GetUserApplication(userId int, applicationId int) (types.Application, error) {
	row := s.db.QueryRow("SELECT * FROM job_applications WHERE user_id = $1 AND id = $2", userId, applicationId)
	a, err := scanApplicationRow(row)
	switch err {
	case nil:
		return a, nil
	case sql.ErrNoRows:
		return types.Application{}, types.RecordDoesNotExist
	default:
		return types.Application{}, err
	}

}

func (s *Store) CreateUserApplication(userId int, applicationData types.Application) (types.Application, error) {
	row := s.db.QueryRow(
		"INSERT INTO job_applications (user_id, label, text) VALUES ($1, $2, $3) RETURNING *",
		userId, applicationData.Label, applicationData.Text,
	)
	a, err := scanApplicationRow(row)
	switch err {
	case nil:
		return a, nil
	case sql.ErrNoRows:
		return types.Application{}, types.RecordDoesNotExist
	default:
		return types.Application{}, err
	}
}

func (s *Store) UpdateUserApplication(userId int, applicationData types.Application) (types.Application, error) {
	row := s.db.QueryRow(
		"UPDATE job_applications SET label = $1, text = $2, updated_at = DEFAULT WHERE user_id = $3 AND id = $4 RETURNING *",
		applicationData.Label, applicationData.Text, userId, applicationData.Id,
	)
	a, err := scanApplicationRow(row)
	switch err {
	case nil:
		return a, nil
	case sql.ErrNoRows:
		return types.Application{}, types.RecordDoesNotExist
	default:
		return types.Application{}, err
	}
}

func (s *Store) DeleteUserApplication(userId int, applicationId int) error {
	_, err := s.db.Exec("DELETE FROM job_applications WHERE user_id = $1 AND id = $2", userId, applicationId)
	return err
}

func (s *Store) GetApplicationSections(userId int, pagination service.PaginationParams) ([]types.ApplicationSection, error) {
	rows, err := s.db.Query(
		"SELECT FROM * FROM job_application_sections WHERE user_id = $1",
		userId, pagination.GetOffset(), pagination.Count)
	if err != nil {
		return []types.ApplicationSection{}, err
	}

	sections := make([]types.ApplicationSection, 0)
	for rows.Next() {
		s, err := scanSectionRow(rows)
		if err != nil {
			return []types.ApplicationSection{}, err
		}
		sections = append(sections, s)
	}

	return sections, nil
}

func (s *Store) CreateApplicationSections(userId int, sectionData types.ApplicationSection) (types.ApplicationSection, error) {
	row := s.db.QueryRow(
		"INSERT INTO job_application_sections (user_id, label, text) VALUES ($1, $2, $3) RETURNING *",
		userId, sectionData.Label, sectionData.Text,
	)
	section, err := scanSectionRow(row)
	switch err {
	case nil:
		return section, nil
	case sql.ErrNoRows:
		return types.ApplicationSection{}, types.RecordDoesNotExist
	default:
		return types.ApplicationSection{}, err
	}
}

func (s *Store) UpdateApplicationSections(userId int, sectionData types.ApplicationSection) (types.ApplicationSection, error) {
	row := s.db.QueryRow(
		"UPDATE job_application_sections SET label = $1, text = $2, updated_at = DEFAULT WHERE user_id = $3 AND id = $4 RETURNING *",
		sectionData.Label, sectionData.Text, userId, sectionData.Id,
	)
	section, err := scanSectionRow(row)
	switch err {
	case nil:
		return section, nil
	case sql.ErrNoRows:
		return types.ApplicationSection{}, types.RecordDoesNotExist
	default:
		return types.ApplicationSection{}, err
	}
}

func (s *Store) DeleteApplicationSections(userId int, sectionId int) error {
	_, err := s.db.Exec("DELETE FROM job_applications WHERE user_id = $1 AND id = $2", userId, sectionId)
	return err
}

func scanApplicationRow(row db.Scannable) (types.Application, error) {
	var a types.Application
	err := row.Scan(
		&a.Id,
		&a.UserId,
		&a.Label,
		&a.Text,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	return a, err
}

func scanSectionRow(row db.Scannable) (types.ApplicationSection, error) {
	var s types.ApplicationSection
	err := row.Scan(
		&s.Id,
		&s.UserId,
		&s.Label,
		&s.Text,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return types.ApplicationSection{}, types.UserDoesNotExistErr
	}
	if err != nil {
		return types.ApplicationSection{}, err
	}

	return s, nil
}
