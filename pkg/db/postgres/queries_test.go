package postgres_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db/postgres"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
)

type TestQueryType struct {
	types.WithId
	types.WithAccountId
	TestField1 string `body:"create,update" db:"test_field1"`
	TestField2 string `body:"create,update" db:"test_field2"`
	TestField3 string `db:""`
	TestField4 string `db:""`
	types.WithTimestamps
}

type QueryTestCase struct {
	expectedQuery string
}

type SpecificTypeQueryTestCase struct {
	shouldFail     bool
	expectedQuery  string
	generatedQuery string
}

func specificTestCase(expected string, generated string) SpecificTypeQueryTestCase {
	return SpecificTypeQueryTestCase{
		shouldFail:     false,
		expectedQuery:  utils.TrimInside(expected),
		generatedQuery: utils.TrimInside(generated),
	}
}

func testQueries(fnName string, testCases []QueryTestCase, generatedQuery string, t *testing.T) {
	generatedQuery = utils.TrimInside(generatedQuery)
	for i, tc := range testCases {
		if generatedQuery != utils.TrimInside(tc.expectedQuery) {
			t.Fatalf("%s failed for test case with index %d\nqueries\nexpected:  %s\ngenerated: %s",
				fnName, i, tc.expectedQuery, generatedQuery,
			)
		}
	}
}

func testSpecificQueries(fnName string, testCases []SpecificTypeQueryTestCase, t *testing.T) {
	t.Log("testing for specific types")
	for i, tc := range testCases {
		if utils.TrimInside(tc.generatedQuery) != utils.TrimInside(tc.expectedQuery) {
			t.Fatalf("%s failed for specific type test case with index %d\nqueries\nexpected:  %s\ngenerated: %s",
				fnName, i, tc.expectedQuery, tc.generatedQuery,
			)
		}
	}
}

func TestSingleRecordQuery(t *testing.T) {
	t.Parallel()

	testCases := []QueryTestCase{
		{expectedQuery: "SELECT id, account_id, test_field1, test_field2, created_at, updated_at FROM table WHERE account_id = $1 AND id = $2"},
	}

	generatedQuery := postgres.SingleRecordQuery[TestQueryType]("table")
	testQueries("SingleRecordQuery()", testCases, generatedQuery, t)
}

func TestManyRecordsQuery(t *testing.T) {
	t.Parallel()

	type genFunction = func(string, bool) string

	createTCs := func(table string, query string, fn genFunction) []SpecificTypeQueryTestCase {
		return []SpecificTypeQueryTestCase{
			{expectedQuery: query, generatedQuery: fn(table, false), shouldFail: false},
			{expectedQuery: query + " OFFSET $2 LIMIT $3", generatedQuery: fn(table, true), shouldFail: false},
		}
	}

	specificTestCases := make([]SpecificTypeQueryTestCase, 0)
	specificTestCases = append(
		specificTestCases,
		createTCs(postgres.TagsTable, `
        SELECT id, account_id, label, color
        FROM tags
        WHERE account_id = $1`, postgres.ManyRecordsQuery[types.Tag])...,
	)
	specificTestCases = append(
		specificTestCases,
		createTCs(postgres.ApplicationPresetsTable, `
        SELECT id, account_id, label, created_at, updated_at
        FROM application_presets
        WHERE account_id = $1`, postgres.ManyRecordsQuery[types.ApplicationPreset])...,
	)
	specificTestCases = append(
		specificTestCases,
		createTCs(postgres.ApplicationSectionsTable, `
        SELECT id, account_id, label, text, created_at, updated_at
        FROM application_sections
        WHERE account_id = $1`, postgres.ManyRecordsQuery[types.ApplicationSection])...,
	)
	specificTestCases = append(
		specificTestCases,
		createTCs(postgres.ResumePresetsTable, `
        SELECT id, account_id, label, created_at, updated_at
        FROM resume_presets
        WHERE account_id = $1`, postgres.ManyRecordsQuery[types.ResumePreset])...,
	)
	specificTestCases = append(
		specificTestCases,
		createTCs(postgres.ResumeSectionsTable, `
        SELECT id, account_id, label, text, created_at, updated_at
        FROM resume_sections
        WHERE account_id = $1`, postgres.ManyRecordsQuery[types.ResumeSection])...,
	)

	testSpecificQueries("ManyRecordsQuery()", specificTestCases, t)
}

