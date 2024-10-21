package postgres_test

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db/postgres"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func createDbAndStore() (*embeddedpostgres.EmbeddedPostgres, *postgres.PostgresStore) {
	database := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Port(54321),
	)
	err := database.Start()
	if err != nil {
		panic(fmt.Sprintf("Could not create database for testing, %s", err.Error()))
	}

	connectionString := "host=localhost port=54321 user=postgres password=postgres dbname=postgres sslmode=disable"
	migrateUp(connectionString)
	return database, postgres.NewPostgresStore(connectionString)
}

func cleanupDb(db *embeddedpostgres.EmbeddedPostgres) {
	if err := db.Stop(); err != nil {
		panic(err.Error())
	}
}

func migrateUp(connectionString string) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err.Error())
	}

	driver, err := pg.WithInstance(db, &pg.Config{})
	if err != nil {
		panic(err.Error())
	}

	path, _ := filepath.Abs("../../../migrations")
	migration, err := migrate.NewWithDatabaseInstance(fmt.Sprint("file://", path), "postgres", driver)
	if err != nil {
		panic(err.Error())
	}

	if err = migration.Up(); err != nil {
		panic(err.Error())
	}
}

func seedAccount(t *testing.T, store *postgres.PostgresStore) types.Account {
	accountData, _ := types.CreateNewAccountData("Display name", "test@test.com", "password")
	account, err := store.CreateAccount(accountData)
	if err != nil {
		t.Fatalf("could not create account for testing, %s", err.Error())
	}
	return account
}
