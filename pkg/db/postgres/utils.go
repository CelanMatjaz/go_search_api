package postgres

import (
	"database/sql"
	"reflect"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

// Function assumes that there is only 1 level of embedding
func reflectDbFields[T any]() []string {
	var val T
	v := reflect.ValueOf(&val).Elem()
	values := make([]string, 0)

	if v.Kind() != reflect.Struct {
		panic("invalid type provided")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		tag := fieldType.Tag.Get("db")
		if tag != "" {
			values = append(values, tag)
			continue
		}

		if fieldType.Anonymous && field.Kind() == reflect.Struct {
			for j := 0; j < field.NumField(); j++ {
				embeddedFieldType := field.Type().Field(j)

				embeddedTag := embeddedFieldType.Tag.Get("db")
				if embeddedTag != "" {
					values = append(values, embeddedTag)
				}
			}
		}
	}

	return values
}

// Function assumes that there is only 1 level of embedding
func getScanFields[T any](val *T) []any {
	v := reflect.ValueOf(val).Elem()

	if reflect.TypeOf(*val).Kind() != reflect.Struct {
		panic("invalid type provided")
	}

	fields := make([]any, 0)
	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		tag := fieldType.Tag.Get("db")
		if tag != "" && field.CanAddr() {
			fields = append(fields, field.Addr().Interface())
			continue
		}

		if !fieldType.Anonymous || field.Kind() != reflect.Struct {
			continue
		}

		for j := range field.NumField() {
			embeddedField := field.Field(j)
			embeddedFieldType := field.Type().Field(j)

			embeddedTag := embeddedFieldType.Tag.Get("db")
			if embeddedTag != "" && embeddedField.CanAddr() {
				fields = append(fields, embeddedField.Addr().Interface())
			}
		}
	}

	return fields
}

func createScanFunc[T any]() func(scannable db.Scannable) (T, error) {
	var temp T
	fields := getScanFields(&temp)

	return func(scannable db.Scannable) (T, error) {
		if err := scannable.Scan(fields...); err != nil {
			return temp, err
		}
		return temp, nil
	}
}

func createScanWithTagsFunc[T any]() func(scannable db.Scannable) (T, types.Tag, error) {
	var temp T
	fields := getScanFields(&temp)

	var tempTag types.Tag
	fields = append(fields, getScanFields(&tempTag)...)

	return func(scannable db.Scannable) (T, types.Tag, error) {
		if err := scannable.Scan(fields...); err != nil {
			return temp, tempTag, err
		}
		return temp, tempTag, nil
	}
}

// ???
func WithTransaction(
	db *sql.DB,
	fn func(tx *sql.Tx, query string, args ...any) error,
	query string, args ...any,
) error {
	tx, err := db.Begin()
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
	db *sql.DB,
	fn func(tx *sql.Tx, scan func(db.Scannable) (T, error), query string, args ...any) (T, error),
	scan func(db.Scannable) (T, error), query string, args ...any,
) (T, error) {
	var temp T
	tx, err := db.Begin()
	if err != nil {
		return temp, err
	}
	defer tx.Rollback()

	record, err := fn(tx, scan, query, args...)
	if err != nil {
		return temp, err
	}

	if err = tx.Commit(); err != nil {
		return temp, err
	}

	return record, nil
}
