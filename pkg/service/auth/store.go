package auth

import (
	"database/sql"
	"errors"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(connection *db.DbConnection) *Store {
	return &Store{db: connection.DB}
}

func (s *Store) GetUserById(id int) (types.User, error) {
	return s.getUser("WHERE u.id = $1", id)
}

func (s *Store) GetUserByEmail(email string) (types.User, error) {
	return s.getUser("WHERE u.email = $1", email)
}

func (s *Store) CreateUser(user types.User) (types.User, error) {
	transaction, err := s.db.Begin()
	if err != nil {
		return types.User{}, err
	}
	defer transaction.Rollback()

	row := transaction.QueryRow(`
        WITH new_user AS (
            INSERT INTO users (display_name, email) 
            VALUES ($1, $2) 
            RETURNING id
        )
        INSERT INTO password_hashes (user_id, password_hash) 
        SELECT id, $3 FROM new_user
        RETURNING user_id as id;`,
		user.DisplayName,
		user.Email,
		user.PasswordHash,
	)

	err = row.Scan(&user.Id)
	if err != nil {
		return types.User{}, err
	}

	if err != nil {
		return types.User{}, err
	}

	err = transaction.Commit()
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (s *Store) getUser(whereClause string, value any) (types.User, error) {
	row := s.db.QueryRow(`
        SELECT
            u.id, 
            u.display_name,
            u.email,
            ph.password_hash,
            u.refresh_token_version,
            u.created_at,
            u.updated_at
        FROM users u
        LEFT JOIN password_hashes ph ON u.id = ph.user_id `+whereClause,
		value)

	user, err := scanUserRow(row)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return types.User{}, types.UserDoesNotExistErr
		default:
			return types.User{}, err
		}
	}

	return user, nil
}

func (s *Store) CreateUserWithOAuth(user types.User, tokenResponse types.TokenResponse, clientId int) error {
	var id int
	err := s.db.QueryRow(`
        INSERT INTO users (display_name, email) 
        VALUES ($1, $2)
        RETURNING id`,
		user.DisplayName,
		user.Email,
	).Scan(&id)
	if err != nil {
		return err
	}

	return s.createAccountOAuthData(id, clientId, tokenResponse)
}

func (s *Store) createAccountOAuthData(userId int, oautClientId int, tokenResponse types.TokenResponse) error {
	_, err := s.db.Exec(
		"INSERT INTO account_oauth_data (user_id, oauth_client_id, access_token, refresh_token) VALUES ($1, $2, $3, $4)",
		userId, oautClientId, tokenResponse.AccessToken, tokenResponse.RefreshToken,
	)
	if err != nil {
		return err
	}

	return nil

}

func (s *Store) UpdateUserToOAuth(user types.User, tokenResponse types.TokenResponse, clientId int) error {
	_, err := s.db.Exec("UPDATE users SET is_oauth = true WHERE id = $1", user.Id)
	if err != nil {
		return err
	}

	return s.createAccountOAuthData(user.Id, clientId, tokenResponse)
}

func (s *Store) GetOauthClientByName(clientId string) (types.OAuthClient, error) {
	row := s.db.QueryRow("SELECT * FROM oauth_clients WHERE name = $1", clientId)
	var client types.OAuthClient
	err := row.Scan(
		&client.Id,
		&client.Name,
		&client.ClientId,
		&client.ClientSecret,
		&client.Scopes,
		&client.CodeEndpoint,
		&client.TokenEndpoint,
		&client.UserDataEndpoint,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return types.OAuthClient{}, types.RecordDoesNotExist
		default:
			return types.OAuthClient{}, err
		}
	}

	return client, nil
}

func scanUserRow(row *sql.Row) (types.User, error) {
	var user types.User
	err := row.Scan(
		&user.Id,
		&user.DisplayName,
		&user.Email,
		&user.PasswordHash,
		&user.TokenVersion,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return types.User{}, types.UserDoesNotExistErr
	}
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}
