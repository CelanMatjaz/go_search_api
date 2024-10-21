package postgres

import (
	"fmt"
	"strings"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
)

var tagQueryHolder = createQueryHolder[types.Tag](TagsTable, false)
var tagScanFunc = createScanFunc[types.Tag]()

func (s *PostgresStore) GetTags(accountId int) ([]types.Tag, error) {
	return getRecords(
		s.Db, tagScanFunc,
		tagQueryHolder.queryMany,
		accountId)
}

func (s *PostgresStore) GetTag(accountId int, tagId int) (types.Tag, error) {
	return WithTransactionScan(s.Db, getRecord, tagScanFunc, tagQueryHolder.querySingle, accountId, tagId)
}

func (s *PostgresStore) CreateTag(accountId int, tag types.Tag) (types.Tag, error) {
	return WithTransactionScan(
		s.Db, getRecord, tagScanFunc,
		tagQueryHolder.createSingle,
		accountId, tag.Label, tag.Color,
	)
}

func (s *PostgresStore) UpdateTag(accountId int, tagId int, tag types.Tag) (types.Tag, error) {
	println(tagQueryHolder.updateSingle)
	return WithTransactionScan(
		s.Db, getRecord, tagScanFunc,
		tagQueryHolder.updateSingle,
		utils.GetValuesFromBody(tag)...,
	)
}

func (s *PostgresStore) DeleteTag(accountId int, tagId int) error {
	return WithTransaction(s.Db, deleteRecord, tagQueryHolder.deleteSingle, tagId, accountId)
}

var (
	applicationPresetTagsQuery  = createTagJoinQuery("mtm_tags_application_presets")
	applicationSectionTagsQuery = createTagJoinQuery("mtm_tags_application_sections")
	resumePresetTagsQuery       = createTagJoinQuery("mtm_tags_resume_presets")
	resumeSectionTagsQuery      = createTagJoinQuery("mtm_tags_resume_sections")
)

func (s *PostgresStore) GetApplicationPresetTags(accountId int, id int) ([]types.Tag, error) {
	return getRecords(
		s.Db, tagScanFunc, applicationPresetTagsQuery,
		accountId, id,
	)
}

func (s *PostgresStore) GetApplicationSectionTags(accountId int, id int) ([]types.Tag, error) {
	return getRecords(
		s.Db, tagScanFunc, applicationSectionTagsQuery,
		accountId, id,
	)
}

func (s *PostgresStore) GetResumePresetTags(accountId int, id int) ([]types.Tag, error) {
	return getRecords(
		s.Db, tagScanFunc, resumePresetTagsQuery,
		accountId, id,
	)
}

func (s *PostgresStore) GetResumeSectionTags(accountId int, id int) ([]types.Tag, error) {
	return getRecords(
		s.Db, tagScanFunc, resumeSectionTagsQuery,
		accountId, id,
	)
}

func generateTagInsertString(start int, tagIds []int) string {
	sections := make([]string, len(tagIds))
	for i := range tagIds {
		sections[i] = fmt.Sprintf("($%d, (SELECT id FROM new_record))", i+start+1)
	}

	return strings.Join(sections, ",")
}
