package postgres

import (
	"database/sql"
	"fmt"
	"strings"

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
	return WithTransactionScan(
		s, createRecord, scanTagRow,
		"INSERT INTO tags (account_id, label, color) VALUES ($1, $2, $3) RETURNING *",
		accountId, tag.Label, tag.Color,
	)
}

func (s *PostgresStore) UpdateTag(accountId int, tagId int, tag types.TagBody) (*types.Tag, error) {
	return WithTransactionScan(
		s, updateRecord, scanTagRow,
		"UPDATE tags SET label = $1, color = $2 WHERE id = $3 AND account_id = $4 RETURNING *",
		tag.Label, tag.Color, tagId, accountId,
	)
}

func (s *PostgresStore) DeleteTag(accountId int, tagId int) error {
	return WithTransaction(s, deleteRecord, "DELETE FROM tags WHERE account_id = $1 AND id = $2", accountId, tagId)
}

func createTagJoinQuery(tableName string) string {
	return fmt.Sprintf(
		`SELECT * FROM tags LEFT JOIN %s t ON t.tag_id = tag.id 
        WHERE tags.account_id = $1 AND t.record_id = $2`, tableName)
}

var (
	applicationPresetTagsQuery  = createTagJoinQuery("mtm_tags_application_presets")
	applicationSectionTagsQuery = createTagJoinQuery("mtm_tags_application_sections")
	resumePresetTagsQuery       = createTagJoinQuery("mtm_tags_resume_presets")
	resumeSectionTagsQuery      = createTagJoinQuery("mtm_tags_resume_sections")
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

func recordWithTagsQuery(recordTable string, mtmTable string) (string, string) {
	left := fmt.Sprintf(`
    WITH new_record AS (
        INSERT INTO %s (account_id, label, text) 
        VALUES ($1, $2, $3)
        RETURNING *
    ),
    _ AS (
        INSERT INTO %s (tag_id, record_id) 
        VALUES 
    `, recordTable, mtmTable)

	right := `) SELECT cs.* FROM new_record cs`

	return left, right
}

func generateTagInsertString(tagIds []int) string {
	sections := make([]string, len(tagIds))
	start := 4 // start is based on above function args
	for i := range tagIds {
		sections[i] = fmt.Sprintf("($%d, (SELECT id FROM new_record))", i+start)
	}

	return strings.Join(sections, ",")
}

func createRecordWithTagsQuery(tableName string, associationTableName string) string {
	return fmt.Sprintf(
		`   WITH limited_records AS (
                SELECT *
                FROM %s r
                WHERE account_id = $1
                ORDER BY r.id
                OFFSET $2 LIMIT $3
            )
            SELECT lr.*, t.*
            FROM limited_records lr
            LEFT JOIN %s ta ON lr.id = ta.record_id
            LEFT JOIN tags t ON ta.tag_id = t.id;`, tableName, associationTableName)
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
