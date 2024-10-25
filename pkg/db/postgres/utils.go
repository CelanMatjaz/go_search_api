package postgres

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
)

func GetDbFields[T any]() []string {
	var val T
	v := reflect.ValueOf(&val).Elem()
	return getDbTagValues(v, false, "")
}

func GetDbFieldsForCreate[T any]() []string {
	var val T
	v := reflect.ValueOf(&val).Elem()
	return getDbTagValues(v, true, "body")
}

func GetDbFieldsForUpdate[T any]() []string {
	var val T
	v := reflect.ValueOf(&val).Elem()
	return getDbTagValues(v, false, "body")
}

func getDbTagValues(value reflect.Value, omit bool, checkForTag string) []string {
	fields := make([]string, 0)

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := value.Type().Field(i)

		tag, ok := fieldType.Tag.Lookup("db")
		if field.Kind() == reflect.Struct && !ok {
			fields = append(fields, getDbTagValues(field, omit, checkForTag)...)
			continue
		}

		if checkForTag != "" {
			_, ok := fieldType.Tag.Lookup(checkForTag)
			if !ok {
				continue
			}
		}

		if !ok || tag == "" || (omit && (strings.Contains("id created_at updated_at", tag))) {
			continue
		}

		fields = append(fields, tag)
	}

	return fields
}

func GetScanFields[T any](val *T) []any {
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
	fields := GetScanFields(&temp)

	return func(scannable db.Scannable) (T, error) {
		if err := scannable.Scan(fields...); err != nil {
			return temp, err
		}
		return temp, nil
	}
}

// func createScanWithTagsFunc[T any]() func(scannable db.Scannable) (*T, types.Tag, error) {
// 	var temp T
// 	fields := getScanFields(&temp)
//
// 	var tempTag types.ScanTag
// 	fields = append(fields, getScanFields(&tempTag)...)
//
// 	return func(scannable db.Scannable) (*T, types.Tag, error) {
// 		if err := scannable.Scan(fields...); err != nil {
// 			return &temp, tempTag.Tag(), err
// 		}
//
// 		return &temp, tempTag.Tag(), nil
// 	}
// }

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
