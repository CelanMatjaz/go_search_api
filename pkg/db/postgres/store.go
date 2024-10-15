package postgres

import (
	"database/sql"
	"log"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	Db *sql.DB

	ApplicationSections StoreWithTags[types.ApplicationSection, types.ApplicationSectionBody]
	ApplicationPresets  StoreWithTags[types.ApplicationPreset, types.ApplicationPresetBody]
	ResumeSections      StoreWithTags[types.ResumeSection, types.ResumeSectionBody]
	ResumePresets       StoreWithTags[types.ResumePreset, types.ResumePresetBody]
}

func CreatePostgresStore(connectionString string) *PostgresStore {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Could not create postgres connection")
	}
	return &PostgresStore{
		Db:                  db,
		ApplicationSections: createGenericStoreWithTags[types.ApplicationSection, types.ApplicationSectionBody](db, application_sections_table, mtm_tags_app_sec),
		ApplicationPresets:  createGenericStoreWithTags[types.ApplicationPreset, types.ApplicationPresetBody](db, application_presets_table, mtm_tags_app_pre),
		ResumeSections:      createGenericStoreWithTags[types.ResumeSection, types.ResumeSectionBody](db, resume_sections_table, mtm_tags_res_sec),
		ResumePresets:       createGenericStoreWithTags[types.ResumePreset, types.ResumePresetBody](db, resume_presets_table, mtm_tags_res_pre),
	}
}

type StoreWithTags[T any, bodyT any] interface {
	GetMany(accountId int, pagination types.PaginationParams) ([]T, error)
	GetSingle(accountId int, recordId int) (T, error)
	CreateSingle(accountId int, body bodyT) (T, error)
	UpdateSingle(accountId int, recordId int, body bodyT) (T, error)
	DeleteSingle(accountId int, recordId int) error
}
