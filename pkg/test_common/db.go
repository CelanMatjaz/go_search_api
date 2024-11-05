package testcommon

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db/postgres"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/golang-migrate/migrate/v4"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func CreateStore(t *testing.T) *postgres.PostgresStore {
	store := postgres.NewPostgresStore(os.Getenv("CONNECTION_STRING_CONTAINER"))
	t.Cleanup(func() {
		CleanupStore(store)
	})

	// https://stackoverflow.com/a/2829485
	_, err := store.Db.Exec(`
        CREATE OR REPLACE FUNCTION truncate_tables(username IN VARCHAR) RETURNS void AS $$
        DECLARE
            statements CURSOR FOR
                SELECT tablename FROM pg_tables
                WHERE tableowner = username AND schemaname = 'public';
        BEGIN
            FOR stmt IN statements LOOP
                EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.tablename) || ' CASCADE;';
            END LOOP;
        END;
        $$ LANGUAGE plpgsql;`)
	if err != nil {
		panic(err)
	}

	if val := recover(); val != nil {
		t.Fatalf("%v", val)
	}

	return store
}

func CleanupStore(s *postgres.PostgresStore) {
	s.Db.Exec("SELECT truncate_tables('postgres')")
	if err := s.Db.Close(); err != nil {
		panic(err)
	}
}

func MigrateUp(connectionString string) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err.Error())
	}

	driver, err := pg.WithInstance(db, &pg.Config{})
	if err != nil {
		panic(err.Error())
	}

	migrationsDir, err := getMigrationsDir()
	if err != nil {
		panic(err.Error())
	}

	migration, err := migrate.NewWithDatabaseInstance(fmt.Sprint("file:///", migrationsDir), "postgres", driver)
	if err != nil {
		panic(err.Error())
	}

	if err = migration.Up(); err != nil {
		panic(err.Error())
	}
}

func SeedAccount(t *testing.T, store *postgres.PostgresStore) (types.Account, string) {
	password := "Password1!"
	accountData, _ := types.CreateNewAccountData("Display name", "test@test.com", password)
	account, err := store.CreateAccount(accountData)
	AssertError(t, err, "could not create account for testing")
	return account, password
}

func getMigrationsDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		potentialPath := filepath.Join(cwd, "migrations")
		info, err := os.Stat(potentialPath)
		if err == nil && info.IsDir() {
			return potentialPath, nil
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			break
		}
		cwd = parent
	}

	return "", os.ErrNotExist
}

func CreateOAuthClient(t *testing.T, store *postgres.PostgresStore) types.OAuthClient {
	query := postgres.CreateRecordQuery[types.OAuthClient]("oauth_clients")
	row := store.Db.QueryRow(query,
		os.Getenv("OAUTH_CLIENT_NAME"),
		os.Getenv("OAUTH_CLIENT_ID"),
		os.Getenv("OAUTH_CLIENT_SECRET"),
		os.Getenv("OAUTH_CLIENT_SCOPES"),
		os.Getenv("OAUTH_CLIENT_CODE_ENDPOINT"),
		os.Getenv("OAUTH_CLIENT_TOKEN_ENDPOINT"),
		os.Getenv("OAUTH_CLIENT_DATA_ENDPOINT"),
	)

	var client types.OAuthClient
	err := row.Scan(postgres.GetScanFields(&client)...)
	AssertError(t, err, "could not create oauth client for testing")

	return client
}
