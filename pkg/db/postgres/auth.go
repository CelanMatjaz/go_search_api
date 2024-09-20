package postgres

import (
	"database/sql"
	"errors"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func (s *PostgresStore) GetAccountById(id int) (*types.Account, error) {
	return s.getAccount("WHERE accounts.id = $1", id)
}

func (s *PostgresStore) GetAccountByEmail(email string) (*types.Account, error) {
	return s.getAccount("WHERE accounts.email = $1", email)
}

func (s *PostgresStore) CreateAccount(account types.Account) (*types.Account, error) {
	transaction, err := s.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()

	row := transaction.QueryRow(`
        WITH new_account AS (
            INSERT INTO accounts (display_name, email) 
            VALUES ($1, $2) 
            RETURNING id
        )
        INSERT INTO password_hashes (account_id, password_hash) 
        SELECT id, $3 FROM new_account
        RETURNING account_id as id;`,
		account.DisplayName,
		account.Email,
		account.PasswordHash,
	)

	newAccount := &types.Account{}
	err = row.Scan(&account.Id)
	if err != nil {
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	return newAccount, nil
}

func (s *PostgresStore) CreateAccountWithOAuth(account types.Account, tokenResponse types.TokenResponse, clientId int) (*types.Account, error) {
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
	acc, err := scanAccountRow(row)
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

func (s *PostgresStore) GetOAuthClientByName(name string) (*types.OAuthClient, error) {
	row := s.Db.QueryRow("SELECT * FROM oauth_clients WHERE name = $1", name)
	client := &types.OAuthClient{}
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
		return nil, nil
	}

	return client, err
}

func (s *PostgresStore) getAccount(whereClause string, value any) (*types.Account, error) {
	row := s.Db.QueryRow(`
        SELECT
            accounts.id, 
            accounts.display_name,
            accounts.email,
            ph.password_hash,
            accounts.refresh_token_version,
            accounts.created_at,
            accounts.updated_at
        FROM accounts
        LEFT JOIN password_hashes ph ON accounts.id = ph.account_id `+whereClause,
		value)

	return scanAccountRow(row)
}

func scanAccountRow(row db.Scannable) (*types.Account, error) {
	account := &types.Account{}
	err := row.Scan(
		&account.Id,
		&account.DisplayName,
		&account.Email,
		&account.PasswordHash,
		&account.TokenVersion,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, err
		}
	}

	return account, nil
}
