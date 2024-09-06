package auth

import (
	"database/sql"
	"errors"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

var UserDoesNotExistErr = errors.New("User does not exist")

type Store struct {
	db *sql.DB
}

func NewStore(connection *db.DbConnection) *Store {
	return &Store{db: connection.DB}
}

func (s *Store) GetInternalUserById(id int) (types.InternalUser, error) {
	return s.getInternalUser("WHERE users.id = $1", id)
}

func (s *Store) GetInternalUserByEmail(email string) (types.InternalUser, error) {
	return s.getInternalUser("WHERE users.email = $1", email)
}

func (s *Store) CreateUser(user types.InternalUser) (types.InternalUser, error) {
	transaction, err := s.db.Begin()
	if err != nil {
		return types.InternalUser{}, err
	}
	defer transaction.Rollback()

	row := transaction.QueryRow(`
        INSERT INTO users (first_name, last_name, email, password_hash)
        VALUES ($1, $2, $3, $4)
        RETURNING id, first_name, last_name, email, password_hash, created_at, updated_at`,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
	)

	newUser, err := scanUserRow(row)
	if err != nil {
		return types.InternalUser{}, err
	}

	err = transaction.Commit()
	if err != nil {
		return types.InternalUser{}, err
	}

	return newUser, nil
}

func (s *Store) getInternalUser(whereClause string, value any) (types.InternalUser, error) {
	row := s.db.QueryRow(`
        SELECT
            id,
            first_name,
            last_name,
            email,
            password_hash,
            created_at,
            updated_at
        FROM users `+whereClause,
		value)

	user, err := scanUserRow(row)
	if err != nil {
		return types.InternalUser{}, err
	}

	return user, nil
}

func scanUserRow(row *sql.Row) (types.InternalUser, error) {
	var user types.InternalUser
	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
    

	if errors.Is(err, sql.ErrNoRows) {
		return types.InternalUser{}, UserDoesNotExistErr
	}
	if err != nil {
		return types.InternalUser{}, err
	}

	return user, nil
}
