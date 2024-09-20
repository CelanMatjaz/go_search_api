package postgres

import (
	"database/sql"
	"fmt"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func (s *PostgresStore) GetTags(accountId int) ([]types.Tag, error) {
	return getRecords(s, scanTagRow, "SELECT * FROM tags WHERE account_id = $1", accountId)
}

func (s *PostgresStore) GetTag(accountId int, tagId int) (*types.Tag, error) {
	return getRecord(s, scanTagRow, "SELECT * FROM tags WHERE account_id = $1 AND id = $2", accountId, tagId)
}

func (s *PostgresStore) CreateTag(accountId int, tag types.TagBody) (*types.Tag, error) {
	return createRecord(
		s, scanTagRow,
		"INSERT INTO tags (account_id, label, color) VALUES ($1, $2, $3) RETURNING *",
		accountId, tag.Label, tag.Color,
	)
}

func (s *PostgresStore) UpdateTag(accountId int, tagId int, tag types.TagBody) (*types.Tag, error) {
	return updateRecord(
		s, scanTagRow,
		"UPDATE tags SET label = $1, color = $2 WHERE id = $3 AND account_id = $4 RETURNING *",
		tag.Label, tag.Color, tagId, accountId,
	)
}

func (s *PostgresStore) DeleteTag(accountId int, tagId int) error {
	return deleteRecord(s, "DELETE FROM tags WHERE account_id = $1 AND id = $2", accountId, tagId)
}

func createTagJoinQuery(tableName string, colName string) string {
	return fmt.Sprintf(
		`SELECT * FROM tags LEFT JOIN %s t 
        ON t.tag_id = tag.id 
        WHERE tags.account_id = $1 AND t.%s = $2`, tableName, colName)
}

var (
	applicationPresetTagsQuery  = createTagJoinQuery("mtm_tags_application_presets", "preset_id")
	applicationSectionTagsQuery = createTagJoinQuery("mtm_tags_application_sections", "section_id")
	resumePresetTagsQuery       = createTagJoinQuery("mtm_tags_resume_presets", "preset_id")
	resumeSectionTagsQuery      = createTagJoinQuery("mtm_tags_resume_sections", "section_id")
)

func (s *PostgresStore) GetApplicationPresetTags(accountId int, applicationId int) ([]types.Tag, error) {
	return getRecords(
		s, scanTagRow, applicationPresetTagsQuery,
		accountId, applicationId,
	)
}

func (s *PostgresStore) GetApplicationSectionTags(accountId int, applicationSectionId int) ([]types.Tag, error) {
	return getRecords(
		s, scanTagRow, applicationSectionTagsQuery,
		accountId, applicationSectionId,
	)
}
func (s *PostgresStore) GetResumePresetTags(accountId int, resumeId int) ([]types.Tag, error) {
	return getRecords(
		s, scanTagRow, resumePresetTagsQuery,
		accountId, resumeId,
	)
}
func (s *PostgresStore) GetResumeSectionTags(accountId int, resumeSectionId int) ([]types.Tag, error) {
	return getRecords(
		s, scanTagRow, resumeSectionTagsQuery,
		accountId, resumeSectionId,
	)
}

func scanTagRow(row db.Scannable) (*types.Tag, error) {
	tag := &types.Tag{}
	err := row.Scan(
		&tag.Id,
		&tag.AccountId,
		&tag.Label,
		&tag.Color,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, err
		}
	}

	return tag, nil
}
