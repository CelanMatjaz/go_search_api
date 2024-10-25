package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func getRecords[T any](db *sql.DB, scan func(db.Scannable) (T, error), query string, args ...any) ([]T, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	records := make([]T, 0)
	for rows.Next() {
		record, err := scan(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func getRecord[T any](tx *sql.Tx, scan func(db.Scannable) (T, error), query string, args ...any) (T, error) {
	var temp T
	row := tx.QueryRow(query, args...)
	record, err := scan(row)
	if err != nil {
		return temp, err
	}

	return record, nil
}

func createRecord[T any](tx *sql.Tx, scan func(db.Scannable) (T, error), query string, args ...any) (T, error) {
	return getRecord(tx, scan, query, args...)
}

// func createRecordWithTags[T any](tx *sql.Tx, scan func(db.Scannable) (T, error), query string, args ...any) (T, error) {
// 	return getRecord(tx, scan, query, args...)
// }
//
// func updateRecord[T any](tx *sql.Tx, scan func(db.Scannable) (T, error), query string, args ...any) (T, error) {
// 	return createRecord(tx, scan, query, args...)
// }

func deleteRecord(tx *sql.Tx, query string, args ...any) error {
	_, err := tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

type HasIdAndTags interface {
	GetId() int
	GetTagCount() int
}

func genericGetRecordsWithTags[T HasIdAndTags](
	db *sql.DB,
	query string, accountId int,
	pagination types.PaginationParams,
) ([]types.RecordWithTags[T], error) {
	rows, err := db.Query(query, accountId, pagination.GetOffset(), pagination.Count)
	if err != nil {
		return nil, err
	}

	recordSet := make(map[int]int)
	recordArr := make([]types.RecordWithTags[T], 0)

	var record T
	scanFields := GetScanFields(&record)
	var tag types.ScanTag
	scanFields = append(scanFields, GetScanFields(&tag)...)

	for rows.Next() {
		err := rows.Scan(scanFields...)
		if err != nil {
			return nil, err
		}

		recordIndex, exists := recordSet[(record).GetId()]
		if !exists {
			recordSet[(record).GetId()] = len(recordArr)
			recordIndex = len(recordArr)
			recordArr = append(recordArr, types.RecordWithTags[T]{Record: record})
		}

		if tag.Id.Valid {
			recordArr[recordIndex].Tags = append(recordArr[recordIndex].Tags, tag.Tag())
		}
	}

	return recordArr, nil
}
