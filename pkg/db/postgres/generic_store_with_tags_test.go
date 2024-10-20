package postgres_test

import (
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db/postgres"
)

func TestGenericStore(t *testing.T) {
	db, conn := createDbAndStore()
	defer cleanupDb(db)

	store := postgres.CreatePostgresStore(conn.Db)

	account := seedAccount(store)

	t.Fatalf("%v\n", account)
}
