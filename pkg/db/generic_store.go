package db

import (
	"database/sql"
)

type Scannable interface {
	Scan(dest ...any) error
}

type ScannableFunction[T any] interface {
	Scan(row Scannable) (T, error)
}

type GenericStoreFunctions[T any] interface {
	GetRecords(args ...any) ([]T, error)
	GetRecord(args ...any) (T, error)
	CreateRecord(args ...any) (T, error)
	UpdateRecord(args ...any) (T, error)
	DeleteRecord(id int) error
}

type GenericStore[T any] struct {
	GenericStoreFunctions[T]
	Scanner ScannableFunction[T]

	SelectManyQuery string
	SelectQuery     string
	CreateQuery     string
	UpdateQuery     string
	DeleteQuery     string

	Db *sql.DB
}

func (s *GenericStore[T]) GetRecords(args ...any) ([]T, error) {
	rows, err := s.Db.Query(s.SelectManyQuery, args...)
	if err != nil {
		return []T{}, err
	}

	var records = make([]T, 0)
	for rows.Next() {
		r, _ := s.Scanner.Scan(rows)
		records = append(records, r)
	}

	return records, nil
}

func (s *GenericStore[T]) GetRecord(args ...any) (T, error) {
	row := s.Db.QueryRow(s.SelectQuery, args...)
	record, err := s.Scanner.Scan(row)
	return record, err
}

func (s *GenericStore[T]) CreateRecord(args ...any) (T, error) {
	row := s.Db.QueryRow(s.CreateQuery, args...)
	record, err := s.Scanner.Scan(row)
	return record, err
}

func (s *GenericStore[T]) UpdateRecord(args ...any) (T, error) {
	row := s.Db.QueryRow(s.UpdateQuery, args...)
	record, err := s.Scanner.Scan(row)
	return record, err
}

func (s *GenericStore[T]) DeleteRecord(id int) error {
	_, err := s.Db.Exec(s.DeleteQuery, id)
	return err
}
