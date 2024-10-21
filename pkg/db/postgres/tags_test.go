package postgres_test

import (
	"database/sql"
	"math/rand/v2"
	"reflect"
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db/postgres"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func TestTags(t *testing.T) {
	db, conn := createDbAndStore()
	defer cleanupDb(db)

	store := postgres.CreatePostgresStore(conn.Db)
	account := seedAccount(t, store)

	tags := []types.Tag{
		types.CreateTag(account.Id, "label1", "#AAAAAA"),
		types.CreateTag(account.Id, "label2", "#BBBBBB"),
		types.CreateTag(account.Id, "label3", "#CCCCCC"),
		types.CreateTag(account.Id, "label4", "#DDDDDD"),
		types.CreateTag(account.Id, "label5", "#EEEEEE"),
		types.CreateTag(account.Id, "label6", "#FFFFFF"),
	}

	createdTags := make([]types.Tag, len(tags))

	for i, tag := range tags {
		newTag, err := store.CreateTag(account.Id, tag)
		if err != nil {
			t.Fatalf("Could not create new tag, %s", err.Error())
		}
		createdTags[i] = newTag
	}

	randomTagIndex := rand.IntN(len(createdTags))
	randomTag := createdTags[randomTagIndex]
	randomTag.Color = "#000000"
	randomTag.Label = "updated label"

	updatedTag, err := store.UpdateTag(account.Id, int(randomTag.AccountId), randomTag)
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
		t.Fatalf("Tag was not deleted")
	}

	updatedTag, err = store.UpdateTag(account.Id, int(randomTag.AccountId), randomTag)
	if err != sql.ErrNoRows {
		t.Fatalf("Wrong error when updating non existing tag")
	}

	err = store.DeleteTag(account.Id, int(singleTag.Id))
	if err != nil {
		t.Fatalf("Wrong error when deleting non existing tag")
	}
}

func checkTagField(t *testing.T, tag string, expected string, value string) {
	if value != expected {
		t.Fatalf("%s value is not correct\nexpected: %s\nactual:   %s", tag, expected, value)
	}
}
