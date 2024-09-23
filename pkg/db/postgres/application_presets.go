package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

var applicationPresetWithTagsQuery = createRecordWithTagsQuery("application_presets", "mtm_tags_application_presets")

func (s *PostgresStore) GetApplicationPresets(accountId int, pagination types.PaginationParams) ([]*types.ApplicationPreset, error) {
	return genericGetRecordsWithTags(s, scanApplicationPresetWithTag, applicationPresetWithTagsQuery, accountId, pagination)
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

func scanApplicationPresetWithTag(row db.Scannable) (*types.ApplicationPreset, *types.Tag, error) {
	preset := &types.ApplicationPreset{}
	tag := &types.Tag{}
	err := row.Scan(
		&preset.Id,
		&preset.AccountId,
		&preset.Label,
		&preset.CreatedAt,
		&preset.UpdatedAt,
		&tag.Id,
		&tag.AccountId,
		&tag.Label,
		&tag.Color,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil, nil
		default:
			return nil, nil, err
		}
	}

	return preset, tag, nil
}
