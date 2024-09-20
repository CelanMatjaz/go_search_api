package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func (s *PostgresStore) GetResumePresets(accountId int) ([]types.ResumePreset, error) {
	return getRecords(s, scanResumePresetRow, "SELECT * FROM resume_presets WHERE account_id = $1", accountId)
}

func (s *PostgresStore) GetResumePreset(accountId int, presetId int) (*types.ResumePreset, error) {
	return getRecord(s, scanResumePresetRow, "SELECT * FROM resume_presets WHERE account_id = $1 AND id = $2", accountId, presetId)
}

func (s *PostgresStore) CreateResumePreset(accountId int, preset types.ResumePresetBody) (*types.ResumePreset, error) {
	return createRecord(
		s, scanResumePresetRow,
		"INSERT INTO resume_presets (account_id, label) VALUES ($1, $2) RETURNING *",
		accountId, preset.Label,
	)
}

func (s *PostgresStore) UpdateResumePreset(accountId int, presetId int, preset types.ResumePresetBody) (*types.ResumePreset, error) {
	return updateRecord(
		s, scanResumePresetRow,
		"UPDATE resume_presets SET label = $1, updated_at = DEFAULT WHERE id = $2 AND account_id = $3 RETURNING *",
		preset.Label, presetId, accountId,
	)
}

func (s *PostgresStore) DeleteResumePreset(accountId int, presetId int) error {
	return deleteRecord(s, "DELETE FROM resume_presets WHERE account_id = $1 AND id = $2", accountId, presetId)
}

func (s *PostgresStore) GetResumeSections(accountId int) ([]types.ResumeSection, error) {
	return getRecords(s, scanResumeSectionRow, "SELECT * FROM resume_sections WHERE account_id = $1", accountId)
}

func (s *PostgresStore) GetResumeSection(accountId int, sectionId int) (*types.ResumeSection, error) {
	return getRecord(s, scanResumeSectionRow, "SELECT * FROM resume_sections WHERE account_id = $1 AND id = $2", accountId, sectionId)
}

func (s *PostgresStore) CreateResumeSections(accountId int, section types.ResumeSectionBody) (*types.ResumeSection, error) {
	return createRecord(
		s, scanResumeSectionRow,
		"INSERT INTO resume_sections (account_id, label, text) VALUES ($1, $2, $3) RETURNING *",
		accountId, section.Label, section.Text,
	)
}

func (s *PostgresStore) UpdateResumeSections(accountId int, sectionId int, section types.ResumeSectionBody) (*types.ResumeSection, error) {
	return updateRecord(
		s, scanResumeSectionRow,
		"UPDATE resume_sections SET label = $1, text = $2, updated_at = DEFAULT WHERE id = $3 AND account_id = $4 RETURNING *",
		section.Label, section.Text, sectionId, accountId,
	)
}

func (s *PostgresStore) DeleteResumeSections(accountId int, sectionId int) error {
	return deleteRecord(s, "DELETE resume_sections WHERE account_id = $1 AND id = $2", accountId, sectionId)
}

func scanResumePresetRow(row db.Scannable) (*types.ResumePreset, error) {
	preset := &types.ResumePreset{}
	err := row.Scan(
		&preset.Id,
		&preset.AccountId,
		&preset.Label,
		&preset.CreatedAt,
		&preset.UpdatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, err
		}
	}

	return preset, nil
}

func scanResumeSectionRow(row db.Scannable) (*types.ResumeSection, error) {
	section := &types.ResumeSection{}
	err := row.Scan(
		&section.Id,
		&section.AccountId,
		&section.Label,
		&section.Text,
		&section.CreatedAt,
		&section.UpdatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, err
		}
	}

	return section, nil
}
