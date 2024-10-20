package postgres

import "github.com/CelanMatjaz/job_application_tracker_api/pkg/types"

func (s *PostgresStore) GetApplicationPresets(accountId int, pagination types.PaginationParams) ([]types.ApplicationPreset, error) {
	return s.ApplicationPresets.GetMany(accountId, pagination)
}

func (s *PostgresStore) GetApplicationPreset(accountId int, id int) (types.ApplicationPreset, error) {
	return s.ApplicationPresets.GetSingle(accountId, id)
}

func (s *PostgresStore) CreateApplicationPreset(body types.ApplicationPreset) (types.ApplicationPreset, error) {
	return s.ApplicationPresets.CreateSingle(body)
}

func (s *PostgresStore) UpdateApplicationPreset(body types.ApplicationPreset) (types.ApplicationPreset, error) {
	return s.ApplicationPresets.UpdateSingle(body)
}

func (s *PostgresStore) DeleteApplicationPreset(accountId int, id int) error {
	return s.ApplicationPresets.DeleteSingle(accountId, id)
}

func (s *PostgresStore) GetApplicationSections(accountId int, pagination types.PaginationParams) ([]types.ApplicationSection, error) {
	return s.ApplicationSections.GetMany(accountId, pagination)
}

func (s *PostgresStore) GetApplicationSection(accountId int, id int) (types.ApplicationSection, error) {
	return s.ApplicationSections.GetSingle(accountId, id)
}

func (s *PostgresStore) CreateApplicationSection(body types.ApplicationSection) (types.ApplicationSection, error) {
	return s.ApplicationSections.CreateSingle(body)
}

func (s *PostgresStore) UpdateApplicationSection(body types.ApplicationSection) (types.ApplicationSection, error) {
	return s.ApplicationSections.UpdateSingle(body)
}

func (s *PostgresStore) DeleteApplicationSection(accountId int, id int) error {
	return s.ApplicationPresets.DeleteSingle(accountId, id)
}

func (s *PostgresStore) GetResumePresets(accountId int, pagination types.PaginationParams) ([]types.ResumePreset, error) {
	return s.ResumePresets.GetMany(accountId, pagination)
}

func (s *PostgresStore) GetResumePreset(accountId int, id int) (types.ResumePreset, error) {
	return s.ResumePresets.GetSingle(accountId, id)
}

func (s *PostgresStore) CreateResumePreset(body types.ResumePreset) (types.ResumePreset, error) {
	return s.ResumePresets.CreateSingle(body)
}

func (s *PostgresStore) UpdateResumePreset(body types.ResumePreset) (types.ResumePreset, error) {
	return s.ResumePresets.UpdateSingle(body)
}

func (s *PostgresStore) DeleteResumePreset(accountId int, id int) error {
	return s.ApplicationPresets.DeleteSingle(accountId, id)
}

func (s *PostgresStore) GetResumeSections(accountId int, pagination types.PaginationParams) ([]types.ResumeSection, error) {
	return s.ResumeSections.GetMany(accountId, pagination)
}

func (s *PostgresStore) GetResumeSection(accountId int, id int) (types.ResumeSection, error) {
	return s.ResumeSections.GetSingle(accountId, id)
}

func (s *PostgresStore) CreateResumeSection(body types.ResumeSection) (types.ResumeSection, error) {
	return s.ResumeSections.CreateSingle(body)
}

func (s *PostgresStore) UpdateResumeSection(body types.ResumeSection) (types.ResumeSection, error) {
	return s.ResumeSections.UpdateSingle(body)
}

func (s *PostgresStore) DeleteResumeSection(accountId int, id int) error {
	return s.ResumeSections.DeleteSingle(accountId, id)
}
