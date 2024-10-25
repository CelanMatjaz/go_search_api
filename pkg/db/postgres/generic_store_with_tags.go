package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
)

type GenericStoreWithTags[T HasIdAndTags] struct {
	db      *sql.DB
	scan    func(scannable db.Scannable) (T, error)
	queries QueryHolderWithCreateTags
}

func CreateGenericStoreWithTags[T HasIdAndTags](db *sql.DB, table string, mtmTable string) DefaultStoreWithTags[T] {
	return GenericStoreWithTags[T]{
		db:      db,
		scan:    createScanFunc[T](),
		queries: createQueryHolderWithTags[T](table, mtmTable),
	}
}

func (s GenericStoreWithTags[T]) GetMany(accountId int, pagination types.PaginationParams) ([]types.RecordWithTags[T], error) {
	return genericGetRecordsWithTags[T](
		s.db, s.queries.queryManyWithTags,
		accountId, pagination)
}

func (s GenericStoreWithTags[T]) GetSingle(accountId int, recordId int) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.querySingle,
		accountId, recordId,
	)
}

func (s GenericStoreWithTags[T]) CreateSingle(accountId int, body T) (T, error) {
	var value T
	tx, err := s.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return value, err
	}

	query := s.queries.createWithTags(body.GetTagCount())
	row := tx.QueryRow(query, utils.GetValuesFromBody(body, []any{accountId})...)
	value, err = s.scan(row)
	if err != nil {
		return value, err
	}

	err = tx.Commit()
	if err != nil {
		return value, err
	}

	return value, err
}

func (s GenericStoreWithTags[T]) UpdateSingle(accountId int, recordId int, body T) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.updateSingle,
		utils.GetValuesFromBody(body, []any{recordId, accountId})...,
	)
}

func (s GenericStoreWithTags[T]) DeleteSingle(accountId int, recordId int) error {
	return WithTransaction(
		s.db, deleteRecord,
		s.queries.deleteSingle,
		recordId, accountId,
	)
}
