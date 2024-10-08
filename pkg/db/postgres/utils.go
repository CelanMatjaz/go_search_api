package postgres

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
)

// ???
func WithTransaction(
	s *PostgresStore,
	fn func(tx *sql.Tx, query string, args ...any) error,
	query string, args ...any,
) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = fn(tx, query, args...); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// / ???
func WithTransactionScan[T any](
	s *PostgresStore,
	fn func(tx *sql.Tx, scan func(db.Scannable) (*T, error), query string, args ...any) (*T, error),
	scan func(db.Scannable) (*T, error), query string, args ...any,
) (*T, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	record, err := fn(tx, scan, query, args...)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return record, nil
}
