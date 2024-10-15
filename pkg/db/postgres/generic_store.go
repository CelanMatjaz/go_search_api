package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

type GenericStore[T CanSetTags, bodyT any] struct {
	db      *sql.DB
	scan    func(scannable db.Scannable) (T, error)
	queries QueryHolder
}

func createGenericStore[T CanSetTags, bodyT any](db *sql.DB, table string, hasPagination bool) StoreWithTags[T, bodyT] {
	return GenericStore[T, bodyT]{
		db:      db,
		scan:    createScanFunc[T](),
		queries: createQueryHolder[T, bodyT](table, hasPagination),
	}
}

func (s GenericStore[T, bodyT]) GetMany(accountId int, pagination types.PaginationParams) ([]T, error) {
	return getRecords(
		s.db, s.scan,
		s.queries.queryMany,
		accountId, pagination)
}

func (s GenericStore[T, bodyT]) GetSingle(accountId int, recordId int) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.querySingle,
		accountId, recordId,
	)
}

func (s GenericStore[T, bodyT]) CreateSingle(accountId int, body bodyT) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.createSingle,
		accountId, body,
	)
}
func (s GenericStore[T, bodyT]) UpdateSingle(accountId int, recordId int, body bodyT) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.updateSingle,
		accountId, body,
	)
}

func (s GenericStore[T, bodyT]) DeleteSingle(accountId int, recordId int) error {
	return WithTransaction(
		s.db, deleteRecord,
		s.queries.deleteSingle,
		accountId, recordId,
	)
}
