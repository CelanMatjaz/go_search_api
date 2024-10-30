package postgres_test

import (
	"reflect"
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db/postgres"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func TestGetDbFieldsSelect(t *testing.T) {
	type Struct struct {
		Field1 any `db:"field1" body:"select,update,create"`
		Field2 any `db:"field2"`
		Field3 any `db:"" body:""`
		Field4 any `d:"field4"`
	}

	expected := []string{"field1", "field2"}
	generated := postgres.GetDbFieldsSelect[Struct]()

	if !reflect.DeepEqual(expected, generated) {
		t.Fatalf("generated fields are not equal to expected fields\nexpected:  %s\ngenerated: %s", expected, generated)
	}
}

func TestGetDbFieldsForCreate(t *testing.T) {
	type Struct struct {
		types.WithId
		types.WithAccountId
		Field1 any `db:"field1" body:"update,create"`
		Field2 any `db:"field2"`
		Field3 any `db:"" body:""`
		Field4 any `d:"field4"`
		Field5 any `db:"field5" body:"create"`
		Field6 any `db:"field6" body:"update"`
		types.WithTimestamps
	}

	expected := []string{"account_id", "field1", "field2", "field5"}
	generated := postgres.GetDbFieldsForCreate[Struct]()

	if !reflect.DeepEqual(expected, generated) {
		t.Fatalf("generated fields are not equal to expected fields\nexpected:  %s\ngenerated: %s", expected, generated)
	}
}

func TestGetDbFieldsForUpdate(t *testing.T) {
	type Struct struct {
		types.WithId
		types.WithAccountId
		Field1 any `db:"field1" body:"update,create"`
        Field2 any `db:"field2" body:"omit"`
		Field3 any `db:"" body:""`
		Field4 any `d:"field4"`
		Field5 any `db:"field5" body:"create"`
		Field6 any `db:"field6" body:"update"`
		types.WithTimestamps
	}

	expected := []string{"field1", "field6", "updated_at"}
	generated := postgres.GetDbFieldsForUpdate[Struct]()

	if !reflect.DeepEqual(expected, generated) {
		t.Fatalf("generated fields are not equal to expected fields\nexpected:  %s\ngenerated: %s", expected, generated)
	}
}

func TestGetScanFields(t *testing.T) {
	type Struct struct {
		types.WithId
		types.WithAccountId
		Field1 any `db:"field1" body:"update,create"`
		Field2 any `db:"field2"`
		Field3 any `db:"" body:""`
		Field4 any `d:"field4"`
		Field5 any `db:"field5" body:"create"`
		Field6 any `db:"field6" body:"update"`
		types.WithTimestamps
	}

	// var value Struct
	// fields := postgres.GetScanFields(&value)

	// t.Fatalf("fields count %d", len(fields))

}
