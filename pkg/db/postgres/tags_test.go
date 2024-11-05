package postgres_test

import (
	"database/sql"
	"fmt"
	"math/rand/v2"
	"reflect"
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db/postgres"
	testcommon "github.com/CelanMatjaz/job_application_tracker_api/pkg/test_common"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func TestTags(t *testing.T) {
	conn := testcommon.CreateStore(t)
	store := postgres.CreatePostgresStore(conn.Db)
	account, _ := testcommon.SeedAccount(t, store)

	createdTags, _ := createTags(t, store, account.Id)

	randomTagIndex := rand.IntN(len(createdTags))
	randomTag := createdTags[randomTagIndex]
	randomTag.Color = "#000000"
	randomTag.Label = "updated label1"

	updatedTag, err := store.UpdateTag(randomTag.Id, account.Id, randomTag)
	if err != nil {
		t.Fatalf("Could not update tag, %s", err.Error())
	}
	if !reflect.DeepEqual(randomTag, updatedTag) {
		t.Fatalf("Updated tag does not match expected tag\nexpected: %v\nupdated:  %v", randomTag, updatedTag)
	}

	singleTag, err := store.GetTag(account.Id, int(randomTag.Id))
	if err != nil {
		t.Fatalf("Could not query tag, %s", err.Error())
	}
	if !reflect.DeepEqual(randomTag, singleTag) {
		t.Fatalf("Queried tag does not match expected tag\nexpected: %v\nupdated:  %v", randomTag, singleTag)
	}

	err = store.DeleteTag(account.Id, int(singleTag.Id))
	if err != nil {
		t.Fatalf("Could not delete tag, %s", err.Error())
	}

	singleTag, err = store.GetTag(account.Id, randomTag.Id)
	if err != sql.ErrNoRows {
		t.Fatalf("Wrong error when deleting non existing tag")
	}

	updatedTag, err = store.UpdateTag(account.Id, randomTag.Id, randomTag)
	if err != sql.ErrNoRows {
		t.Fatalf("Wrong error when updating non existing tag")
	}

	err = store.DeleteTag(account.Id, int(singleTag.Id))
	if err != nil {
		t.Fatalf("Wrong error when deleting non existing tag")
	}
}

func TestTagAssociations(t *testing.T) {
	conn := testcommon.CreateStore(t)

	store := postgres.CreatePostgresStore(conn.Db)
	account, _ := testcommon.SeedAccount(t, store)

	_, tagIds := createTags(t, store, account.Id)
	createdRecords := make([]types.ResumePreset, 0)

	recordCount := len(tagIds)
	tagAssociationCount := 0
	for i := range recordCount {
		tagIdsPerm := make([]int, 0)
		if i%3 == 0 {
			goto skip
		}
		for i := range rand.Perm(i) {
			tagIdsPerm = append(tagIdsPerm, tagIds[i])
		}

	skip:
		newRecord := types.ResumePreset{
			Label: fmt.Sprintf("label%d", i),
			WithTags: &types.WithTags{
				TagIds: tagIdsPerm,
			},
		}

		tagAssociationCount += len(newRecord.TagIds)

		newRecord, err := store.ResumePresets.CreateSingle(account.Id, newRecord)
		if err != nil {
			t.Fatalf("Error creating record, %s", err.Error())
		}
		createdRecords = append(createdRecords, newRecord)
	}

	queryAssociationCount := func(recordId int) int {
		var row *sql.Row
		if recordId > 0 {
			row = conn.Db.QueryRow("SELECT COUNT(*) FROM mtm_tags_resume_presets WHERE record_id = $1")
		} else {
			row = conn.Db.QueryRow("SELECT COUNT(*) FROM mtm_tags_resume_presets")
		}

		var count int
		err := row.Scan(&count)
		if err != nil {
			t.Fatalf("Error scanning tag associations, %s", err.Error())
		}
		return count
	}

	t.Run("association count", func(t *testing.T) {
		if associationCount := queryAssociationCount(0); associationCount != tagAssociationCount {
			t.Fatalf("Generated association count does not match count from db\nexpected: %d\nactual:   %d", tagAssociationCount, associationCount)
		}
	})

	t.Run("record count", func(t *testing.T) {
		records, err := store.ResumePresets.GetMany(account.Id, types.DefaultPagaintion())
		if err != nil {
			t.Fatalf("Could not query many records, %s", err.Error())
		}

		if length := len(records); length != recordCount {
			t.Fatalf("Queried record count does not equal expected count\nexpected: %d\nactual:   %d", recordCount, length)
		}
	})

	t.Run("check associated tag count", func(t *testing.T) {
		records, err := store.ResumePresets.GetMany(account.Id, types.DefaultPagaintion())
		if err != nil {
			t.Fatalf("Could not query many records, %s", err.Error())
		}

		for _, record := range records {
			row := conn.Db.QueryRow("SELECT COUNT(*) FROM mtm_tags_resume_presets WHERE record_id = $1", record.Record.Id)
			var queriedAssociationCount int
			err = row.Scan(&queriedAssociationCount)
			if err != nil {
				t.Fatalf("Error scanning tag associations, %s", err.Error())
			}

			if tagCount := len(record.Tags); tagCount != queriedAssociationCount {
				t.Fatalf("Queried associated tag count does not equal expected count\nexpected: %d\nactual:   %d", tagCount, queriedAssociationCount)
			}
		}
	})

	t.Run("check tag association count for deleted records", func(t *testing.T) {
		records, err := store.ResumePresets.GetMany(account.Id, types.DefaultPagaintion())
		if err != nil {
			t.Fatalf("Could not query many records, %s", err.Error())
		}

		for _, record := range records {
			err := store.ResumePresets.DeleteSingle(account.Id, record.Record.Id)
			if err != nil {
				t.Fatalf("Could not delete record, %s", err.Error())
			}

			row := conn.Db.QueryRow("SELECT COUNT(*) FROM mtm_tags_resume_presets WHERE record_id = $1", record.Record.Id)
			var queriedAssociationCount int
			err = row.Scan(&queriedAssociationCount)
			if err != nil {
				t.Fatalf("Error scanning tag associations, %s", err.Error())
			}

			if queriedAssociationCount != 0 {
				t.Fatalf("Queried associated tag count does not equal expected count\nexpected: %d\nactual:   %d", 0, queriedAssociationCount)
			}
		}

		if count := queryAssociationCount(0); count != 0 {
			t.Fatalf("Tag associations were not deleted after deleting all records\nexpected: %d\nactual:   %d", 0, count)
		}
	})
}

func createTags(t *testing.T, store *postgres.PostgresStore, accountId int) ([]types.Tag, []int) {
	tags := []types.Tag{
		types.CreateTag(accountId, "label1", "#AAAAAA"),
		types.CreateTag(accountId, "label2", "#BBBBBB"),
		types.CreateTag(accountId, "label3", "#CCCCCC"),
		types.CreateTag(accountId, "label4", "#DDDDDD"),
		types.CreateTag(accountId, "label5", "#EEEEEE"),
		types.CreateTag(accountId, "label6", "#FFFFFF"),
	}

	createdTags := make([]types.Tag, len(tags))
	tagIds := make([]int, len(tags))

	for i, tag := range tags {
		newTag, err := store.CreateTag(accountId, tag)
		if err != nil {
			t.Fatalf("Could not create new tag, %s", err.Error())
		}
		createdTags[i] = newTag
		tagIds[i] = newTag.Id
	}

	return createdTags, tagIds
}

func checkTagField(t *testing.T, tag string, expected string, value string) {
	if value != expected {
		t.Fatalf("%s value is not correct\nexpected: %s\nactual:   %s", tag, expected, value)
	}
}