func TestCreateRecordQuery(t *testing.T) {
	t.Parallel()

	testCases := []QueryTestCase{
		{expectedQuery: "INSERT INTO table (account_id, test_field1, test_field2) VALUES ($1, $2, $3) RETURNING *"},
	}

	generatedQuery := postgres.CreateRecordQuery[TestQueryType]("table")
	testQueries("CreateRecordQuery()", testCases, generatedQuery, t)

	specificTestCases := make([]SpecificTypeQueryTestCase, 0)
	specificTestCases = append(
		specificTestCases,
		specificTestCase(
			"INSERT INTO tags (account_id, label, color) VALUES ($1, $2, $3) RETURNING *",
			postgres.CreateRecordQuery[types.Tag](postgres.TagsTable),
		),
		specificTestCase(
			"INSERT INTO application_presets (account_id, label) VALUES ($1, $2) RETURNING *",
			postgres.CreateRecordQuery[types.ApplicationPreset](postgres.ApplicationPresetsTable),
		),
		specificTestCase(
			"INSERT INTO application_sections (account_id, label, text) VALUES ($1, $2, $3) RETURNING *",
			postgres.CreateRecordQuery[types.ApplicationSection](postgres.ApplicationSectionsTable),
		),
		specificTestCase(
			"INSERT INTO resume_presets (account_id, label) VALUES ($1, $2) RETURNING *",
			postgres.CreateRecordQuery[types.ResumePreset](postgres.ResumePresetsTable),
		),
		specificTestCase(
			"INSERT INTO resume_sections (account_id, label, text) VALUES ($1, $2, $3) RETURNING *",
			postgres.CreateRecordQuery[types.ResumeSection](postgres.ResumeSectionsTable),
		),
	)
	testSpecificQueries("CreateRecordQuery()", specificTestCases, t)
}

func TestUpdateRecordQuery(t *testing.T) {
	t.Parallel()

	testCases := []QueryTestCase{
		{expectedQuery: "UPDATE table SET test_field1 = $3, test_field2 = $4, updated_at = DEFAULT WHERE id = $1 AND account_id = $2 RETURNING *"},
	}

	generatedQuery := postgres.UpdateRecordQuery[TestQueryType]("table")
	testQueries("UpdateRecordQuery()", testCases, generatedQuery, t)

	specificTestCases := make([]SpecificTypeQueryTestCase, 0)
	specificTestCases = append(
		specificTestCases,
		specificTestCase(
			"UPDATE tags SET label = $3, color = $4 WHERE id = $1 AND account_id = $2 RETURNING *",
			postgres.UpdateRecordQuery[types.Tag](postgres.TagsTable),
		),
		specificTestCase(
			"UPDATE application_presets SET label = $3, updated_at = DEFAULT WHERE id = $1 AND account_id = $2 RETURNING *",
			postgres.UpdateRecordQuery[types.ApplicationPreset](postgres.ApplicationPresetsTable),
		),
		specificTestCase(
			"UPDATE application_sections SET label = $3, text = $4, updated_at = DEFAULT WHERE id = $1 AND account_id = $2 RETURNING *",
			postgres.UpdateRecordQuery[types.ApplicationSection](postgres.ApplicationSectionsTable),
		),
		specificTestCase(
			"UPDATE resume_presets SET label = $3, updated_at = DEFAULT WHERE id = $1 AND account_id = $2 RETURNING *",
			postgres.UpdateRecordQuery[types.ResumePreset](postgres.ResumePresetsTable),
		),
		specificTestCase(
			"UPDATE resume_sections SET label = $3, text = $4, updated_at = DEFAULT WHERE id = $1 AND account_id = $2 RETURNING *",
			postgres.UpdateRecordQuery[types.ResumeSection](postgres.ResumeSectionsTable),
		),
	)
	testSpecificQueries("UpdateRecordQuery()", specificTestCases, t)
}

