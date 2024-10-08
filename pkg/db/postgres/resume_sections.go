package postgres

import (
	"database/sql"
	"strings"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

var resumeSectionWithTagsQuery = createRecordWithTagsQuery("resume_sections", "mtm_tags_resume_sections")

func (s *PostgresStore) GetResumeSections(accountId int, pagination types.PaginationParams) ([]*types.ResumeSection, error) {
	return genericGetRecordsWithTags(s, scanResumeSectionWithTag, resumeSectionWithTagsQuery, accountId, pagination)
}

func (s *PostgresStore) GetResumeSection(accountId int, sectionId int) (*types.ResumeSection, error) {
	return getRecord(s, scanResumeSectionRow, "SELECT * FROM resume_sections WHERE account_id = $1 AND id = $2", accountId, sectionId)
}

var createRSL, createRSR = recordWithTagsQuery("resume_sections", "mtm_tags_resume_sections")

func (s *PostgresStore) CreateResumeSection(accountId int, section types.ResumeSectionBody) (*types.ResumeSection, error) {
	args := make([]any, len(section.TagIds)+3)
	args[0] = accountId
	args[1] = section.Label
	args[2] = section.Text
	for i, tagId := range section.TagIds {
		args[i+3] = tagId
	}

	query := "INSERT INTO resume_sections (account_id, label, text) VALUES ($1, $2, $3) RETURNING *"
	if len(section.TagIds) > 0 {
		query = strings.Join([]string{createASL, generateTagInsertString(section.TagIds), createASR}, "")
	}

	return WithTransactionScan(s, createRecord,
		scanResumeSectionRow, query, args...,
	)
}

func (s *PostgresStore) UpdateResumeSection(accountId int, sectionId int, section types.ResumeSectionBody) (*types.ResumeSection, error) {
	return WithTransactionScan(
		s, updateRecord, scanResumeSectionRow,
		"UPDATE resume_sections SET label = $1, text = $2, updated_at = DEFAULT WHERE id = $3 AND account_id = $4 RETURNING *",
		section.Label, section.Text, sectionId, accountId,
	)
}

func (s *PostgresStore) DeleteResumeSection(accountId int, sectionId int) error {
	return WithTransaction(s, deleteRecord,
		"DELETE resume_sections WHERE account_id = $1 AND id = $2", accountId, sectionId)
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

func scanResumeSectionWithTag(row db.Scannable) (*types.ResumeSection, *types.Tag, error) {
	section := &types.ResumeSection{}
	tag := &types.Tag{}
	err := row.Scan(
		&section.Id,
		&section.AccountId,
		&section.Label,
		&section.Text,
		&section.CreatedAt,
		&section.UpdatedAt,
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

	return section, tag, nil
}
