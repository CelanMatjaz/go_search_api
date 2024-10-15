package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

type GenericStoreWithTags[T CanSetTags, bodyT any] struct {
	db          *sql.DB
	scan        func(scannable db.Scannable) (T, error)
	scanWithTag func(scannable db.Scannable) (T, types.Tag, error)
	queries     QueryHolderWithCreateTags
}

func createGenericStoreWithTags[T CanSetTags, bodyT any](db *sql.DB, table string, mtmTable string) StoreWithTags[T, bodyT] {
	return GenericStoreWithTags[T, bodyT]{
		db:          db,
		scan:        createScanFunc[T](),
		scanWithTag: createScanWithTagsFunc[T](),
		queries:     createQueryHolderWithTags[T, bodyT](table, mtmTable),
	}
}

func (s GenericStoreWithTags[T, bodyT]) GetMany(accountId int, pagination types.PaginationParams) ([]T, error) {
	return genericGetRecordsWithTags(
		s.db, s.scanWithTag,
		s.queries.queryManyWithTags,
		accountId, pagination)
}

func (s GenericStoreWithTags[T, bodyT]) GetSingle(accountId int, recordId int) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.querySingle,
		accountId, recordId,
	)
}

func (s GenericStoreWithTags[T, bodyT]) CreateSingle(accountId int, body bodyT) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.createSingle,
		accountId, body,
	)
}
func (s GenericStoreWithTags[T, bodyT]) UpdateSingle(accountId int, recordId int, body bodyT) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.updateSingle,
		accountId, body,
	)
}

func (s GenericStoreWithTags[T, bodyT]) DeleteSingle(accountId int, recordId int) error {
	return WithTransaction(
		s.db, deleteRecord,
		s.queries.deleteSingle,
		accountId, recordId,
	)
}