func TestDeleteRecordQuery(t *testing.T) {
	t.Parallel()

	testCases := []QueryTestCase{
		{expectedQuery: "DELETE FROM table WHERE id = $1 AND account_id = $2"},
	}

	generatedQuery := postgres.DeleteRecordQuery("table")
	testQueries("DeleteRecordQuery()", testCases, generatedQuery, t)

	specificTestCases := make([]SpecificTypeQueryTestCase, 0)
	specificTestCases = append(
		specificTestCases,
		specificTestCase(
			"DELETE FROM tags WHERE id = $1 AND account_id = $2",
			postgres.DeleteRecordQuery(postgres.TagsTable),
		),
		specificTestCase(
			"DELETE FROM application_presets WHERE id = $1 AND account_id = $2",
			postgres.DeleteRecordQuery(postgres.ApplicationPresetsTable),
		),
		specificTestCase(
			"DELETE FROM application_sections WHERE id = $1 AND account_id = $2",
			postgres.DeleteRecordQuery(postgres.ApplicationSectionsTable),
		),
		specificTestCase(
			"DELETE FROM resume_presets WHERE id = $1 AND account_id = $2",
			postgres.DeleteRecordQuery(postgres.ResumePresetsTable),
		),
		specificTestCase(
			"DELETE FROM resume_sections WHERE id = $1 AND account_id = $2",
			postgres.DeleteRecordQuery(postgres.ResumeSectionsTable),
		),
	)
	testSpecificQueries("DeleteRecordQuery()", specificTestCases, t)
}

func TestManyRecordsWithTagsQuery(t *testing.T) {
	t.Parallel()

	expectedQuery := `WITH limited_records AS (SELECT id, account_id, test_field1, test_field2, created_at, updated_at FROM table WHERE account_id = $1 OFFSET $2 LIMIT $3) SELECT lr.id, lr.account_id, lr.test_field1, lr.test_field2, lr.created_at, lr.updated_at, t.id, t.account_id, t.label, t.color FROM limited_records lr LEFT JOIN mtm_table ta ON lr.id = ta.record_id LEFT JOIN tags t ON ta.tag_id = t.id ORDER BY lr.updated_at DESC`

	testCases := []QueryTestCase{
		{expectedQuery: expectedQuery},
	}

	generatedQuery := postgres.ManyRecordsWithTagsQuery[TestQueryType]("table", "mtm_table")
	testQueries("ManyRecordsWithTagsQuery()", testCases, generatedQuery, t)

	specificTestCases := make([]SpecificTypeQueryTestCase, 0)
	specificTestCases = append(
		specificTestCases,
		specificTestCase(
			"WITH limited_records AS (SELECT id, account_id, label, created_at, updated_at FROM application_presets WHERE account_id = $1 OFFSET $2 LIMIT $3) SELECT lr.id, lr.account_id, lr.label, lr.created_at, lr.updated_at, t.id, t.account_id, t.label, t.color FROM limited_records lr LEFT JOIN mtm_tags_application_presets ta ON lr.id = ta.record_id LEFT JOIN tags t ON ta.tag_id = t.id ORDER BY lr.updated_at DESC",
			postgres.ManyRecordsWithTagsQuery[types.ApplicationPreset](postgres.ApplicationPresetsTable, postgres.MtmTagsAppPre),
		),
		specificTestCase(
			"WITH limited_records AS (SELECT id, account_id, label, text, created_at, updated_at FROM application_sections WHERE account_id = $1 OFFSET $2 LIMIT $3) SELECT lr.id, lr.account_id, lr.label, lr.text, lr.created_at, lr.updated_at, t.id, t.account_id, t.label, t.color FROM limited_records lr LEFT JOIN mtm_tags_application_sections ta ON lr.id = ta.record_id LEFT JOIN tags t ON ta.tag_id = t.id ORDER BY lr.updated_at DESC",
			postgres.ManyRecordsWithTagsQuery[types.ApplicationSection](postgres.ApplicationSectionsTable, postgres.MtmTagsAppSec),
		),
		specificTestCase(
			"WITH limited_records AS (SELECT id, account_id, label, created_at, updated_at FROM resume_presets WHERE account_id = $1 OFFSET $2 LIMIT $3) SELECT lr.id, lr.account_id, lr.label, lr.created_at, lr.updated_at, t.id, t.account_id, t.label, t.color FROM limited_records lr LEFT JOIN mtm_tags_resume_presets ta ON lr.id = ta.record_id LEFT JOIN tags t ON ta.tag_id = t.id ORDER BY lr.updated_at DESC",
			postgres.ManyRecordsWithTagsQuery[types.ResumePreset](postgres.ResumePresetsTable, postgres.MtmTagsResPre),
		),
		specificTestCase(
			"WITH limited_records AS (SELECT id, account_id, label, text, created_at, updated_at FROM resume_sections WHERE account_id = $1 OFFSET $2 LIMIT $3) SELECT lr.id, lr.account_id, lr.label, lr.text, lr.created_at, lr.updated_at, t.id, t.account_id, t.label, t.color FROM limited_records lr LEFT JOIN mtm_tags_resume_sections ta ON lr.id = ta.record_id LEFT JOIN tags t ON ta.tag_id = t.id ORDER BY lr.updated_at DESC",
			postgres.ManyRecordsWithTagsQuery[types.ResumeSection](postgres.ResumeSectionsTable, postgres.MtmTagsResSec),
		),
	)
	testSpecificQueries("UpdateRecordQuery()", specificTestCases, t)
}

