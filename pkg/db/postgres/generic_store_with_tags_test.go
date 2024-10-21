package postgres_test

import (
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db/postgres"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func TestGenericStore(t *testing.T) {
	db, conn := createDbAndStore()
	defer cleanupDb(db)

	store := postgres.CreatePostgresStore(conn.Db)
	account := seedAccount(t, store)

	_, err := store.CreateResumeSection(types.ResumeSection{
		AccountId: account.Id,
		Label:     "label",
		Text:      "text",
		WithTags: types.WithTags{
			TagIds: []int{},
		},
	})
	if err != nil {
		t.Fatalf("could not create resume section, %s", err.Error())
	}

	sections, err := store.GetResumeSections(account.Id, types.PaginationParams{})

	if err != nil {
		t.Fatalf("could not get resume sections, %s", err.Error())
	}

    if  len(sections) == 0 {
        t.Errorf("no sections found")
    }

	for i, s := range sections {
		t.Errorf("%d %v", i, s)
	}

}
