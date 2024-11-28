package postgres_test

import (
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db/postgres"
	testcommon "github.com/CelanMatjaz/job_application_tracker_api/pkg/test_common"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func TestGenericStoreWithTags(t *testing.T) {
	conn := testcommon.CreateStore(t)

	store := postgres.CreatePostgresStore(conn.Db)
	account, _ := testcommon.SeedAccount(t, store)

	_, err := store.ResumeSections.CreateSingle(account.Id, types.ResumeSection{
		Label:    "label",
		Text:     "text",
		WithTags: &types.WithTags{},
	})
	if err != nil {
		t.Fatalf("could not create resume section, %s", err.Error())
	}

	sections, err := store.GetResumeSections(account.Id, types.DefaultPagaintion())
	if err != nil {
		t.Fatalf("could not query resume sections, %s", err.Error())
	}

	_ = sections

	// for _, s := range sections {
	// 	t.Logf("%v\n", s)
	// }
	//
	// if len(sections) == 0 {
	// 	t.Logf("no sections found")
	// }

	// t.Fatal()
}