func TestCreateCreateManyRecordsWithTags(t *testing.T) {
	t.Parallel()

	expectedQuery := `WITH new_record AS (INSERT INTO table (account_id, test_field1, test_field2) VALUES ($1, $2, $3) RETURNING *), _ AS (INSERT INTO mtm_table (tag_id, record_id) VALUES ($4, (SELECT id FROM new_record))) SELECT cs.* FROM new_record cs`

	testQueries("CreateCreateManyRecordsWithTags()", []QueryTestCase{
		{expectedQuery: expectedQuery},
	}, postgres.CreateCreateSingleWithTags[TestQueryType]("table", "mtm_table")(1), t)

	type queryFunc = func(int) string

	type TestCaseCreateWithTags struct {
		expectedPatterns []string
		queryFunc        queryFunc
		offset           int
	}

	expectedPatters := func(args ...string) []string {
		return args
	}

	generateTagValues := func(tagCount int, start int) string {
		sections := make([]string, tagCount)
		for i := 0; i < tagCount; i++ {
			sections[i] = fmt.Sprintf("($%d, (SELECT id FROM new_record))", i+start+1)
		}
		return strings.Join(sections, ", ")
	}

	testCases := []TestCaseCreateWithTags{
		{
			expectedPatterns: expectedPatters("INSERT INTO resume_presets (account_id, label)"),
			queryFunc:        postgres.CreateCreateSingleWithTags[types.ResumePreset](postgres.ResumePresetsTable, postgres.MtmTagsResPre),
			offset:           len(postgres.GetDbFieldsForCreate[types.ResumePreset]()),
		},
		{
			expectedPatterns: expectedPatters("INSERT INTO resume_sections (account_id, label, text)"),
			queryFunc:        postgres.CreateCreateSingleWithTags[types.ResumeSection](postgres.ResumeSectionsTable, postgres.MtmTagsResPre),
			offset:           len(postgres.GetDbFieldsForCreate[types.ResumeSection]()),
		},
		{
			expectedPatterns: expectedPatters("INSERT INTO application_presets (account_id, label)"),
			queryFunc:        postgres.CreateCreateSingleWithTags[types.ApplicationPreset](postgres.ApplicationPresetsTable, postgres.MtmTagsResPre),
			offset:           len(postgres.GetDbFieldsForCreate[types.ApplicationPreset]()),
		},
		{
			expectedPatterns: expectedPatters("INSERT INTO application_sections (account_id, label, text)"),
			queryFunc:        postgres.CreateCreateSingleWithTags[types.ApplicationSection](postgres.ApplicationSectionsTable, postgres.MtmTagsResPre),
			offset:           len(postgres.GetDbFieldsForCreate[types.ApplicationSection]()),
		},
	}

	for _, tc := range testCases {
		for tagCount := 0; tagCount < 5; tagCount++ {
			generatedQuery := utils.TrimInside(tc.queryFunc(tagCount))

			for patternIndex, pattern := range tc.expectedPatterns {
				if !strings.Contains(generatedQuery, pattern) {
					t.Fatalf("\ngenerated query does not contain pattern\nquery:   %s\npattern: %s\nwith index %d",
						generatedQuery, pattern, patternIndex,
					)
				}
			}

			if tagValues := generateTagValues(tagCount, tc.offset); !strings.Contains(generatedQuery, tagValues) {
				t.Fatalf("\ngenerated query: %s\ndoes not contain %d tag values\nvalues: %s",
					generatedQuery, tagCount, tagValues,
				)
			}
		}
	}
}
