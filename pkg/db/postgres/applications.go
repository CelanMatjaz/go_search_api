package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func (s *PostgresStore) GetApplicationPresets(accountId int) ([]types.ApplicationPreset, error) {
	return getRecords(s, scanApplicationPresetRow, "SELECT * FROM application_presets WHERE account_id = $1", accountId)
}

func (s *PostgresStore) GetApplicationPreset(accountId int, presetId int) (*types.ApplicationPreset, error) {
	return getRecord(s, scanApplicationPresetRow, "SELECT * FROM application_presets WHERE account_id = $1 AND id = $2", accountId, presetId)
}

func (s *PostgresStore) CreateApplicationPreset(accountId int, preset types.ApplicationPresetBody) (*types.ApplicationPreset, error) {
	return createRecord(
		s, scanApplicationPresetRow,
		"INSERT INTO application_presets (account_id, label) VALUES ($1, $2) RETURNING *",
		accountId, preset.Label,
	)
}

func (s *PostgresStore) UpdateApplicationPreset(accountId int, presetId int, preset types.ApplicationPresetBody) (*types.ApplicationPreset, error) {
	return updateRecord(
		s, scanApplicationPresetRow,
		"UPDATE application_presets SET label = $1, updated_at = DEFAULT WHERE id = $2 AND account_id = $3 RETURNING *",
		preset.Label, presetId, accountId,
	)
}

func (s *PostgresStore) DeleteApplicationPreset(accountId int, presetId int) error {
	return deleteRecord(s, "DELETE FROM application_presets WHERE account_id = $1 AND id = $2", accountId, presetId)
}

func (s *PostgresStore) GetApplicationSections(accountId int) ([]types.ApplicationSection, error) {
	return getRecords(s, scanApplicationSectionRow, "SELECT * FROM application_sections WHERE account_id = $1", accountId)
}

func (s *PostgresStore) GetApplicationSection(accountId int, sectionId int) (*types.ApplicationSection, error) {
	return getRecord(s, scanApplicationSectionRow, "SELECT * FROM application_sections WHERE account_id = $1 AND id = $2", accountId, sectionId)
}

func (s *PostgresStore) CreateApplicationSections(accountId int, section types.ApplicationSectionBody) (*types.ApplicationSection, error) {
	return createRecord(
		s, scanApplicationSectionRow,
		"INSERT INTO application_sections (account_id, label, text) VALUES ($1, $2, $3) RETURNING *",
		accountId, section.Label, section.Text,
	)
}

func (s *PostgresStore) UpdateApplicationSections(accountId int, sectionId int, section types.ApplicationSectionBody) (*types.ApplicationSection, error) {
	return updateRecord(
		s, scanApplicationSectionRow,
		"UPDATE application_sections SET label = $1, text = $2, updated_at = DEFAULT WHERE id = $3 AND account_id = $4 RETURNING *",
		section.Label, section.Text, sectionId, accountId,
	)
}

func (s *PostgresStore) DeleteApplicationSections(accountId int, sectionId int) error {
	return deleteRecord(s, "DELETE application_sections WHERE account_id = $1 AND id = $2", accountId, sectionId)
}

func scanApplicationPresetRow(row db.Scannable) (*types.ApplicationPreset, error) {
	preset := &types.ApplicationPreset{}
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

func scanApplicationSectionRow(row db.Scannable) (*types.ApplicationSection, error) {
	section := &types.ApplicationSection{}
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
