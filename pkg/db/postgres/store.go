package postgres

import (
	"database/sql"
	"log"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	Db *sql.DB

	ApplicationSections StoreWithTags[types.ApplicationSection]
	ApplicationPresets  StoreWithTags[types.ApplicationPreset]
	ResumeSections      StoreWithTags[types.ResumeSection]
	ResumePresets       StoreWithTags[types.ResumePreset]
}

func NewPostgresStore(connectionString string) *PostgresStore {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Could not create postgres connection")
	}
	return CreatePostgresStore(db)
}

func CreatePostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		Db:                  db,
		ApplicationSections: CreateGenericStoreWithTags[types.ApplicationSection](db, ApplicationSectionsTable, MtmTagsAppSec),
		ApplicationPresets:  CreateGenericStoreWithTags[types.ApplicationPreset](db, ApplicationPresetsTable, MtmTagsAppPre),
		ResumeSections:      CreateGenericStoreWithTags[types.ResumeSection](db, ResumeSectionsTable, MtmTagsResSec),
		ResumePresets:       CreateGenericStoreWithTags[types.ResumePreset](db, ResumePresetsTable, MtmTagsResPre),
	}
}

type StoreWithTags[T any] interface {
	GetMany(accountId int, pagination types.PaginationParams) ([]T, error)
	GetSingle(accountId int, recordId int) (T, error)
	CreateSingle(body T) (T, error)
	UpdateSingle(body T) (T, error)
	DeleteSingle(accountId int, recordId int) error
}
