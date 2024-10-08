package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func getRecords[T any](s *PostgresStore, scan func(db.Scannable) (*T, error), query string, args ...any) ([]T, error) {
	rows, err := s.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	records := make([]T, 0)
	for rows.Next() {
		record, err := scan(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, *record)
	}

	return records, nil
}

func getRecord[T any](s *PostgresStore, scan func(db.Scannable) (*T, error), query string, args ...any) (*T, error) {
	row := s.Db.QueryRow(query, args...)
	record, err := scan(row)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createRecord[T any](tx *sql.Tx, scan func(db.Scannable) (*T, error), query string, args ...any) (*T, error) {
	row := tx.QueryRow(query, args...)
	record, err := scan(row)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func createRecordWithTags[T any](tx *sql.Tx, scan func(db.Scannable) (*T, error), query string, args ...any) (*T, error) {
	row := tx.QueryRow(query, args...)
	record, err := scan(row)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func updateRecord[T any](tx *sql.Tx, scan func(db.Scannable) (*T, error), query string, args ...any) (*T, error) {
	return createRecord(tx, scan, query, args...)
}

func deleteRecord(tx *sql.Tx, query string, args ...any) error {
	_, err := tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

type CanSetTags interface {
	GetId() int
	AppendTag(newTag *types.Tag)
}

func genericGetRecordsWithTags[T CanSetTags](
	s *PostgresStore,
	scan func(db.Scannable) (T, *types.Tag, error),
	query string, accountId int,
	pagination types.PaginationParams,
) ([]T, error) {
	rows, err := s.Db.Query(query, accountId, pagination.GetOffset(), pagination.Count)
	if err != nil {
		return nil, err
	}

	recordSet := make(map[int]int)
	recordArr := make([]T, 0)

	for rows.Next() {
		scannedRecord, tag, err := scan(rows)

		if err != nil {
			return nil, err
		}

		recordIndex, exists := recordSet[scannedRecord.GetId()]
		if !exists {
			recordSet[scannedRecord.GetId()] = len(recordArr)
			recordIndex = len(recordArr)
			recordArr = append(recordArr, scannedRecord)
		}

		if tag.Id.Valid {
			(recordArr[recordIndex]).AppendTag(tag)
		}
	}

	return recordArr, nil
}
