package postgres

import (
	"database/sql"
	"errors"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

var accountScanFunc = createScanFunc[types.Account]()

func (s *PostgresStore) GetAccountById(id int) (types.Account, bool, error) {
	return s.getAccount("WHERE accounts.id = $1", id)
}

func (s *PostgresStore) GetAccountByEmail(email string) (types.Account, bool, error) {
	return s.getAccount("WHERE accounts.email = $1", email)
}

func (s *PostgresStore) CreateAccount(account types.Account) (types.Account, error) {
	var acc types.Account
	transaction, err := s.Db.Begin()
	if err != nil {
		return acc, err
	}
	defer transaction.Rollback()

	row := transaction.QueryRow(`
        WITH new_account AS (
            INSERT INTO accounts (display_name, email) 
            VALUES ($1, $2) 
            RETURNING *
        ),
        _ AS (
            INSERT INTO password_hashes (account_id, password_hash) 
            VALUES ((SELECT id FROM new_account), $3)
            RETURNING *
        ) 
        SELECT new_account.*, ph.password_hash FROM new_account
        LEFT JOIN _ ph
        ON ph.account_id = new_account.id`,
		account.DisplayName,
		account.Email,
		account.PasswordHash,
	)

	newAccount, err := accountScanFunc(row)
	if err != nil {
		return acc, err
	}

	err = transaction.Commit()
	if err != nil {
		return acc, err
	}

	return newAccount, nil
}

func (s *PostgresStore) CreateAccountWithOAuth(account types.Account, tokenResponse types.TokenResponse, clientId int) (types.Account, error) {
	row := s.Db.QueryRow(`
        INSERT INTO accounts (display_name, email) 
        VALUES ($1, $2)
        RETURNING
            accounts.id, 
            accounts.display_name,
            accounts.email,
            NULL password_hash,
            accounts.refresh_token_version,
            accounts.created_at,
            accounts.updated_at`,
		account.DisplayName,
		account.Email,
	)
	acc, err := accountScanFunc(row)
	if err != nil {
		return acc, err
	}

	err = s.createAccountOAuthData(acc.Id, clientId, tokenResponse)
	return acc, err
}

func (s *PostgresStore) createAccountOAuthData(accountId int, oautClientId int, tokenResponse types.TokenResponse) error {
	_, err := s.Db.Exec(
		"INSERT INTO account_oauth_data (account_id, oauth_client_id, access_token, refresh_token) VALUES ($1, $2, $3, $4)",
		accountId, oautClientId, tokenResponse.AccessToken, tokenResponse.RefreshToken,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccountToOAuth(account types.Account, tokenResponse types.TokenResponse, clientId int) error {
	transaction, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	_, err = transaction.Exec("UPDATE accounts SET is_oauth = true WHERE id = $1", account.Id)
	if err != nil {
		return err
	}

	_, err = transaction.Exec("DELETE FROM password_hashes WHERE account_id = $1", account.Id)
	if err != nil {
		return err
	}

	transaction.Commit()

	return s.createAccountOAuthData(account.Id, clientId, tokenResponse)
}

func (s *PostgresStore) GetOAuthClientByName(name string) (types.OAuthClient, bool, error) {
	row := s.Db.QueryRow("SELECT * FROM oauth_clients WHERE name = $1", name)
	client := types.OAuthClient{}
	err := row.Scan(
		&client.Id,
		&client.Name,
		&client.ClientId,
		&client.ClientSecret,
		&client.Scopes,
		&client.CodeEndpoint,
		&client.TokenEndpoint,
		&client.AccountDataEndpoint,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return client, false, nil
	}

	return client, true, err
}

func (s *PostgresStore) getAccount(whereClause string, value any) (types.Account, bool, error) {
	var account types.Account
	rows, err := s.Db.Query(`
        SELECT
            accounts.id, 
            accounts.display_name,
            accounts.email,
            ph.password_hash,
            accounts.refresh_token_version,
            accounts.created_at,
            accounts.updated_at
        FROM accounts
        LEFT JOIN password_hashes ph ON accounts.id = ph.account_id `+whereClause, value)
	if err != nil {
		return account, false, err
	}

	if rows.Next() {
		account, err := accountScanFunc(rows)
		return account, true, err
	} else {
		return account, false, err
	}
}
