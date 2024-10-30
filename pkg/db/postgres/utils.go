package postgres

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
)

func GetDbFieldsSelect[T any]() []string {
	var val T
	v := reflect.ValueOf(&val).Elem()
	return getDbTagFields(v, "select")
}

func GetDbFieldsForCreate[T any]() []string {
	var val T
	v := reflect.ValueOf(&val).Elem()
	return getDbTagFields(v, "create")
}

func GetDbFieldsForUpdate[T any]() []string {
	var val T
	v := reflect.ValueOf(&val).Elem()
	return getDbTagFields(v, "update")
}

func getDbTagFields(value reflect.Value, bodyTagValue string) []string {
	if bodyTagValue == "" {
		panic("Body tag value should not be an empty string")
	}

	fields := make([]string, 0)

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := value.Type().Field(i)

		dbTag, ok := fieldType.Tag.Lookup("db")
		if field.Kind() == reflect.Struct && !ok {
			fields = append(fields, getDbTagFields(field, bodyTagValue)...)
			continue
		}

		if dbTag == "" {
			continue
		}

		bodyTag := fieldType.Tag.Get("body")
		if bodyTagValue == "select" {
			goto append
		}

		if bodyTag == "omit" {
			continue
		} else if bodyTag != "" && !strings.Contains(bodyTag, bodyTagValue) {
			continue
		}

	append:
		fields = append(fields, dbTag)
	}

	return fields
}

func GetScanFields[T any](val *T) []any {
	if reflect.TypeOf(*val).Kind() != reflect.Struct {
		panic("invalid type provided")
	}

	v := reflect.ValueOf(val).Elem()
	return getScanFields(v)
}

func getScanFields(v reflect.Value) []any {
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

		fields = append(fields, getScanFields(field)...)
	}

	return fields
}

func createScanFunc[T any]() func(scannable db.Scannable) (T, error) {
	var temp T
	fields := GetScanFields(&temp)

	return func(scannable db.Scannable) (T, error) {
		if err := scannable.Scan(fields...); err != nil {
			return temp, err
		}
		return temp, nil
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
