package postgres_test

import (
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db/postgres"
	testcommon "github.com/CelanMatjaz/job_application_tracker_api/pkg/test_common"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func TestCreateAccount(t *testing.T) {
	conn := testcommon.CreateStore(t)
	store := postgres.CreatePostgresStore(conn.Db)

	account, err := types.CreateNewAccountData("Display name", "test@test.com", "password")
	testcommon.AssertNotError(t, err, "could not create new account data")

	_, err = store.CreateAccount(account)
	testcommon.AssertNotError(t, err, "could not create new account")
}

func TestCreateAccountWithOAuth(t *testing.T) {
	conn := testcommon.CreateStore(t)
	store := postgres.CreatePostgresStore(conn.Db)
	oauthClient := testcommon.CreateOAuthClient(t, store)

	account, err := types.CreateNewAccountData("Display name", "test@test.com", "password")
	testcommon.AssertNotError(t, err, "could not create new account data")

	_, err = store.CreateAccountWithOAuth(account, types.TokenResponse{
		AccessToken:  "",
		ExpiresIn:    0,
		RefreshToken: "",
		Scope:        "",
	}, oauthClient.Id)
	testcommon.AssertNotError(t, err, "could not create new oauth account")
}

func TestUpdateAccountToOAuth(t *testing.T) {
	conn := testcommon.CreateStore(t)
	store := postgres.CreatePostgresStore(conn.Db)
	account, _ := testcommon.SeedAccount(t, store)
	oauthClient := testcommon.CreateOAuthClient(t, store)

	err := store.UpdateAccountToOAuth(account, types.TokenResponse{
		AccessToken:  "",
		ExpiresIn:    0,
		RefreshToken: "",
		Scope:        "",
	}, oauthClient.Id)
	testcommon.AssertNotError(t, err, "could not convert normal account to oauth account")
}

func TestGetOAuthClientByName(t *testing.T) {
	conn := testcommon.CreateStore(t)
	store := postgres.CreatePostgresStore(conn.Db)
	oauthClient := testcommon.CreateOAuthClient(t, store)

	_, exists, err := store.GetOAuthClientByName(oauthClient.Name)
	testcommon.AssertNotError(t, err, "error getting oauth client by name")
	testcommon.Assert(t, exists, "oauth client does not exist")

	_, exists, err = store.GetOAuthClientByName("nonexisting name")
	testcommon.AssertNotError(t, err, "error getting oauth client by name")
	testcommon.Assert(t, !exists, "oauth client should not exist")
}
