package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
)

type GenericStoreWithTags[T CanSetTags] struct {
	db          *sql.DB
	scan        func(scannable db.Scannable) (T, error)
	scanWithTag func(scannable db.Scannable) (T, types.Tag, error)
	queries     QueryHolderWithCreateTags
}

func CreateGenericStoreWithTags[T CanSetTags](db *sql.DB, table string, mtmTable string) StoreWithTags[T] {
	return GenericStoreWithTags[T]{
		db:          db,
		scan:        createScanFunc[T](),
		scanWithTag: createScanWithTagsFunc[T](),
		queries:     createQueryHolderWithTags[T](table, mtmTable),
	}
}

func (s GenericStoreWithTags[T]) GetMany(accountId int, pagination types.PaginationParams) ([]T, error) {
	return genericGetRecordsWithTags(
		s.db, s.scanWithTag,
		s.queries.queryManyWithTags,
		accountId, pagination)
}

func (s GenericStoreWithTags[T]) GetSingle(accountId int, recordId int) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.querySingle,
		accountId, recordId,
	)
}

func (s GenericStoreWithTags[T]) CreateSingle(body T) (T, error) {
	var value T
	tx, err := s.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return value, err
	}

	query := s.queries.createWithTags(body.GetTagCount())
	row := tx.QueryRow(query, utils.GetValuesFromBody(body)...)
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

func (s GenericStoreWithTags[T]) UpdateSingle(body T) (T, error) {
	return WithTransactionScan(
		s.db, getRecord, s.scan,
		s.queries.updateSingle,
		utils.GetValuesFromBody(body)...,
	)
}

func (s GenericStoreWithTags[T]) DeleteSingle(accountId int, recordId int) error {
	return WithTransaction(
		s.db, deleteRecord,
		s.queries.deleteSingle,
		accountId, recordId,
	)
}
