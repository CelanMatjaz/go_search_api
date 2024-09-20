package postgres

import "github.com/CelanMatjaz/job_application_tracker_api/pkg/db"

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

func createRecord[T any](s *PostgresStore, scan func(db.Scannable) (*T, error), query string, args ...any) (*T, error) {
	return getRecord[T](s, scan, query, args)
}

func updateRecord[T any](s *PostgresStore, scan func(db.Scannable) (*T, error), query string, args ...any) (*T, error) {
	return getRecord[T](s, scan, query, args)
}

func deleteRecord(s *PostgresStore, query string, args ...any) error {
	_, err := s.Db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
