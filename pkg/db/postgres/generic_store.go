package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
)

type GenericStore[T CanSetTags] struct {
	db      *sql.DB
	scan    func(scannable db.Scannable) (T, error)
	queries QueryHolder
}

func createGenericStore[T CanSetTags](db *sql.DB, table string, hasPagination bool) StoreWithTags[T] {
	return GenericStore[T]{
		db:      db,
		scan:    createScanFunc[T](),
		queries: createQueryHolder[T](table, hasPagination),
	}
}

func (s GenericStore[T]) GetMany(accountId int, pagination types.PaginationParams) ([]T, error) {
	return getRecords(
		s.db, s.scan,
		s.queries.queryMany,
		accountId, pagination)
}

func (s GenericStore[T]) GetSingle(accountId int, recordId int) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.querySingle,
		accountId, recordId,
	)
}

func (s GenericStore[T]) CreateSingle(body T) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.createSingle,
        utils.GetValuesFromBody(body),
	)
}
func (s GenericStore[T]) UpdateSingle(body T) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.updateSingle,
        utils.GetValuesFromBody(body),
	)
}

func (s GenericStore[T]) DeleteSingle(accountId int, recordId int) error {
	return WithTransaction(
		s.db, deleteRecord,
		s.queries.deleteSingle,
		accountId, recordId,
	)
}
