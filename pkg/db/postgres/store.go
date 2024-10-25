package postgres

import (
	"database/sql"
	"log"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	Db *sql.DB

	ApplicationSections DefaultStoreWithTags[types.ApplicationSection]
	ApplicationPresets  DefaultStoreWithTags[types.ApplicationPreset]
	ResumeSections      DefaultStoreWithTags[types.ResumeSection]
	ResumePresets       DefaultStoreWithTags[types.ResumePreset]
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

type DefaultStoreCommon[T any] interface {
	GetSingle(accountId int, recordId int) (T, error)
	CreateSingle(accountId int, body T) (T, error)
	UpdateSingle(accountId int, recordId int, body T) (T, error)
	DeleteSingle(accountId int, recordId int) error
}

type DefaultStore[T any] interface {
	GetMany(accountId int, pagination types.PaginationParams) ([]T, error)
	DefaultStoreCommon[T]
}

type DefaultStoreWithTags[T any] interface {
	GetMany(accountId int, pagination types.PaginationParams) ([]types.RecordWithTags[T], error)
	DefaultStoreCommon[T]
}
