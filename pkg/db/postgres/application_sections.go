package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

var applicationSectionWithTagsQuery = createRecordWithTagsQuery("application_sections", "mtm_tags_application_sections")

func (s *PostgresStore) GetApplicationSections(accountId int, pagination types.PaginationParams) ([]*types.ApplicationSection, error) {
	return genericGetRecordsWithTags(s, scanApplicationSectionWithTag, applicationSectionWithTagsQuery, accountId, pagination)
}

func (s *PostgresStore) GetApplicationSection(accountId int, sectionId int) (*types.ApplicationSection, error) {
	return getRecord(s, scanApplicationSectionRow, "SELECT * FROM application_sections WHERE account_id = $1 AND id = $2", accountId, sectionId)
}

func (s *PostgresStore) CreateApplicationSection(accountId int, section types.ApplicationSectionBody) (*types.ApplicationSection, error) {
	return createRecord(
		s, scanApplicationSectionRow,
		"INSERT INTO application_sections (account_id, label, text) VALUES ($1, $2, $3) RETURNING *",
		accountId, section.Label, section.Text,
	)
}

func (s *PostgresStore) UpdateApplicationSection(accountId int, sectionId int, section types.ApplicationSectionBody) (*types.ApplicationSection, error) {
	return updateRecord(
		s, scanApplicationSectionRow,
		"UPDATE application_sections SET label = $1, text = $2, updated_at = DEFAULT WHERE id = $3 AND account_id = $4 RETURNING *",
		section.Label, section.Text, sectionId, accountId,
	)
}

func (s *PostgresStore) DeleteApplicationSection(accountId int, sectionId int) error {
	return deleteRecord(s, "DELETE application_sections WHERE account_id = $1 AND id = $2", accountId, sectionId)
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

func scanApplicationSectionWithTag(row db.Scannable) (*types.ApplicationSection, *types.Tag, error) {
	section := &types.ApplicationSection{}
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
