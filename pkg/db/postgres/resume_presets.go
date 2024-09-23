package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

var resumePresetWithTagsQuery = createRecordWithTagsQuery("resume_presets", "mtm_tags_resume_presets")

func (s *PostgresStore) GetResumePresets(accountId int, pagination types.PaginationParams) ([]*types.ResumePreset, error) {
	return genericGetRecordsWithTags(s, scanResumePresetWithTag, resumePresetWithTagsQuery, accountId, pagination)
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

func scanResumePresetWithTag(row db.Scannable) (*types.ResumePreset, *types.Tag, error) {
	preset := &types.ResumePreset{}
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
