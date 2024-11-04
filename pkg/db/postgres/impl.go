package postgres

import (
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func (s *PostgresStore) GetApplicationPresets(accountId int, pagination types.PaginationParams) (
	[]types.RecordWithTags[types.ApplicationPreset],
	error,
) {
	return s.ApplicationPresets.GetMany(accountId, pagination)
}

func (s *PostgresStore) GetApplicationPreset(accountId int, id int) (types.ApplicationPreset, error) {
	return s.ApplicationPresets.GetSingle(accountId, id)
}

func (s *PostgresStore) CreateApplicationPreset(accountId int, body types.ApplicationPreset) (types.ApplicationPreset, error) {
	return s.ApplicationPresets.CreateSingle(accountId, body)
}

func (s *PostgresStore) UpdateApplicationPreset(accountId int, id int, body types.ApplicationPreset) (types.ApplicationPreset, error) {
	return s.ApplicationPresets.UpdateSingle(accountId, id, body)
}

func (s *PostgresStore) DeleteApplicationPreset(accountId int, id int) error {
	return s.ApplicationPresets.DeleteSingle(accountId, id)
}

func (s *PostgresStore) GetApplicationSections(accountId int, pagination types.PaginationParams) (
	[]types.RecordWithTags[types.ApplicationSection],
	error,
) {
	return s.ApplicationSections.GetMany(accountId, pagination)
}

func (s *PostgresStore) GetApplicationSection(accountId int, id int) (types.ApplicationSection, error) {
	return s.ApplicationSections.GetSingle(accountId, id)
}

func (s *PostgresStore) CreateApplicationSection(accountId int, body types.ApplicationSection) (types.ApplicationSection, error) {
	return s.ApplicationSections.CreateSingle(accountId, body)
}

func (s *PostgresStore) UpdateApplicationSection(accountId int, id int, body types.ApplicationSection) (types.ApplicationSection, error) {
	return s.ApplicationSections.UpdateSingle(accountId, id, body)
}

func (s *PostgresStore) DeleteApplicationSection(accountId int, id int) error {
	return s.ApplicationPresets.DeleteSingle(accountId, id)
}

func (s *PostgresStore) GetResumePresets(accountId int, pagination types.PaginationParams) (
	[]types.RecordWithTags[types.ResumePreset],
	error,
) {
	return s.ResumePresets.GetMany(accountId, pagination)
}

func (s *PostgresStore) GetResumePreset(accountId int, id int) (types.ResumePreset, error) {
	return s.ResumePresets.GetSingle(accountId, id)
}

func (s *PostgresStore) CreateResumePreset(accountId int, body types.ResumePreset) (types.ResumePreset, error) {
	return s.ResumePresets.CreateSingle(accountId, body)
}

func (s *PostgresStore) UpdateResumePreset(accountId int, id int, body types.ResumePreset) (types.ResumePreset, error) {
	return s.ResumePresets.UpdateSingle(accountId, id, body)
}

func (s *PostgresStore) DeleteResumePreset(accountId int, id int) error {
	return s.ApplicationPresets.DeleteSingle(accountId, id)
}

func (s *PostgresStore) GetResumeSections(accountId int, pagination types.PaginationParams) (
	[]types.RecordWithTags[types.ResumeSection],
	error,
) {
	return s.ResumeSections.GetMany(accountId, pagination)
}

func (s *PostgresStore) GetResumeSection(accountId int, id int) (types.ResumeSection, error) {
	return s.ResumeSections.GetSingle(accountId, id)
}

func (s *PostgresStore) CreateResumeSection(accountId int, body types.ResumeSection) (types.ResumeSection, error) {
	return s.ResumeSections.CreateSingle(accountId, body)
}

func (s *PostgresStore) UpdateResumeSection(accountId int, id int, body types.ResumeSection) (types.ResumeSection, error) {
	return s.ResumeSections.UpdateSingle(accountId, id, body)
}

func (s *PostgresStore) DeleteResumeSection(accountId int, id int) error {
	return s.ResumeSections.DeleteSingle(accountId, id)
}
